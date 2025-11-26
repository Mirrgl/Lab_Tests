package forwardlist

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"encoding/binary"
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
	t.Run("NewForwardList", func(t *testing.T) {
		tests := []struct {
			name   string
			items  []string
			size   int
		}{
			{"Empty list", []string{}, 0},
			{"Single item", []string{"a"}, 1},
			{"Multiple items", []string{"a", "b", "c"}, 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fl := NewForwardList(tt.items...)
				if fl.Size() != tt.size {
					t.Errorf("NewForwardList() size = %d, want %d", fl.Size(), tt.size)
				}

				current := fl.head
				for i, item := range tt.items {
					if current == nil {
						t.Fatalf("Element at %d is nil", i)
					}
					if current.key != item {
						t.Errorf("Element at %d = %s, want %s", i, current.key, item)
					}
					current = current.next
				}
				if current != nil {
					t.Errorf("Expected end of list, but got more elements")
				}
			})
		}
	})
}

func TestCoreOperations(t *testing.T) {
	t.Run("PushOperations", func(t *testing.T) {
		fl := NewForwardList()
		
		fl.PushBack("a")
		if fl.Size() != 1 {
			t.Errorf("Size after PushBack = %d, want 1", fl.Size())
		}
		if fl.head.key != "a" || fl.tail.key != "a" {
			t.Errorf("Head/tail after PushBack = %s/%s, want a/a", fl.head.key, fl.tail.key)
		}
		
		fl.PushBack("b")
		if fl.Size() != 2 {
			t.Errorf("Size after second PushBack = %d, want 2", fl.Size())
		}
		if fl.head.key != "a" || fl.tail.key != "b" {
			t.Errorf("Head/tail after second PushBack = %s/%s, want a/b", fl.head.key, fl.tail.key)
		}
		
		fl.PushFront("c")
		if fl.Size() != 3 {
			t.Errorf("Size after PushFront = %d, want 3", fl.Size())
		}
		if fl.head.key != "c" || fl.tail.key != "b" {
			t.Errorf("Head/tail after PushFront = %s/%s, want c/b", fl.head.key, fl.tail.key)
		}
		
		expected := []string{"c", "a", "b"}
		current := fl.head
		for i, exp := range expected {
			if current == nil {
				t.Fatalf("Element at %d is nil", i)
			}
			if current.key != exp {
				t.Errorf("Element at %d = %s, want %s", i, current.key, exp)
			}
			current = current.next
		}
	})

	t.Run("InsertOperations", func(t *testing.T) {
		tests := []struct {
			name     string
			setup    []string
			insert   func(*ForwardList) error
			expected []string
			size     int
		}{
			{
				"InsertBefore at start",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.InsertBefore("x", 0) },
				[]string{"x", "a", "b", "c"},
				4,
			},
			{
				"InsertBefore at middle",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.InsertBefore("x", 1) },
				[]string{"a", "x", "b", "c"},
				4,
			},
			{
				"InsertAfter at middle",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.InsertAfter("x", 1) },
				[]string{"a", "b", "x", "c"},
				4,
			},
			{
				"InsertAfter at end",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.InsertAfter("x", 2) },
				[]string{"a", "b", "c", "x"},
				4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fl := NewForwardList(tt.setup...)
				err := tt.insert(fl)
				if err != nil {
					t.Fatalf("Insert operation failed: %v", err)
				}

				if fl.Size() != tt.size {
					t.Errorf("Size = %d, want %d", fl.Size(), tt.size)
				}

				current := fl.head
				for i, exp := range tt.expected {
					if current == nil {
						t.Fatalf("Element at %d is nil", i)
					}
					if current.key != exp {
						t.Errorf("Element at %d = %s, want %s", i, current.key, exp)
					}
					current = current.next
				}
				if current != nil {
					t.Errorf("Expected end of list, but got more elements")
				}
			})
		}
	})

	t.Run("PopOperations", func(t *testing.T) {
		tests := []struct {
			name     string
			setup    []string
			popOp    func(*ForwardList) error
			expected []string
			size     int
		}{
			{
				"PopFront from multiple",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.PopFront() },
				[]string{"b", "c"},
				2,
			},
			{
				"PopBack from multiple",
				[]string{"a", "b", "c"},
				func(fl *ForwardList) error { return fl.PopBack() },
				[]string{"a", "b"},
				2,
			},
			{
				"PopFront from single",
				[]string{"a"},
				func(fl *ForwardList) error { return fl.PopFront() },
				[]string{},
				0,
			},
			{
				"PopBack from single",
				[]string{"a"},
				func(fl *ForwardList) error { return fl.PopBack() },
				[]string{},
				0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fl := NewForwardList(tt.setup...)
				err := tt.popOp(fl)
				if err != nil {
					t.Fatalf("Pop operation failed: %v", err)
				}

				if fl.Size() != tt.size {
					t.Errorf("Size = %d, want %d", fl.Size(), tt.size)
				}

				current := fl.head
				for i, exp := range tt.expected {
					if current == nil {
						t.Fatalf("Element at %d is nil", i)
					}
					if current.key != exp {
						t.Errorf("Element at %d = %s, want %s", i, current.key, exp)
					}
					current = current.next
				}
				if current != nil {
					t.Errorf("Expected end of list, but got more elements")
				}
			})
		}
	})

	t.Run("RemoveByValue", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c", "b", "d")

		found := fl.RemoveByValue("b")
		if !found {
			t.Error("RemoveByValue('b') returned false, want true")
		}
		if fl.Size() != 4 {
			t.Errorf("Size after remove = %d, want 4", fl.Size())
		}

		expected := []string{"a", "c", "b", "d"}
		current := fl.head
		for i, exp := range expected {
			if current == nil {
				t.Fatalf("Element at %d is nil", i)
			}
			if current.key != exp {
				t.Errorf("Element at %d = %s, want %s", i, current.key, exp)
			}
			current = current.next
		}

		found = fl.RemoveByValue("x")
		if found {
			t.Error("RemoveByValue('x') returned true, want false")
		}

		fl2 := NewForwardList("a", "b", "c")
		found = fl2.RemoveByValue("c")
		if !found {
			t.Error("RemoveByValue('c') returned false, want true")
		}
		if fl2.tail.key != "b" {
			t.Errorf("Tail after removing last element = %s, want 'b'", fl2.tail.key)
		}
	})

	t.Run("AccessOperations", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")

		t.Run("Front", func(t *testing.T) {
			val, err := fl.Front()
			if err != nil {
				t.Fatalf("Front() failed: %v", err)
			}
			if val != "a" {
				t.Errorf("Front() = %s, want 'a'", val)
			}
		})

		t.Run("Back", func(t *testing.T) {
			val, err := fl.Back()
			if err != nil {
				t.Fatalf("Back() failed: %v", err)
			}
			if val != "c" {
				t.Errorf("Back() = %s, want 'c'", val)
			}
		})

		t.Run("GetAt", func(t *testing.T) {
			tests := []struct {
				index int
				want  string
				err   bool
			}{
				{0, "a", false},
				{1, "b", false},
				{2, "c", false},
				{-1, "", true},
				{3, "", true},
			}

			for _, tt := range tests {
				t.Run(strconv.Itoa(tt.index), func(t *testing.T) {
					val, err := fl.GetAt(tt.index)
					if (err != nil) != tt.err {
						t.Errorf("GetAt(%d) error = %v, wantErr %v", tt.index, err, tt.err)
					}
					if !tt.err && val != tt.want {
						t.Errorf("GetAt(%d) = %s, want %s", tt.index, val, tt.want)
					}
				})
			}
		})
	})

	t.Run("UtilityMethods", func(t *testing.T) {
		empty := NewForwardList()
		single := NewForwardList("a")
		multiple := NewForwardList("a", "b", "c")

		tests := []struct {
			name string
			fl   *ForwardList
			size int
			isEmpty bool
		}{
			{"Empty list", empty, 0, true},
			{"Single element", single, 1, false},
			{"Multiple elements", multiple, 3, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.fl.Size() != tt.size {
					t.Errorf("Size() = %d, want %d", tt.fl.Size(), tt.size)
				}
				if tt.fl.IsEmpty() != tt.isEmpty {
					t.Errorf("IsEmpty() = %v, want %v", tt.fl.IsEmpty(), tt.isEmpty)
				}
			})
		}
	})

	t.Run("Clear", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")
		fl.Clear()
		if fl.Size() != 0 {
			t.Errorf("Size after Clear = %d, want 0", fl.Size())
		}
		if !fl.IsEmpty() {
			t.Error("IsEmpty after Clear = false, want true")
		}
		if fl.head != nil || fl.tail != nil {
			t.Error("Head/tail after Clear != nil, want nil")
		}
	})
}

