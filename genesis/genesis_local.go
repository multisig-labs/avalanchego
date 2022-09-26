// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/avalanchego/vms/platformvm/reward"

	xchainconfig "github.com/ava-labs/avalanchego/vms/avm/config"
	pchainconfig "github.com/ava-labs/avalanchego/vms/platformvm/config"
)

// PrivateKey-vmRQiZeXEXYMyJhEiqdC2z5JhuDbxL8ix9UVvjgMu2Er1NepE => P-local1g65uqn6t77p656w64023nh8nd9updzmxyymev2
// PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN => X-local18jma8ppw3nhx5r4ap8clazz0dps7rv5u00z96u
// 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 => 0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC

const (
	VMRQKeyStr          = "vmRQiZeXEXYMyJhEiqdC2z5JhuDbxL8ix9UVvjgMu2Er1NepE"
	VMRQKeyFormattedStr = crypto.PrivateKeyPrefix + VMRQKeyStr

	EWOQKeyStr          = "ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN"
	EWOQKeyFormattedStr = crypto.PrivateKeyPrefix + EWOQKeyStr
)

var (
	VMRQKey *crypto.PrivateKeySECP256K1R
	EWOQKey *crypto.PrivateKeySECP256K1R

	//go:embed genesis_local.json
	localGenesisConfigJSON []byte

	// LocalParams are the params used for local networks
	LocalParams = Params{
		PChainTxFees: pchainconfig.TxFeeUpgrades{
			InitialFees: pchainconfig.TxFees{
				AddPrimaryNetworkValidator: 0,
				AddPrimaryNetworkDelegator: 0,
				AddPOASubnetValidator:      units.MilliAvax,
				AddPOSSubnetValidator:      units.MilliAvax, // didn't exist
				AddPOSSubnetDelegator:      units.MilliAvax, // didn't exist
				RemovePOASubnetValidator:   units.MilliAvax, // didn't exist
				CreateSubnet:               100 * units.MilliAvax,
				CreateChain:                100 * units.MilliAvax,
				TransformSubnet:            100 * units.MilliAvax, // didn't exist
				Import:                     units.MilliAvax,
				Export:                     units.MilliAvax,
			},
			ApricotPhase3Fees: pchainconfig.TxFees{
				AddPrimaryNetworkValidator: 0,
				AddPrimaryNetworkDelegator: 0,
				AddPOASubnetValidator:      units.MilliAvax,
				AddPOSSubnetValidator:      units.MilliAvax, // didn't exist
				AddPOSSubnetDelegator:      units.MilliAvax, // didn't exist
				RemovePOASubnetValidator:   units.MilliAvax, // didn't exist
				CreateSubnet:               100 * units.MilliAvax,
				CreateChain:                100 * units.MilliAvax,
				TransformSubnet:            100 * units.MilliAvax, // didn't exist
				Import:                     units.MilliAvax,
				Export:                     units.MilliAvax,
			},
			BlueberryFees: pchainconfig.TxFees{
				AddPrimaryNetworkValidator: 0,
				AddPrimaryNetworkDelegator: 0,
				AddPOASubnetValidator:      units.MilliAvax,
				AddPOSSubnetValidator:      units.MilliAvax,
				AddPOSSubnetDelegator:      units.MilliAvax,
				RemovePOASubnetValidator:   units.MilliAvax,
				CreateSubnet:               100 * units.MilliAvax,
				CreateChain:                100 * units.MilliAvax,
				TransformSubnet:            100 * units.MilliAvax,
				Import:                     units.MilliAvax,
				Export:                     units.MilliAvax,
			},
		},
		XChainTxFees: xchainconfig.TxFees{
			Base:        units.MilliAvax,
			CreateAsset: units.MilliAvax,
			Operation:   units.MilliAvax,
			Import:      units.MilliAvax,
			Export:      units.MilliAvax,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloAvax,
			MaxValidatorStake: 3 * units.MegaAvax,
			MinDelegatorStake: 25 * units.Avax,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          720 * units.MegaAvax,
			},
		},
	}
)

func init() {
	errs := wrappers.Errs{}
	vmrqBytes, err := cb58.Decode(VMRQKeyStr)
	errs.Add(err)
	ewoqBytes, err := cb58.Decode(EWOQKeyStr)
	errs.Add(err)

	factory := crypto.FactorySECP256K1R{}
	vmrqIntf, err := factory.ToPrivateKey(vmrqBytes)
	errs.Add(err)
	ewoqIntf, err := factory.ToPrivateKey(ewoqBytes)
	errs.Add(err)

	if errs.Err != nil {
		panic(errs.Err)
	}

	VMRQKey = vmrqIntf.(*crypto.PrivateKeySECP256K1R)
	EWOQKey = ewoqIntf.(*crypto.PrivateKeySECP256K1R)
}
