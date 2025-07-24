package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	openaai "github.com/sashabaranov/go-openai"
	"github.com/sgoal/tide/tool"
)

// SoloAgent is an agent that can work independently to build and deploy projects using ReAct framework.
type SoloAgent struct {
	client       *openaai.Client
	tools        map[string]tool.Tool
	maxLoops     int
	history      []openaai.ChatCompletionMessage
	logWriter    io.Writer
	systemPrompt string
}

// NewSoloAgent creates a new SoloAgent with ReAct framework.
func NewSoloAgent(logWriter io.Writer) (*SoloAgent, error) {
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

	systemPrompt := `You are an autonomous software development agent. You can think, plan, execute, and reflect on your actions.

Your workflow:
1. **THINK**: Analyze the task and break it down into smaller steps
2. **PLAN**: Create a clear plan with specific actions
3. **EXECUTE**: Use the available tools to implement each step
4. **OBSERVE**: Monitor the results of your actions
5. **REFLECT**: Evaluate progress and adjust your plan if needed

Available tools:
%s

Always follow this pattern:
- Start by understanding the task
- Create a detailed plan
- Execute step by step
- Verify each step
- Provide clear status updates
- Deploy the final result

Be autonomous and complete tasks from start to finish.`

	// Discover all available tools dynamically
	availableTools := map[string]tool.Tool{
		"deployer":    &tool.DeployerTool{},
		"code_writer": &tool.CodeWriterTool{},
		"file_editor": &tool.FileEditorTool{},
		"terminal":    &tool.TerminalTool{},
		"search":      &tool.SearchTool{},
	}

	// Build dynamic tool descriptions for system prompt
	toolDescriptions := ""
	for name, t := range availableTools {
		toolDescriptions += fmt.Sprintf("- %s: %s\n", name, t.Description())
	}

	// Update system prompt with dynamic tool list
	systemPrompt = fmt.Sprintf(systemPrompt, toolDescriptions)

	return &SoloAgent{
		client:       client,
		tools:        availableTools,
		maxLoops:     500,
		logWriter:    logWriter,
		systemPrompt: systemPrompt,
	}, nil
}

// Run runs the solo agent to complete the given task using ReAct framework.
func (a *SoloAgent) Run(task string) error {
	fmt.Fprintf(a.logWriter, "🚀 Solo Agent Starting...\n")
	fmt.Fprintf(a.logWriter, "📝 Task: %s\n", task)
	fmt.Fprintf(a.logWriter, "%s\n", strings.Repeat("=", 50))

	// Initialize history with system prompt
	a.history = []openaai.ChatCompletionMessage{
		{
			Role:    openaai.ChatMessageRoleSystem,
			Content: a.systemPrompt,
		},
		{
			Role:    openaai.ChatMessageRoleUser,
			Content: task,
		},
	}

	// Define available tools for the agent dynamically
	tools := []openaai.Tool{}
	for name, t := range a.tools {
		tools = append(tools, openaai.Tool{
			Type: openaai.ToolTypeFunction,
			Function: &openaai.FunctionDefinition{
				Name:        name,
				Description: t.Description(),
				Parameters:  getToolParameters(name),
			},
		})
	}

	// ReAct loop
	for i := 0; i < a.maxLoops; i++ {
		fmt.Fprintf(a.logWriter, "\n🔍 Loop %d/%d\n", i+1, a.maxLoops)

		req := openaai.ChatCompletionRequest{
			Model:    openaai.GPT4o20240806,
			Messages: a.history,
			Tools:    tools,
		}

		fmt.Fprintf(a.logWriter, "🤖 Thinking...\n")
		resp, err := a.client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			return fmt.Errorf("chat completion error: %w", err)
		}

		respMsg := resp.Choices[0].Message
		a.history = append(a.history, respMsg)

		// Display the agent's thought/plan
		if respMsg.Content != "" {
			fmt.Fprintf(a.logWriter, "💭 Agent Thought: %s\n", respMsg.Content)
		}

		// Check if we have a final answer
		if len(respMsg.ToolCalls) == 0 {
			fmt.Fprintf(a.logWriter, "\n✅ Task Completed!\n")
			fmt.Fprintf(a.logWriter, "📝 Final Result: %s\n", respMsg.Content)
			return nil
		}

		// Execute tool calls
		fmt.Fprintf(a.logWriter, "🔧 Executing %d tool(s)...\n", len(respMsg.ToolCalls))
		for _, toolCall := range respMsg.ToolCalls {
			fmt.Fprintf(a.logWriter, "\n📋 Tool: %s\n", toolCall.Function.Name)
			fmt.Fprintf(a.logWriter, "📄 Arguments: %s\n", toolCall.Function.Arguments)

			if tool, exists := a.tools[toolCall.Function.Name]; exists {
				observation, err := tool.Execute(json.RawMessage(toolCall.Function.Arguments))
				if err != nil {
					observation = fmt.Sprintf("❌ Error: %v", err)
				}
				if observation == "" {
					observation = "✅ Operation completed successfully"
				}

				fmt.Fprintf(a.logWriter, "👀 Observation: %s\n", observation)

				// Add tool result to history
				a.history = append(a.history, openaai.ChatCompletionMessage{
					Role:       openaai.ChatMessageRoleTool,
					ToolCallID: toolCall.ID,
					Name:       toolCall.Function.Name,
					Content:    observation,
				})
			} else {
				fmt.Fprintf(a.logWriter, "❌ Tool '%s' not found\n", toolCall.Function.Name)
			}
		}

		fmt.Fprintf(a.logWriter, "%s\n", strings.Repeat("=", 50))
	}

	return fmt.Errorf("⚠️ Maximum loops reached, task may not be fully completed")
}

// getToolParameters returns the JSON schema for tool parameters based on tool name
func getToolParameters(toolName string) json.RawMessage {
	parameters := map[string]json.RawMessage{
		"search": json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": {
					"type": "string",
					"description": "The search query."
				}
			},
			"required": ["query"]
		}`),
		"code_writer": json.RawMessage(`{
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
		"terminal": json.RawMessage(`{
			"type": "object",
			"properties": {
				"command": {
					"type": "string",
					"description": "The command to execute."
				}
			},
			"required": ["command"]
		}`),
		"file_editor": json.RawMessage(`{
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
		"deployer": json.RawMessage(`{
			"type": "object",
			"properties": {
				"project_path": {
					"type": "string",
					"description": "The path to the project directory to deploy."
				}
			},
			"required": ["project_path"]
		}`),
	}

	if params, exists := parameters[toolName]; exists {
		return params
	}

	// Default empty parameters for unknown tools
	return json.RawMessage(`{"type": "object", "properties": {}}`)
}
