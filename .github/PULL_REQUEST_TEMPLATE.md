## Summary

What changed?

## Why

What problem does this solve?

## Testing

- [ ] `go test ./cmd/... ./internal/...`
- [ ] Manual proxy/console check, if applicable

## Cost And Safety

- [ ] Does not label estimated token cost as actual spend
- [ ] Preserves `account_id`, `integration`, and `cost_basis` where relevant
- [ ] Does not add secrets, prompts, invoices, proxy logs, or caches
- [ ] Adds or updates tests for pricing, importers, proxy behavior, or audit math when touched
