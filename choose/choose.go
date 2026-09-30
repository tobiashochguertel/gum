// Package choose provides an interface to choose one option from a given list
// of options. The options can be provided as (new-line separated) stdin or a
// list of arguments.
//
// It is different from the filter command as it does not provide a fuzzy
// finding input, so it is best used for smaller lists of options.
//
// Let's pick from a list of gum flavors:
//
// $ gum choose "Strawberry" "Banana" "Cherry"
package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func defaultKeymap() keymap {
	return keymap{
		Down: key.NewBinding(
			key.WithKeys("down", "j", "ctrl+j", "ctrl+n"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k", "ctrl+k", "ctrl+p"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "ctrl+f"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "ctrl+b"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
		),
		ToggleAll: key.NewBinding(
			key.WithKeys("a", "A", "ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
			key.WithDisabled(),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" ", "tab", "x", "ctrl+@"),
			key.WithHelp("x", "toggle"),
			key.WithDisabled(),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter", "ctrl+q"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type keymap struct {
	Down,
	Up,
	Right,
	Left,
	Home,
	End,
	ToggleAll,
	Toggle,
	Abort,
	Quit,
	Submit key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Toggle,
		key.NewBinding(
			key.WithKeys("up", "down", "right", "left"),
			key.WithHelp("←↓↑→", "navigate"),
		),
		k.Submit,
		k.ToggleAll,
	}
}

type model struct {
	height           int
	padding          []int
	cursor           string
	selectedPrefix   string
	unselectedPrefix string
	disabledPrefix   string
	cursorPrefix     string
	header           string
	items            []item
	quitting         bool
	submitted        bool
	index            int
	limit            int
	numSelected      int
	currentOrder     int
	paginator        paginator.Model
	showHelp         bool
	help             help.Model
	keymap           keymap

	// styles
	cursorStyle       lipgloss.Style
	headerStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	disabledItemStyle lipgloss.Style
}

type item struct {
	text     string
	selected bool
	disabled bool
	order    int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:
		km := m.keymap
		switch {
		case key.Matches(msg, km.Down):
			m.moveCursor(m.index+1, 1, true)
		case key.Matches(msg, km.Up):
			m.moveCursor(m.index-1, -1, true)
		case key.Matches(msg, km.Right):
			m.moveCursor(min(m.index+m.height, len(m.items)-1), 1, false)
		case key.Matches(msg, km.Left):
			m.moveCursor(max(m.index-m.height, 0), -1, false)
		case key.Matches(msg, km.End):
			m.moveCursor(len(m.items)-1, -1, false)
		case key.Matches(msg, km.Home):
			m.moveCursor(0, 1, false)
		case key.Matches(msg, km.ToggleAll):
			if m.limit <= 1 {
				break
			}
			if m.numSelected < m.selectableItems() && m.numSelected < m.limit {
				m = m.selectAll()
			} else {
				m = m.deselectAll()
			}
		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Abort):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, km.Toggle):
			if m.limit == 1 {
				break // no op
			}

			// Don't allow toggling disabled items
			if m.items[m.index].disabled {
				break
			}

			if m.items[m.index].selected {
				m.items[m.index].selected = false
				m.numSelected--
			} else if m.numSelected < m.limit {
				m.items[m.index].selected = true
				m.items[m.index].order = m.currentOrder
				m.numSelected++
				m.currentOrder++
			}
		case key.Matches(msg, km.Submit):
			if m.items[m.index].disabled {
				break
			}
			m.quitting = true
			if m.limit <= 1 && m.numSelected < 1 {
				m.items[m.index].selected = true
			}
			m.submitted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m *model) moveCursor(target, direction int, wrap bool) {
	if len(m.items) == 0 {
		return
	}
	originalIndex := m.index
	for visited := 0; visited < len(m.items); visited++ {
		if target < 0 || target >= len(m.items) {
			if !wrap {
				break
			}
			target = (target + len(m.items)) % len(m.items)
		}
		if !m.items[target].disabled {
			m.index = target
			m.paginator.Page = m.index / m.height
			return
		}
		target += direction
	}
	m.index = originalIndex
	m.paginator.Page = m.index / m.height
}

func (m model) selectableItems() int {
	count := 0
	for _, item := range m.items {
		if !item.disabled {
			count++
		}
	}
	return count
}

func (m model) selectAll() model {
	for i := range m.items {
		if m.numSelected >= m.limit {
			break // do not exceed given limit
		}
		// Skip disabled items
		if m.items[i].disabled {
			continue
		}
		// Skip already selected items
		if m.items[i].selected {
			continue
		}
		m.items[i].selected = true
		m.items[i].order = m.currentOrder
		m.numSelected++
		m.currentOrder++
	}
	return m
}

func (m model) deselectAll() model {
	for i := range m.items {
		m.items[i].selected = false
		m.items[i].order = 0
	}
	m.numSelected = 0
	m.currentOrder = 0
	return m
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i, item := range m.items[start:end] {
		if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursor))
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(m.cursor)))
		}

		if item.disabled {
			s.WriteString(m.disabledItemStyle.Render(m.disabledPrefix + item.text))
		} else if item.selected {
			s.WriteString(m.selectedItemStyle.Render(m.selectedPrefix + item.text))
		} else if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursorPrefix + item.text))
		} else {
			s.WriteString(m.itemStyle.Render(m.unselectedPrefix + item.text))
		}
		if i != m.height {
			s.WriteRune('\n')
		}
	}

	if m.paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.height-m.paginator.ItemsOnPage(len(m.items))+1))
		s.WriteString("  " + m.paginator.View())
	}

	var parts []string

	if m.header != "" {
		parts = append(parts, m.headerStyle.Render(m.header))
	}
	parts = append(parts, s.String())
	if m.showHelp {
		parts = append(parts, "", m.help.View(m.keymap))
	}

	view := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(view)
}
