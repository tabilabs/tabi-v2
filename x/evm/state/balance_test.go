package state_test

import (
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/tabilabs/tabi-v2/testutil/keeper"
	"github.com/tabilabs/tabi-v2/x/evm/state"
	"github.com/tabilabs/tabi-v2/x/evm/types"
)

func TestAddBalance(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	db := state.NewDBImpl(ctx, k, false)
	tabiAddr, evmAddr := testkeeper.MockAddressPair()
	require.Equal(t, big.NewInt(0), db.GetBalance(evmAddr))
	db.AddBalance(evmAddr, big.NewInt(0), tracing.BalanceChangeUnspecified)

	// set association
	k.SetAddressMapping(db.Ctx(), tabiAddr, evmAddr)
	require.Equal(t, big.NewInt(0), db.GetBalance(evmAddr))
	db.AddBalance(evmAddr, big.NewInt(10000000000000), tracing.BalanceChangeUnspecified)
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr), big.NewInt(10000000000000))

	_, evmAddr2 := testkeeper.MockAddressPair()
	db.SubBalance(evmAddr2, big.NewInt(-5000000000000), tracing.BalanceChangeUnspecified) // should redirect to AddBalance
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr), big.NewInt(10000000000000))
	require.Equal(t, db.GetBalance(evmAddr2), big.NewInt(5000000000000))

	_, evmAddr3 := testkeeper.MockAddressPair()
	db.SelfDestruct(evmAddr3)
	db.AddBalance(evmAddr2, big.NewInt(5000000000000), tracing.BalanceChangeUnspecified)
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr3), big.NewInt(0))
}

func TestSubBalance(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	db := state.NewDBImpl(ctx, k, false)
	tabiAddr, evmAddr := testkeeper.MockAddressPair()
	require.Equal(t, big.NewInt(0), db.GetBalance(evmAddr))
	db.SubBalance(evmAddr, big.NewInt(0), tracing.BalanceChangeUnspecified)

	// set association
	k.SetAddressMapping(db.Ctx(), tabiAddr, evmAddr)
	require.Equal(t, big.NewInt(0), db.GetBalance(evmAddr))
	amt := sdk.NewCoins(sdk.NewCoin(k.GetBaseDenom(ctx), sdk.NewInt(20)))
	k.BankKeeper().MintCoins(db.Ctx(), types.ModuleName, amt)
	k.BankKeeper().SendCoinsFromModuleToAccount(db.Ctx(), types.ModuleName, tabiAddr, amt)
	db.SubBalance(evmAddr, big.NewInt(10000000000000), tracing.BalanceChangeUnspecified)
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr), big.NewInt(10000000000000))

	_, evmAddr2 := testkeeper.MockAddressPair()
	amt = sdk.NewCoins(sdk.NewCoin(k.GetBaseDenom(ctx), sdk.NewInt(10)))
	k.BankKeeper().MintCoins(db.Ctx(), types.ModuleName, amt)
	k.BankKeeper().SendCoinsFromModuleToAccount(db.Ctx(), types.ModuleName, sdk.AccAddress(evmAddr2[:]), amt)
	db.AddBalance(evmAddr2, big.NewInt(-5000000000000), tracing.BalanceChangeUnspecified) // should redirect to SubBalance
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr), big.NewInt(10000000000000))
	require.Equal(t, db.GetBalance(evmAddr2), big.NewInt(5000000000000))

	// insufficient balance
	db.SubBalance(evmAddr2, big.NewInt(10000000000000), tracing.BalanceChangeUnspecified)
	require.NotNil(t, db.Err())

	db.WithErr(nil)
	_, evmAddr3 := testkeeper.MockAddressPair()
	db.SelfDestruct(evmAddr3)
	db.SubBalance(evmAddr2, big.NewInt(5000000000000), tracing.BalanceChangeUnspecified)
	require.Nil(t, db.Err())
	require.Equal(t, db.GetBalance(evmAddr3), big.NewInt(0))
}

func TestSetBalance(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	db := state.NewDBImpl(ctx, k, true)
	_, evmAddr := testkeeper.MockAddressPair()
	db.SetBalance(evmAddr, big.NewInt(10000000000000), tracing.BalanceChangeUnspecified)
	require.Equal(t, big.NewInt(10000000000000), db.GetBalance(evmAddr))

	tabiAddr2, evmAddr2 := testkeeper.MockAddressPair()
	k.SetAddressMapping(db.Ctx(), tabiAddr2, evmAddr2)
	db.SetBalance(evmAddr2, big.NewInt(10000000000000), tracing.BalanceChangeUnspecified)
	require.Equal(t, big.NewInt(10000000000000), db.GetBalance(evmAddr2))
}

func TestSurplus(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	_, evmAddr := testkeeper.MockAddressPair()

	// test negative atabisurplus negative wei surplus
	db := state.NewDBImpl(ctx, k, false)
	db.AddBalance(evmAddr, big.NewInt(1_000_000_000_001), tracing.BalanceChangeUnspecified)
	_, err := db.Finalize()
	require.Nil(t, err)

	// test negative atabisurplus positive wei surplus (negative total)
	db = state.NewDBImpl(ctx, k, false)
	db.AddBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.SubBalance(evmAddr, big.NewInt(1), tracing.BalanceChangeUnspecified)
	_, err = db.Finalize()
	require.Nil(t, err)

	// test negative atabisurplus positive wei surplus (positive total)
	db = state.NewDBImpl(ctx, k, false)
	db.AddBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.SubBalance(evmAddr, big.NewInt(2), tracing.BalanceChangeUnspecified)
	db.SubBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	surplus, err := db.Finalize()
	require.Nil(t, err)
	require.Equal(t, sdk.OneInt(), surplus)

	// test positive atabisurplus negative wei surplus (negative total)
	db = state.NewDBImpl(ctx, k, false)
	db.SubBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(2), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	_, err = db.Finalize()
	require.Nil(t, err)

	// test positive atabisurplus negative wei surplus (positive total)
	db = state.NewDBImpl(ctx, k, false)
	db.SubBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	surplus, err = db.Finalize()
	require.Nil(t, err)
	require.Equal(t, sdk.OneInt(), surplus)

	// test snapshots
	db = state.NewDBImpl(ctx, k, false)
	db.SubBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	db.Snapshot()
	db.SubBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	db.Snapshot()
	db.SubBalance(evmAddr, big.NewInt(1_000_000_000_000), tracing.BalanceChangeUnspecified)
	db.AddBalance(evmAddr, big.NewInt(999_999_999_999), tracing.BalanceChangeUnspecified)
	surplus, err = db.Finalize()
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(3), surplus)
}
