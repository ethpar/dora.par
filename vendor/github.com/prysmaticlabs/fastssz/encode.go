package ssz

import (
	"encoding/binary"
	"fmt"
)

// MarshalSSZ marshals an object
func MarshalSSZ(m Marshaler) ([]byte, error) {
	buf := make([]byte, m.SizeSSZ())
	return m.MarshalSSZTo(buf[:0])
}

var (
	ErrOffset                = fmt.Errorf("incorrect offset")
	ErrSize                  = fmt.Errorf("incorrect size")
	ErrBytesLength           = fmt.Errorf("bytes array does not have the correct length")
	ErrVectorLength          = fmt.Errorf("vector does not have the correct length")
	ErrListTooBig            = fmt.Errorf("list length is higher than max value")
	ErrEmptyBitlist          = fmt.Errorf("bitlist is empty")
	ErrInvalidVariableOffset = fmt.Errorf("invalid ssz encoding. first variable element offset indexes into fixed value data")
)

func ErrBytesLengthFn(name string, found, expected int) error {
	return fmt.Errorf("%s (%v): expected %d and %d found", name, ErrBytesLength, expected, found)
}

func ErrVectorLengthFn(name string, found, expected int) error {
	return fmt.Errorf("%s (%v): expected %d and %d found", name, ErrBytesLength, expected, found)
}

func ErrListTooBigFn(name string, found, max int) error {
	return fmt.Errorf("%s (%v): max expected %d and %d found", name, ErrListTooBig, max, found)
}

// MarshalUint64 marshals a little endian uint64 to dst
func MarshalUint64(dst []byte, i uint64) []byte {
	return binary.LittleEndian.AppendUint64(dst, i)
}

// MarshalUint32 marshals a little endian uint32 to dst
func MarshalUint32(dst []byte, i uint32) []byte {
	return binary.LittleEndian.AppendUint32(dst, i)
}

// MarshalUint16 marshals a little endian uint16 to dst
func MarshalUint16(dst []byte, i uint16) []byte {
	return binary.LittleEndian.AppendUint16(dst, i)
}

// MarshalUint8 marshals a little endian uint8 to dst
func MarshalUint8(dst []byte, i uint8) []byte {
	dst = append(dst, byte(i))
	return dst
}

// MarshalBool marshals a boolean to dst
func MarshalBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, 1)
	}
	return append(dst, 0)
}

// WriteOffset writes an offset to dst
func WriteOffset(dst []byte, i int) []byte {
	return MarshalUint32(dst, uint32(i))
}
