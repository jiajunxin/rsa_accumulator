module github.com/jiajunxin/rsa_accumulator

go 1.19

require (
	github.com/consensys/gnark v0.7.0
	github.com/consensys/gnark-crypto v0.7.0
	github.com/iden3/go-iden3-crypto v0.0.13
	github.com/jiajunxin/multiexp v0.1.2
	github.com/remyoudompheng/bigfft v0.0.0-20220927061507-ef77025ab5aa
	github.com/txaty/go-bigcomplex v0.1.6
	lukechampine.com/frand v1.4.2
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/stretchr/testify v1.7.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/consensys/gnark => github.com/bnb-chain/gnark v0.7.1-0.20230203031713-0d81c67d080a
	github.com/consensys/gnark-crypto => github.com/bnb-chain/gnark-crypto v0.7.1-0.20230203031630-7c643ad11891
)
