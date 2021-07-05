package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	bufferLength    = binary.MaxVarintLen64
	maxEncodeLength = binary.MaxVarintLen32
)

// ZskPP is the public parameter of marlin zk-SNARK
// TODO: Adding Verify and proof generation
type ZskPP struct {
	OutputOne uint
	OutputTwo uint
}

// ByteEncode returns the bytes representation of ZskPP
func (pp ZskPP) ByteEncode() []byte {
	buf := make([]byte, bufferLength)
	binary.PutUvarint(buf[:maxEncodeLength], uint64(pp.OutputOne))
	binary.PutUvarint(buf[maxEncodeLength:], uint64(pp.OutputTwo))
	return buf
}

// ByteDecode converts byte encoded ZskPP to a new ZskPP struct
func ByteDecode(encoded []byte) (ZskPP, error) {
	decodeError := errors.New("invalid encoded ZskPP")
	if len(encoded) < bufferLength {
		return ZskPP{}, decodeError
	}
	outputOneReader := bytes.NewReader(encoded[:maxEncodeLength])
	outputTwoReader := bytes.NewReader(encoded[maxEncodeLength:])
	outputOne, err := binary.ReadUvarint(outputOneReader)
	if err != nil {
		return ZskPP{}, decodeError
	}
	outputTwo, err := binary.ReadUvarint(outputTwoReader)
	if err != nil {
		return ZskPP{}, decodeError
	}
	return ZskPP{OutputOne: uint(outputOne), OutputTwo: uint(outputTwo)}, nil
}
