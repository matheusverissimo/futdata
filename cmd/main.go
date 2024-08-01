package main

import (
	"fmt"
	"futdata/pkg/app"
	"futdata/pkg/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cwd, _ := os.Getwd()
	fmt.Printf("%v", cwd)
	if err := db.InitDatabase(); err != nil {
		fmt.Printf("A base de dados não foi encontrada!\n")
		return
	}
	p := tea.NewProgram(app.InitialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
