package main

import (
	"fmt"
	"log"

	service_discovery "rentless-services/service/service_discovery/etcd"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Connect to etcd
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"}, // Assuming etcd is running locally
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	service_discovery.InitializeEtcdClient(client)

	// Register services
	services := map[string]string{
		"user":    "localhost:8080",
		"product": "localhost:8081",
		"review":  "localhost:8082",
		"rental":  "localhost:8083",
	}

	for serviceName, serviceAddr := range services {
		err := service_discovery.RegisterService(service_discovery.Client, serviceName, serviceAddr)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Discover services and print their addresses
	serviceNames := []string{"user", "product", "review", "rental"}
	for _, serviceName := range serviceNames {
		addresses, err := service_discovery.DiscoverServices(service_discovery.Client, serviceName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Discovered %s service addresses: %v\n", serviceName, addresses)
	}
}
