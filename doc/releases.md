# Release Remediation SOP

This repository publishes releases from `.github/workflows/release.yml`, with the release tag derived from
`https://api.nuon.co/version`.

## Invariant

For a release tag `vX.Y.Z`, the tagged commit must embed `spec/doc.json` with `info.version == X.Y.Z`.

If this invariant is broken, binaries can be published under the right tag while containing an older API spec.

## Common Failure Modes

1. Tag exists but release does not.
2. GoReleaser fails with `git tag ... was not made against commit ...`.
3. Release succeeds but `nuon-ext-api --help` reports an API version different from the release tag.

## Quick Diagnosis

1. Check tag vs release:

```bash
gh release view v0.19.821 --repo nuonco/nuon-ext-api
gh api repos/nuonco/nuon-ext-api/tags --paginate --jq '.[].name' | rg '^v0.19.821$'
```

2. Check what commit the tag points to:

```bash
git fetch origin --tags
git rev-list -n1 v0.19.821
git rev-parse origin/main
```

3. Check spec version at the tag:

```bash
git show v0.19.821:spec/doc.json | jq -r '.info.version'
```

## Manual Remediation (Tag Exists, Release Missing, Spec Mismatch)

Use this when the tag points to the wrong commit/spec.

1. Delete the remote/local tag:

```bash
git fetch origin --tags
git tag -d v0.19.821 || true
git push origin :refs/tags/v0.19.821
```

2. Create a clean remediation branch from `main`:

```bash
git checkout main
git pull --ff-only
git checkout -b fix/release-v0.19.821
```

3. Update `spec/doc.json` to the expected API version and commit via PR:

```bash
curl -sS https://api.nuon.co/docs/doc.json -o spec/doc.json
jq -r '.info.version' spec/doc.json
git add spec/doc.json
git commit -m "spec: update to v0.19.821"
git push -u origin fix/release-v0.19.821
```

4. Merge the PR to `main` (required by branch protection).

5. Recreate and push the tag at the merged `main` commit:

```bash
git checkout main
git pull --ff-only
git tag -a v0.19.821 -m "Release v0.19.821" "$(git rev-parse origin/main)"
git push origin refs/tags/v0.19.821
```

6. Re-run the `Release` workflow (`workflow_dispatch`).

## Backfill Release for Existing Good Tag

If a tag already points at the correct spec/version commit, just re-run the workflow. It will build from the tag and
publish the release.

## Post-Release Verification

1. Verify release exists:

```bash
gh release view v0.19.821 --repo nuonco/nuon-ext-api --json tagName,publishedAt,url
```

2. Verify binary-reported API version matches tag:

```bash
TMPDIR=$(mktemp -d)
gh release download v0.19.821 --repo nuonco/nuon-ext-api --pattern 'nuon-ext-api-darwin-arm64.tar.gz' --dir "$TMPDIR"
tar -xzf "$TMPDIR/nuon-ext-api-darwin-arm64.tar.gz" -C "$TMPDIR"
"$TMPDIR/nuon-ext-api" --help | rg 'API version:'
```
