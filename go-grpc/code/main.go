package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/golang_protobuf_brand"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/server"
	"github.com/yurgenua/golang-crud-rest-api/repos"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

func main() {
	LoadAppConfig()

	// Create Brand Repository
	var brandRepo repos.GenericRepo[entities.Brand] = repos.NewBrandRepo()

	// push RPC server as goroutine
	go StartRPCServer(&brandRepo)

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	RegisterBrandRoutes(router, brandRepo)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterBrandRoutes(router *mux.Router, brandRepo repos.GenericRepo[entities.Brand]) {
	NewGenericRouter[entities.Brand, *repos.BrandRepo]("/api/brands", router, &brandRepo)
}

func StartRPCServer(brandRepo *repos.GenericRepo[entities.Brand]) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", AppConfig.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	golang_protobuf_brand.RegisterCrudServer(s, server.NewCRUDServiceServer(brandRepo))

	log.Printf("gRPC server listening on port %v\n", AppConfig.RPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
