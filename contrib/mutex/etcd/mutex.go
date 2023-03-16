package etcd

import (
	"context"
	"github.com/chengyayu/dlock"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var _ dlock.Mutex = (*Mutex)(nil)

type Mutex struct {
	client  *clientv3.Client
	session *concurrency.Session
	*concurrency.Mutex
}

func (m *Mutex) Destructor() {
	defer m.session.Close()
}

func (m *Mutex) Lock(ctx context.Context) error {
	return m.Mutex.Lock(ctx)
}

func (m *Mutex) TryLock(ctx context.Context) error {
	return m.Mutex.TryLock(ctx)
}

func (m *Mutex) Unlock(ctx context.Context) error {
	return m.Mutex.Unlock(ctx)
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
