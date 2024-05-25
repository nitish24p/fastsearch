package pathfinder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func CreateDirectoryStream(rootDir string, maxParsers int) <-chan string {
	out := make(chan string)

	var wg sync.WaitGroup
	visited := make(map[string]struct{})
	var mu sync.Mutex
	parsers := 1
	// Function to walk a directory and send entries to the output channel
	var walkDirectory func(dir string, parentDir string)
	walkDirectory = func(dir string, parentDir string) {
		defer wg.Done()
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					return nil
				}
				return err
			}
			if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
				// Skip parsing hidden files
				return filepath.SkipDir
			}

			mu.Lock()
			defer mu.Unlock()
			// Check if the path has already been visited
			if _, ok := visited[path]; ok {
				// Skip this visit
				// This is mainly for symlinks
				return nil
			}
			visited[path] = struct{}{}

			// If current found file is directory
			if d.IsDir() && path != dir {
				// If max parsers is not met
				// And current directory is at same level as parent directory
				// Spawn a go routine
				// THis is to make the current BFS faster
				if parsers < maxParsers && filepath.Dir(path) == parentDir {
					wg.Add(1)
					parsers++
					fmt.Println("Spawining New Go routine", path, parsers)
					go walkDirectory(path, path)

					// Since a new go routine is already parsing that directoru,
					// Stop current routine from parsing it again.
					return filepath.SkipDir
				} else {
					out <- path
				}
			} else {
				out <- path
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking directory:", err)
		}

	}

	wg.Add(1)
	go func() {
		// Also do cache directory traversal here

		walkDirectory(rootDir, rootDir)
		wg.Wait()
		close(out)
	}()

	return out
}
