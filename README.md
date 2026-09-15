# bitty

Conway's Game of Life, packed two cells high in your terminal. Pause it, draw a
pattern, step it carefully, or set the whole thing loose.

```sh
brew install --cask 0xbenc/tap/bitty
bitty
```

## Controls

| Key | Action |
|---|---|
| arrows or `hjkl` | move the editing cursor |
| `space` or `x` | toggle the selected cell |
| `p` or `enter` | pause/resume |
| `n` | advance one generation |
| `r` | regenerate from the configured density |
| `c` | clear and pause |
| `+` / `-` | change simulation speed |
| `w` | toggle edge wrapping |
| `Esc`, `Ctrl-C`, or `Ctrl-Q` | quit |

The simulation starts immediately with a random field. Use `--paused` to begin
in editing mode and `--density 0` for an empty canvas.

```sh
bitty --density 0.25 --speed 12 --seed 42
bitty --paused --density 0 --width 120 --height 80
bitty --wrap=false
```

Run `bitty --help` for every option. `NO_COLOR` and `BITTY_NO_COLOR` disable
color. `BITTY_THEME_FILE` selects a portable termsystem `.theme` file.

## Build

Bitty supports Linux and macOS and requires Go 1.26.5 or newer.

```sh
go build ./cmd/bitty
go test ./...
```

## License

MIT
