# Work Goal: Integrate proven derived-project lessons

- Status: Complete
- Successor: None
- Owner: Maintainers
- Target: Current change
- Related ADRs: None

## Outcome

Agentic CLI Foundry incorporates the provider-neutral contracts that were
validated by a substantial derived CLI: supported outcomes are operationally
closed, semantic results are bound to their request before presentation,
agent-help workflows remain bounded at a derived catalog scale, generic catalog
values round-trip safely, and repository work-packet claims are mechanically
consistent.

## Why now

The evaluated derived project began from the exact Foundry tree and then
exercised the template through dozens of typed external tasks. Its
implementation, removal, and presentation experiments exposed repeatable gaps
in the template contracts,
including scoped-help workflow multiplication and a completed work packet whose
acceptance boxes remained unchecked. The same project also proved that several
existing Foundry boundaries scale without replacement.

## Non-goals

- Import provider-specific commands, identifiers, authentication choices,
  limits, wire types, notation, or presentation grammar.
- Make exhaustive upstream-operation coverage mandatory for derived projects.
- Add a live-model evaluator, provider SDK, network call, credential flow, or
  third-party dependency.
- Replace the catalog, four-layer architecture, current authentication
  contracts, tool-local fixed targets, or the canonical check script.
- Bootstrap this template repository into a derived product.
- Add passive input-type metadata before a catalog-owned parser can enforce it;
  boolean, integer, repeatable, and multiword inputs need a separate executable
  parser slice with representative commands.
- Add a live upstream-snapshot or evaluator contract before a second derived
  product supplies independent evidence for that abstraction.
- Introduce template-provenance metadata or project schema 2 in this packet.
  That change needs an independent design and migration packet that satisfies
  all provenance, Git-boundary, preview/apply, and fresh-bootstrap checks
  together.
- Change the existing bootstrap backup-cleanup behavior. Its post-commit error
  ambiguity is recorded for the same follow-up rather than mixed into this
  contract slice.

## Acceptance criteria

- [x] Governing documents and `$add-capability` define operational closure for
  supported outcomes and keep raw/export utilities explicitly outside that
  promise when a derived product chooses them.
- [x] Capability guidance requires request-bound semantic results, contextual
  reference-kind validation, empty-collection scope preservation, and negative
  inference tests before presentation.
- [x] Agent-help uses a versioned grouped workflow representation whose encoded
  size is bounded against a derived-scale synthetic catalog without losing
  producer, consumer, invocation, or recovery facts.
- [x] Generic operation enums support strict text/JSON round trips; catalog argv
  grammar supports an exact single literal; and one fault code cannot publish
  conflicting signatures across commands.
- [x] Output metadata separates complete-versus-paged delivery from exhaustive,
  bounded-window, differential-window, or non-collection coverage so a fully
  delivered finite window cannot be mistaken for all provider history.
- [x] A repository check rejects internally inconsistent completed work packets
  and the work templates record supersession, closure, and evidence metadata.
- [x] Project configuration remains on schema 1, provisional template-
  provenance metadata and public-guard exceptions are absent, and the complete
  independent follow-up contract is recorded without claiming implementation.
- [x] No provider-specific or confidential content enters the public template.
- [x] `task check` passes and the required security/public profiles are covered
  by that canonical full gate.

## Governing documents

- Thesis: 2, 5, 6, and 7
- Product contract section: Default supported outcomes, product rules, and
  derived-project completion checklist
- Architecture or security invariant: Catalog source of truth, opaque reference
  boundary, output and terminal safety, claims-to-checks discipline
- Existing ADR: None

## Completion definition

The work is complete when every acceptance criterion has direct repository
evidence, public schema changes are versioned and tested, the new lints run
through the canonical gate, durable decisions are promoted to numbered docs and
the capability Skill, and the working tree contains no unexplained generated or
temporary artifacts.
