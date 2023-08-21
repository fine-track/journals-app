package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	ID          primitive.ObjectID  `bson:"_id" json:"_id"`
	Type        string              `bson:"type" json:"type"`
	Title       string              `bson:"title" json:"title"`
	Amount      int32               `bson:"amount" json:"amount"`
	Description string              `bson:"description" json:"description"`
	Date        string              `bson:"date" json:"date"`
	UserId      primitive.ObjectID  `bson:"user_id" json:"user_id"`
	CreatedAt   primitive.Timestamp `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.Timestamp `bson:"updated_at" json:"updated_at"`
}

type RecordsList []Record

const RECORDS_PER_PAGE int64 = 10

func (r *Record) New() error {
	if err := typeCheck(r.Type); err != nil {
		return err
	}
	payload := bson.M{
		"type":        r.Type,
		"title":       r.Title,
		"amount":      r.Amount,
		"description": r.Description,
		"date":        r.Date,
		"user_id":     r.UserId,
		"create_at":   primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0},
		"updated_at":  primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0},
	}
	if result, err := RecordsColl.InsertOne(context.TODO(), payload); err != nil {
		return err
	} else {
		r.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	}

}

func (r *Record) Get(id string) error {
	if ID, err := primitive.ObjectIDFromHex(id); err != nil {
		return err
	} else {
		err = RecordsColl.FindOne(context.TODO(), ID).Decode(r)
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
			"updated_at":  primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0},
		},
	}
	_, err := RecordsColl.UpdateByID(context.TODO(), r.ID, payload)
	return err
}

func (r *Record) Delete(id string) error {
	if ID, err := primitive.ObjectIDFromHex(id); err != nil {
		return err
	} else {
		_, err = RecordsColl.DeleteOne(context.TODO(), ID)
		return err
	}
}

func (rl *RecordsList) ListByType(userId string, t string, pageIdx int64) error {
	filter := bson.M{"type": t, "user_id": userId}
	opts := options.Find().SetLimit(RECORDS_PER_PAGE).SetSkip(pageIdx * 10)
	cursor, err := RecordsColl.Find(context.TODO(), filter, opts)
	if err != nil {
		return err
	}
	results := RecordsList{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return err
	}
	for _, result := range results {
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
	}
	return nil
}

func typeCheck(t string) error {
	if t != "EXPENSE" && t != "INCOME" {
		return fmt.Errorf("type should be either 'EXPENSE' or 'INCOME'")
	}
	return nil
}
