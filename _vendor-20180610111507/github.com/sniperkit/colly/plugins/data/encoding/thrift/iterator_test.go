package iterator_test

import (
	"testing"

	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/raw"
)

func TestThriftIterator(t *testing.T) {
	// partial decoding
	decoder := thrifter.NewDecoder(reader)
	var msgHeader protocol.MessageHeader
	decoder.Decode(&msgHeader)
	var msgArgs raw.Struct
	decoder.Decode(&msgArgs)

	// modify...

	// encode back
	encoder := thrifter.NewEncoder(writer)
	encoder.Encode(msgHeader)
	encoder.Encode(msgArgs)

}
