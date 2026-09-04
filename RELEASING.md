# Releasing

Releases are built from `vX.Y.Z` tags. Before tagging, ensure `main` CI is green
and configure `TAP_GITHUB_TOKEN` with write access to `0xbenc/homebrew-tap`.

```sh
git tag -a vX.Y.Z -m "bitty vX.Y.Z"
git push origin main
git push origin vX.Y.Z
```

The release workflow builds and attests archives, then GoReleaser updates the
`bitty` cask in the tap.
