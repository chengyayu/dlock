package redis

import (
	"github.com/chengyayu/dlock"
	"github.com/go-redsync/redsync/v4"
)

var _ dlock.MutexFactory = (*Factory)(nil)

type Factory struct {
	rs *redsync.Redsync
}

func NewFactory(rs *redsync.Redsync) *Factory {
	return &Factory{rs: rs}
}

func (f *Factory) MakeMutex(key string) (dlock.Mutex, error) {
	mutex := f.rs.NewMutex(key)
	return &Mutex{mutex}, nil
}
