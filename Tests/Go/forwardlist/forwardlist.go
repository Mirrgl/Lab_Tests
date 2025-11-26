package forwardlist

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type node struct {
	key  string
	next *node
}

type ForwardList struct {
	head *node
	tail *node
	size int
}

func NewForwardList(items ...string) *ForwardList {
	fl := &ForwardList{}
	for _, item := range items {
		fl.PushBack(item)
	}
	return fl
}

func (fl *ForwardList) validatePosition(position int, allowEnd bool) error {
	maxPos := fl.size
	if !allowEnd {
		maxPos = fl.size - 1
	}
	if position < 0 || position > maxPos {
		return errors.New("index out of range")
	}
	return nil
}

func (fl *ForwardList) getNodeAt(position int) (*node, error) {
	if err := fl.validatePosition(position, false); err != nil {
		return nil, err
	}
	current := fl.head
	for i := 0; i < position; i++ {
		current = current.next
	}
	return current, nil
}

func (fl *ForwardList) PushBack(key string) {
	newNode := &node{key: key}
	if fl.head == nil {
		fl.head = newNode
		fl.tail = newNode
	} else {
		fl.tail.next = newNode
		fl.tail = newNode
	}
	fl.size++
}

func (fl *ForwardList) PushFront(key string) {
	newNode := &node{key: key, next: fl.head}
	if fl.head == nil {
		fl.tail = newNode
	}
	fl.head = newNode
	fl.size++
}

func (fl *ForwardList) InsertBefore(key string, position int) error {
	if position == 0 {
		fl.PushFront(key)
		return nil
	}
	if err := fl.validatePosition(position, false); err != nil {
		return err
	}
	prev, err := fl.getNodeAt(position - 1)
	if err != nil {
		return err
	}
	newNode := &node{key: key, next: prev.next}
	prev.next = newNode
	if prev == fl.tail {
		fl.tail = newNode
	}
	fl.size++
	return nil
}

func (fl *ForwardList) InsertAfter(key string, position int) error {
	if err := fl.validatePosition(position, false); err != nil {
		return err
	}
	if position == fl.size-1 {
		fl.PushBack(key)
		return nil
	}
	current, err := fl.getNodeAt(position)
	if err != nil {
		return err
	}
	newNode := &node{key: key, next: current.next}
	current.next = newNode
	fl.size++
	return nil
}

func (fl *ForwardList) PopFront() error {
	if fl.IsEmpty() {
		return errors.New("list is empty")
	}
	fl.head = fl.head.next
	if fl.head == nil {
		fl.tail = nil
	}
	fl.size--
	return nil
}

func (fl *ForwardList) PopBack() error {
	if fl.IsEmpty() {
		return errors.New("list is empty")
	}
	if fl.head == fl.tail {
		fl.head = nil
		fl.tail = nil
		fl.size = 0
		return nil
	}
	current := fl.head
	for current.next != fl.tail {
		current = current.next
	}
	current.next = nil
	fl.tail = current
	fl.size--
	return nil
}

func (fl *ForwardList) RemoveByValue(value string) bool {
	if fl.IsEmpty() {
		return false
	}
	if fl.head.key == value {
		fl.PopFront()
		return true
	}
	prev := fl.head
	for prev.next != nil && prev.next.key != value {
		prev = prev.next
	}
	if prev.next != nil {
		toDelete := prev.next
		prev.next = toDelete.next
		if toDelete == fl.tail {
			fl.tail = prev
		}
		fl.size--
		return true
	}
	return false
}

func (fl *ForwardList) Front() (string, error) {
	if fl.IsEmpty() {
		return "", errors.New("list is empty")
	}
	return fl.head.key, nil
}

func (fl *ForwardList) Back() (string, error) {
	if fl.IsEmpty() {
		return "", errors.New("list is empty")
	}
	return fl.tail.key, nil
}

func (fl *ForwardList) GetAt(index int) (string, error) {
	node, err := fl.getNodeAt(index)
	if err != nil {
		return "", err
	}
	return node.key, nil
}

func (fl *ForwardList) IsEmpty() bool {
	return fl.size == 0
}

func (fl *ForwardList) Size() int {
	return fl.size
}

func (fl *ForwardList) Clear() {
	fl.head = nil
	fl.tail = nil
	fl.size = 0
}

func (fl *ForwardList) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, uint64(fl.size)); err != nil {
		return fmt.Errorf("failed to write size: %w", err)
	}

	current := fl.head
	for current != nil {
		keyLength := uint64(len(current.key))
		if err := binary.Write(file, binary.LittleEndian, keyLength); err != nil {
			return fmt.Errorf("failed to write key length: %w", err)
		}
		if _, err := file.Write([]byte(current.key)); err != nil {
			return fmt.Errorf("failed to write key: %w", err)
		}
		current = current.next
	}
	return nil
}

func (fl *ForwardList) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fl.Clear()

	var size uint64
	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {
		return fmt.Errorf("failed to read size: %w", err)
	}

	for i := uint64(0); i < size; i++ {
		var keyLength uint64
		if err := binary.Read(file, binary.LittleEndian, &keyLength); err != nil {
			return fmt.Errorf("failed to read key length: %w", err)
		}
		keyBytes := make([]byte, keyLength)
		if _, err := file.Read(keyBytes); err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}
		fl.PushBack(string(keyBytes))
	}
	return nil
}

func (fl *ForwardList) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%d\n", fl.size); err != nil {
		return fmt.Errorf("failed to write size: %w", err)
	}

	current := fl.head
	for current != nil {
		if _, err := fmt.Fprintln(file, current.key); err != nil {
			return fmt.Errorf("failed to write element: %w", err)
		}
		current = current.next
	}
	return nil
}

func (fl *ForwardList) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fl.Clear()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return errors.New("file is empty")
	}
	sizeStr := scanner.Text()
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid size format: %w", err)
	}

	for i := 0; i < size; i++ {
		if !scanner.Scan() {
			return errors.New("unexpected end of file")
		}
		fl.PushBack(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

func (fl *ForwardList) Print() {
	if fl.IsEmpty() {
		fmt.Println("List is empty")
		return
	}
	current := fl.head
	for current != nil {
		fmt.Print(current.key, " ")
		current = current.next
	}
	fmt.Println()
}