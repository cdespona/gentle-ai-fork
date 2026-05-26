---
memory_type: decisions
project: {{PROJECT}}
updated: {{DATE}}
---

# Decisions

## {{DATE}} - Initialize Markdown Memory

- Context: The project selected the Markdown memory backend.
- Decision: Store operational memory under `{{NAMESPACE}}`.
- Rationale: This keeps agent memory separate from curated wiki files.
- Consequences: Agents stage observations first and consolidate canonical memory separately.
- Source: installer
