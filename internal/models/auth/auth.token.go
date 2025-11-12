package auth

type ActivationToken struct {
	UserID    string `bson:"user_id" json:"user_id"`
	Token     string `bson:"token" json:"token"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
}
