package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Args struct {
	File string
}

func ReadFile(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return []byte{}
	}

	return bytes
}

func WriteFile(path string, content string) bool {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	defer file.Close()

	file.WriteString(content)
	file.Sync()

	return true
}

func printUsage() {
	fmt.Println("Usage: qoi <path/to/file>")
}

func parseArgs() (result Args, success bool) {
	args := os.Args[1:]

	if len(args) < 1 {
		return Args{}, false
	}

	result = Args{
		File: args[0],
	}

	return result, true
}

func debugImage(image Image) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Width))
	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Height))

	for _, pixel := range image.Data {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", pixel.R, pixel.G, pixel.B, pixel.A))
	}

	WriteFile("test.txt", sb.String())
}

func main() {
	args, success := parseArgs()

	if !success {
		printUsage()
		return
	}

	if path.Ext(args.File) == ".qoi" {
		file := ReadFile(args.File)
		success, image := Decode(file)

		if !success {
			return
		}

		debugImage(image)
	} else {
		fmt.Print("Encoding...")
	}
}
