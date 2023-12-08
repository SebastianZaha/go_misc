# Advent of Code 2023 on Arduino, with TinyGo

To run:
`go run runner.go flash day2_1/main.go day2_1/in.txt`


# Scratchpad

```go
const numBufSize = 128

// just use the leaking gc from tinygo, to dynamically allocate until ram is full
// should provide maximum flexibility in using the RAM.
// you just alloc 2 rows worth of input, and a growing slice of nums, with some reasonable 'cap'
// then panic if that stuff does not fit

// ReadMemStats -> to inspect allocations ?

// taken from
// https://github.com/tinygo-org/tinygo/blob/release/src/machine/buffer.go
// FIXME: extract with generics? needs to be parametrized on a const size though,
// is that possible?

type numRingBuffer struct {
	nums [numBufSize]num
	head byte
	tail byte
}

func (rb *numRingBuffer) put(n num) {
	if !rb.full() {
		rb.head++
		rb.nums[rb.head] = n
	} else {
		panic("full ring buffer")
	}
}

func (rb *numRingBuffer) pop() num {
	if rb.used() != 0 {
		rb.tail++
		return rb.nums[rb.tail-1]
	} else {
		panic("popping empty ring buffer")
	}
}

func (rb *numRingBuffer) full() bool {
	return rb.used() == numBufSize
}

func (rb *numRingBuffer) used() byte {
	return rb.head - rb.tail
}
```