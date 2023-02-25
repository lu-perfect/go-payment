package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"gobank/internal/auth"
	"gobank/internal/util"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := auth.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomName(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.Password, arg.Password)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	// TODO with role
	// createRandomUser(t)
}
