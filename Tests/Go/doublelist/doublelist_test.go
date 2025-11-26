package doublelist

import (
	"bytes"
	"os"
	"path/filepath"
	"encoding/binary"
	"strconv"
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
	t.Run("NewDoubleList", func(t *testing.T) {
		tests := []struct {
			name   string
			items  []string
			length int
		}{
			{"Empty list", []string{}, 0},
			{"Single item", []string{"a"}, 1},
			{"Multiple items", []string{"a", "b", "c"}, 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dl := NewDoubleList(tt.items...)
				if dl.GetLength() != tt.length {
					t.Errorf("NewDoubleList() length = %d, want %d", dl.GetLength(), tt.length)
				}

				for i, item := range tt.items {
					val, err := dl.GetElement(i)
					if err != nil {
						t.Fatalf("GetElement(%d) failed: %v", i, err)
					}
					if val != item {
						t.Errorf("Element at %d = %s, want %s", i, val, item)
					}
				}
			})
		}
	})
}

func TestCoreOperations(t *testing.T) {
	t.Run("AddOperations", func(t *testing.T) {
		tests := []struct {
			name     string
			ops      []func(dl *DoubleList) error
			expected []string
			length   int
		}{
			{
				"AddHead sequence",
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.AddHead("c") },
					func(dl *DoubleList) error { return dl.AddHead("b") },
					func(dl *DoubleList) error { return dl.AddHead("a") },
				},
				[]string{"a", "b", "c"},
				3,
			},
			{
				"AddTail sequence",
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.AddTail("a") },
					func(dl *DoubleList) error { return dl.AddTail("b") },
					func(dl *DoubleList) error { return dl.AddTail("c") },
				},
				[]string{"a", "b", "c"},
				3,
			},
			{
				"AddAfter operations",
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.AddTail("a") },
					func(dl *DoubleList) error { return dl.AddAfter("b", 0) },
					func(dl *DoubleList) error { return dl.AddAfter("c", 1) },
				},
				[]string{"a", "b", "c"},
				3,
			},
			{
				"AddBefore operations",
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.AddTail("c") },
					func(dl *DoubleList) error { return dl.AddBefore("b", 0) },
					func(dl *DoubleList) error { return dl.AddBefore("a", 0) },
				},
				[]string{"a", "b", "c"},
				3,
			},
			{
				"Mixed operations",
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.AddHead("b") },
					func(dl *DoubleList) error { return dl.AddBefore("a", 0) },
					func(dl *DoubleList) error { return dl.AddAfter("c", 1) },
				},
				[]string{"a", "b", "c"},
				3,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dl := NewDoubleList()
				for i, op := range tt.ops {
					if err := op(dl); err != nil {
						t.Fatalf("Operation %d failed: %v", i, err)
					}
				}

				if dl.GetLength() != tt.length {
					t.Errorf("Length = %d, want %d", dl.GetLength(), tt.length)
				}

				for i, expected := range tt.expected {
					val, err := dl.GetElement(i)
					if err != nil {
						t.Fatalf("GetElement(%d) failed: %v", i, err)
					}
					if val != expected {
						t.Errorf("Element at %d = %s, want %s", i, val, expected)
					}
				}
			})
		}
	})

	t.Run("DeleteOperations", func(t *testing.T) {
		tests := []struct {
			name     string
			setup    func() *DoubleList
			ops      []func(dl *DoubleList) error
			expected []string
			length   int
		}{
			{
				"DeleteHead from multiple",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteHead() },
				},
				[]string{"b", "c"},
				2,
			},
			{
				"DeleteTail from multiple",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteTail() },
				},
				[]string{"a", "b"},
				2,
			},
			{
				"DeleteAt middle",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteAt(1) },
				},
				[]string{"a", "c"},
				2,
			},
			{
				"DeleteAt first",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteAt(0) },
				},
				[]string{"b", "c"},
				2,
			},
			{
				"DeleteAt last",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteAt(2) },
				},
				[]string{"a", "b"},
				2,
			},
			{
				"DeleteByValue middle",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteByValue("b") },
				},
				[]string{"a", "c"},
				2,
			},
			{
				"DeleteByValue first",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteByValue("a") },
				},
				[]string{"b", "c"},
				2,
			},
			{
				"DeleteByValue last",
				func() *DoubleList {
					return NewDoubleList("a", "b", "c")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteByValue("c") },
				},
				[]string{"a", "b"},
				2,
			},
			{
				"DeleteHead from single",
				func() *DoubleList {
					return NewDoubleList("a")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteHead() },
				},
				[]string{},
				0,
			},
			{
				"DeleteTail from single",
				func() *DoubleList {
					return NewDoubleList("a")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteTail() },
				},
				[]string{},
				0,
			},
			{
				"DeleteByValue from single",
				func() *DoubleList {
					return NewDoubleList("a")
				},
				[]func(dl *DoubleList) error{
					func(dl *DoubleList) error { return dl.DeleteByValue("a") },
				},
				[]string{},
				0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dl := tt.setup()
				for i, op := range tt.ops {
					if err := op(dl); err != nil {
						t.Fatalf("Operation %d failed: %v", i, err)
					}
				}

				if dl.GetLength() != tt.length {
					t.Errorf("Length = %d, want %d", dl.GetLength(), tt.length)
				}

				for i, expected := range tt.expected {
					val, err := dl.GetElement(i)
					if err != nil {
						t.Fatalf("GetElement(%d) failed: %v", i, err)
					}
					if val != expected {
						t.Errorf("Element at %d = %s, want %s", i, val, expected)
					}
				}

				if tt.length > 0 {
					_, err := dl.GetElement(0)
					if err != nil {
						t.Error("Head element not accessible")
					}
					_, err = dl.GetElement(tt.length - 1)
					if err != nil {
						t.Error("Tail element not accessible")
					}
				}
			})
		}
	})

	t.Run("AccessOperations", func(t *testing.T) {
		dl := NewDoubleList("a", "b", "c", "b", "d")

		t.Run("GetElement", func(t *testing.T) {
			tests := []struct {
				index int
				want  string
				err   bool
			}{
				{0, "a", false},
				{1, "b", false},
				{2, "c", false},
				{3, "b", false},
				{4, "d", false},
				{-1, "", true},
				{5, "", true},
			}

			for _, tt := range tests {
				t.Run(strconv.Itoa(tt.index), func(t *testing.T) {
					val, err := dl.GetElement(tt.index)
					if (err != nil) != tt.err {
						t.Errorf("GetElement(%d) error = %v, wantErr %v", tt.index, err, tt.err)
					}
					if !tt.err && val != tt.want {
						t.Errorf("GetElement(%d) = %s, want %s", tt.index, val, tt.want)
					}
				})
			}
		})

		t.Run("PopElement", func(t *testing.T) {
			tests := []struct {
				name    string
				setup   []string
				index   int
				want    string
				err     bool
				finalLen int
			}{
				{"Pop middle", []string{"a", "b", "c", "b", "d"}, 1, "b", false, 4},
				{"Pop another middle", []string{"a", "c", "b", "d"}, 2, "b", false, 3},
				{"Pop first", []string{"a", "c", "d"}, 0, "a", false, 2},
				{"Pop last", []string{"c", "d"}, 1, "d", false, 1},
				{"Pop only element", []string{"c"}, 0, "c", false, 0},
				{"Pop from empty", []string{}, 0, "", true, 0},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					testDL := NewDoubleList(tt.setup...)
					val, err := testDL.PopElement(tt.index)
					if (err != nil) != tt.err {
						t.Errorf("PopElement(%d) error = %v, wantErr %v", tt.index, err, tt.err)
					}
					if !tt.err && val != tt.want {
						t.Errorf("PopElement(%d) = %s, want %s", tt.index, val, tt.want)
					}
					if testDL.GetLength() != tt.finalLen {
						t.Errorf("Length after pop = %d, want %d", testDL.GetLength(), tt.finalLen)
					}
				})
			}
		})

		t.Run("FindByValue", func(t *testing.T) {
			tests := []struct {
				key  string
				want string
				found bool
			}{
				{"a", "a", true},
				{"b", "b", true},
				{"c", "c", true},
				{"d", "d", true},
				{"x", "", false},
				{"", "", false},
			}

			for _, tt := range tests {
				t.Run(tt.key, func(t *testing.T) {
					node := dl.FindByValue(tt.key)
					if tt.found {
						if node == nil {
							t.Error("FindByValue() = nil, want node")
						} else if node.key != tt.want {
							t.Errorf("FindByValue() = %s, want %s", node.key, tt.want)
						}
					} else {
						if node != nil {
							t.Errorf("FindByValue() = %v, want nil", node)
						}
					}
				})
			}
		})
	})

	t.Run("IndexValidation", func(t *testing.T) {
		dl := NewDoubleList("a", "b", "c")

		t.Run("validateIndex", func(t *testing.T) {
			tests := []struct {
				index   int
				allowEnd bool
				wantErr bool
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
				name := strconv.Itoa(tt.index)
				if tt.allowEnd {
					name += "_allowEnd"
				}
				t.Run(name, func(t *testing.T) {
					err := dl.validateIndex(tt.index, tt.allowEnd)
					if (err != nil) != tt.wantErr {
						t.Errorf("validateIndex(%d, %v) error = %v, wantErr %v", tt.index, tt.allowEnd, err, tt.wantErr)
					}
				})
			}
		})
	})

	t.Run("UtilityMethods", func(t *testing.T) {
		empty := NewDoubleList()
		single := NewDoubleList("a")
		multiple := NewDoubleList("a", "b", "c")

		tests := []struct {
			name string
			dl   *DoubleList
			len  int
			isEmpty bool
		}{
			{"Empty list", empty, 0, true},
			{"Single element", single, 1, false},
			{"Multiple elements", multiple, 3, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.dl.GetLength() != tt.len {
					t.Errorf("GetLength() = %d, want %d", tt.dl.GetLength(), tt.len)
				}
				if tt.dl.IsEmpty() != tt.isEmpty {
					t.Errorf("IsEmpty() = %v, want %v", tt.dl.IsEmpty(), tt.isEmpty)
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
			original := NewDoubleList("hello", "world", "42", "", "with\nnewline")
			if err := original.WriteBinary(binFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			readList := NewDoubleList()
			if err := readList.ReadBinary(binFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}

			if readList.GetLength() != original.GetLength() {
				t.Errorf("Read length = %d, want %d", readList.GetLength(), original.GetLength())
			}

			for i := 0; i < original.GetLength(); i++ {
				origVal, _ := original.GetElement(i)
				readVal, _ := readList.GetElement(i)
				if origVal != readVal {
					t.Errorf("Element at %d = %q, want %q", i, readVal, origVal)
				}
			}
		})

		t.Run("Empty list", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.bin")
			emptyList := NewDoubleList()
			if err := emptyList.WriteBinary(emptyFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			readList := NewDoubleList()
			if err := readList.ReadBinary(emptyFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}
			if readList.GetLength() != 0 {
				t.Errorf("Read length = %d, want 0", readList.GetLength())
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.bin")
			dl := NewDoubleList()
			if err := dl.ReadBinary(nonExistent); err == nil {
				t.Error("ReadBinary() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.bin")
			file, _ := os.Create(invalidFile)
			file.Write([]byte{0x01, 0x02, 0x03})
			file.Close()

			if err := dl.ReadBinary(invalidFile); err == nil {
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

			dl := NewDoubleList()
			if err := dl.ReadBinary(corruptedFile); err == nil {
				t.Error("ReadBinary() expected error for corrupted file, got nil")
			}
		})
	})

	t.Run("TextFileOperations", func(t *testing.T) {
		tempDir := t.TempDir()
		txtFile := filepath.Join(tempDir, "test.txt")

		t.Run("Valid round trip", func(t *testing.T) {
			original := NewDoubleList("hello", "world", "42", "", "with space")
			if err := original.WriteText(txtFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewDoubleList()
			if err := readList.ReadText(txtFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if readList.GetLength() != original.GetLength() {
				t.Errorf("Read length = %d, want %d", readList.GetLength(), original.GetLength())
			}

			for i := 0; i < original.GetLength(); i++ {
				origVal, _ := original.GetElement(i)
				readVal, _ := readList.GetElement(i)
				if origVal != readVal {
					t.Errorf("Element at %d = %q, want %q", i, readVal, origVal)
				}
			}
		})

		t.Run("Empty list", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.txt")
			emptyList := NewDoubleList()
			if err := emptyList.WriteText(emptyFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewDoubleList()
			if err := readList.ReadText(emptyFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}
			if readList.GetLength() != 0 {
				t.Errorf("Read length = %d, want 0", readList.GetLength())
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.txt")
			dl := NewDoubleList()
			if err := dl.ReadText(nonExistent); err == nil {
				t.Error("ReadText() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.txt")
			os.WriteFile(invalidFile, []byte("invalid\nformat"), 0644)

			if err := dl.ReadText(invalidFile); err == nil {
				t.Error("ReadText() expected error for invalid file, got nil")
			}

			invalidLenFile := filepath.Join(tempDir, "invalid_len.txt")
			os.WriteFile(invalidLenFile, []byte("invalid_length\n"), 0644)

			if err := dl.ReadText(invalidLenFile); err == nil {
				t.Error("ReadText() expected error for invalid length format, got nil")
			}
		})

		t.Run("Edge cases", func(t *testing.T) {
			edgeFile := filepath.Join(tempDir, "edge.txt")
			edgeList := NewDoubleList("", " ")
			if err := edgeList.WriteText(edgeFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readList := NewDoubleList()
			if err := readList.ReadText(edgeFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if readList.GetLength() != 2 {
				t.Errorf("Read length = %d, want 2", readList.GetLength())
			}
			expected := []string{"", " "}
			for i := 0; i < 2; i++ {
				val, _ := readList.GetElement(i)
				if val != expected[i] {
					t.Errorf("Element at %d = %q, want %q", i, val, expected[i])
				}
			}
		})
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("EmptyListOperations", func(t *testing.T) {
		dl := NewDoubleList()

		tests := []struct {
			name string
			op   func() error
		}{
			{"GetElement(0)", func() error { _, err := dl.GetElement(0); return err }},
			{"DeleteAt(0)", func() error { return dl.DeleteAt(0) }},
			{"PopElement(0)", func() error { _, err := dl.PopElement(0); return err }},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := tt.op(); err == nil {
					t.Error("Expected error for empty list operation, got nil")
				}
			})
		}

		t.Run("DeleteHead()", func(t *testing.T) {
			if err := dl.DeleteHead(); err == nil {
				t.Error("Expected error for DeleteHead() on empty list, got nil")
			}
		})
		t.Run("DeleteTail()", func(t *testing.T) {
			if err := dl.DeleteTail(); err == nil {
				t.Error("Expected error for DeleteTail() on empty list, got nil")
			}
		})
	})

	t.Run("SingleNodeConsistency", func(t *testing.T) {
		dl := NewDoubleList("a")
		
		if dl.head != dl.tail {
			t.Error("Head and tail should be the same for single node")
		}
		if dl.head.prev != nil {
			t.Error("Head.prev should be nil for single node")
		}
		if dl.tail.next != nil {
			t.Error("Tail.next should be nil for single node")
		}
		
		dl.DeleteHead()
		if dl.head != nil || dl.tail != nil {
			t.Error("Head and tail should be nil after deleting only node")
		}
		if dl.GetLength() != 0 {
			t.Errorf("Length should be 0 after deleting only node, got %d", dl.GetLength())
		}
	})

	t.Run("IndexEdgeCases", func(t *testing.T) {
		dl := NewDoubleList("a", "b", "c")
		
		tests := []struct {
			name string
			op   func() error
		}{
			{"GetElement(-1)", func() error { _, err := dl.GetElement(-1); return err }},
			{"GetElement(3)", func() error { _, err := dl.GetElement(3); return err }},
			{"DeleteAt(-1)", func() error { return dl.DeleteAt(-1) }},
			{"DeleteAt(3)", func() error { return dl.DeleteAt(3) }},
			{"AddAfter(-1)", func() error { return dl.AddAfter("x", -1) }},
			{"AddAfter(3)", func() error { return dl.AddAfter("x", 3) }},
			{"AddBefore(-1)", func() error { return dl.AddBefore("x", -1) }},
			{"AddBefore(4)", func() error { return dl.AddBefore("x", 4) }},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := tt.op(); err == nil {
					t.Error("Expected error for invalid index, got nil")
				}
			})
		}
	})

	t.Run("DuplicateValues", func(t *testing.T) {
		dl := NewDoubleList("a", "b", "a", "c", "a")
		
		if err := dl.DeleteByValue("a"); err != nil {
			t.Fatalf("DeleteByValue failed: %v", err)
		}
		if dl.GetLength() != 4 {
			t.Errorf("Length after delete = %d, want 4", dl.GetLength())
		}
		val, _ := dl.GetElement(0)
		if val != "b" {
			t.Errorf("First element = %s, want 'b'", val)
		}
		
		if err := dl.DeleteByValue("a"); err != nil {
			t.Fatalf("DeleteByValue failed: %v", err)
		}
		if dl.GetLength() != 3 {
			t.Errorf("Length after delete = %d, want 3", dl.GetLength())
		}
		val, _ = dl.GetElement(1)
		if val != "c" {
			t.Errorf("Element at 1 = %s, want 'c'", val)
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
			{"Empty list", []string{}, "Список пуст\n"},
			{"Single element", []string{"a"}, "a\n"},
			{"Multiple elements", []string{"a", "b", "c"}, "a b c\n"},
			{"With spaces", []string{"hello world", "test"}, "hello world test\n"},
			{"Empty strings", []string{"", "", ""}, "  \n"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dl := NewDoubleList(tt.elements...)
				output := captureOutput(dl.Print)
				if output != tt.want {
					t.Errorf("Print() = %q, want %q", output, tt.want)
				}
			})
		}
	})
}