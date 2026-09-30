package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
)

func TestDisabledMatchNavigationAndViewport(t *testing.T) {
	v := viewport.New(10, 2)
	v.SetContent("one\nblocked-one\nblocked-two\ntwo")
	m := model{
		matches: []fuzzy.Match{
			{Str: "one"},
			{Str: "blocked-one"},
			{Str: "blocked-two"},
			{Str: "two"},
		},
		disabledChoices: map[string]bool{
			"blocked-one": true,
			"blocked-two": true,
		},
		viewport: &v,
	}

	m.CursorDown()
	if m.cursor != 3 || m.viewport.YOffset != 2 {
		t.Fatalf("cursor/offset = %d/%d, want 3/2", m.cursor, m.viewport.YOffset)
	}

	m.CursorUp()
	if m.cursor != 0 || m.viewport.YOffset != 0 {
		t.Fatalf("cursor/offset after moving up = %d/%d, want 0/0", m.cursor, m.viewport.YOffset)
	}
}

func TestDisabledMatchesAreExcludedFromToggleAll(t *testing.T) {
	v := viewport.New(10, 3)
	km := defaultKeymap()
	km.ToggleAll.SetEnabled(true)
	m := model{
		matches: []fuzzy.Match{{Str: "one"}, {Str: "blocked"}, {Str: "two"}},
		disabledChoices: map[string]bool{
			"blocked": true,
		},
		textinput: textinput.New(),
		viewport:  &v,
		selected:  make(map[string]struct{}),
		limit:     3,
		keymap:    km,
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlA})
	m = updated.(model)
	if m.numSelected != 2 {
		t.Fatalf("selected %d items, want 2", m.numSelected)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlA})
	m = updated.(model)
	if m.numSelected != 0 {
		t.Fatalf("selected %d items after toggling all again, want 0", m.numSelected)
	}
}

func TestSubmitDisabledMatchDoesNotQuit(t *testing.T) {
	v := viewport.New(10, 1)
	m := model{
		matches:         []fuzzy.Match{{Str: "blocked"}},
		disabledChoices: map[string]bool{"blocked": true},
		textinput:       textinput.New(),
		viewport:        &v,
		selected:        make(map[string]struct{}),
		keymap:          defaultKeymap(),
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if got.submitted || got.quitting {
		t.Fatal("submitting a disabled match should not quit or submit")
	}
}

func TestViewUsesOptionLabel(t *testing.T) {
	v := viewport.New(20, 1)
	m := model{
		matches:        []fuzzy.Match{{Str: "Display label"}},
		choices:        map[string]string{"Display label": "output-value"},
		displayChoices: map[string]string{"Display label": "Display label"},
		textinput:      textinput.New(),
		viewport:       &v,
		selected:       make(map[string]struct{}),
		padding:        []int{0, 0, 0, 0},
	}

	view := m.View()
	if !strings.Contains(view, "Display label") || strings.Contains(view, "output-value") {
		t.Fatalf("view does not display the option label: %q", view)
	}
}

func TestMatchedRanges(t *testing.T) {
	for name, tt := range map[string]struct {
		in  []int
		out [][2]int
	}{
		"empty": {
			in:  []int{},
			out: [][2]int{},
		},
		"one char": {
			in:  []int{1},
			out: [][2]int{{1, 1}},
		},
		"2 char range": {
			in:  []int{1, 2},
			out: [][2]int{{1, 2}},
		},
		"multiple char range": {
			in:  []int{1, 2, 3, 4, 5, 6},
			out: [][2]int{{1, 6}},
		},
		"multiple char ranges": {
			in:  []int{1, 2, 3, 5, 6, 10, 11, 12, 13, 23, 24, 40, 42, 43, 45, 52},
			out: [][2]int{{1, 3}, {5, 6}, {10, 13}, {23, 24}, {40, 40}, {42, 43}, {45, 45}, {52, 52}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			match := matchedRanges(tt.in)
			if !reflect.DeepEqual(match, tt.out) {
				t.Errorf("expected %v, got %v", tt.out, match)
			}
		})
	}
}

func TestByteToChar(t *testing.T) {
	stStr := "\x1b[90m\ue615\x1b[39m \x1b[3m\x1b[32mDow\x1b[0m\x1b[90m\x1b[39m\x1b[3wnloads"
	str := " Downloads"
	rng := [2]int{4, 7}
	expect := "Dow"

	if got := str[rng[0]:rng[1]]; got != expect {
		t.Errorf("expected %q, got %q", expect, got)
	}

	start, stop := bytePosToVisibleCharPos(str, rng)
	if got := ansi.Strip(ansi.Cut(stStr, start, stop)); got != expect {
		t.Errorf("expected %+q, got %+q", expect, got)
	}
}
