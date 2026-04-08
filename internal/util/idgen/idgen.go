// Package idgen generates time-ordered, non-colliding int64 IDs using Twitter
// Snowflake algorithm (via bwmarrin/snowflake).
//
// IDs embed a millisecond timestamp, making them sortable by creation time and
// collision-free within a single node without any coordination.
package idgen

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic("idgen: failed to create snowflake node: " + err.Error())
	}
}

// New returns a time-ordered, non-colliding positive int64 ID.
func New() int64 {
	return node.Generate().Int64()
}
