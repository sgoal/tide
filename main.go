package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openaai "github.com/sashabaranov/go-openai"
)

// Tool defines the interface for a tool that the agent can use.
type Tool interface {
	Name() string
	Description() string
	Execute(args string) (string, error)
}

// CodeWriterTool is a tool for writing code to a file.
type CodeWriterTool struct{}

func (t *CodeWriterTool) Name() string {
	return "code_writer"
}

func (t *CodeWriterTool) Description() string {
	return "A tool for writing code to a file. The input should be a JSON object with 'filepath' and 'code' keys."
}

// Execute expects args to be a JSON string with "filepath" and "code"
func (t *CodeWriterTool) Execute(args string) (string, error) {
	var params struct {
		Filepath string `json:"filepath"`
		Code     string `json:"code"`
	}
	err := json.Unmarshal([]byte(args), &params)
	if err != nil {
		return "", fmt.Errorf("invalid arguments for code_writer tool: %w", err)
	}

	err = os.WriteFile(params.Filepath, []byte(params.Code), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully wrote code to %s", params.Filepath), nil
}

// ReActAgent is an agent that uses the ReAct framework to accomplish tasks.
type ReActAgent struct {
	client   *openaai.Client
	tools    map[string]Tool
	maxLoops int
}

func NewReActAgent() (*ReActAgent, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	var client *openaai.Client
	azureEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	if azureEndpoint != "" {
		config := openaai.DefaultAzureConfig(openAIKey, azureEndpoint)
		deploymentName := os.Getenv("AZURE_OPENAI_DEPLOYMENT")
		if deploymentName != "" {
			config.AzureModelMapperFunc = func(model string) string {
				return deploymentName
			}
		}
		client = openaai.NewClientWithConfig(config)
	} else {
		config := openaai.DefaultConfig(openAIKey)
		baseURL := os.Getenv("OPENAI_BASE_URL")
		if baseURL != "" {
			config.BaseURL = baseURL
		}
		client = openaai.NewClientWithConfig(config)
	}

	return &ReActAgent{
		client: client,
		tools: map[string]Tool{
			"code_writer": &CodeWriterTool{},
		},
		maxLoops: 10,
	}, nil
}

// ProcessCommand processes a command using the ReAct framework with native tool calling.
func (a *ReActAgent) ProcessCommand(command string) (string, error) {
	messages := []openaai.ChatCompletionMessage{
		{
			Role:    openaai.ChatMessageRoleUser,
			Content: command,
		},
	}

	tools := []openaai.Tool{
		{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        "code_writer",
				Description: "A tool for writing code to a file.",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"filepath": {
							"type": "string",
							"description": "The path to the file to write."
						},
						"code": {
							"type": "string",
							"description": "The code to write to the file."
						}
					},
					"required": ["filepath", "code"]
				}`),
			},
		},
	}

	for i := 0; i < a.maxLoops; i++ {
		req := openaai.ChatCompletionRequest{
			Model:    openaai.GPT4o20240806,
			Messages: messages,
			Tools:    tools,
		}

		fmt.Println("--- Sending request to OpenAI ---")
		resp, err := a.client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			return "", fmt.Errorf("chat completion error: %w", err)
		}

		respMsg := resp.Choices[0].Message
		messages = append(messages, respMsg)

		if respMsg.ToolCalls == nil {
			fmt.Println("--- Received final answer ---")
			return respMsg.Content, nil
		}

		fmt.Printf("--- Received tool call: %s ---\n", respMsg.ToolCalls[0].Function.Name)
		for _, toolCall := range respMsg.ToolCalls {
			if tool, exists := a.tools[toolCall.Function.Name]; exists {
				fmt.Printf("Executing tool: %s with args: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
				observation, err := tool.Execute(toolCall.Function.Arguments)
				if err != nil {
					observation = fmt.Sprintf("Error executing tool: %v", err)
				}
				fmt.Printf("Observation: %s\n", observation)
				messages = append(messages, openaai.ChatCompletionMessage{
					Role:       openaai.ChatMessageRoleTool,
					ToolCallID: toolCall.ID,
					Name:       toolCall.Function.Name,
					Content:    observation,
				})
			} else {
				fmt.Printf("Tool '%s' not found.\n", toolCall.Function.Name)
			}
		}
	}

	return "", fmt.Errorf("max loops reached")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("AI Terminal IDE - Type 'exit' to quit")

	agent, err := NewReActAgent()
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		os.Exit(1)
	}

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		response, err := agent.ProcessCommand(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(response)
	}
}
