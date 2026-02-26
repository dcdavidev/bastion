package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dcdavidev/bastion/packages/core/db"
	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock of the db.Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Close()                                         { m.Called() }
func (m *MockDatabase) Ping(ctx context.Context) error                 { return m.Called(ctx).Error(0) }
func (m *MockDatabase) RunMigrations() error                          { return m.Called().Error(0) }
func (m *MockDatabase) GetMigrationStatus() (uint, bool, error)        { 
	args := m.Called()
	return args.Get(0).(uint), args.Bool(1), args.Error(2)
}
func (m *MockDatabase) HasAdmin(ctx context.Context) (bool, error)     { 
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}
func (m *MockDatabase) CreateUser(ctx context.Context, u, e, h, s, r string) (*models.User, error) {
	args := m.Called(ctx, u, e, h, s, r)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockDatabase) GetUserByUsername(ctx context.Context, u string) (*models.User, string, string, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil { return nil, "", "", args.Error(3) }
	return args.Get(0).(*models.User), args.String(1), args.String(2), args.Error(3)
}
func (m *MockDatabase) GetUserByEmail(ctx context.Context, e string) (*models.User, string, string, error) {
	args := m.Called(ctx, e)
	if args.Get(0) == nil { return nil, "", "", args.Error(3) }
	return args.Get(0).(*models.User), args.String(1), args.String(2), args.Error(3)
}
func (m *MockDatabase) GrantProjectAccess(ctx context.Context, u, p uuid.UUID, k string) error {
	return m.Called(ctx, u, p, k).Error(0)
}
func (m *MockDatabase) GetVaultConfig(ctx context.Context) (*db.VaultConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*db.VaultConfig), args.Error(1)
}
func (m *MockDatabase) InitializeVault(ctx context.Context, w, s string) error {
	return m.Called(ctx, w, s).Error(0)
}
func (m *MockDatabase) CreateClient(ctx context.Context, n string) (*models.Client, error) {
	args := m.Called(ctx, n)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Client), args.Error(1)
}
func (m *MockDatabase) GetClients(ctx context.Context) ([]models.Client, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Client), args.Error(1)
}
func (m *MockDatabase) GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Client), args.Error(1)
}
func (m *MockDatabase) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockDatabase) CreateProject(ctx context.Context, c uuid.UUID, n, k string) (*models.Project, error) {
	args := m.Called(ctx, c, n, k)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Project), args.Error(1)
}
func (m *MockDatabase) GetProjectsByClient(ctx context.Context, c uuid.UUID) ([]models.Project, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]models.Project), args.Error(1)
}
func (m *MockDatabase) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Project), args.Error(1)
}
func (m *MockDatabase) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockDatabase) GetProjectKeyForUser(ctx context.Context, p, u uuid.UUID, a bool) (string, error) {
	args := m.Called(ctx, p, u, a)
	return args.String(0), args.Error(1)
}
func (m *MockDatabase) CreateSecret(ctx context.Context, p uuid.UUID, k, v string) (*models.Secret, error) {
	args := m.Called(ctx, p, k, v)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Secret), args.Error(1)
}
func (m *MockDatabase) GetSecretsByProject(ctx context.Context, p uuid.UUID) ([]models.Secret, error) {
	args := m.Called(ctx, p)
	return args.Get(0).([]models.Secret), args.Error(1)
}
func (m *MockDatabase) GetSecretHistory(ctx context.Context, p uuid.UUID, k string) ([]models.Secret, error) {
	args := m.Called(ctx, p, k)
	return args.Get(0).([]models.Secret), args.Error(1)
}
func (m *MockDatabase) LogEvent(ctx context.Context, a, t string, tid uuid.UUID, meta map[string]interface{}) error {
	return m.Called(ctx, a, t, tid, meta).Error(0)
}
func (m *MockDatabase) GetAuditLogs(ctx context.Context, f db.AuditFilter) ([]models.AuditLog, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]models.AuditLog), args.Error(1)
}

func TestStatusHandler(t *testing.T) {
	// Setup environment
	os.Setenv("BASTION_DATABASE_URL", "postgres://test")
	os.Setenv("BASTION_JWT_SECRET", "test-secret")
	defer os.Unsetenv("BASTION_DATABASE_URL")
	defer os.Unsetenv("BASTION_JWT_SECRET")

	mockDB := new(MockDatabase)
	h := NewHandler(mockDB)

	// Mock expectations
	mockDB.On("Ping", mock.Anything).Return(nil)
	mockDB.On("GetMigrationStatus").Return(uint(5), false, nil)
	mockDB.On("HasAdmin", mock.Anything).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v1/status", nil)
	rr := httptest.NewRecorder()
	
	h.StatusHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp StatusResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.True(t, resp.ConnectedToDB)
	assert.Empty(t, resp.MissingEnvVars)
	assert.Equal(t, uint(5), resp.Migrations.CurrentVersion)
	assert.False(t, resp.Migrations.HasPending)
	assert.True(t, resp.HasAdmin)

	mockDB.AssertExpectations(t)
}

func TestStatusHandler_MissingEnv(t *testing.T) {
	os.Unsetenv("BASTION_DATABASE_URL")
	os.Unsetenv("BASTION_JWT_SECRET")

	mockDB := new(MockDatabase)
	h := NewHandler(mockDB)

	// Ping fails if no DB
	mockDB.On("Ping", mock.Anything).Return(assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/status", nil)
	rr := httptest.NewRecorder()
	
	h.StatusHandler(rr, req)

	var resp StatusResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.False(t, resp.ConnectedToDB)
	assert.Contains(t, resp.MissingEnvVars, "BASTION_DATABASE_URL")
	assert.Contains(t, resp.MissingEnvVars, "BASTION_JWT_SECRET")
}
