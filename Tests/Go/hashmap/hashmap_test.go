package hashmap

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestConstructors(t *testing.T) {
	t.Run("NewChainMap", func(t *testing.T) {
		cm := NewChainMap(5)
		if cm.capacity != 5 {
			t.Errorf("NewChainMap(5) capacity = %d, want 5", cm.capacity)
		}
		if cm.size != 0 {
			t.Errorf("NewChainMap(5) size = %d, want 0", cm.size)
		}
		if len(cm.table) != 5 {
			t.Errorf("NewChainMap(5) table length = %d, want 5", len(cm.table))
		}
	})

	t.Run("NewChainNode", func(t *testing.T) {
		node := NewChainNode("test", 42)
		if node.Key != "test" {
			t.Errorf("NewChainNode key = %s, want 'test'", node.Key)
		}
		if node.Data != 42 {
			t.Errorf("NewChainNode data = %d, want 42", node.Data)
		}
		if node.Next != nil {
			t.Errorf("NewChainNode next = %v, want nil", node.Next)
		}
	})

	t.Run("NewBucket", func(t *testing.T) {
		bucket := NewBucket()
		if bucket.Head != nil {
			t.Errorf("NewBucket head = %v, want nil", bucket.Head)
		}
	})
}

func TestHashFunction(t *testing.T) {
	cm := NewChainMap(10)
	
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		hash := cm.hashFunction(key)
		if hash < 0 || hash >= cm.capacity {
			t.Errorf("hashFunction(%s) = %d, out of bounds [0, %d)", key, hash, cm.capacity)
		}
	}
	
	key := "test_key"
	expectedHash := cm.hashFunction(key)
	for i := 0; i < 10; i++ {
		hash := cm.hashFunction(key)
		if hash != expectedHash {
			t.Errorf("hashFunction(%s) inconsistent: got %d, want %d", key, hash, expectedHash)
		}
	}
}

func TestCoreOperations(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		cm := NewChainMap(5)
		
		cm.Add("key1", 10)
		cm.Add("key2", 20)
		cm.Add("key3", 30)
		
		if cm.size != 3 {
			t.Errorf("Size after adds = %d, want 3", cm.size)
		}
		
		cm.Add("key1", 15)
		if cm.size != 3 {
			t.Errorf("Size after update = %d, want 3", cm.size)
		}
		
		data, err := cm.Find("key1")
		if err != nil {
			t.Fatalf("Find('key1') failed: %v", err)
		}
		if data != 15 {
			t.Errorf("Find('key1') = %d, want 15", data)
		}
		
		data, err = cm.Find("key2")
		if err != nil {
			t.Fatalf("Find('key2') failed: %v", err)
		}
		if data != 20 {
			t.Errorf("Find('key2') = %d, want 20", data)
		}
	})

	t.Run("Del", func(t *testing.T) {
		cm := NewChainMap(5)
		cm.Add("key1", 10)
		cm.Add("key2", 20)
		cm.Add("key3", 30)
		
		cm.Del("key2")
		if cm.size != 2 {
			t.Errorf("Size after delete = %d, want 2", cm.size)
		}
		
		if cm.IsContain("key2") {
			t.Error("IsContain('key2') = true after deletion, want false")
		}
		
		initialSize := cm.size
		cm.Del("nonexistent")
		if cm.size != initialSize {
			t.Errorf("Size after deleting non-existent key = %d, want %d", cm.size, initialSize)
		}
		
		cm.Del("key1")
		cm.Del("key3")
		if cm.size != 0 {
			t.Errorf("Size after deleting all = %d, want 0", cm.size)
		}
	})

	t.Run("Find", func(t *testing.T) {
		cm := NewChainMap(5)
		cm.Add("key1", 10)
		cm.Add("key2", 20)
		
		data, err := cm.Find("key1")
		if err != nil {
			t.Fatalf("Find('key1') failed: %v", err)
		}
		if data != 10 {
			t.Errorf("Find('key1') = %d, want 10", data)
		}
		
		_, err = cm.Find("nonexistent")
		if err == nil {
			t.Error("Find('nonexistent') returned nil error, want error")
		}
	})

	t.Run("IsContain", func(t *testing.T) {
		cm := NewChainMap(5)
		cm.Add("key1", 10)
		
		if !cm.IsContain("key1") {
			t.Error("IsContain('key1') = false, want true")
		}
		
		if cm.IsContain("nonexistent") {
			t.Error("IsContain('nonexistent') = true, want false")
		}
	})

	t.Run("Rehashing", func(t *testing.T) {
		cm := NewChainMap(2)
		
		for i := 0; i < 2; i++ {
			cm.Add(fmt.Sprintf("key%d", i), i)
		}
		
		cm.Add("key2", 2)
		
		if cm.capacity != 4 {
			t.Errorf("Capacity after rehash = %d, want 4", cm.capacity)
		}
		
		for i := 0; i < 3; i++ {
			key := fmt.Sprintf("key%d", i)
			data, err := cm.Find(key)
			if err != nil {
				t.Fatalf("Find('%s') failed: %v", key, err)
			}
			if data != i {
				t.Errorf("Find('%s') = %d, want %d", key, data, i)
			}
		}
	})
}

