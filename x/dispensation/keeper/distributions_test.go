package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name)
		err := keeper.SetDistribution(ctx, distribution)
		assert.NoError(t, err)
		_, err = keeper.GetDistribution(ctx, name, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
		assert.NoError(t, err)
	}
	list := keeper.GetDistributions(ctx)
	assert.Len(t, list.Distributions, 10)
}
