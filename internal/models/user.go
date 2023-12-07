package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username        string             `json:"username"`
	Password        string             `json:"password"`
	MetamaskAddress string             `json:"metamaskaddress"`
	Nonce           string             `json:"nonce"`
	Role            string             `json:"role"`
}
