package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

func HandleExecCommand(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.ExecCommandRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Command == "" {
		return errorResult("INVALID_REQUEST", "command cannot be empty")
	}

	cmd := exec.Command(input.Command, input.Args...)

	// Set environment variables
	if input.Env != nil {
		env := os.Environ()
		for key, value := range input.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	if input.Background {
		// Background execution
		sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
		processInfo := storage.CreateSession(sessionID, input.Command+" "+fmt.Sprintf("%v", input.Args))

		// Start command
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errorResult("EXEC_FAILED", "failed to create stdout pipe: "+err.Error())
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errorResult("EXEC_FAILED", "failed to create stderr pipe: "+err.Error())
		}

		if err := cmd.Start(); err != nil {
			return errorResult("EXEC_FAILED", "failed to start command: "+err.Error())
		}

		// Read output in goroutines
		go readPipe(stdout, func(data string) { processInfo.AppendStdout(data) })
		go readPipe(stderr, func(data string) { processInfo.AppendStderr(data) })

		// Wait for completion in goroutine
		go func() {
			exitCode := 0
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}
			processInfo.SetCompleted(exitCode)
		}()

		resp := models.ExecCommandResponse{
			SessionID: sessionID,
			Status:    "running",
		}
		return successResult(resp)
	}

	// Foreground execution
	output, err := cmd.CombinedOutput()

	stdout := string(output)
	stderr := ""
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return errorResult("EXEC_FAILED", "command failed: "+err.Error())
		}
	}

	resp := models.ExecCommandResponse{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
		Status:   "completed",
	}
	return successResult(resp)
}

func HandleManageProcess(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.ManageProcessRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	switch input.Action {
	case "list":
		sessions := storage.ListSessions()
		resp := models.ManageProcessResponse{
			Sessions: sessions,
			Message:  "Listed all sessions",
		}
		return successResult(resp)

	case "poll":
		if input.SessionID == "" {
			return errorResult("INVALID_REQUEST", "session_id required for poll action")
		}

		processInfo, err := storage.GetSession(input.SessionID)
		if err != nil {
			return errorResult("PROCESS_NOT_FOUND", err.Error())
		}

		snapshot := processInfo.GetSnapshot()
		resp := models.ManageProcessResponse{
			Sessions: []models.ProcessSession{snapshot},
			Message:  "Polled session",
		}
		return successResult(resp)

	case "send_keys":
		if input.SessionID == "" || input.Keys == "" {
			return errorResult("INVALID_REQUEST", "session_id and keys required for send_keys action")
		}
		return errorResult("INTERNAL_ERROR", "send_keys not fully implemented")

	case "kill":
		if input.SessionID == "" {
			return errorResult("INVALID_REQUEST", "session_id required for kill action")
		}
		return errorResult("INTERNAL_ERROR", "kill not fully implemented")

	default:
		return errorResult("INVALID_REQUEST", "unknown action: "+input.Action)
	}
}

func readPipe(rc io.ReadCloser, callback func(string)) {
	defer rc.Close()
	buf := make([]byte, 4096)
	for {
		n, err := rc.Read(buf)
		if n > 0 {
			callback(string(buf[:n]))
		}
		if err != nil {
			break
		}
	}
}
