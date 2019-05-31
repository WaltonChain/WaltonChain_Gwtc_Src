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
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/wtc/go-wtc/common"
	"github.com/wtc/go-wtc/consensus"
	"github.com/wtc/go-wtc/core/types"
	"github.com/wtc/go-wtc/log"
)

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
func (ethash *Ethash) SealbyGPU(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}, serverFound chan uint64, t int) (*types.Block, error) {
	oldbalance, coinage, preNumber, preTime := chain.GetBalanceAndCoinAgeByHeaderHash(block.Header().Coinbase)
	balance := new(big.Int).Add(oldbalance, big.NewInt(1e+18))
	//---------------------------------	--------------
	Time := block.Header().Time
	Number := block.Header().Number
	if preTime.Cmp(Time) < 0 && preNumber.Cmp(Number) < 0 {
		t := new(big.Int).Sub(Time, preTime)
		coinage = new(big.Int).Add(new(big.Int).Mul(balance, t), coinage)
	}

	// fmt.Println("set coinage to: ",coinage)

	fmt.Print()
	//-----------------------------------------------
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if ethash.fakeMode {
		header := block.Header()
		header.Nonce, header.MixDigest = types.BlockNonce{}, common.Hash{}
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	if ethash.shared != nil {
		return ethash.shared.SealbyGPU(chain, block, stop, serverFound, 1)
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
	// if threads == 0 {
	// 	threads = runtime.NumCPU()
	// }
	// if threads < 0 {
	// 	threads = 0 // Allows disabling local mining without extra logic around local/remote
	// }
	threads = 1
	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64, serverFound chan uint64) {
			defer pend.Done()
			ethash.minebyGPU(block, id, nonce, abort, found, balance, coinage, serverFound, t)
		}(0, uint64(ethash.rand.Int63()), serverFound)
	}
	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		sendStop(block, ethash.GPUPort)
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		close(abort)
	case <-ethash.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return ethash.SealbyGPU(chain, block, stop, serverFound, 1)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (ethash *Ethash) minebyGPU(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block, balance *big.Int, coinage *big.Int, serverFound chan uint64, t int) {
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
		// attempts = int64(0)
		nonce = seed
	)
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
	var servernonce uint64

	if t == 0 {
		time.Sleep(time.Second * 2)
		send(0, nonce, header.Number, hash, target, order, ethash.GPUPort)
		fmt.Println("send start")
	}
	for {
		select {
		case <-abort:
			//sendStop(block)
			logger.Trace("Ethash nonce search aborted", "attempts", servernonce-seed)
			return
		case servernonce = <-serverFound:
			digest, result := myx11(hash, servernonce, order)
			if Compare(result, FullTo32(target.Bytes()), 32) < 1 {
				// send(0, nonce, header.Number, hash, target, order)

				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(servernonce)
				header.MixDigest = common.BytesToHash(digest)
				header.CoinAge = coinage
				select {
				case found <- block.WithSeal(header):
					logger.Trace("Ethash nonce found and reported", "attempts", servernonce-seed, "nonce", servernonce)
				case <-abort:
					logger.Trace("Ethash nonce found but discarded", "attempts", servernonce-seed, "nonce", servernonce)
				}
				return
			} else {
				send(0, nonce, header.Number, hash, target, order, ethash.GPUPort)
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

func sendStop(block *types.Block, port int64) {
	fmt.Println("send stop")
	var (
		header = block.Header()
		hash   = header.HashNoNonce().Bytes()
		target = big.NewInt(0)
	)
	number := big.NewInt(0)
	order := getX11Order(hash, 11)
	send(1, 0, number, hash, target, order, port)
}

func send(control int, nonce uint64, number *big.Int, input []byte, target *big.Int, order []byte, port int64) {
	server := "127.0.0.1:" + strconv.FormatInt(port, 10)
	fmt.Println("send to ", server)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	sender(conn, control, nonce, number, input, target, order)

}

func sender(conn net.Conn, control int, nonce uint64, number *big.Int, input []byte, target *big.Int, order []byte) {
	words := encodeByte(control, number, input, nonce, target, 9e+18, order)
	conn.Write(words)
}

func encodeByte(control int, blockNumber *big.Int, input []byte, nonce uint64, target *big.Int, count uint64, alg []byte) []byte {
	str := make([]byte, 1)
	str[0] = byte(control)
	str = append(str, Int64ToBytes(blockNumber.Uint64())[4:]...)
	str = append(str, input...)
	str = append(str, Int64ToBytes(nonce)...)
	str = append(str, FullTo32(target.Bytes())...)
	str = append(str, Int64ToBytes(count)...)
	str = append(str, alg...)
	return str
}
func Int64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func FullTo32(word []byte) []byte {
	str := make([]byte, 32-len(word))
	str = append(str, word...)
	return str
}

func Compare(stra []byte, strb []byte, length int) int {
	for i := 0; i < length; i++ {
		if stra[i] > strb[i] {
			return 1
		} else {
			if stra[i] < strb[i] {
				return -1
			} else {
				continue
			}
		}
	}
	return 0
}
