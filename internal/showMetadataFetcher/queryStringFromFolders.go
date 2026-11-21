package showmetadatafetcher

import (
	"os"
	"path/filepath"
)

//Map is "query: path" because when path is queried the map is looped through but when 'query' is queried it isn't.
//See "insertNewShows.go"

func GetQueriesForDirectory(path string) (map[string]string, error) {
	output := make(map[string]string)
	entries, err := os.ReadDir(path)
	if err != nil {
		return map[string]string{}, err
	}

	for _, e := range entries {
		entryPath := filepath.Join(path, e.Name())
		subEntries, err := os.ReadDir(entryPath)
		if err != nil {
			return map[string]string{}, err
		}
		allAreFiles := true
		for _, e2 := range subEntries {
			if e2.IsDir() {
				allAreFiles = false
				currPath := filepath.Join(entryPath, e2.Name())
				currQuery := e.Name() + e2.Name()
				output[currQuery] = currPath
			}
		}

		if allAreFiles {
			output[e.Name()] = entryPath
		}
	}

	return output, nil
}
