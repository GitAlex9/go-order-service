package mapper

import (
	"reflect"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestCustomerMapper(t *testing.T) {
	email, err := valueobjects.NewEmail("cliente@teste.com")
	if err != nil {
		t.Fatalf("NewEmail() error = %v, want nil", err)
	}
	cpf, err := valueobjects.NewCPF("52998224725")
	if err != nil {
		t.Fatalf("NewCPF() error = %v, want nil", err)
	}
	customer, err := entities.NewCustomer("Cliente Teste", email, cpf)
	if err != nil {
		t.Fatalf("NewCustomer() error = %v, want nil", err)
	}

	tests := []struct {
		name      string
		customer  *entities.Customer
		want      dto.CustomerResponse
		wantSlice []dto.CustomerResponse
	}{
		{
			name:     "single customer",
			customer: customer,
			want: dto.CustomerResponse{
				ID:        customer.ID().String(),
				Name:      customer.Name(),
				Email:     customer.Email().String(),
				CPF:       customer.CPF().Formatted(),
				CreatedAt: customer.CreatedAt(),
				UpdatedAt: customer.UpdatedAt(),
			},
			wantSlice: []dto.CustomerResponse{{
				ID:        customer.ID().String(),
				Name:      customer.Name(),
				Email:     customer.Email().String(),
				CPF:       customer.CPF().Formatted(),
				CreatedAt: customer.CreatedAt(),
				UpdatedAt: customer.UpdatedAt(),
			}},
		},
		{
			name:      "empty slice",
			customer:  nil,
			wantSlice: []dto.CustomerResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.customer != nil {
				got := CustomerToResponse(tt.customer)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("CustomerToResponse() got = %+v, want %+v", got, tt.want)
				}
			}

			gotSlice := CustomersToResponse(func() []*entities.Customer {
				if tt.customer == nil {
					return []*entities.Customer{}
				}
				return []*entities.Customer{tt.customer}
			}())
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Fatalf("CustomersToResponse() got = %+v, want %+v", gotSlice, tt.wantSlice)
			}
		})
	}
}

func TestProductMapper(t *testing.T) {
	price, err := valueobjects.NewMoney(5000)
	if err != nil {
		t.Fatalf("NewMoney() error = %v, want nil", err)
	}
	product, err := entities.NewProduct("Notebook", "Descrição", price, 10)
	if err != nil {
		t.Fatalf("NewProduct() error = %v, want nil", err)
	}

	tests := []struct {
		name      string
		product   *entities.Product
		want      dto.ProductResponse
		wantSlice []dto.ProductResponse
	}{
		{
			name:    "single product",
			product: product,
			want: dto.ProductResponse{
				ID:          product.ID().String(),
				Name:        product.Name(),
				Description: product.Description(),
				Price:       product.Price().Amount(),
				Stock:       product.Stock(),
				Active:      product.IsActive(),
				CreatedAt:   product.CreatedAt(),
				UpdatedAt:   product.UpdatedAt(),
			},
			wantSlice: []dto.ProductResponse{{
				ID:          product.ID().String(),
				Name:        product.Name(),
				Description: product.Description(),
				Price:       product.Price().Amount(),
				Stock:       product.Stock(),
				Active:      product.IsActive(),
				CreatedAt:   product.CreatedAt(),
				UpdatedAt:   product.UpdatedAt(),
			}},
		},
		{
			name:      "empty slice",
			product:   nil,
			wantSlice: []dto.ProductResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.product != nil {
				got := ProductToResponse(tt.product)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("ProductToResponse() got = %+v, want %+v", got, tt.want)
				}
			}

			gotSlice := ProductsToResponse(func() []*entities.Product {
				if tt.product == nil {
					return []*entities.Product{}
				}
				return []*entities.Product{tt.product}
			}())
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Fatalf("ProductsToResponse() got = %+v, want %+v", gotSlice, tt.wantSlice)
			}
		})
	}
}

func TestUserMapper(t *testing.T) {
	email, err := valueobjects.NewEmail("user@teste.com")
	if err != nil {
		t.Fatalf("NewEmail() error = %v, want nil", err)
	}
	user, err := entities.NewUser(email, "Senha123!", entities.RoleCustomer)
	if err != nil {
		t.Fatalf("NewUser() error = %v, want nil", err)
	}

	tests := []struct {
		name      string
		user      *entities.User
		want      dto.UserResponse
		wantSlice []dto.UserResponse
	}{
		{
			name: "single user",
			user: user,
			want: dto.UserResponse{
				ID:        user.ID().String(),
				Email:     user.Email().String(),
				Role:      user.Role().String(),
				Active:    user.Active(),
				CreatedAt: user.CreatedAt(),
				UpdatedAt: user.UpdatedAt(),
			},
			wantSlice: []dto.UserResponse{{
				ID:        user.ID().String(),
				Email:     user.Email().String(),
				Role:      user.Role().String(),
				Active:    user.Active(),
				CreatedAt: user.CreatedAt(),
				UpdatedAt: user.UpdatedAt(),
			}},
		},
		{
			name:      "empty slice",
			user:      nil,
			wantSlice: []dto.UserResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.user != nil {
				got := UserToResponse(tt.user)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("UserToResponse() got = %+v, want %+v", got, tt.want)
				}
			}

			gotSlice := UsersToResponse(func() []*entities.User {
				if tt.user == nil {
					return []*entities.User{}
				}
				return []*entities.User{tt.user}
			}())
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Fatalf("UsersToResponse() got = %+v, want %+v", gotSlice, tt.wantSlice)
			}
		})
	}
}

func TestOrderMapper(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, err := valueobjects.NewMoney(1500)
	if err != nil {
		t.Fatalf("NewMoney() error = %v, want nil", err)
	}
	item, err := entities.NewOrderItem(productID, "Produto Teste", price, 2)
	if err != nil {
		t.Fatalf("NewOrderItem() error = %v, want nil", err)
	}
	order, err := entities.NewOrder(customerID, []entities.OrderItem{*item})
	if err != nil {
		t.Fatalf("NewOrder() error = %v, want nil", err)
	}

	wantItem := dto.OrderItemResponse{
		ProductID:   productID.String(),
		ProductName: item.ProductName(),
		UnitPrice:   item.UnitPrice().Amount(),
		Quantity:    item.Quantity(),
		Subtotal:    item.Subtotal().Amount(),
	}

	wantOrder := dto.OrderResponse{
		ID:         order.ID().String(),
		CustomerID: order.CustomerID().String(),
		Status:     order.Status().String(),
		Items:      []dto.OrderItemResponse{wantItem},
		Total:      order.Total().Amount(),
		CreatedAt:  order.CreatedAt(),
		UpdatedAt:  order.UpdatedAt(),
	}

	tests := []struct {
		name      string
		order     *entities.Order
		want      dto.OrderResponse
		wantSlice []dto.OrderResponse
	}{
		{
			name:  "single order",
			order: order,
			want:  wantOrder,
			wantSlice: []dto.OrderResponse{
				wantOrder,
			},
		},
		{
			name:      "empty slice",
			order:     nil,
			wantSlice: []dto.OrderResponse{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.order != nil {
				got := OrderToResponse(tt.order)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("OrderToResponse() got = %+v, want %+v", got, tt.want)
				}
			}

			orders := []*entities.Order{}
			if tt.order != nil {
				orders = []*entities.Order{tt.order}
			}
			gotSlice := OrdersToResponse(orders)
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Fatalf("OrdersToResponse() got = %+v, want %+v", gotSlice, tt.wantSlice)
			}
		})
	}
}
