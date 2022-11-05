package main

import (
	"fmt"
	"strings"
	"testing"
)

func debug(image Image) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Width))
	sb.WriteString(fmt.Sprintf("%d\n", image.Header.Height))

	for _, pixel := range image.Data {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", pixel.R, pixel.G, pixel.B, pixel.A))
	}

	return sb.String()
}

func TestDecode(t *testing.T) {
	testFiles := []string{
		"testdata/dice.qoi",
		"testdata/kodim10.qoi",
		"testdata/kodim23.qoi",
		"testdata/qoi_logo.qoi",
		"testdata/testcard.qoi",
		"testdata/testcard_rgba.qoi",
		"testdata/wikipedia_008.qoi",
	}

	expectedFiles := []string{
		"testdata/dice_expected.txt",
		"testdata/kodim10_expected.txt",
		"testdata/kodim23_expected.txt",
		"testdata/qoi_logo_expected.txt",
		"testdata/testcard_expected.txt",
		"testdata/testcard_rgba_expected.txt",
		"testdata/wikipedia_008_expected.txt",
	}

	for index, path := range testFiles {
		bytes := ReadFile(path)
		expected := string(ReadFile(expectedFiles[index]))

		_, decoded := Decode(bytes)
		result := debug(decoded)

		if result != expected {
			t.Fatalf("Decode %s did not work correctly", path)
		}
	}
}
