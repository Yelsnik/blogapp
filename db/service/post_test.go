package db

import (
	"context"

	"testing"

	"github.com/Yelsnik/blogapp/util"
	"github.com/stretchr/testify/require"
)

func TestGetPostCollection(t *testing.T) {

	test := testClient
	col := getPostCollection(test)
	require.NotEmpty(t, col)
}

func TestCreatePostAndGetPostByID(t *testing.T) {
	test := testClient
	user := createRandomUser(t)

	args := &Post{
		Title: util.RandomString(6),
		Body:  util.RandomString(10),
		User:  user.ID,
	}
	args.MarshalBSONPost()

	id, err := CreatePost(test, context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get post
	post, err := GetPostByID(test, context.Background(), id)

	require.Equal(t, args.Title, post.Title)
	require.Equal(t, args.Body, post.Body)
	require.Equal(t, args.User, post.User)
	require.NotZero(t, post.CreatedAt)
	require.NotZero(t, post.UpdatedAt)
}
