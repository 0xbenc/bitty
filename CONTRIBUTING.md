# Contributing

Before opening a pull request, run:

```sh
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

## Shared TUI stack & drift guards

Bitty consumes the four termsystem modules: `termtheme` for portable semantic
themes, `termnav` for viewport windowing, `termchrome` for the standard footer,
and `termintro` for the once-per-version startup animation. Keep tagged module
pins aligned with the sibling apps and never ship a `replace` directive.

Footers must use `termchrome.Footer`; do not inline its separators. Do not add
inline spinner frames or a golden-update flag. The `tui-conformance` CI job
guards those boundaries. The intro remains TTY-gated, writes to stderr, and
honors `--intro`, `--no-intro`, `BITTY_INTRO_ALWAYS`, and `BITTY_NO_INTRO`.
