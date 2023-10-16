package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	pb "example.com/grpc-demo/proto"
	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

type expenseTracker struct {
	title  string `gorm:"primarykey"`
	amount int
	date   string
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "newDatabase"
	dbUser := "postgres"
	password := "pass1234"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(expenseTracker{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")
}

// creating a grPC server
var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedMovieServiceServer
}

//implementing RPC methods

func (*server) CreateTracker(ctx context.Context, req *pb.CreateTrackerRequest) (*pb.CreateTrackerResponse, error) {
	fmt.Println("Create Tracker")
	expenseTracker := req.GetMovie()
	expenseTracker.Id = uuid.New().String()

	data := expenseTracker{
		Amount: expenseTracker.GetAmount(),
		Title:  expenseTracker.GetTitle(),
		Date:   expenseTracker.GetDate(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("expense tracker unsuccessful")
	}
	return &pb.CreateMovieResponse{
		expenseTracker: &pb.expenseTracker{
			Amount: expenseTracker.GetAmount(),
			Title:  expenseTracker.GetTitle(),
			Date:   expenseTracker.GetDate(),
		},
	}, nil
}

func (*server) GetMovies(ctx context.Context, req *pb.ReadTrackerRequest) (*pb.ReadTrackerResponse, error) {
	fmt.Println("Read the details")
	expenseTracker1 := []*pb.expenseTracker{}
	res := DB.Find(&expenseTracker1)
	if res.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &pb.ReadTrackerResponse{
		expenseTracker: expenseTracker1,
	}, nil
}
