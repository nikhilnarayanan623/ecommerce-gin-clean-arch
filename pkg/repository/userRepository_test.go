package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func TestFindUserByColumnNameAndValue(t *testing.T) {
// 	type columnNameAndValue struct {
// 		columnName string
// 		value      any
// 	}

// 	tests := []struct {
// 		testName       string
// 		input          columnNameAndValue
// 		expectedOutput domain.User
// 		buildStub      func(mock sqlmock.Sqlmock, columnName string)
// 		expectedError  error
// 	}{
// 		{
// 			testName: "invalidColumnNameWillReturnError",
// 			input: columnNameAndValue{
// 				columnName: "nonExistingColumn",
// 			},
// 			expectedOutput: domain.User{},
// 			buildStub: func(mock sqlmock.Sqlmock, columnName string) {
// 				expectedQuery := fmt.Sprintf(`SELECT \* FROM users WHERE %s = \$1`, columnName)
// 				mock.ExpectQuery(expectedQuery).
// 					WillReturnError(fmt.Errorf("%s column not existst", columnName))
// 			},
// 			expectedError: fmt.Errorf("%s column not existst", "nonExistingColumn"),
// 		},
// 		{
// 			testName: "validColumnAndNonExistringValueWillReturnEmptyUser",
// 			input: columnNameAndValue{
// 				columnName: "email",
// 				value:      "nonExistringUser@gmail.com",
// 			},
// 			expectedOutput: domain.User{
// 				ID: 0,
// 			},
// 			buildStub: func(mock sqlmock.Sqlmock, columnName string) {
// 				expectedQuery := fmt.Sprintf(`SELECT \* FROM users WHERE %s = \$1`, columnName)
// 				mock.ExpectQuery(expectedQuery).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
// 			},
// 		},
// 		// {
// 		// 	testName: "validColumnNameAndExistingValueWillReturnUser",
// 		// 	input: columnNameAndValue{
// 		// 		columnName: "user_name",

// 		// 	},
// 		// }
// 	}

// 	for _, test := range tests {
// 		t.Run(test.testName, func(t *testing.T) {
// 			db, mock, err := sqlmock.New()
// 			assert.Nil(t, err, "an error '%s' not expected when opening mock database", err)

// 			gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 				Conn: db,
// 			}), &gorm.Config{})
// 			assert.Nil(t, err, "an error '%s' not expected when opening gorm database", err)

// 			test.buildStub(mock, test.input.columnName)

// 			userRepo := NewUserRepository(gormDB)

// 			user, err := userRepo.FindUserByColumnNameAndValue(context.Background(), test.input.columnName, test.input.value)

// 			if test.expectedError == nil {
// 				assert.Nil(t, err)
// 			} else {
// 				assert.Equal(t, test.expectedError, err)
// 			}

// 			assert.Equal(t, test.expectedOutput, user)
// 		})
// 	}
// }

func TestFindUserByEmail(t *testing.T) {
	tests := []struct {
		testName       string
		inputEmail     string
		expectedOutput domain.User
		buildStub      func(mock sqlmock.Sqlmock)
		expectedError  error
	}{
		{
			testName:       "nonExistingEmailReturnEmptyUser",
			inputEmail:     "nonExistingUser@gmail.com",
			expectedOutput: domain.User{},
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT  \* FROM users WHERE email \= \$1`).
					WithArgs("nonExistingUser@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
						AddRow(0, ""))
			},
			expectedError: nil,
		},
		{
			testName:   "exsitingEmailReturnUser",
			inputEmail: "existingUser@gmail.com",
			expectedOutput: domain.User{
				ID:       1,
				Email:    "existingUser@gmail.com",
				UserName: "existingUserUserName",
				Password: "existingUserHashedPassword",
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT  \* FROM users WHERE email \= \$1`).
					WithArgs("existingUser@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "user_name", "password"}).
						AddRow(1, "existingUser@gmail.com", "existingUserUserName", "existingUserHashedPassword"))
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err, "an error '%s' not expected when opening mock database", err)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.Nil(t, err, "an error '%s' not expected when opening gorm database", err)

			test.buildStub(mock)

			userRepo := NewUserRepository(gormDB)

			user, _ := userRepo.FindUserByEmail(context.Background(), test.inputEmail)

			if test.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err, test.expectedError)
			}

			assert.Equal(t, test.expectedOutput, user)
		})
	}
}
