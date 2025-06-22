package repositories

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoJobRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	mongoInstance *MongoJobRepository
	mongoOnce     sync.Once
)

func NewMongoJobRepository(uri, dbName string) interfaces.JobRepository {
	mongoOnce.Do(func() {
		ctx := context.Background()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
		}

		if err := client.Ping(ctx, nil); err != nil {
			panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
		}

		mongoInstance = &MongoJobRepository{
			client:     client,
			collection: client.Database(dbName).Collection("jobs"),
		}
	})
	return mongoInstance
}

func (mjr *MongoJobRepository) Save(ctx context.Context, job models.Job) (*models.Job, error) {
	job.CreatedOn = time.Now()

	result, err := mjr.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		job.Id = oid.Hex()
	}

	return &job, nil
}

func (mjr *MongoJobRepository) FindByID(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job

	// Try string ID first
	err := mjr.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err == nil {
		return &job, nil
	}

	// Try ObjectID
	if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
		err = mjr.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&job)
		if err == nil {
			return &job, nil
		}
	}

	return nil, fmt.Errorf("job with ID %s not found", id)
}

func (mjr *MongoJobRepository) FindAll(ctx context.Context) ([]models.Job, error) {
	cursor, err := mjr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job
	for cursor.Next(ctx) {
		var job models.Job
		if err := cursor.Decode(&job); err == nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}
