package aerospike

import (
	"fmt"

	"github.com/aerospike/aerospike-client-go/v5"
)

type Store struct {
	asc *aerospike.Client
	wp  *aerospike.WritePolicy
	lp  *aerospike.ListPolicy
}

func NewStore(asc *aerospike.Client) Store {
	return Store{
		asc: asc,
		wp:  aerospike.NewWritePolicy(0, 0),
		lp: aerospike.NewListPolicy(
			aerospike.ListOrderOrdered,
			aerospike.ListWriteFlagsDefault,
		),
	}
}

func (st *Store) Add(projectid, id, name string, score uint64) error {
	key, err := aerospike.NewKey("test", "test", fmt.Sprintf("%s:%s", projectid, "top-scores"))
	if err != nil {
		return fmt.Errorf("when generating keys: %w", err)
	}

	if _, err := st.asc.Operate(
		st.wp,
		key,
		aerospike.ListAppendWithPolicyOp(st.lp, "scores", createScore(name, score))); err != nil {
		return fmt.Errorf("when appending the score: %w", err)
	}

	return nil
}

func createScore(name string, score uint64) string {
	return fmt.Sprintf("%d:%s", score, name)
}
