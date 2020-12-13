/* --------------------------------------------------------------
Project:    Ethereum auto-transfer (accounts to specific address(hotwallet))
Purpose:
Author:     Ho-Jung Kim (godmode2k@hotmail.com)
Date:       Since Dec 4, 2020
Filename:   eth_autotransfer_main.go

Last modified:  Dec 11, 2020
License:

*
* Copyright (C) 2020 Ho-Jung Kim (godmode2k@hotmail.com)
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
-----------------------------------------------------------------
Note:
-----------------------------------------------------------------
1. Build:
	$ go build eth_autotransfer_main.go
    or
	$ go run eth_autotransfer_main.go
-------------------------------------------------------------- */
package main



//! Header
// ---------------------------------------------------------------

import (
    "fmt"
    "log"
    "bytes"
    "strconv"
    "math"
    "math/big"
    "encoding/hex"
    "strings"
    "time"

    "net/http"
    "io/ioutil"
    "encoding/json"

    // eth_auto_transfer
    "eth_auto_transfer/types"

    //"reflect"
)



//! Definition
// --------------------------------------------------------------------

var SERVER_ADDRESS = "127.0.0.1"
var SERVER_PORT = "8544"
var SERVER = SERVER_ADDRESS + ":" + SERVER_PORT
var URL = "http://" + SERVER_ADDRESS + ":" + SERVER_PORT



//! Implementation
// --------------------------------------------------------------------

func eth_get_balance(_address string) {
    // eth: eth_getBalance
    //
    // request:
    // $ curl -X POST --data
    //  '{"jsonrpc":"2.0",
    //  "method":"eth_getBalance",
    //  "params":["0xe6e55eed00218faef27eed24def9208f3878b333", "latest"],"id":0}'
    //  -H "Content-Type: application/json" http://127.0.0.1:8544/

    fmt.Println( "eth_getBalance()" )

    var result types.Result

    var params []interface{}
    //params = append( params, "0xe6e55eed00218faef27eed24def9208f3878b333", "latest" )
    params = append( params, _address, "latest" )
    request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_getBalance", Params: params, Id: 0 }
    //request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_getBalance",
    //	Params: []interface{}{"0xe6e55eed00218faef27eed24def9208f3878b333", "latest"}, Id: 0 }

    message, _ := json.Marshal( request_data )
    //fmt.Println( "message: ", request_data )
    response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
    defer response.Body.Close()
    if err != nil {
        log.Fatal( "http.Post: ", err )
    }

    //fmt.Println( "response: " )
    responseBody, err := ioutil.ReadAll( response.Body )
    if err != nil {
        log.Fatal( "ioutil.ReadAll: ", err )
    }

    //fmt.Println( string(responseBody) )
    err = json.Unmarshal( responseBody, &result )
    if err != nil {
        log.Fatal( "json.Unmarshal: ", err )
    }
    //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )

    // SEE:
    // - https://golang.org/pkg/math/big/
    // - https://golang.org/pkg/strconv/
    // - https://goethereumbook.org/account-balance/
    balance_wei_int := new(big.Int)
    balance_wei_int.SetString( result.Result[2:], 16 )
    fmt.Println( "ether hex-string to int: ", balance_wei_int, "(wei)" )
    balance_wei_float := new(big.Float)
    balance_wei_float.SetString( balance_wei_int.String() )
    balance_eth := new(big.Float).Quo(balance_wei_float, big.NewFloat(math.Pow10(18)))
    fmt.Println( "ether balance: ", balance_eth, "(ether)" )
}

