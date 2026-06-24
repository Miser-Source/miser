# Roadmap

Miser is moving from audit tool to runtime control plane.

## Now

- OpenAI-compatible proxy
- Browser console with playground, decision trace, request inspector, and savings metrics
- Exact cache for repeated non-streaming chat/responses requests
- JSONL audit, plan, and rule generation
- OpenAI usage/cost pullers
- `ccusage` and invoice CSV importers
- Actual invoice reconciliation

## Next

- Semantic cache with quality guardrails
- Model routing policy engine
- Request compression and folded tool output policies
- Provider-specific Anthropic proxy path
- Team/project/API-key views in the console
- Realized savings reports by workflow
- Rule approval workflow before automatic enforcement

## Later

- Local-model fallback adapters
- Replay eval harness
- Generated code/config patches for stable workflows
- Hosted Miser Cloud control plane
- Organization, RBAC, audit trail, and SSO
- Marketplace-style integrations for agent frameworks and provider gateways

## Non-Goals

- Miser should not pretend every workflow can save 90%.
- Miser should not silently route without quality evidence.
- Miser should not store full prompts by default.
- Miser should not mix estimated token cost with actual invoice spend.
