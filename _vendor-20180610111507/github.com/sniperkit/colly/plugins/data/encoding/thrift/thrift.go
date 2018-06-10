package iterator

import (
	"github.com/thrift-iterator/go"
)

type NewOrderRequest struct {
	Lines []NewOrderLine `thrift:",1"`
}

type NewOrderLine struct {
	ProductId string `thrift:",1"`
	Quantity  int    `thrift:",2"`
}

func thriftEncodeBytes(input []*NewOrderLine) (output []byte, err error) {
	// marshal to thrift
	thriftEncodedBytes, err = thrifter.Marshal(
		NewOrderRequest{
			Lines: []NewOrderLine{
				{"apple", 1},
				{"orange", 2},
			},
		},
	)
	return
}

func thriftEncodeBytes(input []byte) (ouput NewOrderRequest, err error) {
	// unmarshal back
	err = thrifter.Unmarshal(input, &ouput)
	return
}
