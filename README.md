# Txsender

Ethereum tx sender

## How to use

- build txsender

```sh
go get github.com/wuxiangzhou2010/txsender
```

- prepare keystore

make a keystore and put the keystore in

```sh
mkdir keystore
```

- edit config
  - `endpoint` is a string array contains ipc or rpc endpoint
  - `rate`: indicate how many txs per seconds
  - `chainAmount`: reserved
  - `silent`: print log or not
  - `signedTxBuffer` : for signed txs
  - `rawTxBuffer`: for raw txs
  - `last`: how many seconds to last

default configuration

```json
{
  "endpoint": "http://127.0.0.1:8545",
  "rate": 200,
  "chainAmount": 2,
  "silent": true,
  "signedTxBuffer": 100,
  "rawTxBuffer": 2000,
  "last": 100
}
```

- send txs

```sh
txsender
```

reference:

- [how to create an ethereum account](https://ethereum.stackexchange.com/questions/39900/create-ethereum-account-using-golang)
