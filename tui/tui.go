package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sgoal/tide/agent"
)

func NewTUI() {
	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	textView.ScrollToEnd()

	inputField := tview.NewInputField().
		SetLabel("Enter a command: ").
		SetFieldWidth(0)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 3, 0, true)

	reactAgent, err := agent.NewReActAgent(textView)
	if err != nil {
		fmt.Fprintf(textView, "Error: %v\n", err)
	} else {
		err = reactAgent.LoadHistory()
		if err != nil {
			fmt.Fprintf(textView, "Error loading history: %v\n", err)
		}
		for _, msg := range reactAgent.GetHistory() {
			fmt.Fprintf(textView, "[yellow]%s:[white] %s\n", msg.Role, msg.Content)
		}
		textView.ScrollToEnd()
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := inputField.GetText()
			if strings.TrimSpace(command) == "" {
				return
			}

			fmt.Fprintf(textView, "[blue]You:[white] %s\n", command)
			inputField.SetText("")
			textView.ScrollToEnd()

			go func() {
				if reactAgent == nil {
					app.QueueUpdateDraw(func() {
						fmt.Fprintf(textView, "Agent not initialized. Please set OPENAI_API_KEY and restart.")
					})
					return
				}
				response, err := reactAgent.ProcessCommand(command)
				app.QueueUpdateDraw(func() {
					if err != nil {
						fmt.Fprintf(textView, "[red]Error:[white] %v\n", err)
					} else {
						fmt.Fprintf(textView, "[green]Agent:[white] %s\n", response)
					}
					reactAgent.SaveHistory()
					textView.ScrollToEnd()
				})
			}()
		}
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
