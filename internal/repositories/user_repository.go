package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/redis_handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUser(username string) (*models.User, error)
	CreateUser(user *models.User) error
	CheckIfUserIsPresent(username string) (bool, error)
	UpdateUserNonce(userID primitive.ObjectID, newNonce string) error
}

type userRepo struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return &userRepo{
		client: client,
	}
}

func (r *userRepo) FindUser(username string) (*models.User, error) {

	var retrievedUser models.User
	collection := r.client.Database("my_database").Collection("users")

	// Logic to fetch user from redis if present
	cachedData, errred := redis_handler.Client.Get(context.Background(), username).Bytes()
	if errred != nil {
		log.Print(errred)
	}

	if cachedData != nil {
		if err := bson.Unmarshal(cachedData, &retrievedUser); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Cached data:", retrievedUser)

	} else {
		//else search on db
		filter := bson.M{"username": username}
		err := collection.FindOne(context.Background(), filter).Decode(&retrievedUser)
		if err != nil {
			return nil, err
		}

		dataBytes, err := bson.Marshal(retrievedUser)
		if err != nil {
			log.Fatal(err)
		}

		// Cache the data in Redis
		err = redis_handler.Client.Set(context.Background(), retrievedUser.Username, dataBytes, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
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

func (r *userRepo) UpdateUserNonce(userID primitive.ObjectID, newNonce string) error {

	collection := r.client.Database("my_database").Collection("users")

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"nonce": newNonce}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error updating document:", err)
		return err
	}

	return nil
}

//for a delete, remember to update the cache
