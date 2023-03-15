package etcd

import (
	"context"
	"fmt"
	"github.com/chengyayu/dlock"
	"github.com/chengyayu/dlock/contrib/mutex/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"testing"
)

func TestExample(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:20002"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	var f dlock.MutexFactory = etcd.NewFactory(cli)
	m1, err := f.MakeMutex("/my-lock")
	if err != nil {
		log.Fatal(err)
	}
	defer m1.Destructor()

	m2, err := f.MakeMutex("/my-lock")
	if err != nil {
		log.Fatal(err)
	}
	defer m2.Destructor()

	// acquire lock for s1
	if err = m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	if err = m2.TryLock(context.TODO()); err == nil {
		log.Fatal("should not acquire lock")
	}
	if err == concurrency.ErrLocked {
		fmt.Println("cannot acquire lock for s2, as already locked in another session")
	}

	if err = m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")
	if err = m2.TryLock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s2")
}

func TestExampleDo(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:20002"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	var f dlock.MutexFactory = etcd.NewFactory(cli)
	m1, err := f.MakeMutex("/my-lock")
	if err != nil {
		log.Fatal(err)
	}

	m2, err := f.MakeMutex("/my-lock")
	if err != nil {
		log.Fatal(err)
	}

	m1.Do(context.TODO(), func() error {
		fmt.Println("acquired lock for s1")
		if err = m2.TryLock(context.TODO()); err == nil {
			log.Fatal("should not acquire lock")
		}
		if err == concurrency.ErrLocked {
			fmt.Println("cannot acquire lock for s2, as already locked in another session")
		}
		fmt.Println("released lock for s1")
		return nil
	})

	if err = m2.TryLock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s2")
}
