package dlock

import "context"

type Mutex interface {
	Destructor()
	Lock(ctx context.Context) error
	TryLock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Do(ctx context.Context, fn func() error) error
}

type MutexFactory interface {
	MakeMutex(key string) (Mutex, error)
}
