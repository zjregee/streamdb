package chunkenc

import (
	"io"
	"encoding/binary"
)

type bit bool

const (
	zero bit = false
	one  bit = true
)

type bstream struct {
	stream []byte
	count  int
}

func (bs *bstream) bytes() []byte {
	return bs.stream
}

func (bs *bstream) writeBit(b bit) {
	if bs.count == 0 {
		bs.stream = append(bs.stream, 0)
		bs.count = 8
	}
	i := len(bs.stream) - 1
	if b {
		bs.stream[i] |= 1 << (bs.count - 1)
	}
	bs.count -= 1
}

func (bs *bstream) wirteByte(b byte) {
	if bs.count == 0 {
		bs.stream = append(bs.stream, 0)
		bs.count = 8
	}
	i := len(bs.stream) - 1
	bs.stream[i] |= b >> (8 - bs.count)
	bs.stream = append(bs.stream, 0)
	i += 1
	bs.stream[i] = b << bs.count
}

func (bs *bstream) writeBits(u uint64, nbits int) {
	u <<= 64 - nbits
	for nbits >= 8 {
		bs.wirteByte(byte(u >> 56))
		u <<= 8
		nbits -= 8
	}
	for nbits > 0 {
		bs.writeBit((u >> 63) == 1)
		u <<= 1
		nbits -= 1
	}
}

func newBReader(b []byte) bstreamReader {
	return bstreamReader{ stream: b }
}

type bstreamReader struct {
	stream       []byte
	streamOffset int
	buffer       uint64
	valid        int
}

func (br *bstreamReader) readBit() (bit, error) {
	if br.valid == 0 {
		if !br.loadNextBuffer(1) {
			return false, io.EOF
		}
	}
	return br.readBitFast()
}

func (br *bstreamReader) readBitFast() (bit, error) {
	if br.valid == 0 {
		return false, io.EOF
	}
	br.valid -= 1
	bitmask := uint64(1) << br.valid
	return (br.buffer & bitmask) != 0, nil
}

func (br *bstreamReader) readBits(nbits int) (uint64, error) {
	if br.valid == 0 {
		if !br.loadNextBuffer(nbits) {
			return 0, io.EOF
		}
	}
	if nbits <= br.valid {
		return br.readBitsFast(nbits)
	}
	bitmask := (uint64(1) << br.valid) - 1
	nbits -= br.valid
	v := (br.buffer & bitmask) << nbits
	br.valid = 0
	if !br.loadNextBuffer(nbits) {
		return 0, io.EOF
	}
	bitmask = (uint64(1) << nbits) - 1
	v |= (br.buffer >> (br.valid - nbits)) & bitmask
	br.valid -= nbits
	return v, nil
}

func (br *bstreamReader) readBitsFast(nbits int) (uint64, error) {
	if nbits > br.valid {
		return 0, io.EOF
	}
	bitmask := (uint64(1) << nbits) - 1
	br.valid -= nbits
	return (br.buffer >> br.valid) & bitmask, nil
}

func (br *bstreamReader) loadNextBuffer(nbits int) bool {
	if br.streamOffset >= len(br.stream) {
		return false
	}
	if br.streamOffset + 8 < len(br.stream) {
		br.buffer = binary.BigEndian.Uint64(br.stream[br.streamOffset:])
		br.streamOffset += 8
		br.valid = 64
		return true
	}
	nbytes := (nbits - 1) / 8 + 1
	if br.streamOffset + nbytes > len(br.stream) {
		nbytes = len(br.stream) - br.streamOffset
	}
	buffer := uint64(0)
	for i := 0; i < nbytes; i++ {
		buffer |= (uint64(br.stream[br.streamOffset + i]) << (8 * (nbytes - i - 1)))
	}
	br.buffer = buffer
	br.streamOffset += nbytes
	br.valid = nbytes * 8
	return true;
}
