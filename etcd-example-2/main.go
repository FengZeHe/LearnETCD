package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"sync"
	"time"
)

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

func main() {
	var endpoint = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoint)
	defer ser.CloseService()
	ser.WatchService("/web/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			fmt.Println(ser.GetService())
		}
	}

}
