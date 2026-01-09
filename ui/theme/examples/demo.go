package examples

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	currentExample int
	examples       []example
}

type example struct {
	name   string
	render func() string
}

func newModel() model {
	return model{
		currentExample: 0,
		examples: []example{
			{"Modal Dialog", ModalExample},
			{"Status Bar", StatusBarExample},
			{"Input Field", InputExample},
			{"List View", ListExample},
			{"Preview Pane", PreviewExample},
			{"Diff Display", DiffExample},
			{"Chat Interface", ChatExample},
			{"Validation Errors", ValidationExample},
			{"Text Highlighting", HighlightExample},
			{"Keyboard Shortcuts", KeyboardExample},
			{"Tool Execution", ToolCallExample},
			{"Status Messages", StatusMessageExample},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
				return m, tea.Quit
			}
		case tea.KeyRight, tea.KeySpace:
			m.currentExample = (m.currentExample + 1) % len(m.examples)
			return m, tick()
		case tea.KeyLeft:
			m.currentExample = (m.currentExample - 1 + len(m.examples)) % len(m.examples)
			return m, tick()
		}
	case tickMsg:
		m.currentExample = (m.currentExample + 1) % len(m.examples)
		return m, tick()
	case tea.WindowSizeMsg:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	ex := m.examples[m.currentExample]
	title := theme.HeaderStyle().Render(fmt.Sprintf("=== OpenCode Design System: %s (%d/%d) ===", ex.name, m.currentExample+1, len(m.examples)))

	help := "\n" + theme.MutedStyle().Render(
		"Auto-advancing... | Left: Previous | Right/Space: Next | Ctrl+C/Esc/Q: Quit",
	)

	return title + "\n\n" + ex.render() + help
}

func RunDemo() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running demo: %v\n", err)
	}
}
