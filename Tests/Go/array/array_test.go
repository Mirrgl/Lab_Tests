package array

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
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
	t.Run("NewArray", func(t *testing.T) {
		tests := []struct {
			name    string
			size    int
			wantErr bool
		}{
			{"Valid size", 5, false},
			{"Minimum size", 1, false},
			{"Zero size", 0, true},
			{"Negative size", -1, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, err := NewArray(tt.size)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewArray() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err != nil {
					return
				}
				if a.len != 0 {
					t.Errorf("NewArray() len = %d, want 0", a.len)
				}
				if a.cap != tt.size {
					t.Errorf("NewArray() cap = %d, want %d", a.cap, tt.size)
				}
				if len(a.data) != tt.size {
					t.Errorf("NewArray() data length = %d, want %d", len(a.data), tt.size)
				}
			})
		}
	})

	t.Run("NewArrayFromList", func(t *testing.T) {
		tests := []struct {
			name   string
			items  []string
			want   []string
			length int
			cap    int
		}{
			{"Empty list", []string{}, []string{}, 0, 1},
			{"Single item", []string{"a"}, []string{"a"}, 1, 1},
			{"Multiple items", []string{"a", "b", "c"}, []string{"a", "b", "c"}, 3, 3},
			{"With empty strings", []string{"", "", ""}, []string{"", "", ""}, 3, 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, err := NewArrayFromList(tt.items)
				if err != nil {
					t.Fatalf("NewArrayFromList() unexpected error: %v", err)
				}
				if a.len != tt.length {
					t.Errorf("NewArrayFromList() length = %d, want %d", a.len, tt.length)
				}
				if a.cap != tt.cap {
					t.Errorf("NewArrayFromList() capacity = %d, want %d", a.cap, tt.cap)
				}
				for i := 0; i < tt.length; i++ {
					if a.data[i] != tt.want[i] {
						t.Errorf("NewArrayFromList() element at %d = %s, want %s", i, a.data[i], tt.want[i])
					}
				}
			})
		}
	})
}

