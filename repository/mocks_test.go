package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockPgxIface is a mock implementation of the PgxIface interface
type MockPgxIface struct {
	mock.Mock
}

func (m *MockPgxIface) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	arguments := m.Called(ctx, sql, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(pgx.Rows), arguments.Error(1)
}

func (m *MockPgxIface) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(pgx.Row)
}

func (m *MockPgxIface) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	arguments := m.Called(ctx, query, args)
	// Handle both pgconn.CommandTag and MockCommandTag
	result := arguments.Get(0)
	if mockTag, ok := result.(MockCommandTag); ok {
		// Convert MockCommandTag to pgconn.CommandTag by creating one
		var tag pgconn.CommandTag
		if mockTag.rowsAffected > 0 {
			tag = pgconn.NewCommandTag("UPDATE " + string(rune(mockTag.rowsAffected+'0')))
		} else {
			tag = pgconn.NewCommandTag("")
		}
		return tag, arguments.Error(1)
	}
	return result.(pgconn.CommandTag), arguments.Error(1)
}

func (m *MockPgxIface) Begin(ctx context.Context) (pgx.Tx, error) {
	arguments := m.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(pgx.Tx), arguments.Error(1)
}

// MockRow implements pgx.Row for testing
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	arguments := m.Called(dest)
	return arguments.Error(0)
}

// MockRows implements pgx.Rows for testing
type MockRows struct {
	mock.Mock
	currentRow int
	data       [][]any
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) Err() error {
	arguments := m.Called()
	return arguments.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	arguments := m.Called()
	return arguments.Get(0).(pgconn.CommandTag)
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	arguments := m.Called()
	return arguments.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Next() bool {
	arguments := m.Called()
	return arguments.Bool(0)
}

func (m *MockRows) Scan(dest ...any) error {
	arguments := m.Called(dest)
	return arguments.Error(0)
}

func (m *MockRows) Values() ([]any, error) {
	arguments := m.Called()
	return arguments.Get(0).([]any), arguments.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	arguments := m.Called()
	return arguments.Get(0).([][]byte)
}

func (m *MockRows) Conn() *pgx.Conn {
	arguments := m.Called()
	if arguments.Get(0) == nil {
		return nil
	}
	return arguments.Get(0).(*pgx.Conn)
}

// MockCommandTag implements pgconn.CommandTag
type MockCommandTag struct {
	rowsAffected int64
}

func (m MockCommandTag) String() string {
	return ""
}

func (m MockCommandTag) RowsAffected() int64 {
	return m.rowsAffected
}

func (m MockCommandTag) Insert() bool {
	return false
}

func (m MockCommandTag) Update() bool {
	return false
}

func (m MockCommandTag) Delete() bool {
	return false
}

func (m MockCommandTag) Select() bool {
	return false
}
