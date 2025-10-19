package antedecorators_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tabilabs/tabi-v2/app/antedecorators"
	"github.com/tabilabs/tabi-v2/utils"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestTracedDecorator(t *testing.T) {
	output = ""
	anteDecorators := []sdk.AnteFullDecorator{
		sdk.DefaultWrappedAnteDecorator(FakeAnteDecoratorOne{}),
		sdk.DefaultWrappedAnteDecorator(FakeAnteDecoratorTwo{}),
		sdk.DefaultWrappedAnteDecorator(FakeAnteDecoratorThree{}),
	}
	tracedDecorators := utils.Map(anteDecorators, func(d sdk.AnteFullDecorator) sdk.AnteFullDecorator {
		return sdk.DefaultWrappedAnteDecorator(antedecorators.NewTracedAnteDecorator(d, nil))
	})
	chainedHandler, _ := sdk.ChainAnteDecorators(tracedDecorators...)
	chainedHandler(sdk.NewContext(nil, tmproto.Header{}, false, nil), FakeTx{}, false)
	require.Equal(t, "onetwothree", output)
}
