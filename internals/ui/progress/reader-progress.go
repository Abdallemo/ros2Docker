package progress

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// messages
type tickMsg struct{}
type doneMsg struct{}

// Model for progress bar
type Model struct {
	progress progress.Model
	pct      float64
	done     bool
}

// NewProgress creates a new Model
func NewProgress() *Model {
	return &Model{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m *Model) Init() tea.Cmd {
	return tickProgress()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ms := msg.(type) {
	case tea.KeyMsg:
		if ms.String() == "q" || ms.String() == "esc" || ms.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		if !m.done {

			m.pct += 0.01
			if m.pct >= 1.0 {
				m.pct = 0.99
			}
			return m, tickProgress()
		}
	case doneMsg:
		m.done = true
		m.pct = 1.0
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) View() string {
	if m.done {
		return "\n✅ Finished installing dependencies!\n"
	}
	return fmt.Sprintf("\n🔧 Installing dependencies...\n\n%s\n", m.progress.ViewAs(m.pct))
}

func tickProgress() tea.Cmd {
	return tea.Tick(time.Millisecond*120, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

// RunProgress runs the progress bar until reader is consumed OR context says exec finished
func RunProgress(r io.Reader, done <-chan struct{}) error {
	m := NewProgress()
	p := tea.NewProgram(m, tea.WithAltScreen())

	// consume output so docker doesn't block
	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := r.Read(buf)
			if err != nil {
				break
			}
		}
	}()

	// wait until "done" signal, then finish
	go func() {
		<-done
		p.Send(doneMsg{})
	}()

	_, err := p.Run()
	return err
}
