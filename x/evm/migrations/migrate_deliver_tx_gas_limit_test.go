package migrations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tabilabs/tabi-v2/testutil/keeper"
	"github.com/tabilabs/tabi-v2/x/evm/migrations"
	"github.com/tabilabs/tabi-v2/x/evm/types"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMigrateDeliverTxHookWasmGasLimitParam(t *testing.T) {
	k := testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.NewContext(false, tmtypes.Header{})

	currParams := k.GetParams(ctx)

	// Keep a copy of the other parameters to compare later
	priorityNormalizer := currParams.PriorityNormalizer
	baseFeePerGas := currParams.BaseFeePerGas
	minimumFeePerGas := currParams.MinimumFeePerGas

	// Perform the migration
	err := migrations.MigrateDeliverTxHookWasmGasLimitParam(ctx, &k)
	require.NoError(t, err)

	keeperParams := k.GetParams(ctx)

	// Ensure that the DeliverTxHookWasmGasLimit was migrated to the default value
	require.Equal(t, keeperParams.GetDeliverTxHookWasmGasLimit(), types.DefaultParams().DeliverTxHookWasmGasLimit)

	// Verify that the other parameters were not changed by the migration
	require.True(t, keeperParams.PriorityNormalizer.Equal(priorityNormalizer))
	require.True(t, keeperParams.BaseFeePerGas.Equal(baseFeePerGas))
	require.True(t, keeperParams.MinimumFeePerGas.Equal(minimumFeePerGas))
}
