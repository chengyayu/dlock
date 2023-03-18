package redis

import (
	"context"
	"github.com/chengyayu/dlock"
	redsync "github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
	"time"
)

var _ dlock.Mutex = (*Mutex)(nil)

type Mutex struct {
	*redsync.Mutex
}

func (m *Mutex) Destructor() {}

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

type Option func(mutex *Mutex)

func OptionFunc(opt redsync.Option) Option {
	return func(m *Mutex) {
		opt.Apply(m.Mutex)
	}
}

// WithExpiry can be used to set the expiry of a mutex to the given value.
func WithExpiry(expiry time.Duration) Option {
	return OptionFunc(redsync.WithExpiry(expiry))
}

// WithTries can be used to set the number of times lock acquire is attempted.
func WithTries(tries int) Option {
	return OptionFunc(redsync.WithTries(tries))
}

// WithRetryDelay can be used to set the amount of time to wait between retries.
func WithRetryDelay(delay time.Duration) Option {
	return OptionFunc(redsync.WithRetryDelay(delay))
}

// WithDriftFactor can be used to set the clock drift factor.
func WithDriftFactor(factor float64) Option {
	return OptionFunc(redsync.WithDriftFactor(factor))
}

// WithTimeoutFactor can be used to set the timeout factor.
func WithTimeoutFactor(factor float64) Option {
	return OptionFunc(redsync.WithTimeoutFactor(factor))
}

// WithGenValueFunc can be used to set the custom value generator.
func WithGenValueFunc(genValueFunc func() (string, error)) Option {
	return OptionFunc(redsync.WithGenValueFunc(genValueFunc))
}

// WithValue can be used to assign the random value without having to call lock.
// This allows the ownership of a lock to be "transferred" and allows the lock to be unlocked from elsewhere.
func WithValue(v string) Option {
	return OptionFunc(redsync.WithValue(v))
}
