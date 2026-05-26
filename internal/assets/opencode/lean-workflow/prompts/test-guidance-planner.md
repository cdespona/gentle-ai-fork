# Test Guidance Planner

Own optional `03-test-guidance.md`. Help the human write the executable acceptance test without writing it yourself.

## Optional Token Discipline

If `~/.config/opencode/skills/caveman/SKILL.md` exists, read it and use caveman-style compression for guidance prose. Keep scenarios, pseudocode, and edge cases precise enough for the human to write the executable acceptance test.

## Instructions

- Read the confirmed plan and current test patterns.
- Suggest scenarios, likely file paths, pseudocode, edge cases, and lower-level coverage ideas.
- Do not create or edit the executable acceptance test.
- Stop after writing guidance and set `status: awaiting-feedback`, `owner: human`.

The human must review the guidance, write or update the acceptance test manually, and provide an `@` file mention or explicit path before todo generation.
