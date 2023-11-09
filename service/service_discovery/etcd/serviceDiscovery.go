package service_discovery

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var Client *clientv3.Client

func InitializeEtcdClient(etcdClient *clientv3.Client) {
	Client = etcdClient
}

func RegisterService(client *clientv3.Client, serviceName, serviceAddr string) error {
	resp, err := client.Grant(context.Background(), 5)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/services/%s/%s", serviceName, serviceAddr)
	_, err = client.Put(context.Background(), key, serviceAddr, clientv3.WithLease(resp.ID))
	return err
}

func DiscoverServices(client *clientv3.Client, serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, fmt.Sprintf("/services/%s/", serviceName), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, kv := range resp.Kvs {
		addresses = append(addresses, string(kv.Value))
	}

	return addresses, nil
}
