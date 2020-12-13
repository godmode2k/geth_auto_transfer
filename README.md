# Ethereum JSON RPC Example (Go version)


Summary
----------
> Test Ethereum JSON RPC for </br>
> eth.getBalance(), eth.sendTransaction(), eth.getBlockByNumber() </br>
> eth_call(): balanceOf(), transfer() </br>
> without Ethereum library.


Environment
----------
> build all and tested on GNU/Linux

    GNU/Linux: Ubuntu 16.04_x64 LTS
    Ethereum: Geth/v1.9.24-stable/linux-amd64/go1.15.5
    Go lang: go1.15.5 linux/amd64
    Network: Ethereum Local Private Network


eth.getBalance()
----------
```sh
$ go run eth_autotransfer_main.go

HOST: http://127.0.0.1:8544
eth_getBalance()
ether hex-string to int:  989444987110000000000 (wei)
ether balance:  989.44498711 (ether)
```


eth.sendTransaction()
----------
```sh
$ go run eth_autotransfer_main.go

HOST: http://127.0.0.1:8544
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
$ go run eth_autotransfer_main.go

HOST: http://127.0.0.1:8544
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
$ go run eth_autotransfer_main.go

HOST: http://127.0.0.1:8544
eth_call(): transfer()
txid:  0x9ba6409566c0c418678aa157da31a4fd8ca280298846f918b105a69b1c238be1
```


eth.getBlockByNumber()
----------
```sh
$ go run eth_autotransfer_main.go

HOST: http://127.0.0.1:8544
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
