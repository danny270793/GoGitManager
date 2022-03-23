package directory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var KB float32 = 1024
var MB float32 = 1024 * KB
var GB float32 = 1024 * MB
var TB float32 = 1024 * GB

func SizeToString(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1000*1000 {
		return fmt.Sprintf("%0.2f KB", float32(size)/KB)
	} else if size < 1000*1000*1000 {
		return fmt.Sprintf("%0.2f MB", float32(size)/MB)
	} else if size < 1000*1000*1000*1000 {
		return fmt.Sprintf("%0.2f GB", float32(size)/GB)
	} else {
		return fmt.Sprintf("%0.2f TB", float32(size)/TB)
	}
}

func GetSizeAsString(basePath string) (string, error) {
	size, err := GetSize(basePath)
	if err != nil {
		return "", err
	}

	return SizeToString(size), nil
}

func GetSize(basePath string) (int64, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%s", err))
	}

	var totalSize int64 = 0
	for _, f := range files {
		fullPath := path.Join(basePath, f.Name())
		if f.IsDir() {
			size, err2 := GetSize(fullPath)
			if err2 != nil {
				return 0, err2
			}
			totalSize += size
		} else {
			totalSize += f.Size()
		}
	}

	return totalSize, nil
}

func IsDirectory(path string) (bool, error) {
	dir, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, errors.New(fmt.Sprintf("%s does not exist", path))
	}

	if err != nil {
		return false, err
	}
	return dir.IsDir(), nil
}
