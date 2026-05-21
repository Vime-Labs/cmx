package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Vime-Labs/cmx/internal/api"
)

const protocolVersion = "2024-11-05"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type initializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    serverCapabilities `json:"capabilities"`
	ServerInfo      serverInfo         `json:"serverInfo"`
}

type serverCapabilities struct {
	Tools *toolCapabilities `json:"tools,omitempty"`
}

type toolCapabilities struct {
	ListChanged bool `json:"listChanged"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type toolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema inputSchema `json:"inputSchema"`
}

type inputSchema struct {
	Type       string          `json:"type"`
	Properties json.RawMessage `json:"properties,omitempty"`
	Required   []string        `json:"required,omitempty"`
}

type callToolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolHandler func(api.API, json.RawMessage) callToolResult

type Server struct {
	api     api.API
	version string
	tools   map[string]toolHandler
	defs    []toolDefinition
}

func NewServer(apiClient api.API, version string) *Server {
	s := &Server{
		api:     apiClient,
		version: version,
		tools:   make(map[string]toolHandler),
	}
	s.registerTools()
	return s
}

func (s *Server) Run() error {
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 64*1024)

	for {
		body, err := readMessage(reader, &buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("reading message: %w", err)
		}

		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(nil, -32700, "Parse error")
			continue
		}

		resp := s.handle(&req)
		writeMessage(resp)
	}
}

func readMessage(r *bufio.Reader, buf *[]byte) ([]byte, error) {
	var length int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			n, err := strconv.Atoi(strings.TrimSpace(line[15:]))
			if err != nil {
				continue
			}
			length = n
		}
	}

	if length == 0 {
		return nil, fmt.Errorf("no Content-Length header")
	}

	if cap(*buf) < length {
		*buf = make([]byte, length)
	}
	*buf = (*buf)[:length]

	_, err := io.ReadFull(r, *buf)
	if err != nil {
		return nil, err
	}

	return *buf, nil
}

func writeMessage(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	os.Stdout.WriteString(header)
	os.Stdout.Write(data)
}

func writeError(id interface{}, code int, msg string) {
	writeMessage(rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcError{Code: code, Message: msg},
	})
}

func (s *Server) handle(req *rpcRequest) rpcResponse {
	var id interface{} = nil
	if len(req.ID) > 0 && string(req.ID) != "null" {
		json.Unmarshal(req.ID, &id)
	}

	switch req.Method {
	case "initialize":
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: initializeResult{
				ProtocolVersion: protocolVersion,
				Capabilities: serverCapabilities{
					Tools: &toolCapabilities{ListChanged: false},
				},
				ServerInfo: serverInfo{
					Name:    "cmx",
					Version: s.version,
				},
			},
		}

	case "tools/list":
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string][]toolDefinition{
				"tools": s.defs,
			},
		}

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{
				JSONRPC: "2.0",
				ID:      id,
				Error:   &rpcError{Code: -32602, Message: "Invalid params"},
			}
		}

		handler, ok := s.tools[params.Name]
		if !ok {
			return rpcResponse{
				JSONRPC: "2.0",
				ID:      id,
				Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Tool %q not found", params.Name)},
			}
		}

		result := handler(s.api, params.Arguments)
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result:  result,
		}

	case "notifications/initialized":
		return rpcResponse{JSONRPC: "2.0", ID: id, Result: map[string]interface{}{}}

	default:
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Method %q not found", req.Method)},
		}
	}
}

func (s *Server) add(def toolDefinition, fn toolHandler) {
	s.defs = append(s.defs, def)
	s.tools[def.Name] = fn
}

func str(s string) *string { return &s }
