package advproduct_mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/JekaTatsiy/advproduct/advproduct"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) ListByPosition(position advproduct.Position) ([]*advproduct.Item, error) {
	args := r.Called(position)
	return args.Get(0).([]*advproduct.Item), args.Error(1)
}

func (r *RepositoryMock) Update(items []*advproduct.Item) error {
	args := r.Called(items)
	return args.Error(0)
}

func NewMock() *RepositoryMock {
	return &RepositoryMock{}
}

