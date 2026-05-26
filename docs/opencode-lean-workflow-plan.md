# Plan: Lean OpenCode Workflow With Markdown Memory

**Generated**: 2026-05-21
**Status**: Draft for discussion

This plan compares the current lightweight OpenCode agent workflow with Gentle-AI's broader SDD setup and proposes a smaller hybrid: create an alternate orchestrated workflow, add Markdown memory and safer permissions, and improve requirements discovery with a bounded grill-style clarification phase.

The first deliverable is a validated alternate workflow for `event-catalog-sync`, not a direct replacement of the current agents. Once the alternate workflow proves useful in real work, promote it into a reusable Gentle-AI/OpenCode recipe or installable component.

## Quick Path

1. Build an explicit pilot install surface for the lean OpenCode workflow instead of piggybacking on the existing `sdd` component.
2. Install the pilot workflow with Markdown memory skills and bash-ask permissions. The exact command should follow the component naming chosen during implementation; conceptually it should install:

   ```bash
   gentle-ai install <experimental-lean-opencode-workflow> \
     <markdown-memory-skills> \
     <bash-ask-permissions>
   ```

3. Add the lean workflow as an alternate agent set. Do not directly change the current orchestrator/planner/implementor agents during validation.
4. Treat Gentle-AI SDD as a source of reusable workflow ideas, not necessarily the final user-facing flow.
5. Replace the alternate planner handoff with a requirements-first flow:
   bounded markdown requirements grill -> plan with test inventory -> human acceptance test and confirmation -> todo generation -> implementation.
6. Keep implementation gated by explicit human confirmation.
7. Promote the workflow to a reusable recipe only after the alternate setup has been validated on real tasks.

Do not add a fast path during validation. The pilot should expose the full workflow cost and value before deciding whether a lower-ceremony path is needed.

Focus the pilot workflow on code-changing work first. Docs-only and design-only tasks can use lighter ad hoc flows until validation shows whether this artifact-and-test workflow should grow a separate pattern for them.

## Validation Bar

Validation is owned by the human user. The alternate workflow should be tested on real `event-catalog-sync` tasks before it replaces the current workflow or becomes reusable.

Promotion should depend on value, not speed. A plan is successful when it makes the work clearer, safer, and easier to execute; it does not need to fit an arbitrary review-time budget.

Suggested validation signals:

- the requirements grill catches meaningful ambiguity before planning
- the plan adds useful clarity rather than ceremonial detail
- the todo is operational enough that implementation can proceed from confirmed inputs
- confirmation gates prevent premature implementation
- memory capture preserves reusable context without clutter
- the workflow feels worth using again after real tasks

Token cost matters, but validation should optimize for value first. The workflow should spend context where it improves decisions, implementation quality, or safety, and avoid wasteful rereading by habit.

## Current Workflow

Your current OpenCode setup in `event-catalog-sync` has three agents:

| Agent | Current role | Strength | Limitation |
| ----- | ------------ | -------- | ---------- |
| `agent-orchestrator` | Dispatcher that detects plan/todo state and delegates | Simple, gated, hard to accidentally implement early | Does not improve requirements quality before planning |
| `implementation-planner` | Deep codebase analysis, plan writing, todo writing | Strong exploration mandate and plan template | Plan can become verbose and review-heavy |
| `implementor` | Executes confirmed todo checklist | Clear implementation boundary | Depends heavily on plan/todo quality |

The best parts to preserve:

- orchestrator does not read or edit source code
- planner does analysis but does not implement
- implementor implements only after confirmation
- todo file, not the plan alone, drives implementation
- user confirmation gates the move from plan to todo/implementation

## Gentle-AI Comparison

