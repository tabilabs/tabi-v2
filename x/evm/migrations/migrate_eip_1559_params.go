package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi-v2/x/evm/keeper"
	"github.com/tabilabs/tabi-v2/x/evm/types"
)

func MigrateEip1559Params(ctx sdk.Context, k *keeper.Keeper) error {
	keeperParams := k.GetParamsIfExists(ctx)
	keeperParams.MaxDynamicBaseFeeUpwardAdjustment = types.DefaultParams().MaxDynamicBaseFeeUpwardAdjustment
	keeperParams.MaxDynamicBaseFeeDownwardAdjustment = types.DefaultParams().MaxDynamicBaseFeeDownwardAdjustment
	keeperParams.TargetGasUsedPerBlock = types.DefaultParams().TargetGasUsedPerBlock
	keeperParams.MinimumFeePerGas = types.DefaultParams().MinimumFeePerGas
	k.SetParams(ctx, keeperParams)
	return nil
}
