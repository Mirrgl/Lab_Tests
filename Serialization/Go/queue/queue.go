package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const MAX_SIZE = 1000

type Node struct {
	Data string
	Next *Node
	Prev *Node
}

type Queue struct {
	head   *Node
	tail   *Node
	size   int
	maxSize int
}

func NewQueue() *Queue {
	return &Queue{
		head:    nil,
		tail:    nil,
		size:    0,
		maxSize: MAX_SIZE,
	}
}

func NewQueueWithItems(items ...string) *Queue {
	q := NewQueue()
	for _, item := range items {
		q.Enqueue(item)
	}
	return q
}

func (q *Queue) Enqueue(value string) error {
	if q.size >= q.maxSize {
		return errors.New("queue overflow")
	}

	newNode := &Node{
		Data: value,
		Next: nil,
		Prev: q.tail,
	}

	if q.size == 0 {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.Next = newNode
		q.tail = newNode
	}
	q.size++
	return nil
}

func (q *Queue) Dequeue() (string, error) {
	if q.size == 0 {
		return "", errors.New("queue underflow")
	}

	data := q.head.Data
	if q.size == 1 {
		q.head = nil
		q.tail = nil
	} else {
		newHead := q.head.Next
		newHead.Prev = nil
		q.head = newHead
	}
	q.size--
	return data, nil
}

func (q *Queue) Del(key string) {
	current := q.head
	for current != nil && current.Data != key {
		current = current.Next
	}

	if current == nil {
		return
	}

	if current.Prev != nil {
		current.Prev.Next = current.Next
	} else {
		q.head = current.Next
	}

	if current.Next != nil {
		current.Next.Prev = current.Prev
	} else {
		q.tail = current.Prev
	}

	q.size--
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) Head() *Node {
	return q.head
}

func (q *Queue) Clear() {
	q.head = nil
	q.tail = nil
	q.size = 0
}

func (q *Queue) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, uint64(q.size)); err != nil {
		return fmt.Errorf("failed to write size: %w", err)
	}

	current := q.head
	for current != nil {
		key := current.Data
		keyLength := uint64(len(key))
		if err := binary.Write(file, binary.LittleEndian, keyLength); err != nil {
			return fmt.Errorf("failed to write key length: %w", err)
		}
		if _, err := file.Write([]byte(key)); err != nil {
			return fmt.Errorf("failed to write key: %w", err)
		}
		current = current.Next
	}

	return nil
}

func (q *Queue) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	q.Clear()

	var size uint64
	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {
		return fmt.Errorf("failed to read size: %w", err)
	}

	if size > uint64(q.maxSize) {
		return fmt.Errorf("queue size in file exceeds maximum size %d", q.maxSize)
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

		if err := q.Enqueue(string(keyBytes)); err != nil {
			return err
		}
	}

	return nil
}

func (q *Queue) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%d\n", q.size); err != nil {
		return fmt.Errorf("failed to write size: %w", err)
	}

	current := q.head
	for current != nil {
		if _, err := fmt.Fprintln(file, current.Data); err != nil {
			return fmt.Errorf("failed to write element: %w", err)
		}
		current = current.Next
	}

	return nil
}

func (q *Queue) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	q.Clear()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return errors.New("file is empty")
	}

	sizeStr := scanner.Text()
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid size format: %w", err)
	}

	if size > q.maxSize {
		return fmt.Errorf("queue size in file exceeds maximum size %d", q.maxSize)
	}

	for i := 0; i < size; i++ {
		if !scanner.Scan() {
			return errors.New("unexpected end of file")
		}
		value := scanner.Text()
		if err := q.Enqueue(value); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (q *Queue) Print() {
	if q.size == 0 {
		fmt.Println("Queue is empty")
		return
	}

	current := q.head
	for current != nil {
		fmt.Print(current.Data, " ")
		current = current.Next
	}
	fmt.Println()
}

func main() {
	a := NewQueueWithItems("HELP", "I CANT", "HOLD IT", "ANYMORE")
	a.Print()
	a.WriteBinary("test")
	a.WriteText("test2")

	b := NewQueue()
	b.ReadBinary("test")
	b.Print()

	c := NewQueue()
	c.ReadText("test2")
	c.Print()
}