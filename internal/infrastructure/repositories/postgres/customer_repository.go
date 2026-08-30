package postgres

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ repositories.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	db DBTX
}

func NewCustomerRepository(db DBTX) *CustomerRepository {
	return &CustomerRepository{db: db}
}

type customerRow struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CPF       string
	UserID    *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *CustomerRepository) Save(ctx context.Context, customer *entities.Customer) error {
	const query = `
		INSERT INTO customers (id, name, email, cpf, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			cpf = EXCLUDED.cpf,
			user_id = EXCLUDED.user_id,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		customer.ID(),
		customer.Name(),
		customer.Email().String(),
		customer.CPF().String(),
		customer.UserID(),
		customer.CreatedAt(),
		customer.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainerrors.ErrDuplicateCustomer
		}
		return err
	}
	return nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	const query = `
		SELECT id, name, email, cpf, user_id, created_at, updated_at
		FROM customers
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return scanCustomer(row)

}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.Customer, error) {
	const query = `
		SELECT id, name, email, cpf, user_id, created_at, updated_at
		FROM customers
		WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email.String())
	return scanCustomer(row)
}

func (r *CustomerRepository) List(ctx context.Context, offset, limit int) ([]*entities.Customer, error) {
	const query = `
		SELECT id, name, email, cpf, user_id, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]*entities.Customer, 0)
	for rows.Next() {
		var row customerRow
		if err := rows.Scan(&row.ID, &row.Name, &row.Email, &row.CPF, &row.UserID, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		customer, err := toDomain(row)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *CustomerRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM customers WHERE id = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM customers WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func scanCustomer(row pgx.Row) (*entities.Customer, error) {
	var r customerRow
	err := row.Scan(&r.ID, &r.Name, &r.Email, &r.CPF, &r.UserID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return toDomain(r)
}

func toDomain(row customerRow) (*entities.Customer, error) {
	email, err := valueobjects.NewEmail(row.Email)
	if err != nil {
		return nil, err
	}
	cpf, err := valueobjects.NewCPF(row.CPF)
	if err != nil {
		return nil, err
	}
	return entities.RebuildCustomer(row.ID, row.Name, email, cpf, row.UserID, row.CreatedAt, row.UpdatedAt), nil
}
