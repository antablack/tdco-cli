package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/antablack/tdco-cli/utils"
	"github.com/urfave/cli"
)

type Task struct {
	Description string
	Type        string
	Tags        []string
	Path        string
}

func processFile(taskChan chan<- Task, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to read file %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile(`(TODO|FIXME)\:(.+)`)
	re1 := regexp.MustCompile(`(\@\w+)`)
	var lineCounter int = 1
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
				Path:        fmt.Sprintf("%s#L%s", filePath, fmt.Sprint(lineCounter)),
			}
		}
		lineCounter = lineCounter + 1
	}
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
				Value: "./TODO.md",
				Usage: "Markdown file.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			directoryPath := cCtx.String("directory")
			if err := utils.ValidDirectory(directoryPath); err != nil {
				fmt.Println("Error:", err)
				return err
			}

			mdFile := cCtx.String("md-file")

			fileChan := make(chan string)
			var wg sync.WaitGroup

			wg.Add(1)
			go utils.GetFileList(directoryPath, fileChan, &wg)

			go func() {
				wg.Wait()
				close(fileChan)
			}()

			var wg2 sync.WaitGroup
			taskChan := make(chan Task)

			go func() {
				for file := range fileChan {
					fileName := filepath.Base(file)
					reservedFiles := []string{"TODO.md"}
					if !utils.Contains(reservedFiles, fileName) {
						wg2.Add(1)
						go processFile(taskChan, file, &wg2)
					}
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
						todo = fmt.Sprintf(`%s <span style="background-color: %s; padding: 5px; border-radius: 5px; margin: 3px">%s</span>`, todo, utils.GetRandomColor(), tag)
					}
					todo = fmt.Sprintf("%s [%s](%s)", todo, task.Path, utils.SanitizeURL(task.Path))
				} else if task.Type == "FIXME" {
					fixme = fmt.Sprintf("%s \n - %s", fixme, task.Description)
					for _, tag := range task.Tags {
						fixme = fmt.Sprintf(`%s <span style="background-color: %s; padding: 5px; border-radius: 5px; margin: 3px">%s</span>`, fixme, utils.GetRandomColor(), tag)
					}
					fixme = fmt.Sprintf("%s [%s](%s)", fixme, task.Path, utils.SanitizeURL(task.Path))
				}
			}
			wg.Wait()
			content := fmt.Sprintf("##### TODO %s \n##### FIXME %s", todo, fixme)
			utils.OverwriteFile(content, mdFile)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
