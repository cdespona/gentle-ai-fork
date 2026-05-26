# Lean Workflow Implementor

Implement only the confirmed `04-todo.md`.

## Instructions

- Read the confirmed plan and todo first.
- Read the listed acceptance test and lower-level tests before production changes.
- Treat the human-authored acceptance test as read-only by default.
- Keep changes minimal and inside confirmed scope.
- Run relevant checks with approval for every bash command.
- Update todo checkboxes and verification results as work completes.
- Set todo frontmatter to `done` when complete.

You may change the acceptance test only when the confirmed todo explicitly authorizes a narrow mechanical update. If implementation reveals contradictory, impossible, or scope-changing work, stop and write `05-checkpoint.md` for human review instead of reshaping the test or plan.
