package compression

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestPackerAndUnpacker(t *testing.T) {

	intTest := 5
	stringTest := "abcdefg"
	bytesTest := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	randoNumbers := 10000000
	seedSource := rand.NewSource(time.Now().UnixNano())
	randomNumberGenerator := rand.New(seedSource)
	numbers := make([]int, randoNumbers)

	var p Packer
	err := p.Add(intTest)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Add(stringTest)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Add(bytesTest)
	if err != nil {
		t.Fatal(err)
	}

	u := Unpacker{p.Bytes()}

	i, err := u.NextInt()
	if err != nil {
		t.Fatal(err)
	}

	if i != intTest {
		t.Fatalf("expected %d got %d", intTest, i)
	}

	s, err := u.NextString()

	if s != stringTest {
		t.Fatalf("expected %q got %q", stringTest, s)
	}

	b, err := u.NextBytes(len(bytesTest))

	if !bytes.Equal(b, bytesTest) {
		t.Fatalf("expected %q got %q", bytesTest, b)
	}

	p.Reset()
	if p.Size() != 0 {
		t.Fatal("expected packer size to be 0")
	}

	sign := 0

	// generate random numbers
	for idx := range numbers {

		if idx%2 == 0 {
			sign = -1
		} else {
			sign = 1
		}

		value := sign * int(randomNumberGenerator.Int31())
		numbers[idx] = value
		p.Add(value)
	}
	b = p.Bytes()
	u.Reset(b)

	if len(b) != u.Size() {
		t.Fatalf("size mismatch source %d unpacker %d", len(b), p.Size())
	}

	for idx, number := range numbers {
		n, err := u.NextInt()
		if err != nil {
			t.Fatal(err)
		}
		if n != number {
			t.Fatalf("idx %d expected %d got %d", idx, number, n)
		}
	}

}
