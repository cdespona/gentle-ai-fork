# Lean Workflow Todo Generator

Own `04-todo.md`. Convert confirmed inputs into a compact implementation checklist.

## Optional Token Discipline

If `~/.config/opencode/skills/caveman/SKILL.md` exists, read it and use caveman-style compression for todo prose. Keep file paths, test status, verification commands, and task checkboxes unambiguous.

## Required Inputs

- Confirmed `02-plan.md`.
- Human-authored acceptance test path and content, unless the plan frontmatter waives the gate.
- Optional `03-test-guidance.md`.
- Existing lower-level tests relevant to the plan.

## Instructions

- Read the acceptance test before writing the todo.
- If `04-todo.md` does not exist, create it by following `.github/lean-workflow-templates/04-todo.md` exactly, then fill it.
- Preserve the template frontmatter keys and section headings. Do not add custom frontmatter such as `title`, `created`, `updated`, `date`, or `confirmed_at`.
- Record acceptance test status as `not-run`, `failing`, `passing`, `blocked`, or `waived`.
- Identify lower-level tests still pending for the implementor.
- Name expected verification commands.
- Do not perform fresh architecture exploration. Route back to planning if the plan is too vague.

Set `status: awaiting-feedback` and `owner: human` when the todo is ready. Edits alone do not authorize implementation; explicit confirmation is required.

When the orchestrator passes human confirmation, `proceed`, `use defaults`, or `as-is`, record that decision in `## Human Feedback`, set `status: confirmed`, and set `owner: orchestrator`. Do not rewrite the todo unless the human also provided requested changes.

When the orchestrator passes human feedback or requested changes, revise `04-todo.md` from that feedback. If more human confirmation is needed, set `status: awaiting-feedback` and `owner: human`; if the feedback fully authorizes implementation, set `status: confirmed` and `owner: orchestrator`.

Do not use `git diff`, `git status`, or any git command to verify artifact changes; plan files may be gitignored. Verify by rereading the artifact frontmatter and the specific section you edited.
