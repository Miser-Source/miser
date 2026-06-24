# Contributing

Miser is early, but the bar for cost and billing logic is high. Keep changes small, testable, and easy to review.

## Workflow

Use a branch for each change:

```bash
git checkout -b feature/name
go test ./cmd/... ./internal/...
git push -u origin feature/name
```

Open a pull request into `main`. CI must pass before merge.

## Good First Contributions

- Add provider/model pricing coverage with tests
- Improve importers for provider billing exports
- Add request inspector fields that help teams debug spend
- Add policy/rule generators for common waste patterns
- Improve docs, examples, and integration guides

## What To Preserve

- Do not call estimated token cost actual spend.
- Keep `account_id`, `integration`, and `cost_basis` on imported rows when possible.
- Prefer inspectable savings logic over opaque automation.
- Add tests when changing importers, audit math, proxy behavior, pricing, or route recommendations.
- Treat prompts, billing files, and logs as sensitive by default.

## Cost Basis

Miser uses these cost basis values:

- `actual_invoice`: billing export or invoice data
- `provider_billing_api`: provider billing API data
- `actual_invoice_allocated`: actual invoice dollars allocated across usage rows
- `reported_log_cost`: request logs with provider-reported cost
- `published_token_price`: token usage priced from a known model catalog
- `estimated_token_cost`: estimated token/API value, not an invoice
- `unpriced_token_usage`: token usage for a model Miser does not know yet
- `miser_exact_cache`: request served by Miser cache with zero provider spend
