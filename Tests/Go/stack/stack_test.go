package stack

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"strconv"
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

func TestNewStack(t *testing.T) {
	s := NewStack()
	if s.head != nil {
		t.Error("NewStack head should be nil")
	}
	if s.size != 0 {
		t.Error("NewStack size should be 0")
	}
}

func TestNewStackFromSlice(t *testing.T) {
	s := NewStackFromSlice("a", "b", "c")
	if s.size != 3 {
		t.Error("NewStackFromSlice size should be 3")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "c" {
		t.Error("First pop should return 'c'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "b" {
		t.Error("Second pop should return 'b'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "a" {
		t.Error("Third pop should return 'a'")
	}
}

func TestPush(t *testing.T) {
	s := NewStack()
	
	err := s.Push("a")
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}
	if s.size != 1 {
		t.Error("Size should be 1 after push")
	}
	if s.head == nil || s.head.key != "a" {
		t.Error("Head should be 'a' after push")
	}
	
	err = s.Push("b")
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}
	if s.size != 2 {
		t.Error("Size should be 2 after second push")
	}
	if s.head.key != "b" {
		t.Error("Head should be 'b' after second push")
	}
}

func TestPushOverflow(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < MAX_SIZE; i++ {
		err := s.Push("item" + strconv.Itoa(i))
		if err != nil {
			t.Fatalf("Push %d failed: %v", i, err)
		}
	}
	
	err := s.Push("overflow")
	if err == nil {
		t.Error("Expected error for stack overflow")
	}
}

func TestPop(t *testing.T) {
	s := NewStack()
	
	_, err := s.Pop()
	if err == nil {
		t.Error("Expected error for pop from empty stack")
	}
	
	s.Push("a")
	s.Push("b")
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "b" {
		t.Error("Pop should return 'b'")
	}
	if s.size != 1 {
		t.Error("Size should be 1 after pop")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "a" {
		t.Error("Pop should return 'a'")
	}
	if s.size != 0 {
		t.Error("Size should be 0 after popping all items")
	}
	if s.head != nil {
		t.Error("Head should be nil after popping all items")
	}
}

func TestIsEmpty(t *testing.T) {
	s := NewStack()
	
	if !s.IsEmpty() {
		t.Error("New stack should be empty")
	}
	
	s.Push("a")
	if s.IsEmpty() {
		t.Error("Stack with items should not be empty")
	}
	
	s.Pop()
	if !s.IsEmpty() {
		t.Error("Stack after popping all items should be empty")
	}
}

func TestGetSize(t *testing.T) {
	s := NewStack()
	
	if s.GetSize() != 0 {
		t.Error("Initial size should be 0")
	}
	
	s.Push("a")
	if s.GetSize() != 1 {
		t.Error("Size should be 1 after push")
	}
	
	s.Push("b")
	if s.GetSize() != 2 {
		t.Error("Size should be 2 after second push")
	}
	
	s.Pop()
	if s.GetSize() != 1 {
		t.Error("Size should be 1 after pop")
	}
}

func TestWriteBinary(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.bin")
	
	s := NewStackFromSlice("a", "b", "c")
	err := s.WriteBinary(filename)
	if err != nil {
		t.Fatalf("WriteBinary failed: %v", err)
	}
}

func TestReadBinary(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.bin")
	
	original := NewStackFromSlice("a", "b", "c")
	original.WriteBinary(filename)
	
	s := NewStack()
	err := s.ReadBinary(filename)
	if err != nil {
		t.Fatalf("ReadBinary failed: %v", err)
	}
	
	if s.GetSize() != 3 {
		t.Error("Size should be 3 after reading")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "c" {
		t.Error("First pop should return 'c'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "b" {
		t.Error("Second pop should return 'b'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "a" {
		t.Error("Third pop should return 'a'")
	}
}

func TestReadBinaryOverflow(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "overflow.bin")
	
	content := make([]byte, 0)
	content = append(content, 0x0B, 0x00, 0x00, 0x00)
	for i := 0; i < 11; i++ {
		item := "item" + strconv.Itoa(i)
		itemBytes := []byte(item)
		content = append(content, byte(len(itemBytes)), 0, 0, 0)
		content = append(content, itemBytes...)
	}
	
	os.WriteFile(filename, content, 0644)
	
	s := NewStack()
	err := s.ReadBinary(filename)
	if err == nil {
		t.Error("Expected error for stack size exceeding MAX_SIZE")
	}
}

func TestWriteText(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.txt")
	
	s := NewStackFromSlice("a", "b", "c")
	err := s.WriteText(filename)
	if err != nil {
		t.Fatalf("WriteText failed: %v", err)
	}
}

func TestReadText(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.txt")
	
	original := NewStackFromSlice("a", "b", "c")
	original.WriteText(filename)
	
	s := NewStack()
	err := s.ReadText(filename)
	if err != nil {
		t.Fatalf("ReadText failed: %v", err)
	}
	
	if s.GetSize() != 3 {
		t.Error("Size should be 3 after reading")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "c" {
		t.Error("First pop should return 'c'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "b" {
		t.Error("Second pop should return 'b'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "a" {
		t.Error("Third pop should return 'a'")
	}
}

func TestReadTextOverflow(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "overflow.txt")
	
	content := "11\n"
	for i := 0; i < 11; i++ {
		content += "item" + strconv.Itoa(i) + "\n"
	}
	os.WriteFile(filename, []byte(content), 0644)
	
	s := NewStack()
	err := s.ReadText(filename)
	if err == nil {
		t.Error("Expected error for stack size exceeding MAX_SIZE")
	}
}

func TestFileErrorHandling(t *testing.T) {
	s := NewStack()
	
	err := s.ReadBinary("nonexistent.bin")
	if err == nil {
		t.Error("Expected error for non-existent binary file")
	}
	
	err = s.ReadText("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent text file")
	}
}

func TestPrint(t *testing.T) {
	s := NewStackFromSlice("a", "b", "c")
	output := captureOutput(s.Print)
	
	if output == "" {
		t.Error("Print should produce output")
	}
	
	s2 := NewStack()
	output2 := captureOutput(s2.Print)
	if output2 == "" {
		t.Error("Print for empty stack should produce output")
	}
}

func TestClear(t *testing.T) {
	s := NewStackFromSlice("a", "b", "c")
	
	if s.GetSize() != 3 {
		t.Error("Size should be 3 before clear")
	}
	
	s.Clear()
	
	if s.GetSize() != 0 {
		t.Error("Size should be 0 after clear")
	}
	if s.IsEmpty() != true {
		t.Error("Stack should be empty after clear")
	}
	if s.head != nil {
		t.Error("Head should be nil after clear")
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("SingleElement", func(t *testing.T) {
		s := NewStack()
		s.Push("a")
		
		if s.GetSize() != 1 {
			t.Error("Size should be 1 after single push")
		}
		if s.IsEmpty() {
			t.Error("Stack should not be empty with single element")
		}
		
		item, err := s.Pop()
		if err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		if item != "a" {
			t.Error("Pop should return 'a'")
		}
		if s.GetSize() != 0 {
			t.Error("Size should be 0 after popping single element")
		}
		if !s.IsEmpty() {
			t.Error("Stack should be empty after popping single element")
		}
	})
	
	t.Run("PushPopSequence", func(t *testing.T) {
		s := NewStack()
		
		for i := 0; i < 5; i++ {
			s.Push("item" + strconv.Itoa(i))
		}
		
		for i := 4; i >= 0; i-- {
			item, err := s.Pop()
			if err != nil {
				t.Fatalf("Pop %d failed: %v", 4-i, err)
			}
			expected := "item" + strconv.Itoa(i)
			if item != expected {
				t.Errorf("Pop %d should return %s, got %s", 4-i, expected, item)
			}
		}
		
		if !s.IsEmpty() {
			t.Error("Stack should be empty after popping all items")
		}
	})
	
	t.Run("MultipleClears", func(t *testing.T) {
		s := NewStackFromSlice("a", "b", "c")
		s.Clear()
		s.Clear()
		
		if s.GetSize() != 0 {
			t.Error("Size should be 0 after multiple clears")
		}
		if !s.IsEmpty() {
			t.Error("Stack should be empty after multiple clears")
		}
	})
}

func TestWriteReadBinaryEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.bin")
	
	s := NewStack()
	s.WriteBinary(filename)
	
	s2 := NewStack()
	s2.ReadBinary(filename)
	
	if s2.GetSize() != 0 {
		t.Error("Size should be 0 after reading empty stack")
	}
	if !s2.IsEmpty() {
		t.Error("Stack should be empty after reading empty file")
	}
}

func TestWriteReadTextEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.txt")
	
	s := NewStack()
	s.WriteText(filename)
	
	s2 := NewStack()
	s2.ReadText(filename)
	
	if s2.GetSize() != 0 {
		t.Error("Size should be 0 after reading empty stack")
	}
	if !s2.IsEmpty() {
		t.Error("Stack should be empty after reading empty file")
	}
}

func TestLargeStack(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < MAX_SIZE; i++ {
		err := s.Push("item" + strconv.Itoa(i))
		if err != nil {
			t.Fatalf("Push %d failed: %v", i, err)
		}
	}
	
	if s.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE after pushing MAX_SIZE items")
	}
	
	for i := MAX_SIZE - 1; i >= 0; i-- {
		item, err := s.Pop()
		if err != nil {
			t.Fatalf("Pop %d failed: %v", MAX_SIZE-1-i, err)
		}
		expected := "item" + strconv.Itoa(i)
		if item != expected {
			t.Errorf("Pop %d should return %s, got %s", MAX_SIZE-1-i, expected, item)
		}
	}
	
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}

func TestWriteReadBinaryLarge(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "large.bin")
	
	s := NewStack()
	for i := 0; i < MAX_SIZE; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	s.WriteBinary(filename)
	
	s2 := NewStack()
	s2.ReadBinary(filename)
	
	if s2.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE after reading large stack")
	}
	
	for i := MAX_SIZE - 1; i >= 0; i-- {
		item, err := s2.Pop()
		if err != nil {
			t.Fatalf("Pop %d failed: %v", MAX_SIZE-1-i, err)
		}
		expected := "item" + strconv.Itoa(i)
		if item != expected {
			t.Errorf("Pop %d should return %s, got %s", MAX_SIZE-1-i, expected, item)
		}
	}
}

func TestWriteReadTextLarge(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "large.txt")
	
	s := NewStack()
	for i := 0; i < MAX_SIZE; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	s.WriteText(filename)
	
	s2 := NewStack()
	s2.ReadText(filename)
	
	if s2.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE after reading large stack")
	}
	
	for i := MAX_SIZE - 1; i >= 0; i-- {
		item, err := s2.Pop()
		if err != nil {
			t.Fatalf("Pop %d failed: %v", MAX_SIZE-1-i, err)
		}
		expected := "item" + strconv.Itoa(i)
		if item != expected {
			t.Errorf("Pop %d should return %s, got %s", MAX_SIZE-1-i, expected, item)
		}
	}
}

