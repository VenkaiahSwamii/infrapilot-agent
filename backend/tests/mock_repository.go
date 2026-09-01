package tests

import (
	"errors"
	"sync"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// mockUserRepository provides an in-memory implementation matching repository.UserRepository's methods
type mockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User // keyed by email
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	user.CreatedAt = time.Now()
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) FindAll() ([]models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *mockUserRepository) Update(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Email]; !exists {
		return gorm.ErrRecordNotFound
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// Helper to create a test user with hashed password
func (m *mockUserRepository) createTestUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "viewer",
		IsActive: true,
	}

	if err := m.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
