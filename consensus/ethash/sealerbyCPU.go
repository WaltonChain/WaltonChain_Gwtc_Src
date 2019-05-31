// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-wtc library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-wtc library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	crand "crypto/rand"
	// "fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"

	"github.com/wtc/go-wtc/common"
	"github.com/wtc/go-wtc/consensus"
	"github.com/wtc/go-wtc/core/types"
	"github.com/wtc/go-wtc/log"
)

func (ethash *Ethash) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}, serverFound chan uint64) (*types.Block, error) {
	if ethash.GPUMode {
		return ethash.SealbyGPU(chain, block, stop, serverFound, 0)
	}
	return ethash.SealbyCPU(chain, block, stop, serverFound)
}

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
func (ethash *Ethash) SealbyCPU(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}, serverFound chan uint64) (*types.Block, error) {

	oldbalance, coinage, preNumber, preTime := chain.GetBalanceAndCoinAgeByHeaderHash(block.Header().Coinbase)
	balance := new(big.Int).Add(oldbalance, big.NewInt(1e+18))
	//---------------------------------	--------------
	Time := block.Header().Time
	Number := block.Header().Number
	if preTime.Cmp(Time) < 0 && preNumber.Cmp(Number) < 0 {
		t := new(big.Int).Sub(Time, preTime)
		coinage = new(big.Int).Add(new(big.Int).Mul(balance, t), coinage)
	}

	// fmt.Println("disy.yin ====>coinage:",coinage, "log2", log2(coinage))

	//-----------------------------------------------
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if ethash.fakeMode {
		header := block.Header()
		header.Nonce, header.MixDigest = types.BlockNonce{}, common.Hash{}
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	if ethash.shared != nil {
		return ethash.shared.SealbyCPU(chain, block, stop, serverFound)
	}
	// Create a runner and the multiple search threads it directs
	abort := make(chan struct{})
	found := make(chan *types.Block)

	ethash.lock.Lock()
	threads := ethash.threads
	if ethash.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			ethash.lock.Unlock()
			return nil, err
		}
		ethash.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	ethash.lock.Unlock()
	if threads == 0 {
		threads = runtime.NumCPU()
	}
	if threads < 0 {
		threads = 0 // Allows disabling local mining without extra logic around local/remote
	}
	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {
			defer pend.Done()
			ethash.minebyCPU(block, id, nonce, abort, found, coinage)
		}(i, uint64(ethash.rand.Int63()))
	}
	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		close(abort)
	case <-ethash.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return ethash.SealbyCPU(chain, block, stop, serverFound)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (ethash *Ethash) minebyCPU(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block, coinage *big.Int) {

	// Extract some data from the header
	var (
		header = block.Header()
		hash   = header.HashNoNonce().Bytes()
		target = new(big.Int).Div(maxUint256, header.Difficulty)

		//number  = header.Number.Uint64()
		//dataset = ethash.dataset(number)
	)

	// Start generating random nonces until we abort or find a good one
	var (
		attempts = int64(0)
		nonce    = seed
	)
	// fmt.Println("disy.yin ====>Difficulty:",header.Difficulty, "log2:", log2(header.Difficulty))
	// fmt.Println("disy.yin ====>target:",target, "log2:", log2(target))

	logger := log.New("miner", id)
	logger.Trace("Started ethash search for new nonces", "seed", seed)
	bn_coinage := new(big.Int).Mul(coinage, big.NewInt(1))
	bn_coinage = Sqrt(bn_coinage, 6)
	bn_txnumber := new(big.Int).Mul(new(big.Int).SetUint64(header.TxNumber), big.NewInt(5e+18))
	bn_txnumber = Sqrt(bn_txnumber, 6)
	if bn_coinage.Cmp(big.NewInt(0)) > 0 {
		target.Mul(bn_coinage, target)
	}
	if bn_txnumber.Cmp(big.NewInt(0)) > 0 {
		target.Mul(bn_txnumber, target)
	}
	order := getX11Order(hash, 11)

	// send(nonce, header.Number, hash, target, order)
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			logger.Trace("Ethash nonce search aborted", "attempts", nonce-seed)
			ethash.hashrate.Mark(attempts)
			return

		default:
			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
			attempts++
			if (attempts % (1 << 15)) == 0 {
				ethash.hashrate.Mark(attempts)
				attempts = 0
			}
			// Compute the PoW value of this nonce
			digest, result := myx11(hash, nonce, order)
			if Compare(result, FullTo32(target.Bytes()), 32) < 1 {

				// Correct nonce found, create a new header with it
				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(nonce)
				header.MixDigest = common.BytesToHash(digest)
				header.CoinAge = coinage
				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					logger.Trace("Ethash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("Ethash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				return
			}
			nonce++
		}
	}
}

func log2(number *big.Int) uint {
	var x = big.NewInt(0)
	var i uint
	for i = 0; ; i++ {
		if x.Rsh(number, i).Cmp(big.NewInt(0)) < 1 {
			return i
		}

	}
	return 0
}
func Sqrt(oldnumber *big.Int, exp uint) *big.Int {
	number := new(big.Int).Div(oldnumber, big.NewInt(1e+14))
	var x = number.BitLen()
	var y = new(big.Int).Rsh(number, uint(x)*(exp-1)/exp)
	return y.Div(y, big.NewInt(32))
}
