package services

import (
	"context"

	"github.com/fine-track/journals-app/db"
	"github.com/fine-track/journals-app/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type recordsServer struct {
	pb.UnimplementedRecordsServiceServer
}

// Create
func (s *recordsServer) Create(ctx context.Context, req *pb.CreateRecordRequest) (*pb.UpdateRecordResponse, error) {
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}
	record := db.Record{
		Type:        req.Type.String(),
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        req.Date,
		UserId:      userId,
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

	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	r := db.Record{
		UserId:      userId,
		ID:          id,
		Type:        req.Type.String(),
		Date:        req.Date,
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
	recordsList, err := db.GetUserRecords(req.UserId, req.Type.String(), req.Page)
	if err != nil {
		return nil, err
	}

	pbRecords := []*pb.Record{}
	for _, record := range recordsList {
		pbRecords = append(pbRecords, pbRecordFromRecord(record))
	}
	res := &pb.GetRecordsResponse{
		Success:  true,
		Records:  pbRecords,
		NextPage: req.Page + 1,
		Message:  "Records found",
	}
	return res, nil
}

func (s *recordsServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	res := &pb.PingResponse{
		Message:  req.Message,
		Response: "Pong",
	}
	return res, nil
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
		Date:        record.Date,
		UserId:      record.UserId.Hex(),
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
	return r
}

func RegisterRecordsService(s *grpc.Server) {
	pb.RegisterRecordsServiceServer(s, &recordsServer{})
}
