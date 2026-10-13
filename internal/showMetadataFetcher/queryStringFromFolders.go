package showmetadatafetcher

import (
	"os"
	"path/filepath"
)

type queryPathPair struct {
	path  string
	query string
}

func GetQueriesForDirectory(path string) ([]queryPathPair, error) {
	var output []queryPathPair
	entries, err := os.ReadDir(path)
	if err != nil {
		return []queryPathPair{}, err
	}

	for _, e := range entries {
		entryPath := filepath.Join(path, e.Name())
		subEntries, err := os.ReadDir(entryPath)
		if err != nil {
			return []queryPathPair{}, err
		}
		allAreFiles := true
		for _, e2 := range subEntries {
			if e2.IsDir() {
				allAreFiles = false
				QPPair := queryPathPair{}
				QPPair.path = filepath.Join(entryPath, e2.Name())
				QPPair.query = e.Name() + e2.Name()
				output = append(output, QPPair)
			}
		}

		if allAreFiles {
			output = append(output, queryPathPair{entryPath, e.Name()})
		}
	}

	return output, nil
}
