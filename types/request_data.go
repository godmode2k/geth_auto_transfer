/* --------------------------------------------------------------
Project:    Ethereum auto-transfer (accounts to specific address(hotwallet))
Purpose:
Author:     Ho-Jung Kim (godmode2k@hotmail.com)
Date:       Since Dec 4, 2020
Filename:   types.go

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
-------------------------------------------------------------- */
package types



//! Header
// ---------------------------------------------------------------

//import (
//)



//! Definition
// --------------------------------------------------------------------

type RequestData struct {
    Jsonrpc string `json:"jsonrpc"`
    Method string `json:"method"`

    Params []interface{} `json:"params"`
    // ["address", "latest"]
    // [{"to": "<contract address>", "data": ""}, "latest"]
    // [{"from": "", "to": "<contract address>", "gas": "", "gasPrice": "", "data": ""}, "latest"]

    Id int `json:"id"`
}

type RequestData_params_erc20 struct {
    // [{"to": "<contract address>", "data": ""}, "latest"]
    To string `json:"to"`
    Data string `json:"data"`
}

type RequestData_params_transaction struct {
    // [{"from": "", "to": "", "value": "<wei>", "gas": "", "gasPrice": "", "data": "", "nonce": ""}, "latest"]
    From string `json:"from"`
    To string `json:"to"`
    Value string `json:"value"`
    Gas string `json:"gas"`
    Gasprice string `json:"gasPrice"`
    //! DO NOT USE
    // except: cancel pending transaction, ...
    //Data string `json:"data"`
    //Nonce string `json:"nonce"`
}

type RequestData_params_erc20_transaction struct {
    // [{"from": "", "to": "<contract address>", "value": "<wei>", "gas": "", "gasPrice": "", "data": ""}, "latest"]
    From string `json:"from"`
    To string `json:"to"`
    Value string `json:"value"`
    Gas string `json:"gas"`
    Gasprice string `json:"gasPrice"`
    Data string `json:"data"`
}

type Result struct {
    Jsonrpc string `json:"jsonrpc"`
    Id int `json:"id"`
    Result string `json:"result"`
}

type Result_block struct {
    Jsonrpc string `json:"jsonrpc"`
    Id int `json:"id"`
    Result interface{} `json:"result"`
}


