package files

import (
	"os"
)

// ListDirectory returns a slice of FileInfo for the given path.
func ListDirectory(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		result = append(result, FileInfo{
			Name:     entry.Name(),
			Path:     path + string(os.PathSeparator) + entry.Name(),
			Size:     info.Size(),
			IsDir:    entry.IsDir(),
			Modified: info.ModTime(),
		})
	}
	return result, nil
}
