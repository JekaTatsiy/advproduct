package advproduct

import (
	"fmt"

	"github.com/pkg/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Position string

const (
	BestOffers Position = "bestoffers"
)

type Item struct {
	ID        int      `gorm:"primaryKey;autoIncrement"`
	ProductID int      `json:"product_id"`
	Sort      int      `json:"sort"`
	Position  Position `json:"position"`
	Active    bool     `gorm:"index" json:"active"`
}

func (i *Item) TableName() string {
	return "adv_products"
}

type RepositoryImpl struct {
	db *gorm.DB
}

type Repository interface {
	ListByPosition(position Position) ([]*Item, error)
	Update(items []*Item) error
}

func (r RepositoryImpl) ListByPosition(position Position) ([]*Item, error) {
	items := make([]*Item, 0)
	result := r.db.Table("adv_products").
		Where("active = true AND position = ?", position).
		Order("sort").Find(&items)
	fmt.Println()
	return items, errors.WithStack(result.Error)
}

func (r RepositoryImpl) Update(items []*Item) error {
	var sort int
	activeIDs := make([]int, 0)
	for _, x := range items {
		activeIDs = append(activeIDs, x.ID)

		x.Sort = sort
		x.Active = true
		sort++
	}

	var result *gorm.DB

	if len(items) > 0 {
		result = r.db.Table("adv_products").Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&items)

		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
	}

	if len(activeIDs) > 0 {
		result = r.db.Table("adv_products").
			Not(map[string]interface{}{"id": activeIDs}).
			Update("active", "false")
	} else {
		result = r.db.Table("adv_products").
			Where("TRUE").
			Update("active", "false")
	}

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}

func New(db *gorm.DB) (Repository, error) {
	repo := &RepositoryImpl{
		db: db,
	}
	return repo, nil
}
