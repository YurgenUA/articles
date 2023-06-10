package repos

import (
	"fmt"
	"os"

	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/protobuf/crud_brand"
	"google.golang.org/protobuf/proto"
)

const STORAGE_FILE = "./brands-storage.pb"

type BrandRepo struct {
	brands []entities.Brand
}

func NewBrandRepo() *BrandRepo {
	var br = BrandRepo{make([]entities.Brand, 0)}
	br.loadFromFileStorage()
	return &br
}

func (b *BrandRepo) Create(partial entities.Brand) entities.Brand {
	newItem := entities.Brand{uint(len(b.brands)) + 1, partial.Name, partial.Year}
	b.brands = append(b.brands, newItem)
	b.saveToFileStorage()
	return newItem
}

func (b *BrandRepo) GetList() []entities.Brand {
	return b.brands
}

func (p *BrandRepo) GetOne(id uint) (entities.Brand, error) {
	for _, it := range p.brands {
		if it.ID == id {
			return it, nil
		}
	}
	return entities.Brand{}, fmt.Errorf("key '%d' not found", id)
}

func (b *BrandRepo) Update(id uint, amended entities.Brand) (entities.Brand, error) {
	for i, it := range b.brands {
		if it.ID == id {
			amended.ID = id
			b.brands = append(b.brands[:i], b.brands[i+1:]...)
			b.brands = append(b.brands, amended)
			b.saveToFileStorage()
			return amended, nil
		}
	}
	return entities.Brand{}, fmt.Errorf("key '%d' not found", amended.ID)
}

func (b *BrandRepo) DeleteOne(id uint) (bool, error) {
	for i, it := range b.brands {
		if it.ID == id {
			b.brands = append(b.brands[:i], b.brands[i+1:]...)
			b.saveToFileStorage()
			return true, nil
		}
	}
	return false, fmt.Errorf("key '%d' not found", id)
}

func (b *BrandRepo) saveToFileStorage() error {

	brandsMessage := &crud_brand.BrandRepo{
		Brands: []*crud_brand.Brand{},
	}

	for _, b := range b.brands {
		brandsMessage.Brands = append(brandsMessage.Brands,
			&crud_brand.Brand{
				Id: uint64(b.ID), Name: b.Name, Year: uint32(b.Year),
			})

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

	for _, brand := range brandsMessage.Brands {
		b.brands = append(b.brands,
			entities.Brand{
				ID: uint(brand.Id), Name: brand.Name, Year: uint(brand.Year),
			})

	}
	return nil
}

func ToProtoBrand(brand entities.Brand) *crud_brand.Brand {
	return &crud_brand.Brand{
		Id: uint64(brand.ID), Name: brand.Name, Year: uint32(brand.Year),
	}
}

func ToBrand(brand *crud_brand.Brand) entities.Brand {
	return entities.Brand{
		ID: uint((*brand).Id), Name: (*brand).Name, Year: uint((*brand).Year),
	}
}
