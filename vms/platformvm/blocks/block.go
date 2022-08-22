// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
)

// Block defines the common stateless interface for all blocks
type Block interface {
	ID() ids.ID
	Parent() ids.ID
	Bytes() []byte
	Height() uint64

	// Txs returns list of transactions contained in the block
	Txs() []*txs.Tx

	// Visit calls [visitor] with this block's concrete type
	Visit(visitor Visitor) error

	initialize(bytes []byte) error
}

func initialize(blk Block) error {
	// We serialize this block as a pointer so that it can be deserialized into
	// a Block
	bytes, err := Codec.Marshal(txs.Version, &blk)
	if err != nil {
		return fmt.Errorf("couldn't marshal block: %w", err)
	}
	return blk.initialize(bytes)
}
