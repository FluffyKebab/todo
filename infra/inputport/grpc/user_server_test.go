package grpc

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"github.com/FluffyKebab/todo/infra/outputadapter/auth/jwt"
	"github.com/FluffyKebab/todo/infra/outputadapter/log/testlog"
	"github.com/FluffyKebab/todo/infra/outputadapter/storage/mock"
	"github.com/stretchr/testify/require"
)

func TestUserServer(t *testing.T) {
	userIdReturned := "good_user_id"

	s := NewServer(
		testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
		jwt.NewAuthenticator("secret_key"),
		nil,
		mock.UserService{
			CreateUserFunc: func(_ context.Context, req todo.User) (string, error) {
				require.Equal(t, "bob", req.Name)
				return userIdReturned, nil
			},
		},
	)

	lis := runTestServer(t, s)
	client := createUserClient(t, lis)

	res, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{Name: "bob"})
	require.NoError(t, err)
	require.Equal(t, userIdReturned, res.UserID, "user id in response wrong")

	userIdInToken, err := s.userServer.auth.GetUserID(res.Token)
	require.NoError(t, err)
	require.Equal(t, userIdReturned, userIdInToken, "user id in token wrong")
}
