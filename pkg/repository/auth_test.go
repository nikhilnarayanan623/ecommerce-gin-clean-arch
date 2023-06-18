package repository

// import (
// 	"context"
// 	"database/sql/driver"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/google/uuid"
// 	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func TestSaveRefreshSession(t *testing.T) {
// 	refreshInserQuery := `INSERT INTO refresh_sessions \(token_id, user_id, refresh_token, expire_at\) 
// VALUES \(\$1, \$2, \$3, \$4\)`
// 	tests := []struct {
// 		testName      string
// 		inputField    domain.RefreshSession
// 		buildStub     func(mock sqlmock.Sqlmock, input domain.RefreshSession)
// 		expectedError error
// 	}{
// 		{
// 			testName:   "EmptyInputFieldShouldReturnError",
// 			inputField: domain.RefreshSession{},
// 			buildStub: func(mock sqlmock.Sqlmock, input domain.RefreshSession) {
// 				mock.ExpectExec(refreshInserQuery).
// 					WithArgs(input.TokenID, input.UserID, input.RefreshToken, input.ExpireAt).
// 					WillReturnError(errors.New("insert into refresh_table violate not null constraints"))
// 			},
// 			expectedError: errors.New("insert into refresh_table violate not null constraints"),
// 		},
// 		{
// 			testName:   "ValidInputShouldExecuteAndNoError",
// 			inputField: domain.RefreshSession{TokenID: uuid.New(), RefreshToken: "refreshTokenString", ExpireAt: time.Now()},
// 			buildStub: func(mock sqlmock.Sqlmock, input domain.RefreshSession) {
// 				mock.ExpectExec(refreshInserQuery).
// 					WithArgs(input.TokenID, input.UserID, input.RefreshToken, input.ExpireAt).
// 					WillReturnResult(driver.ResultNoRows)
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.testName, func(t *testing.T) {

// 			db, mock, err := sqlmock.New()
// 			assert.NoError(t, err)
// 			defer db.Close()
// 			gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 				Conn: db,
// 			}), &gorm.Config{})
// 			assert.NoError(t, err)

// 			test.buildStub(mock, test.inputField)
// 			authRepo := NewAuthRepository(gormDB)

// 			err = authRepo.SaveRefreshSession(context.Background(), test.inputField)

// 			assert.Equal(t, test.expectedError, err)
// 		})
// 	}
// }

// func TestFindRefreshSessionByTokenID(t *testing.T) {
// 	findRefresSessionQuery := `SELECT \* FROM refresh_sessions WHERE token_id \= \$1`
// 	randomUUID := uuid.New()

// 	tests := []struct {
// 		testName               string
// 		inputTokenID           uuid.UUID
// 		expectedRefreshSession domain.RefreshSession
// 		buildStub              func(mock sqlmock.Sqlmock, inputTokenID uuid.UUID, dbValues domain.RefreshSession)
// 		expectedError          error
// 	}{
// 		{
// 			testName:               "NonExistingTokenIDReturnEmptyRefreshSession",
// 			inputTokenID:           randomUUID,
// 			expectedRefreshSession: domain.RefreshSession{},
// 			buildStub: func(mock sqlmock.Sqlmock, inputTokenID uuid.UUID, inputs domain.RefreshSession) {
// 				mock.ExpectQuery(findRefresSessionQuery).
// 					WithArgs(inputTokenID).WillReturnRows(sqlmock.NewRows([]string{"token_id", "refresh_token", "expired_at"}).
// 					AddRow(uuid.NullUUID{}, "", time.Time{}))
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			testName:     "ExistingTokenIDReturnRefreshSession",
// 			inputTokenID: randomUUID,
// 			expectedRefreshSession: domain.RefreshSession{
// 				TokenID:      randomUUID,
// 				RefreshToken: "db_refresh_token",
// 				ExpireAt:     time.Now().Add(time.Hour * 1),
// 			},
// 			buildStub: func(mock sqlmock.Sqlmock, inputTokenID uuid.UUID, inputs domain.RefreshSession) {
// 				mock.ExpectQuery(findRefresSessionQuery).
// 					WithArgs(inputs.TokenID).WillReturnRows(sqlmock.NewRows([]string{"token_id", "refresh_token", "expired_at"}).
// 					AddRow(inputs.TokenID, inputs.RefreshToken, inputs.ExpireAt))
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, test := range tests {

// 		t.Run(test.testName, func(t *testing.T) {

// 			db, mock, err := sqlmock.New()
// 			assert.NoError(t, err)
// 			defer db.Close()

// 			gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 				Conn: db,
// 			}), &gorm.Config{})
// 			assert.NoError(t, err)

// 			test.buildStub(mock, test.inputTokenID, test.expectedRefreshSession)

// 			authRepo := NewAuthRepository(gormDB)
// 			refreshSession, err := authRepo.FindRefreshSessionByTokenID(context.Background(), test.inputTokenID)
// 			assert.Equal(t, test.expectedError, err)
// 			assert.Equal(t, test.expectedRefreshSession.TokenID, refreshSession.TokenID)
// 			assert.Equal(t, test.expectedRefreshSession.RefreshToken, refreshSession.RefreshToken)
// 		})
// 	}

// }
