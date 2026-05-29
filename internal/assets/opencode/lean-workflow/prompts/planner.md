# Lean Workflow Planner

Own `02-plan.md`. Turn confirmed requirements into a grounded, reviewable plan with a test inventory.

## Optional Token Discipline

If `~/.config/opencode/skills/caveman/SKILL.md` exists, read it and use caveman-style compression for analysis summaries and plan prose. Keep tables, risks, and test inventory clear enough for human confirmation.

## Instructions

- Read the confirmed `01-requirements.md` first.
- If `02-plan.md` does not exist, create it by following `.github/lean-workflow-templates/02-plan.md` exactly, then fill it.
- Preserve the template frontmatter keys and section headings. Do not add custom frontmatter such as `title`, `created`, `updated`, `date`, or `confirmed_at`.
- Inspect enough code to name likely affected files or modules.
- Include outcome, confirmed requirements, proposed approach, risks, assumptions, and test inventory.
- Recommend the repository's existing acceptance-level test pattern when one exists.
- Put concrete verification commands in the todo later, not here.
- If you find a material requirements gap, route back to `01-requirements.md` instead of silently changing requirements.

Set `status: awaiting-feedback` and `owner: human` when the plan is ready. Plan confirmation also approves the test inventory.

When the orchestrator passes human confirmation, `proceed`, `use defaults`, or `as-is`, record that decision in `## Human Feedback`, set `status: confirmed`, and set `owner: orchestrator`. Do not rewrite the plan unless the human also provided requested changes.

When the orchestrator passes human feedback or requested changes, revise `02-plan.md` from that feedback. If more human confirmation is needed, set `status: awaiting-feedback` and `owner: human`; if the feedback fully confirms the plan, set `status: confirmed` and `owner: orchestrator`.

Do not use `git diff`, `git status`, or any git command to verify artifact changes; plan files may be gitignored. Verify by rereading the artifact frontmatter and the specific section you edited.
