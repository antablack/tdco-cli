package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli"
)

// ValidFile checks if a path is a valid file
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

// ValidDirectory checks if a path is a valid directory
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

// getFileList sends file paths from the directory to the channel
func getFileList(path string, fileChan chan<- string, wg *sync.WaitGroup) {
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
			go getFileList(fullPath, fileChan, wg)
		} else {
			fileChan <- fullPath
		}
	}
}

type Task struct {
	Description string
	Type        string
	Tags        []string
	Path        string
}

// processFile processes a file path received from the channel
func processFile(taskChan chan<- Task, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("uanble to read file %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile(`(TODO|FIXME)\:(.+)`)
	re1 := regexp.MustCompile(`(\@\w+)`)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			matches := re.FindStringSubmatch(line)
			description := matches[2]
			matches1 := re1.FindAllStringSubmatch(description, -1)
			tags := []string{}
			for _, match := range matches1 {
				tags = append(tags, strings.ReplaceAll(match[1], "@", ""))
			}
			taskChan <- Task{
				Description: re1.ReplaceAllString(description, ""),
				Type:        matches[1],
				Tags:        tags,
				Path:        filePath,
			}
		}

	}
}

func getRandomColor() string {
	colors := []string{
		"#F3CA52",
		"#7ABA78",
		"#F6E9B2",
		"#A0153E",
		"#00224D",
		"#57A6A1",
	}
	rand.Seed(time.Now().UnixNano()) // Semilla para el generador de números aleatorios
	return colors[rand.Intn(len(colors))]
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "directory",
				Value: "./",
				Usage: "Directory to assets.",
			},
			&cli.StringFlag{
				Name:  "md-file",
				Value: "./README.md",
				Usage: "Markdown file.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			directoryPath := cCtx.String("directory")
			if err := ValidDirectory(directoryPath); err != nil {
				fmt.Println("Error:", err)
				return err
			}

			mdFile := cCtx.String("md-file")
			if err := ValidFile(mdFile); err != nil {
				fmt.Println("Error:", err)
				return err
			}

			fileChan := make(chan string)
			var wg sync.WaitGroup

			// Start the file discovery goroutine
			wg.Add(1)
			go getFileList(directoryPath, fileChan, &wg)

			// Start the file processing goroutine
			// wg.Add(1)
			// go processFile(fileChan, &wg)

			// Wait for all the file discovery goroutines to finish and then close the channel
			go func() {
				wg.Wait()
				close(fileChan)
			}()

			var wg2 sync.WaitGroup
			taskChan := make(chan Task)

			go func() {
				for file := range fileChan {
					wg2.Add(1)
					go processFile(taskChan, file, &wg2)
					// Simulate file processing
				}
				wg2.Wait()
				close(taskChan)
			}()

			todo := ""
			fixme := ""

			for task := range taskChan {
				if task.Type == "TODO" {
					todo = fmt.Sprintf("%s \n - %s", todo, task.Description)
					for _, tag := range task.Tags {
						todo = fmt.Sprintf(`%s <span style="background-color: %s; padding: 5px; border-radius: 5px">%s</span>`, todo, getRandomColor(), tag)
					}
					todo = fmt.Sprintf("%s %s", todo, task.Path)
				} else if task.Type == "FIXME" {
					fixme = fmt.Sprintf("%s \n - %s", fixme, task.Description)
					for _, tag := range task.Tags {
						fixme = fmt.Sprintf(`%s <span style="background-color: %s; padding: 5px; border-radius: 5px">%s</span>`, fixme, getRandomColor(), tag)
					}
					fixme = fmt.Sprintf("%s %s", fixme, task.Path)
				}
			}
			fmt.Println(fmt.Sprintf("##### TODO %s \n##### FIXME %s", todo, fixme))
			wg.Wait()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
