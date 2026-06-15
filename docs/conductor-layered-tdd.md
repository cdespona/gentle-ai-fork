# Conductor Layered TDD

This workflow ports the layered TDD flow to Microsoft Conductor so a human can run one selected slice through deterministic gates: requirements, slice selection, layer mapping, layer todo approval, red-test gate, implementation, verification, review, and final memory approval.

Use this when you want the layered TDD workflow outside OpenCode, especially from GitHub Copilot CLI.

## Quick Path

1. Install Conductor and authenticate GitHub Copilot.
2. Install the repo-local layered TDD workflow with Gentle AI.
3. Start the workflow from the repository root.
4. Use the web dashboard to answer each human gate.
5. Review the artifacts under `.github/plans/<slug>/`.
6. Approve each layer only after the red-test gate is acceptable.

```bash
curl -sSfL https://aka.ms/conductor/install.sh | sh
gh auth login
gh skill install microsoft/conductor conductor
gentle-ai install --component conductor-layered-tdd
conductor run workflows/conductor/layered-tdd.yaml \
  --workspace-instructions \
  --web \
  --input request="Describe the code-changing task"
```

Use Markdown memory during requirements recall and final capture by selecting the Markdown backend at install time:

```bash
gentle-ai install \
  --component conductor-layered-tdd \
  --memory-backend markdown \
  --memory-project event-catalog-sync
```

For this Conductor workflow, memory backend choices are intentionally limited to `markdown` or `none`. Engram is not supported by this workflow.

## What Gentle AI Installs

`gentle-ai install --component conductor-layered-tdd` writes only local repository files. It does not download Conductor, install the GitHub Copilot CLI skill, or refresh this workflow from Microsoft/GitHub.

| Path | Purpose | Git behavior |
| --- | --- | --- |
| `workflows/conductor/layered-tdd.yaml` | The runnable Conductor workflow. | Added to `.gitignore`. |
| `workflows/conductor/prompts/*.md` | Specialist prompts used by the workflow. | Added to `.gitignore` with the workflow folder. |
| `.github/skills/conductor-*/` | Repo-local skill copies used only by this Conductor workflow. | Added to `.gitignore` so workflow support files do not pollute PRs. |
| `.github/plans/` | Runtime task artifacts written by workflow runs. | Added to `.gitignore`. |

Installed workflow and prompt files are copied from the `gentle-ai` binary. Re-running install or sync preserves local workflow, prompt, and skill edits unless the file is missing.

The Conductor base skills are always installed. The `conductor-memory-*` skills are installed only when the install selected `--memory-backend markdown`.

The installer does not create or edit `.github/copilot-instructions.md`. That file is repository-specific project guidance. If a repository already has one, `conductor run --workspace-instructions` lets Conductor/Copilot read it; otherwise the layered TDD workflow still works from its own workflow instructions and prompts.

The installed repo-local skills are:

| Skill | Used for |
| --- | --- |
| `conductor-tdd` | Red/green/refactor discipline inside approved layer boundaries. |
| `conductor-test-type-classification` | Choosing the right test level and ownership. |
| `conductor-cognitive-doc-design` | Human-readable artifacts and docs changes. |
| `conductor-work-unit-commits` | Keeping slice/layer work reviewable. |
| `conductor-comment-writer` | Human-facing gate summaries and notes. |
| `conductor-memory-recall` | Reading targeted Markdown memory during requirements discovery. Installed only with `--memory-backend markdown`. |
| `conductor-memory-capture` | Staging approved Markdown memory candidates. Installed only with `--memory-backend markdown`. |
| `conductor-memory-consolidate` | Promoting staged Markdown memory into project memory files. Installed only with `--memory-backend markdown`. |

To restore missing local assets later, run:

```bash
gentle-ai sync --component conductor-layered-tdd
```

Sync uses the embedded Gentle AI assets only; it does not upgrade this workflow from a remote repository.

The default command contract is Makefile-based so every repository can decide
what each check means:

| Input | Default target |
| --- | --- |
| `test_command` | `make test` |
| `lint_command` | `make lint` |
| `security_command` | `make audit` |
| `memory_vault` | none |
| `memory_namespace` | `machine/agent-memory` |
| `memory_project` | none |

Override any target when a repository uses different names:

```bash
conductor run workflows/conductor/layered-tdd.yaml \
  --workspace-instructions \
  --web \
  --input request="Describe the code-changing task" \
  --input test_command="make test" \
  --input lint_command="make check" \
  --input security_command="make cve"
```

Pass Markdown memory inputs when you want the final memory gate to write to your PKM:

```bash
conductor run workflows/conductor/layered-tdd.yaml \
  --workspace-instructions \
  --web \
  --input request="Describe the code-changing task" \
  --input memory_vault="/absolute/path/to/vault" \
  --input memory_namespace="machine/agent-memory" \
  --input memory_project="event-catalog-sync"
```

If the Markdown memory structure already exists, the workflow can recall from it during requirements grilling and capture into it at the final memory gate. If `memory_vault` or `memory_project` is omitted, recall/capture are skipped and the final review records that memory was unavailable.

If the vault root exists but the project memory folder does not, final capture creates only a minimal staging file under:

```text
<memory_vault>/<memory_namespace>/projects/<memory_project>/inbox/staged-observations.md
```

It does not initialize canonical memory files such as `index.md`, `current-state.md`, `decisions.md`, `architecture.md`, `risks.md`, or `open-questions.md`. Full Markdown memory initialization still belongs to the Markdown memory component, for example via an install/sync that uses `--memory-backend markdown`.

