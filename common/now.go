// Copyright 2014 The go-ethereum Authors
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

package common

import (
	"time"
	"math/big"
)


func getNow() *big.Int {
	return big.NewInt(time.Now().Unix())
}

func SecChangeToMine(i *big.Int) (*big.Int,*big.Int) {
	return new(big.Int).Div(i,Big60),new(big.Int).Mod(i,Big60)
}

func mineChangeToHour(i *big.Int) (*big.Int,*big.Int) {
	return new(big.Int).Div(i,Big60),new(big.Int).Mod(i,Big60)
}

func HourChangeToDay(i *big.Int) (*big.Int,*big.Int) {
	return new(big.Int).Div(i,Big24),new(big.Int).Mod(i,Big24)
}