func TestEmptyStackOperations(t *testing.T) {
	s := NewStack()
	
	if s.GetSize() != 0 {
		t.Error("Empty stack size should be 0")
	}
	
	if !s.IsEmpty() {
		t.Error("Empty stack should return true for IsEmpty")
	}
	
	_, err := s.Pop()
	if err == nil {
		t.Error("Pop from empty stack should return error")
	}
	
	s.Clear()
	if s.GetSize() != 0 {
		t.Error("Clear on empty stack should keep size 0")
	}
}

func TestStackWithSpecialCharacters(t *testing.T) {
	s := NewStack()
	
	specialStrings := []string{"", " ", "\n", "\t", "hello\nworld", "test\ttab", "special@#$"}
	
	for _, str := range specialStrings {
		err := s.Push(str)
		if err != nil {
			t.Fatalf("Push special string failed: %v", err)
		}
	}
	
	if s.GetSize() != len(specialStrings) {
		t.Error("Size should match number of special strings pushed")
	}
	
	for i := len(specialStrings) - 1; i >= 0; i-- {
		item, err := s.Pop()
		if err != nil {
			t.Fatalf("Pop special string failed: %v", err)
		}
		if item != specialStrings[i] {
			t.Errorf("Pop special string: expected %s, got %s", specialStrings[i], item)
		}
	}
	
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all special strings")
	}
}

