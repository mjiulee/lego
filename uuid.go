package lego

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

/*
 * 分布式-雪花uuid
 */
func UUID() int64 {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	id := node.Generate()
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
