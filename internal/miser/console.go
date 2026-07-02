package miser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
)

type ConsoleConfig struct {
	Provider    string
	AccountID   string
	Integration string
	LogPath     string
	CachePath   string
}

func RenderConsoleHTML(config ConsoleConfig) string {
	provider := defaultString(config.Provider, "openai")
	account := defaultString(config.AccountID, "—")
	integration := defaultString(config.Integration, "—")

	var b strings.Builder
	fmt.Fprintln(&b, "<!doctype html>")
	fmt.Fprintln(&b, `<html lang="en" data-theme="dark">`)
	fmt.Fprintln(&b, `<head>`)
	fmt.Fprintln(&b, `<meta charset="utf-8">`)
	fmt.Fprintln(&b, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintln(&b, `<title>Miser — AI spend control plane</title>`)
	fmt.Fprintln(&b, `<style>`)
	fmt.Fprintln(&b, consoleCSS())
	fmt.Fprintln(&b, `</style>`)
	fmt.Fprintln(&b, `</head>`)
	fmt.Fprintln(&b, `<body>`)

	fmt.Fprintln(&b, `<div class="app">`)

	// ---------- Sidebar ----------
	fmt.Fprintln(&b, `<aside class="sidebar" id="sidebar">`)
	fmt.Fprintln(&b, `<div class="side-resize" id="sideResize"></div>`)
	fmt.Fprintln(&b, `<div class="side-top">`)
	fmt.Fprintln(&b, `<div class="brand">`+miserMark()+`<span class="brand-name">Miser</span><span class="brand-tag">control plane</span></div>`)
	fmt.Fprintln(&b, `<button class="new-chat" id="newChat"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M5 12h14"/></svg>New playground</button>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<nav class="side-nav">`)
	fmt.Fprintln(&b, `<a class="active" href="/"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>Playground</a>`)
	fmt.Fprintln(&b, `<a href="/miser/api/requests" target="_blank"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="M7 14l4-4 3 3 5-6"/></svg>Requests feed</a>`)
	fmt.Fprintln(&b, `<a href="/healthz" target="_blank"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>Health</a>`)
	fmt.Fprintln(&b, `</nav>`)

	fmt.Fprintln(&b, `<div class="side-history">`)
	fmt.Fprintln(&b, `<p class="side-label">Recent traffic</p>`)
	fmt.Fprintln(&b, `<div id="history" class="history"><p class="history-empty">No intercepted requests yet.</p></div>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<div class="side-foot">`)
	fmt.Fprintln(&b, `<div class="foot-row"><div class="status"><span class="dot"></span>Proxy live</div><button class="icon-btn sm" id="settingsBtn" title="Settings"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg></button></div>`)
	fmt.Fprintf(&b, `<dl class="meta"><div><dt>Provider</dt><dd>%s</dd></div><div><dt>Account</dt><dd>%s</dd></div><div><dt>Integration</dt><dd>%s</dd></div><div><dt>Key</dt><dd id="keyStatus" class="key-status" title="Connect or change provider key">not set</dd></div></dl>`,
		html.EscapeString(provider), html.EscapeString(account), html.EscapeString(integration))

	// settings popover
	fmt.Fprintln(&b, `<div class="settings-menu" id="settingsMenu" hidden>`)
	fmt.Fprintf(&b, `<p class="sm-account">%s</p>`, html.EscapeString(account))
	fmt.Fprintln(&b, `<div class="sm-divider"></div>`)
	fmt.Fprintln(&b, `<p class="sm-label">Theme</p>`)
	fmt.Fprintln(&b, `<div class="sm-themes">`)
	themeSwatch(&b, "midnight", "Midnight", "#161616")
	themeSwatch(&b, "white", "White", "#ffffff")
	themeSwatch(&b, "slate", "Slate", "#c2ccd6")
	themeSwatch(&b, "red", "Red", "#f5436b")
	themeSwatch(&b, "rose", "Rose", "#f5487f")
	themeSwatch(&b, "pink", "Pink", "#f546b3")
	themeSwatch(&b, "magenta", "Magenta", "#f54df5")
	themeSwatch(&b, "purple", "Purple", "#9a4df5")
	themeSwatch(&b, "violet", "Violet", "#7d5cf6")
	themeSwatch(&b, "indigo", "Indigo", "#5a6cf6")
	themeSwatch(&b, "blue", "Blue", "#36a3f7")
	themeSwatch(&b, "sky", "Sky", "#3ac6f5")
	themeSwatch(&b, "cyan", "Cyan", "#36d6e6")
	themeSwatch(&b, "teal", "Teal", "#33d6d6")
	themeSwatch(&b, "green", "Green", "#2fe0a0")
	themeSwatch(&b, "lime", "Lime", "#82dd54")
	themeSwatch(&b, "yellow", "Yellow", "#f5d23a")
	themeSwatch(&b, "amber", "Amber", "#f5b13a")
	themeSwatch(&b, "orange", "Orange", "#f5803a")
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `<div class="sm-divider"></div>`)
	fmt.Fprintln(&b, `<button class="sm-item danger" id="disconnectBtn"><svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><path d="m16 17 5-5-5-5"/><path d="M21 12H9"/></svg>Disconnect provider</button>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</aside>`)

	// ---------- Main ----------
	fmt.Fprintln(&b, `<main class="main">`)
	fmt.Fprintln(&b, `<header class="topbar">`)
	fmt.Fprintln(&b, `<button class="icon-btn" id="sidebarToggle" title="Toggle sidebar"><svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18"/></svg></button>`)
	fmt.Fprintln(&b, `<span class="topbar-spacer"></span>`)
	fmt.Fprintln(&b, `<button class="icon-btn" id="inspectorToggle" title="Toggle Miser inspector"><svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="7"/><path d="m21 21-4.3-4.3"/></svg></button>`)
	fmt.Fprintln(&b, `</header>`)

	// thread
	fmt.Fprintln(&b, `<div class="thread" id="thread">`)
	fmt.Fprintln(&b, `<div class="home" id="welcome">`)
	fmt.Fprintln(&b, `<div class="home-head">`)
	fmt.Fprintf(&b, `<h1><span class="home-mark">%s</span>Spend less on every call.</h1>`, miserMark())
	fmt.Fprintln(&b, `<span class="home-aside" id="homeSaved">$0.00 saved so far</span>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<section class="home-section">`)
	fmt.Fprintln(&b, `<p class="home-label">Quick start</p>`)
	fmt.Fprintln(&b, `<div class="home-list">`)
	homeRow(&b, "accent", "Summarize a support ticket", "Two-sentence summary — watch the cache write and token cost", "send",
		"Summarize this support ticket in two sentences: customer says the export button on the billing page returns a 500 error after the latest release, and they need a fix before month-end close.")
	homeRow(&b, "accent", "Classify an email", "Cheap-model classification — see the decision trace", "send",
		"Classify this email as one of: billing, bug, feature_request, spam. Reply with only the label.\n\n\"Hey, the dashboard keeps logging me out every few minutes since yesterday.\"")
	homeRow(&b, "accent", "Repeat to test the exact cache", "Send it twice and watch the cache HIT — zero new spend", "send",
		"Reply with exactly: Miser exact cache works.")
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</section>`)

	fmt.Fprintln(&b, `<section class="home-section" id="trafficSection" style="display:none">`)
	fmt.Fprintln(&b, `<p class="home-label">Recent traffic</p>`)
	fmt.Fprintln(&b, `<div class="home-list" id="homeTraffic"></div>`)
	fmt.Fprintln(&b, `</section>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</div>`)

	// composer
	fmt.Fprintln(&b, `<div class="composer-wrap">`)
	fmt.Fprintln(&b, `<form id="chatForm" class="composer">`)
	fmt.Fprintln(&b, `<textarea id="prompt" rows="1" placeholder="Message the Miser proxy…  (Enter to send, Shift+Enter for newline)"></textarea>`)
	fmt.Fprintln(&b, `<button type="submit" id="sendBtn" class="send" disabled aria-label="Send"><svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 19V5M5 12l7-7 7 7"/></svg></button>`)
	fmt.Fprintln(&b, `</form>`)
	fmt.Fprintln(&b, `<div class="composer-foot">`)
	dm := defaultModelFor(provider)
	fmt.Fprintln(&b, `<div class="model-pick" id="modelPick">`)
	fmt.Fprintf(&b, `<input type="hidden" id="model" value="%s">`, html.EscapeString(dm))
	fmt.Fprintf(&b, `<button type="button" class="model-btn" id="modelBtn"><svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="3"/><path d="M9 9h6v6H9z"/></svg><span id="modelLabel">%s</span><svg class="caret" viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg></button>`, html.EscapeString(dm))
	fmt.Fprintf(&b, `<div class="model-menu" id="modelMenu" hidden>%s</div>`, modelMenuHTML(provider, dm))
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `<span class="composer-hint" id="cacheHint"></span>`)
	fmt.Fprintln(&b, `<span class="composer-legal">Prompts not stored unless <code>--store-prompts</code></span>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</main>`)

	// ---------- Inspector ----------
	fmt.Fprintln(&b, `<aside class="inspector" id="inspector">`)
	fmt.Fprintln(&b, `<div class="insp-head"><h2>Miser inspector</h2><button class="icon-btn" id="inspClose" title="Close"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg></button></div>`)

	fmt.Fprintln(&b, `<div class="insp-metrics">`)
	inspMetric(&b, "m_requests", "0", "Intercepted")
	inspMetric(&b, "m_saved", "$0.00", "Saved by cache")
	inspMetric(&b, "m_rate", "0%", "Cache hit rate")
	inspMetric(&b, "m_spend", "$0.00", "Spend after Miser")
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<div class="insp-section">`)
	fmt.Fprintln(&b, `<p class="insp-label">Decision trace</p>`)
	fmt.Fprintln(&b, `<div id="decision" class="decision"><div class="dec-action"><span class="badge ghost">Idle</span></div><p class="dec-reason">No intercepted request yet. Send a message to see what Miser does with it.</p></div>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<div class="insp-section grow">`)
	fmt.Fprintln(&b, `<div class="insp-label-row"><p class="insp-label">Request inspector</p><button class="copy-btn" id="copyJson">Copy</button></div>`)
	fmt.Fprintln(&b, `<pre id="inspectorJson" class="json">{
  "state": "waiting_for_request"
}</pre>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</aside>`)

	fmt.Fprintln(&b, `<div class="scrim" id="scrim"></div>`)
	fmt.Fprintln(&b, `</div>`)

	// ---------- Setup gate (terminal) ----------
	fmt.Fprintln(&b, `<div class="setup" id="setup" hidden>`)
	fmt.Fprintln(&b, `<div class="term">`)
	fmt.Fprintln(&b, `<div class="term-bar"><span class="tdot r"></span><span class="tdot y"></span><span class="tdot g"></span><span class="term-title">miser — connect provider</span></div>`)
	fmt.Fprintln(&b, `<div class="term-body">`)
	fmt.Fprintln(&b, `<p class="term-line"><span class="prompt">miser$</span> connect <span id="termProvider">openai</span></p>`)
	fmt.Fprintln(&b, `<p class="term-out">⚠ No API key found in <code id="termEnv">OPENAI_API_KEY</code>.</p>`)
	fmt.Fprintln(&b, `<p class="term-out">Paste your provider API key below to route live requests through the Miser proxy.</p>`)
	fmt.Fprintln(&b, `<p class="term-out dim"># held in memory for this session only · never written to logs or disk</p>`)
	fmt.Fprintln(&b, `<form id="keyForm" class="term-form">`)
	fmt.Fprintln(&b, `<span class="prompt">key&gt;</span>`)
	fmt.Fprintln(&b, `<input id="keyInput" class="term-input" type="password" autocomplete="off" autocapitalize="off" spellcheck="false" placeholder="sk-..." />`)
	fmt.Fprintln(&b, `<button type="button" id="keyReveal" class="term-reveal" title="Show/hide">show</button>`)
	fmt.Fprintln(&b, `<button type="submit" id="keyConnect" class="term-connect">Connect →</button>`)
	fmt.Fprintln(&b, `</form>`)
	fmt.Fprintln(&b, `<p class="term-msg" id="keyMsg"></p>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</div>`)
	fmt.Fprintln(&b, `</div>`)

	fmt.Fprintln(&b, `<script>`)
	fmt.Fprintln(&b, consoleJS())
	fmt.Fprintln(&b, `</script>`)
	fmt.Fprintln(&b, `</body></html>`)
	return b.String()
}

