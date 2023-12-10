package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/rosariocannavo/go_auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUser(username string) (*models.User, error)
	CreateUser(user *models.User) error
	CheckIfUserIsPresent(username string) (bool, error)
	UpdateUserNonce(userID primitive.ObjectID, newNonce string) error
	// Other user-related methods...
}

type userRepo struct {
	client *mongo.Client
}

var redisClient *redis.Client

func NewUserRepository(client *mongo.Client) UserRepository {

	//TODO move to init and unify with ratelimiter
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "localhost", "6379"),
	})
	//
	return &userRepo{
		client: client,
	}
}

func (r *userRepo) FindUser(username string) (*models.User, error) {

	var retrievedUser models.User
	collection := r.client.Database("my_database").Collection("users")

	// Logic to fetch user from redis if present
	cachedData, errred := redisClient.Get(context.Background(), username).Bytes() //change with id
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
		err = redisClient.Set(context.Background(), retrievedUser.Username, dataBytes, 0).Err()
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

// TODO: fix this function
// no redis here cause i need to write
func (r *userRepo) UpdateUserNonce(userID primitive.ObjectID, newNonce string) error {

	collection := r.client.Database("my_database").Collection("users")

	filter := bson.M{"_id": userID} // Assuming userID is the unique identifier for the user
	update := bson.M{"$set": bson.M{"nonce": newNonce}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error updating document:", err)
		return err
	}

	return nil
}

//TO a delete, remember to update the cache
// Implement other methods for user-related operations using MongoDB...
