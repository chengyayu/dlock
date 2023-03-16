package redis

import (
	"context"
	"github.com/chengyayu/dlock"
	redsync "github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
)

var _ dlock.Mutex = (*Mutex)(nil)

type Mutex struct {
	*redsync.Mutex
}

func (m *Mutex) Destructor() {
}

func (m *Mutex) Lock(ctx context.Context) error {
	return errors.Wrap(m.Mutex.LockContext(ctx), "lock failed")
}

func (m *Mutex) TryLock(ctx context.Context) error {
	return errors.Wrap(m.Mutex.LockContext(ctx), "tryLock failed")
}

func (m *Mutex) Unlock(ctx context.Context) error {
	if ok, err := m.Mutex.UnlockContext(ctx); !ok || err != nil {
		return errors.New("unlock failed")
	}
	return nil
}

func (m *Mutex) Do(ctx context.Context, fn func() error) (err error) {
	defer m.Destructor()
	if err := m.TryLock(ctx); err != nil {
		return err
	}
	defer func() {
		if tempErr := m.Unlock(ctx); tempErr != nil {
			err = tempErr
		}
	}()
	return fn()
}
