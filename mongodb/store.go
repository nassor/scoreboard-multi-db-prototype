package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	mongoCollection *mongo.Collection

	updatedOpts *options.UpdateOptions
}

func NewStore(ctx context.Context, mc *mongo.Client) Store {
	t := true
	mongoCollection := mc.Database("gameops").Collection("leaderboards")
	mongoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "projectId", Value: 1}},
		Options: &options.IndexOptions{
			Unique:     &t,
			Background: &t,
		},
	})
	updatedOpts := options.Update().SetUpsert(true)
	return Store{
		mongoCollection: mongoCollection,
		updatedOpts:     updatedOpts,
	}
}

func (st *Store) Add(ctx context.Context, projectID, id, name string, score uint64, ts time.Time) error {
	_, err := st.mongoCollection.UpdateOne(
		ctx,
		bson.M{"projectId": projectID},
		bson.M{
			"$push": bson.M{
				"topScores": bson.M{
					"$each": bson.A{bson.D{
						{Key: "score", Value: score},
						// scoredAt can be use externally to "untie" same scores
						{Key: "scoredAt", Value: ts.UnixNano()},
						{Key: "playerId", Value: id},
						{Key: "playerName", Value: name},
					}},
					"$slice": 4,  // leaderboard max size
					"$sort":  -1, // first array member ("score") in descending order
				},
			},
			"$setOnInsert": bson.D{
				// can be used to decide when clean up the dashboard
				{Key: "createdAt", Value: ts.UnixNano()},
			},
		},
		st.updatedOpts,
	)
	if err != nil {
		return fmt.Errorf("when adding scores: %w", err)
	}

	return nil
}
