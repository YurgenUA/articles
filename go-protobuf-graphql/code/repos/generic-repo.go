package repos

type GenericRepo[T any] interface {
	Create(*T) *T
	GetList() []*T
	GetOne(uint64) (*T, error)
	Update(uint64, *T) (*T, error)
	DeleteOne(uint64) (bool, error)
}
