package user

type RegisterRequest struct {
	FullName        string `bson:"fullname" json:"fullname" validate:"required"`
	Username        string `bson:"username" json:"username" validate:"required"`
	Email           string `bson:"email" json:"email" validate:"required,email"`
	Password        string `bson:"password" json:"password" validate:"required"`
	ConfirmPassword string `bson:"confirmpassword" json:"confirmpassword" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
