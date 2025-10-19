package state

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/tabilabs/tabi-v2/x/evm/types"
)

func (s *DBImpl) SubBalance(evmAddr common.Address, amt *big.Int, reason tracing.BalanceChangeReason) {
	s.k.PrepareReplayedAddr(s.ctx, evmAddr)
	if amt.Sign() == 0 {
		return
	}
	if amt.Sign() < 0 {
		s.AddBalance(evmAddr, new(big.Int).Neg(amt), reason)
		return
	}

	ctx := s.ctx

	// this avoids emitting cosmos events for ephemeral bookkeeping transfers like send_native
	if s.eventsSuppressed {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
	}

	atabi, wei := SplitAtabiWeiAmount(amt)
	addr := s.getTabiAddress(evmAddr)
	err := s.k.BankKeeper().SubUnlockedCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), atabi)), true)
	if err != nil {
		s.err = err
		return
	}
	err = s.k.BankKeeper().SubWei(ctx, addr, wei)
	if err != nil {
		s.err = err
		return
	}

	if s.logger != nil && s.logger.OnBalanceChange != nil {
		// We could modify AddWei instead so it returns us the old/new balance directly.
		newBalance := s.GetBalance(evmAddr)
		oldBalance := new(big.Int).Add(newBalance, amt)

		s.logger.OnBalanceChange(evmAddr, oldBalance, newBalance, reason)
	}

	s.tempStateCurrent.surplus = s.tempStateCurrent.surplus.Add(sdk.NewIntFromBigInt(amt))
}

func (s *DBImpl) AddBalance(evmAddr common.Address, amt *big.Int, reason tracing.BalanceChangeReason) {
	s.k.PrepareReplayedAddr(s.ctx, evmAddr)
	if amt.Sign() == 0 {
		return
	}
	if amt.Sign() < 0 {
		s.SubBalance(evmAddr, new(big.Int).Neg(amt), reason)
		return
	}

	ctx := s.ctx
	// this avoids emitting cosmos events for ephemeral bookkeeping transfers like send_native
	if s.eventsSuppressed {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
	}

	atabi, wei := SplitAtabiWeiAmount(amt)
	addr := s.getTabiAddress(evmAddr)
	err := s.k.BankKeeper().AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), atabi)), true)
	if err != nil {
		s.err = err
		return
	}
	err = s.k.BankKeeper().AddWei(ctx, addr, wei)
	if err != nil {
		s.err = err
		return
	}

	if s.logger != nil && s.logger.OnBalanceChange != nil {
		// We could modify AddWei instead so it returns us the old/new balance directly.
		newBalance := s.GetBalance(evmAddr)
		oldBalance := new(big.Int).Sub(newBalance, amt)

		s.logger.OnBalanceChange(evmAddr, oldBalance, newBalance, reason)
	}

	s.tempStateCurrent.surplus = s.tempStateCurrent.surplus.Sub(sdk.NewIntFromBigInt(amt))
}

func (s *DBImpl) GetBalance(evmAddr common.Address) *big.Int {
	s.k.PrepareReplayedAddr(s.ctx, evmAddr)
	tabiAddr := s.getTabiAddress(evmAddr)
	return s.k.GetBalance(s.ctx, tabiAddr)
}

// should only be called during simulation
func (s *DBImpl) SetBalance(evmAddr common.Address, amt *big.Int, reason tracing.BalanceChangeReason) {
	if !s.simulation {
		panic("should never call SetBalance in a non-simulation setting")
	}
	tabiAddr := s.getTabiAddress(evmAddr)
	moduleAddr := s.k.AccountKeeper().GetModuleAddress(types.ModuleName)
	s.send(tabiAddr, moduleAddr, s.GetBalance(evmAddr))
	if s.err != nil {
		panic(s.err)
	}
	atabi, _ := SplitAtabiWeiAmount(amt)
	coinsAmt := sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), atabi.Add(sdk.OneInt())))
	if err := s.k.BankKeeper().MintCoins(s.ctx, types.ModuleName, coinsAmt); err != nil {
		panic(err)
	}
	s.send(moduleAddr, tabiAddr, amt)
	if s.err != nil {
		panic(s.err)
	}
}

func (s *DBImpl) getTabiAddress(evmAddr common.Address) sdk.AccAddress {
	if s.coinbaseEvmAddress.Cmp(evmAddr) == 0 {
		return s.coinbaseAddress
	}
	return s.k.GetTabiAddressOrDefault(s.ctx, evmAddr)
}

func (s *DBImpl) send(from sdk.AccAddress, to sdk.AccAddress, amt *big.Int) {
	atabi, wei := SplitAtabiWeiAmount(amt)
	err := s.k.BankKeeper().SendCoinsAndWei(s.ctx, from, to, atabi, wei)
	if err != nil {
		s.err = err
	}
}
