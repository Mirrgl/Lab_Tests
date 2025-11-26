package hashmap

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strconv"
	"strings"
)

type ChainNode struct {
	Key  string
	Data int
	Next *ChainNode
}

func NewChainNode(key string, data int) *ChainNode {
	return &ChainNode{
		Key:  key,
		Data: data,
		Next: nil,
	}
}

type Bucket struct {
	Head *ChainNode
}

func NewBucket() *Bucket {
	return &Bucket{Head: nil}
}

type ChainMap struct {
	table    []*Bucket
	capacity int
	size     int
}

func NewChainMap(initialCapacity int) *ChainMap {
	table := make([]*Bucket, initialCapacity)
	for i := range table {
		table[i] = NewBucket()
	}
	return &ChainMap{
		table:    table,
		capacity: initialCapacity,
		size:     0,
	}
}

func (cm *ChainMap) hashFunction(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % cm.capacity
}

func (cm *ChainMap) rehash() {
	newCapacity := cm.capacity * 2
	newTable := make([]*Bucket, newCapacity)
	for i := range newTable {
		newTable[i] = NewBucket()
	}

	for i := 0; i < cm.capacity; i++ {
		currentNode := cm.table[i].Head
		for currentNode != nil {
			nextNode := currentNode.Next

			h := fnv.New32a()
			h.Write([]byte(currentNode.Key))
			newIndex := int(h.Sum32()) % newCapacity

			currentNode.Next = newTable[newIndex].Head
			newTable[newIndex].Head = currentNode

			currentNode = nextNode
		}
	}

	cm.table = newTable
	cm.capacity = newCapacity
}

func (cm *ChainMap) Add(key string, data int) {
	if float64(cm.size) >= float64(cm.capacity)*0.75 {
		cm.rehash()
	}

	index := cm.hashFunction(key)
	currentNode := cm.table[index].Head

	for currentNode != nil {
		if currentNode.Key == key {
			currentNode.Data = data
			return
		}
		currentNode = currentNode.Next
	}

	newNode := NewChainNode(key, data)
	newNode.Next = cm.table[index].Head
	cm.table[index].Head = newNode
	cm.size++
}

func (cm *ChainMap) Del(key string) {
	index := cm.hashFunction(key)
	currentNode := cm.table[index].Head
	var prevNode *ChainNode

	for currentNode != nil {
		if currentNode.Key == key {
			if prevNode == nil {
				cm.table[index].Head = currentNode.Next
			} else {
				prevNode.Next = currentNode.Next
			}
			cm.size--
			return
		}
		prevNode = currentNode
		currentNode = currentNode.Next
	}
}

func (cm *ChainMap) IsContain(key string) bool {
	index := cm.hashFunction(key)
	current := cm.table[index].Head

	for current != nil {
		if current.Key == key {
			return true
		}
		current = current.Next
	}

	return false
}

func (cm *ChainMap) Find(key string) (int, error) {
	index := cm.hashFunction(key)
	current := cm.table[index].Head

	for current != nil {
		if current.Key == key {
			return current.Data, nil
		}
		current = current.Next
	}

	return 0, fmt.Errorf("в словаре нет такого ключа")
}

func (cm *ChainMap) GetAllKeys(result *ChainMap) {
	for i := 0; i < cm.capacity; i++ {
		currentNode := cm.table[i].Head
		for currentNode != nil {
			result.Add(currentNode.Key, 1)
			currentNode = currentNode.Next
		}
	}
}

func (cm *ChainMap) GetAllKeysAsString() string {
	var result strings.Builder
	tempKeys := NewChainMap(cm.capacity)
	cm.GetAllKeys(tempKeys)

	for i := 0; i < tempKeys.capacity; i++ {
		currentNode := tempKeys.table[i].Head
		for currentNode != nil {
			result.WriteString(currentNode.Key)
			currentNode = currentNode.Next
		}
	}

	return result.String()
}

