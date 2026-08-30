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

var _ repositories.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

type userRow struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *UserRepository) Save(ctx context.Context, user *entities.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash(),
		user.Role().String(),
		user.Active(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainerrors.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	const query = `
		SELECT id, email, password_hash, role, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return scanUser(row)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error) {
	const query = `
		SELECT id, email, password_hash, role, active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email.String())
	return scanUser(row)
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	const query = `
		SELECT id, email, password_hash, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*entities.User, 0)
	for rows.Next() {
		var row userRow
		if err := rows.Scan(&row.ID, &row.Email, &row.PasswordHash, &row.Role, &row.Active, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		user, err := toDomainUser(row)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func scanUser(row pgx.Row) (*entities.User, error) {
	var r userRow
	err := row.Scan(&r.ID, &r.Email, &r.PasswordHash, &r.Role, &r.Active, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(r)
}

func toDomainUser(row userRow) (*entities.User, error) {
	email, err := valueobjects.NewEmail(row.Email)
	if err != nil {
		return nil, err
	}

	role := entities.Role(row.Role)
	if !role.IsValid() {
		return nil, domainerrors.ErrInvalidRole
	}

	return entities.RestoreUser(
		row.ID,
		email,
		row.PasswordHash,
		role,
		row.Active,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}
