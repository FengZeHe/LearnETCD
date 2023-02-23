# LearnETCD

## ETCD是什么、名字怎么来的？
ETCD的名字是怎么来的？其实源于两个想法，即unix/etc文件夹和分布式系统"d"istibuted。"/etc"文件夹是单个系统存储配置数据的地方，而ETCD存储大规模分布式系统的配置信息，因此"d"istibuted的"/etc"，为"etcd"。ETCD以一致性和容错的方式存储元数据。分布式系统使用ETCD作为一致性键值存储，用于配置管理，服务发现和协调分布工作。使用ETCD的通用分布式模型包括领导选举，分布式锁和监控机器活动。
[//]: # (ETCD是分布式键值对存储，设计用来可靠而快速的保存关键数据并提供访问。通过分布式锁，leader选举和写屏障&#40;write barrieers&#41;来实现可靠的分布式写作。ETCD集群是为高可用，持久性数据存储和检索而准备。)



## ETCD支持的平台
![](./png/platform.png)

### Support tier
第 1 层：由etcd 维护者全力支持；etcd 保证通过所有测试，包括功能和压力测试。
第 2 层：etcd 保证通过集成和端到端测试，但不一定通过功能或压力测试。
第 3 层：保证 etcd 可以构建，可能会被轻微测试（或不被测试），因此它应该被认为是不稳定的。

## ETCD的使用场景及特点
### 使用场景
1. etcd 在稳定性、可靠性和可伸缩性上表现极佳，同时也为云原生应用系统提供了协调机制。
2. etcd 经常用于服务注册与发现的场景，此外还有键值对存储、消息发布与订阅、分布式锁等场景。
### 特点
1. 简单；安装配置简单，而且提供了 HTTP API 进行交互，使用简单。
2. 键值对存储；数据存储在分层组织的目录中，类似于我们日常使用的文件系统。
3. 监测变更;监测特定的键或目录，并对更改进行通知。
4. 安全;支持 SSL 证书验证。
5. 快速;根据官方提供的 benchmark 数据，单实例支持每秒 2k+ 读操作
6. 可靠;基于 Raft 共识算法，实现分布式系统内部数据存储、服务调用的一致性和高可用性。

## ETCD的键值对存储
- 采用kv型数据存储，比关系型数据库快
- 支持内存动态存储和磁盘静态存储
- 分布式存储，可部署多个节点集群
- 存储方式类似目录结构；只有叶子节点才能真正存储数据，相当于文件。叶子节点的父节点一定是目录，目录不能存储数据
