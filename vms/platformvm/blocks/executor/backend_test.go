// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/vms/platformvm/blocks"
	"github.com/ava-labs/avalanchego/vms/platformvm/state"
)

func TestGetState(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		mockState     = state.NewMockState(ctrl)
		onAcceptState = state.NewMockDiff(ctrl)
		blkID1        = ids.GenerateTestID()
		blkID2        = ids.GenerateTestID()
		b             = &backend{
			state: mockState,
			blkIDToState: map[ids.ID]*blockState{
				blkID1: {
					onAcceptState: onAcceptState,
				},
				blkID2: {},
			},
		}
	)

	{
		// Case: block is in the map and onAcceptState isn't nil.
		gotState, ok := b.GetState(blkID1)
		require.True(ok)
		require.Equal(onAcceptState, gotState)
	}

	{
		// Case: block is in the map and onAcceptState is nil.
		_, ok := b.GetState(blkID2)
		require.False(ok)
	}

	{
		// Case: block is not in the map and block isn't last accepted.
		mockState.EXPECT().GetLastAccepted().Return(ids.GenerateTestID())
		_, ok := b.GetState(ids.GenerateTestID())
		require.False(ok)
	}

	{
		// Case: block is not in the map and block is last accepted.
		blkID := ids.GenerateTestID()
		mockState.EXPECT().GetLastAccepted().Return(blkID)
		gotState, ok := b.GetState(blkID)
		require.True(ok)
		require.Equal(mockState, gotState)
	}
}

func TestBackendGetBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		blkID1       = ids.GenerateTestID()
		statelessBlk = blocks.NewMockBlock(ctrl)
		state        = state.NewMockState(ctrl)
		b            = &backend{
			state: state,
			blkIDToState: map[ids.ID]*blockState{
				blkID1: {
					statelessBlock: statelessBlk,
				},
			},
		}
	)

	{
		// Case: block is in the map.
		gotBlk, err := b.GetBlock(blkID1)
		require.Nil(err)
		require.Equal(statelessBlk, gotBlk)
	}

	{
		// Case: block isn't in the map or database.
		blkID := ids.GenerateTestID()
		state.EXPECT().GetStatelessBlock(blkID).Return(nil, choices.Unknown, database.ErrNotFound)
		_, err := b.GetBlock(blkID)
		require.Equal(database.ErrNotFound, err)
	}

	{
		// Case: block isn't in the map and is in database.
		blkID := ids.GenerateTestID()
		state.EXPECT().GetStatelessBlock(blkID).Return(statelessBlk, choices.Accepted, nil)
		gotBlk, err := b.GetBlock(blkID)
		require.NoError(err)
		require.Equal(statelessBlk, gotBlk)
	}
}

func TestGetTimestamp(t *testing.T) {
	type test struct {
		name              string
		blkF              func(*gomock.Controller) blocks.Block
		backendF          func(*gomock.Controller) *backend
		expectedTimestamp time.Time
	}

	var (
		blueberryAbortBlockTime    = time.Unix(0, 0)
		blueberryCommitBlockTime   = time.Unix(1, 0)
		blueberryProposalBlockTime = time.Unix(2, 0)
		blueberryStandardBlockTime = time.Unix(3, 0)
		blkID                      = ids.GenerateTestID()
	)
	tests := []test{
		{
			name: "blueberry abort block",
			blkF: func(*gomock.Controller) blocks.Block {
				return &blocks.BlueberryAbortBlock{
					Time: uint64(blueberryAbortBlockTime.Unix()),
				}
			},
			backendF: func(*gomock.Controller) *backend {
				return &backend{}
			},
			expectedTimestamp: blueberryAbortBlockTime,
		},
		{
			name: "blueberry commit block",
			blkF: func(*gomock.Controller) blocks.Block {
				return &blocks.BlueberryCommitBlock{
					Time: uint64(blueberryCommitBlockTime.Unix()),
				}
			},
			backendF: func(*gomock.Controller) *backend {
				return &backend{}
			},
			expectedTimestamp: blueberryCommitBlockTime,
		},
		{
			name: "blueberry proposal block",
			blkF: func(*gomock.Controller) blocks.Block {
				return &blocks.BlueberryProposalBlock{
					Time: uint64(blueberryProposalBlockTime.Unix()),
				}
			},
			backendF: func(*gomock.Controller) *backend {
				return &backend{}
			},
			expectedTimestamp: blueberryProposalBlockTime,
		},
		{
			name: "blueberry standard block",
			blkF: func(*gomock.Controller) blocks.Block {
				return &blocks.BlueberryStandardBlock{
					Time: uint64(blueberryStandardBlockTime.Unix()),
				}
			},
			backendF: func(*gomock.Controller) *backend {
				return &backend{}
			},
			expectedTimestamp: blueberryStandardBlockTime,
		},
		{
			name: "apricot block is in map",
			blkF: func(ctrl *gomock.Controller) blocks.Block {
				blk := blocks.NewMockBlock(ctrl)
				blk.EXPECT().ID().Return(blkID)
				return blk
			},
			backendF: func(ctrl *gomock.Controller) *backend {
				return &backend{
					blkIDToState: map[ids.ID]*blockState{
						blkID: {
							timestamp: time.Unix(1337, 0),
						},
					},
				}
			},
			expectedTimestamp: time.Unix(1337, 0),
		},
		{
			name: "apricot block isn't map",
			blkF: func(ctrl *gomock.Controller) blocks.Block {
				blk := blocks.NewMockBlock(ctrl)
				blk.EXPECT().ID().Return(blkID)
				return blk
			},
			backendF: func(ctrl *gomock.Controller) *backend {
				state := state.NewMockState(ctrl)
				state.EXPECT().GetTimestamp().Return(time.Unix(1337, 0))
				return &backend{
					state: state,
				}
			},
			expectedTimestamp: time.Unix(1337, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			backend := tt.backendF(ctrl)
			blk := tt.blkF(ctrl)
			gotTimestamp := backend.getTimestamp(blk)
			require.Equal(tt.expectedTimestamp, gotTimestamp)
		})
	}
}