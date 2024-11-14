package services

import (
	"context"

	"github.com/LombardiDaniel/go-gin-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServiceImpl struct {
	sessionsCol *mongo.Collection
}

func NewAuthServiceImpl(sessionsCol *mongo.Collection) AuthService {
	return &AuthServiceImpl{
		sessionsCol: sessionsCol,
	}
}

func (s *AuthServiceImpl) Authenticate(ctx context.Context, key string) error {
	var token models.Token

	query := bson.M{
		"token": key,
	}

	err := s.sessionsCol.FindOne(ctx, query).Decode(&token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) CreateToken(ctx context.Context, token models.Token) error {
	_, err := s.sessionsCol.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return err
}
