package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	UserId      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Type        string             `bson:"type" json:"type"`
	Date        string             `bson:"date" json:"date"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Amount      int32              `bson:"amount" json:"amount"`
	CreatedAt   string             `bson:"created_at" json:"created_at"`
	UpdatedAt   string             `bson:"updated_at" json:"updated_at"`
}

const RECORDS_PER_PAGE int64 = 10

func (r *Record) New() error {
	if err := typeCheck(r.Type); err != nil {
		return err
	}
	payload := bson.M{
		"user_id":     r.UserId,
		"type":        r.Type,
		"date":        r.Date,
		"title":       r.Title,
		"description": r.Description,
		"amount":      r.Amount,
		"create_at":   time.Now().UTC().String(),
		"updated_at":  time.Now().UTC().String(),
	}
	if result, err := RecordsColl.InsertOne(context.TODO(), payload); err != nil {
		return err
	} else {
		r.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	}

}

func (r *Record) Get(id string) error {
	if Id, err := primitive.ObjectIDFromHex(id); err != nil {
		return err
	} else {
		err = RecordsColl.FindOne(context.TODO(), bson.M{"_id": Id}).Decode(r)
		return err
	}
}

func (r *Record) Update() error {
	if err := typeCheck(r.Type); err != nil {
		return err
	}
	payload := bson.M{
		"$set": bson.M{
			"title":       r.Title,
			"amount":      r.Amount,
			"description": r.Description,
			"date":        r.Date,
			"type":        r.Type,
			"updated_at":  time.Now().UTC().String(),
		},
	}
	_, err := RecordsColl.UpdateByID(context.TODO(), r.ID, payload)
	return err
}

func (r *Record) Delete(id string) error {
	if ID, err := primitive.ObjectIDFromHex(id); err != nil {
		return err
	} else {
		_, err = RecordsColl.DeleteOne(context.TODO(), bson.M{"_id": ID})
		return err
	}
}

func GetUserRecords(userId string, recordType string, pageIdx int32) ([]Record, error) {
	rl := []Record{}

	objUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return rl, err
	}

	if err := typeCheck(recordType); err != nil {
		return rl, err
	}

	filter := bson.M{"user_id": objUserId}
	opts := options.Find().SetLimit(RECORDS_PER_PAGE).SetSkip(RECORDS_PER_PAGE * int64(pageIdx))
	cursor, err := RecordsColl.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &rl); err != nil {
		return nil, err
	}
	return rl, nil
}

func typeCheck(t string) error {
	if t != "EXPENSE" && t != "INCOME" {
		return fmt.Errorf("type should be either 'EXPENSE' or 'INCOME'")
	}
	return nil
}
