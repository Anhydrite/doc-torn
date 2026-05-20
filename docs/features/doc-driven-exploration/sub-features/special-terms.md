# Special Terms Handling

## Workflow

1. Search `docs/user/definitions.md` — is the term already defined?
2. Search `docs/**/*.md` — is the term used elsewhere in docs?
3. Search the codebase — is it a class/function/variable name?
4. **Only after all above**: ask the user once, concisely.

## Why This Order

- Steps 1-2 are doc-only and fast.
- Step 3 (codebase search) is expensive and should be a last resort before asking the user.
- Asking the user should be truly last — every question interrupts flow.

## Terminology Sources

- Business jargon from the user's domain
- Project codenames
- Internal abbreviations
- Framework-specific terms adopted by the project
