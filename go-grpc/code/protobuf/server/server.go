package server

import (
	"context"
	"log"

	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/golang_protobuf_brand"
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

func (c CRUDServiceServer) Create(context.Context, *golang_protobuf_brand.ProtoBrandRepo_ProtoBrand) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, nil
}

func (c CRUDServiceServer) GetOne(_ context.Context, id *wrapperspb.Int64Value) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	brand, err := (*c.repo).GetOne(uint(id.Value))
	if err != nil {
		log.Printf("failed to get Brand: %v", err)
		return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, err
	}
	return repos.ToProtoBrand(brand), nil
}

func (c CRUDServiceServer) GetList(_ *emptypb.Empty, stream golang_protobuf_brand.Crud_GetListServer) error {
	for _, brand := range (*c.repo).GetList() {
		if err := stream.Send(repos.ToProtoBrand(brand)); err != nil {
			return err
		}
	}
	return nil
}

func (c CRUDServiceServer) Update(_ context.Context, message *golang_protobuf_brand.UpdateRequest) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	brand, err := (*c.repo).Update(uint(message.ID.Value), repos.ToBrand(message.Brand))
	if err != nil {
		log.Printf("failed to update Brand: %v", err)
		return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, err
	}
	return repos.ToProtoBrand(brand), nil
}

func (c CRUDServiceServer) Delete(_ context.Context, message *wrapperspb.Int64Value) (*wrapperspb.BoolValue, error) {
	success, err := (*c.repo).DeleteOne(uint(message.Value))
	if err != nil {
		log.Printf("failed to delete Brand: %v", err)
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(success), nil
}
