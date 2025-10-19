package state_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tabilabs/tabi-v2/testutil/keeper"
	"github.com/tabilabs/tabi-v2/x/evm/state"
)

func TestGasRefund(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	statedb := state.NewDBImpl(ctx, k, false)

	require.Equal(t, uint64(0), statedb.GetRefund())
	statedb.AddRefund(2)
	require.Equal(t, uint64(2), statedb.GetRefund())
	statedb.SubRefund(1)
	require.Equal(t, uint64(1), statedb.GetRefund())
	require.Panics(t, func() { statedb.SubRefund(2) })
}
