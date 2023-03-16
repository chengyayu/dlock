package etcd

import (
	"context"
	"fmt"
	"github.com/chengyayu/dlock"
	redismutex "github.com/chengyayu/dlock/contrib/mutex/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"log"
	"testing"
)

func TestExampleDo(t *testing.T) {
	cli := goredislib.NewClient(&goredislib.Options{Addr: "localhost:6379"})
	pool := goredis.NewPool(cli)
	rs := redsync.New(pool)

	var f dlock.MutexFactory = redismutex.NewFactory(rs)
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
		if err != nil {
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
