![](images/wtc_logo.jpg)

# GO WTC
Waltonchain Mainnet User Manual


## Steps

### 1. Run docker container
Install latest distribution of [Go](https://golang.org "Go") if you don't have it already.  
`# sudo apt-get install -y build-essential`  

### 2. Compile source code
`# cd /usr/local/src`  
`# git clone https://github.com/WaltonChain/WaltonChain_Gwtc_Src.git`  
`# cd WaltonChain_Gwtc_Src`  
`# make gwtc`  
`# ./build/bin/gwtc version`  

### 3. Deploy
`# cd /usr/local/src/WaltonChain_Gwtc_Src/gwtc_bin/`  
`# cp ../build/bin/gwtc ./bin/gwtc`  
`# ./backend.sh`

### 4. Enter console
`# cd /usr/local/src/WaltonChain_Gwtc_Src/gwtc_bin/`  
`# ./bin/gwtc attach ./data/gwtc.ipc`

### 5. View information of the connected node
`# admin.peers`

### 6. Create account
`# personal.newAccount()`  
`# ******`  ---- Enter new account password  
`# ******`  ---- Confirm the new account password  

### 7. Mine
`# miner.start()`

### 8. Query
`# wtc.getBalance(wtc.coinbase)`

### 9. Unlock account
`# personal.unlockAccount(wtc.coinbase)`

### 10. Transfer
`# wtc.sendTransaction({from: wtc.accounts[0], to: wtc.accounts[1], value: web3.toWei(1)})`

### 11. Exit console
`# exit`

### 12. View log
`# cd /usr/local/src/WaltonChain_Gwtc_Src/gwtc_bin/`  
`# tail -f gwtc.log`

### 13. Stop gwtc
`# cd /usr/local/src/WaltonChain_Gwtc_Src/gwtc_bin/`  
`# ./stop.sh` 


## Acknowledgement
We hereby thank:  
· [Ethereum](https://www.ethereum.org/ "Ethereum")




