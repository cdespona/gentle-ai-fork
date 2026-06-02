# Layered TDD OpenCode Workflow

Layered TDD is a separate OpenCode workflow for one selected slice of code-changing work. It replaces the standalone plan artifact with a requirements grill, a layer map, human-approved layer todos, red-test gates, implementation inside approved layer boundaries, and review after each layer.

## Quick Path

1. Start OpenCode with the `layered-tdd-orchestrator` agent.
2. Give it one code-changing request, or say `resume <task-slug>` to continue an existing task folder.
3. Confirm `.github/plans/<task-slug>/00-requirements.md`.
4. If the request has multiple independently valuable slices, select exactly one slice from `slice-selection.md`; the selected slice starts a fresh task folder.
5. Confirm `.github/plans/<slice-slug>/01-layer-map.md`.
6. Choose the next layer to run.
7. Review the selected layer todo under `.github/plans/<slice-slug>/layers/` and approve its top-level Gherkin proposals and test ownership.
8. Provide or approve the top-level red-test gate for that layer.
9. Let the implementor work only inside the approved layer boundary.
10. Review the layer result, then select the next layer until the slice is complete.
11. Review `99-final-review.md` and approve any memory candidates worth keeping.

## Artifact Roles

| Artifact | Purpose |
| --- | --- |
| `00-requirements.md` | Captures the selected slice goal, blockers, non-blocker assumptions, out-of-scope work, and confirmation state. |
| `slice-selection.md` | Used only when the original request contains multiple independently valuable slices. The workflow stops until exactly one slice is selected. |
| `01-layer-map.md` | Lists the layers for the selected slice, why each layer matters, the recommended order, and skeleton todo filenames. |
| `layers/<nn>-<layer>.todo.md` | Defines one layer's goal, test ownership, top-level Gherkin proposals, red-test gate, implementation scope, tasks, verification, and checkpoints. |
| `99-final-review.md` | Records final slice review, waived layers, residual risks, and memory candidates. |

Each artifact uses small frontmatter with `status` and `owner`. Human confirmation must be written back into the artifact before the workflow moves forward.

## Slice Selection

One workflow run equals one selected slice. If the requirements grill finds multiple valuable slices, it writes `slice-selection.md` and stops. After you select one slice, the orchestrator starts a new slice-specific task folder with a fresh `00-requirements.md`.

The slice-specific requirements pass should be short. It confirms the selected slice goal, blockers for that slice, assumptions for that slice, and out-of-scope slices.

## Layer Gates

After requirements are confirmed, the layer mapper writes `01-layer-map.md`. The map should follow the repository's real boundaries. Ports-and-adapters and deep-module boundaries are useful when present, but they are not required.

After the layer map is approved, skeleton todos are created for all layers so the slice shape is visible. Only the human-selected layer is fully detailed. The orchestrator stops after every layer and asks which layer to run next.

## Test Ownership

Layer test ownership applies to the whole layer:

| Mode | Meaning |
| --- | --- |
| `human-written` | You write the top-level red tests. The agent implements production code after you confirm the red-test gate. |
| `agent-written-after-approval` | You approve Gherkin proposals first. The agent writes those top-level tests, then stops for confirmation before production code. |
| `waived` | No top-level red test for this layer. The todo must include a human-approved reason. |

The default is `human-written`.

## Red-Test Gate

Top-level tests should be observed red before production implementation by default. Allowed red-test gate states are:

- `observed-red`
- `not-run-human-approved`
- `already-passing-human-approved`
- `waived`

The implementor treats human-written or human-confirmed top-level tests as read-only unless the confirmed layer todo authorizes a narrow mechanical update.

## TDD Boundary

The top-level layer contracts are human-gated. Internal implementation tests are agent-owned when they stay inside the approved layer scope. If implementation discovers new top-level behavior, the implementor stops and records a checkpoint in the active layer todo.

## Review And Memory

The final-reviewer runs a small review after each layer and a broader review when the slice is complete. It must call out waived layers, risky assumptions, and verification gaps. It may propose memory candidates, but memory capture happens only after human approval.

## Local Files

The installer adds narrow gitignore entries for:

- `.github/plans/`
- `.github/layered-tdd-templates/`

Installed templates are local scaffolding and can be edited during validation. Reinstalling does not overwrite them.

## Promotion Notes

Treat this workflow as validation scaffolding. Promote it only after real tasks show that slice selection, layer gates, red-test discipline, final review, memory proposals, and bash approval friction are worth keeping.