func TestCollisionHandling(t *testing.T) {
	cm := NewChainMap(3)
	
	keys := []string{"abc", "def", "ghi", "jkl", "mno"}
	for i, key := range keys {
		cm.Add(key, i*10)
	}
	
	if cm.size != len(keys) {
		t.Errorf("Size = %d, want %d", cm.size, len(keys))
	}
	
	for i, key := range keys {
		data, err := cm.Find(key)
		if err != nil {
			t.Fatalf("Find('%s') failed: %v", key, err)
		}
		if data != i*10 {
			t.Errorf("Find('%s') = %d, want %d", key, data, i*10)
		}
	}
	
	cm.Del("def")
	if cm.IsContain("def") {
		t.Error("IsContain('def') = true after deletion, want false")
	}
	
	data, err := cm.Find("abc")
	if err != nil {
		t.Fatalf("Find('abc') failed after deletion: %v", err)
	}
	if data != 0 {
		t.Errorf("Find('abc') = %d, want 0", data)
	}
}

func TestGetAllKeys(t *testing.T) {
	cm := NewChainMap(5)
	cm.Add("key1", 10)
	cm.Add("key2", 20)
	cm.Add("key3", 30)
	
	result := NewChainMap(5)
	cm.GetAllKeys(result)
	
	if result.size != 3 {
		t.Errorf("GetAllKeys result size = %d, want 3", result.size)
	}
	
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		if !result.IsContain(key) {
			t.Errorf("GetAllKeys result doesn't contain '%s'", key)
		}
	}
}

func TestGetAllKeysAsString(t *testing.T) {
	cm := NewChainMap(5)
	cm.Add("key1", 10)
	cm.Add("key2", 20)
	cm.Add("key3", 30)
	
	result := cm.GetAllKeysAsString()
	
	expectedKeys := []string{"key1", "key2", "key3"}
	for _, key := range expectedKeys {
		if !strings.Contains(result, key) {
			t.Errorf("GetAllKeysAsString result doesn't contain '%s': %s", key, result)
		}
	}
}

