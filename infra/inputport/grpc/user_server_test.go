package grpc

import (
	"context"
	"testing"

	"github.com/FluffyKebab/todo/domain/todo"
	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"github.com/FluffyKebab/todo/infra/outputport/auth/jwt"
	"github.com/FluffyKebab/todo/infra/outputport/log/testlog"
	"github.com/FluffyKebab/todo/infra/outputport/storage/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUserServer(t *testing.T) {
	t.Parallel()

	userIdReturned := "good_user_id"

	userServer := &userServer{
		userService: mock.UserService{
			CreateUserFunc: func(_ context.Context, _ todo.User) (string, error) {
				return userIdReturned, nil
			},
		},
		auth: jwt.NewAuthenticator("secret_key"),
		logger: testlog.Logger{ErrorFunc: func(err error) {
			t.Fatalf("unexpected error: %s", err.Error())
		}},
	}

	go func() {
		err := userServer.run(":8080")
		require.NoError(t, err)
	}()

	cc, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := pb.NewUserServiceClient(cc)
	res, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{Name: "bob"})
	require.NoError(t, err)
	require.Equal(t, userIdReturned, res.UserID, "user id in response wrong")

	userIdInToken, err := userServer.auth.GetUserID(res.Token)
	require.NoError(t, err)
	require.Equal(t, userIdReturned, userIdInToken, "user id in response wrong")
}
