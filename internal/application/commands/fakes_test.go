package commands

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type fakeCustomerRepository struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entities.Customer
}

func newFakeCustomerRepository() *fakeCustomerRepository {
	return &fakeCustomerRepository{byID: make(map[uuid.UUID]*entities.Customer)}
}

func (r *fakeCustomerRepository) Save(ctx context.Context, c *entities.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[c.ID()] = c
	return nil
}

func (r *fakeCustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.byID[id]
	if !ok {
		return nil, domainerrors.ErrNotFound
	}
	return c, nil
}

func (r *fakeCustomerRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.byID {
		if c.Email().Equals(email) {
			return c, nil
		}
	}
	return nil, domainerrors.ErrNotFound
}

func (r *fakeCustomerRepository) List(ctx context.Context, offset, limit int) ([]*entities.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]*entities.Customer, 0, len(r.byID))
	for _, c := range r.byID {
		list = append(list, c)
	}
	return list, nil
}

func (r *fakeCustomerRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.byID[id]
	return ok, nil
}

func (r *fakeCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return domainerrors.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

// fakeProductRepository — mesmo princípio, para repositories.ProductRepository.
type fakeProductRepository struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entities.Product
}

func newFakeProductRepository() *fakeProductRepository {
	return &fakeProductRepository{byID: make(map[uuid.UUID]*entities.Product)}
}

func (r *fakeProductRepository) Save(ctx context.Context, p *entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[p.ID()] = p
	return nil
}

func (r *fakeProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.byID[id]
	if !ok {
		return nil, domainerrors.ErrNotFound
	}
	return p, nil
}

func (r *fakeProductRepository) List(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]*entities.Product, 0, len(r.byID))
	for _, p := range r.byID {
		list = append(list, p)
	}
	return list, nil
}

func (r *fakeProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.byID[id]
	return ok, nil
}

func (r *fakeProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return domainerrors.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

type fakeDispatcher struct {
	dispatched []any
}

func (d *fakeDispatcher) Dispatch(ctx context.Context, evts []any) {
	d.dispatched = append(d.dispatched, evts...)
}
