package queue

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	if q.head != nil {
		t.Error("NewQueue head should be nil")
	}
	if q.tail != nil {
		t.Error("NewQueue tail should be nil")
	}
	if q.size != 0 {
		t.Error("NewQueue size should be 0")
	}
	if q.maxSize != MAX_SIZE {
		t.Error("NewQueue maxSize should be MAX_SIZE")
	}
}

func TestNewQueueWithItems(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	if q.size != 3 {
		t.Error("NewQueueWithItems size should be 3")
	}
	if q.head.Data != "a" {
		t.Error("Head should be 'a'")
	}
	if q.tail.Data != "c" {
		t.Error("Tail should be 'c'")
	}
}

func TestEnqueue(t *testing.T) {
	q := NewQueue()
	
	err := q.Enqueue("a")
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}
	if q.size != 1 {
		t.Error("Size should be 1 after enqueue")
	}
	if q.head.Data != "a" || q.tail.Data != "a" {
		t.Error("Head and tail should be 'a' after first enqueue")
	}
	
	err = q.Enqueue("b")
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}
	if q.size != 2 {
		t.Error("Size should be 2 after second enqueue")
	}
	if q.head.Data != "a" || q.tail.Data != "b" {
		t.Error("Head should be 'a' and tail should be 'b'")
	}
}

func TestEnqueueOverflow(t *testing.T) {
	q := NewQueue()
	q.maxSize = 2
	
	err := q.Enqueue("a")
	if err != nil {
		t.Fatalf("First enqueue failed: %v", err)
	}
	err = q.Enqueue("b")
	if err != nil {
		t.Fatalf("Second enqueue failed: %v", err)
	}
	
	err = q.Enqueue("c")
	if err == nil {
		t.Error("Expected error for queue overflow")
	}
}

func TestDequeue(t *testing.T) {
	q := NewQueue()
	q.Enqueue("a")
	q.Enqueue("b")
	
	data, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if data != "a" {
		t.Error("Dequeue should return 'a'")
	}
	if q.size != 1 {
		t.Error("Size should be 1 after dequeue")
	}
	if q.head.Data != "b" || q.tail.Data != "b" {
		t.Error("Head and tail should be 'b' after dequeue")
	}
	
	data, err = q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if data != "b" {
		t.Error("Dequeue should return 'b'")
	}
	if q.size != 0 {
		t.Error("Size should be 0 after dequeueing all elements")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after dequeueing all elements")
	}
}

func TestDequeueUnderflow(t *testing.T) {
	q := NewQueue()
	_, err := q.Dequeue()
	if err == nil {
		t.Error("Expected error for queue underflow")
	}
}

