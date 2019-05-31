# GoDash - x11

Package x11 implements the hash and required functions in go.

## Install

```bash
$ go get godash.org/libs/cb-crypto/x11
```

## Usage

```go
	package main

	import (
		"fmt"
		"godash.org/libs/cb-crypto/x11"
	)

	func main() {
		hs, out := x11.New(), [32]byte{}
		hs.Hash([]byte("DASH"), out[:])
		fmt.Printf("%x \n", out[:])
	}
```

### Notes

Echo, Simd and Shavite do not have 100% test coverage, a full test on these
requires the test to hash a blob of bytes that is several gigabytes large.


## License

go-x11 is licensed under the [copyfree](http://copyfree.org) ISC license.
