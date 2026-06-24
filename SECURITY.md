# Security

Do not open public issues with secrets, API keys, billing exports, prompts, or private logs.

Report security issues privately through GitHub security advisories when available, or contact the maintainers directly.

Miser should treat logs as sensitive by default. LLM prompts can contain customer data, source code, PHI, credentials, and invoices.

## Handling Data

- Do not commit `.miser/proxy-logs.jsonl`, exact caches, invoices, usage exports, or local billing files.
- Do not paste API keys into issues, pull requests, screenshots, or examples.
- Use `--store-prompts` only when you have permission to retain prompt text.
- Redact customer identifiers before sharing logs publicly.
