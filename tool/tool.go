package tool

import "encoding/json"

// Tool defines the interface for a tool that the agent can use.
type Tool interface {
	Name() string
	Description() string
	Execute(args json.RawMessage) (string, error)
}