func eth_send_transaction(_from string, _to string, _amount string, _gas string, _gasprice string) {
    // eth: eth_sendTransaction
    //
    // request:
    // $ curl -X POST --data
    //  '{"jsonrpc":"2.0",
    //  "method":"eth_sendTransaction",
    //  "params":[{
    //      "from": "0xe6e55eed00218faef27eed24def9208f3878b333",
    //      "to": "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454",
    //      "value": "0x8ca93d72e1ed4000", "gas": "0x11170", "gasPrice": "0x12a05f2000"}],"id":0}'
    //  -H "Content-Type: application/json" http://127.0.0.1:8544/

    fmt.Println( "eth_sendTransaction()" )

    var result types.Result

    from := _from
    to := _to
    gas := _gas
    gasprice := _gasprice
    value := _amount

    //gas := "70000" // 70000
    //gasprice := "100" // 100 * gwei(1e9)
    //value := "11.1357" // 10 * wei(1e18)
    gas_hex := ""
    gasprice_hex := ""
    value_hex := ""
    // ---
    {
        gas_int := new(big.Int)
        gas_float := new(big.Float)
        gasprice_int := new(big.Int)
        gasprice_float := new(big.Float)

        gas_float.SetString( gas )
        gasprice_float.SetString( gasprice )
        gasprice_decimals := big.NewFloat( math.Pow10(9) ) //new(big.Float)( math.Pow10(9) )
        gasprice_float_mul := new(big.Float).Mul( gasprice_float, gasprice_decimals ) // value * decimals(wei: 1e9)

        // float to int for hex
        // SEE: https://stackoverflow.com/questions/47545898/golang-convert-big-float-to-big-int
        gas_float.Int( gas_int )
        gasprice_float_mul.Int( gasprice_int )

        // ---

        value_float := new(big.Float)
        value_float.SetString( value )
        value_decimals := big.NewFloat( math.Pow10(18) ) //new(big.Float)( math.Pow10(18) )
        value_float_mul := new(big.Float).Mul( value_float, value_decimals ) // value * decimals(wei: 1e18)
        // DO NOT USE [
        //value_result := value_float_mul.Text( 'f', 8 ) // precision: 8, no exponent
        //value_result := value_float_mul.Text( 'x', 8 ) // precision: 8, hexadecimal mantissa
        //fmt.Println( "result:", value_result )
        //
        // USE THIS
        // SEE: https://stackoverflow.com/questions/47545898/golang-convert-big-float-to-big-int
        value_int := new(big.Int)
        value_float_mul.Int( value_int ) // float to int for hex
        // ]

        //fmt.Println( "value:" , value, "value_float:", value_float, "value_decimals:", value_decimals )
        //fmt.Printf( "%f\n", value_float_mul )
        //fmt.Printf( "hex = %s\n", hex.EncodeToString([]byte(value_result)) ) // DO NOT USE
        //fmt.Printf( "%s, %s\n", value_int, value_int.Text(16) ) // hex

        // ---

        //gas_hex := "0x" + hex.EncodeToString( []byte(gas) )
        //gasprice_hex := "0x" + hex.EncodeToString( []byte(gasprice) )
        //value_hex := "0x" + hex.EncodeToString( []byte(string(value_result)) )

        gas_hex = "0x" + gas_int.Text( 16 )
        gasprice_hex = "0x" + gasprice_int.Text( 16 )
        value_hex = "0x" + value_int.Text( 16 )
    }
    // ---
    //from := "0xe6e55eed00218faef27eed24def9208f3878b333"
    ////to := "0x1e57f9561600b269a37437f02ce9da31e5b830ce" // erc-20: contract address (abcd token)
    //to := "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454"
    //holder_address := ""
    method := "eth_sendTransaction"
    //! DO NOT USE [
    // except: cancel pending transaction, ...
    //data := ""
    //nonce := ""
    //data_hex := "0x" + data
    //nonce_hex := "0x" + nonce
    // ]
    var params []interface{}

    request_data_param := types.RequestData_params_transaction {
        From: from, To: to, Value: value_hex, Gas: gas_hex, Gasprice: gasprice_hex,
        //! DO NOT USE [
        // except: cancel pending transaction, ...
        //Data: data_hex, Nonce: nonce_hex
        // ]
    }
    params = append( params, request_data_param )
    request_data := types.RequestData { Jsonrpc: "2.0", Method: method, Params: params, Id: 0 }


    {
        // unlock: personal_unlockAccount
        //
        // request:
        // $ curl -X POST --data
        //  '{"jsonrpc":"2.0",
        //  "method":"personal_unlockAccount",
        //  "params": ["0xe6e55eed00218faef27eed24def9208f3878b333","12345678",5], "id":0}'
        //  -H "Content-Type: application/json" http://127.0.0.1:8544/

        type Result struct {
            Jsonrpc string `json:"jsonrpc"`
            Id int `json:"id"`
            Result bool `json:"result"`
        }
        var result Result


        passphrase := "12345678"
        duration := 5
        var params []interface{}
        params = append( params, from, passphrase, duration )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "personal_unlockAccount", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )

        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
    }


    message, _ := json.Marshal( request_data )
    //fmt.Println( "message: ", request_data )

    fmt.Printf( "send ether:\nfrom = %s\nto = %s\nvalue = %s, gas = %s, gasPrice = %s\n", from, to, value, gas, gasprice )

    response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
    defer response.Body.Close()
    if err != nil {
        log.Fatal( "http.Post: ", err )
    }

    //fmt.Println( "response: " )
    responseBody, err := ioutil.ReadAll( response.Body )
    if err != nil {
        log.Fatal( "ioutil.ReadAll: ", err )
    }

    //fmt.Println( string(responseBody) )
    err = json.Unmarshal( responseBody, &result )
    if err != nil {
        log.Fatal( "json.Unmarshal: ", err )
    }
    //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
    fmt.Println( "txid: ", result.Result )
}

