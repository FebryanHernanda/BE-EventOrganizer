package auth

import "time"

type ActivationToken struct {
	Token     string    `bson:"token" json:"token"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Email     string    `bson:"email" json:"email"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	Used      bool      `bson:"used" json:"used"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