func TestWriteReadBinarySpecialChars(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "special.bin")
	
	s := NewStack()
	specialStrings := []string{"", " ", "\n", "\t", "hello\nworld", "test\ttab", "special@#$"}
	
	for _, str := range specialStrings {
		s.Push(str)
	}
	
	s.WriteBinary(filename)
	
	s2 := NewStack()
	s2.ReadBinary(filename)
	
	if s2.GetSize() != len(specialStrings) {
		t.Error("Size should match after reading special chars")
	}
	
	for i := len(specialStrings) - 1; i >= 0; i-- {
		item, err := s2.Pop()
		if err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		if item != specialStrings[i] {
			t.Errorf("Read special string: expected %s, got %s", specialStrings[i], item)
		}
	}
}

func TestWriteReadWithDifferentSizes(t *testing.T) {
	tempDir := t.TempDir()
	binFile := filepath.Join(tempDir, "test.bin")
	txtFile := filepath.Join(tempDir, "test.txt")
	
	sizes := []int{1, 3, 5, MAX_SIZE - 1}
	
	for _, size := range sizes {
		s := NewStack()
		for i := 0; i < size; i++ {
			s.Push("item" + strconv.Itoa(i))
		}
		
		s.WriteBinary(binFile)
		s2 := NewStack()
		s2.ReadBinary(binFile)
		
		if s2.GetSize() != size {
			t.Errorf("Binary read size mismatch for size %d", size)
		}
		
		s.WriteText(txtFile)
		s3 := NewStack()
		s3.ReadText(txtFile)
		
		if s3.GetSize() != size {
			t.Errorf("Text read size mismatch for size %d", size)
		}
	}
}

