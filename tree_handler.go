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
	s.root.addRoute(path, handleFunc)
}

func (s *HandlerBasedOnTree) Match(method string, path string) (HandlerFunc, error) {
	return s.root.matchRoute(path)
}

func (n *node) matchRoute(path string) (HandlerFunc, error) {
	path = strings.Trim(path, "/")
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
	path = strings.Trim(path, "/")
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
	for _, child := range n.children {
		if child.segment == segment {
			return child, nil
		}
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
