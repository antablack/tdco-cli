# tdco CLI

## Description

**tdco CLI** is a command-line tool written in Go that helps you identify and list `TODO` and `FIXME` comments in the files of a specified directory. This tool is perfect for keeping track of pending tasks and areas in the code that need attention by generating a Markdown file with all the relevant comments.

## Features

- Scans all files in a specified directory.
- Searches for `TODO` and `FIXME` comments.
- Generates a Markdown file listing all found comments, including the file path and line number.
- Easy to use and fast.

## Installation

First, make sure you have [Go](https://golang.org/dl/) installed on your system.

1. Clone this repository:

    ```sh
    git clone https://github.com/your_username/bin/tdco.git
    cd bin/tdco
    ```

2. Build the project:

    ```sh
    go build -ldflags="-s -w" -o bin/tdco
    ```

3. (Optional) Move the executable to a directory included in your PATH to use it from anywhere:

    ```sh
    mv bin/tdco /usr/local/bin/
    ```

## Usage

To use the tool, simply run the following command in your terminal:

```sh
bin/tdco --directory /path/to/directory --md-file report.md
```

Where:
- dir specifies the path to the directory you want to scan.
- output is the name of the Markdown file to be generated with the list of TODO and FIXME comments.

### Example
```sh
bin/tdco --directory ./ --md-file TODO.md
```
This will scan all the files in the root directory and generate a TODO.md file with content similar to:

##### TODO  
 -  Replace deprecated function  <span style="background-color: #F3CA52; padding: 5px; border-radius: 5px; margin: 3px">code</span> [utils/file.go#L49](utils/file.go#L49) 
##### FIXME  
 -  Change function name  <span style="background-color: #7ABA78; padding: 5px; border-radius: 5px; margin: 3px">quality</span> [utils/file.go#L36](utils/file.go#L36)
-  

## Contribution
Contributions are welcome! If you want to improve this tool, please follow these steps:

1. Fork the repository.
2. Create a new branch (git checkout -b feature/new-feature).
3. Make your changes and commit them (git commit -am 'Add new feature').
4. Push your changes (git push origin feature/new-feature).
5. Open a Pull Request.

License
This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.

