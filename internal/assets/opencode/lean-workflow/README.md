# Lean OpenCode Workflow Pilot

This pilot installs a separate OpenCode agent set for requirements-first code changes. It is not a replacement for the existing SDD workflow.

## Happy Path

1. Start OpenCode with the `lean-workflow-orchestrator` agent.
2. Give it a code-changing request, or say `resume <task-slug>` to continue an existing local task folder.
3. Review `.github/plans/<task-slug>/01-requirements.md` and reply `confirmed` or add feedback.
4. Review `.github/plans/<task-slug>/02-plan.md` and reply `confirmed` or add feedback.
5. Write or update the executable acceptance test yourself, then reply with an `@` file mention or explicit path.
6. Review `.github/plans/<task-slug>/04-todo.md` and reply `confirmed` before implementation starts.
7. Let the implementor work from the confirmed todo and update verification results.

## Artifact Roles

- `01-requirements.md` records the bounded requirements grill, decisions, assumptions, and open questions.
- `02-plan.md` records the proposed implementation approach and test inventory.
- `03-test-guidance.md` is optional and only appears when you request guidance before writing the acceptance test.
- `04-todo.md` is the implementor contract.
- `05-checkpoint.md` is used when implementation discovers a material surprise.

Each artifact starts with small frontmatter: `status` and `owner`. Human confirmation must be written back to the artifact before the next phase starts.

## Acceptance Test Gate

The acceptance test is human-owned executable coverage. The workflow may suggest scenarios or pseudocode, but it must not write the acceptance test unless you explicitly authorize a narrow mechanical edit later.

You may waive the gate in `02-plan.md` frontmatter when executable acceptance coverage is disproportionate:

```yaml
acceptance_test:
  required: waived
  waiver_reason: "Documentation-only change; no executable behavior changes."
  waived_by: human
```

## Optional Skills

The pilot keeps optional skills separate from the core workflow. Install add-on skills explicitly with the existing skill flags, for example:

```bash
gentle-ai install --agent opencode --component skills --skills tdd,caveman
```

The `tdd` skill is vendored from Matt Pocock's MIT-licensed skills repository and can help the planner, test-guidance planner, todo generator, or implementor follow red-green-refactor habits. It does not bypass workflow gates: the executable acceptance test remains human-owned, and implementation still waits for a confirmed todo.

The `caveman` skill is also optional. When installed, lean workflow prompts tell subagents to read it and use terse, token-conscious communication for routing chatter, summaries, checkpoints, and artifact prose. Required artifact sections, commands, paths, errors, and verification results must remain clear and exact.

## Model Control

The installer writes lean workflow agents into `opencode.json` with the same model-control shape used by SDD OpenCode agents. Explicit model assignments for `lean-workflow-orchestrator`, `requirements-griller`, `planner`, `test-guidance-planner`, `todo-generator`, or `implementor` are written onto those agent definitions. Existing agent model settings are preserved. New lean workflow agents fall back to the root OpenCode `model` when one exists, so subagents do not silently drift to an unintended runtime default.

## Memory Boundary

Use Markdown memory only when context should survive beyond the local task folder: domain terms, recurring constraints, important decisions, or paused-work handoffs. Task artifacts remain the source of truth for the active task.

## Local Files

The installer adds narrow gitignore entries for:

- `.github/plans/`
- `.github/lean-workflow-templates/`

The copied templates are local and editable during validation. Reinstalling does not overwrite them.

## Promotion Notes

Treat this as validation scaffolding. Promote it to a reusable recipe only after real tasks show that the requirements grill, test gates, todo contract, memory use, and bash approval friction are worth keeping.
