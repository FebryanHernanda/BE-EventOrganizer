package repository

import (
	"context"
	"errors"

	modelUser "github.com/FebryanHernanda/BE-EventOrganizer/internal/models/user"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection           *mongo.Collection
	ActivationCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection:           db.Collection("users"),
		ActivationCollection: db.Collection("activation_tokens"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *modelUser.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*modelUser.User, error) {
	var user modelUser.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.WithField("email", email).Warn("User not found in DB")
			return nil, mongo.ErrNoDocuments
		}
		logrus.WithField("email", email).WithError(err).Error("Failed to fetch user by email")
		return nil, err
	}

	logrus.WithField("email", email).Info("User found by email")
	return &user, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*modelUser.User, error) {
	var user modelUser.User
	err := r.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.WithField("username", username).Warn("User not found in DB")
			return nil, mongo.ErrNoDocuments
		}
		logrus.WithField("username", username).WithError(err).Error("Failed to fetch user by email")
		return nil, err
	}
	logrus.WithField("username", username).Info("User found by username")
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*modelUser.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.WithField("userID", id).WithError(err).Error("Invalid ObjectID")
		return nil, err
	}

	var user modelUser.User
	err = r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logrus.WithField("userID", id).Warn("User not found")
			return nil, errors.New("user not found")
		}
		logrus.WithField("userID", id).WithError(err).Error("DB error fetching user")
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ActivateUser(ctx context.Context, userID string) error {
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
