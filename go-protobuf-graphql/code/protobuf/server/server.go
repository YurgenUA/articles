package server

import (
	"context"
	"log"

	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"github.com/yurgenua/golang-crud-rest-api/repos"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CRUDServiceServer struct {
	repo *repos.GenericRepo[crud_brand.Brand]
}

func NewCRUDServiceServer(repo *repos.GenericRepo[crud_brand.Brand]) *CRUDServiceServer {
	server := CRUDServiceServer{
		repo: repo,
	}
	return &server
}

func (c CRUDServiceServer) Create(ctx context.Context, in *crud_brand.CreateRequest) (*crud_brand.CreateResponse, error) {
	brand := (*c.repo).Create(&(*in.Brand))
	response := &crud_brand.CreateResponse{Brand: brand}
	return response, nil
}

func (c CRUDServiceServer) GetOne(_ context.Context, id *wrapperspb.Int64Value) (*crud_brand.GetOneResponse, error) {
	response := &crud_brand.GetOneResponse{}
	brand, err := (*c.repo).GetOne(uint64(id.Value))
	if err != nil {
		log.Printf("failed to get Brand: %v", err)
		return response, err
	}
	response.Brand = brand
	return response, nil
}

func (c CRUDServiceServer) GetList(_ context.Context, _ *emptypb.Empty) (*crud_brand.GetListResponse, error) {
	response := &crud_brand.GetListResponse{}
	for _, brand := range (*c.repo).GetList() {
		response.Brands = append(response.Brands, brand)
	}
	return response, nil
}

func (c CRUDServiceServer) Update(_ context.Context, message *crud_brand.UpdateRequest) (*crud_brand.UpdateResponse, error) {
	response := &crud_brand.UpdateResponse{}
	brand, err := (*c.repo).Update(uint64(message.Id.Value), message.Brand)
	if err != nil {
		log.Printf("failed to update Brand: %v", err)
		return response, err
	}
	response.Brand = brand
	return response, nil
}

func (c CRUDServiceServer) Delete(_ context.Context, in *wrapperspb.Int64Value) (*emptypb.Empty, error) {
	_, err := (*c.repo).DeleteOne(uint64(in.Value))
	if err != nil {
		log.Printf("failed to delete Brand: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
