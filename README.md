# dlock

dlock 的构建目的是为调用者提供一套统一的用户操作界面，来使用分布式锁服务。不同的分布式锁服务供应商以插件的方式接入。

## 如何选择服务商插件

> 业务还在单机就可以搞定的量级时，那么按照需求使用任意的单机锁方案就可以。 
> 
> 如果发展到了分布式服务阶段，但业务规模不大，qps 很小的情况下，使用哪种锁方案都差不多。如果公司内已有可以使用的 ZooKeeper、etcd 或者 Redis 集群，那么就尽量在不引入新的技术栈的情况下满足业务需求。
>
> 业务发展到一定量级的话，就需要从多方面来考虑了。首先是你的锁是否在任何恶劣的条件下都不允许数据丢失，如果不允许，那么就不要使用 Redis 的 setnx 的简单锁。
>
> 对锁数据的可靠性要求极高的话，那只能使用 etcd 或者 ZooKeeper 这种通过一致性协议保证数据可靠性的锁方案。但可靠的背面往往都是较低的吞吐量和较高的延迟。需要根据业务的量级对其进行压力测试，以确保分布式锁所使用的 etcd 或 ZooKeeper 集群可以承受得住实际的业务请求压力。
> 
> 需要注意的是，etcd 和 Zookeeper 集群是没有办法通过增加节点来提高其性能的。要对其进行横向扩展，只能增加搭建多个集群来支持更多的请求。这会进一步提高对运维和监控的要求。多个集群可能需要引入 proxy，没有 proxy 那就需要业务去根据某个业务 id 来做分片。如果业务已经上线的情况下做扩展，还要考虑数据的动态迁移。这些都不是容易的事情。
> 
>  [《Go语言高级编程》](https://chai2010.cn/advanced-go-programming-book/ch6-cloud/ch6-02-lock.html#627-%E5%A6%82%E4%BD%95%E9%80%89%E6%8B%A9%E5%90%88%E9%80%82%E7%9A%84%E9%94%81)

## 使用方式 1

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

## 使用方式 2

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