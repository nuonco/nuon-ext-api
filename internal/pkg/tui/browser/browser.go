package browser

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nuonco/nuon-ext-api/internal/pkg/tui"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// Action describes what the caller should do after the browser exits.
type Action int

const (
	ActionNone    Action = iota // user quit without selecting
	ActionSelect                // user pressed enter — print the route
	ActionExecute               // user pressed x — execute the GET endpoint
)

// Result is returned after the browser exits.
type Result struct {
	Route  *spec.Route
	Action Action
}

// Run launches the interactive endpoint browser and returns the selected route.
func Run(api *spec.API, apiURL string) (*Result, error) {
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

	keys := keyMap{
		Open: key.NewBinding(
			key.WithKeys("B"),
			key.WithHelp("B", "open docs"),
		),
		Execute: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "execute GET"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.Open, keys.Execute}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.Open, keys.Execute}
	}

	m := model{list: l, keys: keys, apiURL: apiURL}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	if final, ok := result.(model); ok && final.result != nil {
		return final.result, nil
	}

	return &Result{Action: ActionNone}, nil
}

type keyMap struct {
	Open    key.Binding
	Execute key.Binding
	Select  key.Binding
}

// model is the bubbletea model for the endpoint browser.
type model struct {
	list   list.Model
	keys   keyMap
	apiURL string
	result *Result
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
		switch {
		case key.Matches(msg, m.keys.Open):
			if item, ok := m.list.SelectedItem().(routeItem); ok {
				openBrowser(item.route.DocsURL(m.apiURL))
			}
			return m, nil
		case key.Matches(msg, m.keys.Execute):
			if item, ok := m.list.SelectedItem().(routeItem); ok {
				if item.route.Method == "GET" {
					r := item.route
					m.result = &Result{Route: &r, Action: ActionExecute}
					return m, tea.Quit
				}
			}
			return m, nil
		case msg.String() == "enter":
			if item, ok := m.list.SelectedItem().(routeItem); ok {
				r := item.route
				m.result = &Result{Route: &r, Action: ActionSelect}
			}
			return m, tea.Quit
		case msg.String() == "q" || msg.String() == "ctrl+c":
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

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	cmd.Start()
}
