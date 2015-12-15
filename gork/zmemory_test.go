package gork

import (
	"encoding/binary"
	"testing"
)

// len must be multiple of 4
var readTestData []byte = []byte{42, 73, 96, 7, 28, 1, 2, 3}
var writeTestData []byte = []byte{3, 2, 1, 28, 7, 96, 73, 42}

var byteOrder binary.ByteOrder = binary.BigEndian

func TestByteAt(t *testing.T) {
	mem := ZMemory(readTestData)

	for i := range readTestData {
		if readTestData[i] != mem.ByteAt(uint32(i)) {
			t.Fail()
		}
	}
}

func TestWordAt(t *testing.T) {
	mem := ZMemory(readTestData)

	for i := uint32(0); i < uint32(len(readTestData)/2); i++ {
		if byteOrder.Uint16(mem[i:i+2]) != mem.WordAt(i) {
			t.Fail()
		}
	}
}

func TestUint32At(t *testing.T) {
	mem := ZMemory(readTestData)

	for i := uint32(0); i < uint32(len(readTestData)/4); i++ {
		if byteOrder.Uint16(mem[i:i+4]) != mem.WordAt(i) {
			t.Fail()
		}
	}
}

func TestWriteByteAt(t *testing.T) {
	mem := ZMemory(readTestData)

	for i := range readTestData {
		mem.WriteByteAt(uint32(i), writeTestData[i])
		if mem.ByteAt(uint32(i)) != writeTestData[i] {
			t.Fail()
		}
	}
}

func TestWriteWordAt(t *testing.T) {
	mem := ZMemory(readTestData)

	for i := uint32(0); i < uint32(len(readTestData)/2); i++ {
		toWrite := byteOrder.Uint16(writeTestData[i : i+2])
		mem.WriteWordAt(i, toWrite)
		if toWrite != mem.WordAt(i) {
			t.Fail()
		}
	}
}

func TestPeekByte(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := range readTestData {
		if seq.PeekByte() != seq.mem.ByteAt(seq.pos) || seq.pos != uint32(i) {
			t.Fail()
		}
		seq.pos += 1
	}
}

func TestPeekWord(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := uint32(0); i < uint32(len(mem)/2); i++ {
		if seq.PeekWord() != seq.mem.WordAt(seq.pos) || seq.pos != uint32(i*2) {
			t.Fail()
		}
		seq.pos += 2
	}
}

func TestPeekUint32(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := uint32(0); i < uint32(len(mem)/4); i++ {
		if seq.PeekUInt32() != seq.mem.UInt32At(seq.pos) || seq.pos != uint32(i*4) {
			t.Fail()
		}
		seq.pos += 4
	}
}

func TestReadByte(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := range readTestData {
		if seq.pos != uint32(i) || seq.ReadByte() != seq.mem.ByteAt(uint32(i)) {
			t.Fail()
		}
	}
}

func TestReadWord(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := uint32(0); i < uint32(len(mem)/2); i++ {
		if seq.pos != uint32(i*2) || seq.ReadWord() != seq.mem.WordAt(i*2) {
			t.Fail()
		}
	}
}

func TestReadUint32(t *testing.T) {
	mem := ZMemory(readTestData)
	seq := mem.GetSequential(0)

	for i := uint32(0); i < uint32(len(mem)/4); i++ {
		if seq.pos != uint32(i*4) || seq.ReadUint32() != seq.mem.UInt32At(i*4) {
			t.Fail()
		}
	}
}

func TestPackedAddres(t *testing.T) {
	for i := uint32(0); i < 10; i++ {
		if !IsPackedAddress(PackedAddress(i)) {
			t.Fail()
		}
	}
}