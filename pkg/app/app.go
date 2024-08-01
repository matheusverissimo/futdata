package app

import (
	"futdata/pkg/db"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list         list.Model
	todosTimes   []string
	cursor       int
	start        int
	end          int
}

type listItem struct {
	nome string
}

func (l listItem) FilterValue() string {
	return l.nome
}

func (i listItem) Title() string { return i.nome }
func (i listItem) Description() string { return "" }

func InitialModel() Model {
	times, _ := db.NewRepository().FindAllTimes()
	var listItems []list.Item

	for _, time := range times {
		listItems = append(listItems, listItem{nome: time})
	}

	m := Model{
		list:         list.New(listItems, list.NewDefaultDelegate(), 0, 0),
		todosTimes:   times,
	}

	m.list.Title = "Futdata - Times"

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:
			// Cool, what was the actual key pressed?
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			
		case tea.WindowSizeMsg:
			m.list.SetSize(msg.Width - 1, msg.Height - 1)
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	}

// View implements tea.Model.
func (m Model) View() string {
	return m.list.View()
}
