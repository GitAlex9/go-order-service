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

var _ repositories.CustomerRepository = (*CustomerRepositoryPostgres)(nil)

type CustomerRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepositoryPostgres {
	return &CustomerRepositoryPostgres{
		pool: pool,
	}
}

func (r *CustomerRepositoryPostgres) Save(customer *entities.Customer) error {

	query := `
		INSERT INTO customers
		(
			id,
			name,
			email,
			active,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,$2,$3,$4,$5,$6
		)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.IsActive(),
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save customer: %w", err)
	}

	return nil
}

func (r *CustomerRepositoryPostgres) FindByID(id string) (*entities.Customer, error) {

	query := `
		SELECT
			id,
			name,
			email,
			active,
			created_at,
			updated_at
		FROM customers
		WHERE id = $1
	`

	var (
		customerID string
		name       string
		email      string
		active     bool
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&customerID,
		&name,
		&email,
		&active,
		&createdAt,
		&updatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	customer := entities.RebuildCustomer(
		customerID,
		name,
		email,
		active,
		createdAt,
		updatedAt,
	)

	return customer, nil
}

func (r *CustomerRepositoryPostgres) List() ([]*entities.Customer, error) {

	query := `
		SELECT
			id,
			name,
			email,
			active,
			created_at,
			updated_at
		FROM customers
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	defer rows.Close()

	customers := []*entities.Customer{}

	for rows.Next() {

		var (
			id        string
			name      string
			email     string
			active    bool
			createdAt time.Time
			updatedAt time.Time
		)

		err := rows.Scan(
			&id,
			&name,
			&email,
			&active,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customer := entities.RebuildCustomer(
			id,
			name,
			email,
			active,
			createdAt,
			updatedAt,
		)

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating customers: %w", err)
	}

	return customers, nil
}

func (r *CustomerRepositoryPostgres) Exists(id string) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM customers
			WHERE id = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed checking customer existence: %w", err)
	}

	return exists, nil
}

func (r *CustomerRepositoryPostgres) Delete(id string) error {

	query := `
		DELETE
		FROM customers
		WHERE id = $1
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}
