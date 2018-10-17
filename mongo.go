package main

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// ObjectID used in MongoDB
type ObjectID [12]byte

// Counter struct
type Counter struct {
	ID      objectid.ObjectID `bson:"_id"`
	Counter int               `bson:"counter"`
}

func mongoConnect() *mongo.Client {
	// Connect to MongoDB
	conn, err := mongo.Connect(context.Background(), "mongodb://localhost:27017", nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
}

// Check if the track already exists in the database
func urlInMongo(url string, trackColl *mongo.Collection) bool {

	// Read the documents where the trackurl field is equal to url parameter
	cursor, err := trackColl.Find(context.Background(),
		bson.NewDocument(bson.EC.String("trackurl", url)))
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the cursor
	defer cursor.Close(context.Background())

	track := igcTrack{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}
	}

	if track.TrackURL == "" { // If there is an empty field, in this case, `trackurl`, it means the track is not on the database
		return false
	}
	return true
}

// Get trackName from URL
func trackNameFromURL(url string, trackColl *mongo.Collection) string {
	// Get the trackName
	cursor, err := trackColl.Find(context.Background(),
		bson.NewDocument(bson.EC.String("trackurl", url)))

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	dbResult := igcTrack{}

	for cursor.Next(context.Background()) {
		err = cursor.Decode(&dbResult)
		if err != nil {
			log.Fatal(err)
		}
	}

	return dbResult.TrackName
}

// Get track counter from DB
func getTrackCounter(db *mongo.Database) int {
	counter := db.Collection("counter") // `counter` Collection

	cursor, err := counter.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	resCounter := Counter{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resCounter)
		if err != nil {
			log.Fatal(err)
		}
	}
	return resCounter.Counter
}

// Increase the track counter
func increaseTrackCounter(cnt int32, db *mongo.Database) {
	collection := db.Collection("counter") // `counter` Collection

	// This is the way to update the counter field in the document
	// Which is storen in the counter collection
	_, err := collection.UpdateOne(context.Background(), nil,
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int32("counter", cnt+1), // Increase the counter by one
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}
