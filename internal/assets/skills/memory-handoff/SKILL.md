---
name: memory-handoff
description: "Trigger: create or update a task-scoped Markdown memory handoff."
license: Apache-2.0
---

# Memory Handoff

Use when another agent or future session needs to resume the same task.

## Protocol

- Write `machine/agent-memory/projects/<project>/handoffs/<task-slug>.md`.
- Link the current task session file.
- Include current state, last completed step, next action, files or areas involved, decisions to preserve, work not to repeat, and blockers.
- Keep one handoff per task slug. Do not create a single project-level `handoff.md`.
