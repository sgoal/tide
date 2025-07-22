package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sgoal/tide/agent"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("AI Terminal IDE - Type 'exit' to quit")

	agent, err := agent.NewReActAgent()
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		os.Exit(1)
	}

	if err := agent.LoadHistory(); err != nil {
		fmt.Printf("Error loading conversation history: %v\n", err)
	}

	for _, msg := range agent.GetHistory() {
		fmt.Printf("%s: %s\n", msg.Role, msg.Content)
	}

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			if err := agent.SaveHistory(); err != nil {
				fmt.Printf("Error saving conversation history: %v\n", err)
			}
			fmt.Println("Exiting AI Terminal IDE.")
			break
		}

		response, err := agent.ProcessCommand(command)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		go func() {
			agent.SaveHistory()
		}()
		fmt.Println(response)
	}
}
