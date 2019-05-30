// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package seed

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DhashABI is the input ABI used to generate the binding from.
const DhashABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"hotList\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"latest\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"}],\"name\":\"getHash\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_list\",\"type\":\"string\"}],\"name\":\"updateHotList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"},{\"name\":\"_hash\",\"type\":\"string\"},{\"name\":\"_code\",\"type\":\"uint256\"}],\"name\":\"updateVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"versionTable\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getHotList\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLatest\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"string\"}],\"name\":\"versionHash\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// DhashBin is the compiled bytecode used for deploying new contracts.
const DhashBin = `0x608060405234801561001057600080fd5b50600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36000600155610bf88061006e6000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80639f2c6d651161008c578063b417985d11610066578063b417985d14610445578063c36af4601461044d578063d299862a146104d4578063f2fde38b14610578576100cf565b80639f2c6d6514610259578063aa0192cc146102fd578063b2242d2914610428576100cf565b806320c2fff7146100d457806352bfe789146101515780635b6beeb91461016b578063715018a61461020f5780638da5cb5b146102195780638f32d59b1461023d575b600080fd5b6100dc61059e565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101165781810151838201526020016100fe565b50505050905090810190601f1680156101435780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61015961062c565b60408051918252519081900360200190f35b6100dc6004803603602081101561018157600080fd5b810190602081018135600160201b81111561019b57600080fd5b8201836020820111156101ad57600080fd5b803590602001918460018302840111600160201b831117156101ce57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610632945050505050565b610217610726565b005b610221610781565b604080516001600160a01b039092168252519081900360200190f35b610245610791565b604080519115158252519081900360200190f35b6102176004803603602081101561026f57600080fd5b810190602081018135600160201b81111561028957600080fd5b82018360208201111561029b57600080fd5b803590602001918460018302840111600160201b831117156102bc57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506107a2945050505050565b6102176004803603606081101561031357600080fd5b810190602081018135600160201b81111561032d57600080fd5b82018360208201111561033f57600080fd5b803590602001918460018302840111600160201b8311171561036057600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b8111156103b257600080fd5b8201836020820111156103c457600080fd5b803590602001918460018302840111600160201b831117156103e557600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050913592506107ca915050565b6100dc6004803603602081101561043e57600080fd5b5035610888565b6100dc6108ee565b610455610984565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610498578181015183820152602001610480565b50505050905090810190601f1680156104c55780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6100dc600480360360208110156104ea57600080fd5b810190602081018135600160201b81111561050457600080fd5b82018360208201111561051657600080fd5b803590602001918460018302840111600160201b8311171561053757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610a2c945050505050565b6102176004803603602081101561058e57600080fd5b50356001600160a01b0316610aa0565b6004805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156106245780601f106105f957610100808354040283529160200191610624565b820191906000526020600020905b81548152906001019060200180831161060757829003601f168201915b505050505081565b60015481565b60606003826040518082805190602001908083835b602083106106665780518252601f199092019160209182019101610647565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188208054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561071a5780601f106106ef5761010080835404028352916020019161071a565b820191906000526020600020905b8154815290600101906020018083116106fd57829003601f168201915b50505050509050919050565b61072e610791565b61073757600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b03165b90565b6000546001600160a01b0316331490565b6107aa610791565b6107b357600080fd5b80516107c6906004906020840190610b2b565b5050565b6107d2610791565b6107db57600080fd5b60015481116107e957600080fd5b60018190556000818152600260209081526040909120845161080d92860190610b2b565b50816003846040518082805190602001908083835b602083106108415780518252601f199092019160209182019101610822565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810190932084516108829591949190910192509050610b2b565b50505050565b600260208181526000928352604092839020805484516001821615610100026000190190911693909304601f81018390048302840183019094528383529192908301828280156106245780601f106105f957610100808354040283529160200191610624565b60048054604080516020601f600260001961010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561097a5780601f1061094f5761010080835404028352916020019161097a565b820191906000526020600020905b81548152906001019060200180831161095d57829003601f168201915b5050505050905090565b600180546000818152600260208181526040808420805482519781161561010002600019011693909304601f81018390048302870183019091528086529294606094939091839190830182828015610a1d5780601f106109f257610100808354040283529160200191610a1d565b820191906000526020600020905b815481529060010190602001808311610a0057829003601f168201915b50505050509050915091509091565b805160208183018101805160038252928201938201939093209190925280546040805160026001841615610100026000190190931692909204601f810185900485028301850190915280825290928301828280156106245780601f106105f957610100808354040283529160200191610624565b610aa8610791565b610ab157600080fd5b610aba81610abd565b50565b6001600160a01b038116610ad057600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b6c57805160ff1916838001178555610b99565b82800160010185558215610b99579182015b82811115610b99578251825591602001919060010190610b7e565b50610ba5929150610ba9565b5090565b61078e91905b80821115610ba55760008155600101610baf56fea265627a7a72305820c248ec223cd1954d3082b046a8f17379f2a71b5cd3695fb14c01c2137bdfba0b64736f6c63430005090032`

// DeployDhash deploys a new Ethereum contract, binding an instance of Dhash to it.
func DeployDhash(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Dhash, error) {
	parsed, err := abi.JSON(strings.NewReader(DhashABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DhashBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Dhash{DhashCaller: DhashCaller{contract: contract}, DhashTransactor: DhashTransactor{contract: contract}, DhashFilterer: DhashFilterer{contract: contract}}, nil
}

// Dhash is an auto generated Go binding around an Ethereum contract.
type Dhash struct {
	DhashCaller     // Read-only binding to the contract
	DhashTransactor // Write-only binding to the contract
	DhashFilterer   // Log filterer for contract events
}

// DhashCaller is an auto generated read-only Go binding around an Ethereum contract.
type DhashCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DhashTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DhashTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DhashFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DhashFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DhashSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DhashSession struct {
	Contract     *Dhash            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DhashCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DhashCallerSession struct {
	Contract *DhashCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DhashTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DhashTransactorSession struct {
	Contract     *DhashTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DhashRaw is an auto generated low-level Go binding around an Ethereum contract.
type DhashRaw struct {
	Contract *Dhash // Generic contract binding to access the raw methods on
}

// DhashCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DhashCallerRaw struct {
	Contract *DhashCaller // Generic read-only contract binding to access the raw methods on
}

// DhashTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DhashTransactorRaw struct {
	Contract *DhashTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDhash creates a new instance of Dhash, bound to a specific deployed contract.
func NewDhash(address common.Address, backend bind.ContractBackend) (*Dhash, error) {
	contract, err := bindDhash(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dhash{DhashCaller: DhashCaller{contract: contract}, DhashTransactor: DhashTransactor{contract: contract}, DhashFilterer: DhashFilterer{contract: contract}}, nil
}

// NewDhashCaller creates a new read-only instance of Dhash, bound to a specific deployed contract.
func NewDhashCaller(address common.Address, caller bind.ContractCaller) (*DhashCaller, error) {
	contract, err := bindDhash(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DhashCaller{contract: contract}, nil
}

// NewDhashTransactor creates a new write-only instance of Dhash, bound to a specific deployed contract.
func NewDhashTransactor(address common.Address, transactor bind.ContractTransactor) (*DhashTransactor, error) {
	contract, err := bindDhash(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DhashTransactor{contract: contract}, nil
}

// NewDhashFilterer creates a new log filterer instance of Dhash, bound to a specific deployed contract.
func NewDhashFilterer(address common.Address, filterer bind.ContractFilterer) (*DhashFilterer, error) {
	contract, err := bindDhash(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DhashFilterer{contract: contract}, nil
}

// bindDhash binds a generic wrapper to an already deployed contract.
func bindDhash(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DhashABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dhash *DhashRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Dhash.Contract.DhashCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dhash *DhashRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dhash.Contract.DhashTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dhash *DhashRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dhash.Contract.DhashTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dhash *DhashCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Dhash.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dhash *DhashTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dhash.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dhash *DhashTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dhash.Contract.contract.Transact(opts, method, params...)
}

// GetHash is a free data retrieval call binding the contract method 0x5b6beeb9.
//
// Solidity: function getHash(string _version) constant returns(string)
func (_Dhash *DhashCaller) GetHash(opts *bind.CallOpts, _version string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "getHash", _version)
	return *ret0, err
}

// GetHash is a free data retrieval call binding the contract method 0x5b6beeb9.
//
// Solidity: function getHash(string _version) constant returns(string)
func (_Dhash *DhashSession) GetHash(_version string) (string, error) {
	return _Dhash.Contract.GetHash(&_Dhash.CallOpts, _version)
}

// GetHash is a free data retrieval call binding the contract method 0x5b6beeb9.
//
// Solidity: function getHash(string _version) constant returns(string)
func (_Dhash *DhashCallerSession) GetHash(_version string) (string, error) {
	return _Dhash.Contract.GetHash(&_Dhash.CallOpts, _version)
}

// GetHotList is a free data retrieval call binding the contract method 0xb417985d.
//
// Solidity: function getHotList() constant returns(string)
func (_Dhash *DhashCaller) GetHotList(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "getHotList")
	return *ret0, err
}

// GetHotList is a free data retrieval call binding the contract method 0xb417985d.
//
// Solidity: function getHotList() constant returns(string)
func (_Dhash *DhashSession) GetHotList() (string, error) {
	return _Dhash.Contract.GetHotList(&_Dhash.CallOpts)
}

// GetHotList is a free data retrieval call binding the contract method 0xb417985d.
//
// Solidity: function getHotList() constant returns(string)
func (_Dhash *DhashCallerSession) GetHotList() (string, error) {
	return _Dhash.Contract.GetHotList(&_Dhash.CallOpts)
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() constant returns(uint256, string)
func (_Dhash *DhashCaller) GetLatest(opts *bind.CallOpts) (*big.Int, string, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Dhash.contract.Call(opts, out, "getLatest")
	return *ret0, *ret1, err
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() constant returns(uint256, string)
func (_Dhash *DhashSession) GetLatest() (*big.Int, string, error) {
	return _Dhash.Contract.GetLatest(&_Dhash.CallOpts)
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() constant returns(uint256, string)
func (_Dhash *DhashCallerSession) GetLatest() (*big.Int, string, error) {
	return _Dhash.Contract.GetLatest(&_Dhash.CallOpts)
}

// HotList is a free data retrieval call binding the contract method 0x20c2fff7.
//
// Solidity: function hotList() constant returns(string)
func (_Dhash *DhashCaller) HotList(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "hotList")
	return *ret0, err
}

// HotList is a free data retrieval call binding the contract method 0x20c2fff7.
//
// Solidity: function hotList() constant returns(string)
func (_Dhash *DhashSession) HotList() (string, error) {
	return _Dhash.Contract.HotList(&_Dhash.CallOpts)
}

// HotList is a free data retrieval call binding the contract method 0x20c2fff7.
//
// Solidity: function hotList() constant returns(string)
func (_Dhash *DhashCallerSession) HotList() (string, error) {
	return _Dhash.Contract.HotList(&_Dhash.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Dhash *DhashCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Dhash *DhashSession) IsOwner() (bool, error) {
	return _Dhash.Contract.IsOwner(&_Dhash.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Dhash *DhashCallerSession) IsOwner() (bool, error) {
	return _Dhash.Contract.IsOwner(&_Dhash.CallOpts)
}

// Latest is a free data retrieval call binding the contract method 0x52bfe789.
//
// Solidity: function latest() constant returns(uint256)
func (_Dhash *DhashCaller) Latest(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "latest")
	return *ret0, err
}

// Latest is a free data retrieval call binding the contract method 0x52bfe789.
//
// Solidity: function latest() constant returns(uint256)
func (_Dhash *DhashSession) Latest() (*big.Int, error) {
	return _Dhash.Contract.Latest(&_Dhash.CallOpts)
}

// Latest is a free data retrieval call binding the contract method 0x52bfe789.
//
// Solidity: function latest() constant returns(uint256)
func (_Dhash *DhashCallerSession) Latest() (*big.Int, error) {
	return _Dhash.Contract.Latest(&_Dhash.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Dhash *DhashCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Dhash *DhashSession) Owner() (common.Address, error) {
	return _Dhash.Contract.Owner(&_Dhash.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Dhash *DhashCallerSession) Owner() (common.Address, error) {
	return _Dhash.Contract.Owner(&_Dhash.CallOpts)
}

// VersionHash is a free data retrieval call binding the contract method 0xd299862a.
//
// Solidity: function versionHash(string ) constant returns(string)
func (_Dhash *DhashCaller) VersionHash(opts *bind.CallOpts, arg0 string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "versionHash", arg0)
	return *ret0, err
}

// VersionHash is a free data retrieval call binding the contract method 0xd299862a.
//
// Solidity: function versionHash(string ) constant returns(string)
func (_Dhash *DhashSession) VersionHash(arg0 string) (string, error) {
	return _Dhash.Contract.VersionHash(&_Dhash.CallOpts, arg0)
}

// VersionHash is a free data retrieval call binding the contract method 0xd299862a.
//
// Solidity: function versionHash(string ) constant returns(string)
func (_Dhash *DhashCallerSession) VersionHash(arg0 string) (string, error) {
	return _Dhash.Contract.VersionHash(&_Dhash.CallOpts, arg0)
}

// VersionTable is a free data retrieval call binding the contract method 0xb2242d29.
//
// Solidity: function versionTable(uint256 ) constant returns(string)
func (_Dhash *DhashCaller) VersionTable(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Dhash.contract.Call(opts, out, "versionTable", arg0)
	return *ret0, err
}

// VersionTable is a free data retrieval call binding the contract method 0xb2242d29.
//
// Solidity: function versionTable(uint256 ) constant returns(string)
func (_Dhash *DhashSession) VersionTable(arg0 *big.Int) (string, error) {
	return _Dhash.Contract.VersionTable(&_Dhash.CallOpts, arg0)
}

// VersionTable is a free data retrieval call binding the contract method 0xb2242d29.
//
// Solidity: function versionTable(uint256 ) constant returns(string)
func (_Dhash *DhashCallerSession) VersionTable(arg0 *big.Int) (string, error) {
	return _Dhash.Contract.VersionTable(&_Dhash.CallOpts, arg0)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dhash *DhashTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dhash.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dhash *DhashSession) RenounceOwnership() (*types.Transaction, error) {
	return _Dhash.Contract.RenounceOwnership(&_Dhash.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dhash *DhashTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Dhash.Contract.RenounceOwnership(&_Dhash.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dhash *DhashTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Dhash.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dhash *DhashSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Dhash.Contract.TransferOwnership(&_Dhash.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dhash *DhashTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Dhash.Contract.TransferOwnership(&_Dhash.TransactOpts, newOwner)
}

// UpdateHotList is a paid mutator transaction binding the contract method 0x9f2c6d65.
//
// Solidity: function updateHotList(string _list) returns()
func (_Dhash *DhashTransactor) UpdateHotList(opts *bind.TransactOpts, _list string) (*types.Transaction, error) {
	return _Dhash.contract.Transact(opts, "updateHotList", _list)
}

// UpdateHotList is a paid mutator transaction binding the contract method 0x9f2c6d65.
//
// Solidity: function updateHotList(string _list) returns()
func (_Dhash *DhashSession) UpdateHotList(_list string) (*types.Transaction, error) {
	return _Dhash.Contract.UpdateHotList(&_Dhash.TransactOpts, _list)
}

// UpdateHotList is a paid mutator transaction binding the contract method 0x9f2c6d65.
//
// Solidity: function updateHotList(string _list) returns()
func (_Dhash *DhashTransactorSession) UpdateHotList(_list string) (*types.Transaction, error) {
	return _Dhash.Contract.UpdateHotList(&_Dhash.TransactOpts, _list)
}

// UpdateVersion is a paid mutator transaction binding the contract method 0xaa0192cc.
//
// Solidity: function updateVersion(string _version, string _hash, uint256 _code) returns()
func (_Dhash *DhashTransactor) UpdateVersion(opts *bind.TransactOpts, _version string, _hash string, _code *big.Int) (*types.Transaction, error) {
	return _Dhash.contract.Transact(opts, "updateVersion", _version, _hash, _code)
}

// UpdateVersion is a paid mutator transaction binding the contract method 0xaa0192cc.
//
// Solidity: function updateVersion(string _version, string _hash, uint256 _code) returns()
func (_Dhash *DhashSession) UpdateVersion(_version string, _hash string, _code *big.Int) (*types.Transaction, error) {
	return _Dhash.Contract.UpdateVersion(&_Dhash.TransactOpts, _version, _hash, _code)
}

// UpdateVersion is a paid mutator transaction binding the contract method 0xaa0192cc.
//
// Solidity: function updateVersion(string _version, string _hash, uint256 _code) returns()
func (_Dhash *DhashTransactorSession) UpdateVersion(_version string, _hash string, _code *big.Int) (*types.Transaction, error) {
	return _Dhash.Contract.UpdateVersion(&_Dhash.TransactOpts, _version, _hash, _code)
}

// DhashOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Dhash contract.
type DhashOwnershipTransferredIterator struct {
	Event *DhashOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DhashOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DhashOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DhashOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DhashOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DhashOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DhashOwnershipTransferred represents a OwnershipTransferred event raised by the Dhash contract.
type DhashOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Dhash *DhashFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DhashOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Dhash.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DhashOwnershipTransferredIterator{contract: _Dhash.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Dhash *DhashFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DhashOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Dhash.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DhashOwnershipTransferred)
				if err := _Dhash.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBin is the compiled bytecode used for deploying new contracts.
const OwnableBin = `0x`

// DeployOwnable deploys a new Ethereum contract, binding an instance of Ownable to it.
func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableCallerSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
