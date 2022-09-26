// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/vms"
	"github.com/ava-labs/avalanchego/vms/avm/config"
)

var _ vms.Factory = &Factory{}

type Factory struct {
	config.Config
}

func (f *Factory) New(*snow.Context) (interface{}, error) {
	return &VM{Config: f.Config}, nil
}
