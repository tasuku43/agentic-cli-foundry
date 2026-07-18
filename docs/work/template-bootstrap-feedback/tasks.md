# Work Tasks: Template bootstrap feedback

## Understand

- [x] Read governing theses, product, architecture, security, harness, public, release, and agent-readiness documents.
- [x] Confirm `profile: template`, clean baseline, and protected defaults blob `5beb6c11f2f92ad4ebbdaf5f44e70eb66224ecb0`.
- [x] Inspect `ReadyProblems`, bootstrap replacement/protection, repository path enumeration, and full-profile composition.
- [x] Reproduce Go selection split and measure template agent-help byte sizes.
- [x] Record verified facts, constraints, unknowns, and thesis evidence.

## Decide

- [x] Select field-level meaningful identity with reusable owner/license.
- [x] Select Git-success-only path enumeration with absent deletion filtering.
- [x] Select protected-defaults documentation rather than a broad replacement exclusion.
- [x] Preserve schema-v1 `ready` and recommend a reviewed bootstrap/product profile split.
- [x] Preserve complete agent contracts and defer schema/budget changes until an evaluation corpus exists.

## Implement

- [x] Add and run failing project identity regressions.
- [x] Add and run failing temporary-Git rename regression.
- [x] Implement `ReadyProblems` and repository path fixes.
- [x] Add bootstrap documentation consistency test and update durable prose.
- [x] Add early Go/ShellCheck preflight and regression fixtures.
- [x] Clarify identity-only readiness and full-profile prerequisites in durable documentation.
- [x] Record scoped-help measurements and follow-up evaluation in agent-readiness documentation.

## Verify

- [x] Focused Go and shell tests pass. Evidence: three focused Go packages plus sanitized-environment and mixed-toolchain fixtures.
- [x] `task check` passes. Evidence: final full run includes security, release, and public; vulnerability result was clean.
- [x] `task public:check` passes. Evidence: independent `repoguard (public): OK` and `contractlint: OK`.
- [x] Template `.harness/project.json` still reports `template`. Evidence: `projectmeta --field profile`.
- [x] `tools/internal/projectconfig/defaults.go` blob remains `5beb6c11f2f92ad4ebbdaf5f44e70eb66224ecb0`. Evidence: `git hash-object`.
- [x] Worktree diff contains no unrelated changes. Evidence: final status/diff review is limited to this packet, harness implementation, tests, Skills, and durable documentation.

## Hand off

- [x] Acceptance criteria have evidence.
- [x] Implemented decisions are promoted out of the work packet.
- [x] Unresolved product-profile and scoped-help schema work is explicit.
- [x] Final report includes changes, reproduction evidence, open decisions, and gate results.
