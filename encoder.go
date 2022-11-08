package main

import (
	"bytes"
)

func Encode(data []Pixel, width uint32, height uint32) []byte {
	var bb bytes.Buffer

	encodeHeader(&bb, width, height)

	lastPixel := Pixel{R: 0, G: 0, B: 0, A: 255}
	runningPixels := make([]Pixel, 64)

	for i := 0; i < len(data); i++ {
		index := hashPixel(&data[i])
		if data[i].IsEqualTo(&runningPixels[index]) {
			if lastPixel.IsEqualTo(&data[i]) {
				var count uint8 = 0
				for i < len(data) && data[i].IsEqualTo(&lastPixel) {
					count++
					i++

					if count == 62 {
						break
					}
				}

				encodeRunChunk(&bb, count)

				i--
			} else {
				encodeIndexChunk(&bb, index)
				lastPixel = data[i]
			}
		} else if data[i].A == lastPixel.A {
			encodeRGBChunk(&bb, &data[i])
			lastPixel = data[i]
		} else if data[i].A != lastPixel.A {
			encodeRGBAChunk(&bb, &data[i])
			lastPixel = data[i]
		}

		runningPixels[index] = data[i]
	}

	encodeStreamEnd(&bb)

	return bb.Bytes()
}

func encodeHeader(bb *bytes.Buffer, width uint32, height uint32) {
	// Magic number
	bb.WriteByte(113)
	bb.WriteByte(111)
	bb.WriteByte(105)
	bb.WriteByte(102)

	// Width
	bb.WriteByte(byte((width & 0xFF000000) >> 24))
	bb.WriteByte(byte((width & 0x00FF0000) >> 16))
	bb.WriteByte(byte((width & 0x0000FF00) >> 8))
	bb.WriteByte(byte(width & 0x000000FF))

	// Height
	bb.WriteByte(byte((height & 0xFF000000) >> 24))
	bb.WriteByte(byte((height & 0x00FF0000) >> 16))
	bb.WriteByte(byte((height & 0x0000FF00) >> 8))
	bb.WriteByte(byte(height & 0x000000FF))

	// Channels
	bb.WriteByte(4)

	// Colorspace
	bb.WriteByte(0)
}

func encodeIndexChunk(bb *bytes.Buffer, index uint8) {
	bb.WriteByte(0b00000000 | index)
}

func encodeRunChunk(bb *bytes.Buffer, count uint8) {
	bb.WriteByte(0b11000000 | (count - 1))
}

func encodeRGBChunk(bb *bytes.Buffer, pixel *Pixel) {
	bb.WriteByte(0b11111110)
	bb.WriteByte(pixel.R)
	bb.WriteByte(pixel.G)
	bb.WriteByte(pixel.B)
}

func encodeRGBAChunk(bb *bytes.Buffer, pixel *Pixel) {
	bb.WriteByte(0b11111111)
	bb.WriteByte(pixel.R)
	bb.WriteByte(pixel.G)
	bb.WriteByte(pixel.B)
	bb.WriteByte(pixel.A)
}

func encodeStreamEnd(bb *bytes.Buffer) {
	for i := 0; i < 7; i++ {
		bb.WriteByte(0)
	}

	bb.WriteByte(1)
}