func erc20_get_balance(_to string, _holder_address string) {
    // eth erc-20: eth_call
    //params := Params_ERC20 { "0xe6e55eed00218faef27eed24def9208f3878b333", "0x70a08231" }


    // eth erc-20: balanceOf(address)
    //
    // request:
    // $ curl -X POST --data
    //  '{"jsonrpc":"2.0",
    //  "method":"eth_call",
    //  "params":[{"to": "0x1e57f9561600b269a37437f02ce9da31e5b830ce", // ABCD token contract address
    //  "data":"0x70a08231000000000000000000000000e6e55eed00218faef27eed24def9208f3878b333"}, "latest"],"id":67}'
    //  -H "Content-Type: application/json" http://127.0.0.1:8544/
    //
    // method name:
    // > web3.sha3('balanceOf(address)')
    // "0x70a08231b98ef4ca268c9cc3f6b4590e4bfec28280db06bb5d45e689f2a360be"
    //
    // data:
    // <method name> + '0 x 24' + <token holder address>
    // <70a08231> 000000000000000000000000 <token holder address>

    fmt.Println( "eth_call(): balanceOf()" )

    var result types.Result

    //gas := "70000"
    //gasprice := "100"
    //value := ""
    //from := ""
    to := _to // erc-20 contract address
    holder_address := _holder_address
    //to := "0x1e57f9561600b269a37437f02ce9da31e5b830ce" // erc-20 contract address
    //holder_address := "0xe6e55eed00218faef27eed24def9208f3878b333"
    method := "0x70a08231"
    data := method + "000000000000000000000000" + holder_address[2:]

    var params []interface{}
    //request_data_param := types.RequestData_params_erc20_transaction { From: from, To: to, Value: value, Gas: gas, Gasprice: gasprice, Data: data }
    request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
    params = append( params, request_data_param, "latest" )
    request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

    message, _ := json.Marshal( request_data )
    //fmt.Println( "message: ", request_data )
    response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
    defer response.Body.Close()
    if err != nil {
        log.Fatal( "http.Post: ", err )
    }

    //fmt.Println( "response: " )
    responseBody, err := ioutil.ReadAll( response.Body )
    if err != nil {
        log.Fatal( "ioutil.ReadAll: ", err )
    }

    //fmt.Println( string(responseBody) )
    err = json.Unmarshal( responseBody, &result )
    if err != nil {
        log.Fatal( "json.Unmarshal: ", err )
    }
    //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )

    // SEE:
    // - https://golang.org/pkg/math/big/
    // - https://golang.org/pkg/strconv/
    // - https://goethereumbook.org/account-balance/
    //balance_wei_int := new(big.Int)
    //balance_wei_int.SetString( result.Result[2:], 16 )
    //fmt.Println( "hex-string to int: ", balance_wei_int, "(wei)" )
    //balance_wei_float := new(big.Float)
    //balance_wei_float.SetString( balance_wei_int.String() )
    //balance_token := new(big.Float).Quo(balance_wei_float, big.NewFloat(math.Pow10(18)))
    //fmt.Printf( "erc-20 token: %f ()\n", balance_token )


    balance_wei := result.Result
    _token_name := ""
    _token_symbol := ""
    _token_decimals := ""
    _token_total_supply := ""

    {
        // Token: name
        method = "0x06fdde03"
        data = method + "000000000000000000000000" + holder_address[2:]

        var params []interface{}
        request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
        params = append( params, request_data_param, "latest" )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )
        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
        _token_name = result.Result
    }


    {
        // Token: symbol
        method = "0x95d89b41"
        data = method + "000000000000000000000000" + holder_address[2:]

        var params []interface{}
        request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
        params = append( params, request_data_param, "latest" )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )
        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
        _token_symbol = result.Result
    }


    {
        // Token: decimals
        method = "0x313ce567"
        data = method + "000000000000000000000000" + holder_address[2:]

        var params []interface{}
        request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
        params = append( params, request_data_param, "latest" )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )
        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
        _token_decimals = result.Result
    }


    {
        // Token: total_supply 
        method = "0x18160ddd"
        data = method + "000000000000000000000000" + holder_address[2:]

        var params []interface{}
        request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
        params = append( params, request_data_param, "latest" )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )
        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
        _token_total_supply = result.Result
    }

    //-----{
    // token name: 0x + [60 bytes] + [4 bytes] + [60 bytes] + [4 bytes]:
    token_name, _ := hex.DecodeString( _token_name[2 + 60 + 4 + 60 + 4:] )

    // token symbol: 0x + [60 bytes] + [4 bytes] + [60 bytes] + [4 bytes]:
    token_symbol, _ := hex.DecodeString( _token_symbol[2 + 60 + 4 + 60 + 4:] )

    // token decimals: 0x + [60 bytes] + [4 bytes]
    token_decimals_int := new(big.Int)
    token_decimals_int.SetString( _token_decimals[2:], 16 )
    token_decimals := token_decimals_int.String()
    token_decimals_int32, _ := strconv.Atoi( token_decimals )

    // token total supply:
    token_total_supply_int := new(big.Int)
    token_total_supply_int.SetString( _token_total_supply[2:], 16 )

    //token_total_supply := token_total_supply_int.String()
    token_total_supply_float := new(big.Float)
    token_total_supply_float.SetString( token_total_supply_int.String() )
    token_total_supply := new(big.Float).Quo(token_total_supply_float, big.NewFloat(math.Pow10(token_decimals_int32)))

    fmt.Println( "token name:", string(token_name) )
    fmt.Println( "token_symbol:", string(token_symbol) )
    fmt.Println( "token_decimals:", token_decimals )
    fmt.Printf( "token_total_supply: %f\n", token_total_supply )
    //-----}


    // SEE:
    // - https://golang.org/pkg/math/big/
    // - https://golang.org/pkg/strconv/
    // - https://goethereumbook.org/account-balance/
    balance_wei_int := new(big.Int)
    balance_wei_int.SetString( balance_wei[2:], 16 )
    fmt.Println( "erc-20 token balance hex-string to int: ", balance_wei_int, "(wei)" )
    balance_wei_float := new(big.Float)
    balance_wei_float.SetString( balance_wei_int.String() )
    balance_token := new(big.Float).Quo(balance_wei_float, big.NewFloat(math.Pow10(18)))
    fmt.Printf( "erc-20 token balance: %.8f (%s)\n", balance_token, token_symbol )
}

