# LearnETCD

## ETCD是什么、名字怎么来的？
ETCD的名字是怎么来的？其实源于两个想法，即unix/etc文件夹和分布式系统"d"istibuted。"/etc"文件夹是单个系统存储配置数据的地方，而ETCD存储大规模分布式系统的配置信息，因此"d"istibuted的"/etc"，为"etcd"。ETCD以一致性和容错的方式存储元数据。分布式系统使用ETCD作为一致性键值存储，用于配置管理，服务发现和协调分布工作。使用ETCD的通用分布式模型包括领导选举，分布式锁和监控机器活动。
[//]: # (ETCD是分布式键值对存储，设计用来可靠而快速的保存关键数据并提供访问。通过分布式锁，leader选举和写屏障&#40;write barrieers&#41;来实现可靠的分布式写作。ETCD集群是为高可用，持久性数据存储和检索而准备。)



## ETCD支持的平台
![](./png/platform.png)

### Support tier
- 第 1 层：由etcd 维护者全力支持；etcd 保证通过所有测试，包括功能和压力测试。
- 第 2 层：etcd 保证通过集成和端到端测试，但不一定通过功能或压力测试。
- 第 3 层：保证 etcd 可以构建，可能会被轻微测试（或不被测试），因此它应该被认为是不稳定的。

## ETCD的使用场景及特点
### 使用场景
1. etcd 在稳定性、可靠性和可伸缩性上表现极佳，同时也为云原生应用系统提供了协调机制。
2. etcd 经常用于服务注册与发现的场景，此外还有键值对存储、消息发布与订阅、分布式锁等场景。
### 特点
1. 简单；安装配置简单，而且提供了 HTTP API 进行交互，使用简单。
2. 键值对存储；数据存储在分层组织的目录中，类似于我们日常使用的文件系统。
3. 监测变更；监测特定的键或目录，并对更改进行通知。
4. 安全；支持 SSL 证书验证。
5. 快速；根据官方提供的 benchmark 数据，单实例支持每秒 2k+ 读操作
6. 可靠；基于 Raft 共识算法，实现分布式系统内部数据存储、服务调用的一致性和高可用性。

## ETCD的键值对存储
- 采用kv型数据存储，比关系型数据库快
- 支持内存动态存储和磁盘静态存储
- 分布式存储，可部署多个节点集群
- 存储方式类似目录结构；只有叶子节点才能真正存储数据，相当于文件。叶子节点的父节点一定是目录，目录不能存储数据

## ETCD安装

### 下载安装包

- 下载地址：https://github.com/etcd-io/etcd/releases/

1. 解压  `tar zxvf etcd-v3.5.0-linux-amd64.tar.gz ` 并进入文件夹

   ```shell
   $ ls
   Documentation  etcdctl  README-etcdctl.md  README.md
   etcd           etcdutl  README-etcdutl.md  READMEv2-etcdctl.md
   ```

2. 将 `etcd` 、`etcdctl`复制到 `/usr/local/bin`目录下 ` sudo cp etcd etcdctl /usr/local/bin`

   ```shell
   $ ls /usr/local/bin
   docker-compose  etcd     jsonschema  normalizer  npx  pip3     __pycache__
   dotenv          etcdctl  node        npm         pip  pip3.10  wsdump.py
   ```

3. 创建文件夹并将etcd解压后的文件夹复制进去

   ```shell
   sudo mkdir /usr/local/etcd
   sudo cp -r etcd-v3.5.0-linux-amd64 /usr/local/etcd
   ```

4. 验证

   ```shell
   $ etcd --version
   etcd Version: 3.5.0
   Git SHA: 946a5a6f2
   Go Version: go1.16.3
   Go OS/Arch: linux/amd64
   ```



## API



### 查看API版本

```shell
$ etcdctl version
etcdctl version: 3.5.0
API version: 3.5
```

### 写入、读取、删除key值＆监控值变化

```shell
$ etcdctl put hello "hello,etcd"
OK
$ etcdctl get hello
hello
hello,etcd


$ etcdctl put h1 1
$ etcdctl put h2 2
$ etcdctl put h3 3

// 获取从h1-h4的值
$ etcdctl get h1 h4 --print-value-only
1
2
3
// 获取前缀为h的值
$ etcdctl get --prefix h --print-value-only
1
2
3

// 监控单个值 key的值改变或是key被删除都会被监听到，不存在的key也能监听
$ etcdctl watch xxx

// 监控多个值
$ etcdctl watch -i
watch h1
watch h2

// 监控xx前缀的key
etcdctl watch --prefix xx
```

### 设置、撤销、续租约＆获取租约信息

```shell
//设置租约: 一个key被绑定到一个租约上时，它的生命周期与租约的生命周期绑定
// 一个租约可以绑定多个key
$ etcdctl lease grant 120
lease 694d867d086ac245 granted with TTL(120s)

$ etcdctl put --lease=694d867d086ac245 hello "world"
OK

$ etcdctl get hello --print-value-only
world

// 过了120s 就获取不到hello了

// 主动撤销租约 撤销租约将删除其所有绑定的key
$ etcdctl put hello hi --lease=694d867d086ac24d
OK
$ etcdctl get hello
hello
hi
$ etcdctl lease revoke 694d867d086ac24d
lease 694d867d086ac24d revoked
$ etcdctl get hello
// 无返回
$ 


// 续租约 通过刷新其TTL来保持租约有效性，使其不会过期
$ etcdctl put hello hi --lease=694d867d086ac256
OK
// 自动定时执行续租约，续约成功后每次租约为60秒
$ etcdctl lease keep-alive 694d867d086ac256
lease 694d867d086ac256 keepalived with TTL(60)
lease 694d867d086ac256 keepalived with TTL(60)

// 获取租约信息
$ etcdctl lease grant 999
lease 694d867d086ac259 granted with TTL(999s)
$ etcdctl put hello hi 694d867d086ac259
OK
$ etcdctl lease timetolive --keys 694d867d086ac259
lease 694d867d086ac259 granted with TTL(999s), remaining(972s), attached keys([])

```


## ETCD实现服务发现

### 概述

服务发现需要实现的基本功能：

1. 服务注册：同个service的所有节点注册到相同目录下，节点启动后将自己的信息注册到所属服务目录中。
2. 健康检查：服务节点定时进行健康检查。注册到服务目录中的信息设置一个较短的TTL，运行正常的服务节点每隔一段时间会去更新信息的TTL ，从而达到健康检查效果。
3. 服务发现： 通过服务节点能查询到服务提供外部访问的IP和端口号。

### 服务注册＆健康检查

我们启动一个服务并注册到ETCD中，同时通过绑定租约和自动续租约的方式实现健康检查。
```
// 创建租约服务注册
type ServiceResgiter struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

// 创建注册服务
func NewRegisterService(endpoints []string, key, val string, lease int64) (*ServiceResgiter, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	ser := &ServiceResgiter{
		cli: cli,
		key: key,
		val: val,
	}
  
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}


// 设置租约
func (s *ServiceResgiter) putKeyWithLease(lease int64) error {
	// set lease time
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		log.Fatal(err)
	}

	// 注册&绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// 续约
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	fmt.Println(s.leaseID)
	s.keepAliveChan = leaseRespChan
	fmt.Printf("Put key %s val %s", s.key, s.cli)
	return nil
}

// 关闭服务
func (s *ServiceResgiter) CloseService() error {
	// 撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	return s.cli.Close()
}
```
完整代码 [点这里](https://github.com/FengZeHe/LearnETCD/tree/main/etcd-example-1)

### 服务发现
```
// 服务发现
type ServiceDescovery struct {
	cli        *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
}

// 新建服务发现
func NewServiceDiscovery(endpoints []string) *ServiceDescovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDescovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

// WatchService 查看已有服务＆监听
func (s *ServiceDescovery) WatchService(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}
	go s.watcher(prefix)
	return nil
}

// 根据前缀监听
func (s *ServiceDescovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix : %s", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				s.DelServicelist(string(ev.Kv.Key))
			}
		}
	}
}

// add service address
func (s *ServiceDescovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)
	fmt.Printf("put key: %s val: %s", key, val)
}

// delete service address
func (s *ServiceDescovery) DelServicelist(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	fmt.Printf("delete key %s", key)
}

// Get service address
func (s *ServiceDescovery) GetService() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	address := make([]string, 0)
	for _, v := range s.serverList {
		address = append(address, v)
	}
	return address
}

// Close Service
func (s *ServiceDescovery) CloseService() error {
	return s.cli.Close()
}
```

完整代码 [点这里](https://github.com/FengZeHe/LearnETCD/tree/main/etcd-example-2)
### 实践
1. 运行服务注册代码（第一部分代码）

   ![](./png/step1.png)

2. 运行服务发现代码（第二部分代码）

   当前显示只有 `localhost:8000`的服务注册进来。

   ![](./png/step2.png)

3. 使用`etcdctl`手动添加两个服务。

   添加了 /web/node3和 /web/node4

   ![](./png/step3.png)

4. 查看控制台；看到 `/web/node3`和`/web/node4/`显示已经注册进来

   ![](./png/step5.png)



### 可能遇到的问题

#### 使用go get go.etcd.io/etcd/clientv3使出现的问题

```
go: [go.etcd.io/etcd/client/v3@v3.5.7:](http://go.etcd.io/etcd/client/v3@v3.5.7:) verifying go.mod: [go.etcd.io/etcd/client/v3@v3.5.7/go.mod:](http://go.etcd.io/etcd/client/v3@v3.5.7/go.mod:) checking tree#15675525 against tree#15981595: reading [https://goproxy.io/sumdb/sum.golang.org/tile/8/1/239:](https://goproxy.io/sumdb/sum.golang.org/tile/8/1/239:) 404 Not Found server response: not found
```

#### 解决办法

```shell
go get go.etcd.io/etcd/clientv3@release-3.4
```



### 引用

[1] https://github.com/etcd-io/etcd/issues/12484

[2] https://pkg.go.dev/go.etcd.io/etcd/client/v3#section-readme

[3] https://www.cnblogs.com/FireworksEasyCool/p/12890649.html