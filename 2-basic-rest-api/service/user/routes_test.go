package user

import (
	"basic-rest-api/model/domain"
	"basic-rest-api/model/web"
	"basic-rest-api/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock User Repository
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepo) GetUserByID(id int) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserRepo) CreateUser(user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Unit Test
func TestUserHandler(t *testing.T) {
	userRepo := new(mockUserRepo)
	validate := utils.NewValidate()
	handler := NewHandler(userRepo, validate)

	// Test case: Invalid payload
	t.Run("should fail if user payload invalid", func(t *testing.T) {
		payload := web.RegisterUser{
			Username: "dummy",
			Password: "dummy",
			Email:    "", // Invalid email (empty)
			Name:     "dummy",
		}

		data, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	// Test case: User successfully created
	t.Run("should create user if payload is valid", func(t *testing.T) {
		payload := web.RegisterUser{
			Username: "dummy2",
			Password: "dummy2",
			Email:    "dummy4@gmail.com",
			Name:     "dummy2",
		}

		// Mock repository responses
		userRepo.On("GetUserByEmail", "dummy4@gmail.com").Return(nil, sql.ErrNoRows)
		userRepo.On("CreateUser", mock.AnythingOfType("domain.User")).Return(nil)

		data, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, req)

		// Assertion for the correct response status
		assert.Equal(t, http.StatusCreated, recorder.Code)

		// Optionally, check response body for confirmation message
		expectedBody := `{"message":"user created successfully"}`
		assert.JSONEq(t, expectedBody, recorder.Body.String())
	})
}