func erc20_send_transaction(_contract_address string, _from string, _to string, _amount string, _gas string, _gasprice string) {
    // eth erc-20: transfer(address,uint256)
    //
    // request:
    // $ curl -X POST --data
    //  '{"jsonrpc":"2.0",
    //  "method":"eth_sendTransaction",
    //  "params":[{
    //  "from":"0xe6e55eed00218faef27eed24def9208f3878b333",
    //  "to":"0x1e57f9561600b269a37437f02ce9da31e5b830ce",
    //  "gas":"0x11170","gasPrice":"0x174876e800",
    //  "data":"0xa9059cbb0000000000000000000000008f5b2b7608e3e3a3dc0426c3396420fbf18494540000000000000000000000000000000000000000000000000fc2d121ff694000"}],"id":0}'
    //  -H "Content-Type: application/json" http://127.0.0.1:8544/
    //
    // method name:
    // > web3.sha3('transfer(address,uint256)')
    // "0xa9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b"
    //
    // data:
    // <method name>           + // 4 bytes
    // '0 x 24' + <to address> + // 32 bytes
    // '0 x X' + <amount>        // 32 bytes
    //
    // <0xa9059cbb> 000000000000000000000000 <to address>
    // <0 x X> + <amount>

    fmt.Println( "eth_call(): transfer()" )

    var result types.Result

    contract_address := _contract_address // contract address
    from := _from
    to := _to
    //contract_address := "0x1e57f9561600b269a37437f02ce9da31e5b830ce" // contract address
    //from := "0xe6e55eed00218faef27eed24def9208f3878b333"
    //to := "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454"
    //holder_address := ""

    gas := _gas
    gasprice := _gasprice // 100 * gwei(1e9)
    value_amount := _amount // 10 * (erc-20 token decimals)

    value := "0" // for Ether, "0" fixed if ERC-20 transfer()
    //gas := "70000" // 70000
    //gasprice := "100" // 100 * gwei(1e9)
    //value_amount := "1.1" // 10 * (erc-20 token decimals)
    gas_hex := ""
    gasprice_hex := ""
    value_hex := ""
    value_amount_hex := ""
    // ---
    {
        gas_int := new(big.Int)
        gas_float := new(big.Float)
        gasprice_int := new(big.Int)
        gasprice_float := new(big.Float)

        gas_float.SetString( gas )
        gasprice_float.SetString( gasprice )
        gasprice_decimals := big.NewFloat( math.Pow10(9) ) //new(big.Float)( math.Pow10(9) )
        gasprice_float_mul := new(big.Float).Mul( gasprice_float, gasprice_decimals ) // value * decimals(wei: 1e9)

        // float to int for hex
        // SEE: https://stackoverflow.com/questions/47545898/golang-convert-big-float-to-big-int
        gas_float.Int( gas_int )
        gasprice_float_mul.Int( gasprice_int )

        // ---

        var token_decimals_int32 = 0
        {
            // Token: decimals
            method := "0x313ce567"
            data := method + "000000000000000000000000" + from[2:] // holder address

            var params []interface{}
            request_data_param := types.RequestData_params_erc20 { To: contract_address, Data: data }
            params = append( params, request_data_param, "latest" )
            request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

            message, _ := json.Marshal( request_data )
            //fmt.Println( "message: ", request_data )
            response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
            defer response.Body.Close()
            if err != nil {
                log.Fatal( "http.Post: ", err )
            }

            //fmt.Println( "response: " )
            responseBody, err := ioutil.ReadAll( response.Body )
            if err != nil {
                log.Fatal( "ioutil.ReadAll: ", err )
            }

            //fmt.Println( string(responseBody) )
            err = json.Unmarshal( responseBody, &result )
            if err != nil {
                log.Fatal( "json.Unmarshal: ", err )
            }
            //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
            _token_decimals := result.Result


            // token decimals: 0x + [60 bytes] + [4 bytes]
            token_decimals_int := new(big.Int)
            token_decimals_int.SetString( _token_decimals[2:], 16 )
            token_decimals := token_decimals_int.String()
            token_decimals_int32, _ = strconv.Atoi( token_decimals )
        }

        // ---

        value_amount_float := new(big.Float)
        value_amount_float.SetString( value_amount )
        value_amount_decimals := big.NewFloat( math.Pow10(token_decimals_int32) ) //new(big.Float)( math.Pow10(18) )
        value_amount_float_mul := new(big.Float).Mul( value_amount_float, value_amount_decimals ) // value * decimals(wei: 1e18)
        // DO NOT USE [
        //value_amount_result := value_amount_float_mul.Text( 'f', 8 ) // precision: 8, no exponent
        //value_amount_result := value_amount_float_mul.Text( 'x', 8 ) // precision: 8, hexadecimal mantissa
        //fmt.Println( "result:", value_amount_result )
        //
        // USE THIS
        // SEE: https://stackoverflow.com/questions/47545898/golang-convert-big-float-to-big-int
        value_amount_int := new(big.Int)
        value_amount_float_mul.Int( value_amount_int ) // float to int for hex
        // ]

        //fmt.Println( "value_amount:" , value_amount, "value_amount_float:", value_amount_float, "value_amount_decimals:", value_amount_decimals )
        //fmt.Printf( "%f\n", value_amount_float_mul )
        //fmt.Printf( "hex = %s\n", hex.EncodeToString([]byte(value_amount_result)) ) // DO NOT USE
        //fmt.Printf( "%s, %s\n", value_amount_int, value_amount_int.Text(16) ) // hex

        // ---

        //gas_hex := "0x" + hex.EncodeToString( []byte(gas) )
        //gasprice_hex := "0x" + hex.EncodeToString( []byte(gasprice) )
        //value_hex := "0x" + hex.EncodeToString( []byte(string(value_result)) )

        gas_hex = "0x" + gas_int.Text( 16 )
        gasprice_hex = "0x" + gasprice_int.Text( 16 )
        value_hex = "0x" + value // always '0x0' for erc-20
        value_amount_hex = "0x" + value_amount_int.Text( 16 )
    }
    // ---
    method := "0xa9059cbb"
    data := method + "000000000000000000000000" + to[2:] +
            strings.Repeat("0", 64 - len(value_amount_hex[2:])) + value_amount_hex[2:]

    var params []interface{}
    request_data_param := types.RequestData_params_erc20_transaction {
        From: from, To: contract_address, Value: value_hex, Gas: gas_hex, Gasprice: gasprice_hex,
        Data: data,
    }
    //request_data_param := types.RequestData_params_erc20 { To: to, Data: data }
    params = append( params, request_data_param )
    request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_sendTransaction", Params: params, Id: 0 }


    {
        // unlock: personal_unlockAccount
        //
        // request:
        // $ curl -X POST --data
        //  '{"jsonrpc":"2.0",
        //  "method":"personal_unlockAccount",
        //  "params": ["0xe6e55eed00218faef27eed24def9208f3878b333","12345678",5], "id":0}'
        //  -H "Content-Type: application/json" http://127.0.0.1:8544/

        type Result struct {
            Jsonrpc string `json:"jsonrpc"`
            Id int `json:"id"`
            Result bool `json:"result"`
        }
        var result Result


        passphrase := "12345678"
        duration := 5
        var params []interface{}
        params = append( params, from, passphrase, duration )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "personal_unlockAccount", Params: params, Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )

        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
    }


    message, _ := json.Marshal( request_data )
    //fmt.Println( "message: ", request_data )

    response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
    defer response.Body.Close()
    if err != nil {
        log.Fatal( "http.Post: ", err )
    }

    //fmt.Println( "response: " )
    responseBody, err := ioutil.ReadAll( response.Body )
    if err != nil {
        log.Fatal( "ioutil.ReadAll: ", err )
    }

    //fmt.Println( string(responseBody) )
    err = json.Unmarshal( responseBody, &result )
    if err != nil {
        log.Fatal( "json.Unmarshal: ", err )
    }
    //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
    fmt.Println( "txid: ", result.Result )
}

