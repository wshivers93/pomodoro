package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput redirects stdout and captures what a function prints.
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestRenderDisplayContainsLabel(t *testing.T) {
	tests := []struct {
		label string
	}{
		{"Work"},
		{"Break"},
	}
	for _, tt := range tests {
		out := captureOutput(func() {
			renderDisplay(tt.label, 60, 120, false)
		})
		if !strings.Contains(out, tt.label) {
			t.Errorf("renderDisplay output should contain %q", tt.label)
		}
	}
}

func TestRenderDisplayContainsTime(t *testing.T) {
	out := captureOutput(func() {
		renderDisplay("Work", 90, 300, false)
	})
	// 90 seconds = 01:30
	if !strings.Contains(out, "01:30") {
		t.Errorf("renderDisplay output should contain \"01:30\", got:\n%s", out)
	}
}

func TestRenderDisplayPaused(t *testing.T) {
	out := captureOutput(func() {
		renderDisplay("Work", 60, 120, true)
	})
	if !strings.Contains(out, "PAUSED") {
		t.Error("renderDisplay with paused=true should contain \"PAUSED\"")
	}
}

func TestRenderDisplayNotPaused(t *testing.T) {
	out := captureOutput(func() {
		renderDisplay("Work", 60, 120, false)
	})
	if strings.Contains(out, "PAUSED") {
		t.Error("renderDisplay with paused=false should not contain \"PAUSED\"")
	}
}

func TestRenderDisplayProgressBar(t *testing.T) {
	// Half elapsed: 60 remaining out of 120 total
	out := captureOutput(func() {
		renderDisplay("Work", 60, 120, false)
	})
	if !strings.Contains(out, "█") {
		t.Error("renderDisplay should contain filled bar characters")
	}
	if !strings.Contains(out, "░") {
		t.Error("renderDisplay should contain empty bar characters")
	}
	if !strings.Contains(out, "50%") {
		t.Errorf("renderDisplay at halfway should show 50%%")
	}
}

func TestRenderDisplayZeroTotal(t *testing.T) {
	// Should not panic
	out := captureOutput(func() {
		renderDisplay("Work", 0, 0, false)
	})
	if !strings.Contains(out, "00:00") {
		t.Errorf("renderDisplay with zero total should show 00:00")
	}
}

func TestRenderDisplayFullyElapsed(t *testing.T) {
	out := captureOutput(func() {
		renderDisplay("Work", 0, 120, false)
	})
	if !strings.Contains(out, "100%") {
		t.Errorf("renderDisplay at completion should show 100%%")
	}
}

func TestRenderInlineContainsLabel(t *testing.T) {
	out := captureOutput(func() {
		renderInline("Work", 90, 300, false)
	})
	if !strings.Contains(out, "Work") {
		t.Error("renderInline output should contain \"Work\"")
	}
}

func TestRenderInlineContainsTime(t *testing.T) {
	out := captureOutput(func() {
		renderInline("Break", 65, 300, false)
	})
	// 65 seconds = 01:05
	if !strings.Contains(out, "01:05") {
		t.Errorf("renderInline output should contain \"01:05\"")
	}
}

func TestRenderInlinePaused(t *testing.T) {
	out := captureOutput(func() {
		renderInline("Work", 60, 120, true)
	})
	if !strings.Contains(out, "⏸") {
		t.Error("renderInline with paused=true should contain pause icon")
	}
}

func TestRenderInlineNotPaused(t *testing.T) {
	out := captureOutput(func() {
		renderInline("Work", 60, 120, false)
	})
	if strings.Contains(out, "⏸") {
		t.Error("renderInline with paused=false should not contain pause icon")
	}
}

func TestRenderDone(t *testing.T) {
	out := captureOutput(func() {
		renderDone()
	})
	if !strings.Contains(out, "All done") {
		t.Error("renderDone should contain \"All done\"")
	}
}

func TestRenderInlineDone(t *testing.T) {
	out := captureOutput(func() {
		renderInlineDone()
	})
	if !strings.Contains(out, "All done") {
		t.Error("renderInlineDone should contain \"All done\"")
	}
}
