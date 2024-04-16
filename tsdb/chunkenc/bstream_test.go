package chunkenc

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestBstream(t *testing.T) {
	w := bstream{}
	for _, bit := range []bit{true, false} {
		w.writeBit(bit)
	}
	for v := 0; v < 256; v++ {
		w.wirteByte(byte(v))
	}
	for v := 0; v < 10000; v += 123 {
		w.writeBits(uint64(v), 29)
	}
	r := newBReader(w.bytes())
	for _, bit := range []bit{true, false} {
		v, err := r.readBitFast()
		if err != nil {
			v, err = r.readBit()
		}
		require.NoError(t, err)
		require.Equal(t, bit, v)
	}
	for v := 0; v < 256; v++ {
		actual, err := r.readBitsFast(8)
		if err != nil {
			actual, err = r.readBits(8)
		}
		require.NoError(t, err)
		require.Equal(t, uint64(v), actual)
	}
	for v := 0; v < 10000; v += 123 {
		actual, err := r.readBitsFast(29)
		if err != nil {
			actual, err = r.readBits(29)
		}
		require.NoError(t, err)
		require.Equal(t, uint64(v), actual)
	}
}
