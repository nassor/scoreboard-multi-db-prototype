package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadDatabase(t *testing.T) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"), &options.ClientOptions{Auth: &options.Credential{
		Username: "test",
		Password: "test",
	}})
	if err != nil {
		t.FailNow()
	}
	mc.Database("gameops").Drop(ctx)
	return mc
}

func TestStore_Add(t *testing.T) {
	ctx := context.Background()
	mc := loadDatabase(t)
	st := NewStore(ctx, mc)

	assert.NoError(t, st.Add(ctx, "p1", "id1", "name1", 100203, time.Now()))
	assert.NoError(t, st.Add(ctx, "p1", "id2", "name2", 10003, time.Now()))
	assert.NoError(t, st.Add(ctx, "p1", "id3", "name3", 400203, time.Now()))
	assert.NoError(t, st.Add(ctx, "p1", "id4", "name4", 40203, time.Now()))
	assert.NoError(t, st.Add(ctx, "p1", "id5", "name5", 874003, time.Now()))
	assert.NoError(t, st.Add(ctx, "p1", "id6", "name6", 574003, time.Now()))

	var something interface{}
	err := st.mongoCollection.FindOne(ctx, bson.M{
		"projectId": "p1",
	}).Decode(&something)

	if err != mongo.ErrNoDocuments {
		assert.NoError(t, err)
	}
	t.Logf("\nResults: %+v (the order is asc)", something)

}