func TestIndexValidation(t *testing.T) {
	fl := NewForwardList("a", "b", "c")

	t.Run("validatePosition", func(t *testing.T) {
		tests := []struct {
			position int
			allowEnd bool
			wantErr  bool
		}{
			{-1, false, true},
			{0, false, false},
			{1, false, false},
			{2, false, false},
			{3, false, true},
			{3, true, false},
			{4, true, true},
		}

		for _, tt := range tests {
			name := strconv.Itoa(tt.position)
			if tt.allowEnd {
				name += "_allowEnd"
			}
			t.Run(name, func(t *testing.T) {
				err := fl.validatePosition(tt.position, tt.allowEnd)
				if (err != nil) != tt.wantErr {
					t.Errorf("validatePosition(%d, %v) error = %v, wantErr %v", tt.position, tt.allowEnd, err, tt.wantErr)
				}
			})
		}
	})

	t.Run("getNodeAt", func(t *testing.T) {
		tests := []struct {
			position int
			want     string
			err      bool
		}{
			{0, "a", false},
			{1, "b", false},
			{2, "c", false},
			{-1, "", true},
			{3, "", true},
		}

		for _, tt := range tests {
			t.Run(strconv.Itoa(tt.position), func(t *testing.T) {
				node, err := fl.getNodeAt(tt.position)
				if (err != nil) != tt.err {
					t.Errorf("getNodeAt(%d) error = %v, wantErr %v", tt.position, err, tt.err)
				}
				if !tt.err && node.key != tt.want {
					t.Errorf("getNodeAt(%d) = %s, want %s", tt.position, node.key, tt.want)
				}
			})
		}
	})
}

