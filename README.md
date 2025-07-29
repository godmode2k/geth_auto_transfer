# Ethereum JSON RPC Example (Go version)


Summary
----------
> Test Ethereum JSON RPC for ETH, ERC-20, ERC-1155 (NFT) </br>
> without Ethereum library. </br>
> APIs: eth.getBalance(), eth.sendTransaction(), eth.call(), eth.getBlockByNumber() </br>
> ABIs: balanceOf(), transfer(), safeTransferFrom(), ... </br>


Environment
----------
> build all and tested on GNU/Linux

    GNU/Linux: Ubuntu 20.04_x64 LTS
    Ethereum: Geth/v1.10.15-stable/linux-amd64/go1.10.15
    Golang: go1.15.5 linux/amd64
    Node.js: node-v16.13.2
    Network: Ethereum Local Private Network


Tools installation
----------
```sh
Golang
$ wget https://dl.google.com/go/go1.15.5.linux-amd64.tar.gz
$ tar xzvf go1.15.5.linux-amd64.tar.gz -C /usr/local/
$ echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile

Ethereum
$ git clone https://github.com/ethereum/go-ethereum.git -b v1.10.15
$ cd go-ethereum-1.10.15 && make all
$ cd .. && ln -s go-ethereum-1.10.15/build/bin/geth .


You can see also:

1. Setup an Ethereum Private Network Node
 - https://github.com/godmode2k/blockchain/tree/master/build_guide

2. Ethereum Smart-Contract (ERC-20, ERC-1155) Tools (Truffle, Hardhat, Foundry), Create & Deploy,
   Dockerfile for Ethereum Private Network Node
 - https://github.com/godmode2k/blockchain/tree/master/build_guide/ethereum

3. Ethereum Block Explorer
 - https://github.com/godmode2k/eth_block_explorer
```


Run
----------
```sh
$ go run eth_autotransfer_main.go
```


eth.getBalance()
----------
```sh
eth_getBalance()
ether hex-string to int:  989444987110000000000 (wei)
ether balance:  989.44498711 (ether)
```


eth.sendTransaction()
----------
```sh
eth_sendTransaction()
send ether:
from = 0xe6e55eed00218faef27eed24def9208f3878b333
to = 0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454
value = 11.1357, gas = 70000, gasPrice = 100
txid:  0xac870ca0abff42862d375636629a06984bba24217d0a196c1e15b3f9aad969b9
```


eth.call(): balanceOf()
----------
```sh
eth_call(): balanceOf()
token name: abcd_coin
token_symbol: abcd
token_decimals: 18
token_total_supply: 300000000.000000
erc-20 token balance hex-string to int:  299990178552087110000000000 (wei)
erc-20 token balance: 299990178.55208711 (abcd)
```


eth.call(): transfer()
----------
```sh
eth_call(): transfer()
txid:  0x9ba6409566c0c418678aa157da31a4fd8ca280298846f918b105a69b1c238be1
```


eth.getBlockByNumber()
----------
```sh
eth_call(): eth_getBlockByNumber()
block start =  0
block end =  584

hash = 0xdc58a5ee5507cff6c2c7df684a73293959646316b2835c8ff36b4eaccd4731af
timestamp = 2019-08-11 02:28:48 +0900 KST
block_number = 398
from = 0xe6e55eed00218faef27eed24def9208f3878b333
ERC-20
method = 0xa9059cbb
token_contract address = 0x1249cda86774bc170cab843437dd37484f173ca8
token_to = 0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454
token_name = AbcdefCoin
token_symbol = ABCD
token_decimals = 18
token_total_supply = 100000.00000000
token_value_wei = 3000000000000000000000 (wei)
token_value_ABCD = 3000.00000000 (ABCD)

hash = 0xdb7799420ff7b9129a623732b4f620229da9ba89682e7865f23b2c43012e8a5f
timestamp = 2020-12-07 22:07:15 +0900 KST
block_number = 408
from = 0xe6e55eed00218faef27eed24def9208f3878b333
Ether
to = 0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454
ether hex-string to int:  15876512890000000000 (wei)
value_wei = 15876512890000000000 (wei)
value_ether = 15.87651289 (ether)

...
```


eth.call(): uri()
----------
```sh
eth_call(): uri()
erc-1155 URI hex-string to str:  http://127.0.0.1/api/token/{id}.json
erc-1155 URI:  http://127.0.0.1/api/token/0.json
token_id str:  0
token_id hex (from str literally):  30
token_id ASCII:  0
erc-1155 URI:  http://127.0.0.1/api/token/0000000000000000000000000000000000000000000000000000000000000030.json
```


eth_call(): eth_getBlockByNumber()
----------
```sh
eth_call(): eth_getBlockByNumber()
block start =  0
block end =  1353
...

ERC-1155 safeTransferFrom() transaction
hash = 0xa2cc536fb9c5a63943478335351054d77396f4f7bb242c5140c2d41de73f76c6
timestamp = 2022-01-31 23:50:53 +0900 KST
block_number = 397
from = 0xe6e55eed00218faef27eed24def9208f3878b333
token_contract address = 0x1249cda86774bc170cab843437dd37484f173ca8
token_from = 0xe6e55eed00218faef27eed24def9208f3878b333
token_to = 0x8f5b2b7608e3e3a3
token_id =  0x0000000000000000000000000000000000000000000000000000000000000000
token_amount =  1
token_uri (ASCII) =  http://127.0.0.1/api/token/0.json
token_uri (Hexadecimal) =  http://127.0.0.1/api/token/0000000000000000000000000000000000000000000000000000000000000030.json
token_data_length =  0x00000000000000000000000000000000000000000000000000000000000000a0
token_data =  0x0000000000000000000000000000000000000000000000000000000000000000

Ether
hash = 0x6e06b1a5fac9411df31e3a0204c1d7bf2c599ffdd647c44490a3386fb3e8eb44
timestamp = 2022-02-02 19:27:26 +0900 KST
block_number = 1324
from = 0xe6e55eed00218faef27eed24def9208f3878b333
to = 0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454
ether hex-string to int:  11135700000000000000 (wei)
value_wei = 11135700000000000000 (wei)
value_ether = 11.13570000 (ether)

ERC-20 transfer() transaction
hash = 0x516ef91be8d560fcb6d2bab8a0f1eab8efdb2a8d7ccfb0159a47b3985d4f13e6
timestamp = 2022-02-02 19:37:26 +0900 KST
block_number = 1345
from = 0xe6e55eed00218faef27eed24def9208f3878b333
token_contract address = 0xb5accfe1b7a59317a9f5a100dc1105ed66b2058c
token_to = 0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454
token_name = ABCToken
token_symbol = ABC
token_decimals = 18
token_total_supply = 10000.00000000
token_value_wei = 11135700000000000000 (wei)
token_value_ABC = 11.13570000 (ABC)

...
```



## Donation
If this project help you reduce time to develop, you can give me a cup of coffee :)

(BitcoinCash) -> bitcoincash:qqls8jsln7w5vzd32g4yrwprstu57aa8rgf4yvsm3m <br>
(Bitcoin) -> 16kC7PUd75rvmwom4oftXRyg3gR9KTPb4m <br>
(Ethereum) -> 0x90B45D2CBBB0367D50590659845C486497F89cBB <br>