## What The Human Does

| Gate | Human decision |
| --- | --- |
| Slice selection | If multiple valuable slices exist, select exactly one. |
| Requirements | Confirm one slice goal, blockers, assumptions, and out-of-scope work. |
| Layer selection | Choose the next layer after reading `01-layer-map.md`. |
| Layer todo | Approve Gherkin, test ownership, and red-test gate state. |
| Checkpoint | Decide what to do if implementation discovers new top-level behavior. |
| Layer approval | Approve the layer result after tests, lint, and audit/CVE checks, request fixes, or move to final review. |
| Memory | Approve, revise, or skip final memory candidates. Approve routes to a Markdown memory capture agent before completion. |

The workflow is intentionally human-gated at top-level behavior. The agent owns internal TDD details only inside the approved layer boundary.

Markdown memory is used in two places when `memory_vault` and `memory_project` are provided. The requirements griller loads `conductor-memory-recall` to read targeted current-state, decisions, risks, open questions, and relevant handoffs before grilling. If the memory root or project folder does not exist, recall is skipped. Current repository files, command output, the current request, and human edits in workflow artifacts always override memory; conflicting or likely superseded memory is recorded as a risk/assumption, not followed blindly.

Memory capture is also Markdown-only. The Conductor workflow does not use Engram or MCP `mem_save`; it loads `.github/skills/conductor-memory-capture/SKILL.md` and `.github/skills/conductor-memory-consolidate/SKILL.md` when the human chooses capture. If `memory_vault`, `memory_project`, or the vault root is missing, the capture agent records that memory capture was skipped instead of writing to an guessed location.

If you edit or comment in `.github/plans/<slice-slug>/00-requirements.md`, choose the requirements gate option that reruns the griller. The workflow routes back to `requirements_griller`, which rereads the artifact and revises it before layer mapping starts.

## Artifacts

| Path | Purpose |
| --- | --- |
| `.github/plans/<discovery-slug>/slice-selection.md` | Written only when the request contains multiple independently valuable slices. |
| `.github/plans/<slice-slug>/00-requirements.md` | Selected slice goal, blockers, assumptions, and out-of-scope slices. |
| `.github/plans/<slice-slug>/01-layer-map.md` | Real repo layers, recommended order, and skeleton layer todos. |
| `.github/plans/<slice-slug>/layers/<nn>-<layer>.todo.md` | One layer's Gherkin, test ownership, red-test gate, scope, tasks, and review notes. |
| `.github/plans/<slice-slug>/99-final-review.md` | Slice summary, waived gates, verification gaps, risks, and memory candidates. |

These artifacts are local working files. The Conductor layered TDD installer adds `.github/plans/` to the target repository's `.gitignore`.

## Preflight Checks

The first workflow step runs `test_command`, which defaults to `make test`.
Requirements analysis does not start while the existing test suite is failing.

If tests pass, the workflow runs `lint_command` and `security_command`, defaulting
to `make lint` and `make audit`. A failure opens a human gate with three choices:

- retry after fixing the repository
- continue with an explicit waiver
- stop before analysis

The same three checks run again after every implementation pass before layer
review.

## Status Schema

Use simple frontmatter so Conductor agents and humans can inspect state without guessing.

| Artifact | Allowed `status` values |
| --- | --- |
| `slice-selection.md` | `slice-selection-required`, `selected`, `stopped` |
| `00-requirements.md` | `needs-human-confirmation`, `blocked`, `confirmed` |
| `01-layer-map.md` | `needs-human-layer-selection`, `confirmed`, `revising` |
| `layers/*.todo.md` | `skeleton`, `needs-human-test-gate`, `ready-for-implementation`, `in-progress`, `checkpoint`, `reviewed`, `approved` |
| `99-final-review.md` | `needs-human-memory-decision`, `complete` |

Each artifact should also include:

```yaml
---
workflow: layered-tdd
owner: human
status: needs-human-confirmation
---
```

## Red-Test Gate

Implementation can proceed only when the selected layer todo records one of these gate states:

| State | Meaning |
| --- | --- |
| `observed-red` | The top-level test failed before production code. |
| `not-run-human-approved` | The human approved proceeding without observing the run. |
| `already-passing-human-approved` | The human approved because the behavior was already covered/passing. |
| `waived` | The human explicitly waived the top-level red-test gate. |

Any other state blocks production implementation and routes back to the layer todo.

## Test Ownership

| Mode | Meaning |
| --- | --- |
| `human-written` | The human writes top-level red tests, then confirms the gate. |
| `agent-written-after-approval` | The human approves Gherkin first; the agent writes top-level tests and stops for confirmation. |
| `waived` | No top-level red test for the layer; the todo must record a human-approved reason. |

Top-level tests are read-only for the implementor by default. Internal tests may be added by the agent under TDD when they stay inside the approved layer boundary.

## Resume

Use `resume_path` to point the workflow at an existing task folder:

```bash
conductor run workflows/conductor/layered-tdd.yaml \
  --workspace-instructions \
  --web \
  --input request="Resume the active layered TDD slice" \
  --input resume_path=".github/plans/<slice-slug>"
```

## Validate Before Running

```bash
conductor validate workflows/conductor/layered-tdd.yaml
conductor run workflows/conductor/layered-tdd.yaml \
  --dry-run \
  --input request="Smoke test the workflow"
```

## References

- [Layered TDD grill notes](lean-workflow-v2-grill-notes.md)
- [Microsoft Conductor](https://github.com/microsoft/conductor)
