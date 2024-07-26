package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func ValidFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.New("path does not exist")
	}
	if info.IsDir() {
		return errors.New("path is a directory, not a file")
	}
	return nil
}

func ValidDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.New("directory does not exist")
	}
	if !info.IsDir() {
		return errors.New("path is not a directory")
	}
	return nil
}

func SanitizeURL(url string) string { // FIXME: Change function name @quality
	return strings.ReplaceAll(url, " ", "%20")
}

func GetRandomColor() string {
	colors := []string{
		"#F3CA52",
		"#7ABA78",
		"#F6E9B2",
		"#A0153E",
		"#00224D",
		"#57A6A1",
	}
	rand.Seed(time.Now().UnixNano()) // TODO: Replace deprecated function @code
	return colors[rand.Intn(len(colors))]
}

func OverwriteFile(content string, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err1 := f.WriteString(content)
	if err1 != nil {
		panic(err1)
	}
	f.Sync()
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetFileList(path string, fileChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			wg.Add(1)
			go GetFileList(fullPath, fileChan, wg)
		} else {
			fileChan <- fullPath
		}
	}
}
