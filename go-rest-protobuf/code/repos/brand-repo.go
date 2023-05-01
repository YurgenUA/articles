package repos

import (
	"fmt"

	"github.com/yurgenua/golang-crud-rest-api/entities"
)

const STORAGE_FILE = "./../brands-storage.pd"

type BrandRepo struct {
	brands []entities.Brand
}

func NewBrandRepo() *BrandRepo {
	var br = BrandRepo{make([]entities.Brand, 0)}
	return &br
}

func (b *BrandRepo) Create(partial entities.Brand) entities.Brand {
	newItem := entities.Brand{uint(len(b.brands)) + 1, partial.Name, partial.Year}
	b.brands = append(b.brands, newItem)
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
			return amended, nil
		}
	}
	return entities.Brand{}, fmt.Errorf("key '%d' not found", amended.ID)
}

func (b *BrandRepo) DeleteOne(id uint) (bool, error) {
	for i, it := range b.brands {
		if it.ID == id {
			b.brands = append(b.brands[:i], b.brands[i+1:]...)
			return true, nil
		}
	}
	return false, fmt.Errorf("key '%d' not found", id)
}
