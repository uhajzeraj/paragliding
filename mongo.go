package main

import (
	"context"
	"log"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// ObjectID used in MongoDB
type ObjectID [12]byte

// Counter struct
type Counter struct {
	ID             objectid.ObjectID `bson:"_id"`
	TrackCounter   int               `bson:"trackCounter"`
	WebhookCounter int               `bson:"webhookCounter"`
}

// Webhook struct
type Webhook struct {
	ID              objectid.ObjectID `bson:"_id"`
	WebhookID       string            `bson:"webhookID"`
	WebhookURL      string            `bson:"webhookURL"`
	MinTriggerValue int               `bson:"minTriggerValue"`
}

func mongoConnect() *mongo.Client {
	// Connect to MongoDB
	conn, err := mongo.Connect(context.Background(), "mongodb://uranh:connecttome123@ds137913.mlab.com:37913/paragliding", nil)
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
func getCounter(db *mongo.Database, whichCounter string) int {
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
	if whichCounter == "trackCounter" {
		return resCounter.TrackCounter
	} else if whichCounter == "webhookCounter" {
		return resCounter.WebhookCounter
	}
	return -1
}

// Increase the track counter
func increaseCounter(cnt int32, db *mongo.Database, whichCounter string) {
	collection := db.Collection("counter") // `counter` Collection

	// This is the way to update the counter field in the document
	// Which is storen in the counter collection
	_, err := collection.UpdateOne(context.Background(), nil,
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int32(whichCounter, cnt+1), // Increase the counter by one
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// Get all tracks
func getAllTracks(client *mongo.Client) []igcTrack {
	db := client.Database("paragliding") // `paragliding` Database
	collection := db.Collection("track") // `track` Collection

	var cursor mongo.Cursor
	var err error
	cursor, err = collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	resTracks := []igcTrack{}
	resTrack := igcTrack{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
		resTracks = append(resTracks, resTrack) // Append each resTrack to resTracks slice
	}

	return resTracks
}

// Get all webhooks
func getAllWebhooks(client *mongo.Client) []Webhook {
	db := client.Database("paragliding")   // `paragliding` Database
	collection := db.Collection("webhook") // `webhook` Collection

	var cursor mongo.Cursor
	var err error
	cursor, err = collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	resWebhooks := []Webhook{}
	resWebhook := Webhook{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resWebhook)
		if err != nil {
			log.Fatal(err)
		}
		resWebhooks = append(resWebhooks, resWebhook) // Append each resTrack to resTracks slice
	}

	return resWebhooks
}

// Get track
func getTrack(client *mongo.Client, url string) igcTrack {
	db := client.Database("paragliding") // `paragliding` Database
	collection := db.Collection("track") // `track` Collection

	cursor, err := collection.Find(context.Background(),
		bson.NewDocument(bson.EC.String("trackname", url)))

	if err != nil {
		log.Fatal(err)
	}

	resTrack := igcTrack{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resTrack

}

// Get webhook
func getWebhook(client *mongo.Client, webhookID string) Webhook {
	db := client.Database("paragliding")   // `paragliding` Database
	collection := db.Collection("webhook") // `webhook` Collection

	cursor, err := collection.Find(context.Background(),
		bson.NewDocument(bson.EC.String("webhookID", webhookID)))

	if err != nil {
		log.Fatal(err)
	}

	resWebhook := Webhook{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resWebhook)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resWebhook

}

// Delete webhook
func deleteWebhook(client *mongo.Client, webhookID string) {
	db := client.Database("paragliding")   // `paragliding` Database
	collection := db.Collection("webhook") // `webhook` Collection

	// Delete the webhook
	collection.DeleteOne(
		context.Background(), bson.NewDocument(
			bson.EC.String("webhookID", webhookID),
		),
	)
}

// Delete all tracks
func deleteAllTracks(client *mongo.Client) {
	db := client.Database("paragliding") // `paragliding` Database
	collection := db.Collection("track") // `track` Collection

	// Delete the tracks
	collection.DeleteMany(context.Background(), bson.NewDocument())

	// Reset the track counter
	increaseCounter(int32(0), db, "trackCounter")

}

// Insert or Update the webhook
func insertUpdateWebhook(data map[string]interface{}) string {

	conn := mongoConnect()
	db := conn.Database("paragliding") // `paragliding` Database
	coll := db.Collection("webhook")   // `webhook` Collection

	// Check if Webhook exists
	cursor, err := coll.Find(context.Background(),
		bson.NewDocument(bson.EC.String("webhookURL", data["webhookURL"].(string))))
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the cursor
	defer cursor.Close(context.Background())

	var paraglide map[string]interface{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&paraglide)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check the counter of webhook
	whCounter := getCounter(db, "webhookCounter")

	webhookName := "webhook" + strconv.Itoa(whCounter)

	// If it is nil, it means we can add the webhook
	if paraglide["webhookURL"] == nil {

		// Insert webhook
		_, err := coll.InsertOne(context.Background(),
			bson.NewDocument(
				bson.EC.String("webhookID", webhookName),
				bson.EC.String("webhookURL", data["webhookURL"].(string)),
				bson.EC.Int32("minTriggerValue", int32(data["minTriggerValue"].(float64))),
			))

		if err != nil {
			log.Fatal(err)
		}

		// Increase webhook counter number
		increaseCounter(int32(whCounter), db, "webhookCounter")

	} else { // If the webhook exists, just update it

		_, err := coll.UpdateOne(context.Background(),
			bson.NewDocument(
				bson.EC.String("webhookURL", data["webhookURL"].(string)),
			),
			bson.NewDocument(
				bson.EC.SubDocumentFromElements("$set",
					bson.EC.Int32("minTriggerValue", int32(data["minTriggerValue"].(float64))),
				),
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		// WebhookID is different in this case
		webhookName = paraglide["webhookID"].(string)

	}

	// Return webhookID
	return webhookName
}