func get_blocks() {
    // eth_getBlockByNumber
    //
    // request:
    // $
    //

    fmt.Println( "eth_call(): eth_getBlockByNumber()" )

    //block_num_start_uint64 := uint64(502)
    //block_num_end_uint64 := uint64(503)
    block_num_start_uint64 := uint64(0)
    block_num_end_uint64 := uint64(0)

    {
        var result types.Result

        //var params []interface{}
        //request_data_param := types.RequestData {  }
        //params = append( params, "latest", true )
        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_blockNumber", Id: 0 }

        message, _ := json.Marshal( request_data )
        //fmt.Println( "message: ", request_data )
        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
        defer response.Body.Close()
        if err != nil {
            log.Fatal( "http.Post: ", err )
        }

        //fmt.Println( "response: " )
        responseBody, err := ioutil.ReadAll( response.Body )
        if err != nil {
            log.Fatal( "ioutil.ReadAll: ", err )
        }

        //fmt.Println( string(responseBody) )
        err = json.Unmarshal( responseBody, &result )
        if err != nil {
            log.Fatal( "json.Unmarshal: ", err )
        }
        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )

        block_num_end_int := new(big.Int)
        block_num_end_int.SetString( result.Result[2:], 16 )
        block_num_end := block_num_end_int.String()
        //block_num_end_int32, _ = strconv.Atoi( block_num_end )
        block_num_end_uint64, _ = strconv.ParseUint( block_num_end, 10, 64 )
    }

    {
        /*
        type Result struct {
            Jsonrpc string `json:"jsonrpc"`
            Id int `json:"id"`
            Result string `json:"result"`
        }
        var result Result

        type Result_block struct {
            Jsonrpc string `json:"jsonrpc"`
            Id int `json:"id"`
            Result interface{} `json:"result"`
        }
        var result_block Result_block
        */

        var result types.Result
        var result_block types.Result_block




        fmt.Println( "block start = ", block_num_start_uint64 )
        fmt.Println( "block end = ", block_num_end_uint64 )
        fmt.Println()
        for i := block_num_start_uint64; i < uint64(block_num_end_uint64); i++ {
            request_block_num_hex := ""
            request_block_num_int := new(big.Int)
            request_block_num_int.SetUint64( uint64(i) )
            request_block_num_hex = "0x" + request_block_num_int.Text( 16 )

            //fmt.Println( i, "(" + request_block_num_hex + ")", "----------" )

            var params []interface{}
            //request_data_param := types.RequestData {  }
            //params = append( params, "latest", true )
            params = append( params, request_block_num_hex, true )
            request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_getBlockByNumber", Params: params, Id: 0 }

            message, _ := json.Marshal( request_data )
            //fmt.Println( "message: ", request_data )
            response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
            defer response.Body.Close()
            if err != nil {
                log.Fatal( "http.Post: ", err )
            }

            //fmt.Println( "response: " )
            responseBody, err := ioutil.ReadAll( response.Body )
            if err != nil {
                log.Fatal( "ioutil.ReadAll: ", err )
            }

            //fmt.Println( string(responseBody) )
            err = json.Unmarshal( responseBody, &result_block )
            if err != nil {
                log.Fatal( "json.Unmarshal: ", err )
            }
            //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result_block.Result )


            //fmt.Println( "data =", result_block.Result )
            //fmt.Println( reflect.TypeOf(result_block.Result) )
            _txn := result_block.Result.(map[string]interface{})["transactions"]
            //fmt.Println( "size = ", len(_txn.([]interface{})) )
            if len(_txn.([]interface{})) <= 0 {
                //fmt.Println( "no transaction: size = 0" )
                continue
            }
            _txn = _txn.([]interface{})[0]
            txn := _txn.(map[string]interface{})
            //if txn["from"] != "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454" && txn["to"] != "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454" {
            //    continue
            //}

            timestamp_hex := result_block.Result.(map[string]interface{})["timestamp"]
            timestamp_int := new(big.Int)
            timestamp_int.SetString( timestamp_hex.(string)[2:], 16 )
            timestamp_unixtime := timestamp_int.String()
            //timestamp_int32, _ := strconv.Atoi( timestamp_unixtime )
            timestamp_int64, _ := strconv.ParseInt( timestamp_unixtime, 10, 64 )
            tx_timestamp_date := time.Unix( timestamp_int64, 0 )
            //tx_timestamp_date_rfc3339 := timestamp_date.Format( time.RFC3339 )


            tx_hash := txn["hash"]
            tx_block_number_hex := txn["blockNumber"]
            tx_block_number := ""
            {
                block_number_int := new(big.Int)
                block_number_int.SetString( tx_block_number_hex.(string)[2:], 16 )
                //fmt.Println( "ether hex-string to int: ", block_number_int )
                tx_block_number = block_number_int.String()
            }
            tx_from := txn["from"]
            tx_to := txn["to"]
            tx_value_wei_hex := txn["value"]
            tx_value_wei := ""
            tx_value := "" // Ether
            tx_input := txn["input"]

            tx_token_to := "" // for ERC-20
            tx_token_name := ""
            tx_token_symbol := ""
            tx_token_decimals := ""
            tx_token_total_supply := ""
            tx_token_amount_wei_hex := ""
            tx_token_amount_wei := ""
            tx_token_amount := ""

            if tx_to == nil {
                continue
            }

            //fmt.Println( "transaction =", _txn )
            fmt.Println( "hash =", tx_hash )
            fmt.Println( "timestamp =", tx_timestamp_date ) // "Y/m/d/ H:i:s"
            fmt.Println( "block_number =", tx_block_number )
            fmt.Println( "from =", tx_from )

            if txn["input"] == "0x" {
                fmt.Println( "Ether" )

                fmt.Println( "to =", tx_to )

                // SEE:
                // - https://golang.org/pkg/math/big/
                // - https://golang.org/pkg/strconv/
                // - https://goethereumbook.org/account-balance/
                amount_wei_int := new(big.Int)
                amount_wei_int.SetString( tx_value_wei_hex.(string)[2:], 16 )
                fmt.Println( "ether hex-string to int: ", amount_wei_int, "(wei)" )
                amount_wei_float := new(big.Float)
                amount_wei_float.SetString( amount_wei_int.String() )
                tx_value_float := new(big.Float).Quo(amount_wei_float, big.NewFloat(math.Pow10(18)))
                tx_value = fmt.Sprintf( "%.8f", tx_value_float )
                tx_value_wei = amount_wei_int.String()

                fmt.Println( "value_wei =", tx_value_wei, "(wei)" )
                fmt.Println( "value_ether =", tx_value, "(ether)" )
                fmt.Println()
            } else {
                fmt.Println( "ERC-20" )

                //fmt.Println( "input data =", tx_input )

                // token to: [2: 0x] + [8: method] + [0 x 24]
                tx_token_to = "0x" + tx_input.(string)[2 + 8 + 24:(2+8+24 + 40)]

                // amount: 32 bytes (64 chars): [2: 0x] + [8: method] + [0 x 24] + [40: to address]
                tx_token_amount_wei_hex = "0x" + tx_input.(string)[2 + 8 + 24 + 40:]

                method := ""
                data := ""

                fmt.Println( "method =", tx_input.(string)[:10] )
                if tx_input.(string)[:10] != "0xa9059cbb" {
                    fmt.Println( "Not ERC-20 transfer transaction" )
                    continue
                }

                {
                    _token_name := ""
                    _token_symbol := ""
                    _token_decimals := ""
                    _token_total_supply := ""

                    {
                        // Token: name
                        method = "0x06fdde03"
                        data = method + "000000000000000000000000" + tx_from.(string)[2:]

                        var params []interface{}
                        request_data_param := types.RequestData_params_erc20 { To: tx_to.(string), Data: data }
                        params = append( params, request_data_param, "latest" )
                        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

                        message, _ := json.Marshal( request_data )
                        //fmt.Println( "message: ", request_data )
                        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
                        defer response.Body.Close()
                        if err != nil {
                            log.Fatal( "http.Post: ", err )
                        }

                        //fmt.Println( "response: " )
                        responseBody, err := ioutil.ReadAll( response.Body )
                        if err != nil {
                            log.Fatal( "ioutil.ReadAll: ", err )
                        }

                        //fmt.Println( string(responseBody) )
                        err = json.Unmarshal( responseBody, &result )
                        if err != nil {
                            log.Fatal( "json.Unmarshal: ", err )
                        }
                        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
                        _token_name = result.Result
                    }

                    {
                        // Token: symbol
                        method = "0x95d89b41"
                        data = method + "000000000000000000000000" + tx_from.(string)[2:]

                        var params []interface{}
                        request_data_param := types.RequestData_params_erc20 { To: tx_to.(string), Data: data }
                        params = append( params, request_data_param, "latest" )
                        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

                        message, _ := json.Marshal( request_data )
                        //fmt.Println( "message: ", request_data )
                        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
                        defer response.Body.Close()
                        if err != nil {
                            log.Fatal( "http.Post: ", err )
                        }

                        //fmt.Println( "response: " )
                        responseBody, err := ioutil.ReadAll( response.Body )
                        if err != nil {
                            log.Fatal( "ioutil.ReadAll: ", err )
                        }

                        //fmt.Println( string(responseBody) )
                        err = json.Unmarshal( responseBody, &result )
                        if err != nil {
                            log.Fatal( "json.Unmarshal: ", err )
                        }
                        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
                        _token_symbol = result.Result
                    }

                    {
                        // Token: decimals
                        method = "0x313ce567"
                        data = method + "000000000000000000000000" + tx_from.(string)[2:]

                        var params []interface{}
                        request_data_param := types.RequestData_params_erc20 { To: tx_to.(string), Data: data }
                        params = append( params, request_data_param, "latest" )
                        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

                        message, _ := json.Marshal( request_data )
                        //fmt.Println( "message: ", request_data )
                        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
                        defer response.Body.Close()
                        if err != nil {
                            log.Fatal( "http.Post: ", err )
                        }

                        //fmt.Println( "response: " )
                        responseBody, err := ioutil.ReadAll( response.Body )
                        if err != nil {
                            log.Fatal( "ioutil.ReadAll: ", err )
                        }

                        //fmt.Println( string(responseBody) )
                        err = json.Unmarshal( responseBody, &result )
                        if err != nil {
                            log.Fatal( "json.Unmarshal: ", err )
                        }
                        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
                        _token_decimals = result.Result
                    }

                    {
                        // Token: total_supply 
                        method = "0x18160ddd"
                        data = method + "000000000000000000000000" + tx_from.(string)[2:]

                        var params []interface{}
                        request_data_param := types.RequestData_params_erc20 { To: tx_to.(string), Data: data }
                        params = append( params, request_data_param, "latest" )
                        request_data := types.RequestData { Jsonrpc: "2.0", Method: "eth_call", Params: params, Id: 0 }

                        message, _ := json.Marshal( request_data )
                        //fmt.Println( "message: ", request_data )
                        response, err := http.Post( URL, "application/json", bytes.NewBuffer(message) )
                        defer response.Body.Close()
                        if err != nil {
                            log.Fatal( "http.Post: ", err )
                        }

                        //fmt.Println( "response: " )
                        responseBody, err := ioutil.ReadAll( response.Body )
                        if err != nil {
                            log.Fatal( "ioutil.ReadAll: ", err )
                        }

                        //fmt.Println( string(responseBody) )
                        err = json.Unmarshal( responseBody, &result )
                        if err != nil {
                            log.Fatal( "json.Unmarshal: ", err )
                        }
                        //fmt.Println( "jsonrpc =" , result.Jsonrpc, ", id =", result.Id, ", result =", result.Result )
                        _token_total_supply = result.Result
                    }



                    //-----{
                    // token name: 0x + [60 bytes] + [4 bytes] + [60 bytes] + [4 bytes]:
                    __token_name, _ := hex.DecodeString( _token_name[2 + 60 + 4 + 60 + 4:] )
                    tx_token_name = string(__token_name)

                    // token symbol: 0x + [60 bytes] + [4 bytes] + [60 bytes] + [4 bytes]:
                    __token_symbol, _ := hex.DecodeString( _token_symbol[2 + 60 + 4 + 60 + 4:] )
                    tx_token_symbol = string(__token_symbol)

                    // token decimals: 0x + [60 bytes] + [4 bytes]
                    token_decimals_int := new(big.Int)
                    token_decimals_int.SetString( _token_decimals[2:], 16 )
                    __token_decimals := token_decimals_int.String()
                    token_decimals_int32, _ := strconv.Atoi( __token_decimals )
                    tx_token_decimals = __token_decimals

                    // token total supply:
                    token_total_supply_int := new(big.Int)
                    token_total_supply_int.SetString( _token_total_supply[2:], 16 )
                    token_total_supply_float := new(big.Float)
                    token_total_supply_float.SetString( token_total_supply_int.String() )
                    __token_total_supply := new(big.Float).Quo(token_total_supply_float, big.NewFloat(math.Pow10(token_decimals_int32)))
                    tx_token_total_supply = fmt.Sprintf( "%.8f", __token_total_supply )

                    //fmt.Println( "token name:", string(__token_name) )
                    //fmt.Println( "token_symbol:", string(__token_symbol) )
                    //fmt.Println( "token_decimals:", __token_decimals )
                    //fmt.Printf( "token_total_supply: %f\n", __token_total_supply )
                    //-----}


                    // SEE:
                    // - https://golang.org/pkg/math/big/
                    // - https://golang.org/pkg/strconv/
                    // - https://goethereumbook.org/account-balance/
                    token_amount_wei_int := new(big.Int)
                    token_amount_wei_int.SetString( tx_token_amount_wei_hex[2:], 16 )
                    //fmt.Println( "erc-20 token amount hex-string to int: ", token_amount_wei_int, "(wei)" )
                    token_amount_wei_float := new(big.Float)
                    token_amount_wei_float.SetString( token_amount_wei_int.String() )
                    token_amount := new(big.Float).Quo(token_amount_wei_float, big.NewFloat(math.Pow10(18)))
                    tx_token_amount = fmt.Sprintf( "%.8f", token_amount )
                    tx_token_amount_wei = token_amount_wei_int.String()
                    //fmt.Printf( "erc-20 token amount: %s (%s)\n", tx_token_amount, tx_token_symbol )
                }

                fmt.Println( "token_contract address =", tx_to )
                fmt.Println( "token_to =", tx_token_to )
                fmt.Println( "token_name =", tx_token_name )
                fmt.Println( "token_symbol =", tx_token_symbol )
                fmt.Println( "token_decimals =", tx_token_decimals )
                fmt.Println( "token_total_supply =", tx_token_total_supply )
                fmt.Println( "token_value_wei =", tx_token_amount_wei, "(wei)" )
                fmt.Println( "token_value_" + tx_token_symbol + " =", tx_token_amount, "(" + tx_token_symbol + ")" )
                fmt.Println()
            }


        } // for ()
    }
}



