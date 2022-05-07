package binarysearchtree

import (
	"errors"
	"fmt"
	"time"
)

type BinaryNode struct {
	dentist interface{}
	patient interface{}
	date    string
	session int
	left    *BinaryNode
	right   *BinaryNode
}

type Binarysearchtree struct {
	root *BinaryNode
}

func New() *Binarysearchtree {
	bst := &Binarysearchtree{nil}
	return bst
}

func (bn *BinaryNode) GetDentist() interface{} {
	return bn.dentist
}

func (bn *BinaryNode) GetPatient() interface{} {
	return bn.patient
}

func (bn *BinaryNode) GetDate() string {
	return bn.date
}

func (bn *BinaryNode) GetSession() int {
	return bn.session
}

func (bn *BinaryNode) SetSession(session int) {
	bn.session = session
}

func (bn *BinaryNode) SetDentist(dentist interface{}) {
	bn.dentist = dentist
}

func (bst *Binarysearchtree) Add(date string, session int, dentist interface{}, patient interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic, recovered value: %v\n", r)
		}
	}()

	bst.insertNode(&bst.root, date, session, dentist, patient)
}

func (bst *Binarysearchtree) insertNode(t **BinaryNode, date string, session int, dentist interface{}, patient interface{}) {
	if (*t) == nil {
		newNode := &BinaryNode{dentist, patient, date, session, nil, nil}
		(*t) = newNode
	} else {
		if date < (*t).date {
			bst.insertNode(&(*t).left, date, session, dentist, patient) // dereferencing
		} else {
			bst.insertNode(&(*t).right, date, session, dentist, patient) // dereferencing
		}
	}
}

func (bst *Binarysearchtree) findSuccessor(t *BinaryNode) *BinaryNode {
	for t.right != nil { // Find node on extreme right
		t = t.right
	}
	return t
}

func (bst *Binarysearchtree) removeNode(t **BinaryNode, removeNode *BinaryNode) (*BinaryNode, error) {
	if *t == nil {
		return nil, errors.New("error: tree is empty")
	} else if removeNode.date < (*t).date {
		(*t).left, _ = bst.removeNode(&(*t).left, removeNode)
	} else if removeNode.date > (*t).date {
		(*t).right, _ = bst.removeNode(&(*t).right, removeNode)
	} else {
		if (*t).left == nil {
			return (*t).right, nil
		} else if (*t).right == nil {
			return (*t).left, nil
		} else { // 3rd case of 2 children
			(*t) = bst.findSuccessor((*t).left)
			(*t).left, _ = bst.removeNode(&(*t).left, removeNode)
		}
	}
	return *t, nil

}

func (bst *Binarysearchtree) Remove(removeNode *BinaryNode) error {
	bst.root, _ = bst.removeNode(&bst.root, removeNode)
	return nil
}

func (bst *Binarysearchtree) GetSize() int {
	return size(bst.root)
}

func size(node *BinaryNode) int {
	if node == nil {
		return 0
	} else {
		return (size(node.left) + 1 + size(node.right))
	}
}

func (bst *Binarysearchtree) Contains(date string, session int, dentist interface{}, patient interface{}) bool {
	newNode := BinaryNode{dentist, patient, date, session, nil, nil}
	return containsTraversal(bst.root, newNode)
}

func containsTraversal(t *BinaryNode, n BinaryNode) bool {
	if t != nil {
		containsTraversal(t.left, n)
		if t.patient == n.patient && t.dentist == n.dentist && t.date == n.date && t.session == n.session {
			return true
		}
		containsTraversal(t.right, n)
	}
	return false
}

func (bst *Binarysearchtree) GetSchedule(date string, searchInterface interface{}) []*BinaryNode {
	list := []*BinaryNode{}
	bst.searchSchedule(bst.root, date, searchInterface, &list)
	return list
}

