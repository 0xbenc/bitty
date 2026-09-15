package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xbenc/bitty/internal/app"
	"github.com/0xbenc/bitty/internal/state"
	"github.com/0xbenc/bitty/internal/terminal"
	"github.com/0xbenc/bitty/internal/termstyle"
	"github.com/0xbenc/termintro"
	"github.com/0xbenc/termtheme"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type options struct {
	width, height int
	density       float64
	speed         int
	seed          uint64
	wrap          bool
	paused        bool
	intro         bool
	noIntro       bool
	noColor       bool
	themeFile     string
	showVersion   bool
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "bitty:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 1 && args[0] == "version" {
		printVersion()
		return nil
	}
	opts, err := parse(args)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}
	if opts.showVersion {
		printVersion()
		return nil
	}
	theme, err := termstyle.Resolve(opts.themeFile, os.Environ(), opts.noColor)
	if err != nil {
		return err
	}
	maybePlayIntro(opts, theme.NoColor)
	return app.Run(app.Options{
		Width: opts.width, Height: opts.height, Density: opts.density,
		Speed: opts.speed, Seed: opts.seed, Wrap: opts.wrap, Paused: opts.paused,
		Theme: theme, Input: os.Stdin, Output: os.Stdout,
	})
}

func parse(args []string) (options, error) {
	o := options{density: .25, speed: 12, seed: uint64(time.Now().UnixNano()), wrap: true}
	fs := newFlagSet(&o)
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	if fs.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	if o.width < 0 || o.height < 0 {
		return options{}, fmt.Errorf("width and height cannot be negative")
	}
	if o.density < 0 || o.density > 1 {
		return options{}, fmt.Errorf("density must be between 0 and 1")
	}
	if o.speed < 1 || o.speed > 60 {
		return options{}, fmt.Errorf("speed must be between 1 and 60")
	}
	if o.intro && o.noIntro {
		return options{}, fmt.Errorf("--intro and --no-intro cannot be used together")
	}
	return o, nil
}

func newFlagSet(o *options) *flag.FlagSet {
	fs := flag.NewFlagSet("bitty", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.IntVar(&o.width, "width", 0, "world width (default: terminal width)")
	fs.IntVar(&o.height, "height", 0, "world height (default: twice terminal height)")
	fs.Float64Var(&o.density, "density", o.density, "initial live-cell density, 0..1")
	fs.IntVar(&o.speed, "speed", o.speed, "generations per second, 1..60")
	fs.Func("seed", "random seed (default: current time)", func(value string) error {
		seed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid seed %q", value)
		}
		o.seed = seed
		return nil
	})
	fs.BoolVar(&o.wrap, "wrap", o.wrap, "connect opposite world edges")
	fs.BoolVar(&o.paused, "paused", false, "start paused")
	fs.BoolVar(&o.intro, "intro", false, "play the startup intro")
	fs.BoolVar(&o.noIntro, "no-intro", false, "skip the startup intro")
	fs.BoolVar(&o.noColor, "no-color", false, "disable ANSI colors")
	fs.StringVar(&o.themeFile, "theme-file", "", "portable .theme file")
	fs.BoolVar(&o.showVersion, "version", false, "print version information")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "bitty — set life loose in your terminal")
		fmt.Fprintln(fs.Output(), "\nUsage: bitty [options]\n       bitty version")
		fmt.Fprintln(fs.Output(), "\nOptions:")
		fs.PrintDefaults()
		fmt.Fprintln(fs.Output(), "\nKeys: arrows/hjkl move, space toggles, p pauses, n steps, r randomizes,")
		fmt.Fprintln(fs.Output(), "      c clears, +/- changes speed, w toggles wrapping, Esc/Ctrl-C/Ctrl-Q quit")
	}
	return fs
}

func maybePlayIntro(opts options, noColor bool) {
	if !terminal.IsTTY(os.Stderr) {
		return
	}
	env := termtheme.EnvMap(os.Environ())
	if opts.noIntro || termtheme.EnvTruthy(env["BITTY_NO_INTRO"]) {
		return
	}
	dir, err := state.ResolveDir(os.Environ())
	if err != nil {
		return
	}
	seen, _ := state.Read(dir)
	force := opts.intro || termtheme.EnvTruthy(env["BITTY_INTRO_ALWAYS"])
	if !force && seen.LastSeen == version {
		return
	}
	termintro.Play(termintro.Options{
		Title: "BITTY", Tagline: "SET LIFE LOOSE", Credits: []string{"0xbenc"},
		Version: version, Output: os.Stderr, NoColor: noColor,
	})
	_ = state.Write(dir, version)
}

func printVersion() {
	parts := []string{"bitty", version}
	if commit != "none" {
		parts = append(parts, commit)
	}
	if date != "unknown" {
		parts = append(parts, date)
	}
	fmt.Println(strings.Join(parts, " "))
}
