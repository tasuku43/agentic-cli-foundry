# Work Plan: Generalize target binding and authentication setup contracts

- Status: Implemented
- Goal: [goal.md](goal.md)
- Context: [context.md](context.md)

## Chosen approach

Generalize actions to exactly one target-binding mode. Ordinary actions remain reference-bound. A `RoleAct` may instead carry one `FixedTarget` whose only supported scope is `tool_local`; this declaration is complete, stable, reference-free, and projected through scoped agent help. Mutation validation branches only after the common safety contract is checked, leaving the established reference-bound branch intact.

Add a provider-neutral `authn.UserConfiguration` vocabulary, an application resolver with environment-over-persistent precedence and no invalid-value fallback, and a reusable infrastructure file store. The store owns strict bounded JSON decoding, schema rejection, safe permissions, symlink/non-regular rejection, atomic replacement, and read-only status. It stores no credentials and does not perform authentication.

## Alternatives considered

### Model the singleton as a synthetic opaque reference

Rejected because it creates meaningless discovery or ceremonial fixed input, misstates who chooses the target, and increases agent/human work without adding target certainty.

### Document configuration storage without executable support

Rejected because strict decoding, atomic replacement, and filesystem rejection are security claims that need types and tests. A generic injected-path store is reusable without selecting a provider.

### Add concrete OAuth setup commands

Rejected because grant, browser, callback, registration, credential storage, and account policy belong to a derived product and its reviewed security model.

## Design

### Public contract

Scoped agent-help schema advances to version 4 because `AgentContract` gains an optional `fixed_target`. Root index shape remains compact. A fixed-target act declares no produced or consumed references and therefore creates no reference workflow edge.

### Layer changes

- Domain: bounded, schema-versioned secret-free user authentication configuration and status vocabulary.
- Application: deterministic two-source resolver; invalid or corrupt higher-priority values stop resolution.
- Infrastructure: injected-path non-secret JSON store with strict codec and atomic file semantics.
- CLI and catalog: `FixedTarget`, validation, clone, JSON/help projection, and mutation binding branch.

### Data and control flow

Catalog validation establishes either reference binding or command binding before dispatch. Runtime mutation code still constructs and validates `operation.Intent`, `TargetRef`, and `Impact` before policy and side effects.

Environment configuration, when present, is decoded and validated first. Absence permits persistent configuration. Any present invalid source fails closed. The selected `Method` is returned exactly once; the downstream authenticator may not probe another method after failure.

### Error and cancellation behavior

Configuration errors contain stable classifications and no rejected value. Corrupt, unsafe, or unknown-schema files are invalid, not absent. Save failure leaves the previous target intact when atomic rename has not completed. Status only observes and never repairs.

### Security and public boundary

The file contains public setup metadata only, uses owner-only permissions, and is distinct from every credential store. Browser launch and callback receipt are separate optional infrastructure boundaries. A launcher may receive an authorization URL in argv only after documenting process-list/log exposure; authorization codes, PKCE verifiers, tokens, and secrets never enter argv.

## Implementation slices

1. Promote thesis, architecture, security, authentication, readiness, harness, and Skill decisions.
2. Add fixed-target catalog type, schema projection, validation, clone, and regression tests.
3. Add non-secret authentication configuration vocabulary, resolver, file store, fixtures, and tests.
4. Run focused and repository gates; record evidence.

## Verification

- Unit and contract tests: fixed target, deep copy, JSON projection, config validation/resolution/codec/store.
- Negative side-effect tests: catalog rejection precedes dispatch; unsafe config never replaces/loads.
- Opaque-reference tests: existing sample and reference-bound mutation suite unchanged and passing.
- Agent-readiness: document and exercise target-binding certainty and human-handoff scorecard.
- Required profiles: `task check`, `task security`, `task public:check`.

## Rollout and rollback

The optional field is additive but intentionally versioned as agent-help schema 4. Existing catalogs require no migration. Removing fixed targets returns the catalog to reference-only behavior. Configuration files are opt-in; no default path or state is created.

## Documentation promotion

Promote all target-binding, non-secret configuration, OAuth UX-boundary, and human-handoff decisions into numbered durable documents, `AGENTS.md`, and `$add-capability` before completion.
