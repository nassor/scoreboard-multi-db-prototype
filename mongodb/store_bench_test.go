package mongodb

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BenchmarkStore_Add(b *testing.B) {
	ctx := context.Background()
	mc := loadDatabaseBench(b)
	st := NewStore(ctx, mc)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		st.Add(ctx, "p1", "id1", "name1", 100000+uint64(n), time.Now())
	}
}

func loadDatabaseBench(b *testing.B) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"), &options.ClientOptions{Auth: &options.Credential{
		Username: "test",
		Password: "test",
	}})
	if err != nil {
		b.FailNow()
	}
	mc.Database("gameops").Drop(ctx)
	return mc
}