| Area | Current local workflow | Gentle-AI SDD | Proposed hybrid |
| ---- | ---------------------- | ------------- | --------------- |
| Requirements | Mostly captured from initial user request and plan iteration | Multiple SDD phases can clarify and refine scope | Add a dedicated grill/requirements phase before planning |
| Planning | One rich plan document | Proposal, spec, design, tasks phases | One value-focused plan with explicit test inventory |
| Tasks | Todo generated after confirmation | `sdd-tasks` phase | Keep todo generation after confirmation |
| Implementation | Implementor subagent | `sdd-apply` phase | Keep implementor, but feed it confirmed todo plus tests |
| Memory | None or project-specific ad hoc context | Engram or Markdown Memory | Use Markdown memory skills only |
| Permissions | Current implementor allows bash | Gentle-AI permissions can be broad unless tuned | Ask before every bash command in the alternate workflow |
| Cognitive load | Lower than full SDD | Potentially too much for everyday work | Use SDD internals selectively |

## Proposed Workflow

### Context Budget

Every phase should read and pass forward the context it needs to add value, not every artifact that exists.

Default rules:

- the orchestrator reads phase frontmatter and only enough body content to route work
- subagents receive the minimum phase artifacts needed for their responsibility
- the implementor reads the confirmed plan, todo, listed tests, and relevant code before changing production code
- memory recall is targeted to context that may affect the current phase
- checkpoints summarize current state instead of copying whole prior artifacts

This budget is a relevance rule, not a token-savings-first rule. AI work has a cost, but the workflow should preserve context when it materially improves clarity, quality, or safety.

### Phase 1 - Bounded Requirements Grill

Goal: improve requirements before any implementation plan exists.

This borrows from Matt Pocock's `grill-me` and `grill-with-docs` skills:

- explore code/docs before asking questions
- prefer a bounded markdown question batch over a long chat loop
- ask one chat question only when one answer determines the next branch of discovery
- recommend a default answer with each question
- challenge vague or overloaded domain language
- cross-check user statements against existing code and docs

Default loop:

1. The requirements subagent explores existing code and docs.
2. It writes `01-requirements.md` with current understanding, recommended defaults, blocking questions, non-blocking assumptions, and a place for human feedback.
3. The orchestrator waits for feedback in the artifact or in chat, then writes any chat feedback back into the artifact.
4. The requirements subagent may write one follow-up question batch if feedback opens a meaningful new branch.
5. Further rounds happen only when the human explicitly asks to keep exploring. Otherwise unresolved non-blockers become explicit assumptions.

Questions should be limited to decisions that materially change scope, architecture, safety, testing, or user-visible behavior.

The requirements subagent should inspect docs, code, and targeted memory when that prevents avoidable questions or catches a contradiction in domain language, current behavior, test patterns, or project constraints. Keep that discovery focused; implementation-oriented code inspection belongs to planning.

Output should be concise:

```markdown
## Requirements Summary

- Goal:
- Non-goals:
- Confirmed decisions:
- Open questions:
- Domain terms:
- Important edge cases:
```

The grill phase should stop when either:

- all blocking questions are resolved, or
- remaining unknowns are explicitly accepted as assumptions.

Keep the requirements phase mandatory during validation, even when the request is detailed or an acceptance test already exists. In those cases `01-requirements.md` may be small: record the clear goal, checks performed, confirmed assumptions, and lack of blocking questions, then move to planning.

Require human confirmation of `01-requirements.md` before planning during validation, even when the requirements subagent reports no blocking questions.

### Phase 2 - Plan

Goal: produce a reviewable plan that is valuable and grounded in code without ceremonial detail.

The plan should include:

- request summary
- confirmed requirements
- affected files or modules
- proposed approach
- risks and assumptions
- test inventory
- user confirmation gate

The plan should describe test intent and needed coverage. Concrete verification commands belong in the todo so the implementor's expected command surface is visible before implementation.

Plan confirmation includes approval of the test inventory. Do not add a separate confirmation inside the plan for test inventory approval; the later acceptance-test gate remains separate.

If the planner discovers a material requirements gap while building the plan, it must route back to `01-requirements.md` for revision instead of silently changing confirmed requirements inside `02-plan.md`.

