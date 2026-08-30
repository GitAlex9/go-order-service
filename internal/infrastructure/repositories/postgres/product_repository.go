package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ repositories.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) *ProductRepository {
	return &ProductRepository{db: db}
}

type productRow struct {
	ID          uuid.UUID
	Name        string
	Description string
	PriceCents  int64
	Stock       int
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *ProductRepository) Save(ctx context.Context, product *entities.Product) error {
	const query = `
		INSERT INTO products (id, name, description, price_cents, stock, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price_cents = EXCLUDED.price_cents,
			stock = EXCLUDED.stock,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		product.ID(),
		product.Name(),
		product.Description(),
		product.Price().Cents(),
		product.Stock(),
		product.IsActive(),
		product.CreatedAt(),
		product.UpdatedAt(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	const query = `
		SELECT id, name, description, price_cents, stock, active, created_at, updated_at
		FROM products WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return scanProduct(row)
}

func (r *ProductRepository) List(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
	const query = `
		SELECT id, name, description, price_cents, stock, active, created_at, updated_at
		FROM products ORDER BY created_at DESC OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*entities.Product, 0)
	for rows.Next() {
		var pr productRow
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.Description, &pr.PriceCents, &pr.Stock, &pr.Active, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		product, err := productToDomain(pr)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM products WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domainerrors.ErrProductInUse
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func scanProduct(row pgx.Row) (*entities.Product, error) {
	var pr productRow
	err := row.Scan(&pr.ID, &pr.Name, &pr.Description, &pr.PriceCents, &pr.Stock, &pr.Active, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return productToDomain(pr)
}

func productToDomain(pr productRow) (*entities.Product, error) {
	price, err := valueobjects.NewMoney(pr.PriceCents)
	if err != nil {
		return nil, err
	}
	return entities.RebuildProduct(pr.ID, pr.Name, pr.Description, price, pr.Stock, pr.Active, pr.CreatedAt, pr.UpdatedAt), nil
}
