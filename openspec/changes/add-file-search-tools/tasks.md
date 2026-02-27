## 1. Boyer-Moore String Search Implementation

- [x] 1.1 Create pkg/tools/filesearch/ package directory
- [x] 1.2 Implement boyer_moore.go with stringFinder struct (badCharSkip, goodSuffixSkip tables)
- [x] 1.3 Implement MakeStringFinder() to initialize skip tables
- [x] 1.4 Implement longestCommonSuffix() helper function
- [x] 1.5 Implement stringFinder.next() method for substring search
- [x] 1.6 Add unit tests for Boyer-Moore edge cases (empty text, pattern longer than text, single char)

## 2. Concurrent Search Engine

- [x] 2.1 Create search.go with SearchOptions struct (kind, lines, regex, finder)
- [x] 2.2 Implement searchJob struct and worker pool pattern
- [x] 2.3 Implement Search() orchestration function with configurable workers (default 128)
- [x] 2.4 Implement dirTraversal() for recursive directory walking
- [x] 2.5 Implement searchWorker() with binary file detection (NUL byte check)
- [x] 2.6 Add result collection with file path, line number, hash, and content
- [x] 2.7 Integrate hash.HashLine() for CRC32 line hashing on matches

## 3. Directory Operations

- [x] 3.1 Create directory.go for list/tree functionality
- [x] 3.2 Implement listDirectory() returning files with metadata (name, type, size, permissions)
- [x] 3.3 Implement findFiles() with glob pattern matching
- [x] 3.4 Implement treeDirectory() with ASCII box-drawing (├──, └──, │)
- [x] 3.5 Add depth limit support for tree generation
- [x] 3.6 Add symlink detection to prevent infinite loops

## 4. Request/Response Models

- [x] 4.1 Add SearchCodeRequest struct to pkg/models/models.go (path, query, regex bool, workers int)
- [x] 4.2 Add SearchCodeResponse struct with Results array (file, line, hash, content)
- [x] 4.3 Add FindFilesRequest struct (path, pattern string)
- [x] 4.4 Add FindFilesResponse struct with Files array (path, size, modified)
- [x] 4.5 Add ListDirectoryRequest struct (path string)
- [x] 4.6 Add ListDirectoryResponse struct with Entries array (name, type, size, permissions)
- [x] 4.7 Add TreeDirectoryRequest struct (path string, max_depth int)
- [x] 4.8 Add TreeDirectoryResponse struct (tree string)

## 5. MCP Tool Handlers

- [x] 5.1 Create pkg/tools/file_search.go
- [x] 5.2 Implement HandleSearchCode() calling filesearch.Search() with literal/regex mode
- [x] 5.3 Implement HandleFindFiles() calling filesearch.findFiles()
- [x] 5.4 Implement HandleListDirectory() calling filesearch.listDirectory()
- [x] 5.5 Implement HandleTreeDirectory() calling filesearch.treeDirectory()
- [x] 5.6 Add path validation and error handling for all handlers
- [x] 5.7 Add logging with request ID and duration tracking

## 6. Tool Registration

- [x] 6.1 Register search_code tool in internal/server.go registerTools()
- [x] 6.2 Register find_files tool in internal/server.go registerTools()
- [x] 6.3 Register list_directory tool in internal/server.go registerTools()
- [x] 6.4 Register tree_directory tool in internal/server.go registerTools()
- [x] 6.5 Update tool count in server initialization log from 22 to 26

## 7. Testing

- [x] 7.1 Write unit tests for Boyer-Moore stringFinder (pattern matching correctness)
- [x] 7.2 Write integration tests for search_file with literal and regex queries
- [x] 7.3 Write tests for hash consistency between search_file and read_file
- [x] 7.4 Write tests for binary file detection
- [x] 7.5 Write tests for find_files glob pattern matching
- [x] 7.6 Write tests for list_directory metadata accuracy
- [x] 7.7 Write tests for tree_directory ASCII formatting
- [x] 7.8 Write tests for path validation and error cases

## 8. Documentation and Validation

- [x] 8.1 Add tool descriptions with example requests/responses
- [x] 8.2 Test all 4 tools manually via MCP client
- [x] 8.3 Verify worker pool terminates without goroutine leaks
- [x] 8.4 Benchmark search performance on large codebases (>1000 files)
- [x] 8.5 Verify no external dependencies added to go.mod
