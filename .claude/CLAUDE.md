# Project Guidelines

## Testing Requirements
- Every requirement in openspec specs must have corresponding unit or integration tests
- Test/use case name must match the test name (e.g., requirement "User can login" → test "TestUserCanLogin" or "test_user_can_login")
- Tests should cover both happy path and edge cases
- Test coverage should validate each requirement is implemented correctly

## Code Style
- Encapsulate complexity in functions
- One function should be at most 9 if statements, or equivalent complexity
- Prefer minimalism - smallest amount of code to get the job done
- One class/struct/file should be at most 5 functions, or equivalent complexity
- Prefer strict types and validations
- All errors must be handled, at minimum with an error log
- Prefer functional tools like map, reduce over for loops

## Planning
- Be concise - prioritize conciseness over grammar
- Add a section for questions at the end, listing any unresolved questions
- Break down the plan into small steps

## Execution
- Tests should be run to verify implementation
- Provide a plan before implementation
- After completing a step, wait for input before continuing

## Naming
- Call the user by the name "coworker", when necessary