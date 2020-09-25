package keeper

import (
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// BlacklistedAddr checks if a given address is blacklisted (i.e restricted from
// receiving funds)
func (a AssetKeeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	// TODO: node black list
	return false
}

// GetAllBalances get all coins for a account
func (a AssetKeeper) GetAllBalances(ctx sdk.Context, ID types.AccountID) Coins {
	res, _ := a.GetCoins(ctx, ID)
	return res
}

// GetBalance get coins balance for account
func (a AssetKeeper) GetBalance(ctx sdk.Context, ID types.AccountID, denom string) Coin {
	creator, symbol, err := types.CoinAccountsFromDenom(denom)
	if err != nil {
		panic(err)
	}

	res, _ := a.GetCoin(ctx, ID, creator, symbol)
	return res
}

func (a AssetKeeper) Approve(ctx sdk.Context, id, spender types.AccountID, amt types.Coins, isLock bool) error {
	logger := a.Logger(ctx)

	logger.Debug("approve coins", "id", id, "spender", spender, "amount", amt)

	apporveCoins, err := a.getApprove(ctx, id, spender)
	if apporveCoins != nil {
		if apporveCoins.IsLock != isLock {
			return types.ErrAssetApporveCannotChangeLock
		}

		if apporveCoins.IsLock && apporveCoins.Amount.IsAnyGT(amt) {
			return sdkerrors.Wrap(types.ErrAssetApporveCannotChangeLock, "amount in lock apporve cannot be less")
		}
	} else {
		apporveCoins = NewApproveData(amt)
	}

	apporveCoins.IsLock = isLock
	apporveCoins.Amount = amt

	err = a.setApprove(ctx, id, spender, *apporveCoins)
	if err != nil {
		return sdkerrors.Wrapf(err, "approve %s to %s by %s error", id, spender, amt)
	}

	return nil
}

func (a AssetKeeper) ApplyApporve(ctx sdk.Context, from, to types.AccountID, amount Coins) error {
	apporveCoins, err := a.getApprove(ctx, from, to)
	if err != nil {
		return sdkerrors.Wrap(err, "apply apporveCoins get error")
	}

	if apporveCoins == nil {
		return types.ErrAssetApporveNotEnough
	}

	newApporves, hasNeg := apporveCoins.Amount.SafeSub(amount)
	if hasNeg {
		return types.ErrAssetApporveNotEnough
	}

	apporveCoins.Amount = newApporves

	if err := a.setApprove(ctx, from, to, *apporveCoins); err != nil {
		return sdkerrors.Wrap(err, "apply apporveCoins set new error")
	}

	return nil
}
