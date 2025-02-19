package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
)

var (
	// Db queries & statements
	GetProductQuery            = "SELECT id, description, price FROM products WHERE id = ?"
	GetProductsMostSelledQuery = "SELECT COUNT(products.id) as count_total, products.description, ROUND(SUM(products.price), 1) as total FROM products INNER JOIN sales ON sales.product_id = products.id GROUP BY products.id ORDER BY count_total DESC LIMIT 5;"
	StoreProductStatement      = "INSERT INTO products(description, price) VALUES(?, ?)"

	// Errors
	ErrorProductNotFound              = errors.New("product not found")
	ErrorProductPrepareStoreStatement = errors.New("can not prepare store statement")
	ErrorProductExecStoreStatement    = errors.New("error executing store statement")
)

type ProductRepository interface {
	Get(ctx context.Context, id int) (domain.Product, error)
	StoreBulk(ctx context.Context, products []domain.Product) ([]domain.Product, error)
	ProductsMostSelled(ctx context.Context) ([]domain.ProductMostSelledDTO, error)
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

type productRepository struct {
	db *sql.DB
}

func (r *productRepository) Get(ctx context.Context, id int) (domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRowContext(ctx, GetProductQuery, id).Scan(&product.Id, &product.Description, &product.Price)

	if err != nil {
		return domain.Product{}, ErrorProductNotFound
	}

	return product, nil
}

func (r *productRepository) StoreBulk(ctx context.Context, products []domain.Product) ([]domain.Product, error) {
	valueStrings := make([]string, 0, len(products))
	valueArgs := make([]interface{}, 0, len(products)*3)

	for _, product := range products {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, product.Id)
		valueArgs = append(valueArgs, product.Description)
		valueArgs = append(valueArgs, product.Price)
	}

	stmtString := fmt.Sprintf("INSERT INTO products (id, description, price) VALUES %s", strings.Join(valueStrings, ","))
	stmt, err := r.db.PrepareContext(ctx, stmtString)
	if err != nil {
		return nil, ErrorProductPrepareStoreStatement
	}

	defer stmt.Close()

	result, err := r.db.Exec(stmtString, valueArgs...)

	if err != nil {
		return nil, ErrorProductExecStoreStatement
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) ProductsMostSelled(ctx context.Context) ([]domain.ProductMostSelledDTO, error) {
	rows, err := r.db.QueryContext(ctx, GetProductsMostSelledQuery)

	if err != nil {
		return nil, err
	}

	var productsMostSelled []domain.ProductMostSelledDTO

	for rows.Next() {
		var productMostSelled domain.ProductMostSelledDTO
		var id int
		err = rows.Scan(&id, &productMostSelled.Description, &productMostSelled.Total)
		if err != nil {
			return nil, err
		}

		productsMostSelled = append(productsMostSelled, productMostSelled)
	}

	return productsMostSelled, nil
}
