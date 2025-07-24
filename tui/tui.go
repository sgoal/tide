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

	// Create a form for mode selection
	form := tview.NewForm().
		AddButton("Builder Mode", func() {
			showBuilderMode(app)
		}).
		AddButton("SOLO Mode", func() {
			showSoloMode(app)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Select a mode").SetTitleAlign(tview.AlignCenter)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func showBuilderMode(app *tview.Application) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	textView.SetBorder(true).SetTitle("Tide")
	textView.ScrollToEnd()
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)

	agent, err := agent.NewReActAgent(textView)
	if err != nil {
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(textView, "[red]Error:[white] %v\n", err)
		})
	} else {
		err = agent.LoadHistory()
		if err != nil {
			fmt.Fprintf(textView, "Error loading history: %v\n", err)
		}
		for _, msg := range agent.GetHistory() {
			fmt.Fprintf(textView, "[yellow]%s:[white] %s\n", msg.Role, msg.Content)
		}
		textView.ScrollToEnd()
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			prompt := inputField.GetText()
			if strings.TrimSpace(prompt) == "" {
				return
			}
			inputField.SetText("")
			textView.ScrollToEnd()
			go func() {
				response, err := agent.ProcessCommand(prompt)
				if err != nil {
					app.QueueUpdateDraw(func() {
						fmt.Fprintf(textView, "[red]Error:[white] %v\n", err)
					})
				} else {
					app.QueueUpdateDraw(func() {
						fmt.Fprintf(textView, "[green]Agent:[white] %s\n", response)
					})
				}
			}()
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 3, 0, true)

	app.SetRoot(flex, true)
}

func showSoloMode(app *tview.Application) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	textView.SetBorder(true).SetTitle("SOLO Mode Output")

	inputField := tview.NewInputField().
		SetLabel("Enter your project requirement: ").
		SetFieldWidth(0)

	soloAgent, err := agent.NewSoloAgent(textView)
	if err != nil {
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(textView, "[red]Error:[white] %v\n", err)
		})
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			requirement := inputField.GetText()
			if strings.TrimSpace(requirement) == "" {
				return
			}
			inputField.SetText("")
			go func() {
				if err := soloAgent.Run(requirement); err != nil {
					app.QueueUpdateDraw(func() {
						fmt.Fprintf(textView, "[red]Error:[white] %v\n", err)
					})
				}
			}()
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 3, 0, true)

	app.SetRoot(flex, true)
}
