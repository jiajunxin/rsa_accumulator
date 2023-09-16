// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PairingMetaData contains all meta data concerning the Pairing contract.
var PairingMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566050600b82828239805160001a6073146043577f4e487b7100000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220810cfc148b7b46391992a313f584aa54823a6351215631e8dd39d2c65c67b1ed64736f6c63430008150033",
}

// PairingABI is the input ABI used to generate the binding from.
// Deprecated: Use PairingMetaData.ABI instead.
var PairingABI = PairingMetaData.ABI

// PairingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PairingMetaData.Bin instead.
var PairingBin = PairingMetaData.Bin

// DeployPairing deploys a new Ethereum contract, binding an instance of Pairing to it.
func DeployPairing(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Pairing, error) {
	parsed, err := PairingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PairingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pairing{PairingCaller: PairingCaller{contract: contract}, PairingTransactor: PairingTransactor{contract: contract}, PairingFilterer: PairingFilterer{contract: contract}}, nil
}

// Pairing is an auto generated Go binding around an Ethereum contract.
type Pairing struct {
	PairingCaller     // Read-only binding to the contract
	PairingTransactor // Write-only binding to the contract
	PairingFilterer   // Log filterer for contract events
}

// PairingCaller is an auto generated read-only Go binding around an Ethereum contract.
type PairingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PairingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PairingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PairingSession struct {
	Contract     *Pairing          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PairingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PairingCallerSession struct {
	Contract *PairingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PairingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PairingTransactorSession struct {
	Contract     *PairingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PairingRaw is an auto generated low-level Go binding around an Ethereum contract.
type PairingRaw struct {
	Contract *Pairing // Generic contract binding to access the raw methods on
}

// PairingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PairingCallerRaw struct {
	Contract *PairingCaller // Generic read-only contract binding to access the raw methods on
}

// PairingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PairingTransactorRaw struct {
	Contract *PairingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPairing creates a new instance of Pairing, bound to a specific deployed contract.
func NewPairing(address common.Address, backend bind.ContractBackend) (*Pairing, error) {
	contract, err := bindPairing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pairing{PairingCaller: PairingCaller{contract: contract}, PairingTransactor: PairingTransactor{contract: contract}, PairingFilterer: PairingFilterer{contract: contract}}, nil
}

// NewPairingCaller creates a new read-only instance of Pairing, bound to a specific deployed contract.
func NewPairingCaller(address common.Address, caller bind.ContractCaller) (*PairingCaller, error) {
	contract, err := bindPairing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PairingCaller{contract: contract}, nil
}

// NewPairingTransactor creates a new write-only instance of Pairing, bound to a specific deployed contract.
func NewPairingTransactor(address common.Address, transactor bind.ContractTransactor) (*PairingTransactor, error) {
	contract, err := bindPairing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PairingTransactor{contract: contract}, nil
}

// NewPairingFilterer creates a new log filterer instance of Pairing, bound to a specific deployed contract.
func NewPairingFilterer(address common.Address, filterer bind.ContractFilterer) (*PairingFilterer, error) {
	contract, err := bindPairing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PairingFilterer{contract: contract}, nil
}

// bindPairing binds a generic wrapper to an already deployed contract.
func bindPairing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PairingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pairing *PairingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pairing.Contract.PairingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pairing *PairingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pairing.Contract.PairingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pairing *PairingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pairing.Contract.PairingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pairing *PairingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pairing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pairing *PairingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pairing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pairing *PairingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pairing.Contract.contract.Transact(opts, method, params...)
}

// VerifierMetaData contains all meta data concerning the Verifier contract.
var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[7]\",\"name\":\"input\",\"type\":\"uint256[7]\"}],\"name\":\"verifyProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"r\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611fa8806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063c894e75714610030575b600080fd5b61004a60048036038101906100459190611802565b610060565b6040516100579190611886565b60405180910390f35b600061006a61139e565b604051806040016040528087600060028110610089576100886118a1565b5b60200201518152602001876001600281106100a7576100a66118a1565b5b6020020151815250816000018190525060405180604001604052806040518060400160405280886000600281106100e1576100e06118a1565b5b60200201516000600281106100f9576100f86118a1565b5b6020020151815260200188600060028110610117576101166118a1565b5b602002015160016002811061012f5761012e6118a1565b5b6020020151815250815260200160405180604001604052808860016002811061015b5761015a6118a1565b5b6020020151600060028110610173576101726118a1565b5b6020020151815260200188600160028110610191576101906118a1565b5b60200201516001600281106101a9576101a86118a1565b5b602002015181525081525081602001819052506040518060400160405280856000600281106101db576101da6118a1565b5b60200201518152602001856001600281106101f9576101f86118a1565b5b602002015181525081604001819052506000610213610735565b90506000604051806040016040528060008152602001600081525090507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478360000151600001511061029a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102919061192d565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4783600001516020015110610304576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102fb90611999565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47836020015160000151600060028110610341576103406118a1565b5b602002015110610386576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161037d90611a05565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478360200151602001516000600281106103c3576103c26118a1565b5b602002015110610408576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103ff90611a71565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47836020015160000151600160028110610445576104446118a1565b5b60200201511061048a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161048190611add565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478360200151602001516001600281106104c7576104c66118a1565b5b60200201511061050c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050390611b49565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4783604001516000015110610576576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161056d90611bb5565b60405180910390fd5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47836040015160200151106105e0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d790611c21565b60405180910390fd5b60005b60078110156106cb577f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000018682600781106106205761061f6118a1565b5b602002015110610665576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161065c90611c8d565b60405180910390fd5b6106b6826106b1856080015160018561067e9190611cdc565b6008811061068f5761068e6118a1565b5b60200201518985600781106106a7576106a66118a1565b5b6020020151610d94565b610e6c565b915080806106c390611d10565b9150506105e3565b506106f28183608001516000600881106106e8576106e76118a1565b5b6020020151610e6c565b90506107286107048460000151610f6a565b84602001518460000151856020015185876040015189604001518960600151611028565b9350505050949350505050565b61073d6113d1565b60405180604001604052807f08d94f29090c5a2514a42e7710bc460da13350c5d953906c888218123dc735e081526020017f11913b7462a18827b73eb8defff4d69f654d29d937e5a8d5bcd67b5cd9c609ad8152508160000181905250604051806040016040528060405180604001604052807f286819b403cf450448c0b5f17ded11aea34763a51897d9a1f0540808cf4c396481526020017f2aa9fe033abb9ed7ed19a93492e53e7ee0a5a68142ffac68641320176e9d7aef815250815260200160405180604001604052807f23a52d495c1765297877f8e5663bfdcbe726fccba1f18c76880d5f4ed8f0024d81526020017f10cd48dbf552bb1d2b71588dbd63a97504e2d2b5d5a93c8dd6853059d8cc0f6d8152508152508160200181905250604051806040016040528060405180604001604052807f2b03e06618909e861c980d411a147039db9519cb388469126a81a1696f38409881526020017f11cba7526c8ad6151e133a63818c17b5ee75e00bc15c2ad66d8a474710b413bd815250815260200160405180604001604052807f2e3f3077f89211208b874f98f90afa76e2f3d63223fdd90e315b340e3f85442581526020017f1038bbe9ea771c0091c05ee4a5d40ff5683c1e8179b5532275a26a0937f916968152508152508160400181905250604051806040016040528060405180604001604052807f2ec580d5cdee6a2df48f463fd8c3191d6170ef3969be2cd196308b9218b4cd9381526020017f01236faf829813621a5bae61429f8827a8f0abd198eac6b26b89ac0503e28b52815250815260200160405180604001604052807f0f580c1189e7bc1abd09b2925bac8a8b378259053f0e2230a93ce5b6821bc49381526020017f178af6471510ee7cbc62410718b73e1f0cf85a6f5a824f0a4a673e8ec5db379d815250815250816060018190525060405180604001604052807f034682702455c9df8e2727c57c00e953779d7321a59c0e9d542c191a200bdf2781526020017f25b3163d06ef6ea221feb2115ff1ce5c998d6db235ae7d24eda7d2e9f3e64fa78152508160800151600060088110610a5657610a556118a1565b5b602002018190525060405180604001604052807f1ebf61a3bfa170bc0467ef07b9dab95bf975823b0f2cd4a7ce0720f29df049dc81526020017f27e04ea03b074880d3d342893f564b81b00cb4108d0618e48ed185632f887b668152508160800151600160088110610acb57610aca6118a1565b5b602002018190525060405180604001604052807f1a95d71e69c7fbe3af71e97d11e8b191b79e6b4414060a8c65d428910e29f9f081526020017f1ab9ddfdde14b26c23df2dae5dc01a43499aeaf60797162eb43396e2577d57e68152508160800151600260088110610b4057610b3f6118a1565b5b602002018190525060405180604001604052807f2933e999f26108e6504790e50504f4ae8699a9303f3fe73cbc9d01e6720eb24a81526020017f01126c2413f0dde513738efcafa452d18022228473d49172f30b434212011dcb8152508160800151600360088110610bb557610bb46118a1565b5b602002018190525060405180604001604052807f1c6b525c45434cd601375cb0e3fe00c3d4612209d153cd151456b15fe4c3195381526020017f2e7701fe54b2b8a7683776372e33787b296802e90c5032782d239a59a6476ca58152508160800151600460088110610c2a57610c296118a1565b5b602002018190525060405180604001604052807f0a0bd2e51d659c86934d0ac4d7e502f88af7e7b2d198bdea058f0c7fbb2457e981526020017f143087a2d0875b17e0c60f8785dcfe053ffa84edaefb5e853bf63547cb1154628152508160800151600560088110610c9f57610c9e6118a1565b5b602002018190525060405180604001604052807f05fc44c911c601af4217173948d60895f5b50b59c013e96921b5a9d1f4ddc41581526020017f1613c26b77228e6e59dbaf7659c0fd7f628f9528d10570321039f1c38b723de98152508160800151600660088110610d1457610d136118a1565b5b602002018190525060405180604001604052807f2a7c271b3e5676fadc08afdd0015d235e36eb28a4e6fdcf8edecb2c8a1b6050e81526020017f22241f97cdefc9d35c9e629592d3cb5f67d36f8788446522e9643d25fc6dbafb8152508160800151600760088110610d8957610d886118a1565b5b602002018190525090565b610d9c61141e565b610da4611438565b836000015181600060038110610dbd57610dbc6118a1565b5b602002018181525050836020015181600160038110610ddf57610dde6118a1565b5b6020020181815250508281600260038110610dfd57610dfc6118a1565b5b602002018181525050600060608360808460076107d05a03fa90508060008103610e2357fe5b5080610e64576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e5b90611da4565b60405180910390fd5b505092915050565b610e7461141e565b610e7c61145a565b836000015181600060048110610e9557610e946118a1565b5b602002018181525050836020015181600160048110610eb757610eb66118a1565b5b602002018181525050826000015181600260048110610ed957610ed86118a1565b5b602002018181525050826020015181600360048110610efb57610efa6118a1565b5b602002018181525050600060608360c08460066107d05a03fa90508060008103610f2157fe5b5080610f62576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f5990611e10565b60405180910390fd5b505092915050565b610f7261141e565b60008260000151148015610f8a575060008260200151145b15610fad5760405180604001604052806000815260200160008152509050611023565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151610ff29190611e5f565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761101d9190611e90565b81525090505b919050565b60008060405180608001604052808b8152602001898152602001878152602001858152509050600060405180608001604052808b815260200189815260200187815260200185815250905060006018905060008167ffffffffffffffff81111561109557611094611539565b5b6040519080825280602002602001820160405280156110c35781602001602082028036833780820191505090505b50905060005b60048110156113015760006006826110e19190611ec4565b90508582600481106110f6576110f56118a1565b5b6020020151600001518360008361110d9190611cdc565b8151811061111e5761111d6118a1565b5b60200260200101818152505085826004811061113d5761113c6118a1565b5b602002015160200151836001836111549190611cdc565b81518110611165576111646118a1565b5b602002602001018181525050848260048110611184576111836118a1565b5b6020020151600001516000600281106111a05761119f6118a1565b5b6020020151836002836111b39190611cdc565b815181106111c4576111c36118a1565b5b6020026020010181815250508482600481106111e3576111e26118a1565b5b6020020151600001516001600281106111ff576111fe6118a1565b5b6020020151836003836112129190611cdc565b81518110611223576112226118a1565b5b602002602001018181525050848260048110611242576112416118a1565b5b60200201516020015160006002811061125e5761125d6118a1565b5b6020020151836004836112719190611cdc565b81518110611282576112816118a1565b5b6020026020010181815250508482600481106112a1576112a06118a1565b5b6020020151602001516001600281106112bd576112bc6118a1565b5b6020020151836005836112d09190611cdc565b815181106112e1576112e06118a1565b5b6020026020010181815250505080806112f990611d10565b9150506110c9565b5061130a61147c565b6000602082602086026020860160086107d05a03fa9050806000810361132c57fe5b508061136d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161136490611f52565b60405180910390fd5b600082600060018110611383576113826118a1565b5b60200201511415965050505050505098975050505050505050565b60405180606001604052806113b161141e565b81526020016113be61149e565b81526020016113cb61141e565b81525090565b6040518060a001604052806113e461141e565b81526020016113f161149e565b81526020016113fe61149e565b815260200161140b61149e565b81526020016114186114c4565b81525090565b604051806040016040528060008152602001600081525090565b6040518060600160405280600390602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b6040518060200160405280600190602082028036833780820191505090505090565b60405180604001604052806114b16114f2565b81526020016114be6114f2565b81525090565b6040518061010001604052806008905b6114dc61141e565b8152602001906001900390816114d45790505090565b6040518060400160405280600290602082028036833780820191505090505090565b6000604051905090565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61157182611528565b810181811067ffffffffffffffff821117156115905761158f611539565b5b80604052505050565b60006115a3611514565b90506115af8282611568565b919050565b600067ffffffffffffffff8211156115cf576115ce611539565b5b602082029050919050565b600080fd5b6000819050919050565b6115f2816115df565b81146115fd57600080fd5b50565b60008135905061160f816115e9565b92915050565b6000611628611623846115b4565b611599565b90508060208402830185811115611642576116416115da565b5b835b8181101561166b57806116578882611600565b845260208401935050602081019050611644565b5050509392505050565b600082601f83011261168a57611689611523565b5b6002611697848285611615565b91505092915050565b600067ffffffffffffffff8211156116bb576116ba611539565b5b602082029050919050565b60006116d96116d4846116a0565b611599565b905080604084028301858111156116f3576116f26115da565b5b835b8181101561171c57806117088882611675565b8452602084019350506040810190506116f5565b5050509392505050565b600082601f83011261173b5761173a611523565b5b60026117488482856116c6565b91505092915050565b600067ffffffffffffffff82111561176c5761176b611539565b5b602082029050919050565b600061178a61178584611751565b611599565b905080602084028301858111156117a4576117a36115da565b5b835b818110156117cd57806117b98882611600565b8452602084019350506020810190506117a6565b5050509392505050565b600082601f8301126117ec576117eb611523565b5b60076117f9848285611777565b91505092915050565b6000806000806101e0858703121561181d5761181c61151e565b5b600061182b87828801611675565b945050604061183c87828801611726565b93505060c061184d87828801611675565b92505061010061185f878288016117d7565b91505092959194509250565b60008115159050919050565b6118808161186b565b82525050565b600060208201905061189b6000830184611877565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082825260208201905092915050565b7f76657269666965722d61582d6774652d7072696d652d71000000000000000000600082015250565b60006119176017836118d0565b9150611922826118e1565b602082019050919050565b600060208201905081810360008301526119468161190a565b9050919050565b7f76657269666965722d61592d6774652d7072696d652d71000000000000000000600082015250565b60006119836017836118d0565b915061198e8261194d565b602082019050919050565b600060208201905081810360008301526119b281611976565b9050919050565b7f76657269666965722d6258302d6774652d7072696d652d710000000000000000600082015250565b60006119ef6018836118d0565b91506119fa826119b9565b602082019050919050565b60006020820190508181036000830152611a1e816119e2565b9050919050565b7f76657269666965722d6259302d6774652d7072696d652d710000000000000000600082015250565b6000611a5b6018836118d0565b9150611a6682611a25565b602082019050919050565b60006020820190508181036000830152611a8a81611a4e565b9050919050565b7f76657269666965722d6258312d6774652d7072696d652d710000000000000000600082015250565b6000611ac76018836118d0565b9150611ad282611a91565b602082019050919050565b60006020820190508181036000830152611af681611aba565b9050919050565b7f76657269666965722d6259312d6774652d7072696d652d710000000000000000600082015250565b6000611b336018836118d0565b9150611b3e82611afd565b602082019050919050565b60006020820190508181036000830152611b6281611b26565b9050919050565b7f76657269666965722d63582d6774652d7072696d652d71000000000000000000600082015250565b6000611b9f6017836118d0565b9150611baa82611b69565b602082019050919050565b60006020820190508181036000830152611bce81611b92565b9050919050565b7f76657269666965722d63592d6774652d7072696d652d71000000000000000000600082015250565b6000611c0b6017836118d0565b9150611c1682611bd5565b602082019050919050565b60006020820190508181036000830152611c3a81611bfe565b9050919050565b7f76657269666965722d6774652d736e61726b2d7363616c61722d6669656c6400600082015250565b6000611c77601f836118d0565b9150611c8282611c41565b602082019050919050565b60006020820190508181036000830152611ca681611c6a565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611ce7826115df565b9150611cf2836115df565b9250828201905080821115611d0a57611d09611cad565b5b92915050565b6000611d1b826115df565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611d4d57611d4c611cad565b5b600182019050919050565b7f70616972696e672d6d756c2d6661696c65640000000000000000000000000000600082015250565b6000611d8e6012836118d0565b9150611d9982611d58565b602082019050919050565b60006020820190508181036000830152611dbd81611d81565b9050919050565b7f70616972696e672d6164642d6661696c65640000000000000000000000000000600082015250565b6000611dfa6012836118d0565b9150611e0582611dc4565b602082019050919050565b60006020820190508181036000830152611e2981611ded565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611e6a826115df565b9150611e75836115df565b925082611e8557611e84611e30565b5b828206905092915050565b6000611e9b826115df565b9150611ea6836115df565b9250828203905081811115611ebe57611ebd611cad565b5b92915050565b6000611ecf826115df565b9150611eda836115df565b9250828202611ee8816115df565b91508282048414831517611eff57611efe611cad565b5b5092915050565b7f70616972696e672d6f70636f64652d6661696c65640000000000000000000000600082015250565b6000611f3c6015836118d0565b9150611f4782611f06565b602082019050919050565b60006020820190508181036000830152611f6b81611f2f565b905091905056fea26469706673582212203040049f10fa166032bd8ec711f8b85087c0364df5615a22a5f0d94d61ef916b64736f6c63430008150033",
}

// VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierMetaData.ABI instead.
var VerifierABI = VerifierMetaData.ABI

// VerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierMetaData.Bin instead.
var VerifierBin = VerifierMetaData.Bin

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// VerifyProof is a free data retrieval call binding the contract method 0xc894e757.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[7] input) view returns(bool r)
func (_Verifier *VerifierCaller) VerifyProof(opts *bind.CallOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [7]*big.Int) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "verifyProof", a, b, c, input)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0xc894e757.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[7] input) view returns(bool r)
func (_Verifier *VerifierSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [7]*big.Int) (bool, error) {
	return _Verifier.Contract.VerifyProof(&_Verifier.CallOpts, a, b, c, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0xc894e757.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[7] input) view returns(bool r)
func (_Verifier *VerifierCallerSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [7]*big.Int) (bool, error) {
	return _Verifier.Contract.VerifyProof(&_Verifier.CallOpts, a, b, c, input)
}
