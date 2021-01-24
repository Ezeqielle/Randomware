package file

import (
	"os"
	"path/filepath"
)

// Discover : returns array of strings of either files or folders in a root directory
func Discover(root string) ([]string, []string, error) {
	var files []string
	var folders []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		} else {
			folders = append(folders, path)
		}
		return nil
	})
	return files, folders, err
}

// DiscoverFiles : returns array of strings of  files or folders in a root directory
func DiscoverFiles(root string) ([]string, error) {
	files, _, err := Discover(root)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// DiscoverFolders : returns array of strings of  files or folders in a root directory
func DiscoverFolders(root string) ([]string, error) {
	_, folders, err := Discover(root)
	if err != nil {
		return nil, err
	}
	return folders, nil
}
