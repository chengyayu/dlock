package etcd

import (
	"github.com/chengyayu/dlock"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var _ dlock.MutexFactory = (*Factory)(nil)

type Factory struct {
	client *clientv3.Client
}

func NewFactory(client *clientv3.Client) *Factory {
	return &Factory{client: client}
}

func (f *Factory) MakeMutex(key string) (dlock.Mutex, error) {
	ss, err := concurrency.NewSession(f.client)
	if err != nil {
		return nil, errors.Wrap(err, "new etcd session err")
	}

	mutex := concurrency.NewMutex(ss, key)

	return &Mutex{
		client:  f.client,
		session: ss,
		Mutex:   mutex,
	}, nil
}
