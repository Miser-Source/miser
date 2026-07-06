package miser

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func anthropicMockUpstream(t *testing.T, gotHeaders *http.Header) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if gotHeaders != nil {
			*gotHeaders = r.Header.Clone()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "msg_mock",
			"type":  "message",
			"role":  "assistant",
			"model": "claude-haiku-4-5",
			"content": []map[string]interface{}{
				{"type": "text", "text": "Ticket triaged."},
			},
			"usage": map[string]interface{}{
				"input_tokens":             900,
				"output_tokens":            100,
				"cache_read_input_tokens":  100,
				"cache_write_input_tokens": 0,
			},
		})
	}))
}

func newTestProxy(t *testing.T, opts ProxyOptions) (*ProxyServer, *httptest.Server) {
	t.Helper()
	dir := t.TempDir()
	if opts.LogPath == "" {
		opts.LogPath = filepath.Join(dir, "proxy-logs.jsonl")
	}
	if opts.CachePath == "" {
		opts.CachePath = filepath.Join(dir, "exact-cache.json")
	}
	server, err := NewProxyServer(opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })
	proxy := httptest.NewServer(server.Handler())
	t.Cleanup(proxy.Close)
	return server, proxy
}

func TestProxyAnthropicMessagesCacheAndPricing(t *testing.T) {
	var upstreamHeaders http.Header
	upstream := anthropicMockUpstream(t, &upstreamHeaders)
	defer upstream.Close()

	server, proxy := newTestProxy(t, ProxyOptions{
		Provider: "anthropic",
		Upstream: upstream.URL,
		APIKey:   "sk-ant-test",
	})

	body := `{"model":"claude-haiku-4-5","max_tokens":1024,"messages":[{"role":"user","content":"Triage this ticket"}]}`

	resp, err := http.Post(proxy.URL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if got := resp.Header.Get("X-Miser-Cache"); got != "MISS" {
		t.Fatalf("first /v1/messages call: expected cache MISS, got %q", got)
	}
	if got := upstreamHeaders.Get("x-api-key"); got != "sk-ant-test" {
		t.Fatalf("expected upstream x-api-key to be injected, got %q", got)
	}
	if got := upstreamHeaders.Get("anthropic-version"); got == "" {
		t.Fatal("expected anthropic-version header to be injected")
	}

	resp, err = http.Post(proxy.URL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if got := resp.Header.Get("X-Miser-Cache"); got != "HIT" {
		t.Fatalf("repeated /v1/messages call: expected cache HIT, got %q", got)
	}
	if !strings.Contains(string(payload), "Ticket triaged.") {
		t.Fatalf("cached response body mismatch: %s", payload)
	}

	rows, err := loadProxyLogRows(server.opts.LogPath, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 log rows, got %d", len(rows))
	}
	hit, miss := rows[0], rows[1]
	if miss["cache_status"] != "miss" || hit["cache_status"] != "hit" {
		t.Fatalf("unexpected cache statuses: %v then %v", miss["cache_status"], hit["cache_status"])
	}
	if miss["provider"] != "anthropic" || miss["workflow"] != "proxy_messages" {
		t.Fatalf("unexpected provider/workflow: %v/%v", miss["provider"], miss["workflow"])
	}
	// haiku-4.5 at $1/M in, $5/M out with 100 cached reads at 10%:
	// 900*1.00/1M + 100*0.10/1M + 100*5.00/1M = 0.00141
	cost, _ := miss["cost_usd"].(float64)
	if cost < 0.00140 || cost > 0.00142 {
		t.Fatalf("expected priced anthropic cost ~0.00141, got %v (basis %v)", cost, miss["cost_basis"])
	}
	if miss["cost_basis"] != "published_token_price" {
		t.Fatalf("expected published_token_price basis, got %v", miss["cost_basis"])
	}
	saved, _ := hit["cache_saved_usd"].(float64)
	if saved <= 0 {
		t.Fatalf("expected cache hit to record savings, got %v", saved)
	}
	// Anthropic input_tokens excludes cache reads; Miser normalizes to include them.
	if got := intFromAny(miss["input_tokens"]); got != 1000 {
		t.Fatalf("expected normalized input_tokens 1000, got %d", got)
	}
}

func TestProxyRuntimeProviderSwitch(t *testing.T) {
	server, proxy := newTestProxy(t, ProxyOptions{Provider: "openai", APIKey: "sk-openai"})

	if got := server.currentUpstream(); got.Host != "api.openai.com" {
		t.Fatalf("expected default openai upstream, got %s", got.Host)
	}

	resp, err := http.Post(proxy.URL+"/miser/api/key", "application/json",
		strings.NewReader(`{"key":"sk-ant-new","provider":"anthropic"}`))
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&cfg)
	_ = resp.Body.Close()
	if cfg["provider"] != "anthropic" {
		t.Fatalf("expected provider anthropic after switch, got %v", cfg["provider"])
	}
	if cfg["key_env"] != "ANTHROPIC_API_KEY" {
		t.Fatalf("expected ANTHROPIC_API_KEY key env, got %v", cfg["key_env"])
	}
	if got := server.currentUpstream(); got.Host != "api.anthropic.com" {
		t.Fatalf("expected upstream to follow provider, got %s", got.Host)
	}

	resp, err = http.Post(proxy.URL+"/miser/api/key", "application/json",
		strings.NewReader(`{"key":"sk-x","provider":"bogus"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown provider, got %d", resp.StatusCode)
	}
}

func TestProxyCustomUpstreamSurvivesProviderSwitch(t *testing.T) {
	upstream := anthropicMockUpstream(t, nil)
	defer upstream.Close()

	server, proxy := newTestProxy(t, ProxyOptions{
		Provider: "openai",
		Upstream: upstream.URL,
		APIKey:   "sk-test",
	})
	resp, err := http.Post(proxy.URL+"/miser/api/key", "application/json",
		strings.NewReader(`{"key":"sk-ant","provider":"anthropic"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if got := server.currentUpstream(); !strings.Contains(upstream.URL, got.Host) {
		t.Fatalf("custom upstream should survive provider switch, got %s", got.Host)
	}
}

func TestProxyBusinessProfileAndBudget(t *testing.T) {
	upstream := anthropicMockUpstream(t, nil)
	defer upstream.Close()

	_, proxy := newTestProxy(t, ProxyOptions{
		Provider:  "anthropic",
		Upstream:  upstream.URL,
		APIKey:    "sk-ant",
		Mode:      "business",
		Workspace: "Acme Inc",
		BudgetUSD: 0.001, // tiny budget so a single priced call exceeds it
	})

	resp, err := http.Get(proxy.URL + "/miser/api/config")
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&cfg)
	_ = resp.Body.Close()
	if cfg["mode"] != "business" || cfg["workspace"] != "Acme Inc" {
		t.Fatalf("expected business profile in config, got %v/%v", cfg["mode"], cfg["workspace"])
	}
	if budget, _ := cfg["budget_usd"].(float64); budget != 0.001 {
		t.Fatalf("expected budget 0.001, got %v", cfg["budget_usd"])
	}

	// First call accrues spend past the tiny budget.
	body := `{"model":"claude-haiku-4-5","max_tokens":64,"messages":[{"role":"user","content":"hello"}]}`
	resp, err = http.Post(proxy.URL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if got := resp.Header.Get("X-Miser-Budget-Status"); got != "ok" {
		t.Fatalf("first call should be within budget, got status %q", got)
	}

	resp, err = http.Post(proxy.URL+"/v1/messages", "application/json",
		strings.NewReader(strings.Replace(body, "hello", "hello again", 1)))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if got := resp.Header.Get("X-Miser-Budget-Status"); got != "exceeded" {
		t.Fatalf("second call should exceed budget, got status %q", got)
	}

	// Profile updates via /miser/api/profile.
	resp, err = http.Post(proxy.URL+"/miser/api/profile", "application/json",
		strings.NewReader(`{"workspace":"Acme Labs","budget_usd":100}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = json.NewDecoder(resp.Body).Decode(&cfg)
	_ = resp.Body.Close()
	if cfg["workspace"] != "Acme Labs" {
		t.Fatalf("expected workspace update, got %v", cfg["workspace"])
	}
	if exceeded, _ := cfg["budget_exceeded"].(bool); exceeded {
		t.Fatal("raising the budget should clear exceeded state")
	}
	if spend, _ := cfg["spend_month_usd"].(float64); spend <= 0 {
		t.Fatalf("expected month spend > 0, got %v", spend)
	}
}

func TestMonthSpendSeedsFromLog(t *testing.T) {
	upstream := anthropicMockUpstream(t, nil)
	defer upstream.Close()

	dir := t.TempDir()
	logPath := filepath.Join(dir, "logs.jsonl")

	server, proxy := newTestProxy(t, ProxyOptions{
		Provider: "anthropic", Upstream: upstream.URL, APIKey: "k",
		LogPath: logPath, CachePath: filepath.Join(dir, "c1.json"),
	})
	body := `{"model":"claude-haiku-4-5","max_tokens":64,"messages":[{"role":"user","content":"seed"}]}`
	resp, err := http.Post(proxy.URL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_, firstSpend, _ := server.budgetStatus()
	if firstSpend <= 0 {
		t.Fatalf("expected spend accrued, got %v", firstSpend)
	}

	if _, err := os.Stat(logPath); err != nil {
		t.Fatal(err)
	}
	restarted, err := NewProxyServer(ProxyOptions{
		Provider: "anthropic", Upstream: upstream.URL, APIKey: "k",
		LogPath: logPath, CachePath: filepath.Join(dir, "c2.json"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer restarted.Close()
	_, seeded, _ := restarted.budgetStatus()
	if seeded != firstSpend {
		t.Fatalf("expected restart to seed month spend %v from log, got %v", firstSpend, seeded)
	}
}
