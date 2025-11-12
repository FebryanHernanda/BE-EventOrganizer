package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName   string             `bson:"fullname" json:"fullname"`
	Username   string             `bson:"username" json:"username"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"-"`
	Role       string             `bson:"role" json:"role"`
	ProfilePic string             `bson:"profile_pic,omitempty" json:"profile_pic,omitempty"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
