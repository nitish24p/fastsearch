package main

import (
	"fastsearch/pkg/pathfinder"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Fastsearch is in progress...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return

	}
	// Pattern / file name to search
	pattern := os.Args[1]
	fmt.Println(pattern)
	start := time.Now()

	/**
	  number of workers is no of go routines to spawn
	  to process file names

	  In general the algo is
	  1. Create a read stream of file names
	  2. Pass these file names to workers, the worker will do string matching
	  3. Merge all the processed responses in output channel
	*/
	const numWorkers = 10
	const maxParsers = 10

	inputStream := pathfinder.CreateDirectoryStream(homeDir, maxParsers)

	workerChannels := make([]<-chan string, numWorkers)

	// Dynamically create workers and store their channels
	for i := 0; i < numWorkers; i++ {
		workerChannels[i] = pathfinder.CreateWorklets(inputStream, i+1, pattern)
	}

	// Merge the workers
	resultStream := pathfinder.MergeWorkers(workerChannels...)

	// Collect results
	for result := range resultStream {
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("Found file: %s in %d ms \n", result, elapsed.Milliseconds())
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Search completed in %d ms\n", elapsed.Milliseconds())

}
