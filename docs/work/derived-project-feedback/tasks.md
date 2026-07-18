# Work Tasks: Integrate proven derived-project lessons

- Goal: [goal.md](goal.md)
- Plan: [plan.md](plan.md)

## Understand

- [x] Read governing theses, product, architecture, security, harness, Skill, and
  agent-readiness sections. Evidence: repository review at `b37ec8369b2c`.
- [x] Reproduce or observe current behavior. Evidence: exact Git-tree lineage,
  scoped-help measurements, and catalog/work-packet inspection in `context.md`.
- [x] Record verified facts, unknowns, repeated friction, and thesis evidence.
- [x] Confirm the public outcome and non-goals in `goal.md`.

## Decide

- [x] Compare selective promotion, wholesale copy, prose-only, and full evaluator
  alternatives in `plan.md`.
- [x] Identify public-contract and compatibility impact: agent-help schema bump.
- [x] Confirm command roles, reference flow, capabilities, effects, authentication,
  and side-effect behavior do not change.
- [x] Confirm no external API or new dependency is involved.
- [x] Record that no ADR is required unless implementation reveals a new mature
  thesis conflict.

## Implement

- [x] Promote operational closure and semantic-result rules to durable docs and
  `$add-capability`.
- [x] Add strict operation enum text/JSON round-trip tests and implementation.
- [x] Add exact single-literal argv grammar tests and implementation.
- [x] Add catalog-wide fault signature tests and validation.
- [x] Split output delivery from collection coverage and migrate built-in and
  synthetic catalog contracts.
- [x] Version and implement grouped agent-help workflows.
- [x] Add derived-scale help edge-equivalence and whole-response budget tests.
- [x] Add work-packet completion/supersession validation and update templates.
- [x] Withdraw the provisional template-provenance/schema-2 implementation and
  public claims; retain schema 1 and record the independent combined follow-up
  contract for root evidence, migration, Git binding, preview/apply binding,
  fresh-bootstrap public E2E, and the pre-existing backup-cleanup ambiguity.
- [x] Add lightweight presentation-evidence and capability-retirement guidance.
- [x] Add task-owned sample evidence for successful empty JSON output,
  same-label distinct identity, exact selected-ID reads, and target-mismatch
  rejection before presentation.

## Verify

- [x] Focused tests pass. Evidence: schema-1/project/bootstrap and work-packet
  packages pass with `go test ./tools/internal/projectconfig
  ./tools/projectmeta ./tools/bootstrap ./tools/repoguard`; focused `gofmt -l`,
  `git diff --check`, and `go run ./tools/repoguard --scope hygiene` are clean
  with writable isolated Go caches.
- [x] `task check` passes. Evidence: the Go 1.26.5 full profile passes hygiene,
  architecture/catalog contracts, unit/race/vet/tidy/diff checks, gosec,
  govulncheck, release lint, actionlint, public guard, and contract lint.
- [x] Manual root/sample agent-help probes match the versioned contract.
  Evidence: schema 5 root and `sample` output are 1,517 and 8,222 bytes; the
  grouped synthetic fixture is 24,493 bytes for 324 implicit edges.
- [x] Generated diff and repository status are understood. Evidence: 25 tracked
  files are intentionally modified and six Markdown files are intentionally
  new; no generated, credential, build, or temporary artifact is present.

## Hand off

- [x] Acceptance criteria have evidence and checkboxes are complete.
- [x] Durable decisions are promoted from this packet.
- [x] Temporary diagnostics and sensitive artifacts are absent.
- [x] Follow-up work is explicit and non-blocking.
- [x] Handoff summary explains outcome, why, checks, compatibility, and risks.
