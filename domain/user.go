package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Email string        `bson:"email" json:"email"`
}

// func (u User) Marshal() ([]byte, error) {
// 	return bson.Marshal(u)
// }
