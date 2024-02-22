package commons

import (
	"os"
	"path/filepath"
)

func isYAMLFile(filePath string) bool {
	st, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	if st.IsDir() {
		return false
	}

	ext := filepath.Ext(filePath)
	return ext == ".yaml" || ext == ".yml"
}

func isJSONFile(filePath string) bool {
	st, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	if st.IsDir() {
		return false
	}

	ext := filepath.Ext(filePath)
	return ext == ".json"
}
