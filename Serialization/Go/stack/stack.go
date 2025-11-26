package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"bufio"
	"os"
	"strconv"
)

const MAX_SIZE = 10

type SNode struct {
	key  string
	next *SNode
}

type Stack struct {
	head *SNode
	size int
}

func NewStack() *Stack {
	return &Stack{
	 head: nil,
	 size: 0,
	}
}

func NewStackFromSlice(items ...string) *Stack {
	s := NewStack()
	for _, item := range items {
		s.Push(item)
	}
	return s
}

func (s *Stack) Push(data string) error {
	if s.size >= MAX_SIZE {
	 return errors.New("stack overflow: maximum size reached")
	}
	
	newNode := &SNode{
	 key:  data,
	 next: s.head,
	}
	s.head = newNode
	s.size++
	return nil
}

func (s *Stack) Pop() (string, error) {
	if s.head == nil {
	 return "", errors.New("stack underflow: stack is empty")
	}
	
	data := s.head.key
	s.head = s.head.next
	s.size--
	return data, nil
}

func (s *Stack) IsEmpty() bool {
	return s.head == nil
}

func (s *Stack) GetSize() int {
	return s.size
}

func (s *Stack) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
	 return fmt.Errorf("не удалось открыть файл для записи: %s", filename)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, int32(s.size)); err != nil {
	 return err
	}

	current := s.head
	for current != nil {
	 keyBytes := []byte(current.key)
	 if err := binary.Write(file, binary.LittleEndian, int32(len(keyBytes))); err != nil {
	  return err
	 }
	 if _, err := file.Write(keyBytes); err != nil {
	  return err
	 }
	 current = current.next
	}
	
	return nil
}

func (s *Stack) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
	 return fmt.Errorf("не удалось открыть файл: %s", filename)
	}
	defer file.Close()

	s.Clear()

	var fileSize int32
	if err := binary.Read(file, binary.LittleEndian, &fileSize); err != nil {
	 return err
	}

	if fileSize > MAX_SIZE {
	 return errors.New("размер стека в файле превышает максимально допустимый")
	}

	tempArray := make([]string, fileSize)
	
	for i := int(fileSize) - 1; i >= 0; i-- {
	 var keyLength int32
	 if err := binary.Read(file, binary.LittleEndian, &keyLength); err != nil {
	  return errors.New("ошибка чтения длины строки из файла")
	 }
	 
	 keyBytes := make([]byte, keyLength)
	 if _, err := io.ReadFull(file, keyBytes); err != nil {
	  return errors.New("ошибка чтения строки из файла")
	 }
	 tempArray[i] = string(keyBytes)
	}

	for i := 0; i < int(fileSize); i++ {
	 if err := s.Push(tempArray[i]); err != nil {
	  return err
	 }
	}

	return nil
}

func (s *Stack) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
	 return fmt.Errorf("не удалось открыть файл для записи: %s", filename)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%d\n", s.size)); err != nil {
	 return err
	}

	current := s.head
	for current != nil {
	 if _, err := file.WriteString(current.key + "\n"); err != nil {
	  return err
	 }
	 current = current.next
	}
	
	return nil
}

func (s *Stack) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %s", filename)
	}
	defer file.Close()

	s.Clear()

	scanner := bufio.NewScanner(file)
	
	if !scanner.Scan() {
		return errors.New("не удалось прочитать размер стека")
	}
	
	sizeStr := scanner.Text()
	fileSize, err := strconv.Atoi(sizeStr)
	if err != nil {
		return err
	}

	if fileSize > MAX_SIZE {
		return errors.New("размер стека в файле превышает максимально допустимый")
	}

	tempArray := make([]string, fileSize)
	
	for i := fileSize - 1; i >= 0; i-- {
		if !scanner.Scan() {
			return errors.New("не удалось прочитать элемент стека")
		}
		tempArray[i] = scanner.Text()
	}

	for i := 0; i < fileSize; i++ {
		if err := s.Push(tempArray[i]); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (s *Stack) Print() {
	if s.IsEmpty() {
	 fmt.Println("Стек пуст")
	 return
	}
	
	current := s.head
	for current != nil {
	 fmt.Print(current.key, " ")
	 current = current.next
	}
	fmt.Println()
}

func (s *Stack) Clear() {
	for !s.IsEmpty() {
	 s.Pop()
	}
}

func main() {	
	a := NewStackFromSlice("HELP", "I CANT", "HOLD IT", "ANYMORE")
	a.Print()
	a.WriteBinary("test")
	a.WriteText("test2")

	b := NewStack()
	b.ReadBinary("test")
	b.Print()

	c := NewStack()
	c.ReadText("test2")
	c.Print()
}