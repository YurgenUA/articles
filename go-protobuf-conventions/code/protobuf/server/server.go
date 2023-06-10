package server

import (
	"context"
	"log"

	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"github.com/yurgenua/golang-crud-rest-api/repos"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CRUDServiceServer struct {
	repo *repos.GenericRepo[entities.Brand]
}

func NewCRUDServiceServer(repo *repos.GenericRepo[entities.Brand]) *CRUDServiceServer {
	server := CRUDServiceServer{
		repo: repo,
	}
	return &server
}

func (c CRUDServiceServer) Create(ctx context.Context, in *crud_brand.CreateRequest) (*crud_brand.CreateResponse, error) {
	brand := (*c.repo).Create(repos.ToBrand(in.Brand))
	response := &crud_brand.CreateResponse{Brand: repos.ToProtoBrand(brand)}
	return response, nil
}

func (c CRUDServiceServer) GetOne(_ context.Context, id *wrapperspb.Int64Value) (*crud_brand.GetOneResponse, error) {
	response := &crud_brand.GetOneResponse{}
	brand, err := (*c.repo).GetOne(uint(id.Value))
	if err != nil {
		log.Printf("failed to get Brand: %v", err)
		return response, err
	}
	response.Brand = repos.ToProtoBrand(brand)
	return response, nil
}

func (c CRUDServiceServer) GetList(_ context.Context, _ *emptypb.Empty) (*crud_brand.GetListResponse, error) {
	response := &crud_brand.GetListResponse{}
	for _, brand := range (*c.repo).GetList() {
		response.Brands = append(response.Brands, repos.ToProtoBrand(brand))
	}
	return response, nil
}

func (c CRUDServiceServer) Update(_ context.Context, message *crud_brand.UpdateRequest) (*crud_brand.UpdateResponse, error) {
	response := &crud_brand.UpdateResponse{}
	brand, err := (*c.repo).Update(uint(message.Id.Value), repos.ToBrand(message.Brand))
	if err != nil {
		log.Printf("failed to update Brand: %v", err)
		return response, err
	}
	response.Brand = repos.ToProtoBrand(brand)
	return response, nil
}

func (c CRUDServiceServer) Delete(_ context.Context, in *wrapperspb.Int64Value) (*emptypb.Empty, error) {
	_, err := (*c.repo).DeleteOne(uint(in.Value))
	if err != nil {
		log.Printf("failed to delete Brand: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
