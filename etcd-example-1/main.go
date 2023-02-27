package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type ServiceResgiter struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

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

func (s *ServiceResgiter) putKeyWithLease(lease int64) error {
	// set lease time
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		log.Fatal(err)
	}

	// register
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//
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

func (s *ServiceResgiter) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Println("renew lease success", leaseKeepResp)
	}
	log.Println("close renow lease")
}

func (s *ServiceResgiter) CloseService() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	return s.cli.Close()
}

func main() {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewRegisterService(endpoints, "/web/node1", "localhost:8000", 5)
	if err != nil {
		log.Fatal(err)
	}
	go ser.ListenLeaseRespChan()
	select {
	case <-time.After(20 * time.Second):
		ser.CloseService()
	}

}
