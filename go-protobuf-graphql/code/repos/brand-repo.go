package repos

import (
	"fmt"
	"os"

	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"google.golang.org/protobuf/proto"
)

const STORAGE_FILE = "./brands-storage.pb"

type BrandRepo struct {
	brands []*crud_brand.Brand
}

func NewBrandRepo() *BrandRepo {
	var br = BrandRepo{make([]*crud_brand.Brand, 0)}
	br.loadFromFileStorage()
	return &br
}

func (b *BrandRepo) Create(partial *crud_brand.Brand) *crud_brand.Brand {
	newItem := &crud_brand.Brand{Id: uint64(len(b.brands)) + 1, Name: partial.Name, Year: partial.Year}
	b.brands = append(b.brands, newItem)
	b.saveToFileStorage()
	return newItem
}

func (b *BrandRepo) GetList() []*crud_brand.Brand {
	return b.brands
}

func (p *BrandRepo) GetOne(id uint64) (*crud_brand.Brand, error) {
	for _, it := range p.brands {
		if (*it).Id == id {
			return it, nil
		}
	}
	return &crud_brand.Brand{}, fmt.Errorf("key '%d' not found", id)
}

func (b *BrandRepo) Update(id uint64, amended *crud_brand.Brand) (*crud_brand.Brand, error) {
	for i, it := range b.brands {
		if (*it).Id == id {
			amended.Id = id
			b.brands = append(b.brands[:i], b.brands[i+1:]...)
			b.brands = append(b.brands, amended)
			b.saveToFileStorage()
			return amended, nil
		}
	}
	return &crud_brand.Brand{}, fmt.Errorf("key '%d' not found", amended.Id)
}

func (b *BrandRepo) DeleteOne(id uint64) (bool, error) {
	for i, it := range b.brands {
		if it.Id == id {
			b.brands = append(b.brands[:i], b.brands[i+1:]...)
			b.saveToFileStorage()
			return true, nil
		}
	}
	return false, fmt.Errorf("key '%d' not found", id)
}

func (b *BrandRepo) saveToFileStorage() error {

	brandsMessage := &crud_brand.BrandRepo{
		Brands: b.brands,
	}

	data, err := proto.Marshal(brandsMessage)
	if err != nil {
		return fmt.Errorf("cannot marshal to binary: %w", err)
	}

	err = os.WriteFile(STORAGE_FILE, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

func (b *BrandRepo) loadFromFileStorage() error {

	_, err := os.Stat(STORAGE_FILE)
	if err != nil {
		fmt.Println("storage file is not found, starting with empty storage")
		return nil
	}

	data, err := os.ReadFile(STORAGE_FILE)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	var brandsMessage crud_brand.BrandRepo
	err = proto.Unmarshal(data, &brandsMessage)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary data to protobuf: %w", err)
	}
	b.brands = brandsMessage.Brands
	return nil
}

/*
func ToProtoBrand(brand *crud_brand.Brand) *crud_brand.Brand {
	return &crud_brand.Brand{
		Id: uint64(brand.ID), Name: brand.Name, Year: uint32(brand.Year),
	}
}

func ToBrand(brand *crud_brand.Brand) entities.Brand {
	return entities.Brand{
		ID: uint((*brand).Id), Name: (*brand).Name, Year: uint((*brand).Year),
	}
}
*/
