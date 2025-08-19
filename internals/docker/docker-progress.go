package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImagePullMsg ImageBuildMessage
type DoneMsg struct{}
type ErrMsg struct{ Err error }

type PullModel struct {
	spinner  spinner.Model
	progress progress.Model
	status   string
	layer    string
	done     bool
	err      error
}

func NewPullModel() PullModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return PullModel{
		spinner:  s,
		progress: p,
	}
}

func (m PullModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m PullModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case ImagePullMsg:
		m.status = msg.Status
		m.layer = msg.ID
		if msg.ProgressDetail.Total > 0 {
			percent := float64(msg.ProgressDetail.Current) / float64(msg.ProgressDetail.Total)
			cmd := m.progress.SetPercent(percent)
			return m, cmd
		}
	case DoneMsg:
		m.done = true
		return m, tea.Quit
	case ErrMsg:
		m.err = msg.Err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newProg, ok := newModel.(progress.Model); ok {
			m.progress = newProg
		}
		return m, cmd
	}
	return m, nil
}

func (m PullModel) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error: %v", m.err))
	}
	if m.done {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✅ Image pull complete\n")
	}
	spin := m.spinner.View()
	info := fmt.Sprintf("Pulling %s: %s", m.layer, m.status)
	bar := m.progress.View()
	gap := strings.Repeat(" ", max(0, 60-len(info)))
	return fmt.Sprintf("%s %s%s %s", spin, info, gap, bar)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func StreamLogsTui(reader io.ReadCloser) error {
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	prog := tea.NewProgram(NewPullModel())

	// Run Bubble Tea in a goroutine
	go func() {
		for {
			var msg ImageBuildMessage
			if err := decoder.Decode(&msg); err != nil {
				if errors.Is(err, io.EOF) {
					prog.Send(DoneMsg{})
					break
				}
				prog.Send(ErrMsg{Err: fmt.Errorf("decode error: %w", err)})
				break
			}
			prog.Send(ImagePullMsg(msg))
		}
	}()

	// Block until UI finishes
	_, err := prog.Run()
	return err
}
