package set

import (
	"sync"

	"github.com/linfangrong/carrier/models/command/command"
	"github.com/linfangrong/carrier/utils/util"
)

type CommandTree interface {
	AddCommand(name []byte, c command.Command)
	SearchCommand(name []byte) (command.Command, bool)
}

type commandNode struct {
	cmd       command.Command
	childNode map[byte]*commandNode
}

func newCommandNode() *commandNode {
	return &commandNode{
		cmd:       nil,
		childNode: make(map[byte]*commandNode),
	}
}

type commandTree struct {
	sync.RWMutex
	node *commandNode
}

func NewCommandTree() CommandTree {
	return &commandTree{
		node: newCommandNode(),
	}
}

func (ct *commandTree) AddCommand(name []byte, cmd command.Command) {
	ct.Lock()
	defer ct.Unlock()

	var (
		key         byte
		currentNode *commandNode
		nextNode    *commandNode
		found       bool

		currentNodes []*commandNode = []*commandNode{ct.node}
		nextNodes    []*commandNode = nil
	)
	for _, key = range name {
		for _, currentNode = range currentNodes {
			// 加大写
			if nextNode, found = currentNode.childNode[util.ToUpper(key)]; !found {
				nextNode = newCommandNode()
				currentNode.childNode[util.ToUpper(key)] = nextNode
			}
			nextNodes = append(nextNodes, nextNode)
			// 加小写
			if nextNode, found = currentNode.childNode[util.ToLower(key)]; !found {
				nextNode = newCommandNode()
				currentNode.childNode[util.ToLower(key)] = nextNode
			}
			nextNodes = append(nextNodes, nextNode)
		}
		currentNodes = nextNodes
		nextNodes = nil
	}
	// 写command
	for _, currentNode = range currentNodes {
		currentNode.cmd = cmd
	}
}

func (ct *commandTree) SearchCommand(name []byte) (c command.Command, ok bool) {
	ct.RLock()
	defer ct.RUnlock()

	var (
		key         byte
		currentNode *commandNode = ct.node
		nextNode    *commandNode
	)
	for _, key = range name {
		if nextNode, ok = currentNode.childNode[key]; !ok {
			return
		}
		currentNode = nextNode
	}
	return currentNode.cmd, currentNode.cmd != nil
}