func miserMark() string {
	return `<svg class="mark" viewBox="0 0 32 32" width="26" height="26" aria-hidden="true"><rect width="32" height="32" rx="9" fill="var(--accent)"/><path d="M9 22V10l5 6 5-6v12" fill="none" stroke="var(--on-accent)" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"/><circle cx="22.5" cy="20.5" r="2.2" fill="var(--on-accent)"/></svg>`
}

func homeRow(b *strings.Builder, dot, title, desc, meta, prompt string) {
	fmt.Fprintf(b, `<button class="row" type="button" data-prompt="%s">`, html.EscapeString(prompt))
	fmt.Fprintf(b, `<span class="row-dot %s"></span>`, html.EscapeString(dot))
	fmt.Fprintf(b, `<span class="row-main"><strong>%s</strong><span class="row-sub">%s</span></span>`, html.EscapeString(title), html.EscapeString(desc))
	fmt.Fprintf(b, `<span class="row-meta">%s%s</span>`, html.EscapeString(meta), chevron())
	fmt.Fprintln(b, `</button>`)
}

func themeSwatch(b *strings.Builder, name, label, color string) {
	fmt.Fprintf(b, `<button class="swatch" type="button" data-theme="%s" title="%s" style="background:%s"></button>`,
		html.EscapeString(name), html.EscapeString(label), html.EscapeString(color))
}

func chevron() string {
	return `<svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>`
}

func inspMetric(b *strings.Builder, id, value, label string) {
	fmt.Fprintf(b, `<div class="metric"><span id="%s" class="metric-val">%s</span><span class="metric-lbl">%s</span></div>`,
		html.EscapeString(id), html.EscapeString(value), html.EscapeString(label))
}

type modelGroup struct {
	label  string
	models []string
}

func modelGroups(provider string) []modelGroup {
	if provider == "anthropic" {
		return []modelGroup{
			{"Claude 3.7", []string{"claude-3-7-sonnet-latest"}},
			{"Claude 3.5", []string{"claude-3-5-sonnet-latest", "claude-3-5-haiku-latest"}},
			{"Claude 3", []string{"claude-3-opus-latest", "claude-3-haiku-20240307"}},
		}
	}
	// Full OpenAI API lineup (chat/completions + responses).
	return []modelGroup{
		{"GPT-5", []string{"gpt-5", "gpt-5-mini", "gpt-5-nano", "gpt-5-chat-latest"}},
		{"GPT-4.1", []string{"gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano"}},
		{"GPT-4o", []string{"gpt-4o", "gpt-4o-mini", "chatgpt-4o-latest", "gpt-4o-2024-11-20", "gpt-4o-search-preview", "gpt-4o-mini-search-preview"}},
		{"Reasoning (o-series)", []string{"o3", "o3-pro", "o3-mini", "o4-mini", "o1", "o1-pro", "o1-mini"}},
		{"GPT-4.5 / legacy", []string{"gpt-4.5-preview", "gpt-4-turbo", "gpt-4", "gpt-3.5-turbo"}},
	}
}

func defaultModelFor(provider string) string {
	if provider == "anthropic" {
		return "claude-3-5-haiku-latest"
	}
	return "gpt-4o-mini"
}

func modelMenuHTML(provider, selected string) string {
	var b strings.Builder
	for _, g := range modelGroups(provider) {
		fmt.Fprintf(&b, `<p class="model-group">%s</p>`, html.EscapeString(g.label))
		for _, m := range g.models {
			cls := "model-opt"
			if m == selected {
				cls += " selected"
			}
			fmt.Fprintf(&b, `<button type="button" class="%s" data-model="%s"><svg class="check" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg><span>%s</span></button>`,
				cls, html.EscapeString(m), html.EscapeString(m))
		}
	}
	return b.String()
}

