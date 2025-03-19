package auth

import (
	"time"
	"tranquil-pages/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthStateRepository struct {
	db *database.Database
}

func NewOAuthStateRepository(db *database.Database) *OAuthStateRepository {
	return &OAuthStateRepository{db: db}
}

func (r *OAuthStateRepository) Create(state *OAuthState) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	state.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	state.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(15 * time.Minute))

	_, err := r.db.GetCollection("oauth_states").InsertOne(ctx, state)
	return err
}

func (r *OAuthStateRepository) FindAndDelete(state string) (*OAuthState, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	var oauthState OAuthState
	err := r.db.GetCollection("oauth_states").FindOneAndDelete(ctx, bson.M{
		"state": state,
		"expires_at": bson.M{
			"$gt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}).Decode(&oauthState)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &oauthState, err
}

func (r *OAuthStateRepository) CleanupExpired() error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	_, err := r.db.GetCollection("oauth_states").DeleteMany(ctx, bson.M{
		"expires_at": bson.M{
			"$lte": primitive.NewDateTimeFromTime(time.Now()),
		},
	})
	return err
}
