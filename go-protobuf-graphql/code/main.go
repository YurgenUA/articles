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
	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/server"
	"github.com/yurgenua/golang-crud-rest-api/repos"

	graphql_runtime "github.com/ysugimoto/grpc-graphql-gateway/runtime"

	"google.golang.org/grpc"
)

func main() {
	LoadAppConfig()

	// Create Brand Repository
	var brandRepo repos.GenericRepo[crud_brand.Brand] = repos.NewBrandRepo()

	// push RPC server as goroutine
	go StartRPCServer(&brandRepo)

	// push gRPC-Gateway generated server as goroutine
	go StartRPCGatewayServer()

	// push gRPC-GrsaphQL generated server as goroutine
	go StartRPCGraphQLServer()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	log.Println("Termination signal received. Exiting...")
}

func StartRPCGraphQLServer() {
	mux := graphql_runtime.NewServeMux()
	err := crud_brand.RegisterCrudServiceGraphql(mux)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", mux)
	log.Printf("Serving gRPC-GraphQL on http://localhost:%s\n", AppConfig.GraphqlPort)
	log.Fatalln(http.ListenAndServe(":"+AppConfig.GraphqlPort, nil))
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

func StartRPCServer(brandRepo *repos.GenericRepo[crud_brand.Brand]) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", AppConfig.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	crud_brand.RegisterCrudServiceServer(s, server.NewCRUDServiceServer(brandRepo))

	log.Printf("Serving gRPC on http://localhost:%s\n", AppConfig.RPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
