package gofile

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// IsExists determine if a file or folder exists
func IsExists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// GetRunPath get the absolute path of the program execution
func GetRunPath() string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}

	return filepath.Dir(dir)
}

// GetFilename get file name
func GetFilename(filePath string) string {
	_, name := filepath.Split(filePath)
	return name
}

// GetFileDir get dir
func GetFileDir(filePath string) string {
	dir, _ := filepath.Split(filePath)
	return dir
}

// CreateDir create dir
func CreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0766)
	}
	return nil
}

// GetFilenameWithoutSuffix get file name without suffix
func GetFilenameWithoutSuffix(filePath string) string {
	_, name := filepath.Split(filePath)

	return strings.TrimSuffix(name, path.Ext(name))
}

// Join joins any number of path elements into a single path
func Join(elem ...string) string {
	dir := strings.Join(elem, "/")

	if IsWindows() {
		return strings.ReplaceAll(dir, "/", "\\")
	}

	return dir
}

// IsWindows determining whether a window environment
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// GetPathDelimiter get separator by system type
func GetPathDelimiter() string {
	delimiter := "/"
	if IsWindows() {
		delimiter = "\\"
	}

	return delimiter
}

// ListFiles iterates over all files in the specified directory, returning the absolute path to the file
func ListFiles(dirPath string, opts ...Option) ([]string, error) {
	files := []string{}
	err := error(nil)

	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return files, err
	}

	o := defaultOptions()
	o.apply(opts...)

	switch o.filter {
	case prefix:
		return files, walkDirWithFilter(dirPath, &files, matchPrefix(o.name))
	case suffix:
		return files, walkDirWithFilter(dirPath, &files, matchSuffix(o.name))
	case contain:
		return files, walkDirWithFilter(dirPath, &files, matchContain(o.name))
	}

	return files, walkDir(dirPath, &files)
}

// ListDirsAndFiles iterates through all subdirectories of the specified directory, returning the absolute path to the file
func ListDirsAndFiles(dirPath string) (map[string][]string, error) {
	df := make(map[string][]string, 2)

	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return df, err
	}

	dirs := []string{}
	files := []string{}
	err = walkDir2(dirPath, &dirs, &files)
	if err != nil {
		return df, err
	}

	df["dirs"] = dirs
	df["files"] = files

	return df, nil
}

// FuzzyMatchFiles fuzzy matching of documents, * only
func FuzzyMatchFiles(f string) []string {
	var files []string
	dir, filenameReg := filepath.Split(f)
	if !strings.Contains(filenameReg, "*") {
		files = append(files, f)
		return files
	}

	lFiles, err := ListFiles(dir)
	if err != nil {
		return files
	}
	for _, file := range lFiles {
		_, filename := filepath.Split(file)
		isMatch, _ := path.Match(filenameReg, filename)
		if isMatch {
			files = append(files, file)
		}
	}

	return files
}

// iterative traversal of documents with filter conditions
func walkDirWithFilter(dirPath string, allFiles *[]string, filter filterFn) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			err = walkDirWithFilter(deepFile, allFiles, filter)
			if err != nil {
				return err
			}
			continue
		}
		if filter(deepFile) {
			*allFiles = append(*allFiles, deepFile)
		}
	}

	return nil
}

func walkDir2(dirPath string, allDirs *[]string, allFiles *[]string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			*allDirs = append(*allDirs, deepFile)
			err = walkDir2(deepFile, allDirs, allFiles)
			if err != nil {
				return err
			}
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

type filterFn func(string) bool

// suffix matching
func matchSuffix(suffixName string) filterFn {
	return func(filename string) bool {
		if suffixName == "" {
			return false
		}

		size := len(filename) - len(suffixName)
		if size >= 0 && filename[size:] == suffixName {
			return true
		}
		return false
	}
}

// prefix Matching
func matchPrefix(prefixName string) filterFn {
	return func(filePath string) bool {
		if prefixName == "" {
			return false
		}
		filename := GetFilename(filePath)
		size := len(filename) - len(prefixName)
		if size >= 0 && filename[:len(prefixName)] == prefixName {
			return true
		}
		return false
	}
}

// contains the string
func matchContain(containName string) filterFn {
	return func(filePath string) bool {
		if containName == "" {
			return false
		}
		filename := GetFilename(filePath)
		return strings.Contains(filename, containName)
	}
}

// traversing the document by iteration
func walkDir(dirPath string, allFiles *[]string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			err = walkDir(deepFile, allFiles)
			if err != nil {
				return err
			}
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}