The planner should inspect enough code to name likely affected files or modules. When exact impact is intentionally deferred or still uncertain after focused inspection, record that uncertainty as a plan risk instead of pretending certainty.

The workflow should not impose one acceptance-test location or framework convention across repositories. The planner should discover and recommend the project's existing acceptance-level test pattern; if no credible pattern exists, record that as a plan risk or human decision. The plan does not need to propose a specific acceptance-test file path unless that detail adds value.

Recommended plan shape:

```markdown
# Plan: <Feature>

## Outcome

<One paragraph.>

## Confirmed Requirements

- <Requirement>

## Proposed Change

| Area | Change | Why |
| ---- | ------ | --- |

## Test Inventory

| Type | Test | Purpose | Owner |
| ---- | ---- | ------- | ----- |
| Acceptance | <scenario> | Proves user-visible behavior | Human |
| Integration | <module interaction> | Proves wiring | Agent |
| Unit | <unit behavior> | Proves edge case | Agent |

## Risks

- <Risk and mitigation>

## Confirmation

Reply `confirmed` to generate implementation tasks, or add feedback.
```

Important change from current behavior: the plan should list the tests needed, and the human writes the executable acceptance test before implementation begins unless that gate is explicitly waived.

### Phase 3 - Human Acceptance Test Gate

Goal: make the desired behavior concrete before tasks are generated.

The acceptance layer is human-owned and executable. The acceptance test acts as a specification boundary for stochastic LLM implementation work, so no subagent should implement the acceptance test itself. Prose scenarios and pseudocode may guide the human, but they do not satisfy the acceptance test gate.

The human may explicitly waive the acceptance test gate in the plan artifact when executable acceptance coverage is not appropriate or is disproportionate for the task. The waiver changes workflow routing, so it belongs in `02-plan.md` frontmatter and must not exist only in chat:

```markdown
---
status: confirmed
owner: human
acceptance_test:
  required: waived
  waiver_reason: "Documentation-only change; no executable behavior changes."
  waived_by: human
---
```

When the plan carries a human-approved waiver, todo generation may proceed without a human-authored acceptance test and must carry that waiver into the todo inputs or constraints.

After plan confirmation, the default path is:

1. The human writes or updates the acceptance test.
2. The human confirms that the acceptance test is ready for the workflow to continue and references it with an OpenCode `@` file mention or explicit path.
3. Todo generation reads the confirmed plan, the acceptance test, and any test guidance artifact that exists.

The acceptance test must exist and be read before todo generation unless the human waives it in the plan. It does not have to be executed before todo generation. When its observed status is known, record it as one of:

- `not-run`
- `failing`
- `passing`
- `blocked`

A passing acceptance test should make the todo generator and implementor inspect scope carefully instead of assuming production changes are still required.

Writing the acceptance test may reveal that an upstream artifact needs revision. Route backward based on what changed: revise the todo later for operational task changes, revise the plan when approach or test inventory changes, and reopen requirements when desired behavior changes.

Test guidance is opt-in. A simple plan confirmation should not generate it.

If the human requests test guidance after confirming the plan, trigger a test-planning subagent before todo generation. That subagent should write a markdown guidance artifact with scenarios, likely file paths, pseudocode, edge cases, or lower-level coverage suggestions. It should not write the acceptance test code.

After test guidance is written, the orchestrator must stop again. The human reviews the guidance, writes or updates the acceptance test manually, and confirms that the acceptance test is ready before todo generation continues.

The implementation todo should be based on:

- confirmed plan
- human-authored acceptance test
- required test inventory from the plan
- optional test guidance artifact
- actual lower-level tests already present

Todo generation should identify which lower-level unit or integration tests are still pending for the implementor. It should not treat guidance pseudocode as implemented test code.

### Phase 4 - Todo Generation

Goal: generate a compact task checklist for the implementor.

Todo generation is another human gate. After `04-todo.md` is written, the orchestrator waits for human confirmation before dispatching the implementor.

