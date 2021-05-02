package lego

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"sync"
)

var _snowflakeNodeOnce sync.Once
var _node *snowflake.Node

func nodeInstance() *snowflake.Node {
	_snowflakeNodeOnce.Do(func() {
		anode, err := snowflake.NewNode(1)
		if err != nil {
			panic(err)
		}
		_node = anode
	})
	return _node
}

/*
 * 分布式-雪花uuid
 */
func UUID() int64 {
	id := nodeInstance().Generate()
	return id.Int64()
}

/*
 * 分布式-雪花uuid
 */
func UUIDWithNoteId(nid int64) int64 {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(nid)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	id := node.Generate()
	return id.Int64()
}
