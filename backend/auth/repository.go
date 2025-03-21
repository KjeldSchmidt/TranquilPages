package auth

import (
	"errors"
	"fmt"
	"time"

	"tranquil-pages/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthStateRepository struct {
	collection *mongo.Collection
}

func NewOAuthStateRepository(db *database.Database) *OAuthStateRepository {
	return &OAuthStateRepository{
		collection: db.GetCollection("oauth_states"),
	}
}

func (r *OAuthStateRepository) Create(state *OAuthState) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	now := time.Now()
	state.CreatedAt = primitive.NewDateTimeFromTime(now)
	state.ExpiresAt = primitive.NewDateTimeFromTime(now.Add(15 * time.Minute))

	_, err := r.collection.InsertOne(ctx, state)
	if err != nil {
		return fmt.Errorf("failed to insert state: %w", err)
	}

	return nil
}

func (r *OAuthStateRepository) FindAndDelete(state string) (*OAuthState, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	filter := bson.M{
		"state": state,
		"expires_at": bson.M{
			"$gt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	var result OAuthState
	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find and delete state: %w", err)
	}

	return &result, nil
}

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(db *database.Database) *TokenRepository {
	return &TokenRepository{
		collection: db.GetCollection("blacklisted_jwts"),
	}
}

func (r *TokenRepository) Blacklist(token string) error {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	blacklistedToken := &BlacklistedToken{
		Token:     token,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := r.collection.InsertOne(ctx, blacklistedToken)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (r *TokenRepository) IsBlacklisted(token string) (bool, error) {
	ctx, cancel := database.WithTimeout()
	defer cancel()

	filter := bson.M{"token": token}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check blacklisted token: %w", err)
	}

	return count > 0, nil
}
