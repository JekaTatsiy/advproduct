package advproduct_test

import (
	"testing"

	"products/advproduct"
	"products/advproduct/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBrandzonesRepository(t *testing.T) {
	format.MaxLength = 0

	RegisterFailHandler(Fail)
	RunSpecs(t, "AdvProductRepository")
}

func BuildProduct(ID int, position advproduct.Position) *advproduct.Item {
	return &advproduct.Item{
		ID:        ID,
		ProductID: 0,
		Sort:      0,
		Position:  position,
		Active:    true,
	}
}

func ProductArray(product *advproduct.Item) []*advproduct.Item {
	return append(make([]*advproduct.Item, 0), product)
}
func CreateProducts(start, count int, position advproduct.Position) ([]*advproduct.Item, int) {
	prod := make([]*advproduct.Item, 0)
	lastID := 1
	for i := start; i < start+count; i++ {
		prod = append(prod, BuildProduct(i, position))
	}
	return prod, lastID + 1
}

type RepTest struct {
	rep *advproduct.Repository
	db  *gorm.DB
}

func NewRepTest(r advproduct.Repository) *RepTest {
	return &RepTest{rep: &r}
}
func (t *RepTest) addDB(db *gorm.DB) { t.db = db }

func (t *RepTest) mocks(f func(*mock.Mock)) {
	m, e := (*t.rep).(*advproduct_mock.RepositoryMock)
	if !e {
		return
	}
	f(&m.Mock)
}
func (t *RepTest) fixtures(f func(*gorm.DB)) {
	_, e := (*t.rep).(advproduct.RepositoryImpl)
	if !e {
		return
	}
	f(t.db)
}
func (t *RepTest) Update(items []*advproduct.Item) error {
	m, e := (*t.rep).(*advproduct_mock.RepositoryMock)
	if e {
		return m.Update(items)
	}
	r, e := (*t.rep).(advproduct.Repository)
	if e {
		return r.Update(items)
	}
	return nil
}
func (t *RepTest) ListByPosition(position advproduct.Position) ([]*advproduct.Item, error) {
	m, e := (*t.rep).(*advproduct_mock.RepositoryMock)
	if e {
		return m.ListByPosition(position)
	}
	r, e := (*t.rep).(advproduct.Repository)
	if e {
		return r.ListByPosition(position)
	}
	return nil, nil
}

var _ = Describe("AdvProductRepository", func() {
	db, _ := gorm.Open(postgres.Open("host=192.168.10.136 port=5432 user=postgres password=postgres dbname=product sslmode=disable"))
	//r_mock := advproduct_mock.NewMock()
	r_gorm, _ := advproduct.New(db)
	repo := NewRepTest(r_gorm)
	repo.addDB(db)

	repo.mocks(func(m *mock.Mock) {
		// update olways returns nil
		m.On("Update", mock.AnythingOfType("[]*advproduct.Item")).Return(nil)

		// get products
		m.On("ListByPosition", advproduct.Position("b")).Return(ProductArray(BuildProduct(1, advproduct.Position("b"))), nil)

		// new products
		next := 1
		productsX, next := CreateProducts(next, 3, "x")
		productsY, next := CreateProducts(next, 2, "y")

		m.On("ListByPosition", advproduct.Position("x")).Return(productsX, nil).Once()            // first call for x
		m.On("ListByPosition", advproduct.Position("x")).Return(make([]*advproduct.Item, 0), nil) // second call for x
		m.On("ListByPosition", advproduct.Position("y")).Return(productsY, nil)
	})

	repo.fixtures(func(db *gorm.DB) {
		old := make([]*advproduct.Item, 0)
		BeforeSuite(func() {
			db.Table("adv_products").Order("id").Find(&old)
			db.Where("TRUE").Delete(&advproduct.Item{})
		})
		AfterSuite(func() {
			db.Create(&old)
		})
		AfterEach(func() {
			db.Where("TRUE").Delete(&advproduct.Item{})
		})
	})

	Context("Public functions", func() {
		When("add new product", func() {
			It("Success", func() {
				err := repo.Update(ProductArray(BuildProduct(0, advproduct.Position("а"))))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("get products", func() {
			It("Success", func() {
				err := repo.Update(ProductArray(BuildProduct(1, advproduct.Position("b"))))
				Expect(err).ShouldNot(HaveOccurred())

				items, err := repo.ListByPosition(advproduct.Position("b"))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(items)).ShouldNot(Equal(0))
			})
		})
		When("new products", func() {
			It("Success", func() {
				next := 1
				prods, next := CreateProducts(next, 3, "x")
				err := repo.Update(prods)
				Expect(err).ShouldNot(HaveOccurred())

				items, err := repo.ListByPosition(advproduct.Position("x"))
				Expect(len(items)).Should(Equal(3))

				prods, next = CreateProducts(next, 2, "y")
				err = repo.Update(prods)
				Expect(err).ShouldNot(HaveOccurred())

				items, err = repo.ListByPosition("x")
				Expect(len(items)).Should(Equal(0))
				items, err = repo.ListByPosition(advproduct.Position("y"))
				Expect(len(items)).Should(Equal(2))
			})
		})
	})
})
