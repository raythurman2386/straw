# API Reference

Straw's IPC API allows programmatic interaction with the daemon.

## Overview

The daemon exposes a JSON-RPC API over Unix Domain Sockets (or named pipes on Windows). The TUI client uses this API to communicate with the daemon.

## Connection

### Socket Location

Default socket paths by platform:

| Platform | Default Path |
|----------|--------------|
| Linux | `~/.config/straw/straw.sock` or `/tmp/straw.sock` |
| macOS | `~/Library/Application Support/straw/straw.sock` |
| Windows | `\\.\pipe\straw` |

### Connecting

**Go Example:**
```go
conn, err := net.Dial("unix", "/tmp/straw.sock")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

**Python Example:**
```python
import socket

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
sock.connect("/tmp/straw.sock")
```

## Protocol

### JSON-RPC 2.0

All communication uses JSON-RPC 2.0 format:

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "MethodName",
  "params": {
    "param1": "value1"
  },
  "id": 1
}
```

**Success Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "key": "value"
  },
  "id": 1
}
```

**Error Response:**
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32600,
    "message": "Invalid Request"
  },
  "id": 1
}
```

## Methods

### GetStatus

Get the current daemon status.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "GetStatus",
  "id": 1
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "status": "running",
    "version": "1.0.0",
    "uptime": 3600,
    "watches": 3,
    "rules": 10
  },
  "id": 1
}
```

### GetEvents

Get recent file system events.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "GetEvents",
  "params": {
    "limit": 50,
    "since": "2024-01-01T00:00:00Z"
  },
  "id": 2
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "events": [
      {
        "timestamp": "2024-01-15T10:30:00Z",
        "type": "created",
        "path": "/home/user/Downloads/file.pdf",
        "rule": "Organize PDFs",
        "action": "move"
      }
    ]
  },
  "id": 2
}
```

### GetRules

Get all configured rules.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "GetRules",
  "id": 3
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "rules": [
      {
        "name": "Organize PDFs",
        "enabled": true,
        "matches": 42,
        "match": {
          "extension": ".pdf"
        },
        "actions": [
          {
            "type": "move",
            "target": "~/Documents/PDFs"
          }
        ]
      }
    ]
  },
  "id": 3
}
```

### AddRule

Add a new rule dynamically.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "AddRule",
  "params": {
    "name": "New Rule",
    "enabled": true,
    "match": {
      "extension": ".txt"
    },
    "actions": [
      {
        "type": "move",
        "target": "~/Documents/Text"
      }
    ]
  },
  "id": 4
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true,
    "ruleId": "rule-123"
  },
  "id": 4
}
```

### RemoveRule

Remove a rule by name or ID.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "RemoveRule",
  "params": {
    "name": "New Rule"
  },
  "id": 5
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true
  },
  "id": 5
}
```

### EnableRule / DisableRule

Enable or disable a rule.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "EnableRule",
  "params": {
    "name": "New Rule"
  },
  "id": 6
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true
  },
  "id": 6
}
```

### ReloadConfig

Reload the configuration from disk.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "ReloadConfig",
  "id": 7
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true,
    "message": "Configuration reloaded"
  },
  "id": 7
}
```

### GetWatches

Get all watched directories.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "GetWatches",
  "id": 8
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "watches": [
      {
        "path": "~/Downloads",
        "recursive": true,
        "active": true
      }
    ]
  },
  "id": 8
}
```

### TriggerAction

Manually trigger an action on a file.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "TriggerAction",
  "params": {
    "file": "/home/user/test.pdf",
    "action": {
      "type": "move",
      "target": "~/Documents"
    }
  },
  "id": 9
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true,
    "message": "Action completed"
  },
  "id": 9
}
```

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON |
| -32600 | Invalid Request | Invalid request object |
| -32601 | Method not found | Unknown method |
| -32602 | Invalid params | Invalid method parameters |
| -32603 | Internal error | Internal daemon error |
| -32000 | Config error | Configuration error |
| -32001 | Watch error | Filesystem watch error |
| -32002 | Action error | Action execution error |

## Example Client

**Go Client:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "net"
)

type Client struct {
    conn net.Conn
}

func NewClient(socketPath string) (*Client, error) {
    conn, err := net.Dial("unix", socketPath)
    if err != nil {
        return nil, err
    }
    return &Client{conn: conn}, nil
}

func (c *Client) Call(method string, params interface{}) (map[string]interface{}, error) {
    req := map[string]interface{}{
        "jsonrpc": "2.0",
        "method":  method,
        "params":  params,
        "id":      1,
    }
    
    encoder := json.NewEncoder(c.conn)
    if err := encoder.Encode(req); err != nil {
        return nil, err
    }
    
    var resp map[string]interface{}
    decoder := json.NewDecoder(c.conn)
    if err := decoder.Decode(&resp); err != nil {
        return nil, err
    }
    
    return resp, nil
}

func main() {
    client, err := NewClient("/tmp/straw.sock")
    if err != nil {
        panic(err)
    }
    defer client.conn.Close()
    
    result, err := client.Call("GetStatus", nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %+v\n", result)
}
```

## WebSocket Bridge

For browser-based clients, a WebSocket bridge can be used:

```go
// WebSocket to Unix Socket bridge
// Run this alongside the daemon to enable web UI
```

See `examples/websocket-bridge/` for a complete implementation.
