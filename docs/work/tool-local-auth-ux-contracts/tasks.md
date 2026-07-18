# Work Tasks: Generalize target binding and authentication setup contracts

## Understand

- [x] Read governing theses, product, architecture, security, harness, authentication, external API, readiness, ADR, and Skill documents.
- [x] Reproduce the current universal required-reference validation from code and tests.
- [x] Record derived-project evidence, constraints, and unknowns.

## Decide

- [x] Choose exclusive reference-bound or command-bound target binding.
- [x] Keep provider/flow/browser/callback/credential choices out of template core.
- [x] Choose schema-versioned bounded configuration, deterministic precedence, and fail-closed invalid behavior.
- [x] Advance scoped agent-help schema for the new public field.
- [x] Record thesis and architecture changes before mechanism implementation.

## Implement

- [x] Promote durable documentation and Skill guidance.
- [x] Add fixed-target type, validation, help/JSON projection, and deep copy.
- [x] Add every requested negative catalog test and existing-workflow regression coverage.
- [x] Add non-secret configuration types and application resolver.
- [x] Add strict codec, atomic file store, status, fixtures, and filesystem tests.
- [x] Add human-handoff readiness scorecard and OAuth UX checklist.

## Verify

- [x] Focused tests pass. Evidence: `go test ./internal/domain/authn ./internal/app/authn ./internal/infra/authconfigfile ./internal/cli`
- [x] `task check` passes. Evidence: Go 1.26.5, 2026-07-18.
- [x] `task security` passes. Evidence: module verification, repoguard security, and govulncheck reported no vulnerabilities.
- [x] `task public:check` passes. Evidence: repoguard public and contractlint succeeded.
- [x] Repository status and commits are understood. Evidence: contract, fixed-target mechanism, and authentication configuration are separate commits; final evidence is a documentation-only commit.

## Hand off

- [x] Acceptance criteria have evidence.
- [x] Durable decisions were promoted out of the work packet.
- [x] Follow-up provider-specific choices remain explicit.
- [x] Summary explains outcome, checks, compatibility, and risks.
