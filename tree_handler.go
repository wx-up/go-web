package go_web

import (
	"errors"
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root *node
}

// 按照 / 分割 path 每一段是一个 node 节点
type node struct {
	segment string

	children []*node

	handler HandlerFunc
}

func (s *HandlerBasedOnTree) ServeHTTP(ctx *Context) {
	handlerFunc, err := s.Match(ctx.R.Method, ctx.R.URL.Path)
	if err != nil {
		ctx.W.WriteHeader(http.StatusNotFound)
		_, _ = ctx.W.Write([]byte("not found"))
		return
	}

	handlerFunc(ctx)

}

func (s *HandlerBasedOnTree) Route(method string, path string, handleFunc HandlerFunc) {
	path = strings.Trim(path, "/")
	if err := s.validatePattern(path); err != nil {
		panic(err)
	}
	s.root.addRoute(path, handleFunc)
}

func (s *HandlerBasedOnTree) Match(method string, path string) (HandlerFunc, error) {
	path = strings.Trim(path, "/")
	return s.root.matchRoute(path)
}

var ErrInvalidRouterPattern = errors.New("invalid router pattern")

func (s *HandlerBasedOnTree) validatePattern(path string) error {
	// 当存在 * 的时候，它必须是最后一个，并且前面必须是 /
	pos := strings.Index(path, "*")
	if pos < 0 {
		return nil
	}
	if pos != len(path)-1 {
		return ErrInvalidRouterPattern
	}

	if pos != 0 && path[pos-1] != '/' {
		return ErrInvalidRouterPattern
	}
	return nil
}

func (n *node) matchRoute(path string) (HandlerFunc, error) {
	segments := strings.Split(path, "/")
	current := n
	for _, segment := range segments {
		matchNode, err := current.matchChildNode(segment)
		if err != nil {
			return nil, err
		}
		current = matchNode
	}

	// 这个场景：比如注册了 /user/profile 但是访问了 /user
	if current.handler == nil {
		return nil, errors.New("路由不存在")
	}
	return current.handler, nil
}

func (n *node) addRoute(path string, handleFunc HandlerFunc) {
	segments := strings.Split(path, "/")

	currentNode := n
	for index, segment := range segments {
		matchNode, err := currentNode.matchChildNode(segment)
		if err != nil {
			// 创建子树
			currentNode.makeSubTree(segments[index:], handleFunc)
			return
		}
		currentNode = matchNode
	}

	currentNode.handler = handleFunc
}

func (n *node) matchChildNode(segment string) (*node, error) {
	var selectNode *node
	for _, child := range n.children {
		// 精确匹配
		if child.segment == segment && child.segment != "*" {
			return child, nil
		}

		// 如果 segment 为 * 则缓存当前 node，继续匹配，可能有精确的 segment 可以匹配
		if child.segment == "*" {
			selectNode = child
		}
	}

	// 说明只有通配符匹配上了
	if selectNode != nil {
		return selectNode, nil
	}

	return nil, errors.New("node not found")
}

// makeSubTree 创建子树
func (n *node) makeSubTree(segments []string, handleFunc HandlerFunc) {
	current := n
	for _, segment := range segments {
		newNode := &node{
			segment:  segment,
			children: make([]*node, 0, 2),
		}
		current.children = append(current.children, newNode)
		current = newNode
	}
	current.handler = handleFunc
}