The human may edit the todo directly, but edits alone do not authorize implementation. The todo must still receive explicit human confirmation after any direct changes or revision feedback.

The todo generator may inspect the acceptance test, existing lower-level tests, and modules named by the confirmed plan to produce operational tasks and verification commands. It must not perform fresh architectural exploration; if the confirmed plan is too vague to become a todo, route back to planning.

The todo should be shorter and more operational than the plan:

```markdown
# Todo: <Feature>

## Inputs

- Plan: `.github/plans/<task-slug>/02-plan.md`
- Acceptance test read: `<path or human-approved waiver from plan>`
- Acceptance test status: `<not-run, failing, passing, blocked, or waived>`
- Test guidance read: `<path or none>`
- Existing lower-level tests read: `<paths>`

## Tasks

- [ ] Add or update pending lower-level tests
- [ ] Implement smallest production change
- [ ] Wire integration points
- [ ] Run verification commands
- [ ] Update memory/handoff if needed
```

The todo should avoid large prose blocks. It should name files and expected verification commands.

After confirmation, the implementor may update todo completion checkboxes and record verification results. It must not rewrite task meaning or expand scope inside the todo; material changes route through a checkpoint or todo revision.

Use `04-todo.md` as the terminal task artifact during validation. When implementation is complete, mark its tasks completed or explicitly skipped and set its frontmatter status to `done`; summarize the outcome in chat and capture memory only when warranted.

### Phase 5 - Implementation

Goal: implement only the confirmed todo.

Implementor rules to preserve:

- read plan and todo first
- read listed tests before changing code
- treat the human-authored acceptance test as read-only by default
- do not add scope not in todo
- keep changes minimal
- run relevant checks
- update checklist as work completes

The implementor may change the acceptance test only when the confirmed todo explicitly authorizes a narrow mechanical update, such as repairing an import after an approved file move. If implementation reveals that the acceptance test is contradictory, impossible, or no longer matches a confirmed requirement, the implementor should stop and write a checkpoint for human review instead of changing the test to fit the implementation.

Implementation-local decisions that still satisfy confirmed scope may stay with the implementor. A discovery that changes behavior, scope, acceptance expectations, or architecture must stop implementation and create `05-checkpoint.md`.

The checkpoint should explain what was attempted, what was discovered, why it changes confirmed work, and the implementor's recommended next route:

- resume implementation from the current todo
- revise the todo while preserving the confirmed plan
- revise the plan and re-evaluate downstream gates
- reopen requirements and re-evaluate downstream artifacts

The checkpoint is the first stop, not an automatic restart. The orchestrator waits for human confirmation, then routes to the recommended or human-selected next subagent.

Keep a single stable `05-checkpoint.md` for active checkpoint state. Once a checkpoint is resolved, fold durable decisions back into the relevant requirements, plan, todo, or memory artifact instead of accumulating numbered checkpoint files.

With the updated permissions policy, every bash command from every alternate-workflow subagent should ask for approval before running. Validation should include the permission friction during discovery, planning, todo generation, and implementation.

The dumb orchestrator follows the same bash approval rule and should normally avoid bash entirely. If it needs shell discovery to route work, treat that as a signal that orchestration logic is becoming too substantive.

### Phase 6 - Memory Capture

Goal: persist useful decisions without creating noise.

Use only the Markdown memory skills:

- `memory-recall` before planning or implementation when context may matter
- `memory-capture` for candidate observations during work
- `memory-consolidate` after a meaningful decision or recurring pattern
- `memory-handoff` when work is paused or transferred

Do not force memory writes for every task or phase boundary. Task artifacts already carry current-work state. Capture memory only when a subagent identifies context worth preserving beyond the task folder, such as a resolved domain term, recurring pattern, non-obvious project constraint, or paused-work handoff.

Memory remains a skill used by the specialized subagent that has the relevant live context. Do not add a dedicated memory subagent during validation.

