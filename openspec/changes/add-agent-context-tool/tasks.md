## 1. Documentation Creation

- [x] 1.1 Create `docs/` directory if it doesn't exist
- [x] 1.2 Write `docs/AGENT_CONTEXT.md` with imperative guidance covering Claw identity, personality, token efficiency, and skills system

## 2. Tool Handler Implementation

- [x] 2.1 Create `pkg/tools/agent_context.go` with `HandleGetAgentContext` function
- [x] 2.2 Implement handler to read `docs/AGENT_CONTEXT.md` and return content as text
- [x] 2.3 Add error handling for missing documentation file
- [x] 2.4 Add logging for context retrieval operations

## 3. Models and Types

- [x] 3.1 Add `GetAgentContextRequest` type to `pkg/models/` (empty struct, no parameters needed)
- [x] 3.2 Add `GetAgentContextResponse` type to `pkg/models/` with `Content` field

## 4. Server Registration
g
- [x] 4.1 Register `get_agent_context` tool in `internal/server.go` registerTools method
- [x] 4.2 Write tool description with essential context (under 500 chars, imperative tone)
- [x] 4.3 Update tool count in server initialization log message

## 5. Testing

- [x] 5.1 Create `pkg/tools/agent_context_test.go` with unit tests
- [x] 5.2 Test successful context retrieval
- [x] 5.3 Test error handling when documentation file is missing
- [x] 5.4 Verify tool description length is under 500 characters
- [ ] 5.5 Manual test: verify tool appears in MCP tool listings with correct description

## 6. Integration

- [x] 6.1 Build and run server locally
- [ ] 6.2 Connect MCP client and verify `get_agent_context` appears in tool list
- [ ] 6.3 Invoke tool and verify full documentation is returned
- [ ] 6.4 Verify documentation content covers all required areas (identity, personality, token efficiency, skills)
