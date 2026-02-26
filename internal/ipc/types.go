package ipc

import (
	"encoding/json"

	"straw/internal/config"
)

// Method defines the type of request method.
type Method string

const (
	MethodGetStatus     Method = "GET_STATUS"
	MethodGetRules      Method = "GET_RULES"
	MethodAddRule       Method = "ADD_RULE"
	MethodUpdateRule    Method = "UPDATE_RULE"
	MethodDeleteRule    Method = "DELETE_RULE"
	MethodTriggerReload Method = "TRIGGER_RELOAD"
	MethodDryRun        Method = "DRY_RUN_REQUEST"
)

// DeleteRuleParams defines the parameters for deleting a rule.
type DeleteRuleParams struct {
	Name string `json:"name"`
}

// UpdateRuleParams defines the parameters for updating a rule.
type UpdateRuleParams struct {
	OriginalName string      `json:"original_name"`
	Rule         config.Rule `json:"rule"`
}

// Request represents a JSON command sent from client to daemon.
type Request struct {
	ID     string          `json:"id"`
	Method Method          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON response from daemon to client.
type Response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// EventType defines the type of asynchronous event.
type EventType string

const (
	EventNotification EventType = "EVENT_NOTIFICATION"
)

// Event represents an asynchronous message from daemon to client.
type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
