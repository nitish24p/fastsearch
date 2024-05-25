package pathfinder

import "sync"

func MergeWorkers(worklets ...<-chan string) <-chan string {
	waitGroups := len(worklets)
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(waitGroups)
	for _, worklet := range worklets {
		go func(d <-chan string) {
			defer wg.Done()
			for val := range d {
				out <- val
			}

		}(worklet)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}
