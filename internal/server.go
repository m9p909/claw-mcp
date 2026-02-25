package internal

import (
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/tools"
)

type Server struct {
	mcpServer *mcp.Server
}

func NewServer() (*Server, error) {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "claw-mcp-server",
		Version: "1.0.0",
	}, nil)

	s := &Server{
		mcpServer: mcpServer,
	}

	if err := s.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	log.Println("MCP Server initialized with 8 tools")
	return s, nil
}

func (s *Server) registerTools() error {
	// Filesystem tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "read_file", Description: "Read the contents of a file"},
		tools.HandleReadFile)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "write_file", Description: "Write or create a file with content"},
		tools.HandleWriteFile)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "edit_file", Description: "Edit a file by replacing a range of lines (identified by hashes)"},
		tools.HandleEditFile)

	// Execution tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "exec_command", Description: "Execute a command (foreground or background)"},
		tools.HandleExecCommand)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "manage_process", Description: "Manage background processes (list, poll, send_keys, kill)"},
		tools.HandleManageProcess)

	// Memory tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "write_memory", Description: "Write to persistent memory"},
		tools.HandleWriteMemory)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "query_memory", Description: "Query memory with SQL SELECT"},
		tools.HandleQueryMemory)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "memory_search", Description: "Search memory with substring matching"},
		tools.HandleMemorySearch)

	return nil
}

func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
