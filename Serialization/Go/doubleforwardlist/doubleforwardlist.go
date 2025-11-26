package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DFNode struct {
	key  string
	next *DFNode
	prev *DFNode
}

type DoubleList struct {
	head   *DFNode
	tail   *DFNode
	length int
}

func NewDoubleList(items ...string) *DoubleList {
	dl := &DoubleList{}
	for _, item := range items {
		dl.AddTail(item)
	}
	return dl
}

func (dl *DoubleList) validateIndex(index int, allowEnd bool) error {
	maxIndex := dl.length
	if !allowEnd {
		maxIndex--
	}
	if index < 0 || index > maxIndex {
		return errors.New("Индекс больше возможного")
	}
	return nil
}

func (dl *DoubleList) getNodeAt(index int) (*DFNode, error) {
	if err := dl.validateIndex(index, false); err != nil {
		return nil, err
	}

	if index <= dl.length/2 {
		current := dl.head
		for i := 0; i < index; i++ {
			current = current.next
		}
		return current, nil
	} else {
		current := dl.tail
		for i := dl.length - 1; i > index; i-- {
			current = current.prev
		}
		return current, nil
	}
}

func (dl *DoubleList) AddAfter(key string, index int) error {
	if err := dl.validateIndex(index, false); err != nil {
		return err
	}

	current, err := dl.getNodeAt(index)
	if err != nil {
		return err
	}

	newNode := &DFNode{key: key}
	newNode.next = current.next
	newNode.prev = current

	if current.next != nil {
		current.next.prev = newNode
	} else {
		dl.tail = newNode
	}
	current.next = newNode

	dl.length++
	return nil
}

func (dl *DoubleList) AddBefore(key string, index int) error {
	if index == 0 {
		return dl.AddHead(key)
	}
	return dl.AddAfter(key, index-1)
}

func (dl *DoubleList) AddHead(key string) error {
	newNode := &DFNode{key: key}
	newNode.next = dl.head

	if dl.head != nil {
		dl.head.prev = newNode
	} else {
		dl.tail = newNode
	}
	dl.head = newNode
	dl.length++
	return nil
}

func (dl *DoubleList) AddTail(key string) error {
	newNode := &DFNode{key: key}
	newNode.prev = dl.tail

	if dl.tail != nil {
		dl.tail.next = newNode
	} else {
		dl.head = newNode
	}
	dl.tail = newNode
	dl.length++
	return nil
}

func (dl *DoubleList) DeleteAt(index int) error {
	if err := dl.validateIndex(index, false); err != nil {
		return err
	}

	if index == 0 {
		return dl.DeleteHead()
	}
	if index == dl.length-1 {
		return dl.DeleteTail()
	}

	toDelete, err := dl.getNodeAt(index)
	if err != nil {
		return err
	}

	toDelete.prev.next = toDelete.next
	toDelete.next.prev = toDelete.prev
	dl.length--
	return nil
}

func (dl *DoubleList) DeleteHead() error {
	if dl.head == nil {
		return nil
	}

	toDelete := dl.head
	dl.head = toDelete.next

	if dl.head != nil {
		dl.head.prev = nil
	} else {
		dl.tail = nil
	}
	dl.length--
	return nil
}

func (dl *DoubleList) DeleteTail() error {
	if dl.tail == nil {
		return nil
	}

	toDelete := dl.tail
	dl.tail = toDelete.prev

	if dl.tail != nil {
		dl.tail.next = nil
	} else {
		dl.head = nil
	}
	dl.length--
	return nil
}

func (dl *DoubleList) DeleteByValue(key string) error {
	current := dl.head
	index := 0

	for current != nil {
		if current.key == key {
			return dl.DeleteAt(index)
		}
		current = current.next
		index++
	}
	return errors.New("Ключ не найден")
}

func (dl *DoubleList) GetElement(index int) (string, error) {
	node, err := dl.getNodeAt(index)
	if err != nil {
		return "", err
	}
	return node.key, nil
}

func (dl *DoubleList) PopElement(index int) (string, error) {
	value, err := dl.GetElement(index)
	if err != nil {
		return "", err
	}
	if err := dl.DeleteAt(index); err != nil {
		return "", err
	}
	return value, nil
}

func (dl *DoubleList) FindByValue(key string) *DFNode {
	current := dl.head
	for current != nil {
		if current.key == key {
			return current
		}
		current = current.next
	}
	return nil
}

func (dl *DoubleList) IsEmpty() bool {
	return dl.length == 0
}

func (dl *DoubleList) GetLength() int {
	return dl.length
}

func (dl *DoubleList) clear() {
	dl.head = nil
	dl.tail = nil
	dl.length = 0
}

func (dl *DoubleList) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл для записи: %w", err)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, uint64(dl.length)); err != nil {
		return err
	}

	current := dl.head
	for current != nil {
		keyLength := uint64(len(current.key))
		if err := binary.Write(file, binary.LittleEndian, keyLength); err != nil {
			return err
		}
		if _, err := file.Write([]byte(current.key)); err != nil {
			return err
		}
		current = current.next
	}
	return nil
}

func (dl *DoubleList) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл: %w", err)
	}
	defer file.Close()

	dl.clear()

	var newLength uint64
	if err := binary.Read(file, binary.LittleEndian, &newLength); err != nil {
		return err
	}

	for i := uint64(0); i < newLength; i++ {
		var keyLength uint64
		if err := binary.Read(file, binary.LittleEndian, &keyLength); err != nil {
			return err
		}

		keyBytes := make([]byte, keyLength)
		if _, err := io.ReadFull(file, keyBytes); err != nil {
			return err
		}

		if err := dl.AddTail(string(keyBytes)); err != nil {
			return err
		}
	}
	return nil
}

func (dl *DoubleList) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл для записи: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	if _, err := fmt.Fprintf(writer, "%d\n", dl.length); err != nil {
		return err
	}

	current := dl.head
	for current != nil {
		if _, err := fmt.Fprintln(writer, current.key); err != nil {
			return err
		}
		current = current.next
	}
	return nil
}

func (dl *DoubleList) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл: %w", err)
	}
	defer file.Close()

	dl.clear()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return io.EOF
	}

	lengthStr := strings.TrimSpace(scanner.Text())
	newLength, err := strconv.Atoi(lengthStr)
	if err != nil {
		return fmt.Errorf("invalid length line: %w", err)
	}

	for i := 0; i < newLength; i++ {
		if !scanner.Scan() {
			return errors.New("unexpected EOF in file")
		}
		key := strings.TrimSpace(scanner.Text())
		if key != "" || dl.length == 0 {
			if err := dl.AddTail(key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (dl *DoubleList) Print() {
	if dl.IsEmpty() {
		fmt.Println("Список пуст")
		return
	}

	current := dl.head
	for current != nil {
		fmt.Print(current.key)
		if current.next != nil {
			fmt.Print(" ")
		}
		current = current.next
	}
	fmt.Println()
}

func main() {
	a := NewDoubleList("HELP", "I CANT", "HOLD IT", "ANYMORE")
	a.Print()

	a.WriteBinary("test")
	a.WriteText("test2")

	b := NewDoubleList()
	b.ReadBinary("test")
	b.Print()

	d := NewDoubleList()
	d.ReadText("test2")
	d.Print()
}