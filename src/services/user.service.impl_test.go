package services

import (
	"context"
	"testing"
	"time"

	"github.com/LombardiDaniel/go-gin-template/helpers"
	"github.com/LombardiDaniel/go-gin-template/models"
)

func TestUserServicePgImpl_CreateUser(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := helpers.NewPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	s := NewUserServicePgImpl(pgContainer.DB)

	now := time.Now()

	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"insert new user",
			args{
				ctx,
				models.User{
					Email:        "test1@email.com",
					PasswordHash: "hashtest",
					FirstName:    "Test",
					LastName:     "One",
					DateOfBirth:  &now,
				},
			},
			false,
		},
		{
			"reinsert same user",
			args{
				ctx,
				models.User{
					Email:        "test1@email.com",
					PasswordHash: "hashtest",
					FirstName:    "Test",
					LastName:     "One",
					DateOfBirth:  &now,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserServicePgImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