func TestFileOperations(t *testing.T) {
	t.Run("BinaryFileOperations", func(t *testing.T) {
		tempDir := t.TempDir()
		binFile := filepath.Join(tempDir, "test.bin")

		t.Run("Valid round trip", func(t *testing.T) {
			original := NewForwardList("hello", "world", "42", "", "with\nnewline")
			if err := original.WriteBinary(binFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			readList := NewForwardList()
			if err := readList.ReadBinary(binFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}

			if readList.Size() != original.Size() {
				t.Errorf("Read size = %d, want %d", readList.Size(), original.Size())
			}

			origCurrent := original.head
			readCurrent := readList.head
			for i := 0; i < original.Size(); i++ {
				if origCurrent == nil || readCurrent == nil {
					t.Fatalf("Element at %d is nil", i)
				}
				if origCurrent.key != readCurrent.key {
					t.Errorf("Element at %d = %q, want %q", i, readCurrent.key, origCurrent.key)
				}
				origCurrent = origCurrent.next
				readCurrent = readCurrent.next
			}
		})

		t.Run("Empty list", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.bin")
			emptyList := NewForwardList()
			if err := emptyList.WriteBinary(emptyFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			readList := NewForwardList()
			if err := readList.ReadBinary(emptyFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}
			if readList.Size() != 0 {
				t.Errorf("Read size = %d, want 0", readList.Size())
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.bin")
			fl := NewForwardList()
			if err := fl.ReadBinary(nonExistent); err == nil {
				t.Error("ReadBinary() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.bin")
			file, _ := os.Create(invalidFile)
			file.Write([]byte{0x01, 0x02, 0x03})
			file.Close()

			if err := fl.ReadBinary(invalidFile); err == nil {
				t.Error("ReadBinary() expected error for invalid file, got nil")
			}
		})

		t.Run("Corrupted file handling", func(t *testing.T) {
			corruptedFile := filepath.Join(tempDir, "corrupted.bin")
			file, _ := os.Create(corruptedFile)
			binary.Write(file, binary.LittleEndian, uint64(2))
			binary.Write(file, binary.LittleEndian, uint64(3))
			file.Write([]byte("abc"))
			file.Close()

			fl := NewForwardList()
			if err := fl.ReadBinary(corruptedFile); err == nil {
				t.Error("ReadBinary() expected error for corrupted file, got nil")
			}
		})
	})

	t.Run("TextFileOperations", func(t *testing.T) {
		tempDir := t.TempDir()
		txtFile := filepath.Join(tempDir, "test.txt")

		t.Run("Valid round trip", func(t *testing.T) {
			original := NewForwardList("hello", "world", "42", "", "with space")
			if err := original.WriteText(txtFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewForwardList()
			if err := readList.ReadText(txtFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if readList.Size() != original.Size() {
				t.Errorf("Read size = %d, want %d", readList.Size(), original.Size())
			}

			origCurrent := original.head
			readCurrent := readList.head
			for i := 0; i < original.Size(); i++ {
				if origCurrent == nil || readCurrent == nil {
					t.Fatalf("Element at %d is nil", i)
				}
				if origCurrent.key != readCurrent.key {
					t.Errorf("Element at %d = %q, want %q", i, readCurrent.key, origCurrent.key)
				}
				origCurrent = origCurrent.next
				readCurrent = readCurrent.next
			}
		})

		t.Run("Empty list", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.txt")
			emptyList := NewForwardList()
			if err := emptyList.WriteText(emptyFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewForwardList()
			if err := readList.ReadText(emptyFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}
			if readList.Size() != 0 {
				t.Errorf("Read size = %d, want 0", readList.Size())
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.txt")
			fl := NewForwardList()
			if err := fl.ReadText(nonExistent); err == nil {
				t.Error("ReadText() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.txt")
			os.WriteFile(invalidFile, []byte("invalid\nformat"), 0644)

			if err := fl.ReadText(invalidFile); err == nil {
				t.Error("ReadText() expected error for invalid file, got nil")
			}

			invalidSizeFile := filepath.Join(tempDir, "invalid_size.txt")
			os.WriteFile(invalidSizeFile, []byte("invalid_size\n"), 0644)

			if err := fl.ReadText(invalidSizeFile); err == nil {
				t.Error("ReadText() expected error for invalid size format, got nil")
			}
		})

		t.Run("Edge cases", func(t *testing.T) {
			edgeFile := filepath.Join(tempDir, "edge.txt")
			edgeList := NewForwardList("", " ")
			if err := edgeList.WriteText(edgeFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewForwardList()
			if err := readList.ReadText(edgeFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if readList.Size() != 2 {
				t.Errorf("Read size = %d, want 2", readList.Size())
			}
			expected := []string{"", " "}
			current := readList.head
			for i, exp := range expected {
				if current == nil {
					t.Fatalf("Element at %d is nil", i)
				}
				if current.key != exp {
					t.Errorf("Element at %d = %q, want %q", i, current.key, exp)
				}
				current = current.next
			}
		})
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("EmptyListOperations", func(t *testing.T) {
		fl := NewForwardList()

		t.Run("Pop operations", func(t *testing.T) {
			err := fl.PopFront()
			if err == nil {
				t.Error("PopFront() on empty list expected error, got nil")
			}

			err = fl.PopBack()
			if err == nil {
				t.Error("PopBack() on empty list expected error, got nil")
			}
		})

		t.Run("Front/Back access", func(t *testing.T) {
			_, err := fl.Front()
			if err == nil {
				t.Error("Front() on empty list expected error, got nil")
			}

			_, err = fl.Back()
			if err == nil {
				t.Error("Back() on empty list expected error, got nil")
			}
		})

		t.Run("RemoveByValue", func(t *testing.T) {
			found := fl.RemoveByValue("x")
			if found {
				t.Error("RemoveByValue('x') on empty list returned true, want false")
			}
		})
	})

	t.Run("SingleNodeConsistency", func(t *testing.T) {
		fl := NewForwardList("a")
		
		if fl.head != fl.tail {
			t.Error("Head and tail should be the same for single node")
		}
		if fl.head.next != nil {
			t.Error("Head.next should be nil for single node")
		}
		
		fl.PopFront()
		if fl.head != nil || fl.tail != nil {
			t.Error("Head and tail should be nil after deleting only node")
		}
		if fl.Size() != 0 {
			t.Errorf("Size should be 0 after deleting only node, got %d", fl.Size())
		}
	})

	t.Run("InsertValidation", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")
		
		err := fl.InsertBefore("x", -1)
		if err == nil {
			t.Error("InsertBefore(-1) expected error, got nil")
		}
		err = fl.InsertBefore("x", 4)
		if err == nil {
			t.Error("InsertBefore(4) expected error, got nil")
		}
		
		err = fl.InsertAfter("x", -1)
		if err == nil {
			t.Error("InsertAfter(-1) expected error, got nil")
		}
		err = fl.InsertAfter("x", 3)
		if err == nil {
			t.Error("InsertAfter(3) expected error, got nil")
		}
	})

	t.Run("LargeListOperations", func(t *testing.T) {
		fl := NewForwardList()
		for i := 0; i < 100; i++ {
			fl.PushBack(strconv.Itoa(i))
		}
		
		if fl.Size() != 100 {
			t.Errorf("Size of large list = %d, want 100", fl.Size())
		}
		
		val, err := fl.GetAt(50)
		if err != nil {
			t.Fatalf("GetAt(50) failed: %v", err)
		}
		if val != "50" {
			t.Errorf("GetAt(50) = %s, want '50'", val)
		}
		
		fl.RemoveByValue("25")
		fl.RemoveByValue("75")
		if fl.Size() != 98 {
			t.Errorf("Size after removals = %d, want 98", fl.Size())
		}
	})
}

func TestOutput(t *testing.T) {
	t.Run("Print", func(t *testing.T) {
		tests := []struct {
			name     string
			elements []string
			want     string
		}{
			{"Empty list", []string{}, "List is empty\n"},
			{"Single element", []string{"a"}, "a \n"},
			{"Multiple elements", []string{"a", "b", "c"}, "a b c \n"},
			{"With spaces", []string{"hello world", "test"}, "hello world test \n"},
			{"Empty strings", []string{"", "", ""}, "   \n"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fl := NewForwardList(tt.elements...)
				output := captureOutput(fl.Print)
				if output != tt.want {
					t.Errorf("Print() = %q, want %q", output, tt.want)
				}
			})
		}
	})
}

func TestErrorConditions(t *testing.T) {
	t.Run("InsertAfter at last position", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")
		err := fl.InsertAfter("x", 2)
		if err != nil {
			t.Fatalf("InsertAfter at last position failed: %v", err)
		}
		if fl.Size() != 4 {
			t.Errorf("Size after InsertAfter at last = %d, want 4", fl.Size())
		}
		current := fl.head
		for i := 0; i < 3; i++ {
			current = current.next
		}
		if current.key != "x" {
			t.Errorf("Last element after InsertAfter = %s, want 'x'", current.key)
		}
	})

	t.Run("InsertAfter beyond last position", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")
		err := fl.InsertAfter("x", 3)
		if err == nil {
			t.Error("InsertAfter beyond last position expected error, got nil")
		}
	})

	t.Run("InsertBefore beyond last position", func(t *testing.T) {
		fl := NewForwardList("a", "b", "c")
		err := fl.InsertBefore("x", 3)
		if err == nil {
			t.Error("InsertBefore beyond last position expected error, got nil")
		}
	})

	t.Run("InsertBefore at start", func(t *testing.T) {
		fl := NewForwardList("b", "c")
		err := fl.InsertBefore("a", 0)
		if err != nil {
			t.Fatalf("InsertBefore at start failed: %v", err)
		}
		if fl.Size() != 3 {
			t.Errorf("Size after InsertBefore at start = %d, want 3", fl.Size())
		}
		if fl.head.key != "a" {
			t.Errorf("Head after InsertBefore at start = %s, want 'a'", fl.head.key)
		}
	})
}