For the current task, confirmed task artifacts and human feedback take precedence over recalled memory. Current code and tests should validate both. If memory conflicts with a confirmed plan, acceptance test, or user decision, surface the conflict instead of silently treating memory as an override.

## Recommended Agent Shape

Keep the current three-agent workflow unchanged during validation. Add a parallel, alternate workflow centered on a deliberately dumb orchestrator:

- the orchestrator detects workflow state, gates progress, and dispatches
- subagents own task-specific work such as requirements grilling, planning, todo generation, implementation, and memory capture
- markdown artifacts are the shared interface between the orchestrator, subagents, and human feedback

Keep orchestrator edits mechanical. It may create task-folder scaffolding, update frontmatter after human confirmation, copy chat feedback into the relevant artifact feedback section, and record which artifact is awaiting input. It must not rewrite requirements, revise plan content, generate todo tasks, or turn test guidance into new technical decisions.

Use an explicit start contract. `resume <task-slug>` tells the orchestrator to look for that local task folder and continue from its artifact state. Any other coding request starts a new task folder from scratch.

Do not add orchestrator task-listing behavior during validation. The human can inspect local task folders to choose the slug to resume.

Resume routing should follow artifact frontmatter first. Specialized subagents inspect code, tests, or worktree state when their resumed phase needs that technical context.

Normal `resume <task-slug>` should stop when the task todo is already `done`. Continuing from a completed task folder requires explicit reopen intent from the human so a finished local workflow state does not silently become the base for a new change.

Explicit reopen intent should start a new task folder for follow-up work rather than reuse a completed folder. The new task may reference the old slug when useful, but completed artifacts stay as local validation evidence.

| Agent | Responsibility |
| ----- | -------------- |
| `workflow-orchestrator` | Selectable pilot agent that detects state, asks for confirmation, waits while human feedback is pending, and dispatches to subagents |
| `requirements-griller` | Inspect docs/code, write the bounded requirements grill artifact, and revise requirements from human feedback |
| `planner` | Turn confirmed requirements into a plan with test inventory and revise the plan from human feedback |
| `todo-generator` | Turn confirmed plan and test inputs into a compact implementation todo |
| `implementor` | Execute confirmed todo and verification only |

Optional later split:

| Agent | When to add |
| ----- | ----------- |
| `test-guidance-planner` | When the human explicitly requests guidance before writing the acceptance test |

The current three-agent shape remains unchanged. The alternate workflow should use focused subagents where they reduce context breadth and make artifact ownership clearer, then prove that behavior before replacing or generalizing anything.

OpenCode invocation should start by selecting the pilot orchestrator agent. The specialized workflow agents should be installed as subagents for orchestrator dispatch, not as user-selected entrypoints.

During validation, keep the alternate workflow in separate agent files instead of immediately sharing prompt partials with the current workflow. This makes behavioral comparison cleaner. After validation, extract reusable prompt pieces such as confirmation gates, implementation boundaries, todo format, and memory rules.

The alternate implementor should be its own agent file from day one. It may borrow proven constraints from the current implementor, but it must be written for the new artifact contract: confirmed todo as operational input, human-authored acceptance test as read-only by default, checkpoint routing for material surprises, focused context reads, and bash approval under the proposed permissions model.

### Markdown Feedback Gates

When a subagent needs human feedback, it should write a markdown artifact that makes the feedback request explicit. The main orchestrator then stops and waits for the user to edit or answer against that artifact before continuing.

This applies to:

- requirements summaries that need user decisions
- plans that need confirmation or revision
- optional test guidance that needs human acceptance-test input
- todos that need final confirmation before implementation
- memory handoffs or checkpoints when work pauses

The orchestrator should treat the markdown artifact as the source of truth after feedback is provided. Subagents should resume from the confirmed artifact rather than relying on chat history alone.

Use a task folder with separate phase files:

```text
.github/plans/<task-slug>/
  01-requirements.md
  02-plan.md
  03-test-guidance.md
  04-todo.md
  05-checkpoint.md
```

