package users

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersRepository_CreateUser(t *testing.T) {
	type args struct {
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		want    *User
		wantErr bool
	}{
		{
			name: "Successful user creation",
			args: args{
				name:  "John Doe",
				email: "john.doe@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				expectedUUID := uuid.NewString() // This will be the UUID we expect to be generated
				createdAt := time.Now().Truncate(time.Millisecond)
				updatedAt := createdAt

				rows := sqlmock.NewRows([]string{"id", "uuid", "name", "email", "created_at", "updated_at"}).
					AddRow(1, expectedUUID, "John Doe", "john.doe@example.com", createdAt, updatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO users (uuid, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, uuid, name, email, created_at, updated_at`,
				)).
					WithArgs(sqlmock.AnyArg(), "John Doe", "john.doe@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want: &User{
				ID:    1,
				Name:  "John Doe",
				Email: "john.doe@example.com",
				// UUID, CreatedAt, UpdatedAt will be checked dynamically
			},
			wantErr: false,
		},
		{
			name: "Error on scanning the queried row",
			args: args{
				name:  "Jane Smith",
				email: "jane.smith@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO users (uuid, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, uuid, name, email, created_at, updated_at`,
				)).
					WithArgs(sqlmock.AnyArg(), "Jane Smith", "jane.smith@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows) // Simulate a scan error like no rows found
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error on inserting the user (simulating UUID error by DB constraint)",
			args: args{
				name:  "Bob Johnson",
				email: "bob.johnson@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO users (uuid, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, uuid, name, email, created_at, updated_at`,
				)).
					WithArgs(sqlmock.AnyArg(), "Bob Johnson", "bob.johnson@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db: unique constraint violation on uuid")) // Simulate a DB error related to UUID
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			r := NewUsersRepository(db)

			tt.mock(mock)

			got, err := r.CreateUser(context.Background(), tt.args.name, tt.args.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.NotEmpty(t, got.UUID)
				assert.False(t, got.CreatedAt.IsZero())
				assert.False(t, got.UpdatedAt.IsZero())
				assert.Equal(t, got.CreatedAt, got.UpdatedAt)

				// Ensure that the mock expectations were met
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
