package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repositories.ProductRepository = (*ProductRepositoryPostgres)(nil)

type ProductRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepositoryPostgres {
	return &ProductRepositoryPostgres{
		pool: pool,
	}
}

func (r *ProductRepositoryPostgres) Save(product *entities.Product) error {

	query := `
		INSERT INTO products
		(
			id,
			name,
			description,
			price,
			stock,
			active,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,$2,$3,$4,$5,$6,$7,$8
		)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock(),
		product.IsActive(),
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("saving product: %w", err)
	}
	return err
}

func (r *ProductRepositoryPostgres) FindByID(id string) (*entities.Product, error) {

	query := `
		SELECT
			id,
			name,
			description,
			price,
			stock,
			active,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
	`

	var (
		productID   string
		name        string
		description string
		price       float64
		stock       int
		active      bool
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := r.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&productID,
		&name,
		&description,
		&price,
		&stock,
		&active,
		&createdAt,
		&updatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	product := entities.RebuildProduct(
		productID,
		name,
		description,
		price,
		stock,
		active,
		createdAt,
		updatedAt,
	)

	return product, nil
}

func (r *ProductRepositoryPostgres) List() ([]*entities.Product, error) {
	query := `
		SELECT
			id,
			name,
			description,
			price,
			stock,
			active,
			created_at,
			updated_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*entities.Product

	for rows.Next() {

		var (
			id          string
			name        string
			description string
			price       float64
			stock       int
			active      bool
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&id,
			&name,
			&description,
			&price,
			&stock,
			&active,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		product := entities.RebuildProduct(
			id,
			name,
			description,
			price,
			stock,
			active,
			createdAt,
			updatedAt,
		)

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepositoryPostgres) Exists(id string) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM products
			WHERE id = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ProductRepositoryPostgres) Delete(id string) error {

	query := `
		DELETE
		FROM products
		WHERE id = $1
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		id,
	)

	return err
}
