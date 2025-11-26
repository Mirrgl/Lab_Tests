package array

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Array struct {
	data []string 
	len  int      
	cap  int      
}

func NewArray(size int) (*Array, error) {
	if size < 1 {
		return nil, errors.New("cannot create array of zero size")
	}
	return &Array{
		data: make([]string, size),
		len:  0,
		cap:  size,
	}, nil
}

func NewArrayFromList(items []string) (*Array, error) {
	if len(items) < 1 {
		return NewArray(1)
	}
	a, _ := NewArray(len(items))
	for _, item := range items {
		a.AddElementEnd(item)
	}
	return a, nil
}

func (a *Array) grow() {
	newCap := a.cap * 2
	newData := make([]string, newCap)
	copy(newData, a.data[:a.len])
	a.data = newData
	a.cap = newCap
}

func (a *Array) GetElement(index int) (string, error) {
	if index < 0 || index >= a.len {
		return "", errors.New("index out of bounds")
	}
	return a.data[index], nil
}

func (a *Array) SetElement(key string, index int) error {
	if index < 0 || index >= a.len {
		return errors.New("index out of bounds")
	}
	a.data[index] = key
	return nil
}

func (a *Array) DeleteElement(index int) error {
	if index < 0 || index >= a.len {
		return errors.New("index out of bounds")
	}
	for i := index; i < a.len-1; i++ {
		a.data[i] = a.data[i+1]
	}
	a.len--
	return nil
}

func (a *Array) AddElementAtIndex(key string, index int) error {
	if index < 0 || index > a.len {
		return errors.New("index out of bounds")
	}
	if a.len >= a.cap {
		a.grow()
	}
	for i := a.len; i > index; i-- {
		a.data[i] = a.data[i-1]
	}
	a.data[index] = key
	a.len++
	return nil
}

func (a *Array) AddElementEnd(key string) {
	if a.len >= a.cap {
		a.grow()
	}
	a.data[a.len] = key
	a.len++
}

func (a *Array) GetLength() int {
	return a.len
}

func (a *Array) GetCapacity() int {
	return a.cap
}

func (a *Array) IsInArray(key string) int {
	for i := 0; i < a.len; i++ {
		if a.data[i] == key {
			return i
		}
	}
	return -1
}

func (a *Array) Print() {
	for i := 0; i < a.len; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(a.data[i])
	}
	fmt.Println()
}

func (a *Array) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, uint32(a.len)); err != nil {
		return err
	}

	for i := 0; i < a.len; i++ {
		s := a.data[i]
		length := uint32(len(s))
		if err := binary.Write(file, binary.LittleEndian, length); err != nil {
			return err
		}
		if _, err := file.Write([]byte(s)); err != nil {
			return err
		}
	}
	return nil
}

func (a *Array) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var length uint32
	if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
		return err
	}

	newCap := int(length)
	if newCap == 0 {
		newCap = 1
	}
	a.data = make([]string, newCap)
	a.cap = newCap
	a.len = 0

	for i := uint32(0); i < length; i++ {
		var strLen uint32
		if err := binary.Read(file, binary.LittleEndian, &strLen); err != nil {
			return err
		}
		buf := make([]byte, strLen)
		if _, err := file.Read(buf); err != nil {
			return err
		}
		a.data[i] = string(buf)
		a.len++
	}
	return nil
}

func (a *Array) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	if _, err := fmt.Fprintf(writer, "%d\n", a.len); err != nil {
		return err
	}

	for i := 0; i < a.len; i++ {
		if _, err := fmt.Fprintln(writer, a.data[i]); err != nil {
			return err
		}
	}
	return nil
}

func (a *Array) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return errors.New("empty file")
	}
	lengthStr := scanner.Text()
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return fmt.Errorf("invalid length line: %w", err)
	}

	newCap := length
	if newCap == 0 {
		newCap = 1
	}
	a.data = make([]string, newCap)
	a.cap = newCap
	a.len = 0

	for i := 0; i < length; i++ {
		if !scanner.Scan() {
			return errors.New("unexpected EOF")
		}
		a.data[i] = scanner.Text()
		a.len++
	}
	return nil
}