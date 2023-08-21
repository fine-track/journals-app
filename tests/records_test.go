package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fine-track/journals-app/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ADDRESS = "localhost:8082"
	USER_ID = "64d1b7e92b3de19c6a478936"
)

// Tests creating a new record on the db
func TestCreateRecord(t *testing.T) {
	conn, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("unable to connect to the service\n%v\n", err)
	}
	t.Cleanup(func() { conn.Close() })
	journalsService := pb.NewRecordsServiceClient(conn)

	// now create a new journal in the db
	payload := &pb.CreateRecordRequest{
		UserId:      USER_ID,
		Type:        pb.RecordType_EXPENSE,
		Amount:      1100,
		Title:       "Testing Records",
		Description: fmt.Sprintf("Testing records on %s", time.Now().Format("DD-MM-YYYY")),
		Date:        time.Now().Format("YYYY-MM-DD"),
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}
	result, err := journalsService.Create(context.TODO(), payload)
	if err != nil {
		t.Errorf("unable to create record\npayload: %v\n%v\n", payload, err)
	}
	if !result.Success {
		t.Errorf("unable to create record\npayload: %v\n%v\n", payload, result.Message)
	}
}

// Tests getting records from the db
func TestGetRecords(t *testing.T) {
	conn, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("unable to connect to the service\n%v\n", err)
	}
	t.Cleanup(func() { conn.Close() })
	journalsService := pb.NewRecordsServiceClient(conn)

	expensesPayload := &pb.GetRecordsRequest{
		UserId: USER_ID,
		Type:   pb.RecordType_EXPENSE,
		Page:   0,
	}
	expensesResult, err := journalsService.GetRecords(context.TODO(), expensesPayload)
	if err != nil {
		t.Errorf("unable to get record\npayload: %v\n%v\n", expensesPayload, err)
	}
	if !expensesResult.Success {
		t.Errorf("unable to get record\npayload: %v\n%v\n", expensesPayload, expensesResult.Message)
	}
	// verify the records have correct types
	for _, r := range expensesResult.Records {
		if r.Type.String() != expensesPayload.Type.String() {
			t.Errorf("type mismatch, requested '%s' got '%s'", expensesPayload.Type.String(), r.Type.String())
		}
	}

	incomesPayload := &pb.GetRecordsRequest{
		UserId: USER_ID,
		Type:   pb.RecordType_EXPENSE,
		Page:   0,
	}
	incomesResult, err := journalsService.GetRecords(context.TODO(), incomesPayload)
	if err != nil {
		t.Errorf("unable to get record\npayload: %v\n%v\n", incomesPayload, err)
	}
	if !incomesResult.Success {
		t.Errorf("unable to get record\npayload: %v\n%v\n", incomesPayload, incomesResult.Message)
	}
	// verify the records have correct types
	for _, r := range incomesResult.Records {
		if r.Type.String() != incomesPayload.Type.String() {
			t.Errorf("type mismatch, requested '%s' got '%s'", incomesPayload.Type.String(), r.Type.String())
		}
	}
}

func TestUpdateARecord(t *testing.T) {
	// TODO
	t.Fail()
}

func TestDeleteARecord(t *testing.T) {
	// TODO
	t.Fail()
}
