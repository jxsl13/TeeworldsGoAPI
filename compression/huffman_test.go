package compression

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_bubbleSort(t *testing.T) {
	type testStruct struct {
		name string
		list []*huffmanConstructNode
	}

	tests := []testStruct{}

	// 10! = 3628800 unique lists
	initialList := []*huffmanConstructNode{
		{0, 1},
		{0, 2},
		{0, 3},
		{0, 4},
		{0, 5},
		{0, 6},
		{0, 7},
		{0, 8},
		{0, 9},
		{0, 10},
	}

	// create all possible permutations for the list
	allPermutations := permutate(initialList)

	// put all permutations into tests
	for idx, permutation := range allPermutations {
		tests = append(tests,
			testStruct{
				fmt.Sprintf("#%d", idx+1),
				permutation,
			})
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// sort permutation
			bubbleSort(tt.list)

			// resulting permutation must have ordered frequencies
			for idx, value := range tt.list {
				if idx > 0 && tt.list[idx-1].Frequency < value.Frequency {
					t.Errorf("idx: %d = %d, idx: %d = %d",
						idx-1,
						tt.list[idx-1].Frequency,
						idx,
						value.Frequency,
					)

				}
			}
		})
	}
}

func permutate(a []*huffmanConstructNode) [][]*huffmanConstructNode {
	var res [][]*huffmanConstructNode
	calPermutation(a, &res, 0)
	return res
}
func calPermutation(arr []*huffmanConstructNode, res *[][]*huffmanConstructNode, k int) {
	for i := k; i < len(arr); i++ {
		swap(arr, i, k)
		calPermutation(arr, res, k+1)
		swap(arr, k, i)
	}
	if k == len(arr)-1 {
		r := make([]*huffmanConstructNode, len(arr))
		copy(r, arr)
		*res = append(*res, r)
		return
	}
}
func swap(arr []*huffmanConstructNode, i, k int) {
	arr[i], arr[k] = arr[k], arr[i]
}

func TestHuffman_Compress_Decompress(t *testing.T) {
	huffman := Huffman{}
	huffman.Init(nil)

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"#1", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := huffman.Compress(tt.input)
			if err != nil {
				t.Errorf("Huffman.Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(compressed) == 0 {
				t.Error("Huffman.Compress() : Compressed length is 0")
				return
			}

			decompressed, err := huffman.Decompress(compressed)
			if (err != nil) != tt.wantErr {
				t.Errorf("\nInput        = %v\nCompressed   = %v\n", tt.input, compressed)
				t.Errorf("Huffman.Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(tt.input, decompressed) {
				t.Errorf("\nInput        = %v\nCompressed   = %v\nDecompressed = %v", tt.input, compressed, decompressed)
				return
			}

		})
	}
}

func Test_initFrequenciesTable(t *testing.T) {
	sb := strings.Builder{}
	sb.WriteString("{\n")
	for _, f := range frequenciesTable {
		sb.WriteString(fmt.Sprintf("%d, ", byte(f)))
	}
	sb.WriteString("\n}")

	t.Errorf("\nTable: %v\n", sb.String())
}
