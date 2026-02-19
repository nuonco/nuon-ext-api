package selector

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nuonco/nuon-ext-api/internal/pkg/tui"
)

// Result is returned after the selector exits.
type Result struct {
	ID       string
	Name     string
	Selected bool
}

// Resource represents one item from an API list response.
type Resource struct {
	ID   string
	Name string
}

func (r Resource) Title() string       { return r.Name }
func (r Resource) Description() string { return r.ID }
func (r Resource) FilterValue() string { return r.Name + " " + r.ID }

// Run launches an interactive selector for the given resources.
func Run(paramName string, resources []Resource) (*Result, error) {
	if len(resources) == 0 {
		return nil, fmt.Errorf("no resources available for {%s}", paramName)
	}

	items := make([]list.Item, len(resources))
	for i, r := range resources {
		items[i] = r
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(tui.PrimaryColor).
		BorderLeftForeground(tui.PrimaryColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(tui.SubtleColor).
		BorderLeftForeground(tui.PrimaryColor)

	l := list.New(items, delegate, 60, 16)
	l.Title = fmt.Sprintf("Select {%s}", paramName)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(tui.TextColor).
		Background(tui.PrimaryColor).
		Padding(0, 1)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	if final, ok := result.(model); ok && final.selected != nil {
		return &Result{
			ID:       final.selected.ID,
			Name:     final.selected.Name,
			Selected: true,
		}, nil
	}

	return &Result{}, nil
}

type model struct {
	list     list.Model
	selected *Resource
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(Resource); ok {
				m.selected = &item
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
