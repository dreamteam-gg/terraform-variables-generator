package utils

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// GetAllFiles will get all files in current directory
func GetAllFiles(dir string, ext string) ([]string, error) {
	var files []string
	log.Infof("Finding files in %q directory", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		CheckError(err)
	}
	files, err := filepath.Glob(path.Join(dir, ext))
	CheckError(err)

	if len(files) == 0 {
		log.Infof("No files with %q extensions found in %q", ext, dir)
	}
	return files, nil
}

// FileExists checks if file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}
