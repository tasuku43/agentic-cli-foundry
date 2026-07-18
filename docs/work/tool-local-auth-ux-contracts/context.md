# Work Context: Generalize target binding and authentication setup contracts

## Current behavior

- `validateCommandReferenceRole` rejects every `RoleAct` with no consumed opaque reference, even when the command path already identifies the sole tool-owned target.
- `validateAgentContract` requires reference-bound create/write mutation inputs and has no command-bound branch.
- `authn.Requirement`, `authn.Session`, and `app/authn.Gate` intentionally exclude credentials, but no typed persistent non-secret user-configuration boundary exists.
- Agent-readiness counts discovery invocations and notes authentication ceremonies separately, but does not record their human handoffs.

## Relevant structure

- Domain rule: `internal/domain/authn`
- Application use case: `internal/app/authn`
- Infrastructure boundary: new provider-neutral non-secret file store only; provider authentication remains derived
- CLI catalog or presentation: `internal/cli/catalog.go` and `internal/cli/help.go`
- Existing tests and harness checks: catalog, help, authn, architecture, security, and public profiles

## Constraints

- A fixed target is allowed only when the command path identifies exactly one CLI-owned local object; external or ambiguous targets stay reference-bound.
- Fixed and reference binding are mutually exclusive and fixed-target acts produce and consume no opaque references.
- Credentials never enter the new configuration types or plaintext configuration file.
- Invalid higher-priority configuration fails closed instead of falling back.
- No OAuth protocol or platform process implementation is added to the template.

## External facts

No external source is required. The evidence is a derived-project implementation result supplied by a maintainer and current repository behavior reproduced from local catalog validation.

## Unknowns

- [ ] A derived project must decide its concrete configuration fields, provider registration, account model, browser policy, callback listener, and credential store.
- [ ] A future public auth status/setup command must define its own catalog output and fault contract.

## Thesis evidence

- Repeated design decision or point of agent confusion: a universal discover-to-act rule forces an agent to select a target when no selection exists.
- User outcome or friction observed in the minimal slice: local authentication state needed a ceremonial reference/fixed argument before it could be read or changed.
- Code workaround or exception being considered: weakening `RoleAct` reference validation locally.
- Proposed thesis revision: every action has exactly one explicit target-binding mode—opaque reference or catalog-declared tool-local singleton.
- Downstream impact: theses, product, architecture, security, Skill, catalog schema/validation, mutation rules, harness evidence, and agent-readiness.

## Reproduction or observation

```sh
go test ./internal/cli -run TestCatalogRejectsActWithoutRequiredReference
```

The existing test and validator reject a no-reference act. The new regression suite will retain that rejection unless a complete fixed target is declared.

## Security and public-boundary notes

- Assets and side effects involved: tool-local state metadata and a non-secret user-configuration file.
- Credentials involved: none; tokens, PATs, refresh tokens, authorization codes, PKCE verifiers, and client secrets remain infrastructure-private.
- New dependencies or destinations: none.
- OAuth process/browser integration remains a derived, bounded infrastructure decision with manual URL fallback.

## Glossary

- **reference-bound target**: a target supplied as a required opaque input produced by a catalog workflow.
- **fixed target**: a stable catalog declaration identifying one command-bound object.
- **tool-local**: owned by this CLI installation and uniquely selected by the command path, not by external discovery or caller input.
- **user authentication configuration**: persistable non-secret setup metadata; never a credential.
