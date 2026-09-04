// Package termstyle is bitty's small adapter around the shared theme engine.
package termstyle

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/0xbenc/termtheme"
)

type Theme = termtheme.Theme
type Role = termtheme.Role

const (
	RoleTitle       = termtheme.RoleTitle
	RolePrimary     = termtheme.RolePrimary
	RoleMuted       = termtheme.RoleMuted
	RoleForeground  = termtheme.RoleForeground
	RoleSelected    = termtheme.RoleSelected
	RoleSelectedBar = termtheme.RoleSelectedBar
	RoleBorder      = termtheme.RoleBorder
	RoleSuccess     = termtheme.RoleSuccess
	RoleWarning     = termtheme.RoleWarning
)

var (
	VisibleWidth = termtheme.VisibleWidth
	Sanitize     = termtheme.Sanitize
	Truncate     = termtheme.Truncate
)

// TerminalTheme uses ordinary ANSI colors and remains legible in light and dark terminals.
func TerminalTheme() Theme {
	return Theme{Name: "terminal", Codes: map[termtheme.Role]string{
		RoleTitle:       "1;36",
		RolePrimary:     "36",
		RoleMuted:       "2",
		RoleForeground:  "39",
		RoleSelected:    "1;7",
		RoleSelectedBar: "7",
		RoleBorder:      "2;36",
		RoleSuccess:     "32",
		RoleWarning:     "33",
	}}
}

// Resolve loads BITTY_THEME_FILE (or the shared default path) and applies
// portable role overrides. Missing default files and unknown base names fail open.
func Resolve(file string, env []string, noColor bool) (Theme, error) {
	theme := TerminalTheme()
	path, explicit := termtheme.ResolveThemeFile("bitty", file, env, false)
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if explicit || !errors.Is(err, fs.ErrNotExist) {
				return Theme{}, fmt.Errorf("load theme %s: %w", path, err)
			}
		} else {
			cfg, err := termtheme.ParseThemeConfig(data)
			if err != nil {
				return Theme{}, fmt.Errorf("load theme %s: %w", path, err)
			}
			for role, code := range cfg.Codes {
				theme.Codes[role] = code
			}
		}
	}
	theme.NoColor = termtheme.EnvNoColor("bitty", env, noColor)
	return theme, nil
}
