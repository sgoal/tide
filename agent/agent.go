package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	openaai "github.com/sashabaranov/go-openai"
	"github.com/sgoal/tide/tool"
)

// ReActAgent is an agent that uses the ReAct framework to accomplish tasks.
type ReActAgent struct {
	client    *openaai.Client
	tools     map[string]tool.Tool
	maxLoops  int
	history   []openaai.ChatCompletionMessage
	logWriter io.Writer
}

const historyFilePath = "conversation_history.json"

func (a *ReActAgent) SaveHistory() error {
	data, err := json.Marshal(a.history)
	if err != nil {
		return err
	}
	return os.WriteFile(historyFilePath, data, 0644)
}

func (a *ReActAgent) LoadHistory() error {
	if _, err := os.Stat(historyFilePath); os.IsNotExist(err) {
		return nil // No history file yet
	}
	data, err := os.ReadFile(historyFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.history)
}

func NewReActAgent(logWriter io.Writer) (*ReActAgent, error) {
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

	if logWriter == nil {
		logWriter = io.Discard
	}

	return &ReActAgent{
		client: client,
		tools: map[string]tool.Tool{
			"code_writer": &tool.CodeWriterTool{},
			"file_editor": &tool.FileEditorTool{},
			"terminal":    &tool.TerminalTool{},
			"search":      &tool.SearchTool{},
		},
		maxLoops:  10,
		logWriter: logWriter,
	}, nil
}

func getTools() []openaai.Tool {
	tools := []openaai.Tool{
		{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        "search",
				Description: "A tool for searching the web using DuckDuckGo.",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The search query."
						}
					},
					"required": ["query"]
				}`),
			},
		},
		{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        "code_writer",
				Description: "A tool for writing code to a file.",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"dir_path": {
							"type": "string",
							"description": "The directory path to write the file to."
						},
						"file_name": {
							"type": "string",
							"description": "The name of the file to write."
						},
						"code": {
							"type": "string",
							"description": "The code to write to the file."
						}
					},
					"required": ["dir_path", "file_name", "code"]
				}`),
			},
		},
		{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        "terminal",
				Description: "Executes shell commands. Use this to run scripts, execute programs, or perform any other command-line operations. For example, to run a python script, you would use 'python your_script.py'.",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"command": {
							"type": "string",
							"description": "The command to execute."
						}
					},
					"required": ["command"]
				}`),
			},
		},
		{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        "file_editor",
				Description: "A tool for reading files from a directory, modifying their content, and writing them back.",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"dir_path": {
							"type": "string",
							"description": "The directory path containing the file to modify."
						},
						"file_name": {
							"type": "string",
							"description": "The name of the file to modify."
						},
						"search_text": {
							"type": "string",
							"description": "The text to search for in the file content."
						},
						"replace_text": {
							"type": "string",
							"description": "The text to replace the searched content with."
						}
					},
					"required": ["dir_path", "file_name", "search_text", "replace_text"]
				}`),
			},
		},
	}

	return tools
}

// ProcessCommand processes a command using the ReAct framework with native tool calling.
func (a *ReActAgent) ProcessCommand(command string) (string, error) {
	a.history = append(a.history, openaai.ChatCompletionMessage{
		Role:    openaai.ChatMessageRoleUser,
		Content: command,
	})
	tools := getTools()
	for i := 0; i < a.maxLoops; i++ {
		req := openaai.ChatCompletionRequest{
			Model:    openaai.GPT4o20240806,
			Messages: a.history,
			Tools:    tools,
		}

		fmt.Fprintln(a.logWriter, "--- Sending request to OpenAI ---")
		resp, err := a.client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			return "", fmt.Errorf("chat completion error: %w", err)
		}

		respMsg := resp.Choices[0].Message
		a.history = append(a.history, respMsg)

		if respMsg.ToolCalls == nil {
			fmt.Fprintln(a.logWriter, "--- Received final answer ---")
			return respMsg.Content, nil
		}

		fmt.Fprintf(a.logWriter, "--- Received tool call: %s ---\n", respMsg.ToolCalls[0].Function.Name)
		for _, toolCall := range respMsg.ToolCalls {
			if tool, exists := a.tools[toolCall.Function.Name]; exists {
				fmt.Fprintf(a.logWriter, "Executing tool: %s with args: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
				observation, err := tool.Execute(json.RawMessage(toolCall.Function.Arguments))
				if err != nil {
					observation = fmt.Sprintf("Error executing tool: %v", err)
				}
				if observation == "" {
					observation = "No result found."
				}
				fmt.Fprintf(a.logWriter, "Observation: %s\n", observation)
				a.history = append(a.history, openaai.ChatCompletionMessage{
					Role:       openaai.ChatMessageRoleTool,
					ToolCallID: toolCall.ID,
					Name:       toolCall.Function.Name,
					Content:    observation,
				})
			} else {
				fmt.Fprintf(a.logWriter, "Tool '%s' not found.\n", toolCall.Function.Name)
			}
		}
	}

	return "", fmt.Errorf("max loops reached")
}

func (a *ReActAgent) GetHistory() []openaai.ChatCompletionMessage {
	return a.history
}
