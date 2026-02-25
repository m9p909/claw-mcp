package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

func HandleExecCommand(ctx context.Context, req *mcp.CallToolRequest, input models.ExecCommandRequest) (*mcp.CallToolResult, models.ExecCommandResponse, error) {
	logger := pkglog.NewLogger()

	if input.Command == "" {
		return errorResult(ctx, "INVALID_REQUEST", "command cannot be empty"), models.ExecCommandResponse{}, nil
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
		start := time.Now()
		sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
		processInfo := storage.CreateSession(sessionID, input.Command+" "+fmt.Sprintf("%v", input.Args))

		logger.Info(ctx, "Starting background command",
			"session_id", sessionID,
			"command", input.Command,
			"args", fmt.Sprintf("%v", input.Args))

		// Start command
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errorResult(ctx, "EXEC_FAILED", "failed to create stdout pipe: "+err.Error()), models.ExecCommandResponse{}, nil
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errorResult(ctx, "EXEC_FAILED", "failed to create stderr pipe: "+err.Error()), models.ExecCommandResponse{}, nil
		}

		if err := cmd.Start(); err != nil {
			return errorResult(ctx, "EXEC_FAILED", "failed to start command: "+err.Error()), models.ExecCommandResponse{}, nil
		}

		// Read output in goroutines
		go func() {
			logger.Debug(ctx, "Background goroutine started", "session_id", sessionID, "stream", "stdout")
			readPipe(stdout, func(data string) {
				logger.Trace(ctx, "Received stdout chunk", "session_id", sessionID, "size", len(data))
				processInfo.AppendStdout(data)
			})
		}()

		go func() {
			logger.Debug(ctx, "Background goroutine started", "session_id", sessionID, "stream", "stderr")
			readPipe(stderr, func(data string) {
				logger.Trace(ctx, "Received stderr chunk", "session_id", sessionID, "size", len(data))
				processInfo.AppendStderr(data)
			})
		}()

		// Wait for completion in goroutine
		go func() {
			exitCode := 0
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}
			processInfo.SetCompleted(exitCode)

			logger.Info(ctx, "Background command completed",
				"session_id", sessionID,
				"exit_code", exitCode,
				pkglog.Duration(time.Since(start)))
		}()

		resp := models.ExecCommandResponse{
			SessionID: sessionID,
			Status:    "running",
		}
		return nil, resp, nil
	}

	// Foreground execution
	start := time.Now()
	logger.Info(ctx, "Starting foreground command",
		"command", input.Command,
		"args", fmt.Sprintf("%v", input.Args))

	output, err := cmd.CombinedOutput()

	stdout := string(output)
	stderr := ""
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			logger.Warn(ctx, "Command exited with non-zero code",
				"exit_code", exitCode,
				"stdout_size", len(stdout),
				pkglog.Duration(time.Since(start)))
		} else {
			return errorResult(ctx, "EXEC_FAILED", "command failed: "+err.Error()), models.ExecCommandResponse{}, nil
		}
	} else {
		logger.Info(ctx, "Foreground command completed",
			"exit_code", exitCode,
			"stdout_size", len(stdout),
			"stderr_size", len(stderr),
			pkglog.Duration(time.Since(start)))
	}

	resp := models.ExecCommandResponse{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
		Status:   "completed",
	}
	return nil, resp, nil
}

func HandleManageProcess(ctx context.Context, req *mcp.CallToolRequest, input models.ManageProcessRequest) (*mcp.CallToolResult, models.ManageProcessResponse, error) {
	logger := pkglog.NewLogger()

	switch input.Action {
	case "list":
		logger.Info(ctx, "Process management action", "action", "list")
		sessions := storage.ListSessions()
		logger.Debug(ctx, "Process list result", "session_count", len(sessions))

		resp := models.ManageProcessResponse{
			Sessions: sessions,
			Message:  "Listed all sessions",
		}
		return nil, resp, nil

	case "poll":
		if input.SessionID == "" {
			return errorResult(ctx, "INVALID_REQUEST", "session_id required for poll action"), models.ManageProcessResponse{}, nil
		}

		logger.Info(ctx, "Process management action", "action", "poll", "session_id", input.SessionID)

		processInfo, err := storage.GetSession(input.SessionID)
		if err != nil {
			return errorResult(ctx, "PROCESS_NOT_FOUND", err.Error()), models.ManageProcessResponse{}, nil
		}

		snapshot := processInfo.GetSnapshot()
		logger.Debug(ctx, "Poll result", "session_id", input.SessionID,
			"status", snapshot.Status,
			"exit_code", snapshot.ExitCode)

		resp := models.ManageProcessResponse{
			Sessions: []models.ProcessSession{snapshot},
			Message:  "Polled session",
		}
		return nil, resp, nil

	case "send_keys":
		if input.SessionID == "" || input.Keys == "" {
			return errorResult(ctx, "INVALID_REQUEST", "session_id and keys required for send_keys action"), models.ManageProcessResponse{}, nil
		}
		return errorResult(ctx, "INTERNAL_ERROR", "send_keys not fully implemented"), models.ManageProcessResponse{}, nil

	case "kill":
		if input.SessionID == "" {
			return errorResult(ctx, "INVALID_REQUEST", "session_id required for kill action"), models.ManageProcessResponse{}, nil
		}
		return errorResult(ctx, "INTERNAL_ERROR", "kill not fully implemented"), models.ManageProcessResponse{}, nil

	default:
		return errorResult(ctx, "INVALID_REQUEST", "unknown action: "+input.Action), models.ManageProcessResponse{}, nil
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
