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

package params

import "math/big"

const (
	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
	ExpByteGas            uint64 = 2    // Times ceil(log256(exponent)) for the EXP instruction.
	SloadGas              uint64 = 12    // Multiplied by the number of 32-byte words that are copied (round up) for any *COPY operation and added.
	CallValueTransferGas  uint64 = 2000  // Paid for CALL when the value transfer is non-zero.
	CallNewAccountGas     uint64 = 6000 // Paid for CALL when the destination address didn't exist prior.
	TxGas                 uint64 = 5000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	TxGasContractCreation uint64 = 10000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.
	TxDataZeroGas         uint64 = 1     // Per byte of data attached to a transaction that equals zero. NOTE: Not payable on data of calls between transactions.
	QuadCoeffDiv          uint64 = 512   // Divisor for the quadratic particle of the memory cost equation.
	SstoreSetGas          uint64 = 5000 // Once per SLOAD operation.
	LogDataGas            uint64 = 2     // Per byte in a LOG* operation's data.
	CallStipend           uint64 = 600  // Free gas given at beginning of call.

	Sha3Gas          uint64 = 8    // Once per SHA3 operation.
	Sha3WordGas      uint64 = 2     // Once per word of the SHA3 operation's data.
	SstoreResetGas   uint64 = 1200  // Once per SSTORE operation if the zeroness changes from zero.
	SstoreClearGas   uint64 = 1200  // Once per SSTORE operation if the zeroness doesn't change.
	SstoreRefundGas  uint64 = 3000 // Once per SSTORE operation if the zeroness changes to zero.
	JumpdestGas      uint64 = 1     // Refunded gas, once per SSTORE operation if the zeroness changes to zero.
	EpochDuration    uint64 = 30000 // Duration between proof-of-work epochs.
	CallGas          uint64 = 10    // Once per CALL operation & message call transaction.
	CreateDataGas    uint64 = 50   //
	CallCreateDepth  uint64 = 1024  // Maximum depth of call/create stack.
	ExpGas           uint64 = 3    // Once per EXP instruction
	LogGas           uint64 = 80   // Per LOG* operation.
	CopyGas          uint64 = 1     //
	StackLimit       uint64 = 1024  // Maximum size of VM stack allowed.
	TierStepGas      uint64 = 0     // Once per operation, for a selection of them.
	LogTopicGas      uint64 = 80   // Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.
	CreateGas        uint64 = 8000 // Once per CREATE operation & contract-creation transaction.
	SuicideRefundGas uint64 = 6000 // Refunded following a suicide operation.
	MemoryGas        uint64 = 1     // Times the address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL.
	TxDataNonZeroGas uint64 = 20    // Per byte of data attached to a transaction that is not equal to zero. NOTE: Not payable on data of calls between transactions.

	MaxCodeSize = 196608 // Maximum bytecode to permit for a contract

	// Precompiled contract gas prices

	EcrecoverGas            uint64 = 800   // Elliptic curve sender recovery gas price
	Sha256BaseGas           uint64 = 15     // Base price for a SHA256 operation
	Sha256PerWordGas        uint64 = 3     // Per-word price for a SHA256 operation
	Ripemd160BaseGas        uint64 = 150    // Base price for a RIPEMD160 operation
	Ripemd160PerWordGas     uint64 = 30    // Per-word price for a RIPEMD160 operation
	IdentityBaseGas         uint64 = 4     // Base price for a data copy operation
	IdentityPerWordGas      uint64 = 1      // Per-work price for a data copy operation
	ModExpQuadCoeffDiv      uint64 = 20     // Divisor for the quadratic particle of the big int modular exponentiation
	Bn256AddGas             uint64 = 120    // Gas needed for an elliptic curve addition
	Bn256ScalarMulGas       uint64 = 10000  // Gas needed for an elliptic curve scalar multiplication
	Bn256PairingBaseGas     uint64 = 25000 // Base price for an elliptic curve pairing check
	Bn256PairingPerPointGas uint64 = 20000  // Per-point price for an elliptic curve pairing check
)

var (
	GasLimitBoundDivisor   = big.NewInt(2048)                  // The bound divisor of the gas limit, used in update calculations.
	MinGasLimit            = big.NewInt(78000000000)                  // Minimum the gas limit may ever be.
	GenesisGasLimit        = big.NewInt(78000000000)               // Gas limit of the Genesis block.
	TargetGasLimit         = new(big.Int).Set(GenesisGasLimit) // The artificial target
	DifficultyBoundDivisor = big.NewInt(2048)                  // The bound divisor of the difficulty, used in the update calculations.
	GenesisDifficulty      = big.NewInt(131072)                // Difficulty of the Genesis block.
	MinimumDifficulty      = big.NewInt(131072)                // The minimum that the difficulty may ever be. 
	DurationLimit          = big.NewInt(6)                    // The decision boundary on the blocktime duration used to determine whether difficulty should go up or not.
	MinGasFloorCreateContract = big.NewInt(5000)               // The floor of gas to create contract. add by disy.yin disy.yin@gmail.com 2018-12-10
	AddedRewardForMN	   = big.NewInt(5000)               	//  added reward for MN
)
