# dlock

dlock 的构建目的是为调用者提供一套统一的用户操作界面，来使用分布式锁服务。不同的分布式锁服务供应商以插件的方式接入。

## example 1

你可以手动调用 `m.Destructor()` 来关闭当前锁持有的 session。

```go
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

// acquire lock for s1
if err = m1.Lock(context.TODO()); err != nil {
    log.Fatal(err)
}
fmt.Println("acquired lock for s1")

// do something
doSomething()

// release lock for s1
if err = m1.Unlock(context.TODO()); err != nil {
log.Fatal(err)
}
fmt.Println("released lock for s1")
```

## example 2

当然你也可以调用 `m.Do()` 来执行业务逻辑，这样 m 会自动关闭当前锁持有的 session。

```go
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

m1.Do(context.TODO(), func() error {
    fmt.Println("acquired lock for s1")
	// do something
    doSomething()
    fmt.Println("released lock for s1")
    return nil
})
```