func main() {
    fmt.Println( "HOST: " + URL )

    /*
    {
        // Ether: balance
        address := "0xe6e55eed00218faef27eed24def9208f3878b333"
        eth_get_balance( address )
    }
    */


    /*
    {
        // Ether: send transaction
        from := "0xe6e55eed00218faef27eed24def9208f3878b333"
        to := "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454"
        amount := "11.1357" // 10 * wei(1e18)
        gas := "70000"
        gasprice := "100"
        eth_send_transaction( from, to, amount, gas, gasprice )
    }
    */


    /*
    {
        // ERC-20: balance
        contract_address := "0x1e57f9561600b269a37437f02ce9da31e5b830ce"
        address := "0xe6e55eed00218faef27eed24def9208f3878b333"
        erc20_get_balance( contract_address, address )
    }
    */


    /*
    {
        // ERC-20: send transaction
        contract_address := "0x1e57f9561600b269a37437f02ce9da31e5b830ce"
        from := "0xe6e55eed00218faef27eed24def9208f3878b333"
        to := "0x8f5b2b7608e3e3a3dc0426c3396420fbf1849454"
        amount := "11.1357" // 10 * wei(1e18)
        gas := "70000"
        gasprice := "100"
        erc20_send_transaction( contract_address,  from, to, amount, gas, gasprice )
    }
    */


    /*
    {
        // Get blocks
        //eth_get_block_by_number
        get_blocks()
    }
    */
}
