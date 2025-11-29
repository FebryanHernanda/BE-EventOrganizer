package repository

import (
	"context"
	"time"

	modelAuth "github.com/FebryanHernanda/BE-EventOrganizer/internal/models/auth"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (r *AuthRepository) SaveActivationToken(ctx context.Context, token *modelAuth.ActivationToken) error {
	token.CreatedAt = time.Now()

	_, err := r.ActivationCollection.InsertOne(ctx, token)
	if err != nil {
		logrus.WithError(err).Error("failed to save activation token")
	}

	return err
}

func (r *AuthRepository) FindActivationToken(ctx context.Context, token string) (*modelAuth.ActivationToken, error) {
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

func (r *AuthRepository) FindValidActivationToken(ctx context.Context, token string) (*modelAuth.ActivationToken, error) {
	var t modelAuth.ActivationToken

	err := r.ActivationCollection.FindOne(ctx, bson.M{
		"token": token,
		"used":  false,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&t)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.WithField("token", token).Warn("activation token not found or invalid")
			return nil, err
		}
		return nil, err
	}

	return &t, nil
}

func (r *AuthRepository) MarkTokenUsed(ctx context.Context, token string) error {
	_, err := r.ActivationCollection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)

	if err != nil {
		logrus.WithError(err).Error("failed to mark activation token as used")
	}

	return err
}

func (r *AuthRepository) DeleteActivationToken(ctx context.Context, userID string) error {
	_, err := r.ActivationCollection.DeleteMany(ctx, bson.M{"user_id": userID})

	if err != nil {
		logrus.WithError(err).Error("failed to delete user activation tokens")
	}

	return err
}

func (r *AuthRepository) ActivateUser(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"is_active": true}})
	if err != nil {
		logrus.WithError(err).Error("Failed to activate user")
	}
	return err
}
