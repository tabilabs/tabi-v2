package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi-v2/x/evm/keeper"
)

func MigrateDisableRegisterPointer(ctx sdk.Context, k *keeper.Keeper) error {
	params := k.GetParams(ctx)
	params.RegisterPointerDisabled = true
	k.SetParams(ctx, params)
	return nil
}
