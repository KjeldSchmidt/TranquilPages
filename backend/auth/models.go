package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuthState struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	State     string             `bson:"state"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	ExpiresAt primitive.DateTime `bson:"expires_at"`
}
