package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

const (
	pathDelimiter = "/"
)

type node struct {
	label    string
	actions  map[string]*action
	children []*node
}

type action struct {
	handler http.Handler
}

type RadixTraverse struct {
	node *node
	tab  int
}

func NewRadixTraverse(n *node) *RadixTraverse {
	return &RadixTraverse{
		node: n,
		tab:  0,
	}
}

func newNode() *node {
	return &node{
		label:    "",
		actions:  make(map[string]*action),
		children: make([]*node, 0),
	}
}

func (n *node) Insert(path string) {
	labels := splitPath(path)
	n.insert(labels)
}

func (n *node) insert(labels []string) {
	currNode := n

	if len(labels) == 0 {
		return
	}

	for _, child := range currNode.children {
		generalPart := getPrefix(child.label, labels[0])
		lenGeneralPart := len(generalPart)

		if labels[0] == child.label {
			child.insert(labels[1:])
			return
		}

		if generalPart == child.label {
			labels[0] = labels[0][lenGeneralPart:]
			child.insert(labels)
			return
		}

		if lenGeneralPart > 0 {
			child.divideLabel(lenGeneralPart)

			labels[0] = labels[0][lenGeneralPart:]
			child.insert(labels)
			return
		}
	}

	currNode.children = append(currNode.children, &node{
		label:    labels[0],
		actions:  make(map[string]*action),
		children: make([]*node, 0),
	})
	currNode = currNode.children[len(currNode.children)-1]
	currNode.insert(labels[1:])
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

func (n *node) divideLabel(lenGeneralPart int) {
	label := n.label
	children := n.children

	n.label = label[:lenGeneralPart]
	n.children = []*node{{
		label:    label[lenGeneralPart:],
		actions:  make(map[string]*action),
		children: children,
	}}

}

func (n *node) Search(path string) bool {
	labels := splitPath(path)
	return n.search(labels)
}

func (n *node) search(labels []string) bool {
	for _, child := range n.children {
		prefix := getPrefix(child.label, labels[0])

		if prefix == labels[0] {
			if len(labels) == 1 {
				return true
			}
			return child.search(labels[1:])
		}
		if len(prefix) > 0 {
			labels[0] = labels[0][len(prefix):]
			return child.search(labels)
		}
	}
	return false
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
