// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package x11

import (
	"github.com/wtc/go-wtc/crypto/hash"

	"github.com/wtc/go-wtc/crypto/x11/blake"
	"github.com/wtc/go-wtc/crypto/x11/bmw"
	"github.com/wtc/go-wtc/crypto/x11/cubed"
	"github.com/wtc/go-wtc/crypto/x11/echo"
	"github.com/wtc/go-wtc/crypto/x11/groest"
	"github.com/wtc/go-wtc/crypto/x11/jhash"
	"github.com/wtc/go-wtc/crypto/x11/keccak"
	"github.com/wtc/go-wtc/crypto/x11/luffa"
	"github.com/wtc/go-wtc/crypto/x11/shavite"
	"github.com/wtc/go-wtc/crypto/x11/simd"
	"github.com/wtc/go-wtc/crypto/x11/skein"
)

////////////////

// Hash contains the state objects
// required to perform the x11.Hash.
type Hash struct {
	tha [64]byte
	thb [64]byte

	blake   hash.Digest
	bmw     hash.Digest
	cubed   hash.Digest
	echo    hash.Digest
	groest  hash.Digest
	jhash   hash.Digest
	keccak  hash.Digest
	luffa   hash.Digest
	shavite hash.Digest
	simd    hash.Digest
	skein   hash.Digest
}

// New returns a new object to compute a x11 hash.
func New() *Hash {
	ref := &Hash{}
	ref.blake = blake.New()
	ref.bmw = bmw.New()
	ref.cubed = cubed.New()
	ref.echo = echo.New()
	ref.groest = groest.New()
	ref.jhash = jhash.New()
	ref.keccak = keccak.New()
	ref.luffa = luffa.New()
	ref.shavite = shavite.New()
	ref.simd = simd.New()
	ref.skein = skein.New()
	return ref
}

// Hash computes the hash from the src bytes and stores the result in dst.
func (ref *Hash) Hash(src []byte, dst []byte, order []byte) {
	in := ref.tha[:]
	out := ref.thb[:]
	in = src[:]

	for i := 0; i < len(order); i++ {
		switch order[i] {
			case 'A':
				ref.blake.Write(in)
				ref.blake.Close(out, 0, 0)
				in = out[:]
			case 'B':
				ref.cubed.Write(in)
				ref.cubed.Close(out, 0, 0)
				in = out[:]
			case 'C':
				ref.echo.Write(in)
				ref.echo.Close(out, 0, 0)
				in = out[:]
			case 'D':
				ref.bmw.Write(in)
				ref.bmw.Close(out, 0, 0)
				in = out[:]
			case 'E':
				ref.jhash.Write(in)
				ref.jhash.Close(out, 0, 0)
				in = out[:]
			case 'F':
				ref.groest.Write(in)
				ref.groest.Close(out, 0, 0)
				in = out[:]
			case 'G':
				ref.luffa.Write(in)
				ref.luffa.Close(out, 0, 0)
				in = out[:]
			case 'H':
				ref.skein.Write(in)
				ref.skein.Close(out, 0, 0)
				in = out[:]
			case 'I':
				ref.simd.Write(in)
				ref.simd.Close(out, 0, 0)
				in = out[:]
			case 'J':
				ref.keccak.Write(in)
				ref.keccak.Close(out, 0, 0)
				in = out[:]
			case 'K':
				ref.shavite.Write(in)
				ref.shavite.Close(out, 0, 0)
				in = out[:]

		}
		copy(dst, out)
	}
}






func (ref *Hash) Hash1(src []byte, dst []byte, order [11]byte) {

	ta := ref.tha[:]
	tb := ref.thb[:]

	ref.blake.Write(src)
	ref.blake.Close(tb, 0, 0)

	ref.cubed.Write(tb)
	ref.cubed.Close(ta, 0, 0)

	ref.echo.Write(ta)
	ref.echo.Close(tb, 0, 0)

	ref.bmw.Write(tb)
	ref.bmw.Close(ta, 0, 0)

	ref.jhash.Write(ta)
	ref.jhash.Close(tb, 0, 0)
	ref.groest.Write(tb)
	ref.groest.Close(ta, 0, 0)
	ref.luffa.Write(ta)
	ref.luffa.Close(tb, 0, 0)

	ref.skein.Write(tb)
	ref.skein.Close(ta, 0, 0)
	ref.simd.Write(ta)
	ref.simd.Close(tb, 0, 0)
	ref.keccak.Write(tb)
	ref.keccak.Close(ta, 0, 0)

	ref.shavite.Write(ta)
	ref.shavite.Close(tb, 0, 0)

	copy(dst, tb)
}
