package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/server"
	"github.com/yurgenua/golang-crud-rest-api/repos"

	"google.golang.org/grpc"
)

func main() {
	LoadAppConfig()

	// Create Brand Repository
	var brandRepo repos.GenericRepo[entities.Brand] = repos.NewBrandRepo()

	// push RPC server as goroutine
	go StartRPCServer(&brandRepo)

	// push gRPC-Gateway generated server as goroutine
	go StartRPCGatewayServer()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	log.Println("Termination signal received. Exiting...")
}

func StartRPCGatewayServer() {
	gwmux := runtime.NewServeMux()
	err := crud_brand.RegisterCrudServiceHandlerFromEndpoint(context.Background(), gwmux, ":"+AppConfig.RPCPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal(err)
	}
	gwServer := &http.Server{
		Addr:    ":" + AppConfig.Port,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://localhost:%s\n", AppConfig.Port)
	log.Fatalln(gwServer.ListenAndServe())
}

func StartRPCServer(brandRepo *repos.GenericRepo[entities.Brand]) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", AppConfig.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	crud_brand.RegisterCrudServiceServer(s, server.NewCRUDServiceServer(brandRepo))

	log.Printf("gRPC server listening on port %v\n", AppConfig.RPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
