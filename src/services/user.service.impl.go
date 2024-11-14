package services

import (
	"context"
	"log/slog"

	"github.com/LombardiDaniel/go-gin-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceMongoImpl struct {
	usersCol *mongo.Collection
}

func NewUserServiceMongoImpl(col *mongo.Collection) UserService {
	return &UserServiceMongoImpl{
		usersCol: col,
	}
}

func (s *UserServiceMongoImpl) CreateUser(ctx context.Context, user models.User) error {
	_, err := s.usersCol.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (s *UserServiceMongoImpl) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User

	query := bson.M{
		"username": username,
	}

	err := s.usersCol.FindOne(ctx, query).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserServiceMongoImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	cur, err := s.usersCol.Find(ctx, bson.M{})
	if err != nil {
		return users, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *UserServiceMongoImpl) NoUsersRegistered(ctx context.Context) (bool, error) {

	_, err := s.usersCol.Find(ctx, bson.M{})
	if err == mongo.ErrNilDocument {
		return true, nil
	} else if err != nil {
		slog.Error(err.Error())
		return false, err
	}

	return false, nil
}
