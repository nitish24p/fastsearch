// package app

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/fs"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"time"
// )

// // FileInfo holds the necessary file information
// type FileInfo struct {
// 	Path    string    `json:"path"`
// 	ModTime time.Time `json:"mod_time"`
// 	IsDir   bool      `json:"is_dir"`
// }

// // CacheData holds a map of file paths to FileInfo
// type CacheData struct {
// 	Files map[string]FileInfo `json:"files"`
// 	mu    sync.Mutex
// }

// // Walk directory, build cache, and search for pattern concurrently
// func cacheDirectoryStructureAndSearch(rootDir, cacheFile, pattern string, results chan<- string, done chan<- bool) {
// 	files := make(map[string]FileInfo)
// 	foundFiles := make(chan string)

// 	go func() {
// 		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
// 			if err != nil {
// 				if os.IsPermission(err) {
// 					//fmt.Printf("Skipping directory due to permission error: %s\n", path)
// 					return nil
// 				}
// 				return err
// 			}
// 			// Skip hidden directories
// 			if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
// 				return filepath.SkipDir
// 			}

// 			// Add file and directory information to the map
// 			info, err := d.Info()
// 			if err != nil {
// 				return err
// 			}
// 			files[path] = FileInfo{
// 				Path:    path,
// 				ModTime: info.ModTime(),
// 				IsDir:   d.IsDir(),
// 			}

// 			// Check if the current file matches the pattern
// 			if !d.IsDir() && strings.Contains(filepath.Base(path), pattern) {
// 				foundFiles <- path
// 			}

// 			return nil
// 		})

// 		if err != nil {
// 			fmt.Println("Error walking directory:", err)
// 		}

// 		// Serialize to JSON and write to cache file
// 		cacheData := CacheData{Files: files}
// 		cacheBytes, err := json.Marshal(cacheData)
// 		if err != nil {
// 			fmt.Println("Error marshalling cache:", err)
// 		}

// 		err = os.WriteFile(cacheFile, cacheBytes, 0644)
// 		if err != nil {
// 			fmt.Println("Error writing cache to file:", err)
// 		}

// 		close(foundFiles)
// 		done <- true
// 	}()

// 	for path := range foundFiles {
// 		results <- path
// 	}
// }

// // Load cache from file
// func loadCache(cacheFile string) (CacheData, error) {
// 	var cacheData CacheData

// 	data, err := os.ReadFile(cacheFile)
// 	if err != nil {
// 		return cacheData, err
// 	}

// 	err = json.Unmarshal(data, &cacheData)
// 	if err != nil {
// 		return cacheData, err
// 	}

// 	return cacheData, nil
// }

// // Check if cache is valid, remove invalid entries
// func isCacheValid(rootDir string, cacheData *CacheData) {
// 	cacheData.mu.Lock()
// 	defer cacheData.mu.Unlock()

// 	for path, file := range cacheData.Files {
// 		info, err := os.Stat(path)
// 		if os.IsNotExist(err) {
// 			delete(cacheData.Files, path) // Remove invalid entry
// 			continue
// 		} else if err != nil {
// 			fmt.Printf("Error checking file: %s\n", err)
// 			continue
// 		}
// 		if info.ModTime().After(file.ModTime) {
// 			delete(cacheData.Files, path) // Remove invalid entry
// 			fmt.Printf("File modified: %s\n", path)
// 			continue
// 		}
// 	}
// }

// // Search in cached data
// func searchInCache(files map[string]FileInfo, pattern string, results chan<- string) {
// 	for path, file := range files {
// 		if !file.IsDir && strings.Contains(filepath.Base(file.Path), pattern) {
// 			results <- path
// 		}
// 	}
// 	close(results)
// }

// // Search Non Cached Results
// func searchNonCachedResults(files map[string]FileInfo, pattern string, results chan<- string) {

// }

// // Main function
// func main() {
// 	cacheFile := "file_cache.json"
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Println("Error getting home directory:", err)
// 		return
// 	}

// 	pattern := "App.tsx"
// 	results := make(chan string)
// 	done := make(chan bool)

// 	// Check if cache file exists
// 	var cacheData CacheData
// 	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
// 		// Cache does not exist, build cache and search for pattern
// 		fmt.Println("Creating cache and searching for pattern...")
// 		go cacheDirectoryStructureAndSearch(homeDir, cacheFile, pattern, results, done)
// 	} else {
// 		// Load cache
// 		cacheData, err = loadCache(cacheFile)
// 		if err != nil {
// 			fmt.Println("Error loading cache:", err)
// 			return
// 		}

// 		// Check cache validity, remove invalid entries
// 		fmt.Println("Checking cache validity...")
// 		isCacheValid(homeDir, &cacheData)

// 		// Search in cached data
// 		go searchInCache(cacheData.Files, pattern, results)
// 		//go func() { done <- true }()
// 	}

// 	start := time.Now()

// 	for {
// 		select {
// 		case path := <-results:
// 			if path != "" {
// 				fmt.Printf("Found file: %s\n", path)
// 			}
// 		case <-done:
// 			// Rebuild cache for missed or invalid entries
// 			fmt.Println("Rebuilding missed or invalid cache entries...")
// 			go cacheDirectoryStructureAndSearch(homeDir, cacheFile, pattern, results, done)

// 			// Wait for cache to rebuild
// 			<-done

// 			// Reload cache after rebuilding
// 			cacheData, err = loadCache(cacheFile)
// 			if err != nil {
// 				fmt.Println("Error loading cache after rebuild:", err)
// 				return
// 			}

// 			// Search in the updated cache
// 			fmt.Println("Searching in updated cache...")
// 			go searchInCache(cacheData.Files, pattern, results)
// 		}

// 		// Check if all results are received
// 		if len(cacheData.Files) == len(results) {
// 			end := time.Now()
// 			elapsed := end.Sub(start)
// 			fmt.Printf("Search completed in %d ms\n", elapsed.Milliseconds())
// 			return
// 		}
// 	}
// }
