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

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Wtc network.
var MainnetBootnodes = []string{
	"enode://e1adaa6427fa63b7440f03188ab2f2e894593f41727a799d05bb08542b185412b96aca0af115d9b5e7628a94b1ce20ae2a30362fe95409ccaa03c7c12faf2d67@27.102.130.163:10101",
	"enode://bc75d1861480101805e2f9277543421c45821dfeb26a19f64a899e7752135f624b7e72374086dbbdeac346ee85302db33ae67b357758bfa390474938353e7429@47.88.226.40:10101",
	"enode://9bf63165d038b4739db680fbb510e9596164b0fc65f19bf53d15c112b485b09328f6c63925a3cfc5aff5ff44d5f8eecd1ec72a26c746a3d34fa20aa944996730@47.88.62.219:10101",

	"enode://aaaae77d11539c63d16c43dbc35734a255e9f97676241d866135a50afe40acc54d566f61c657d152b5caa9039fe9a076ac24d8ad01abd57cc4571279806a9a82@45.76.135.125:10101",
	"enode://cc9c8e97a669cce1af7d3e9acf8c7bf51505334abda1581a2a7cf811290ee7cdc91971561b3261b7560ec1919b30ea7c74239dbaee2bafce6d6b4b319fe5ab05@208.167.249.242:10101",
	"enode://6559b12fda8f90242236e43afe5fe978b822467074260f7f6583d01fb0d31de7db87aa776d3c6f25f1b0bf42f3415641943ea7fae0af0c5b78647a29c3d69168@209.250.237.179:10101",

	"enode://18fb4bbd45745468ff6a89375aa33687aa205741092f62c858c999a2df6ef9156d94585e4cfa0dd17a9c5d949a1f583d8bd7c75e591b6571161937c9e39d81ff@47.95.205.144:10101",
	"enode://7248a6f713007c575e4de3fbf7bc2f8595ad1e974003591b2d216b75bd4ddfeb6e55b56d0fbda89d703e2c5da31a837d91b4ce9b77b0ed7a930a0826150bab08@47.111.4.61:10101",
	"enode://0a36aac3e150d251dedb71164b0efd10fb13c2b140f0102b3b24d0f58ea8d894879d47f70012cc7d0041f4ae4fa624709872e598996e41cbaea218c24752c480@47.112.118.220:10101",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
}

// RinkebyV5Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network for the experimental RLPx v5 topic-discovery network.
var RinkebyV5Bootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
}