func (cm *ChainMap) PrintContents() {
	fmt.Println("Содержимое хеш-таблицы:")
	for i := 0; i < cm.capacity; i++ {
		fmt.Printf("[%d]: ", i)
		currentNode := cm.table[i].Head
		if currentNode != nil {
			for currentNode.Next != nil {
				fmt.Printf("%s -> %d, ", currentNode.Key, currentNode.Data)
				currentNode = currentNode.Next
			}
			fmt.Printf("%s -> %d", currentNode.Key, currentNode.Data)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (cm *ChainMap) appendNode(index int, node *ChainNode) {
    bucket := cm.table[index]
    if bucket.Head == nil {
        bucket.Head = node
        return
    }
    cur := bucket.Head
    for cur.Next != nil {
        cur = cur.Next
    }
    cur.Next = node
}

func (cm *ChainMap) WriteBinary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл для записи: %s", filename)
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, int64(cm.capacity)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, int64(cm.size)); err != nil {
		return err
	}

	for i := 0; i < cm.capacity; i++ {
		currentNode := cm.table[i].Head
		for currentNode != nil {
			keyLength := int64(len(currentNode.Key))
			if err := binary.Write(file, binary.LittleEndian, keyLength); err != nil {
				return err
			}
			if _, err := file.Write([]byte(currentNode.Key)); err != nil {
				return err
			}

			if err := binary.Write(file, binary.LittleEndian, int32(currentNode.Data)); err != nil {
				return err
			}

			currentNode = currentNode.Next
		}
	}

	return nil
}

func (cm *ChainMap) ReadBinary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %s", filename)
	}
	defer file.Close()

	cm.table = nil

	var capacity, size int64
	if err := binary.Read(file, binary.LittleEndian, &capacity); err != nil {
		return err
	}
	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {
		return err
	}

	cm.capacity = int(capacity)
	cm.size = int(size)

	cm.table = make([]*Bucket, cm.capacity)
	for i := range cm.table {
		cm.table[i] = NewBucket()
	}

	for i := 0; i < int(size); i++ {
		var keyLength int64
		if err := binary.Read(file, binary.LittleEndian, &keyLength); err != nil {
			return err
		}

		keyBytes := make([]byte, keyLength)
		if _, err := io.ReadFull(file, keyBytes); err != nil {
			return err
		}
		key := string(keyBytes)

		var data int32
		if err := binary.Read(file, binary.LittleEndian, &data); err != nil {
			return err
		}

		index := cm.hashFunction(key)
		newNode := NewChainNode(key, int(data))
		cm.appendNode(index, newNode)
	}

	return nil
}

func (cm *ChainMap) WriteText(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл для записи: %s", filename)
	}
	defer file.Close()

	fmt.Fprintf(file, "%d %d\n", cm.capacity, cm.size)

	for i := 0; i < cm.capacity; i++ {
		currentNode := cm.table[i].Head
		for currentNode != nil {
			fmt.Fprintf(file, "%s %d\n", currentNode.Key, currentNode.Data)
			currentNode = currentNode.Next
		}
	}

	return nil
}

func (cm *ChainMap) ReadText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %s", filename)
	}
	defer file.Close()

	cm.table = nil

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return fmt.Errorf("неверный формат файла")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) != 2 {
		return fmt.Errorf("неверный формат файла")
	}

	cm.capacity, _ = strconv.Atoi(fields[0])
	cm.size, _ = strconv.Atoi(fields[1])

	cm.table = make([]*Bucket, cm.capacity)
	for i := range cm.table {
		cm.table[i] = NewBucket()
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		spacePos := strings.LastIndex(line, " ")
		if spacePos == -1 {
			return fmt.Errorf("неверный формат файла")
		}

		key := line[:spacePos]
		data, err := strconv.Atoi(line[spacePos+1:])
		if err != nil {
			return fmt.Errorf("неверный формат файла")
		}

		index := cm.hashFunction(key)
		newNode := NewChainNode(key, data)
		cm.appendNode(index, newNode)
	}

	return scanner.Err()
}