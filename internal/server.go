package internal

import (
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/browser"
	browsertools "awesomeProject/pkg/browser/tools"
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

	// Initialize BrowserManager singleton
	_ = browser.NewBrowserManager()

	log.Println("MCP Server initialized with 22 tools")
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

	// Browser automation tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_navigate", Description: "Navigate to URL"},
		browsertools.HandleBrowserNavigate)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_snapshot", Description: "Get accessibility snapshot of page"},
		browsertools.HandleBrowserSnapshot)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_click", Description: "Click element on page"},
		browsertools.HandleBrowserClick)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_type", Description: "Type text into element"},
		browsertools.HandleBrowserType)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_fill_form", Description: "Fill multiple form fields"},
		browsertools.HandleBrowserFillForm)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_select_option", Description: "Select option from dropdown"},
		browsertools.HandleBrowserSelectOption)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_press_key", Description: "Press keyboard key"},
		browsertools.HandleBrowserPressKey)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_wait_for", Description: "Wait for text, element, or time"},
		browsertools.HandleBrowserWaitFor)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_handle_dialog", Description: "Handle JavaScript dialogs"},
		browsertools.HandleBrowserHandleDialog)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_navigate_back", Description: "Navigate backward in history"},
		browsertools.HandleBrowserNavigateBack)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_hover", Description: "Hover over element"},
		browsertools.HandleBrowserHover)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "browser_close", Description: "Close browser"},
		browsertools.HandleBrowserClose)

	// Agent skills tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "list_skills", Description: "List all available Agent Skills"},
		tools.HandleListSkills)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "get_skill", Description: "Retrieve a specific Agent Skill with full content"},
		tools.HandleGetSkill)

	return nil
}

func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
