# Components, Skills & Presets

← [Back to README](../README.md)

---

## Components

| Component | ID | Description |
|-----------|-----|-------------|
| Engram | `engram` | Persistent cross-session memory via MCP — auto-detection of project name, full-text search, git sync, project consolidation. See [engram repo](https://github.com/Gentleman-Programming/engram) |
| Markdown Memory | `markdown-memory` | Plain Markdown memory in an agent-owned vault namespace. Selected with `--memory-backend markdown`; currently injects protocol for OpenCode and Codex. |
| SDD | `sdd` | Spec-Driven Development workflow (9 phases) — the agent handles SDD organically when the task warrants it, or when you ask; you don't need to learn the commands |
| Skills | `skills` | Curated coding skill library |
| Context7 | `context7` | MCP server for live framework/library documentation |
| Persona | `persona` | Managed Gentleman/neutral persona injection, or unmanaged custom persona mode |
| Permissions | `permissions` | Security-first defaults and guardrails |
| GGA | `gga` | Gentleman Guardian Angel — AI provider switcher |
| Theme | `theme` | Gentleman Kanagawa theme overlay |
| Claude Gentleman Theme | `claude-theme` | Claude Code Gentleman custom theme |
| OpenCode Gentle Logo | `opencode-gentle-logo` | OpenCode home logo TUI plugin with Braille rose |
| OpenCode Lean Workflow | `opencode-lean-workflow` | Pilot requirements-first OpenCode workflow with local artifacts |
| OpenCode Layered TDD | `opencode-layered-tdd` | Layer-gated OpenCode TDD workflow with slice selection and human approvals |
| Conductor Layered TDD | `conductor-layered-tdd` | Repo-local Microsoft Conductor workflow for layered TDD with Copilot CLI, project skills, and human gates. See [Conductor Layered TDD](conductor-layered-tdd.md). |

## GGA Behavior

`gentle-ai --component gga` installs/provisions the `gga` binary globally on your machine.

It does **not** run project-level hook setup automatically (`gga init` / `gga install`) because that should be an explicit decision per repository.

After global install, enable GGA per project with:

```bash
gga init
gga install
```

---

## Skills

### Included Skills (installed by gentle-ai)

24 skill files organized by category, embedded in the binary and injected into your agent's configuration:

#### SDD (Spec-Driven Development)

| Skill | ID | Description |
|-------|-----|-------------|
| SDD Init | `sdd-init` | Bootstrap SDD context in a project |
| SDD Explore | `sdd-explore` | Investigate codebase before committing to a change |
| SDD Propose | `sdd-propose` | Create change proposal with intent, scope, approach |
| SDD Spec | `sdd-spec` | Write specifications with requirements and scenarios |
| SDD Design | `sdd-design` | Technical design with architecture decisions |
| SDD Tasks | `sdd-tasks` | Break down a change into implementation tasks |
| SDD Apply | `sdd-apply` | Implement tasks following specs and design |
| SDD Verify | `sdd-verify` | Validate implementation matches specs |
| SDD Archive | `sdd-archive` | Sync delta specs to main specs and archive |
| SDD Onboard | `sdd-onboard` | Guided end-to-end SDD walkthrough on the real codebase |
| Judgment Day | `judgment-day` | Parallel adversarial review — two independent judges review the same target |

#### Foundation

| Skill | ID | Description |
|-------|-----|-------------|
| Go Testing | `go-testing` | Go testing patterns including Bubbletea TUI testing |
| Skill Creator | `skill-creator` | Create new AI agent skills following the Agent Skills spec |
| Branch & PR | `branch-pr` | PR creation workflow with conventional commits, branch naming, and issue-first enforcement |
| Issue Creation | `issue-creation` | Issue filing workflow with bug report and feature request templates |
| Skill Registry | `skill-registry` | Build a compact project standards registry from installed skills |
| Chained PR | `chained-pr` | Plan and create reviewable stacked/chained pull requests |
| Cognitive Doc Design | `cognitive-doc-design` | Write docs that reduce review and onboarding cognitive load |
| Comment Writer | `comment-writer` | Draft warm, direct collaboration comments and review replies |
| Work Unit Commits | `work-unit-commits` | Split implementation into reviewable work units |

These foundation skills are installed by default with both `full-gentleman` and `ecosystem-only` presets.

#### Markdown Memory

These skills are installed when `--memory-backend markdown` is selected.

| Skill | ID | Description |
|-------|-----|-------------|
| Memory Recall | `memory-recall` | Read the Markdown memory hot cache, project index, and targeted project files |
| Memory Capture | `memory-capture` | Stage concise candidate memories during active work |
| Memory Consolidate | `memory-consolidate` | Promote staged observations into canonical memory files |
| Memory Handoff | `memory-handoff` | Maintain task-scoped handoff files |

### Coding Skills (separate repository)

For framework-specific skills (React 19, Angular, TypeScript, Tailwind 4, Zod 4, Playwright, etc.), see [Gentleman-Programming/Gentleman-Skills](https://github.com/Gentleman-Programming/Gentleman-Skills). These are maintained by the community and installed separately by cloning the repo and copying skills to your agent's skills directory.

---

## Presets

| Preset | ID | What's Included |
|--------|-----|-----------------|
| Full Gentleman | `full-gentleman` | All components (Engram + SDD + Skills + Context7 + GGA + Persona + Permissions + Theme) + all skills + gentleman persona |
| Ecosystem Only | `ecosystem-only` | Core components (Engram + SDD + Skills + Context7 + GGA) + all skills + gentleman persona |
| Minimal | `minimal` | Engram + SDD skills only |
| Custom | `custom` | You choose components and skills manually while keeping any existing persona/settings unmanaged |

Use `--memory-backend markdown` to replace Engram with Markdown Memory in presets. Use `--memory-backend none` to skip persistent memory injection entirely.
