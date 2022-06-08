package advproduct_mock

import (
	"github.com/stretchr/testify/mock"
	advProductsRepo "products/advproduct"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) ListByPosition(position advProductsRepo.Position) ([]*advProductsRepo.Item, error) {
	args := r.Called(position)
	return args.Get(0).([]*advProductsRepo.Item), args.Error(1)
}

func (r *RepositoryMock) Update(items []*advProductsRepo.Item) error {
	args := r.Called(items)
	return args.Error(0)
}

func NewMock() *RepositoryMock {
	return &RepositoryMock{}
}

