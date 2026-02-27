package main

import (
	"os"
	"testing"
)

func TestInTmux(t *testing.T) {
	// Save and restore original value
	orig := os.Getenv("TMUX")
	defer os.Setenv("TMUX", orig)

	os.Setenv("TMUX", "/tmp/tmux-501/default,12345,0")
	if !inTmux() {
		t.Error("inTmux() should return true when TMUX is set")
	}

	os.Unsetenv("TMUX")
	if inTmux() {
		t.Error("inTmux() should return false when TMUX is unset")
	}
}

func TestResolveModeAuto(t *testing.T) {
	orig := os.Getenv("TMUX")
	defer os.Setenv("TMUX", orig)

	// In tmux
	os.Setenv("TMUX", "/tmp/tmux-501/default,12345,0")
	if got := resolveMode("auto"); got != "tmux" {
		t.Errorf("resolveMode(\"auto\") in tmux = %q, want \"tmux\"", got)
	}

	// Not in tmux
	os.Unsetenv("TMUX")
	if got := resolveMode("auto"); got != "window" {
		t.Errorf("resolveMode(\"auto\") outside tmux = %q, want \"window\"", got)
	}
}

func TestResolveModeExplicit(t *testing.T) {
	modes := []string{"inline", "tmux", "window"}
	for _, mode := range modes {
		if got := resolveMode(mode); got != mode {
			t.Errorf("resolveMode(%q) = %q, want %q", mode, got, mode)
		}
	}
}

func TestSelfPath(t *testing.T) {
	path := selfPath()
	if path == "" {
		t.Error("selfPath() returned empty string")
	}
}