Use `.github/plans/<task-slug>/` as the artifact home during validation. Revisit artifact-root naming only when promoting the workflow beyond the pilot.

Preserve completed task folders locally during validation. Use them to judge which artifacts, gates, checkpoints, and memory interactions were valuable before deciding on archival or cleanup policy for a reusable workflow.

When starting a task folder, derive a conservative slug from an issue or branch identifier when one is available. Otherwise derive it from the user request and let the human correct it if needed.

Keep phase filenames stable across revisions. Use frontmatter and a compact decision log for state and history instead of encoding versions or confirmation state into filenames.

Stable phase filenames are enough for normal task-folder routing; do not add upstream dependency references to every artifact frontmatter. The todo should still list the specific inputs it consumed because it is the implementor contract and may reference tests outside the task folder.

Keep timestamps out of pilot artifact frontmatter. Add timing metadata later only if validation shows a need for stale-task detection or workflow analytics.

Optional phases may leave numeric gaps. For example, a task with no test guidance may move from `02-plan.md` to `04-todo.md`. The prefixes describe the route map; they do not require every task folder to contain every artifact.

Phase files coordinate the current task. They are the reviewable, editable artifacts the human and orchestrator use to move through gates.

Task artifact folders are local workflow state during validation and should be gitignored in the target repository so they do not add noise to unrelated work or reviews. Durable cross-task context should flow through selective memory capture instead of committed task folders.

Do not rely on gitignored task artifacts as PR review context. When review-facing documentation is needed, generate a deliberate PR description or review note from confirmed artifacts and actual code changes instead of exposing internal workflow state.

The human-authored acceptance test is not workflow scratch state. It belongs in the normal test tree and should be committed with the code change as executable specification and regression coverage.

Use compact templates during validation for requirements, plan, optional test guidance, todo, and checkpoint artifacts. Templates should keep frontmatter, feedback inboxes, confirmation language, and decision records consistent while the workflow behavior is being evaluated.

Keep generic artifact templates in this repository outside the agent prompts. Installation should copy them into the target repository as local workflow scaffolding, add the installed template copies to that target repository's gitignore, and have prompts reference the installed template paths instead of embedding large markdown skeletons inline.

Installed template copies should be locally editable during validation. Reinstall or update flows must not overwrite those gitignored template copies unless the human explicitly requests or confirms a template refresh.

Pilot installation should update the target repository gitignore with narrow entries for installed local templates and local task artifact folders.

The workflow needs a dedicated human guide outside this repository's main README. Add a focused pilot-user workflow README beside the future reusable workflow assets and link it from this plan so it evolves with prompts and templates. The guide should explain the happy path, required confirmations, artifact roles, `@` acceptance-test reference at the test gate, optional test guidance branch, checkpoint routing, memory boundary, and local gitignored artifact/template behavior. Keep reusable install guidance provisional in a short promotion-notes section until validation is complete.

Pilot guide asset: [`internal/assets/opencode/lean-workflow/README.md`](../internal/assets/opencode/lean-workflow/README.md).

Memory handoff files have a different job: they preserve context that should survive outside the current task flow, such as why a decision was made, what domain term was resolved, where work paused, or what a future agent needs to know without reading every phase artifact. They should not replace the phase files, and phase files should not become long-term memory dumps.

Each phase file should start with a small frontmatter state block so the orchestrator can route from explicit state instead of inferring progress from prose:

```markdown
---
status: awaiting-feedback
owner: human
---
```

Suggested statuses:

- `draft`
- `awaiting-feedback`
- `confirmed`
- `needs-revision`
- `superseded`
- `done`

The status block should stay small. Use `status` as the single confirmation state instead of a duplicate confirmation boolean. The markdown body remains the place for requirements, decisions, feedback prompts, and task detail.

Use specific workflow agent names in `owner` when an artifact is waiting on agent work, and `human` when the workflow is gated on human input. For example: `requirements-griller`, `planner`, `todo-generator`, `implementor`, or `human`.

