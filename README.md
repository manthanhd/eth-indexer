# Eth Indexer
Super simple indexer for eth that indexes token transactions

## Features
1. Indexes simple token transactions
2. Indexes to sqlite database (other db types can be configured easily)

## How to use
Create `config.yml` file with following:
```yaml
eth:
  url: "https://mainnet.infura.io/v3/..."
```
replacing the url with path to your node (or with valid infura key)

Compile the project from source:
```shell
go build -o indexer.exe .\main\main.go
```

Run the application:
```shell
.\indexer.exe
```

The application will create database called `index.database` (if using sqlite).

## Todo
1. Index contract creation transactions
2. Index multi token transfer transactions
3. Index erc20 transferFrom transactions
4. Index all erc20 operations
5. Index all erc721 (NFT) transactions
6. Index all erc1155 (NFT) transactions

## Develop
To generate/update erc20.go from abi
```shell
abigen.exe --abi
.\erc20.abi --pkg erc20 --out erc20\erc20.go
```

## Contribute