func TestClearAfterOperations(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < 5; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	if s.GetSize() != 5 {
		t.Error("Size should be 5 after adding items")
	}
	
	for i := 0; i < 2; i++ {
		_, err := s.Pop()
		if err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
	}
	
	if s.GetSize() != 3 {
		t.Error("Size should be 3 after removing 2 items")
	}
	
	s.Clear()
	
	if s.GetSize() != 0 {
		t.Error("Size should be 0 after clear")
	}
	if !s.IsEmpty() {
		t.Error("Stack should be empty after clear")
	}
	
	s.Push("new")
	if s.GetSize() != 1 {
		t.Error("Size should be 1 after adding after clear")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "new" {
		t.Error("Should pop 'new' after clear and add")
	}
}

func TestPopUntilEmpty(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < 3; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	for i := 0; i < 3; i++ {
		_, err := s.Pop()
		if err != nil {
			t.Fatalf("Pop %d failed: %v", i, err)
		}
	}
	
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
	
	_, err := s.Pop()
	if err == nil {
		t.Error("Expected error for pop from empty stack")
	}
}

func TestWriteReadEmptyStack(t *testing.T) {
	tempDir := t.TempDir()
	binFile := filepath.Join(tempDir, "empty.bin")
	txtFile := filepath.Join(tempDir, "empty.txt")
	
	s := NewStack()
	
	s.WriteBinary(binFile)
	s.WriteText(txtFile)
	
	s2 := NewStack()
	s2.ReadBinary(binFile)
	
	if s2.GetSize() != 0 {
		t.Error("Binary read should result in empty stack")
	}
	
	s3 := NewStack()
	s3.ReadText(txtFile)
	
	if s3.GetSize() != 0 {
		t.Error("Text read should result in empty stack")
	}
}

func TestStackOrder(t *testing.T) {
	s := NewStack()
	
	s.Push("a")
	s.Push("b")
	s.Push("c")
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "c" {
		t.Error("First pop should return 'c'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "b" {
		t.Error("Second pop should return 'b'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "a" {
		t.Error("Third pop should return 'a'")
	}
	
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}

func TestWriteReadInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()
	binFile := filepath.Join(tempDir, "invalid.bin")
	txtFile := filepath.Join(tempDir, "invalid.txt")
	
	os.WriteFile(binFile, []byte{0x01, 0x02, 0x03}, 0644)
	os.WriteFile(txtFile, []byte("invalid\ncontent"), 0644)
	
	s := NewStack()
	
	err := s.ReadBinary(binFile)
	if err == nil {
		t.Error("Expected error for invalid binary file")
	}
	
	err = s.ReadText(txtFile)
	if err == nil {
		t.Error("Expected error for invalid text file")
	}
}

func TestStackWithLargeStrings(t *testing.T) {
	s := NewStack()
	
	largeString := "a" + string(make([]byte, 1000)) + "z"
	err := s.Push(largeString)
	if err != nil {
		t.Fatalf("Push large string failed: %v", err)
	}
	
	if s.GetSize() != 1 {
		t.Error("Size should be 1 after pushing large string")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop large string failed: %v", err)
	}
	
	if item != largeString {
		t.Error("Popped string should match original large string")
	}
	
	if s.GetSize() != 0 {
		t.Error("Size should be 0 after popping large string")
	}
}

func TestMultipleWriteReadCycles(t *testing.T) {
	tempDir := t.TempDir()
	binFile := filepath.Join(tempDir, "cycle.bin")
	txtFile := filepath.Join(tempDir, "cycle.txt")
	
	for cycle := 0; cycle < 3; cycle++ {
		s := NewStack()
		
		for i := 0; i < cycle+1; i++ {
			s.Push("cycle" + string(rune(cycle)) + "_item" + strconv.Itoa(i))
		}
		
		s.WriteBinary(binFile)
		s2 := NewStack()
		s2.ReadBinary(binFile)
		
		if s2.GetSize() != cycle+1 {
			t.Errorf("Binary cycle %d: size mismatch", cycle)
		}
		
		s.WriteText(txtFile)
		s3 := NewStack()
		s3.ReadText(txtFile)
		
		if s3.GetSize() != cycle+1 {
			t.Errorf("Text cycle %d: size mismatch", cycle)
		}
	}
}

func TestNodeStructure(t *testing.T) {
	s := NewStack()
	
	s.Push("a")
	s.Push("b")
	s.Push("c")
	
	if s.head.key != "c" {
		t.Error("Head should be 'c'")
	}
	if s.head.next.key != "b" {
		t.Error("Head.next should be 'b'")
	}
	if s.head.next.next.key != "a" {
		t.Error("Head.next.next should be 'a'")
	}
	if s.head.next.next.next != nil {
		t.Error("Last node should have nil next")
	}
}

func TestStackCapacity(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < MAX_SIZE; i++ {
		err := s.Push("item" + strconv.Itoa(i))
		if err != nil {
			t.Fatalf("Push %d failed: %v", i, err)
		}
	}
	
	if s.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE after filling")
	}
	
	err := s.Push("overflow")
	if err == nil {
		t.Error("Expected error for push when at MAX_SIZE")
	}
	
	_, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	
	err = s.Push("new")
	if err != nil {
		t.Error("Push after pop should succeed")
	}
	
	if s.GetSize() != MAX_SIZE {
		t.Error("Size should still be MAX_SIZE after pop and push")
	}
}

func TestPushAfterClear(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < 3; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	s.Clear()
	
	for i := 0; i < 2; i++ {
		s.Push("new" + strconv.Itoa(i))
	}
	
	if s.GetSize() != 2 {
		t.Error("Size should be 2 after clear and add")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "new1" {
		t.Errorf("First pop should return 'new1', got %s", item)
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "new0" {
		t.Errorf("Second pop should return 'new0', got %s", item)
	}
}

func TestStackWithMaxSize(t *testing.T) {
	s := NewStack()
	
	for i := 0; i < MAX_SIZE; i++ {
		s.Push("item" + strconv.Itoa(i))
	}
	
	if s.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE")
	}
	
	err := s.Push("overflow")
	if err == nil {
		t.Error("Expected overflow error")
	}
	
	_, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	
	err = s.Push("new")
	if err != nil {
		t.Error("Push after pop should succeed")
	}
	
	if s.GetSize() != MAX_SIZE {
		t.Error("Size should be MAX_SIZE again")
	}
}

func TestStackWithEmptyStrings(t *testing.T) {
	s := NewStack()
	
	s.Push("")
	s.Push("non-empty")
	s.Push("")
	
	if s.GetSize() != 3 {
		t.Error("Size should be 3 after pushing empty strings")
	}
	
	item, err := s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "" {
		t.Error("Should pop empty string")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "non-empty" {
		t.Error("Should pop 'non-empty'")
	}
	
	item, err = s.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if item != "" {
		t.Error("Should pop empty string")
	}
	
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}