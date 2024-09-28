package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/escalopa/mongo-playground/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	client          *mongo.Client
	database        *mongo.Database
	usersCollection *mongo.Collection
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Create the 'users' collection if it doesn't exist
	_ = client.Database("mydatabase").CreateCollection(ctx, "users", nil)

	s := &Storage{
		client:          client,
		database:        client.Database("mydatabase"),
		usersCollection: client.Database("mydatabase").Collection("users"),
	}

	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *Storage) CreateUser(ctx context.Context, user domain.User) (string, error) {
	resp, err := s.usersCollection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	id := resp.InsertedID.(bson.ObjectID).Hex()

	return id, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = s.usersCollection.
		FindOne(ctx, bson.D{{
			Key:   "_id",
			Value: objID,
		}}).
		Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id string, updatedUser domain.User) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.usersCollection.UpdateOne(ctx, bson.D{{
		Key:   "_id",
		Value: objID,
	}}, bson.D{{
		Key:   "$set",
		Value: updatedUser,
	}})

	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.usersCollection.DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: objID,
	}})

	return err
}

func (s *Storage) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	cur, err := s.usersCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var user domain.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("Found %d users\n", len(users))
	for _, user := range users {
		fmt.Printf(" - %d: [%s | %s]\n", user.ID, user.Name, user.Email)
	}

	return users, nil
}
