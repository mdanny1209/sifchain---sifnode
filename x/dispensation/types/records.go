package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

func NewDistributionRecord(distributionName string, recipientAddress string, coins sdk.Coins, start sdk.Int, end sdk.Int) DistributionRecord {
	return DistributionRecord{DistributionName: distributionName, RecipientAddress: recipientAddress, Coins: coins, DistributionStartHeight: start, DistributionCompletedHeight: end}
}

func (dr DistributionRecord) Validate() bool {
	if len(dr.RecipientAddress) == 0 {
		return false
	}
	if !dr.Coins.IsValid() {
		return false
	}
	if !dr.Coins.IsAllPositive() {
		return false
	}
	return true
}

func (dr DistributionRecord) Add(dr2 DistributionRecord) DistributionRecord {
	dr.Coins = dr.Coins.Add(dr2.Coins...)
	return dr
}

func NewDistribution(t DistributionType, name string) Distribution {
	return Distribution{DistributionType: t, DistributionName: name}
}

func (ar Distribution) Validate() bool {
	if ar.DistributionName == "" {
		return false
	}
	return true
}