func TestFileOperations(t *testing.T) {
	tempDir := t.TempDir()
	
	t.Run("BinaryFileOperations", func(t *testing.T) {
		binFile := filepath.Join(tempDir, "test.bin")
		
		original := NewChainMap(5)
		original.Add("key1", 10)
		original.Add("key2", 20)
		original.Add("key3", 30)
		original.Add("key with spaces", 40)
		original.Add("", 50)
		
		if err := original.WriteBinary(binFile); err != nil {
			t.Fatalf("WriteBinary() failed: %v", err)
		}
		
		readMap := NewChainMap(1)
		if err := readMap.ReadBinary(binFile); err != nil {
			t.Fatalf("ReadBinary() failed: %v", err)
		}
		
		if readMap.size != original.size {
			t.Errorf("ReadBinary size = %d, want %d", readMap.size, original.size)
		}
		
		keys := []string{"key1", "key2", "key3", "key with spaces", ""}
		expectedValues := []int{10, 20, 30, 40, 50}
		for i, key := range keys {
			if !readMap.IsContain(key) {
				t.Errorf("ReadBinary missing key '%s'", key)
			}
			
			data, err := readMap.Find(key)
			if err != nil {
				t.Fatalf("Find('%s') failed: %v", key, err)
			}
			if data != expectedValues[i] {
				t.Errorf("Find('%s') = %d, want %d", key, data, expectedValues[i])
			}
		}
	})
	
	t.Run("TextFileOperations", func(t *testing.T) {
		txtFile := filepath.Join(tempDir, "test.txt")
		
		original := NewChainMap(5)
		original.Add("key1", 10)
		original.Add("key2", 20)
		original.Add("key3", 30)
		original.Add("key with spaces", 40)
		original.Add("", 50)
		
		if err := original.WriteText(txtFile); err != nil {
			t.Fatalf("WriteText() failed: %v", err)
		}
		
		readMap := NewChainMap(1)
		if err := readMap.ReadText(txtFile); err != nil {
			t.Fatalf("ReadText() failed: %v", err)
		}
		
		if readMap.size != original.size {
			t.Errorf("ReadText size = %d, want %d", readMap.size, original.size)
		}
		
		keys := []string{"key1", "key2", "key3", "key with spaces", ""}
		expectedValues := []int{10, 20, 30, 40, 50}
		for i, key := range keys {
			if !readMap.IsContain(key) {
				t.Errorf("ReadText missing key '%s'", key)
			}
			
			data, err := readMap.Find(key)
			if err != nil {
				t.Fatalf("Find('%s') failed: %v", key, err)
			}
			if data != expectedValues[i] {
				t.Errorf("Find('%s') = %d, want %d", key, data, expectedValues[i])
			}
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		cm := NewChainMap(5)
		
		if cm.size != 0 {
			t.Errorf("Empty map size = %d, want 0", cm.size)
		}
		
		if cm.IsContain("nonexistent") {
			t.Error("IsContain on empty map returned true, want false")
		}
		
		_, err := cm.Find("nonexistent")
		if err == nil {
			t.Error("Find on empty map returned nil error, want error")
		}
		
		cm.Del("nonexistent")
		if cm.size != 0 {
			t.Errorf("Size after deleting from empty map = %d, want 0", cm.size)
		}
	})
	
	t.Run("LargeValues", func(t *testing.T) {
		cm := NewChainMap(5)
		
		cm.Add("large", 2147483647)
		data, err := cm.Find("large")
		if err != nil {
			t.Fatalf("Find('large') failed: %v", err)
		}
		if data != 2147483647 {
			t.Errorf("Find('large') = %d, want 2147483647", data)
		}
		
		cm.Add("negative", -1000)
		data, err = cm.Find("negative")
		if err != nil {
			t.Fatalf("Find('negative') failed: %v", err)
		}
		if data != -1000 {
			t.Errorf("Find('negative') = %d, want -1000", data)
		}
	})
	
	t.Run("EmptyStrings", func(t *testing.T) {
		cm := NewChainMap(5)
		
		cm.Add("", 42)
		if !cm.IsContain("") {
			t.Error("Empty key not found after insertion")
		}
		
		data, err := cm.Find("")
		if err != nil {
			t.Fatalf("Find('') failed: %v", err)
		}
		if data != 42 {
			t.Errorf("Find('') = %d, want 42", data)
		}
		
		cm.Add("empty_value", 0)
		data, err = cm.Find("empty_value")
		if err != nil {
			t.Fatalf("Find('empty_value') failed: %v", err)
		}
		if data != 0 {
			t.Errorf("Find('empty_value') = %d, want 0", data)
		}
	})
	
	t.Run("CollisionHandling", func(t *testing.T) {
		cm := NewChainMap(1)
		
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key%d", i)
			cm.Add(key, i)
		}
		
		if cm.size != 10 {
			t.Errorf("Size after collisions = %d, want 10", cm.size)
		}
		
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key%d", i)
			data, err := cm.Find(key)
			if err != nil {
				t.Fatalf("Find('%s') failed: %v", key, err)
			}
			if data != i {
				t.Errorf("Find('%s') = %d, want %d", key, data, i)
			}
		}
		
		cm.Del("key5")
		if cm.IsContain("key5") {
			t.Error("key5 found after deletion in collision chain")
		}
		
		for i := 0; i < 10; i++ {
			if i == 5 {
				continue
			}
			key := fmt.Sprintf("key%d", i)
			if !cm.IsContain(key) {
				t.Errorf("Key '%s' missing after deletion in collision chain", key)
			}
		}
	})
}

