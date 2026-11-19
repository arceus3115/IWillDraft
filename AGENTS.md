## Project layout
- `app/` contains routes and server-rendered pages.
- `components/` contains reusable UI components.
- `lib/` contains shared business logic and integrations.
- Tests live next to the code they cover as `*.test.*` 
 
## Coding conventions
- Prefer small, focused modules with explicit types.
- Keep business logic in pure functions where possible.
- Match existing component and file naming patterns.
- Do not add a new dependency unless the existing stack cannot solve the problem.
- Do not write fallback code, if tests/functionality fails re-approach the problem.
 
## Verification
- Before calling work finished, run verification test.
- If a command fails, fix the issue and rerun it.
