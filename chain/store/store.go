package store

import (
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// NewStore creates a new store for chain
func NewStore(ctx sdk.Context, key sdk.StoreKey) KVStore {
	return &storeWapper{
		store:  ctx.KVStore(key),
		num:    ctx.BlockHeight(),
		logger: ctx.Logger().With("module", "kuStoreLog", "storeKey", key, "num", ctx.BlockHeight()),
	}
}

type storeWapper struct {
	num    int64
	logger log.Logger
	store  KVStore
}

func (s *storeWapper) GetStoreType() StoreType {
	return s.store.GetStoreType()
}

func (s *storeWapper) CacheWrap() CacheWrap {
	return s.store.CacheWrap()
}
func (s *storeWapper) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return s.store.CacheWrapWithTrace(w, tc)
}

func (s *storeWapper) Get(key []byte) []byte {
	if s.logger != nil {
		s.logger.Info("Get", "key", key)
	}
	return s.store.Get(key)
}

func (s *storeWapper) Has(key []byte) bool {
	if s.logger != nil {
		s.logger.Info("Has", "key", key)
	}
	return s.store.Has(key)
}

func (s *storeWapper) Set(key, value []byte) {
	if s.logger != nil {
		s.logger.Info("Set", "key", key, "value", value)
	}
	s.store.Set(key, value)
}

func (s *storeWapper) Delete(key []byte) {
	if s.logger != nil {
		s.logger.Info("Delete", "key", key)
	}
	s.store.Delete(key)
}

func (s *storeWapper) Iterator(start, end []byte) Iterator {
	if s.logger != nil {
		s.logger.Info("Iterator", "start", start, "end", end)
	}
	return s.store.Iterator(start, end)
}

func (s *storeWapper) ReverseIterator(start, end []byte) Iterator {
	if s.logger != nil {
		s.logger.Info("ReverseIterator", "start", start, "end", end)
	}
	return s.store.ReverseIterator(start, end)
}
