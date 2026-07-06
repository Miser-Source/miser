# Changelog

All notable changes to Miser will be documented here.

## Unreleased

### Added

- Live Claude (Anthropic) proxy path: `/v1/messages` exact caching, automatic `x-api-key` / `anthropic-version` injection, prompt-cache-aware token normalization, and published pricing for Claude 4.x models.
- Runtime provider switching (OpenAI ⇄ Claude) from the console setup gate, with per-provider model pickers and playground payloads.
- Individual and business deployment profiles: `--mode`, `--workspace`, and `--budget` flags, a `POST /miser/api/profile` endpoint, month-to-date spend tracking that survives restarts, `X-Miser-Budget` response headers, and a live budget bar in the console inspector.
- Redesigned browser console: ChatGPT/Codex-style chat thread, collapsible Miser inspector (decision trace + live request JSON), and a Claude-Code-style home with Quick start and live Recent traffic.
- Terminal-style setup gate to connect a provider API key from the browser; the proxy now starts without a key and accepts it at runtime via `POST /miser/api/key` (held in memory only).
- Disconnect action and a settings menu with selectable color themes (Midnight, White, plus vivid full-screen themes) that recolor the entire console.
- Resizable, draggable sidebar and the full live OpenAI model lineup in the playground model picker.
- Live OpenAI-compatible proxy.
- Browser console with playground, request inspector, decision trace, and savings metrics.
- Exact cache for repeated non-streaming chat/responses requests.
- JSONL audit, analysis, plan, rules, and apply commands.
- OpenAI organization usage and costs pullers.
- `ccusage` and invoice CSV importers.
- Actual invoice reconciliation.
- Published pricing for known OpenAI and Claude models.
- GitHub Actions CI and Dependabot.
