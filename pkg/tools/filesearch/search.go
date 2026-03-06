package filesearch

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

const (
	LITERAL = iota
	REGEX
)

type SearchOptions struct {
	Kind   int
	Regex  *regexp.Regexp
	Finder *stringFinder
}

type SearchResult struct {
	File    string
	Line    int
	Content string
}

type searchJob struct {
	path string
	opts *SearchOptions
}

// Search performs concurrent text search across files.
// Workers defaults to 128 if not specified (0 or negative).
func Search(paths []string, opts *SearchOptions, workers int, results chan<- SearchResult) error {
	if workers <= 0 {
		workers = 128
	}

	searchJobs := make(chan *searchJob, workers)
	var wg sync.WaitGroup

	// Start worker pool
	for w := 0; w < workers; w++ {
		go searchWorker(searchJobs, &wg, results)
	}

	// Traverse and queue jobs
	for _, path := range paths {
		if err := dirTraversal(path, opts, searchJobs, &wg); err != nil {
			close(searchJobs)
			wg.Wait()
			return err
		}
	}

	close(searchJobs)
	wg.Wait()
	return nil
}

func dirTraversal(path string, opts *SearchOptions, searchJobs chan *searchJob, wg *sync.WaitGroup) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("couldn't lstat path %s: %w", path, err)
	}

	if !info.IsDir() {
		wg.Add(1)
		searchJobs <- &searchJob{path, opts}
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open path %s: %w", path, err)
	}
	defer f.Close()

	dirNames, err := f.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("couldn't read dir names for path %s: %w", path, err)
	}

	for _, deeperPath := range dirNames {
		if err := dirTraversal(filepath.Join(path, deeperPath), opts, searchJobs, wg); err != nil {
			return err
		}
	}
	return nil
}

func searchWorker(jobs chan *searchJob, wg *sync.WaitGroup, results chan<- SearchResult) {
	for job := range jobs {
		processFile(job, results)
		wg.Done()
	}
}

func processFile(job *searchJob, results chan<- SearchResult) {
	f, err := os.Open(job.path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	isBinary := false
	lineNum := 1

	for scanner.Scan() {
		text := scanner.Bytes()

		// Check first buffer for NUL byte
		if lineNum == 1 {
			isBinary = bytes.IndexByte(text, 0) != -1
		}

		matched := false
		if job.opts.Kind == LITERAL {
			matched = job.opts.Finder.next(text) != -1
		} else if job.opts.Kind == REGEX {
			matched = job.opts.Regex.Match(text)
		}

		if matched {
			if isBinary {
				results <- SearchResult{
					File:    job.path,
					Line:    0,
					Content: fmt.Sprintf("Binary file %s matches", job.path),
				}
				return
			}

			results <- SearchResult{
				File:    job.path,
				Line:    lineNum,
				Content: string(text),
			}
		}
		lineNum++
	}
}
