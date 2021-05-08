package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

func NewDistributionRecord(status DistributionStatus, distributionType DistributionType, distributionName string, recipientAddress string, coins sdk.Coins, start int64, end int64) DistributionRecord {
	return DistributionRecord{DistributionStatus: status, DistributionType: distributionType, DistributionName: distributionName, RecipientAddress: recipientAddress, Coins: coins, DistributionStartHeight: start, DistributionCompletedHeight: end}
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

func NewUserClaim(userAddress string, userClaimType DistributionType, time string) UserClaim {
	return UserClaim{UserAddress: userAddress, UserClaimType: userClaimType, Locked: false, UserClaimTime: time}
}

func (uc UserClaim) Validate() bool {
	if len(uc.UserAddress) == 0 {
		return false
	}
	return true
}

func (uc UserClaim) IsLocked() bool {
	return uc.Locked
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
func GetDistributionStatus(status string) DistributionStatus {
	switch status {
	case "Completed":
		return DistributionStatus_DISTRIBUTION_STATUS_COMPLETED
	case "Pending":
		return DistributionStatus_DISTRIBUTION_STATUS_PENDING
	default:
		return DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED
	}
}

func IsValidClaim(claimType string) (DistributionType, bool) {
	switch claimType {
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}

func IsValidDistribution(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return DistributionType_DISTRIBUTION_TYPE_AIRDROP, true
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}