func TestPrintContents(t *testing.T) {
	cm := NewChainMap(3)
	cm.Add("key1", 10)
	cm.Add("key2", 20)
	cm.Add("key3", 30)
	
	output := captureOutput(cm.PrintContents)
	
	if !strings.Contains(output, "Содержимое хеш-таблицы:") {
		t.Errorf("PrintContents output doesn't contain header: %s", output)
	}
	
	if !strings.Contains(output, "key1 -> 10") {
		t.Errorf("PrintContents output doesn't contain key1: %s", output)
	}
	
	if !strings.Contains(output, "key2 -> 20") {
		t.Errorf("PrintContents output doesn't contain key2: %s", output)
	}
	
	if !strings.Contains(output, "key3 -> 30") {
		t.Errorf("PrintContents output doesn't contain key3: %s", output)
	}
}

func TestAppendNode(t *testing.T) {
	cm := NewChainMap(1)
	
	cm.Add("key1", 10)
	cm.Add("key2", 20)
	cm.Add("key3", 30)
	
	if !cm.IsContain("key1") || !cm.IsContain("key2") || !cm.IsContain("key3") {
		t.Error("Not all keys found in map")
	}
}

func TestFindError(t *testing.T) {
	cm := NewChainMap(5)
	
	_, err := cm.Find("nonexistent")
	if err == nil {
		t.Error("Find() expected error for non-existent key, got nil")
	}
	
	expectedError := "в словаре нет такого ключа"
	if err.Error() != expectedError {
		t.Errorf("Find() error = %s, want %s", err.Error(), expectedError)
	}
}

func TestRehashingComplex(t *testing.T) {
	cm := NewChainMap(2)
	
	for i := 0; i < 5; i++ {
		cm.Add(fmt.Sprintf("key%d", i), i)
	}
	
	if cm.capacity < 4 {
		t.Errorf("Capacity after rehash = %d, want at least 4", cm.capacity)
	}
	
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key%d", i)
		data, err := cm.Find(key)
		if err != nil {
			t.Fatalf("Find('%s') failed: %v", key, err)
		}
		if data != i {
			t.Errorf("Find('%s') = %d, want %d", key, data, i)
		}
	}
	
	cm.Del("key2")
	cm.Add("key5", 5)
	
	remainingKeys := []struct{key string; value int}{ {"key0", 0}, {"key1", 1}, {"key3", 3}, {"key4", 4}, {"key5", 5} }
	for _, kv := range remainingKeys {
		data, err := cm.Find(kv.key)
		if err != nil {
			t.Fatalf("Find('%s') failed: %v", kv.key, err)
		}
		if data != kv.value {
			t.Errorf("Find('%s') = %d, want %d", kv.key, data, kv.value)
		}
	}
}