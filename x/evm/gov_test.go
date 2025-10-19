package evm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tabilabs/tabi-v2/testutil/keeper"
	"github.com/tabilabs/tabi-v2/x/evm"
	"github.com/tabilabs/tabi-v2/x/evm/artifacts/native"
	"github.com/tabilabs/tabi-v2/x/evm/types"
)

func TestAddERCNativePointerProposalsV2(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx(nil)
	require.Nil(t, evm.HandleAddERCNativePointerProposalV2(ctx, k, &types.AddERCNativePointerProposalV2{
		Token:    "test",
		Name:     "NAME",
		Symbol:   "SYMBOL",
		Decimals: 6,
	}))
	pointer, _, exists := k.GetERC20NativePointer(ctx, "test")
	require.True(t, exists)
	qName, _ := native.GetParsedABI().Pack("name")
	resName, err := k.StaticCallEVM(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), &pointer, qName)
	require.Nil(t, err)
	oName, _ := native.GetParsedABI().Unpack("name", resName)
	require.Equal(t, "NAME", oName[0].(string))
	qSymbol, _ := native.GetParsedABI().Pack("symbol")
	resSymbol, err := k.StaticCallEVM(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), &pointer, qSymbol)
	require.Nil(t, err)
	oSymbol, _ := native.GetParsedABI().Unpack("symbol", resSymbol)
	require.Equal(t, "SYMBOL", oSymbol[0].(string))
	qDecimals, _ := native.GetParsedABI().Pack("decimals")
	resDecimals, err := k.StaticCallEVM(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), &pointer, qDecimals)
	require.Nil(t, err)
	oDecimals, _ := native.GetParsedABI().Unpack("decimals", resDecimals)
	require.Equal(t, uint8(6), oDecimals[0].(uint8))

	// make sure pointers deployed this way won't collide in address
	require.Nil(t, evm.HandleAddERCNativePointerProposalV2(ctx, k, &types.AddERCNativePointerProposalV2{
		Token:    "test2",
		Name:     "NAME2",
		Symbol:   "SYMBOL2",
		Decimals: 6,
	}))
	pointer2, _, exists2 := k.GetERC20NativePointer(ctx, "test2")
	require.True(t, exists2)
	require.NotEqual(t, pointer, pointer2)
}
