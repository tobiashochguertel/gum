package choose

import (
	"testing"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

func TestMoveCursorSkipsDisabledItemsAndSynchronizesPage(t *testing.T) {
	pager := paginator.New()
	pager.SetTotalPages(5)
	pager.PerPage = 1
	m := model{
		items: []item{
			{text: "one"},
			{text: "two", disabled: true},
			{text: "three", disabled: true},
			{text: "four"},
			{text: "five", disabled: true},
		},
		height:    1,
		paginator: pager,
	}

	m.moveCursor(1, 1, false)
	if m.index != 3 || m.paginator.Page != 3 {
		t.Fatalf("cursor/page = %d/%d, want 3/3", m.index, m.paginator.Page)
	}

	m.moveCursor(4, 1, false)
	if m.index != 3 || m.paginator.Page != 3 {
		t.Fatalf("cursor/page after blocked move = %d/%d, want 3/3", m.index, m.paginator.Page)
	}

	m.moveCursor(2, -1, false)
	if m.index != 0 || m.paginator.Page != 0 {
		t.Fatalf("cursor/page after moving backward = %d/%d, want 0/0", m.index, m.paginator.Page)
	}
}

func TestMoveCursorKeepsPositionWhenAllItemsAreDisabled(t *testing.T) {
	pager := paginator.New()
	pager.SetTotalPages(2)
	pager.PerPage = 1
	m := model{
		items:     []item{{disabled: true}, {disabled: true}},
		index:     1,
		height:    1,
		paginator: pager,
	}

	m.moveCursor(0, 1, true)
	if m.index != 1 || m.paginator.Page != 1 {
		t.Fatalf("cursor/page = %d/%d, want 1/1", m.index, m.paginator.Page)
	}
}

func TestSubmitDisabledItemDoesNotQuit(t *testing.T) {
	m := model{
		items:  []item{{text: "blocked", disabled: true}},
		height: 1,
		keymap: defaultKeymap(),
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if got.submitted || got.quitting {
		t.Fatal("submitting a disabled item should not quit or submit")
	}
}

func TestToggleAllUsesSelectableItems(t *testing.T) {
	m := model{
		items: []item{
			{text: "one"},
			{text: "blocked", disabled: true},
			{text: "two"},
		},
		limit: 3,
		keymap: func() keymap {
			km := defaultKeymap()
			km.ToggleAll.SetEnabled(true)
			return km
		}(),
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