func loadProxyLogRows(path string, limit int) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]interface{}{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var rows []map[string]interface{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 10*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var row map[string]interface{}
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			continue
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	if limit > 0 && len(rows) > limit {
		rows = rows[:limit]
	}
	return rows, nil
}

func consoleCSS() string {
	return `
:root {
  color-scheme: dark;
  --th-hue: 0; --th-sat: 0%;
  --bg:        hsl(var(--th-hue) var(--th-sat) 5%);
  --bg-2:      hsl(var(--th-hue) var(--th-sat) 3.5%);
  --sidebar:   hsl(var(--th-hue) var(--th-sat) 4.5%);
  --surface:   hsl(var(--th-hue) var(--th-sat) 10%);
  --surface-2: hsl(var(--th-hue) var(--th-sat) 13%);
  --hover:     hsl(var(--th-hue) var(--th-sat) 15.5%);
  --user-bubble: hsl(var(--th-hue) var(--th-sat) 14%);
  --side-w: 230px;
  --line: rgba(255,255,255,.07);
  --line-2: rgba(255,255,255,.13);
  --text: #ececec;
  --muted: #9a9a9a;
  --faint: #6e6e6e;
  --accent: #4d9cf6;
  --accent-soft: color-mix(in srgb, var(--accent) 16%, transparent);
  --on-accent: #07182e;
  --good: #3fcf8e;
  --warn: #e0b341;
  --elev: #181818;
  --code-bg: #0b0c0e;
  --code-fg: #c5cad3;
  --radius: 14px;
  --shadow: 0 12px 40px rgba(0,0,0,.5);
  --mono: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
  --sans: "Söhne", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}
/* Vivid full-screen color themes (light-on-color, like Linear) */
body[data-theme] {
  color-scheme: light;
  --bg:        hsl(var(--h) 88% 60%);
  --bg-2:      hsl(var(--h) 88% 55%);
  --sidebar:   hsl(var(--h) 86% 57%);
  --surface:   hsla(var(--h) 70% 18% / .10);
  --surface-2: hsla(var(--h) 70% 16% / .16);
  --hover:     hsla(var(--h) 70% 16% / .13);
  --user-bubble: hsla(var(--h) 70% 16% / .16);
  --text:      hsl(var(--h) 85% 13%);
  --muted:     hsla(var(--h) 80% 14% / .68);
  --faint:     hsla(var(--h) 80% 14% / .46);
  --line:      hsla(var(--h) 80% 14% / .15);
  --line-2:    hsla(var(--h) 80% 14% / .26);
  --on-accent: #ffffff;
  --good:      hsl(152 75% 26%);
  --warn:      hsl(34 90% 30%);
  --elev:      hsl(var(--h) 80% 52%);
  --code-bg:   hsla(var(--h) 85% 12% / .14);
  --code-fg:   hsl(var(--h) 80% 16%);
  --accent-soft: hsla(var(--h) 80% 14% / .12);
  --shadow:    0 14px 44px hsla(var(--h) 80% 18% / .4);
}
body[data-theme="red"]     { --h: 353; --accent: #6b3df5; }
body[data-theme="rose"]    { --h: 340; --accent: #6b3df5; }
body[data-theme="pink"]    { --h: 325; --accent: #6b3df5; }
body[data-theme="magenta"] { --h: 300; --accent: #ffd24d; }
body[data-theme="purple"]  { --h: 270; --accent: #ff5dba; }
body[data-theme="violet"]  { --h: 258; --accent: #ffd24d; }
body[data-theme="indigo"]  { --h: 245; --accent: #ffd24d; }
body[data-theme="blue"]    { --h: 212; --accent: #ff5d8f; }
body[data-theme="sky"]     { --h: 200; --accent: #ff4d6d; }
body[data-theme="cyan"]    { --h: 190; --accent: #ff4d6d; }
body[data-theme="teal"]    { --h: 180; --accent: #ff4d6d; }
body[data-theme="green"]   { --h: 150; --accent: #d6009b; }
body[data-theme="lime"]    { --h: 95;  --accent: #d6009b; }
body[data-theme="yellow"]  { --h: 52;  --accent: #6b3df5; }
body[data-theme="amber"]   { --h: 40;  --accent: #da004b; }
body[data-theme="orange"]  { --h: 25;  --accent: #2563eb; }
/* Light themes */
body[data-theme="white"] {
  color-scheme: light;
  --bg: #ffffff; --bg-2: #f3f4f6; --sidebar: #f7f8fa;
  --surface: rgba(0,0,0,.035); --surface-2: rgba(0,0,0,.07); --hover: rgba(0,0,0,.05);
  --user-bubble: rgba(0,0,0,.06);
  --text: #1c1d1f; --muted: rgba(0,0,0,.56); --faint: rgba(0,0,0,.4);
  --line: rgba(0,0,0,.1); --line-2: rgba(0,0,0,.16);
  --accent: #3b82f6; --on-accent: #ffffff;
  --good: hsl(152 68% 32%); --warn: hsl(34 90% 34%);
  --elev: #ffffff; --code-bg: #f3f4f6; --code-fg: #24292f;
  --accent-soft: rgba(59,130,246,.12);
  --shadow: 0 14px 44px rgba(0,0,0,.14);
}
body[data-theme="slate"] {
  color-scheme: light;
  --bg: #e9edf2; --bg-2: #dfe4ea; --sidebar: #e4e9ef;
  --surface: rgba(20,30,50,.05); --surface-2: rgba(20,30,50,.09); --hover: rgba(20,30,50,.07);
  --user-bubble: rgba(20,30,50,.08);
  --text: #1c232e; --muted: rgba(28,35,46,.6); --faint: rgba(28,35,46,.42);
  --line: rgba(20,30,50,.12); --line-2: rgba(20,30,50,.18);
  --accent: #4f46e5; --on-accent: #ffffff;
  --good: hsl(152 60% 30%); --warn: hsl(34 85% 32%);
  --elev: #eef1f5; --code-bg: #dde3ea; --code-fg: #2a3340;
  --accent-soft: rgba(79,70,229,.12);
  --shadow: 0 14px 44px rgba(30,40,60,.16);
}
* { box-sizing: border-box; }
html, body { height: 100%; }
body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font: 15px/1.55 var(--sans);
  -webkit-font-smoothing: antialiased;
}
::selection { background: var(--accent-soft); }
::-webkit-scrollbar { width: 10px; height: 10px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb {
  background: color-mix(in srgb, var(--text) 14%, transparent);
  border-radius: 999px; border: 3px solid transparent; background-clip: padding-box;
}
::-webkit-scrollbar-thumb:hover {
  background: color-mix(in srgb, var(--text) 26%, transparent);
  border: 3px solid transparent; background-clip: padding-box;
}
:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; border-radius: 6px; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: none; } }
@media (prefers-reduced-motion: reduce) { *, *::before, *::after { animation: none !important; transition: none !important; } }
.app {
  height: 100vh;
  display: grid;
  grid-template-columns: var(--side-w) minmax(0,1fr) 0;
  transition: grid-template-columns .22s ease;
}
.app.resizing { transition: none; user-select: none; }
.app.show-inspector { grid-template-columns: var(--side-w) minmax(0,1fr) 372px; }
.app.hide-sidebar { grid-template-columns: 0 minmax(0,1fr) 0; }
.app.hide-sidebar.show-inspector { grid-template-columns: 0 minmax(0,1fr) 372px; }

/* ---------- Sidebar ---------- */
.sidebar {
  position: relative;
  background: var(--sidebar);
  border-right: 1px solid var(--line);
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}
.side-resize {
  position: absolute; top: 0; right: -5px; z-index: 50;
  width: 10px; height: 100%; cursor: col-resize;
  display: flex; align-items: center; justify-content: center;
}
.side-resize::after {
  content: ""; width: 4px; height: 44px; border-radius: 999px;
  background: rgba(255,255,255,.16); transition: background .12s, height .12s;
}
.side-resize:hover::after, .app.resizing .side-resize::after { background: rgba(255,255,255,.32); height: 60px; }
.side-top { padding: 14px 12px 8px; display: grid; gap: 12px; }
.brand { display: flex; align-items: center; gap: 9px; padding: 4px 6px; }
.brand .mark { border-radius: 8px; box-shadow: 0 2px 12px rgba(77,156,246,.35); flex: none; }
.brand-name { font-weight: 650; font-size: 17px; letter-spacing: -.01em; }
.brand-tag { color: var(--faint); font-size: 11px; margin-left: -2px; align-self: flex-end; padding-bottom: 3px; }
.new-chat {
  display: flex; align-items: center; gap: 10px;
  width: 100%; padding: 9px 12px;
  background: transparent; color: var(--text);
  border: 0; border-radius: 9px;
  font: inherit; font-size: 14px; cursor: pointer;
  transition: background .12s;
}
.new-chat svg { color: var(--faint); flex: none; }
.new-chat:hover { background: var(--hover); }
.side-nav { padding: 6px 8px; display: grid; gap: 2px; }
.side-nav a {
  display: flex; align-items: center; gap: 10px;
  padding: 9px 12px; border-radius: 9px;
  color: var(--muted); text-decoration: none; font-size: 14px;
  transition: background .12s, color .12s;
}
.side-nav a svg { color: var(--faint); flex: none; }
.side-nav a:hover { background: var(--hover); color: var(--text); }
.side-nav a.active { background: var(--surface-2); color: var(--text); box-shadow: inset 2.5px 0 0 var(--accent); }
.side-nav a.active svg { color: var(--accent); }
.side-history { flex: 1; min-height: 0; display: flex; flex-direction: column; padding: 10px 8px 4px; }
.side-label { margin: 4px 8px 8px; font-size: 11px; text-transform: uppercase; letter-spacing: .07em; color: var(--faint); }
.history { overflow-y: auto; display: grid; gap: 2px; padding-right: 2px; }
.history-empty { color: var(--faint); font-size: 13px; padding: 4px 10px; }
.hist-item {
  width: 100%; text-align: left; border: 0; background: transparent;
  color: var(--muted); font: inherit; cursor: pointer;
  padding: 8px 10px; border-radius: 8px; display: grid; gap: 3px;
  transition: background .12s;
}
.hist-item:hover { background: var(--hover); color: var(--text); }
.hist-line { display: flex; align-items: center; gap: 7px; font-size: 13px; }
.hist-line .ico { flex: none; }
.hist-model { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.hist-sub { font-size: 11px; color: var(--faint); display: flex; gap: 8px; }
.side-foot { position: relative; border-top: 1px solid var(--line); padding: 12px; display: grid; gap: 10px; }
.foot-row { display: flex; align-items: center; justify-content: space-between; }
.icon-btn.sm { width: 28px; height: 28px; border-radius: 7px; }
.status { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--muted); }
.settings-menu {
  position: absolute; bottom: calc(100% - 4px); left: 12px; right: 12px; z-index: 70;
  background: color-mix(in srgb, var(--elev) 88%, transparent);
  backdrop-filter: blur(18px); -webkit-backdrop-filter: blur(18px);
  border: 1px solid var(--line-2); border-radius: 12px;
  padding: 8px; box-shadow: var(--shadow);
  animation: fadeUp .16s ease both;
}
.settings-menu[hidden] { display: none; }
.sm-account { margin: 4px 6px 8px; font-size: 12px; color: var(--muted); overflow: hidden; text-overflow: ellipsis; }
.sm-divider { height: 1px; background: var(--line); margin: 6px 2px; }
.sm-label { margin: 4px 6px 8px; font-size: 11px; color: var(--faint); }
.sm-themes { display: flex; gap: 8px; padding: 0 4px 4px; flex-wrap: wrap; }
.swatch {
  width: 24px; height: 24px; border-radius: 7px; cursor: pointer;
  border: 1px solid rgba(255,255,255,.18); box-shadow: inset 0 0 0 0 transparent;
  transition: transform .1s, box-shadow .12s;
}
.swatch:hover { transform: translateY(-1px); }
.swatch.active { box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--accent); }
.sm-item {
  display: flex; align-items: center; gap: 9px; width: 100%; text-align: left;
  border: 0; background: transparent; color: var(--text); font: inherit; font-size: 13px;
  padding: 8px 8px; border-radius: 8px; cursor: pointer;
}
.sm-item:hover { background: var(--hover); }
.sm-item.danger { color: #ff7a7a; }
.sm-item.danger svg { color: #ff7a7a; }
.dot { width: 8px; height: 8px; border-radius: 50%; background: var(--good); box-shadow: 0 0 0 0 rgba(31,209,139,.6); animation: pulse 2.4s infinite; }
@keyframes pulse { 0%{box-shadow:0 0 0 0 rgba(31,209,139,.5);} 70%{box-shadow:0 0 0 7px rgba(31,209,139,0);} 100%{box-shadow:0 0 0 0 rgba(31,209,139,0);} }
.meta { margin: 0; display: grid; gap: 5px; }
.meta div { display: flex; justify-content: space-between; gap: 10px; font-size: 12px; }
.meta dt { color: var(--faint); margin: 0; }
.meta dd { margin: 0; color: var(--muted); font-variant-numeric: tabular-nums; }
.key-status { cursor: pointer; font-family: var(--mono); font-size: 11px; }
.key-status:hover { color: var(--accent); }

/* ---------- Main ---------- */
.main { min-width: 0; min-height: 0; display: flex; flex-direction: column; background: var(--bg); position: relative; }
.topbar {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 10px;
}
.topbar-spacer { flex: 1; }
.icon-btn {
  flex: none; display: grid; place-items: center;
  width: 36px; height: 36px; border-radius: 9px;
  background: transparent; border: 1px solid transparent; color: var(--muted);
  cursor: pointer; transition: background .12s, color .12s, border-color .12s;
}
.icon-btn:hover { background: var(--hover); color: var(--text); border-color: var(--line); }

/* ---------- Thread ---------- */
.thread { flex: 1; min-height: 0; overflow-y: auto; scroll-behavior: smooth; padding: 24px 0 8px; }
/* ---------- Home (empty state, Claude-Code style) ---------- */
.home { position: relative; max-width: 760px; margin: 22px auto 0; padding: 0 24px; animation: fadeUp .3s ease both; }
.home::before {
  content: ""; position: absolute; top: -120px; left: 50%; transform: translateX(-50%);
  width: 620px; height: 300px; pointer-events: none;
  background: radial-gradient(closest-side, var(--accent-soft), transparent 72%);
  filter: blur(24px);
}
.home > * { position: relative; }
.home-head { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 36px; }
.home-head h1 { display: flex; align-items: center; gap: 13px; margin: 0; font-size: 27px; font-weight: 650; letter-spacing: -.022em; }
.home-mark { display: inline-grid; place-items: center; }
.home-mark .mark { width: 32px; height: 32px; border-radius: 10px; box-shadow: 0 4px 20px var(--accent-soft), 0 1px 3px rgba(0,0,0,.3); }
.home-aside {
  color: var(--good); font-size: 12.5px; font-variant-numeric: tabular-nums; white-space: nowrap;
  padding: 5px 12px; border-radius: 999px;
  border: 1px solid color-mix(in srgb, var(--good) 28%, transparent);
  background: color-mix(in srgb, var(--good) 8%, transparent);
}
.home-section { margin-bottom: 28px; }
.home-label { margin: 0 0 10px; font-size: 12px; color: var(--faint); letter-spacing: .02em; }
.home-list { display: grid; gap: 4px; }
.row {
  width: 100%; text-align: left; font: inherit; cursor: pointer; color: var(--text);
  display: flex; align-items: center; gap: 10px;
  border: 1px solid var(--line); background: var(--bg-2); border-radius: 7px; padding: 9px 12px;
  transition: background .13s, border-color .13s, transform .13s, box-shadow .13s;
}
.row:hover {
  background: var(--surface); border-color: var(--line-2);
  transform: translateY(-1px);
  box-shadow: 0 6px 18px color-mix(in srgb, var(--text) 7%, transparent);
}
.row:active { transform: none; box-shadow: none; }
.row-meta svg { transition: transform .15s, color .15s; }
.row:hover .row-meta svg { transform: translateX(2px); color: var(--accent); }
.row-dot { width: 6px; height: 6px; border-radius: 50%; flex: none; background: var(--muted); }
.row-dot.accent { background: var(--accent); }
.row-dot.hit { background: var(--good); }
.row-dot.miss { background: var(--warn); }
.row-main { min-width: 0; flex: 1; display: flex; align-items: baseline; gap: 9px; }
.row-main strong { font-weight: 600; font-size: 13px; white-space: nowrap; }
.row-sub { min-width: 0; color: var(--muted); font-size: 12px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.row-meta { flex: none; display: flex; align-items: center; gap: 6px; color: var(--faint); font-size: 11px; white-space: nowrap; font-variant-numeric: tabular-nums; }
.row:hover .row-meta { color: var(--muted); }

.msg { max-width: 768px; margin: 0 auto; padding: 14px 24px; display: flex; gap: 14px; animation: fadeUp .25s ease both; }
.msg .av {
  flex: none; width: 30px; height: 30px; border-radius: 9px; display: grid; place-items: center;
  font-size: 13px; font-weight: 700; margin-top: 1px;
}
.msg.user { justify-content: flex-end; }
.msg.user .body {
  background: color-mix(in srgb, var(--accent) 7%, var(--user-bubble));
  border: 1px solid color-mix(in srgb, var(--accent) 14%, transparent);
  padding: 11px 15px; border-radius: 16px 16px 5px 16px; max-width: 78%;
}
.msg.assistant .av {
  background: linear-gradient(135deg, color-mix(in srgb, var(--accent) 72%, white), var(--accent));
  color: var(--on-accent);
  box-shadow: 0 2px 10px var(--accent-soft);
}
.msg.assistant .av svg { width: 18px; height: 18px; }
.body { min-width: 0; }
.body p { margin: 0 0 10px; }
.body p:last-child { margin-bottom: 0; }
.body .typing { display: inline-flex; gap: 4px; padding: 4px 0; }
.body .typing i { width: 7px; height: 7px; border-radius: 50%; background: var(--muted); animation: blink 1.3s infinite both; }
.body .typing i:nth-child(2){ animation-delay:.18s; } .body .typing i:nth-child(3){ animation-delay:.36s; }
@keyframes blink { 0%,80%,100%{opacity:.25;transform:translateY(0);} 40%{opacity:1;transform:translateY(-2px);} }
.body pre.code {
  background: var(--code-bg); border: 1px solid var(--line); border-radius: 12px;
  padding: 14px 16px; overflow-x: auto; margin: 12px 0; font: 13px/1.6 var(--mono); color: var(--code-fg);
  box-shadow: inset 0 1px 0 color-mix(in srgb, var(--text) 4%, transparent);
}
.body code.inline { background: var(--surface-2); border: 1px solid var(--line); border-radius: 6px; padding: 1px 6px; font: 13px var(--mono); }
.body strong { font-weight: 650; }
.msg-meta { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
.tag {
  display: inline-flex; align-items: center; gap: 5px;
  font-size: 11.5px; padding: 3px 9px; border-radius: 999px;
  border: 1px solid var(--line); background: var(--surface-2); color: var(--muted);
  cursor: default; font-variant-numeric: tabular-nums;
}
.tag.click { cursor: pointer; transition: border-color .12s, color .12s; }
.tag.click:hover { border-color: var(--accent); color: var(--text); }
.tag.hit { color: var(--good); border-color: rgba(31,209,139,.35); background: rgba(31,209,139,.08); }
.tag.miss { color: var(--warn); border-color: rgba(227,179,65,.3); }
.tag .d { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }

/* ---------- Composer ---------- */
.composer-wrap { padding: 6px 24px 12px; }
.composer {
  max-width: 768px; margin: 0 auto;
  display: flex; align-items: flex-end; gap: 8px;
  background: var(--surface); border: 1px solid var(--line-2); border-radius: 15px;
  padding: 6px 6px 6px 14px;
  box-shadow: inset 0 1px 0 color-mix(in srgb, var(--text) 4%, transparent), var(--shadow);
  transition: border-color .15s, box-shadow .15s;
}
.composer:focus-within { border-color: color-mix(in srgb, var(--text) 22%, transparent); }
.composer textarea {
  flex: 1; min-width: 0; max-height: 220px; resize: none; border: 0; outline: 0;
  background: transparent; color: var(--text); font: inherit; line-height: 1.5;
  padding: 6px 0;
}
.composer textarea::placeholder { color: var(--faint); }
.send {
  flex: none; width: 32px; height: 32px; border-radius: 50%;
  display: grid; place-items: center; border: 0; cursor: pointer;
  background: linear-gradient(180deg, color-mix(in srgb, var(--accent) 82%, white), var(--accent));
  color: var(--on-accent);
  box-shadow: 0 3px 12px var(--accent-soft);
  transition: opacity .15s, transform .1s, background .15s, box-shadow .15s;
}
.send:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 18px var(--accent-soft); }
.send:active:not(:disabled) { transform: none; }
.send:disabled { background: var(--surface-2); color: var(--faint); cursor: not-allowed; box-shadow: none; }
.composer-foot {
  max-width: 768px; margin: 8px auto 0; padding: 0 2px;
  display: flex; align-items: center; gap: 10px;
}
.model-pick { position: relative; display: inline-flex; }
.model-btn {
  display: inline-flex; align-items: center; gap: 7px;
  padding: 5px 10px; border: 1px solid var(--line); border-radius: 9px;
  background: transparent; color: var(--text); font: inherit; font-size: 13px; cursor: pointer;
  transition: border-color .12s, background .12s;
}
.model-btn:hover { border-color: var(--line-2); background: var(--surface); }
.model-btn > svg:first-child { color: var(--faint); }
.model-btn .caret { color: var(--faint); transition: transform .15s; }
.model-pick.open .model-btn { border-color: var(--line-2); }
.model-pick.open .model-btn .caret { transform: rotate(180deg); }
.model-menu {
  position: absolute; bottom: calc(100% + 6px); left: 0; z-index: 60;
  width: 264px; max-height: 232px; overflow-y: auto;
  background: color-mix(in srgb, var(--elev) 88%, transparent);
  backdrop-filter: blur(18px); -webkit-backdrop-filter: blur(18px);
  border: 1px solid var(--line-2); border-radius: 12px;
  padding: 5px; box-shadow: var(--shadow);
  animation: fadeUp .16s ease both;
}
.model-menu[hidden] { display: none; }
.model-group { margin: 7px 8px 2px; font-size: 10.5px; color: var(--faint); letter-spacing: .02em; }
.model-group:first-child { margin-top: 2px; }
.model-opt {
  display: flex; align-items: center; gap: 7px; width: 100%; text-align: left;
  border: 0; background: transparent; color: var(--text); font: inherit; font-size: 13px;
  padding: 5px 8px; border-radius: 7px; cursor: pointer;
}
.model-opt:hover { background: var(--hover); }
.model-opt.selected { background: var(--surface-2); }
.model-opt .check { flex: none; color: var(--text); opacity: 0; }
.model-opt.selected .check { opacity: 1; }
.composer-hint { font-size: 12px; color: var(--accent); min-height: 1em; }
.composer-legal { margin-left: auto; font-size: 11.5px; color: var(--faint); }
.composer-legal code { font: 11px var(--mono); color: var(--muted); }

/* ---------- Inspector ---------- */
.inspector {
  background: var(--bg-2); border-left: 1px solid var(--line);
  display: flex; flex-direction: column; min-width: 0; overflow: hidden;
}
.app:not(.show-inspector) .inspector { border-left: 0; }
.insp-head { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid var(--line); }
.insp-head h2 { margin: 0; font-size: 14px; font-weight: 600; }
.insp-metrics { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; padding: 14px; border-bottom: 1px solid var(--line); }
.metric {
  background: linear-gradient(180deg, var(--surface-2), var(--surface));
  border: 1px solid var(--line); border-radius: 11px; padding: 12px; display: grid; gap: 4px;
  box-shadow: inset 0 1px 0 color-mix(in srgb, var(--text) 4%, transparent);
}
.metric-val { font-size: 21px; font-weight: 650; font-variant-numeric: tabular-nums; letter-spacing: -.01em; }
#m_saved.metric-val { color: var(--good); text-shadow: 0 0 18px color-mix(in srgb, var(--good) 35%, transparent); }
.metric-lbl { font-size: 11px; color: var(--faint); }
.insp-section { padding: 14px 16px; border-bottom: 1px solid var(--line); }
.insp-section.grow { flex: 1; min-height: 0; display: flex; flex-direction: column; border-bottom: 0; }
.insp-label { margin: 0 0 10px; font-size: 11px; text-transform: uppercase; letter-spacing: .07em; color: var(--faint); }
.insp-label-row { display: flex; align-items: center; justify-content: space-between; }
.decision { display: grid; gap: 8px; }
.dec-action { display: flex; align-items: center; gap: 8px; }
.dec-reason { margin: 0; color: var(--muted); font-size: 13px; }
.dec-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-top: 4px; }
.dec-cell { background: var(--surface); border: 1px solid var(--line); border-radius: 9px; padding: 9px 10px; }
.dec-cell span { display: block; font-size: 10.5px; text-transform: uppercase; letter-spacing: .05em; color: var(--faint); margin-bottom: 3px; }
.dec-cell strong { font-size: 13px; font-weight: 600; word-break: break-word; }
.dec-cell strong.mono { font: 12px var(--mono); color: var(--muted); }
.badge {
  display: inline-flex; align-items: center; gap: 6px;
  font-size: 12px; font-weight: 600; padding: 4px 11px; border-radius: 999px;
  border: 1px solid var(--line-2);
}
.badge.ghost { color: var(--faint); }
.badge.hit { background: rgba(31,209,139,.12); color: var(--good); border-color: rgba(31,209,139,.4); }
.badge.write { background: var(--accent-soft); color: var(--accent); border-color: rgba(77,156,246,.4); }
.badge.pass { background: rgba(227,179,65,.1); color: var(--warn); border-color: rgba(227,179,65,.3); }
.copy-btn { background: transparent; border: 1px solid var(--line); color: var(--muted); border-radius: 7px; font: inherit; font-size: 12px; padding: 3px 9px; cursor: pointer; transition: all .12s; }
.copy-btn:hover { border-color: var(--accent); color: var(--text); }
.json {
  margin: 0; flex: 1; min-height: 0; overflow: auto;
  background: var(--code-bg); border: 1px solid var(--line); border-radius: 10px;
  padding: 13px; font: 12px/1.55 var(--mono); color: var(--code-fg); white-space: pre;
}

/* ---------- Setup gate (terminal) ---------- */
.setup {
  position: fixed; inset: 0; z-index: 200;
  display: grid; place-items: center; padding: 24px;
  background: rgba(5,5,5,.72); backdrop-filter: blur(6px);
}
.setup[hidden] { display: none; }
.term {
  width: min(560px, 94vw);
  background: var(--elev); border: 1px solid var(--line-2); border-radius: 12px;
  box-shadow: var(--shadow); overflow: hidden;
  animation: termIn .18s ease;
}
@keyframes termIn { from { opacity: 0; transform: translateY(6px) scale(.99); } to { opacity: 1; transform: none; } }
.term-bar {
  display: flex; align-items: center; gap: 7px;
  padding: 9px 12px; background: var(--surface-2); border-bottom: 1px solid var(--line);
}
.tdot { width: 11px; height: 11px; border-radius: 50%; }
.tdot.r { background: #ff5f57; } .tdot.y { background: #febc2e; } .tdot.g { background: #28c840; }
.term-title { margin-left: 8px; font: 12px var(--mono); color: var(--faint); }
.term-body { padding: 18px 16px 16px; font: 13px/1.7 var(--mono); }
.term-line, .term-out { margin: 0 0 4px; color: var(--text); }
.term-out { color: var(--muted); }
.term-out.dim { color: var(--faint); }
.term-out code, .term-line code { color: var(--accent); }
.prompt { color: var(--good); margin-right: 8px; user-select: none; }
.term-form { display: flex; align-items: center; gap: 8px; margin-top: 14px; }
.term-input {
  flex: 1; min-width: 0; border: 0; outline: 0; background: transparent;
  color: var(--text); font: 14px var(--mono); caret-color: var(--good);
}
.term-input::placeholder { color: #4a4a4a; }
.term-reveal {
  border: 0; background: transparent; color: var(--faint); font: 12px var(--mono);
  cursor: pointer; padding: 4px 6px;
}
.term-reveal:hover { color: var(--muted); }
.term-connect {
  flex: none; border: 1px solid var(--line-2); background: var(--surface); color: var(--text);
  font: inherit; font-size: 13px; padding: 7px 12px; border-radius: 8px; cursor: pointer;
  transition: background .12s, border-color .12s;
}
.term-connect:hover { background: var(--hover); border-color: var(--accent); }
.term-msg { margin: 12px 0 0; font: 12px var(--mono); min-height: 1em; }
.term-msg.err { color: #ff6b6b; }
.term-msg.ok { color: var(--good); }

/* ---------- Responsive ---------- */
.scrim { display: none; }
@media (max-width: 1100px) {
  .app.show-inspector { grid-template-columns: 264px minmax(0,1fr) 340px; }
}
@media (max-width: 880px) {
  .app, .app.show-inspector, .app.hide-sidebar { grid-template-columns: 1fr; }
  .sidebar, .inspector {
    position: fixed; top: 0; bottom: 0; z-index: 40; width: 300px; box-shadow: var(--shadow);
  }
  .sidebar { left: 0; transform: translateX(-100%); transition: transform .22s ease; }
  .inspector { right: 0; transform: translateX(100%); transition: transform .22s ease; }
  .app.show-sidebar-m .sidebar { transform: none; }
  .app.show-inspector .inspector { transform: none; }
  .app.show-sidebar-m .scrim, .app.show-inspector .scrim {
    display: block; position: fixed; inset: 0; z-index: 30; background: rgba(0,0,0,.5);
  }
  .side-resize { display: none; }
  .home-head { flex-direction: column; align-items: flex-start; gap: 8px; }
  .row-main { flex-direction: column; align-items: flex-start; gap: 2px; }
  .row-sub { white-space: normal; }
}
`
}

func consoleJS() string {
	return `
const $ = (id) => document.getElementById(id);
const app = document.querySelector('.app');
const state = { rows: [], selected: null, pending: new Map() };

const money = (v, d = 4) => '$' + Number(v || 0).toFixed(d);
const pct = (v) => Math.round(Number(v || 0) * 100) + '%';
const escapeHTML = (v) => String(v ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));

/* ---- markdown-lite ---- */
const BT = String.fromCharCode(96);
const FENCE = BT + BT + BT;
function renderMarkdown(text) {
  const blocks = [];
  let s = String(text ?? '');
  s = s.replace(new RegExp(FENCE + '(\\w*)\\n?([\\s\\S]*?)' + FENCE, 'g'), (m, lang, code) => {
    blocks.push('<pre class="code"><code>' + escapeHTML(code.replace(/\n$/, '')) + '</code></pre>');
    return ' B' + (blocks.length - 1) + ' ';
  });
  s = escapeHTML(s);
  s = s.replace(new RegExp(BT + '([^' + BT + ']+)' + BT, 'g'), (m, c) => '<code class="inline">' + c + '</code>');
  s = s.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
  s = s.split(/\n{2,}/).map(p => '<p>' + p.replace(/\n/g, '<br>') + '</p>').join('');
  s = s.replace(/ B(\d+) /g, (m, i) => blocks[Number(i)]);
  return s;
}

/* ---- decision logic mirrors proxy semantics ---- */
function decisionFor(row) {
  if (!row) return { kind: 'idle', action: 'Idle', cls: 'ghost', reason: 'No intercepted request yet.', saved: 0 };
  const hit = row.cache_status === 'hit';
  const cacheable = row.miser_cache_eligible === true;
  if (hit) return { kind: 'hit', action: 'Exact cache hit', cls: 'hit',
    reason: 'Miser matched the request fingerprint and returned the stored provider response with zero new spend.',
    saved: Number(row.cache_saved_usd || 0) };
  if (cacheable) return { kind: 'write', action: 'Pass-through + cache write', cls: 'write',
    reason: 'First time seeing this fingerprint. Miser forwarded it upstream, priced the tokens, and cached the response for an identical repeat.',
    saved: 0 };
  return { kind: 'pass', action: 'Pass-through', cls: 'pass',
    reason: 'Not cache-eligible — usually a streaming call or an unsupported endpoint. Forwarded untouched.',
    saved: 0 };
}

/* ---- thread ---- */
function appendUser(text) {
  $('welcome')?.remove();
  const el = document.createElement('div');
  el.className = 'msg user';
  el.innerHTML = '<div class="body">' + renderMarkdown(text) + '</div>';
  $('thread').appendChild(el);
  scrollThread();
}
function appendAssistant() {
  const el = document.createElement('div');
  el.className = 'msg assistant';
  el.innerHTML = '<div class="av">' + markSVG() + '</div><div class="body"><div class="typing"><i></i><i></i><i></i></div></div>';
  $('thread').appendChild(el);
  scrollThread();
  return el.querySelector('.body');
}
function markSVG() {
  return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round"><path d="M5 18V7l4 5 4-5v11"/><circle cx="18.5" cy="15.5" r="1.6" fill="currentColor" stroke="none"/></svg>';
}
function scrollThread() { const t = $('thread'); t.scrollTop = t.scrollHeight; }

function metaTags(body, row) {
  const d = decisionFor(row);
  const tags = [];
  tags.push('<span class="tag ' + (d.kind === 'hit' ? 'hit' : (d.kind === 'pass' ? 'miss' : '')) + '"><span class="d"></span>' + escapeHTML(d.action) + '</span>');
  if (row.model) tags.push('<span class="tag">' + escapeHTML(row.model) + '</span>');
  if (d.saved > 0) tags.push('<span class="tag hit">saved ' + money(d.saved) + '</span>');
  else tags.push('<span class="tag">cost ' + money(row.cost_usd) + '</span>');
  if (row.latency_ms != null) tags.push('<span class="tag">' + Math.round(row.latency_ms) + ' ms</span>');
  tags.push('<span class="tag click" data-inspect="1">inspect →</span>');
  const wrap = document.createElement('div');
  wrap.className = 'msg-meta';
  wrap.innerHTML = tags.join('');
  wrap.querySelector('[data-inspect]')?.addEventListener('click', () => { select(row); openInspector(); });
  body.appendChild(wrap);
}

/* ---- data ---- */
async function refresh() {
  try {
    const res = await fetch('/miser/api/requests', { cache: 'no-store' });
    state.rows = await res.json() || [];
  } catch (e) { return; }
  renderStats();
  renderHistory();
  renderHome();
  if (!state.selected && state.rows[0]) select(state.rows[0]);
  else if (state.selected) {
    const match = state.rows.find(r => r.id === state.selected.id);
    if (match) select(match);
  }
}

function renderStats() {
  const rows = state.rows;
  const reqs = rows.length;
  const hits = rows.filter(r => r.cache_status === 'hit').length;
  const spend = rows.reduce((s, r) => s + Number(r.cost_usd || 0), 0);
  const saved = rows.filter(r => r.cache_status === 'hit').reduce((s, r) => s + Number(r.cache_saved_usd || 0), 0);
  $('m_requests').textContent = String(reqs);
  $('m_saved').textContent = money(saved, 2);
  $('m_rate').textContent = reqs ? pct(hits / reqs) : '0%';
  $('m_spend').textContent = money(spend, 2);
  const hs = $('homeSaved'); if (hs) hs.textContent = money(saved, 2) + ' saved so far';
}

function chevronSVG() {
  return '<svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>';
}

function renderHome() {
  const sec = $('trafficSection'), list = $('homeTraffic');
  if (!sec || !list) return;
  const rows = state.rows.slice(0, 5);
  if (!rows.length) { sec.style.display = 'none'; return; }
  sec.style.display = '';
  list.innerHTML = rows.map((r, i) => {
    const d = decisionFor(r);
    const t = r.timestamp ? new Date(r.timestamp).toLocaleTimeString([], {hour:'2-digit',minute:'2-digit'}) : '';
    const dot = d.kind === 'hit' ? 'hit' : (d.kind === 'pass' ? 'miss' : 'accent');
    const right = d.saved > 0 ? ('saved ' + money(d.saved, 4)) : (money(r.cost_usd, 4) + ' cost');
    return '<button class="row" data-i="' + i + '">' +
      '<span class="row-dot ' + dot + '"></span>' +
      '<span class="row-main"><strong>' + escapeHTML(r.model || 'unknown') + '</strong><span class="row-sub">' + escapeHTML(d.action) + '</span></span>' +
      '<span class="row-meta">' + right + ' · ' + t + chevronSVG() + '</span>' +
    '</button>';
  }).join('');
  list.querySelectorAll('.row').forEach(n => n.addEventListener('click', () => {
    select(rows[Number(n.dataset.i)]); openInspector();
  }));
}

function renderHistory() {
  const el = $('history');
  const rows = state.rows.slice(0, 40);
  if (!rows.length) { el.innerHTML = '<p class="history-empty">No intercepted requests yet.</p>'; return; }
  el.innerHTML = rows.map((r, i) => {
    const d = decisionFor(r);
    const t = r.timestamp ? new Date(r.timestamp).toLocaleTimeString([], {hour:'2-digit',minute:'2-digit'}) : '';
    return '<button class="hist-item" data-i="' + i + '">' +
      '<span class="hist-line"><span class="tag ' + (d.kind==="hit"?"hit":(d.kind==="pass"?"miss":"")) + '" style="padding:1px 6px"><span class="d"></span></span>' +
      '<span class="hist-model">' + escapeHTML(r.model || 'unknown') + '</span></span>' +
      '<span class="hist-sub"><span>' + escapeHTML(d.action) + '</span><span>' + t + '</span></span>' +
    '</button>';
  }).join('');
  el.querySelectorAll('.hist-item').forEach(n => n.addEventListener('click', () => {
    const r = rows[Number(n.dataset.i)]; select(r); if (window.innerWidth <= 880) openInspector();
  }));
}

function select(row) {
  state.selected = row;
  const d = decisionFor(row);
  $('decision').innerHTML =
    '<div class="dec-action"><span class="badge ' + d.cls + '">' + escapeHTML(d.action) + '</span></div>' +
    '<p class="dec-reason">' + escapeHTML(d.reason) + '</p>' +
    '<div class="dec-grid">' +
      cell('Route', escapeHTML((row.provider||'provider') + ' / ' + (row.model||'model'))) +
      cell('Saved', money(d.saved)) +
      cell('Status', escapeHTML(String(row.http_status || '—'))) +
      cell('Cost basis', escapeHTML(row.cost_basis || '—')) +
      cell('Fingerprint', escapeHTML(row.request_fingerprint || '—'), true) +
      cell('Latency', (row.latency_ms != null ? Math.round(row.latency_ms) + ' ms' : '—')) +
    '</div>';
  $('inspectorJson').textContent = JSON.stringify(inspectorPayload(row, d), null, 2);
}
function cell(label, value, mono) {
  return '<div class="dec-cell"><span>' + label + '</span><strong' + (mono?' class="mono"':'') + '>' + value + '</strong></div>';
}
function inspectorPayload(row, d) {
  return {
    original_request: { path: row.http_path, method: row.http_method, model: row.model, prompt_fingerprint: row.request_fingerprint, prompt_chars: row.prompt_chars },
    miser_decision: { action: d.action, reason: d.reason, cache_status: row.cache_status, cache_eligible: row.miser_cache_eligible },
    final_route: { provider: row.provider, model: row.model, http_status: row.http_status },
    result: { cost_after_miser: Number(row.cost_usd || 0), estimated_saved: d.saved, input_tokens: row.input_tokens, output_tokens: row.output_tokens, input_cached_tokens: row.input_cached_tokens, latency_ms: row.latency_ms, cost_basis: row.cost_basis }
  };
}

/* ---- chat ---- */
function assistantText(p) {
  const c = p?.choices?.[0];
  return c?.message?.content || c?.text || p?.output_text ||
    (Array.isArray(p?.content) ? p.content.map(x => x.text).filter(Boolean).join('') : '') ||
    FENCE + 'json\n' + JSON.stringify(p, null, 2) + '\n' + FENCE;
}

$('chatForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  const ta = $('prompt');
  const prompt = ta.value.trim();
  if (!prompt) return;
  ta.value = ''; autogrow();
  appendUser(prompt);
  const body = appendAssistant();
  $('sendBtn').disabled = true;
  try {
    const res = await fetch('/v1/chat/completions', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: $('model').value, messages: [{ role: 'user', content: prompt }] })
    });
    const cacheHdr = res.headers.get('X-Miser-Cache');
    const payload = await res.json();
    body.innerHTML = renderMarkdown(assistantText(payload));
    if (cacheHdr) {
      $('cacheHint').textContent = cacheHdr === 'HIT' ? '✓ served from Miser cache' : '';
      setTimeout(() => $('cacheHint').textContent = '', 2600);
    }
    await refresh();
    const latest = state.rows[0];
    if (latest) metaTags(body, latest);
  } catch (err) {
    body.innerHTML = '<p style="color:var(--warn)">Request failed: ' + escapeHTML(err.message) + '</p>' +
      '<p style="color:var(--faint);font-size:13px">Is the proxy pointed at a provider with a valid API key?</p>';
  }
  scrollThread();
});

/* ---- composer UX ---- */
function autogrow() {
  const ta = $('prompt');
  ta.style.height = 'auto';
  ta.style.height = Math.min(ta.scrollHeight, 220) + 'px';
  $('sendBtn').disabled = ta.value.trim() === '';
}
$('prompt').addEventListener('input', autogrow);
$('prompt').addEventListener('keydown', (e) => {
  if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); $('chatForm').requestSubmit(); }
});
document.querySelectorAll('.row[data-prompt]').forEach(c => c.addEventListener('click', () => {
  $('prompt').value = c.dataset.prompt; autogrow(); $('chatForm').requestSubmit();
}));
$('newChat').addEventListener('click', () => {
  document.querySelectorAll('.msg').forEach(m => m.remove());
  $('cacheHint').textContent = '';
  if (!$('welcome')) location.reload();
});

/* ---- panels ---- */
function openInspector() { app.classList.add('show-inspector'); }
function toggleInspector() { app.classList.toggle('show-inspector'); }
$('inspectorToggle').addEventListener('click', toggleInspector);
$('inspClose').addEventListener('click', () => app.classList.remove('show-inspector'));
$('sidebarToggle').addEventListener('click', () => {
  if (window.innerWidth <= 880) app.classList.toggle('show-sidebar-m');
  else app.classList.toggle('hide-sidebar');
});
$('scrim').addEventListener('click', () => app.classList.remove('show-sidebar-m', 'show-inspector'));

/* ---- model picker (custom upward, scrollable, collapsible) ---- */
(() => {
  const pick = $('modelPick'), btn = $('modelBtn'), menu = $('modelMenu'),
        input = $('model'), label = $('modelLabel');
  if (!pick) return;
  const open = () => {
    menu.hidden = false; pick.classList.add('open');
    const sel = menu.querySelector('.model-opt.selected');
    if (sel) sel.scrollIntoView({ block: 'nearest' });
  };
  const close = () => { menu.hidden = true; pick.classList.remove('open'); };
  btn.addEventListener('click', (e) => { e.stopPropagation(); menu.hidden ? open() : close(); });
  menu.querySelectorAll('.model-opt').forEach(o => o.addEventListener('click', () => {
    input.value = o.dataset.model;
    label.textContent = o.dataset.model;
    menu.querySelectorAll('.model-opt').forEach(x => x.classList.toggle('selected', x === o));
    close();
  }));
  document.addEventListener('click', (e) => { if (!pick.contains(e.target)) close(); });
  document.addEventListener('keydown', (e) => { if (e.key === 'Escape') close(); });
})();

/* ---- sidebar resize (drag the right edge) ---- */
const SIDE_MIN = 180, SIDE_MAX = 460;
function setSideWidth(px) {
  const w = Math.max(SIDE_MIN, Math.min(SIDE_MAX, Math.round(px)));
  document.documentElement.style.setProperty('--side-w', w + 'px');
  return w;
}
(() => {
  const saved = Number(localStorage.getItem('miser_side_w'));
  if (saved >= SIDE_MIN && saved <= SIDE_MAX) document.documentElement.style.setProperty('--side-w', saved + 'px');
})();
(() => {
  const handle = $('sideResize');
  if (!handle) return;
  let dragging = false;
  const onMove = (e) => { if (dragging) setSideWidth(e.clientX); };
  const stop = () => {
    if (!dragging) return;
    dragging = false; app.classList.remove('resizing');
    const w = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--side-w'), 10);
    if (w) localStorage.setItem('miser_side_w', w);
  };
  handle.addEventListener('mousedown', (e) => {
    if (window.innerWidth <= 880) return;
    dragging = true; app.classList.add('resizing'); e.preventDefault();
  });
  handle.addEventListener('dblclick', () => { setSideWidth(230); localStorage.setItem('miser_side_w', 230); });
  window.addEventListener('mousemove', onMove);
  window.addEventListener('mouseup', stop);
})();

$('copyJson').addEventListener('click', async () => {
  try { await navigator.clipboard.writeText($('inspectorJson').textContent);
    const b = $('copyJson'); b.textContent = 'Copied'; setTimeout(() => b.textContent = 'Copy', 1400);
  } catch (e) {}
});

/* ---- setup gate: connect a provider API key ---- */
async function checkConfig() {
  let cfg;
  try { cfg = await (await fetch('/miser/api/config', { cache: 'no-store' })).json(); }
  catch (e) { return false; }
  if (cfg.key_env) $('termEnv').textContent = cfg.key_env;
  if (cfg.provider) $('termProvider').textContent = cfg.provider;
  const ks = $('keyStatus');
  if (ks) ks.textContent = cfg.configured ? (cfg.masked || 'connected') : 'not set';
  $('setup').hidden = !!cfg.configured;
  if (!cfg.configured) setTimeout(() => $('keyInput').focus(), 60);
  return !!cfg.configured;
}
function openSetup() { $('setup').hidden = false; setTimeout(() => $('keyInput').focus(), 60); }
$('keyStatus')?.addEventListener('click', openSetup);

(() => {
  const form = $('keyForm');
  if (!form) return;
  const input = $('keyInput'), msg = $('keyMsg'), reveal = $('keyReveal'), connect = $('keyConnect');
  reveal.addEventListener('click', () => {
    const show = input.type === 'password';
    input.type = show ? 'text' : 'password';
    reveal.textContent = show ? 'hide' : 'show';
    input.focus();
  });
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const key = input.value.trim();
    if (!key) return;
    msg.className = 'term-msg'; msg.textContent = 'connecting…';
    connect.disabled = true;
    try {
      const res = await fetch('/miser/api/key', {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'failed to set key');
      msg.className = 'term-msg ok'; msg.textContent = '✓ connected ' + (data.masked || '');
      input.value = '';
      const ks = $('keyStatus'); if (ks) ks.textContent = data.masked || 'connected';
      setTimeout(() => { $('setup').hidden = true; msg.textContent = ''; refresh(); }, 550);
    } catch (err) {
      msg.className = 'term-msg err'; msg.textContent = '✗ ' + err.message;
    } finally { connect.disabled = false; }
  });
})();

/* ---- theme + settings menu ---- */
function applyTheme(name) {
  if (name && name !== 'midnight') document.body.dataset.theme = name;
  else { delete document.body.dataset.theme; name = 'midnight'; }
  localStorage.setItem('miser_theme', name);
  document.querySelectorAll('.swatch').forEach(s => s.classList.toggle('active', s.dataset.theme === name));
}
applyTheme(localStorage.getItem('miser_theme') || 'midnight');
(() => {
  const btn = $('settingsBtn'), menu = $('settingsMenu');
  if (!btn || !menu) return;
  const close = () => { menu.hidden = true; };
  btn.addEventListener('click', (e) => { e.stopPropagation(); menu.hidden = !menu.hidden; });
  document.querySelectorAll('.swatch').forEach(s => s.addEventListener('click', () => applyTheme(s.dataset.theme)));
  $('disconnectBtn')?.addEventListener('click', async () => {
    try { await fetch('/miser/api/key', { method: 'DELETE' }); } catch (e) {}
    close();
    const ks = $('keyStatus'); if (ks) ks.textContent = 'not set';
    openSetup();
  });
  document.addEventListener('click', (e) => { if (!menu.contains(e.target) && !btn.contains(e.target)) close(); });
  document.addEventListener('keydown', (e) => { if (e.key === 'Escape') close(); });
})();

autogrow();
checkConfig();
refresh();
setInterval(refresh, 3500);
`
}
