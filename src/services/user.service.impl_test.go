package services

import (
	"context"
	"testing"
	"time"

	"github.com/LombardiDaniel/gopherbase/helpers"
	"github.com/LombardiDaniel/gopherbase/models"
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
			s := &UserServicePgImpl{
				db: pgContainer.DB,
			}
			if err := s.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserServicePgImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserServicePgImpl_CreateUnconfirmedUser(t *testing.T) {
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

	now := time.Now()

	type args struct {
		ctx             context.Context
		unconfirmedUser models.UnconfirmedUser
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"unconfirmed user",
			args{
				ctx,
				models.UnconfirmedUser{
					Email:        "test1@email.com",
					PasswordHash: "hashtest",
					FirstName:    "Test",
					LastName:     "One",
					DateOfBirth:  &now,
					Otp:          "opt-test",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserServicePgImpl{
				db: pgContainer.DB,
			}
			if err := s.CreateUnconfirmedUser(tt.args.ctx, tt.args.unconfirmedUser); (err != nil) != tt.wantErr {
				t.Errorf("UserServicePgImpl.CreateUnconfirmedUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
