package postgres

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserStorage(t *testing.T) {
	s, err := connectToDatabase()
	require.NoError(t, err)

	id, err := s.CreateUser(context.Background(), todo.User{Name: "bob"})
	require.NoError(t, err)
	_, err = uuid.Parse(id)
	require.NoError(t, err)
}
