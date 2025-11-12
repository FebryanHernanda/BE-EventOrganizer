package repository

import (
	"context"
	"time"

	modelAuth "github.com/FebryanHernanda/BE-EventOrganizer/internal/models/auth"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	Collection           *mongo.Collection
	ActivationCollection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		Collection:           db.Collection("users"),
		ActivationCollection: db.Collection("activation_tokens"),
	}
}

/* Activation Account */
func (r *UserRepository) SaveActivationToken(ctx context.Context, token *modelAuth.ActivationToken) error {
	token.CreatedAt = time.Now().Unix()
	_, err := r.ActivationCollection.InsertOne(ctx, token)
	if err != nil {
		logrus.WithError(err).Error("failed to save activation token")
	}
	return err
}

func (r *UserRepository) FindActivationToken(ctx context.Context, token string) (*modelAuth.ActivationToken, error) {
	var t modelAuth.ActivationToken
	err := r.ActivationCollection.FindOne(ctx, bson.M{"token": token}).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.WithField("token", token).Warn("activation token not found")
			return nil, err
		}
		logrus.WithError(err).Error("failed to find activation token")
		return nil, err
	}
	return &t, nil
}

func (r *UserRepository) DeleteActivationToken(ctx context.Context, token string) error {
	_, err := r.ActivationCollection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		logrus.WithError(err).Error("failed to delete activation token")
	}
	return err
}
