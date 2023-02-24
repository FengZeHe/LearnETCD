package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, key string, val string, leaseID int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}

	// 申请租约设置时间keepalive
	if err := ser.putKeyWithLease(leaseID); err != nil {
		return nil, err
	}
	return ser, nil
}

// 设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 设置租约 定期发送需求请求

	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	log.Println(s.leaseID)
	s.keepAliveChan = leaseRespChan
	log.Printf("Put key %s val %s success!", s.key, s.val)
	return nil
}

// ListenLeaseRespdChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRspChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续约")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//	撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}

func main() {
	var endpoint = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoint, "/web/node1", "localhost:8000", 5)

	if err != nil {
		log.Fatal(err)
	}

	// 监听续租相应chan
	go ser.ListenLeaseRspChan()
	select {}
}
