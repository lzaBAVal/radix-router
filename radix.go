package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
)

const (
	pathDelimiter = "/"
)

var (
	ErrMethodNotAllowed = errors.New("Method not allowed")
	ErrNotFound         = errors.New("Path not found")
)

type RadixTree struct {
	label    string
	actions  map[string]http.Handler
	children []*RadixTree
}

type action struct {
	handler http.Handler
}

type RadixTraverse struct {
	node *RadixTree
	tab  int
}

func NewRadixTraverse(n *RadixTree) *RadixTraverse {
	return &RadixTraverse{
		node: n,
		tab:  0,
	}
}

func newRadixTree() *RadixTree {
	return &RadixTree{
		label:    "/",
		actions:  make(map[string]http.Handler),
		children: make([]*RadixTree, 0),
	}
}

func (n *RadixTree) Insert(methods []string, path string, handler http.Handler) {
	labels := splitPath(path)

	node := n
	if len(labels) != 0 {
		node = n.insert(labels)
	}

	for _, method := range methods {
		node.actions[method] = handler
	}
}

func (n *RadixTree) insert(labels []string) *RadixTree {
	currNode := n

	if len(labels) == 0 {
		return currNode
	}

	for _, child := range currNode.children {
		generalPart := getPrefix(child.label, labels[0])
		lenGeneralPart := len(generalPart)

		if labels[0] == child.label {
			return child.insert(labels[1:])
		}

		if generalPart == child.label {
			labels[0] = labels[0][lenGeneralPart:]
			return child.insert(labels)
		}

		if lenGeneralPart > 0 {
			child.divideLabel(lenGeneralPart)

			labels[0] = labels[0][lenGeneralPart:]
			return child.insert(labels)

		}
	}

	currNode.children = append(currNode.children, &RadixTree{
		label:    labels[0],
		actions:  make(map[string]http.Handler),
		children: make([]*RadixTree, 0),
	})
	currNode = currNode.children[len(currNode.children)-1]
	return currNode.insert(labels[1:])
}

func splitPath(path string) []string {
	splited := strings.Split(path, pathDelimiter)[1:]

	if splited[len(splited)-1] == "" {
		splited = splited[:len(splited)-1]
	}

	labels := make([]string, 0)
	for _, label := range splited {
		trimmedLabel := strings.Trim(label, " ")
		if len(trimmedLabel) != 0 {
			labels = append(labels, trimmedLabel)
			continue
		}
		break
	}
	return labels
}

func getPrefix(l1, l2 string) string {
	minLength := minLength(l1, l2)

	i := 0
	for i = 0; i < minLength; i += 1 {
		if l1[i] != l2[i] {
			break
		}
	}

	return l1[:i]
}

func minLength(s1, s2 string) int {
	return int(math.Min(float64(len(s1)), float64(len(s2))))
}

func (n *RadixTree) divideLabel(lenGeneralPart int) {
	label := n.label
	children := n.children

	n.label = label[:lenGeneralPart]
	n.children = []*RadixTree{{
		label:    label[lenGeneralPart:],
		actions:  make(map[string]http.Handler),
		children: children,
	}}

}

func (n *RadixTree) Search(path, method string) (http.Handler, error) {
	labels := splitPath(path)
	node := n
	if len(labels) != 0 {
		node = n.search(labels)
	}
	if node == nil {
		return nil, ErrNotFound
	}
	if handler, ok := node.actions[method]; ok {
		return handler, nil
	}
	return nil, ErrMethodNotAllowed

}

func (n *RadixTree) search(labels []string) *RadixTree {
	if len(labels) == 0 {
		return nil
	}

	for _, child := range n.children {
		fmt.Println(child.label, labels)
		prefix := getPrefix(child.label, labels[0])

		if labels[0] == child.label {

			if len(labels) == 1 && labels[0] != "" {
				return child
			}
			return child.search(labels[1:])
		}
		if len(prefix) > 0 {
			labels[0] = labels[0][len(prefix):]
			return child.search(labels)
		}
	}
	return nil
}

func (rt *RadixTraverse) Traverse() {
	n := rt.node
	if n.label == "" && rt.tab != 0 || len(n.children) == 0 {
		return
	}

	for _, child := range n.children {
		if child.label == "" {
			continue
		}

		fmt.Println(strings.Repeat(" ", rt.tab) + child.label)

		rt.node = child
		rt.tab += 4
		rt.Traverse()
		rt.tab -= 4

	}
}
