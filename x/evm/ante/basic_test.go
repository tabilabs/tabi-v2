package ante_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/tabilabs/tabi-v2/testutil/keeper"
	"github.com/tabilabs/tabi-v2/x/evm/ante"
	"github.com/tabilabs/tabi-v2/x/evm/types"
	"github.com/tabilabs/tabi-v2/x/evm/types/ethtx"
)

func TestBasicDecorator(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	a := ante.NewBasicDecorator(k)
	msg, _ := types.NewMsgEVMTransaction(&ethtx.LegacyTx{})
	ctx, err := a.AnteHandle(ctx, &mockTx{msgs: []sdk.Msg{msg}}, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NotNil(t, err) // expect out of gas err
	dataTooLarge := make([]byte, params.MaxInitCodeSize+1)
	for i := 0; i <= params.MaxInitCodeSize; i++ {
		dataTooLarge[i] = 1
	}
	msg, _ = types.NewMsgEVMTransaction(&ethtx.LegacyTx{Data: dataTooLarge})
	ctx, err = a.AnteHandle(ctx, &mockTx{msgs: []sdk.Msg{msg}}, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "code size")
	negAmount := sdk.NewInt(-1)
	msg, _ = types.NewMsgEVMTransaction(&ethtx.LegacyTx{Amount: &negAmount})
	ctx, err = a.AnteHandle(ctx, &mockTx{msgs: []sdk.Msg{msg}}, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.Equal(t, sdkerrors.ErrInvalidCoins, err)
	data := make([]byte, 10)
	for i := 0; i < 10; i++ {
		dataTooLarge[i] = 1
	}
	msg, _ = types.NewMsgEVMTransaction(&ethtx.LegacyTx{Data: data})
	ctx, err = a.AnteHandle(ctx, &mockTx{msgs: []sdk.Msg{msg}}, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.Equal(t, sdkerrors.ErrOutOfGas, err)

	msg, _ = types.NewMsgEVMTransaction(&ethtx.BlobTx{GasLimit: 21000})
	ctx, err = a.AnteHandle(ctx, &mockTx{msgs: []sdk.Msg{msg}}, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NotNil(t, err)
	require.Error(t, err, sdkerrors.ErrUnsupportedTxType)
}
