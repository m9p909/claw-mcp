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

	log.Println("MCP Server initialized with 27 tools")
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

	// Agent context tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "get_agent_context", Description: "Claw is a personal agent MCP server in real Linux. Be professional, concise, token-efficient. Skills at ~/.mcpclaw/skills/ (use list_skills/get_skill). Call for full guide."},
		tools.HandleGetAgentContext)

	// File search tools
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "search_file", Description: "Search files with literal or regex patterns, returns results with line hashes"},
		tools.HandleSearchFile)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "find_files", Description: "Find files matching glob patterns recursively"},
		tools.HandleFindFiles)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "list_directory", Description: "List directory contents with metadata"},
		tools.HandleListDirectory)
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{Name: "tree_directory", Description: "Generate ASCII tree visualization of directory structure"},
		tools.HandleTreeDirectory)

	return nil
}

func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