func (bst *Binarysearchtree) searchSchedule(t *BinaryNode, date string, searchInterface interface{}, list *[]*BinaryNode) []*BinaryNode {
	if t == nil {
		return nil
	} else {
		if t.date == date {
			if t.right != nil {
				if fmt.Sprintf("%T", searchInterface) == "*main.Dentist" {
					if t.dentist == searchInterface {
						*list = append(*list, t)
					}
				}
				if fmt.Sprintf("%T", searchInterface) == "*main.Patient" {
					if t.patient == searchInterface {
						*list = append(*list, t)
					}
				}
				return bst.searchSchedule(t.right, date, searchInterface, list)
			} else {
				if fmt.Sprintf("%T", searchInterface) == "*main.Dentist" {
					if t.dentist == searchInterface {
						*list = append(*list, t)
					}
				}
				if fmt.Sprintf("%T", searchInterface) == "*main.Patient" {
					if t.patient == searchInterface {
						*list = append(*list, t)
					}
				}
				return *list
			}
		} else {
			if t.date > date {
				return bst.searchSchedule(t.left, date, searchInterface, list)
			} else {
				return bst.searchSchedule(t.right, date, searchInterface, list)
			}
		}
	}
}

func (bst *Binarysearchtree) GetAllSchedule(searchInterface interface{}) []*BinaryNode {
	list := []*BinaryNode{}
	oldDate := time.Now().AddDate(-100, 0, 0)
	bst.searchAllSchedule(bst.root, oldDate.Format("2006-01-02"), searchInterface, &list)
	return list
}

func (bst *Binarysearchtree) GetUpComingSchedule(searchInterface interface{}) []*BinaryNode {
	list := []*BinaryNode{}
	currentTime := time.Now()
	bst.searchAllSchedule(bst.root, currentTime.Format("2006-01-02"), searchInterface, &list)
	return list
}

func (bst *Binarysearchtree) searchAllSchedule(t *BinaryNode, date string, searchInterface interface{}, list *[]*BinaryNode) []*BinaryNode {
	if t != nil {
		bst.searchAllSchedule(t.left, date, searchInterface, list)
		if t.date >= date {
			if fmt.Sprintf("%T", searchInterface) == "*main.Dentist" {
				if t.dentist == searchInterface {
					*list = append(*list, t)
				}
			}
			if fmt.Sprintf("%T", searchInterface) == "*main.Patient" {
				if t.patient == searchInterface {
					*list = append(*list, t)
				}
			}
		}
		bst.searchAllSchedule(t.right, date, searchInterface, list)
	}
	return *list
}

func (bst *Binarysearchtree) GetScheduleByDate(date string) []*BinaryNode {
	var list []*BinaryNode
	bst.searchScheduleByDate(bst.root, date, &list)
	return list
}

func (bst *Binarysearchtree) searchScheduleByDate(t *BinaryNode, date string, list *[]*BinaryNode) []*BinaryNode {
	if t == nil {
		return nil
	} else {
		if t.date == date {
			if t.right != nil {
				*list = append(*list, t)
				return bst.searchScheduleByDate(t.right, date, list)
			} else {
				*list = append(*list, t)
				return *list
			}
		} else {
			if t.date > date {
				return bst.searchScheduleByDate(t.left, date, list)
			} else {
				return bst.searchScheduleByDate(t.right, date, list)
			}
		}
	}
}

func (bst *Binarysearchtree) GetAllAppointment() []*BinaryNode {
	list := []*BinaryNode{}
	bst.inOrderTraversal(bst.root, &list)
	return list
}

func (bst *Binarysearchtree) inOrderTraversal(t *BinaryNode, list *[]*BinaryNode) []*BinaryNode {
	if t != nil {
		bst.inOrderTraversal(t.left, list)
		*list = append(*list, t)
		bst.inOrderTraversal(t.right, list)
	}
	return *list
}

func (bst *Binarysearchtree) SearchAllNodeByField(field string, value interface{}, channel chan []*BinaryNode) {
	list := []*BinaryNode{}
	bst.searchInOrderTraversal(bst.root, field, value, &list)
	channel <- list
}

func (bst *Binarysearchtree) searchInOrderTraversal(t *BinaryNode, field string, value interface{}, list *[]*BinaryNode) []*BinaryNode {
	if t != nil {
		bst.searchInOrderTraversal(t.left, field, value, list)
		switch field {
		case "date":
			if t.GetDate() == value.(string) {
				*list = append(*list, t)
			}
		case "patient":
			if t.GetPatient() == value {
				*list = append(*list, t)
			}
		case "dentist":
			if t.GetDentist() == value {
				*list = append(*list, t)
			}
		case "session":
			if t.GetSession() == value.(int) {
				*list = append(*list, t)
			}
		}
		bst.searchInOrderTraversal(t.right, field, value, list)
	}
	return *list
}
