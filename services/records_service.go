package services

import (
	"context"
	"time"

	"github.com/fine-track/records-app/db"
	"github.com/fine-track/records-app/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type recordsServer struct {
	pb.UnimplementedRecordsServiceServer
}

// Create
func (s *recordsServer) Create(ctx context.Context, req *pb.CreateRecordRequest) (*pb.UpdateRecordResponse, error) {
	record := db.Record{
		Type:        req.Type.String(),
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
	}
	if err := record.New(); err != nil {
		return nil, err
	} else {
		return &pb.UpdateRecordResponse{
			Success: true,
			Record:  pbRecordFromRecord(record),
		}, nil
	}
}

// Delete
func (s *recordsServer) Delete(ctx context.Context, req *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	r := db.Record{}
	if err := r.Delete(req.Id); err != nil {
		return nil, err
	} else {
		return &pb.DeleteRecordResponse{Success: true}, nil
	}
}

func (s *recordsServer) Update(ctx context.Context, req *pb.Record) (*pb.UpdateRecordResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	r := db.Record{
		ID:          id,
		Type:        req.Type.String(),
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
	}
	if err := r.Update(); err != nil {
		return nil, err
	} else {
		return &pb.UpdateRecordResponse{
			Success: true,
			Record:  pbRecordFromRecord(r),
		}, nil
	}
}

// GetRecords
func (s *recordsServer) GetRecords(ctx context.Context, req *pb.GetRecordsRequest) (*pb.GetRecordsResponse, error) {
	recordsList := db.RecordsList{}
	if err := recordsList.ListByType(req.Type.String(), int64(req.Page)); err != nil {
		return nil, err
	}
	pbRecords := []*pb.Record{}
	for _, record := range recordsList {
		pbRecords = append(pbRecords, pbRecordFromRecord(record))
	}
	res := &pb.GetRecordsResponse{
		Success: true,
		Records: pbRecords,
	}
	return res, nil
}

func timestampToPbTime(ts primitive.Timestamp) *timestamppb.Timestamp {
	return timestamppb.New(time.Unix(int64(ts.T), 0))
}

func strToEnumType(t string) pb.RecordType {
	if t == "EXPENSE" {
		return pb.RecordType_EXPENSE
	} else {
		return pb.RecordType_INCOME
	}
}

func pbRecordFromRecord(record db.Record) *pb.Record {
	r := &pb.Record{
		Id:          record.ID.Hex(),
		Type:        strToEnumType(record.Type),
		Amount:      record.Amount,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   timestampToPbTime(record.CreatedAt),
		UpdatedAt:   timestampToPbTime(record.UpdatedAt),
	}
	return r
}

func RegisterRecordsServer(s *grpc.Server) {
	pb.RegisterRecordsServiceServer(s, &recordsServer{})
}
