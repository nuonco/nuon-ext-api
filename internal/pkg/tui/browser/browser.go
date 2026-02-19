package browser

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nuonco/nuon-ext-api/internal/pkg/tui"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// Result is returned after the browser exits.
type Result struct {
	Route    *spec.Route // nil if user quit without selecting
	Selected bool
}

// Run launches the interactive endpoint browser and returns the selected route.
func Run(api *spec.API) (*Result, error) {
	items := make([]list.Item, len(api.Routes))
	for i, r := range api.Routes {
		items[i] = routeItem{route: r}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(tui.PrimaryColor).
		BorderLeftForeground(tui.PrimaryColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(tui.SubtleColor).
		BorderLeftForeground(tui.PrimaryColor)

	l := list.New(items, delegate, 80, 24)
	l.Title = fmt.Sprintf("Nuon API v%s — %d endpoints", api.Version, len(api.Routes))
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(tui.TextColor).
		Background(tui.PrimaryColor).
		Padding(0, 1)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		}
	}

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	if final, ok := result.(model); ok && final.selected != nil {
		return &Result{Route: final.selected, Selected: true}, nil
	}

	return &Result{}, nil
}

// model is the bubbletea model for the endpoint browser.
type model struct {
	list     list.Model
	selected *spec.Route
}

func (m model) Init() tea.Cmd {
	return nil
}

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
			if item, ok := m.list.SelectedItem().(routeItem); ok {
				r := item.route
				m.selected = &r
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