func TestCoreOperations(t *testing.T) {
	t.Run("GetSetElement", func(t *testing.T) {
		a, _ := NewArrayFromList([]string{"a", "b", "c"})

		tests := []struct {
			name    string
			index   int
			setVal  string
			getVal  string
			wantErr bool
		}{
			{"Valid index", 1, "x", "x", false},
			{"First index", 0, "start", "start", false},
			{"Last index", 2, "end", "end", false},
			{"Negative index", -1, "", "", true},
			{"Index equal to length", 3, "", "", true},
			{"Index beyond length", 4, "", "", true},
			{"Empty array", 0, "test", "", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.name == "Empty array" {
					a, _ = NewArray(1)
				}

				err := a.SetElement(tt.setVal, tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetElement() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					return
				}
				val, err := a.GetElement(tt.index)
				if err != nil {
					t.Fatalf("GetElement() unexpected error: %v", err)
				}
				if val != tt.getVal {
					t.Errorf("GetElement() = %s, want %s", val, tt.getVal)
				}
			})
		}
	})

	t.Run("DeleteElement", func(t *testing.T) {
		tests := []struct {
			name    string
			initial []string
			index   int
			want    []string
			wantErr bool
		}{
			{"Middle element", []string{"a", "b", "c"}, 1, []string{"a", "c"}, false},
			{"First element", []string{"a", "b", "c"}, 0, []string{"b", "c"}, false},
			{"Last element", []string{"a", "b", "c"}, 2, []string{"a", "b"}, false},
			{"Single element", []string{"a"}, 0, []string{}, false},
			{"Empty array", []string{}, 0, nil, true},
			{"Negative index", []string{"a", "b"}, -1, nil, true},
			{"Index equal to length", []string{"a", "b"}, 2, nil, true},
			{"Index beyond length", []string{"a", "b"}, 3, nil, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, _ := NewArrayFromList(tt.initial)
				err := a.DeleteElement(tt.index)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteElement() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}
				if a.len != len(tt.want) {
					t.Errorf("DeleteElement() length = %d, want %d", a.len, len(tt.want))
				}
				for i := 0; i < a.len; i++ {
					if a.data[i] != tt.want[i] {
						t.Errorf("DeleteElement() element at %d = %s, want %s", i, a.data[i], tt.want[i])
					}
				}
			})
		}
	})

	t.Run("AddElements", func(t *testing.T) {
		tests := []struct {
			name       string
			initial    []string
			addVal     string
			addIndex   int
			want       []string
			wantErr    bool
			expectGrow bool
		}{
			{"Add at end", []string{"a", "b"}, "c", -1, []string{"a", "b", "c"}, false, false},
			{"Add at start", []string{"a", "b"}, "x", 0, []string{"x", "a", "b"}, false, false},
			{"Add in middle", []string{"a", "c"}, "b", 1, []string{"a", "b", "c"}, false, false},
			{"Add to empty", []string{}, "a", -1, []string{"a"}, false, false},
			{"Index out of bounds low", []string{"a"}, "b", -2, nil, true, false},
			{"Index out of bounds high", []string{"a"}, "b", 2, nil, true, false},
			{"Grow on add end", []string{"a"}, "b", -1, []string{"a", "b"}, false, true},
			{"Grow on add index", []string{"a"}, "x", 0, []string{"x", "a"}, false, true},
			{"Grow multiple times", []string{"a"}, "b", -1, []string{"a", "b"}, false, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, _ := NewArray(1)
				for _, item := range tt.initial {
					a.AddElementEnd(item)
				}
				initialCap := a.cap

				var err error
				if tt.addIndex == -1 {
					a.AddElementEnd(tt.addVal)
				} else {
					err = a.AddElementAtIndex(tt.addVal, tt.addIndex)
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("AddElement error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					return
				}

				if tt.expectGrow && a.cap <= initialCap {
					t.Errorf("Expected growth but capacity remained %d", a.cap)
				}

				if a.len != len(tt.want) {
					t.Errorf("AddElement() length = %d, want %d", a.len, len(tt.want))
					return
				}

				for i := 0; i < a.len; i++ {
					if a.data[i] != tt.want[i] {
						t.Errorf("AddElement() element at %d = %s, want %s", i, a.data[i], tt.want[i])
					}
				}
			})
		}
	})

	t.Run("GrowthBehavior", func(t *testing.T) {
		a, _ := NewArray(2)
		a.AddElementEnd("a")
		a.AddElementEnd("b")
		
		if a.cap != 2 {
			t.Errorf("Initial capacity = %d, want 2", a.cap)
		}
		
		a.AddElementEnd("c")
		if a.cap != 4 {
			t.Errorf("Capacity after growth = %d, want 4", a.cap)
		}
		
		val, _ := a.GetElement(0)
		if val != "a" {
			t.Errorf("Element at 0 = %s, want 'a'", val)
		}
		val, _ = a.GetElement(1)
		if val != "b" {
			t.Errorf("Element at 1 = %s, want 'b'", val)
		}
		val, _ = a.GetElement(2)
		if val != "c" {
			t.Errorf("Element at 2 = %s, want 'c'", val)
		}
	})

	t.Run("IsInArray", func(t *testing.T) {
		tests := []struct {
			name  string
			items []string
			key   string
			want  int
		}{
			{"Found first", []string{"apple", "banana", "cherry"}, "apple", 0},
			{"Found middle", []string{"apple", "banana", "cherry"}, "banana", 1},
			{"Found last", []string{"apple", "banana", "cherry"}, "cherry", 2},
			{"Not found", []string{"apple", "banana", "cherry"}, "date", -1},
			{"Empty array", []string{}, "anything", -1},
			{"With empty string", []string{"", "a", ""}, "", 0},
			{"Multiple matches", []string{"a", "b", "a"}, "a", 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, _ := NewArrayFromList(tt.items)
				if got := a.IsInArray(tt.key); got != tt.want {
					t.Errorf("IsInArray() = %d, want %d", got, tt.want)
				}
			})
		}
	})

	t.Run("LengthAndCapacity", func(t *testing.T) {
		a, _ := NewArray(3)
		if a.GetLength() != 0 {
			t.Errorf("GetLength() = %d, want 0", a.GetLength())
		}
		if a.GetCapacity() != 3 {
			t.Errorf("GetCapacity() = %d, want 3", a.GetCapacity())
		}
		
		a.AddElementEnd("a")
		a.AddElementEnd("b")
		if a.GetLength() != 2 {
			t.Errorf("GetLength() after adds = %d, want 2", a.GetLength())
		}
		
		a.DeleteElement(0)
		if a.GetLength() != 1 {
			t.Errorf("GetLength() after delete = %d, want 1", a.GetLength())
		}
	})
}