func TestDel(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c", "d")
	
	q.Del("b")
	if q.size != 3 {
		t.Error("Size should be 3 after deletion")
	}
	
	current := q.head
	expected := []string{"a", "c", "d"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
	
	q.Del("a")
	q.Del("d")
	if q.size != 1 {
		t.Error("Size should be 1 after deleting head and tail")
	}
	if q.head.Data != "c" || q.tail.Data != "c" {
		t.Error("Head and tail should be 'c' after deleting head and tail")
	}
	
	q.Del("c")
	if q.size != 0 {
		t.Error("Size should be 0 after deleting last element")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after deleting last element")
	}
}

func TestDelNonExistent(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	initialSize := q.size
	q.Del("nonexistent")
	if q.size != initialSize {
		t.Error("Size should not change when deleting non-existent element")
	}
}

func TestSize(t *testing.T) {
	q := NewQueue()
	if q.Size() != 0 {
		t.Error("Initial size should be 0")
	}
	
	q.Enqueue("a")
	if q.Size() != 1 {
		t.Error("Size should be 1 after enqueue")
	}
	
	q.Enqueue("b")
	if q.Size() != 2 {
		t.Error("Size should be 2 after second enqueue")
	}
	
	q.Dequeue()
	if q.Size() != 1 {
		t.Error("Size should be 1 after dequeue")
	}
}

func TestHead(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	head := q.Head()
	if head == nil {
		t.Error("Head should not be nil")
	}
	if head.Data != "a" {
		t.Error("Head data should be 'a'")
	}
}

func TestClear(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Clear()
	if q.size != 0 {
		t.Error("Size should be 0 after clear")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after clear")
	}
}

func TestWriteBinary(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.bin")
	
	q := NewQueueWithItems("hello", "world", "test")
	err := q.WriteBinary(filename)
	if err != nil {
		t.Fatalf("WriteBinary failed: %v", err)
	}
}

func TestReadBinary(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.bin")
	
	original := NewQueueWithItems("hello", "world", "test")
	original.WriteBinary(filename)
	
	q := NewQueue()
	err := q.ReadBinary(filename)
	if err != nil {
		t.Fatalf("ReadBinary failed: %v", err)
	}
	
	if q.size != 3 {
		t.Error("Size should be 3 after reading")
	}
	
	current := q.head
	expected := []string{"hello", "world", "test"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestReadBinaryOverflow(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.bin")
	
	q := NewQueue()
	q.maxSize = 2
	
	original := NewQueueWithItems("a", "b", "c")
	original.WriteBinary(filename)
	
	err := q.ReadBinary(filename)
	if err == nil {
		t.Error("Expected error for queue size exceeding max")
	}
}

func TestWriteText(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.txt")
	
	q := NewQueueWithItems("hello", "world", "test")
	err := q.WriteText(filename)
	if err != nil {
		t.Fatalf("WriteText failed: %v", err)
	}
}

func TestReadText(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.txt")
	
	original := NewQueueWithItems("hello", "world", "test")
	original.WriteText(filename)
	
	q := NewQueue()
	err := q.ReadText(filename)
	if err != nil {
		t.Fatalf("ReadText failed: %v", err)
	}
	
	if q.size != 3 {
		t.Error("Size should be 3 after reading")
	}
	
	current := q.head
	expected := []string{"hello", "world", "test"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestReadTextOverflow(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.txt")
	
	q := NewQueue()
	q.maxSize = 2
	
	original := NewQueueWithItems("a", "b", "c")
	original.WriteText(filename)
	
	err := q.ReadText(filename)
	if err == nil {
		t.Error("Expected error for queue size exceeding max")
	}
}

func TestFileErrorHandling(t *testing.T) {
	q := NewQueue()
	
	err := q.ReadBinary("nonexistent.bin")
	if err == nil {
		t.Error("Expected error for non-existent binary file")
	}
	
	err = q.ReadText("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent text file")
	}
}

func TestPrint(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	output := captureOutput(q.Print)
	expected := "a b c \n"
	if output != expected {
		t.Errorf("Print output = %q, want %q", output, expected)
	}
	
	q.Clear()
	output = captureOutput(q.Print)
	expected = "Queue is empty\n"
	if output != expected {
		t.Errorf("Print output = %q, want %q", output, expected)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("SingleElement", func(t *testing.T) {
		q := NewQueue()
		q.Enqueue("a")
		
		if q.head != q.tail {
			t.Error("Head and tail should be same for single element")
		}
		
		data, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
		if data != "a" {
			t.Error("Dequeue should return 'a'")
		}
		if q.size != 0 {
			t.Error("Size should be 0 after dequeueing single element")
		}
	})
	
	t.Run("DelHead", func(t *testing.T) {
		q := NewQueueWithItems("a", "b", "c")
		q.Del("a")
		if q.head.Data != "b" {
			t.Error("New head should be 'b'")
		}
		if q.size != 2 {
			t.Error("Size should be 2 after deleting head")
		}
	})
	
	t.Run("DelTail", func(t *testing.T) {
		q := NewQueueWithItems("a", "b", "c")
		q.Del("c")
		if q.tail.Data != "b" {
			t.Error("New tail should be 'b'")
		}
		if q.size != 2 {
			t.Error("Size should be 2 after deleting tail")
		}
	})
	
	t.Run("DelMiddle", func(t *testing.T) {
		q := NewQueueWithItems("a", "b", "c")
		q.Del("b")
		if q.size != 2 {
			t.Error("Size should be 2 after deleting middle")
		}
		
		current := q.head
		if current.Data != "a" {
			t.Error("First element should be 'a'")
		}
		current = current.Next
		if current.Data != "c" {
			t.Error("Second element should be 'c'")
		}
	})
	
	t.Run("MaxSize", func(t *testing.T) {
		q := NewQueue()
		q.maxSize = 1
		
		err := q.Enqueue("a")
		if err != nil {
			t.Fatalf("First enqueue failed: %v", err)
		}
		
		err = q.Enqueue("b")
		if err == nil {
			t.Error("Expected error for queue overflow")
		}
	})
}

func TestNodeStructure(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	
	if q.head.Prev != nil {
		t.Error("Head should have nil Prev")
	}
	if q.tail.Next != nil {
		t.Error("Tail should have nil Next")
	}
	
	current := q.head
	for current.Next != nil {
		if current.Next.Prev != current {
			t.Error("Next node's Prev should point back to current")
		}
		current = current.Next
	}
}

func TestDelFromEmpty(t *testing.T) {
	q := NewQueue()
	initialSize := q.size
	q.Del("nonexistent")
	if q.size != initialSize {
		t.Error("Size should not change when deleting from empty queue")
	}
}

func TestDelFromSingleElement(t *testing.T) {
	q := NewQueueWithItems("a")
	q.Del("a")
	if q.size != 0 {
		t.Error("Size should be 0 after deleting single element")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after deleting single element")
	}
}

func TestEmptyQueueOperations(t *testing.T) {
	q := NewQueue()
	
	if q.Size() != 0 {
		t.Error("Empty queue size should be 0")
	}
	
	head := q.Head()
	if head != nil {
		t.Error("Empty queue head should be nil")
	}
}

func TestLargeQueue(t *testing.T) {
	q := NewQueue()
	
	for i := 0; i < 100; i++ {
		q.Enqueue("item" + string(rune(i)))
	}
	
	if q.Size() != 100 {
		t.Error("Size should be 100 after enqueuing 100 items")
	}
	
	for i := 0; i < 50; i++ {
		_, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
	}
	
	if q.Size() != 50 {
		t.Error("Size should be 50 after dequeuing 50 items")
	}
}

func TestDelAllElements(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	
	q.Del("a")
	q.Del("b")
	q.Del("c")
	
	if q.Size() != 0 {
		t.Error("Size should be 0 after deleting all elements")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after deleting all elements")
	}
}

func TestDelMultipleSameElements(t *testing.T) {
	q := NewQueueWithItems("a", "b", "a", "c", "a")
	
	q.Del("a")
	if q.Size() != 4 {
		t.Error("Size should be 4 after deleting first 'a'")
	}
	
	current := q.head
	expected := []string{"b", "a", "c", "a"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestDelFirstElement(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Del("a")
	
	if q.head.Data != "b" {
		t.Error("Head should be 'b' after deleting first element")
	}
	if q.size != 2 {
		t.Error("Size should be 2 after deleting first element")
	}
}

func TestDelLastElement(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Del("c")
	
	if q.tail.Data != "b" {
		t.Error("Tail should be 'b' after deleting last element")
	}
	if q.size != 2 {
		t.Error("Size should be 2 after deleting last element")
	}
}

func TestEnqueueDequeueSequence(t *testing.T) {
	q := NewQueue()
	
	for i := 0; i < 10; i++ {
		q.Enqueue("item" + string(rune(i)))
	}
	
	for i := 0; i < 5; i++ {
		_, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
	}
	
	for i := 10; i < 15; i++ {
		q.Enqueue("item" + string(rune(i)))
	}
	
	if q.Size() != 10 {
		t.Error("Size should be 10 after mixed operations")
	}
}

func TestBinaryReadWriteEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.bin")
	
	q := NewQueue()
	q.WriteBinary(filename)
	
	q2 := NewQueue()
	q2.ReadBinary(filename)
	
	if q2.Size() != 0 {
		t.Error("Size should be 0 after reading empty queue")
	}
}

func TestTextReadWriteEmpty(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.txt")
	
	q := NewQueue()
	q.WriteText(filename)
	
	q2 := NewQueue()
	q2.ReadText(filename)
	
	if q2.Size() != 0 {
		t.Error("Size should be 0 after reading empty queue")
	}
}

func TestBinaryReadWriteWithEmptyStrings(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty_strings.bin")
	
	q := NewQueueWithItems("", "a", "", "b", "")
	q.WriteBinary(filename)
	
	q2 := NewQueue()
	q2.ReadBinary(filename)
	
	if q2.Size() != 5 {
		t.Error("Size should be 5 after reading with empty strings")
	}
	
	current := q2.head
	expected := []string{"", "a", "", "b", ""}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestTextReadWriteWithEmptyStrings(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty_strings.txt")
	
	q := NewQueueWithItems("", "a", "", "b", "")
	q.WriteText(filename)
	
	q2 := NewQueue()
	q2.ReadText(filename)
	
	if q2.Size() != 5 {
		t.Error("Size should be 5 after reading with empty strings")
	}
	
	current := q2.head
	expected := []string{"", "a", "", "b", ""}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestBinaryReadWriteWithSpecialChars(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "special_chars.bin")
	
	q := NewQueueWithItems("hello\nworld", "test\ttab", "space test", "special@#$")
	q.WriteBinary(filename)
	
	q2 := NewQueue()
	q2.ReadBinary(filename)
	
	if q2.Size() != 4 {
		t.Error("Size should be 4 after reading with special chars")
	}
	
	current := q2.head
	expected := []string{"hello\nworld", "test\ttab", "space test", "special@#$"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestWriteBinaryError(t *testing.T) {
	q := NewQueueWithItems("test")
	
	err := q.WriteBinary("/invalid/path/file.bin")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestWriteTextError(t *testing.T) {
	q := NewQueueWithItems("test")
	
	err := q.WriteText("/invalid/path/file.txt")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestReadBinaryInvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "invalid.bin")
	
	os.WriteFile(filename, []byte{0x01, 0x02, 0x03}, 0644)
	
	q := NewQueue()
	err := q.ReadBinary(filename)
	if err == nil {
		t.Error("Expected error for invalid binary file")
	}
}

func TestReadTextInvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "invalid.txt")
	
	os.WriteFile(filename, []byte("invalid\nformat\nhere"), 0644)
	
	q := NewQueue()
	err := q.ReadText(filename)
	if err == nil {
		t.Error("Expected error for invalid text file")
	}
}

func TestReadTextEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.txt")
	
	os.WriteFile(filename, []byte(""), 0644)
	
	q := NewQueue()
	err := q.ReadText(filename)
	if err == nil {
		t.Error("Expected error for empty text file")
	}
}

func TestReadTextInvalidSize(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "invalid_size.txt")
	
	os.WriteFile(filename, []byte("invalid_size\n"), 0644)
	
	q := NewQueue()
	err := q.ReadText(filename)
	if err == nil {
		t.Error("Expected error for invalid size format")
	}
}

func TestDequeueUntilEmpty(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	
	for i := 0; i < 3; i++ {
		_, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
	}
	
	if q.Size() != 0 {
		t.Error("Size should be 0 after dequeuing all elements")
	}
	
	_, err := q.Dequeue()
	if err == nil {
		t.Error("Expected error when dequeuing from empty queue")
	}
}

func TestEnqueueAfterClear(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Clear()
	
	err := q.Enqueue("new")
	if err != nil {
		t.Fatalf("Enqueue after clear failed: %v", err)
	}
	
	if q.Size() != 1 {
		t.Error("Size should be 1 after enqueueing after clear")
	}
	if q.head.Data != "new" || q.tail.Data != "new" {
		t.Error("Head and tail should be 'new' after enqueueing after clear")
	}
}

func TestDelNonAdjacent(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c", "d", "e")
	q.Del("c")
	
	if q.Size() != 4 {
		t.Error("Size should be 4 after deleting middle element")
	}
	
	current := q.head
	expected := []string{"a", "b", "d", "e"}
	for i, exp := range expected {
		if current == nil || current.Data != exp {
			t.Errorf("Element %d should be %s, got %s", i, exp, current.Data)
		}
		current = current.Next
	}
}

func TestDelAllButOne(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c", "d", "e")
	
	q.Del("a")
	q.Del("b")
	q.Del("d")
	q.Del("e")
	
	if q.Size() != 1 {
		t.Error("Size should be 1 after deleting all but one")
	}
	if q.head.Data != "c" || q.tail.Data != "c" {
		t.Error("Head and tail should be 'c' after deleting all but one")
	}
}

func TestPrintEmpty(t *testing.T) {
	q := NewQueue()
	output := captureOutput(q.Print)
	expected := "Queue is empty\n"
	if output != expected {
		t.Errorf("Empty queue print output = %q, want %q", output, expected)
	}
}

func TestPrintSingleElement(t *testing.T) {
	q := NewQueueWithItems("single")
	output := captureOutput(q.Print)
	expected := "single \n"
	if output != expected {
		t.Errorf("Single element print output = %q, want %q", output, expected)
	}
}

func TestNodePrevNextConsistency(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c", "d")
	
	current := q.head
	for current != nil {
		if current.Next != nil && current.Next.Prev != current {
			t.Error("Next node's Prev should point back to current")
		}
		if current.Prev != nil && current.Prev.Next != current {
			t.Error("Prev node's Next should point to current")
		}
		current = current.Next
	}
}

func TestMaxSizeBoundary(t *testing.T) {
	q := NewQueue()
	q.maxSize = 3
	
	for i := 0; i < 3; i++ {
		err := q.Enqueue("item" + string(rune(i)))
		if err != nil {
			t.Fatalf("Enqueue %d failed: %v", i, err)
		}
	}
	
	if q.Size() != 3 {
		t.Error("Size should be 3 at max size")
	}
	
	err := q.Enqueue("overflow")
	if err == nil {
		t.Error("Expected error for enqueue at max size")
	}
}

func TestDequeueUntilEmptyThenEnqueue(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	
	for i := 0; i < 3; i++ {
		_, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
	}
	
	err := q.Enqueue("new")
	if err != nil {
		t.Fatalf("Enqueue after empty failed: %v", err)
	}
	
	if q.Size() != 1 {
		t.Error("Size should be 1 after enqueueing after empty")
	}
	if q.head.Data != "new" || q.tail.Data != "new" {
		t.Error("Head and tail should be 'new' after enqueueing after empty")
	}
}

func TestDelFromBeginning(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Del("a")
	
	if q.Size() != 2 {
		t.Error("Size should be 2 after deleting from beginning")
	}
	if q.head.Data != "b" {
		t.Error("Head should be 'b' after deleting from beginning")
	}
}

func TestDelFromEnd(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Del("c")
	
	if q.Size() != 2 {
		t.Error("Size should be 2 after deleting from end")
	}
	if q.tail.Data != "b" {
		t.Error("Tail should be 'b' after deleting from end")
	}
}

func TestMultipleClears(t *testing.T) {
	q := NewQueueWithItems("a", "b", "c")
	q.Clear()
	q.Clear()
	
	if q.Size() != 0 {
		t.Error("Size should be 0 after multiple clears")
	}
	if q.head != nil || q.tail != nil {
		t.Error("Head and tail should be nil after multiple clears")
	}
	
	err := q.Enqueue("new")
	if err != nil {
		t.Fatalf("Enqueue after multiple clears failed: %v", err)
	}
	
	if q.Size() != 1 {
		t.Error("Size should be 1 after enqueueing after multiple clears")
	}
}

func TestLargeStringElements(t *testing.T) {
	q := NewQueue()
	
	largeString := "a" + string(make([]byte, 1000)) + "z"
	err := q.Enqueue(largeString)
	if err != nil {
		t.Fatalf("Enqueue large string failed: %v", err)
	}
	
	if q.Size() != 1 {
		t.Error("Size should be 1 after enqueuing large string")
	}
	
	data, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue large string failed: %v", err)
	}
	
	if data != largeString {
		t.Error("Dequeued string should match original large string")
	}
	
	if q.Size() != 0 {
		t.Error("Size should be 0 after dequeuing large string")
	}
}

func TestWriteReadBinaryLargeQueue(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "large.bin")
	
	q := NewQueue()
	for i := 0; i < 100; i++ {
		q.Enqueue("item" + string(rune(i)))
	}
	
	q.WriteBinary(filename)
	
	q2 := NewQueue()
	q2.ReadBinary(filename)
	
	if q2.Size() != 100 {
		t.Error("Size should be 100 after reading large queue")
	}
}

func TestWriteReadTextLargeQueue(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "large.txt")
	
	q := NewQueue()
	for i := 0; i < 100; i++ {
		q.Enqueue("item" + string(rune(i)))
	}
	
	q.WriteText(filename)
	
	q2 := NewQueue()
	q2.ReadText(filename)
	
	if q2.Size() != 100 {
		t.Error("Size should be 100 after reading large queue")
	}
}