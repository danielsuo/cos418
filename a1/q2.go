package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`
  total := 0
	for num := range nums {
    total += num
	}
  out <- total
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	ints, err := readInts(file)
	if err != nil {
		panic(err)
	}

	// Conservatively, we size the buffer to number of workers
	out := make(chan int, num)
	ins := make([]chan int, num)

	for i, _ := range ins {
		ins[i] = make(chan int, len(ints)/num+1)
		go sumWorker(ins[i], out)
	}

	// The fastest thing to do would be to allocate slices of ints array
	// but since we can't modify the signature, use a round robin
	for i, val := range ints {
		ins[i%num] <- val
	}

  for i, _ := range ins {
    close(ins[i])
  }

	total := 0

  for i := 0; i < num; i++ {
    total += <-out
  }

	return total
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