Backward routing must invalidate affected downstream artifacts. When an upstream artifact is revised, the orchestrator or responsible subagent should mark stale downstream artifacts as `needs-revision` or `superseded` so implementation cannot proceed from outdated confirmed state.

Human confirmation can arrive either by editing the artifact or by replying in chat. Before dispatching the next subagent, the orchestrator must write that confirmation back into the relevant phase file so markdown remains the workflow source of truth.

Subagents must not self-confirm artifacts that require human approval.

Each phase artifact should include a dedicated human feedback section by default:

```markdown
## Human Feedback

- Decision:
- Requested changes:
- Questions:
```

The orchestrator should copy chat feedback into this section before dispatching the next subagent. The human may still edit an artifact inline when that is clearer, but normal feedback flow should not depend on subagents inferring intent from arbitrary prose changes.

Treat the human feedback section as an active inbox, not a durable transcript. After a subagent consumes feedback, it should preserve important decisions in the artifact body or a compact decision log and clear or reduce consumed working feedback so the artifact stays readable.

## What To Use From Gentle-AI

Use now:

- Markdown memory backend
- memory skills only
- OpenCode permissions with bash set to ask
- SDD assets as references for phase discipline

Use later only if needed:

- full SDD command sequence
- persona injection
- Context7 MCP
- GGA pre-commit review
- broader skill library

Avoid for now:

- default full preset
- Gentle Logo
- Context7 MCP
- GGA
- managed persona
- broad foundation skills

## Implementation Sequence

1. Inspect the existing component and install layout, then choose the pilot asset directory for the lean OpenCode workflow.
2. Add the pilot-user workflow README beside those assets.
3. Add generic templates for requirements, plan, optional test guidance, todo, and checkpoint artifacts.
4. Add separate alternate agent files for the dumb orchestrator and specialized subagents.
5. Configure alternate-workflow permissions so every bash command asks for approval.
6. Add install behavior that copies local template scaffolding into the target repository, references it from prompts, and adds installed templates plus task artifacts to the target gitignore.
7. Install the pilot workflow into `event-catalog-sync` through the intended install path.
8. Validate the installed workflow on code-changing work in `event-catalog-sync`.
9. Revisit fast paths, archival policy, shared prompt extraction, and reusable Gentle-AI/OpenCode promotion after validation.

Expose the pilot through an explicit install component or option during validation. Do not piggyback it on the existing `sdd` install behavior while the workflow is still experimental.

Choose the exact pilot install name during implementation after inspecting existing component naming. The name should clearly signal an experimental lean OpenCode workflow rather than a stable replacement for `sdd`.

Install the required Markdown memory skills with the pilot workflow when the installer can express that dependency cleanly. Separate manual memory-skill guidance is only a fallback if implementation constraints require it.

Configure the Markdown memory backend as part of the pilot install so validation uses one memory path. Keep any required memory project name as an install input if the existing installer needs it.

## Implementation Acceptance Bar

The pilot is ready for installation validation when:

- source assets exist in this repository for the pilot guide, templates, orchestrator, and specialized subagents
- the pilot has an explicit experimental install surface
- installation configures Markdown memory and installs the required memory skills
- installation copies local templates into the target repo without silently overwriting local edits
- installation updates target gitignore entries for local templates and local task artifacts
- OpenCode exposes the pilot orchestrator as the selectable entrypoint and keeps the specialized agents as subagents
- alternate-workflow permissions ask before every bash command
- a test installation into `event-catalog-sync` proves the expected agents, templates, ignores, and permissions are present

## References

- Current local OpenCode agents reviewed from `/Users/cdespona/code/event-catalog-sync/.opencode/agents/`.
- Matt Pocock `grill-me`: https://github.com/mattpocock/skills/tree/main/skills/productivity/grill-me
- Matt Pocock `grill-with-docs`: https://github.com/mattpocock/skills/tree/main/skills/engineering/grill-with-docs
- Gentle-AI local README section: `OpenCode with Markdown memory`
