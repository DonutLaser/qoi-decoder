package main

import (
	"fmt"
)

type ChannelsType uint8
type ColorspaceType uint8
type ChunkType uint8

const (
	CHANNELS_RGB  ChannelsType = 3
	CHANNELS_RGBA ChannelsType = 4
)

const (
	COLORSPACE_SRGB_WITH_LINEAR_ALPHA ColorspaceType = 0
	COLORSPACE_ALL_CHANNELS_LINEAR    ColorspaceType = 1
)

const (
	CHUNK_RGB ChunkType = iota
	CHUNK_RGBA
	CHUNK_INDEX
	CHUNK_DIFF
	CHUNK_LUMA
	CHUNK_RUN
)

type Header struct {
	Width      uint32
	Height     uint32
	Channels   ChannelsType
	Colorspace ColorspaceType
}

type Image struct {
	Header Header
	Data   []Pixel
}

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (pixel *Pixel) IsEqualTo(other *Pixel) bool {
	return pixel.R == other.R && pixel.G == other.G && pixel.B == other.B && pixel.A == other.A
}

func Decode(data []byte) (success bool, result Image) {
	buffer := NewBuffer(data)

	var isQoi = checkMagic(&buffer)
	if !isQoi {
		fmt.Println("Image has a .qoi extension, but it is not actually a .qoi image")
		return false, Image{}
	}

	header := parseHeader(&buffer)

	pixels := make([]Pixel, 0)

	lastPixel := Pixel{R: 0, G: 0, B: 0, A: 255}
	runningPixels := make([]Pixel, 64)

	endMarker := buffer.PeekBytes(8)
	for !isEnd(endMarker) {
		chunkType := getNextChunkType(&buffer)

		switch chunkType {
		case CHUNK_RGB:
			bytes := buffer.ConsumeBytes(4)

			lastPixel.R = bytes[1]
			lastPixel.G = bytes[2]
			lastPixel.B = bytes[3]
			pixels = append(pixels, lastPixel)

			runningPixels[hashPixel(&lastPixel)] = lastPixel
		case CHUNK_RGBA:
			bytes := buffer.ConsumeBytes(5)

			lastPixel.R = bytes[1]
			lastPixel.G = bytes[2]
			lastPixel.B = bytes[3]
			lastPixel.A = bytes[4]
			pixels = append(pixels, lastPixel)

			runningPixels[hashPixel(&lastPixel)] = lastPixel
		case CHUNK_INDEX:
			b := buffer.ConsumeBytes(1)[0]

			lastPixel = runningPixels[b]
			pixels = append(pixels, lastPixel)
		case CHUNK_DIFF:
			b := buffer.ConsumeBytes(1)[0]

			dr := (b&0b00110000)>>4 - 2
			dg := (b&0b00001100)>>2 - 2
			db := b&0b00000011 - 2

			lastPixel.R += dr
			lastPixel.G += dg
			lastPixel.B += db
			pixels = append(pixels, lastPixel)

			runningPixels[hashPixel(&lastPixel)] = lastPixel
		case CHUNK_LUMA:
			bytes := buffer.ConsumeBytes(2)

			drdg := (bytes[1]&0b11110000)>>4 - 8
			dbdg := bytes[1]&0b00001111 - 8

			dg := bytes[0]&0b00111111 - 32
			dr := dg + drdg
			db := dg + dbdg

			lastPixel.R += dr
			lastPixel.G += dg
			lastPixel.B += db
			pixels = append(pixels, lastPixel)

			runningPixels[hashPixel(&lastPixel)] = lastPixel
		case CHUNK_RUN:
			b := buffer.ConsumeBytes(1)[0]
			runLength := (b & 0b00111111) + 1

			for i := 0; i < int(runLength); i++ {
				pixels = append(pixels, lastPixel)
			}
		}

		endMarker = buffer.PeekBytes(8)
	}

	return true, Image{Header: header, Data: pixels}
}

func parseHeader(buffer *Buffer) Header {
	// We already checked the magic number, so not checking it here again

	widthBytes := buffer.ConsumeBytes(4)
	heightBytes := buffer.ConsumeBytes(4)
	channelsBytes := buffer.ConsumeBytes(1)
	colorspaceBytes := buffer.ConsumeBytes(1)

	return Header{
		Width:      uint32(widthBytes[0])<<24 | uint32(widthBytes[1])<<16 | uint32(widthBytes[2])<<8 | uint32(widthBytes[3]),
		Height:     uint32(heightBytes[0])<<24 | uint32(heightBytes[1])<<16 | uint32(heightBytes[2])<<8 | uint32(heightBytes[3]),
		Channels:   ChannelsType(channelsBytes[0]),
		Colorspace: ColorspaceType(colorspaceBytes[0]),
	}
}

func getNextChunkType(buffer *Buffer) ChunkType {
	number := buffer.PeekBytes(1)[0]

	if number == 0b11111110 {
		return CHUNK_RGB
	}
	if number == 0b11111111 {
		return CHUNK_RGBA
	}

	n := (number & 0b11000000) >> 6
	if n == 0 {
		return CHUNK_INDEX
	}
	if n == 1 {
		return CHUNK_DIFF
	}
	if n == 2 {
		return CHUNK_LUMA
	}
	if n == 3 {
		return CHUNK_RUN
	}

	panic("Unknown chunk")
}

func checkMagic(buffer *Buffer) bool {
	bytes := buffer.ConsumeBytes(4)
	return bytes[0] == 113 && bytes[1] == 111 && bytes[2] == 105 && bytes[3] == 102
}

func isEnd(bytes []byte) bool {
	for i := 0; i < 7; i++ {
		if bytes[i] != 0 {
			return false
		}
	}

	return bytes[7] == 1
}

func hashPixel(pixel *Pixel) uint8 {
	return (pixel.R*3 + pixel.G*5 + pixel.B*7 + pixel.A*11) % 64
}
