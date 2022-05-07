package doublylinkedlist

import (
	"errors"
	"reflect"

	util "GoSchool_Assignment2/utility"
)

type node struct {
	value    interface{}
	previous *node
	next     *node
}

type DoublyLinkedlist struct {
	head *node
	tail *node
	size int
}

// New will return a newly created instance of a doubly linked list
func New() *DoublyLinkedlist {
	list := &DoublyLinkedlist{nil, nil, 0}
	return list
}

// Get Size of linked list
func (list *DoublyLinkedlist) GetSize() int {
	return list.size
}

func (list *DoublyLinkedlist) Add(elm interface{}) error {
	newNode := &node{
		value:    elm,
		previous: nil,
		next:     nil,
	}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		currentNode := list.head
		previousNode := list.head
		for currentNode.next != nil {
			previousNode = currentNode
			currentNode = currentNode.next
		}
		currentNode.previous = previousNode
		currentNode.next = newNode
		list.tail = newNode
	}
	list.size++
	return nil
}

func (list *DoublyLinkedlist) GetList() []interface{} {
	var values []interface{}
	currentNode := list.head
	if currentNode == nil {
		return values
	}
	values = append(values, currentNode.value)
	for currentNode.next != nil {
		currentNode = currentNode.next
		values = append(values, currentNode.value)
	}
	return values
}

func (list *DoublyLinkedlist) Get(index int) (interface{}, error) {
	var value interface{}
	currentNode := list.head
	if currentNode == nil {
		return nil, errors.New("linked list is empty")
	}
	if index == 1 {
		value = currentNode.value
	} else {
		for i := 1; i <= index-1; i++ {
			currentNode = currentNode.next
		}
		value = currentNode.value
	}
	return value, nil
}

func (list *DoublyLinkedlist) Clear() {
	list.head = nil
	list.tail = nil
	list.size = 0
}

func getFieldValue(value interface{}, field string) interface{} {
	v := reflect.ValueOf(value)
	method := v.MethodByName(field)
	if method.IsValid() {
		return method.Call([]reflect.Value{})[0].Interface()
	}
	return nil
}

func (list *DoublyLinkedlist) SearchByMobileNumber(mobileNum int) interface{} {
	return list.recursiveBinarySearchByMobileNumber(list.head, list.tail, mobileNum, list.size)
}

func middleNode(start *node, mid int) *node {
	if start == nil {
		return nil
	}
	for i := 1; i < mid; i++ {
		start = start.next
	}
	return start
}

func (list *DoublyLinkedlist) recursiveBinarySearchByMobileNumber(firstNode *node, lastNode *node, value int, size int) interface{} {
	if firstNode == nil {
		return nil
	}
	firstNodeNum := getFieldValue(firstNode.value, "GetMobileNum").(int)
	lastNodeNum := getFieldValue(lastNode.value, "GetMobileNum").(int)
	if firstNodeNum > lastNodeNum {
		return nil
	} else {
		mid := size / 2
		midNode := middleNode(firstNode, mid)
		midNodeNum, _ := getFieldValue(midNode.value, "GetMobileNum").(int)
		if midNodeNum == value {
			return midNode.value
		} else {
			if value < midNodeNum {
				return list.recursiveBinarySearchByMobileNumber(firstNode, midNode.previous, value, mid)
			} else {
				return list.recursiveBinarySearchByMobileNumber(midNode.next, lastNode, value, mid)
			}
		}
	}
}

func (list *DoublyLinkedlist) SearchByName(value string) []interface{} {
	var searchResult []interface{}
	list.recursiveSeqSearchOfName(list.head, value, &searchResult)
	return searchResult
}

func (list *DoublyLinkedlist) recursiveSeqSearchOfName(node *node, value string, searchResult *[]interface{}) []interface{} {
	if node == nil {
		return nil
	} else {
		strNode := getFieldValue(node.value, "GetName").(string)
		if i := util.LevenshteinDistance(strNode, value); i <= 3 {
			*searchResult = append(*searchResult, node.value)
		}
		return list.recursiveSeqSearchOfName(node.next, value, searchResult)
	}
}
