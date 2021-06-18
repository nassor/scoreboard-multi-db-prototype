package aerospike

import (
	"testing"

	"github.com/aerospike/aerospike-client-go/v5"
	"github.com/stretchr/testify/assert"
)

func loadDatabase(t *testing.T) *aerospike.Client {
	asc, err := aerospike.NewClient("127.0.0.1", 3000)
	if err != nil {
		t.FailNow()
	}
	rs, err := asc.ScanAll(asc.DefaultScanPolicy, "test", "test", "scores")
	assert.NoError(t, err)
	for res := range rs.Results() {
		asc.Delete(aerospike.NewWritePolicy(0, 0), res.Record.Key)
	}
	return asc
}

func TestStore_Add(t *testing.T) {
	asc := loadDatabase(t)
	st := NewStore(asc)
	assert.NoError(t, st.Add("p1", "id1", "name1", 100203))
	assert.NoError(t, st.Add("p1", "id2", "name2", 10003))
	assert.NoError(t, st.Add("p1", "id3", "name3", 400203))

	rs, err := asc.ScanAll(asc.DefaultScanPolicy, "test", "test", "scores")
	assert.NoError(t, err)
	for res := range rs.Results() {
		t.Logf("\nResults: %+v (the order is asc)", res.Record.Bins)
	}
}
