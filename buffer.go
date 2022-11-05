package main

type Buffer struct {
	Data   []byte
	Cursor uint32
}

func NewBuffer(data []byte) Buffer {
	return Buffer{
		Data:   data,
		Cursor: 0,
	}
}

func (buffer *Buffer) ConsumeBytes(count uint32) (result []byte) {
	return buffer.getBytes(count, true)
}

func (buffer *Buffer) ConsumeBytesAsUint32(count uint32) uint32 {
	bytes := buffer.getBytes(count, true)
	return buffer.combineBytesIntoUint32(bytes)
}

func (buffer *Buffer) PeekBytes(count uint32) []byte {
	return buffer.getBytes(count, false)
}

func (buffer *Buffer) PeekBytesAsUint32(count uint32) uint32 {
	bytes := buffer.getBytes(count, false)
	return buffer.combineBytesIntoUint32(bytes)
}

func (buffer *Buffer) getBytes(count uint32, consume bool) (result []byte) {
	result = buffer.Data[buffer.Cursor : buffer.Cursor+count]
	if consume {
		buffer.Cursor += count
	}

	return
}

func (buffer *Buffer) combineBytesIntoUint32(bytes []byte) (result uint32) {
	for i := 0; i < len(bytes); i++ {
		result |= uint32(bytes[i] << ((len(bytes)-1)*8 - i*8))
	}

	return
}
