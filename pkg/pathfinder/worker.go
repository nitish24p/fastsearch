package pathfinder

import (
	"path/filepath"
	"strings"
)

func CreateWorklets(in <-chan string, id int, matchingString string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		matchingStringLower := strings.ToLower(matchingString)
		for path := range in {
			//fmt.Println("Found Match in %d worker %s", id)
			// if strings.Contains(filepath.Base(path), matchingString) {

			// 	out <- path
			// }
			if strings.Contains(strings.ToLower(filepath.Base(path)), matchingStringLower) {
				out <- path
			}
		}
	}()

	return out
}
