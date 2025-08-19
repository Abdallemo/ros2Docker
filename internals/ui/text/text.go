package text

import (
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

type LogLevel int

const (
	Info LogLevel = iota
	Warn
	Error
	Ok
)

type logMsg string
type stopMsg struct{}

type Model struct {
	mu   sync.Mutex
	logs []string
	done bool
	prog *tea.Program
}

func New() *Model {
	return &Model{}
}

// Start the TUI
func (m *Model) Run() {
	m.prog = tea.NewProgram(m)
	go func() { _ = m.prog.Start() }() // run async
}

// Append a message from anywhere
func (m *Model) Append(msg string, level LogLevel) {

	if m.prog == nil {
		return
	}
	if strings.Contains(strings.ToLower(msg), "error") {
		m.prog.Send(logMsg(errorStyle.Render(msg)))
	}
	var style lipgloss.Style
	switch level {
	case Error:
		style = errorStyle
	case Warn:
		style = warnStyle
	case Ok:
		style = okStyle
	default:
		style = infoStyle
	}

	m.prog.Send(logMsg(style.Render(msg)))

}

// Stop the UI
func (m *Model) Stop() {
	if m.prog != nil {
		m.prog.Send(stopMsg{})
	}
}

// --- Bubble Tea impl ---
func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case logMsg:
		m.mu.Lock()
		m.logs = append(m.logs, string(msg))
		m.mu.Unlock()
	case stopMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder
	m.mu.Lock()
	for _, l := range m.logs {
		b.WriteString(l + "\n")
	}
	m.mu.Unlock()

	if m.done {
		b.WriteString("\n" + okStyle.Bold(true).Render("✔ Finished!") + "\n")
	}
	return b.String()
}
