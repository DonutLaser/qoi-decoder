package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

func debugWriteImage(path string, image Image) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Width))
	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Height))

	for index, pixel := range image.Data {
		if index == len(image.Data)-1 {
			sb.WriteString(fmt.Sprintf("%d %d %d %d", pixel.R, pixel.G, pixel.B, pixel.A))
		} else {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", pixel.R, pixel.G, pixel.B, pixel.A))
		}
	}

	WriteFile(fmt.Sprintf("%s.txt", path), sb.String())
}

func debugReadImage(data string) (width uint32, height uint32, result []Pixel) {
	lines := strings.Split(data, "\n")

	parsedWidth, _ := strconv.ParseUint(lines[0], 10, 32)
	width = uint32(parsedWidth)

	parsedHeight, _ := strconv.ParseUint(lines[1], 10, 32)
	height = uint32(parsedHeight)

	for i := 2; i < len(lines); i++ {
		rgbaStrings := strings.Split(lines[i], " ")

		var rgba []uint8
		for j := 0; j < len(rgbaStrings); j++ {
			parsed, _ := strconv.ParseUint(rgbaStrings[j], 10, 8)
			rgba = append(rgba, uint8(parsed))
		}

		result = append(result, Pixel{R: rgba[0], G: rgba[1], B: rgba[2], A: rgba[3]})
	}

	return
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

		filename := filepath.Base(args.File)
		filenameNoExtension := filename[:len(filename)-len(filepath.Ext(filename))]

		debugWriteImage(filenameNoExtension, image)
	} else {
		file := ReadFile(args.File)
		width, height, data := debugReadImage(string(file))

		encoded := Encode(data, width, height)

		filename := filepath.Base(args.File)
		filenameNoExtension := filename[:len(filename)-len(filepath.Ext(filename))]

		WriteFile(fmt.Sprintf("%s.qoi", filenameNoExtension), string(encoded))
	}
}
