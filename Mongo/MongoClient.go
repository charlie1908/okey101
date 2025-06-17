package Mongo

// go get go.mongodb.org/mongo-driver/mongo
import (
	"context"
	"okey101/Core"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	shared "okey101/Shared"
)

// MongoOpen connects to MongoDB using mongo.Connect directly
func MongoOpen() (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(shared.Config.MONGOURL))
	if err != nil {
		cancel()
		return nil, nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		cancel()
		return nil, nil, err
	}

	return client, ctx, nil
}

// EnsureDatabaseAndCollection ensures that 101Okey database and Log_Entry collection exist.
func EnsureDatabaseAndCollection(client *mongo.Client, ctx context.Context) error {
	db := client.Database("101Okey")

	// Check if collection exists
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	exists := false
	for _, name := range collections {
		if name == "Log_Entry" {
			exists = true
			break
		}
	}

	if !exists {
		// Create collection (no special options for now)
		err = db.CreateCollection(ctx, "Log_Entry", &options.CreateCollectionOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertLogEntry inserts a LogEntry document into 101Okey.Log_Entry collection.
func InsertLogEntry(client *mongo.Client, ctx context.Context, entry LogEntry) (interface{}, error) {
	collection := client.Database("101Okey").Collection("Log_Entry")

	res, err := collection.InsertOne(ctx, entry)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func ConvertCoreTilesToMongoTiles(coreTiles Core.TileBag) []Tile {
	mongoTiles := make([]Tile, len(coreTiles))
	for i, t := range coreTiles {
		mongoTiles[i] = Tile{
			ID:      t.ID,
			Number:  t.Number,
			Color:   t.Color,
			IsJoker: t.IsJoker,
			IsOkey:  t.IsOkey,
			IsOpend: t.IsOpend,
			GroupID: t.GroupID,
			OrderID: t.OrderID,
			X:       t.X,
			Y:       t.Y,
		}
	}
	return mongoTiles
}
