# Work Goal: Generalize target binding and authentication setup contracts

- Status: Complete
- Owner: repository maintainers
- Target: template feedback from a derived authentication CLI
- Related ADRs: [ADR 0001](../../decisions/0001-oauth-library-boundary.md)

## Outcome

Derived CLIs can model an action against a command-bound, tool-local singleton without inventing an opaque reference, can persist bounded non-secret authentication configuration without mixing it with credentials, and can compare authentication UX designs using explicit human-handoff measurements.

## Why now

A real derived authentication design exposed three template gaps: the universal action-reference rule created ceremonial discovery and fixed arguments for local singleton state; the authentication contract had no safe home for persistable public setup values; and agent-readiness counted CLI discovery while leaving human/browser handoffs invisible.

## Non-goals

- Add a provider, OAuth grant, browser launcher, callback server, credential store, or public authentication command.
- Relax opaque-reference rules for external, ambiguous, or caller-selected targets.
- Weaken mutation intent, impact, failure, policy, or post-outcome guarantees.
- Treat fewer human steps as automatically safer or better.

## Acceptance criteria

- [x] `AgentContract` represents one validated `tool_local` fixed target and scoped agent help publishes its stable kind, ID, scope, and description.
- [x] Catalog validation rejects incomplete, role-inconsistent, reference-bearing, or mutation-inconsistent fixed-target declarations.
- [x] Existing reference-bound action and mutation validation remains covered by regression tests.
- [x] Provider-neutral, schema-versioned non-secret authentication configuration has bounded strict decoding, fail-closed resolution, atomic file replacement, file-safety checks, and read-only status.
- [x] Authentication and security guidance separates public user configuration from credentials and records safe OAuth browser/callback extension points.
- [x] Agent-readiness records human handoffs and ceremonial inputs as comparison evidence.
- [x] `task check`, `task security`, and `task public:check` pass.

## Governing documents

- Thesis: target binding in Thesis 3; side effects in Thesis 4
- Product contract section: command roles and derived-project decisions
- Architecture or security invariant: catalog target binding; authentication and controlled boundaries
- Existing ADR: OAuth library boundary

## Completion definition

The work is complete when the executable catalog and authentication contracts, negative tests, durable documents, Skill guidance, and required gates agree, without adding a provider or weakening reference-bound workflows.
