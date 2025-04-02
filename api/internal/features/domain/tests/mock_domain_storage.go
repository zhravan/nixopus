package tests

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/domain/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun"
)

type MockDomainStorage struct {
	mock.Mock
}

type mockTx struct {
	mock.Mock
}

func (m *mockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTx) NewSelect() *bun.SelectQuery {
	args := m.Called()
	return args.Get(0).(*bun.SelectQuery)
}

func (m *mockTx) NewInsert() *bun.InsertQuery {
	args := m.Called()
	return args.Get(0).(*bun.InsertQuery)
}

func (m *mockTx) NewUpdate() *bun.UpdateQuery {
	args := m.Called()
	return args.Get(0).(*bun.UpdateQuery)
}

func (m *mockTx) NewDelete() *bun.DeleteQuery {
	args := m.Called()
	return args.Get(0).(*bun.DeleteQuery)
}

func (m *mockTx) NewCreateTable() *bun.CreateTableQuery {
	args := m.Called()
	return args.Get(0).(*bun.CreateTableQuery)
}

func (m *mockTx) NewDropTable() *bun.DropTableQuery {
	args := m.Called()
	return args.Get(0).(*bun.DropTableQuery)
}

func (m *mockTx) NewCreateIndex() *bun.CreateIndexQuery {
	args := m.Called()
	return args.Get(0).(*bun.CreateIndexQuery)
}

func (m *mockTx) NewDropIndex() *bun.DropIndexQuery {
	args := m.Called()
	return args.Get(0).(*bun.DropIndexQuery)
}

func (m *mockTx) NewTruncateTable() *bun.TruncateTableQuery {
	args := m.Called()
	return args.Get(0).(*bun.TruncateTableQuery)
}

func (m *mockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *mockTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(*sql.Rows), mockArgs.Error(1)
}

func (m *mockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(*sql.Row)
}

func (m *mockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.ExecContext(context.Background(), query, args...)
}

func (m *mockTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *mockTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.QueryRowContext(context.Background(), query, args...)
}

func (m *mockTx) Begin() (*sql.Tx, error) {
	return m.BeginTx(context.Background(), nil)
}

func (m *mockTx) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	mockArgs := m.Called(ctx, opts)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(*sql.Tx), mockArgs.Error(1)
}

func (m *mockTx) Prepare(query string) (*sql.Stmt, error) {
	return m.PrepareContext(context.Background(), query)
}

func (m *mockTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	mockArgs := m.Called(ctx, query)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(*sql.Stmt), mockArgs.Error(1)
}

func NewMockDomainStorage() *MockDomainStorage {
	return &MockDomainStorage{}
}

func (m *MockDomainStorage) CreateDomain(domain *types.Domain) error {
	args := m.Called(domain)
	return args.Error(0)
}

func (m *MockDomainStorage) GetDomain(id string) (*types.Domain, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.Domain), args.Error(1)
}

func (m *MockDomainStorage) UpdateDomain(ID string, Name string) error {
	args := m.Called(ID, Name)
	return args.Error(0)
}

func (m *MockDomainStorage) DeleteDomain(domain *types.Domain) error {
	args := m.Called(domain)
	return args.Error(0)
}

func (m *MockDomainStorage) GetDomains(OrganizationID string, UserID uuid.UUID) ([]types.Domain, error) {
	args := m.Called()
	return args.Get(0).([]types.Domain), args.Error(1)
}

func (m *MockDomainStorage) GetDomainByName(name string) (*types.Domain, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.Domain), args.Error(1)
}

func (m *MockDomainStorage) IsDomainExists(ID string) (bool, error) {
	args := m.Called(ID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDomainStorage) GetDomainOwnerByID(ID string) (string, error) {
	args := m.Called(ID)
	return args.String(0), args.Error(1)
}

func (m *MockDomainStorage) BeginTx() (bun.Tx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return bun.Tx{}, nil
	}
	return args.Get(0).(bun.Tx), args.Error(1)
}

func (m *MockDomainStorage) WithTx(tx bun.Tx) storage.DomainStorageInterface {
	args := m.Called(tx)
	return args.Get(0).(storage.DomainStorageInterface)
}

func (m *MockDomainStorage) WithGetDomainError(id string, err error) *MockDomainStorage {
	m.On("GetDomain", id).Return(nil, err)
	return m
}

func (m *MockDomainStorage) WithGetDomainByNameError(name string, err error) *MockDomainStorage {
	m.On("GetDomainByName", name).Return(nil, err)
	return m
}

func (m *MockDomainStorage) WithGetDomain(id string, domain *types.Domain, err error) *MockDomainStorage {
	m.On("GetDomain", id).Return(domain, err)
	return m
}

func (m *MockDomainStorage) WithUpdateDomain(id string, name string, err error) *MockDomainStorage {
	m.On("UpdateDomain", id, name).Return(err)
	return m
}

func (m *MockDomainStorage) WithGetDomainsError(err error) *MockDomainStorage {
	m.On("GetDomains").Return(nil, err)
	return m
}
