package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompletionsMentionFlags guards the hand-maintained shell completion
// scripts against drifting from the FlagSet built by newFlagSet: every flag
// the FlagSet defines must appear textually in each completion file.
func TestCompletionsMentionFlags(t *testing.T) {
	var flags []string
	newFlagSet(&options{}).VisitAll(func(f *flag.Flag) {
		flags = append(flags, f.Name)
	})
	if len(flags) == 0 {
		t.Fatal("no flags collected from FlagSet")
	}
	for _, file := range []string{"bitty.bash", "bitty.zsh", "bitty.fish"} {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", "completions", file))
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			text := string(data)
			for _, name := range flags {
				// Fish completions spell flags as "-l name" rather than "--name".
				want := "--" + name
				if strings.HasSuffix(file, ".fish") {
					want = "-l " + name
				}
				if !strings.Contains(text, want) {
					t.Errorf("%s missing flag %q", file, want)
				}
			}
		})
	}
}
