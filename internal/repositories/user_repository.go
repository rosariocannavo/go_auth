package repositories

import (
	"context"
	"fmt"

	"github.com/rosariocannavo/go_auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUser(user *models.User) (*models.User, error)
	CreateUser(user *models.User) error
	CheckIfUserIsPresent(username string) (bool, error)
	// Other user-related methods...
}

type userRepo struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return &userRepo{
		client: client,
	}
}

func (r *userRepo) FindUser(user *models.User) (*models.User, error) {
	collection := r.client.Database("my_database").Collection("users")
	// Logic to fetch user from the database using MongoDB
	var retrievedUser models.User
	filter := bson.M{"username": user.Username}
	err := collection.FindOne(context.Background(), filter).Decode(&retrievedUser)
	if err != nil {
		return nil, err
	}

	return &retrievedUser, nil
}

func (r *userRepo) CreateUser(user *models.User) error {
	collection := r.client.Database("my_database").Collection("users")

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error inserting document:", err)
	}

	return err
}

func (r *userRepo) CheckIfUserIsPresent(username string) (bool, error) {
	collection := r.client.Database("my_database").Collection("users")

	var user models.User
	filter := bson.M{"username": username}
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // User not found
		}
		return false, err // Other error occurred
	}

	return true, nil // User found
}

// Implement other methods for user-related operations using MongoDB...
