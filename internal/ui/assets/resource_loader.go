package assets

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

func LoadResource(relativePath string) fyne.Resource {
	for _, candidate := range candidatePaths(relativePath) {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}
		resource, err := fyne.LoadResourceFromPath(candidate)
		if err == nil {
			return resource
		}
	}
	return nil
}

func candidatePaths(relativePath string) []string {
	seen := make(map[string]struct{})
	var paths []string

	add := func(path string) {
		clean := filepath.Clean(path)
		if _, ok := seen[clean]; ok {
			return
		}
		seen[clean] = struct{}{}
		paths = append(paths, clean)
	}

	add(relativePath)
	add(filepath.Join("..", relativePath))
	add(filepath.Join("..", "..", relativePath))

	if wd, err := os.Getwd(); err == nil {
		add(filepath.Join(wd, relativePath))
		add(filepath.Join(wd, "..", relativePath))
		add(filepath.Join(wd, "..", "..", relativePath))
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		add(filepath.Join(exeDir, relativePath))
		add(filepath.Join(exeDir, "..", relativePath))
		add(filepath.Join(exeDir, "..", "..", relativePath))
	}

	return paths
}
