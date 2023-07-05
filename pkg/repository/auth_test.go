package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSaveRefreshSession(t *testing.T) {
	refreshInsertQuery := `INSERT INTO refresh_sessions \(token_id, user_id, refresh_token, expire_at\) 
VALUES \(\$1, \$2, \$3, \$4\)`
	tests := []struct {
		testName      string
		inputField    domain.RefreshSession
		buildStub     func(mock sqlmock.Sqlmock, input domain.RefreshSession)
		expectedError error
	}{
		{
			testName:   "EmptyInputFieldShouldReturnError",
			inputField: domain.RefreshSession{},
			buildStub: func(mock sqlmock.Sqlmock, input domain.RefreshSession) {
				mock.ExpectExec(refreshInsertQuery).
					WithArgs(input.TokenID, input.UserID, input.RefreshToken, input.ExpireAt).
					WillReturnError(errors.New("insert into refresh_table violate not null constraints"))
			},
			expectedError: errors.New("insert into refresh_table violate not null constraints"),
		},
		{
			testName:   "ValidInputShouldExecuteAndNoError",
			inputField: domain.RefreshSession{TokenID: "token_id", RefreshToken: "refreshTokenString", ExpireAt: time.Now()},
			buildStub: func(mock sqlmock.Sqlmock, input domain.RefreshSession) {
				mock.ExpectExec(refreshInsertQuery).
					WithArgs(input.TokenID, input.UserID, input.RefreshToken, input.ExpireAt).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			test.buildStub(mock, test.inputField)
			authRepo := NewAuthRepository(gormDB)

			err = authRepo.SaveRefreshSession(context.Background(), test.inputField)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestFindRefreshSessionByTokenID(t *testing.T) {
	findRefreshSessionQuery := `SELECT \* FROM refresh_sessions WHERE token_id \= \$1`
	expireAt := time.Now().Add(time.Hour * 1)
	tests := []struct {
		testName               string
		tokenID                string
		expectedRefreshSession domain.RefreshSession
		buildStub              func(mock sqlmock.Sqlmock)
		expectedError          error
	}{
		{
			testName:               "NonExistingTokenIDReturnEmptyRefreshSession",
			tokenID:                "non_existing_token_id",
			expectedRefreshSession: domain.RefreshSession{},
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(findRefreshSessionQuery).
					WithArgs("non_existing_token_id").WillReturnRows(sqlmock.NewRows([]string{"token_id", "refresh_token", "expired_at"}).
					AddRow("", "", time.Time{}))
			},
			expectedError: nil,
		},
		{
			testName: "ExistingTokenIDReturnRefreshSession",
			tokenID:  "existing_token_id",
			expectedRefreshSession: domain.RefreshSession{
				TokenID:      "existing_token_id",
				RefreshToken: "db_refresh_token_token",
				ExpireAt:     expireAt,
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(findRefreshSessionQuery).
					WithArgs("existing_token_id").WillReturnRows(sqlmock.NewRows([]string{"token_id", "refresh_token", "expired_at"}).
					AddRow("existing_token_id", "db_refresh_token_token", expireAt))
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			test.buildStub(mock)

			authRepo := NewAuthRepository(gormDB)
			refreshSession, err := authRepo.FindRefreshSessionByTokenID(context.Background(), test.tokenID)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedRefreshSession.TokenID, refreshSession.TokenID)
			assert.Equal(t, test.expectedRefreshSession.RefreshToken, refreshSession.RefreshToken)
		})
	}

}
