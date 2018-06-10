// Go support for Protocol Buffers RPC which compatiable with https://github.com/Baidu-ecom/Jprotobuf-rpc-socket
//
// Copyright 2002-2007 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package pbrpc_test

import (
	"bytes"
	"fmt"
	"testing"

	pbrpc "github.com/baidu-golang/pbrpc"
	"github.com/stretchr/testify/assert"
)

func TestRpcDataWriteReader(t *testing.T) {

	h := pbrpc.Header{}
	h.SetMagicCode([]byte("PRPB"))
	h.SetMessageSize(12300)
	h.SetMetaSize(59487)

	bs, _ := h.Write()

	if len(bs) != pbrpc.SIZE {
		t.Errorf("current head size is '%d', should be '%d'", len(bs), pbrpc.SIZE)
	}

	h2 := pbrpc.Header{}
	h2.Read(bs)
	if !bytes.Equal(h.GetMagicCode(), h2.GetMagicCode()) {
		t.Errorf("magic code is not same. expect '%b' actual is '%b'", h.GetMagicCode(), h2.GetMagicCode())
	}

	assert.Equal(t, h.GetMessageSize(), h2.GetMessageSize(), fmt.Sprintf("expect message size is %d, acutal value is %d", h.GetMessageSize(), h2.GetMessageSize()))

}
