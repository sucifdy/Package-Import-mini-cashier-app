package service

import (
	"a21hc3NpZ25tZW50/database"
	"a21hc3NpZ25tZW50/entity"
	"fmt"
)

// ServiceInterface defines the methods for the service
type ServiceInterface interface {
	AddCart(productName string, quantity int) error
	RemoveCart(productName string) error
	ShowCart() ([]entity.CartItem, error)
	ResetCart() error
	GetAllProduct() ([]entity.Product, error)
	Pay(money int) (entity.PaymentInformation, error)
}

// Service struct contains a reference to the database
type Service struct {
	database database.DatabaseInterface
}

// NewService creates a new service instance
func NewService(database database.DatabaseInterface) *Service {
	return &Service{
		database: database,
	}
}

// AddCart adds a product to the cart with validation
func (s *Service) AddCart(productName string, quantity int) error {
	// Validate quantity
	if quantity <= 0 {
		return fmt.Errorf("invalid quantity")
	}

	// Fetch the product from the database
	product, err := s.database.GetProductByName(productName)
	if err != nil {
		return err // Return error if product not found
	}

	// Fetch existing cart items
	carts, err := s.database.GetCartItems()
	if err != nil {
		return err // Return error if failed to get cart items
	}

	// Check if the product already exists in the cart
	for i, item := range carts {
		if item.ProductName == productName {
			// If it exists, update the quantity
			carts[i].Quantity += quantity
			return s.database.SaveCartItems(carts) // Save the updated cart
		}
	}

	// Add new item to the cart
	newItem := entity.CartItem{
		ProductName: product.Name,
		Price:       product.Price,
		Quantity:    quantity,
	}
	carts = append(carts, newItem)

	// Save the new cart
	err = s.database.SaveCartItems(carts)
	if err != nil {
		return err // Return error if saving fails
	}

	return nil
}

// RemoveCart removes a product from the cart
func (s *Service) RemoveCart(productName string) error {
	// Fetch existing cart items
	carts, err := s.database.GetCartItems()
	if err != nil {
		return err // Return error if failed to get cart items
	}

	// Find and remove the product from the cart
	for i, item := range carts {
		if item.ProductName == productName {
			carts = append(carts[:i], carts[i+1:]...) // Remove product from cart
			return s.database.SaveCartItems(carts)    // Save changes
		}
	}

	return fmt.Errorf("product not found") // Return error if product not found
}

// ShowCart displays the contents of the cart
func (s *Service) ShowCart() ([]entity.CartItem, error) {
	carts, err := s.database.GetCartItems()
	if err != nil {
		return nil, err
	}

	return carts, nil
}

// ResetCart clears all items in the cart
func (s *Service) ResetCart() error {
	// Empty the cart
	carts := []entity.CartItem{}

	// Save empty cart to database
	err := s.database.SaveCartItems(carts)
	if err != nil {
		return err // Return error if saving fails
	}

	return nil
}

// GetAllProduct returns all available products
func (s *Service) GetAllProduct() ([]entity.Product, error) {
	// Fetch all products from the database
	products := s.database.GetProductData() // No error returned, so directly return products
	return products, nil
}

// Pay calculates the total price of the cart, checks payment, and returns payment information
func (s *Service) Pay(money int) (entity.PaymentInformation, error) {
	// Fetch the cart from the database
	carts, err := s.database.GetCartItems()
	if err != nil {
		return entity.PaymentInformation{}, err // Return error if failed to get cart items
	}

	// Calculate total price
	totalPrice := 0
	for _, item := range carts {
		totalPrice += item.Price * item.Quantity
	}

	// Check if the paid amount is sufficient
	if money < totalPrice {
		return entity.PaymentInformation{}, fmt.Errorf("money is not enough")
	}

	// Create payment information
	paymentInfo := entity.PaymentInformation{
		ProductList: carts,
		TotalPrice:  totalPrice,
		MoneyPaid:   money,
		Change:      money - totalPrice,
	}

	// Reset cart after payment
	err = s.ResetCart()
	if err != nil {
		return entity.PaymentInformation{}, err // Return error if failed to reset cart
	}

	return paymentInfo, nil
}