func TestIOOperations(t *testing.T) {
	t.Run("Print", func(t *testing.T) {
		tests := []struct {
			name     string
			elements []string
			want     string
		}{
			{"Empty array", []string{}, "\n"},
			{"Single element", []string{"a"}, "a\n"},
			{"Multiple elements", []string{"a", "b", "c"}, "a b c\n"},
			{"With spaces", []string{"hello world", "test"}, "hello world test\n"},
			{"Empty strings", []string{"", "", ""}, "  \n"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, _ := NewArrayFromList(tt.elements)
				output := captureOutput(a.Print)
				if output != tt.want {
					t.Errorf("Print() = %q, want %q", output, tt.want)
				}
			})
		}
	})

	t.Run("BinaryFileOperations", func(t *testing.T) {
		tempDir := t.TempDir()
		binFile := filepath.Join(tempDir, "test.bin")

		t.Run("Valid round trip", func(t *testing.T) {
			original, _ := NewArrayFromList([]string{"hello", "world", "42", "", "with\nnewline"})
			if err := original.WriteBinary(binFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			binRead, _ := NewArray(1)
			if err := binRead.ReadBinary(binFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}

			if binRead.len != original.len {
				t.Errorf("Read length = %d, want %d", binRead.len, original.len)
			}
			if !reflect.DeepEqual(original.data[:original.len], binRead.data[:binRead.len]) {
				t.Errorf("Binary read mismatch: got %v, want %v", binRead.data[:binRead.len], original.data[:original.len])
			}
		})

		t.Run("Empty array", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.bin")
			emptyArray, _ := NewArray(1)
			if err := emptyArray.WriteBinary(emptyFile); err != nil {
				t.Fatalf("WriteBinary() failed: %v", err)
			}

			readArray, _ := NewArray(1)
			if err := readArray.ReadBinary(emptyFile); err != nil {
				t.Fatalf("ReadBinary() failed: %v", err)
			}
			if readArray.len != 0 {
				t.Errorf("Read length = %d, want 0", readArray.len)
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.bin")
			a, _ := NewArray(1)
			if err := a.ReadBinary(nonExistent); err == nil {
				t.Error("ReadBinary() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.bin")
			file, _ := os.Create(invalidFile)
			file.Write([]byte{0x01, 0x02, 0x03})
			file.Close()

			if err := a.ReadBinary(invalidFile); err == nil {
				t.Error("ReadBinary() expected error for invalid file, got nil")
			}
		})
	})

	t.Run("TextFileOperations", func(t *testing.T) {
		tempDir := t.TempDir()
		txtFile := filepath.Join(tempDir, "test.txt")

		t.Run("Valid round trip", func(t *testing.T) {
			original, _ := NewArrayFromList([]string{"hello", "world", "42", "", "with space"})
			if err := original.WriteText(txtFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			txtRead, _ := NewArray(1)
			if err := txtRead.ReadText(txtFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if txtRead.len != original.len {
				t.Errorf("Read length = %d, want %d", txtRead.len, original.len)
			}
			if !reflect.DeepEqual(original.data[:original.len], txtRead.data[:txtRead.len]) {
				t.Errorf("Text read mismatch: got %v, want %v", txtRead.data[:txtRead.len], original.data[:original.len])
			}
		})

		t.Run("Empty array", func(t *testing.T) {
			emptyFile := filepath.Join(tempDir, "empty.txt")
			emptyArray, _ := NewArray(1)
			if err := emptyArray.WriteText(emptyFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readArray, _ := NewArray(1)
			if err := readArray.ReadText(emptyFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}
			if readArray.len != 0 {
				t.Errorf("Read length = %d, want 0", readArray.len)
			}
		})

		t.Run("File errors", func(t *testing.T) {
			nonExistent := filepath.Join(tempDir, "nonexistent.txt")
			a, _ := NewArray(1)
			if err := a.ReadText(nonExistent); err == nil {
				t.Error("ReadText() expected error for non-existent file, got nil")
			}

			invalidFile := filepath.Join(tempDir, "invalid.txt")
			os.WriteFile(invalidFile, []byte("invalid\nformat"), 0644)

			if err := a.ReadText(invalidFile); err == nil {
				t.Error("ReadText() expected error for invalid file, got nil")
			}

			invalidLenFile := filepath.Join(tempDir, "invalid_len.txt")
			os.WriteFile(invalidLenFile, []byte("invalid_length\n"), 0644)

			if err := a.ReadText(invalidLenFile); err == nil {
				t.Error("ReadText() expected error for invalid length format, got nil")
			}
		})

		t.Run("Edge cases", func(t *testing.T) {
			edgeFile := filepath.Join(tempDir, "edge.txt")
			edgeArray, _ := NewArrayFromList([]string{"", " "})
			if err := edgeArray.WriteText(edgeFile); err != nil {
				t.Fatalf("WriteText() failed: %v", err)
			}

			readArray, _ := NewArray(1)
			if err := readArray.ReadText(edgeFile); err != nil {
				t.Fatalf("ReadText() failed: %v", err)
			}

			if readArray.len != 2 {
				t.Errorf("Read length = %d, want 2", readArray.len)
			}
			expected := []string{"", " "}
			for i := 0; i < 2; i++ {
				if readArray.data[i] != expected[i] {
					t.Errorf("Element at %d = %q, want %q", i, readArray.data[i], expected[i])
				}
			}
		})
	})
}

func TestErrorConditions(t *testing.T) {
	t.Run("Constructor errors", func(t *testing.T) {
		_, err := NewArray(0)
		if err == nil {
			t.Error("NewArray(0) expected error, got nil")
		}
		
		_, err = NewArray(-5)
		if err == nil {
			t.Error("NewArray(-5) expected error, got nil")
		}
	})

	t.Run("Operation errors", func(t *testing.T) {
		a, _ := NewArray(1)
		
		_, err := a.GetElement(0)
		if err == nil {
			t.Error("GetElement(0) on empty array expected error, got nil")
		}
		
		err = a.SetElement("test", 0)
		if err == nil {
			t.Error("SetElement(0) on empty array expected error, got nil")
		}
		
		err = a.DeleteElement(0)
		if err == nil {
			t.Error("DeleteElement(0) on empty array expected error, got nil")
		}
		
		err = a.AddElementAtIndex("test", 1)
		if err == nil {
			t.Error("AddElementAtIndex(1) on empty array expected error, got nil")
		}
	})

	t.Run("File operation errors", func(t *testing.T) {
		a, _ := NewArray(1)
		
		invalidPath := "/invalid/path/test.bin"
		if err := a.WriteBinary(invalidPath); err == nil {
			t.Error("WriteBinary() to invalid path expected error, got nil")
		}
		
		if err := a.WriteText(invalidPath); err == nil {
			t.Error("WriteText() to invalid path expected error, got nil")
		}
		
		t.Run("ReadBinary errors", func(t *testing.T) {
			tempDir := t.TempDir()
			
			emptyFile := filepath.Join(tempDir, "empty.bin")
			os.Create(emptyFile)
			if err := a.ReadBinary(emptyFile); err == nil {
				t.Error("ReadBinary(empty file) expected error, got nil")
			}
			
			partialFile := filepath.Join(tempDir, "partial.bin")
			file, _ := os.Create(partialFile)
			binary.Write(file, binary.LittleEndian, uint32(1))
			file.Close()
			if err := a.ReadBinary(partialFile); err == nil {
				t.Error("ReadBinary(partial file) expected error, got nil")
			}
		})
		
		t.Run("ReadText errors", func(t *testing.T) {
			tempDir := t.TempDir()
			
			emptyFile := filepath.Join(tempDir, "empty.txt")
			os.Create(emptyFile)
			if err := a.ReadText(emptyFile); err == nil {
				t.Error("ReadText(empty file) expected error, got nil")
			}
			
			onlyLength := filepath.Join(tempDir, "only_length.txt")
			os.WriteFile(onlyLength, []byte("1\n"), 0644)
			if err := a.ReadText(onlyLength); err == nil {
				t.Error("ReadText(file with only length) expected error, got nil")
			}
			
			invalidCount := filepath.Join(tempDir, "invalid_count.txt")
			os.WriteFile(invalidCount, []byte("3\none\ntwo\n"), 0644)
			if err := a.ReadText(invalidCount); err == nil {
				t.Error("ReadText(file with invalid element count) expected error, got nil")
			}
		})
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("Zero-length strings", func(t *testing.T) {
		a, _ := NewArray(3)
		a.AddElementEnd("")
		a.AddElementEnd("")
		a.AddElementEnd("")
		
		if a.len != 3 {
			t.Errorf("Length after adding empty strings = %d, want 3", a.len)
		}
		
		for i := 0; i < 3; i++ {
			val, err := a.GetElement(i)
			if err != nil {
				t.Fatalf("GetElement(%d) failed: %v", i, err)
			}
			if val != "" {
				t.Errorf("Element at %d = %q, want empty string", i, val)
			}
		}
	})

	t.Run("Maximum growth", func(t *testing.T) {
		a, _ := NewArray(1)
		for i := 0; i < 5; i++ {
			a.AddElementEnd(strconv.Itoa(i))
		}
		
		if a.cap != 8 {
			t.Errorf("Capacity after growth = %d, want 8", a.cap)
		}
		
		for i := 0; i < 5; i++ {
			val, _ := a.GetElement(i)
			if val != strconv.Itoa(i) {
				t.Errorf("Element at %d = %s, want %s", i, val, strconv.Itoa(i))
			}
		}
	})

	t.Run("Large string handling", func(t *testing.T) {
		a, _ := NewArray(2)
		largeString := string(make([]byte, 10000))
		a.AddElementEnd(largeString)
		
		if a.len != 1 {
			t.Errorf("Length after large string = %d, want 1", a.len)
		}
		
		val, _ := a.GetElement(0)
		if len(val) != 10000 {
			t.Errorf("Large string length = %d, want 10000", len(val))
		}
	})
}