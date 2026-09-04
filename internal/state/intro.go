// Package state stores bitty's small per-user runtime bookkeeping.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/0xbenc/termtheme"
)

type Intro struct {
	SchemaVersion int    `json:"schema_version"`
	LastSeen      string `json:"last_seen_intro_version"`
}

func ResolveDir(env []string) (string, error) {
	values := termtheme.EnvMap(env)
	if dir := strings.TrimSpace(values["BITTY_STATE_DIR"]); dir != "" {
		return filepath.Clean(dir), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "bitty"), nil
	}
	if dir := strings.TrimSpace(values["XDG_STATE_HOME"]); dir != "" {
		return filepath.Join(dir, "bitty"), nil
	}
	return filepath.Join(home, ".local", "state", "bitty"), nil
}

func Read(dir string) (Intro, error) {
	data, err := os.ReadFile(filepath.Join(dir, "intro.json"))
	if errors.Is(err, os.ErrNotExist) {
		return Intro{}, nil
	}
	if err != nil {
		return Intro{}, err
	}
	var value Intro
	if err := json.Unmarshal(data, &value); err != nil {
		return Intro{}, err
	}
	return value, nil
}

func Write(dir, version string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(Intro{SchemaVersion: 1, LastSeen: version}, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(dir, ".intro-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, filepath.Join(dir, "intro.json"))
}
