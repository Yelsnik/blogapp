package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/blogapp/util"
	"github.com/stretchr/testify/require"
)

func TestGetUserCollection(t *testing.T) {

	test := testClient

	col := getUserCollection(test)
	require.NotEmpty(t, col)
	require.NotNil(t, col)
}

func createRandomUser(t *testing.T) User {

	test := testClient

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	args := &User{
		Name:     util.RandomString(6),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
	}
	args.MarshalBSONUser()

	id, err := CreateUser(test, context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get user
	user, err := GetUserByID(test, context.Background(), id)

	require.Equal(t, args.Password, user.Password)
	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
