// Copyright 2015 The go-ethereum Authors
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

package miner

import (
	"encoding/binary"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/wtc/go-wtc/consensus"
	"github.com/wtc/go-wtc/log"
)

type CpuAgent struct {
	mu sync.Mutex

	workCh        chan *Work
	stop          chan struct{}
	quitCurrentOp chan struct{}
	returnCh      chan<- *Result
	serverFound   chan uint64

	chain  consensus.ChainReader
	engine consensus.Engine

	isMining int32 // isMining indicates whether the agent is currently mining
}

func NewCpuAgent(chain consensus.ChainReader, engine consensus.Engine) *CpuAgent {
	miner := &CpuAgent{
		chain:  chain,
		engine: engine,
		stop:   make(chan struct{}, 1),
		workCh: make(chan *Work, 1),
		serverFound: make(chan uint64, 1),
	}
	isGpu, _, getport := engine.IsGPU()
	if isGpu {
		go startServer(miner.serverFound, getport)
	}
	return miner
}

func (self *CpuAgent) Work() chan<- *Work {
	return self.workCh
}
func (self *CpuAgent) SetReturnCh(ch chan<- *Result) {
	self.returnCh = ch
}

func (self *CpuAgent) Stop() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 1, 0) {
		return // agent already stopped
	}
	self.stop <- struct{}{}
done:
	// Empty work channel
	for {
		select {
		case <-self.workCh:
		default:
			break done
		}
	}
}

func (self *CpuAgent) Start() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 0, 1) {
		return // agent already started
	}
	go self.update()
}

func (self *CpuAgent) update() {
out:
	for {
		select {
		case work := <-self.workCh:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
			}
			self.quitCurrentOp = make(chan struct{})
			go self.mine(work, self.quitCurrentOp)
			self.mu.Unlock()
		case <-self.stop:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
				self.quitCurrentOp = nil
			}
			self.mu.Unlock()
			break out
		}
	}
}

func (self *CpuAgent) mine(work *Work, stop <-chan struct{}) {
	if result, err := self.engine.Seal(self.chain, work.Block, stop, self.serverFound); result != nil {
		log.Info("a new block seal finish.", "blockheight", result.Number(), "hash", result.Hash())
		self.returnCh <- &Result{work, result}
	} else {
		if err != nil {
			log.Warn("Block sealing failed", "err", err)
		}
		self.returnCh <- nil
	}
}

func (self *CpuAgent) GetHashRate() int64 {
	if pow, ok := self.engine.(consensus.PoW); ok {
		return int64(pow.Hashrate())
	}
	return 0
}

func startServer(serverFound chan uint64, listenPort int64) {
	listenP := strconv.FormatInt(listenPort, 10)
	//set up socket, listen port
	netListen, err := net.Listen("tcp", "localhost:"+listenP)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		panic(err.Error())
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		//Log(conn.RemoteAddr().String(), " tcp connect success")
		handleConnection(conn, serverFound)
	}
}

//process connect
func handleConnection(conn net.Conn, serverFound chan uint64) {
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {

		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		t, nonce := handeByte(buffer[:n])
		if t == 1 {
			GPUHashrate = int64(nonce)
		} else {
			serverFound <- nonce
		}

	}

}

func handeByte(res []byte) (byte, uint64) {
	if res[0] == 1 {
		rate := make([]byte, 0)
		for _, b := range res[1:] {
			if b == 0 {
				break
			}
			rate = append(rate, b)
		}
		i, err := strconv.Atoi(string(rate))
		if err != nil {
			panic(err)
		}
		return res[0], uint64(i)
	} else {
		//blockByte := res[:4]
		//blockNumber := new(big.Int).SetBytes(blockByte)
		//input := res[4:36]
		// nonce := new(big.Int).SetBytes(res[36:])
		nonce := binary.BigEndian.Uint64(res[36:])
		//fmt.Println("blockNumber =", blockNumber, "\nand input =", input, "\nand nonce =", res[36:], nonce)
		return res[0], nonce
	}

}
