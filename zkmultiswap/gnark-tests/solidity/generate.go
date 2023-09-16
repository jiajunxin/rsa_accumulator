package solidity

//go:generate go run contract/main.go
//go:generate solc --evm-version paris --combined-json abi,bin Notuscontract_g16.sol -o abi --overwrite
//go:generate abigen --combined-json abi/combined.json --pkg solidity --out solidity_groth16.go
