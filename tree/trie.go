package tree

import (
	"errors"
	"git.woa.com/alanclzhu/go-web/context"
	"strings"
)

/**
 * @Description $
 * @Date 2024/9/23 09:51
 **/
type Tree struct {
	root *node
}
type node struct {
	isLast  bool
	segment string
	handler []context.ControllerHandler
	childs  []*node
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  make([]*node, 0),
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{
		root: root,
	}
}

// 判断一个segment是否是通用的segment,是否是以:开头
func IsWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	if IsWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, child := range n.childs {
		if IsWildSegment(child.segment) {
			nodes = append(nodes, child)
		} else if child.segment == segment {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(uri string) *node {
	//uri用"/"分割
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !IsWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	//如果只有一个segment，则是最后一个标记
	if len(segments) == 1 {
		for _, cn := range cnodes {
			if cn.isLast {
				return cn
			}
		}
		return nil
	}
	// 如果有2个segment, 递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}

	return nil
}

func (tree *Tree) AddRoute(uri string, handler []context.ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist:" + uri)
	}
	segments := strings.Split(uri, "/")

	for index, segment := range segments {
		if !IsWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1
		var objNode *node
		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}
		if objNode == nil {
			cnode := new(node)
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}
	return nil
}

func (tree *Tree) FindHandler(uri string) []context.ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}
