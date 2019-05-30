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

// BangumiDataABI is the input ABI used to generate the binding from.
const BangumiDataABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryHash\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryTotalSeason\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"querySharpness\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"querySarmid\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"querySwarmAdd\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryTotalEpisode\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryPoster\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryVideoType\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryRole\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bangumi\",\"type\":\"string\"},{\"name\":\"_poster\",\"type\":\"string\"},{\"name\":\"_role\",\"type\":\"string\"},{\"name\":\"_hash\",\"type\":\"string\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_sharpness\",\"type\":\"string\"},{\"name\":\"_episode\",\"type\":\"string\"},{\"name\":\"_totalEpisode\",\"type\":\"string\"},{\"name\":\"_season\",\"type\":\"string\"},{\"name\":\"_videoType\",\"type\":\"string\"},{\"name\":\"_swarmID\",\"type\":\"string\"},{\"name\":\"_swarmAdd\",\"type\":\"string\"}],\"name\":\"_infoInput\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sn\",\"type\":\"string\"}],\"name\":\"queryEpisode\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// BangumiDataBin is the compiled bytecode used for deploying new contracts.
const BangumiDataBin = `0x60806040819052600080546001600160a01b03191633178082556001600160a01b0316917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3611bff806100576000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638da5cb5b11610097578063d5de372311610066578063d5de372314610788578063e08805901461082c578063f2fde38b14610e88578063fd4e902514610eae57610100565b80638da5cb5b146106005780638f32d59b1461062457806395463ae014610640578063a496a427146106e457610100565b8063557bafc0116100d3578063557bafc01461040a578063702e8170146104ae578063715018a6146105525780638bff2ad71461055c57610100565b806314d407061461010557806316f16be21461021e5780631fccd5ff146102c25780634b9722f514610366575b600080fd5b6101a96004803603602081101561011b57600080fd5b810190602081018135600160201b81111561013557600080fd5b82018360208201111561014757600080fd5b803590602001918460018302840111600160201b8311171561016857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610f52945050505050565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101e35781810151838201526020016101cb565b50505050905090810190601f1680156102105780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101a96004803603602081101561023457600080fd5b810190602081018135600160201b81111561024e57600080fd5b82018360208201111561026057600080fd5b803590602001918460018302840111600160201b8311171561028157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611049945050505050565b6101a9600480360360208110156102d857600080fd5b810190602081018135600160201b8111156102f257600080fd5b82018360208201111561030457600080fd5b803590602001918460018302840111600160201b8311171561032557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611109945050505050565b6101a96004803603602081101561037c57600080fd5b810190602081018135600160201b81111561039657600080fd5b8201836020820111156103a857600080fd5b803590602001918460018302840111600160201b831117156103c957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506111c9945050505050565b6101a96004803603602081101561042057600080fd5b810190602081018135600160201b81111561043a57600080fd5b82018360208201111561044c57600080fd5b803590602001918460018302840111600160201b8311171561046d57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611289945050505050565b6101a9600480360360208110156104c457600080fd5b810190602081018135600160201b8111156104de57600080fd5b8201836020820111156104f057600080fd5b803590602001918460018302840111600160201b8311171561051157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611349945050505050565b61055a611409565b005b6101a96004803603602081101561057257600080fd5b810190602081018135600160201b81111561058c57600080fd5b82018360208201111561059e57600080fd5b803590602001918460018302840111600160201b831117156105bf57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611464945050505050565b610608611521565b604080516001600160a01b039092168252519081900360200190f35b61062c611531565b604080519115158252519081900360200190f35b6101a96004803603602081101561065657600080fd5b810190602081018135600160201b81111561067057600080fd5b82018360208201111561068257600080fd5b803590602001918460018302840111600160201b831117156106a357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611542945050505050565b6101a9600480360360208110156106fa57600080fd5b810190602081018135600160201b81111561071457600080fd5b82018360208201111561072657600080fd5b803590602001918460018302840111600160201b8311171561074757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611602945050505050565b6101a96004803603602081101561079e57600080fd5b810190602081018135600160201b8111156107b857600080fd5b8201836020820111156107ca57600080fd5b803590602001918460018302840111600160201b831117156107eb57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506116c3945050505050565b61055a600480360361018081101561084357600080fd5b810190602081018135600160201b81111561085d57600080fd5b82018360208201111561086f57600080fd5b803590602001918460018302840111600160201b8311171561089057600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b8111156108e257600080fd5b8201836020820111156108f457600080fd5b803590602001918460018302840111600160201b8311171561091557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561096757600080fd5b82018360208201111561097957600080fd5b803590602001918460018302840111600160201b8311171561099a57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b8111156109ec57600080fd5b8201836020820111156109fe57600080fd5b803590602001918460018302840111600160201b83111715610a1f57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610a7157600080fd5b820183602082011115610a8357600080fd5b803590602001918460018302840111600160201b83111715610aa457600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610af657600080fd5b820183602082011115610b0857600080fd5b803590602001918460018302840111600160201b83111715610b2957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610b7b57600080fd5b820183602082011115610b8d57600080fd5b803590602001918460018302840111600160201b83111715610bae57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610c0057600080fd5b820183602082011115610c1257600080fd5b803590602001918460018302840111600160201b83111715610c3357600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610c8557600080fd5b820183602082011115610c9757600080fd5b803590602001918460018302840111600160201b83111715610cb857600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610d0a57600080fd5b820183602082011115610d1c57600080fd5b803590602001918460018302840111600160201b83111715610d3d57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610d8f57600080fd5b820183602082011115610da157600080fd5b803590602001918460018302840111600160201b83111715610dc257600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b811115610e1457600080fd5b820183602082011115610e2657600080fd5b803590602001918460018302840111600160201b83111715610e4757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611783945050505050565b61055a60048036036020811015610e9e57600080fd5b50356001600160a01b0316611803565b6101a960048036036020811015610ec457600080fd5b810190602081018135600160201b811115610ede57600080fd5b820183602082011115610ef057600080fd5b803590602001918460018302840111600160201b83111715610f1157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611820945050505050565b60606001826040518082805190602001908083835b60208310610f865780518252601f199092019160209182019101610f67565b518151600019602094850361010090810a8201928316921993909316919091179092529490920196875260408051978890038201882060029081018054601f60018216159098029095019094160494850182900482028801820190528387529094509192505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b820191906000526020600020905b81548152906001019060200180831161102057829003601f168201915b50505050509050919050565b60606001826040518082805190602001908083835b6020831061107d5780518252601f19909201916020918201910161105e565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206007018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b6020831061113d5780518252601f19909201916020918201910161111e565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206004018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b602083106111fd5780518252601f1990920191602091820191016111de565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206009018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b602083106112bd5780518252601f19909201916020918201910161129e565b518151600019602094850361010090810a82019283169219939093169190911790925294909201968752604080519788900382018820600a018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b6020831061137d5780518252601f19909201916020918201910161135e565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206006018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b611411611531565b61141a57600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b60606001826040518082805190602001908083835b602083106114985780518252601f199092019160209182019101611479565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188208054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b6000546001600160a01b03165b90565b6000546001600160a01b0316331490565b60606001826040518082805190602001908083835b602083106115765780518252601f199092019160209182019101611557565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206008018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b602083106116365780518252601f199092019160209182019101611617565b518151600019602094850361010090810a8201928316921993909316919091179092529490920196875260408051978890038201882060019081018054601f6002938216159098029095019094160494850182900482028801820190528387529094509192505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b60606001826040518082805190602001908083835b602083106116f75780518252601f1990920191602091820191016116d8565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206003018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b61178b611531565b61179457600080fd5b61179c611ae1565b6040518061016001604052808d81526020018c81526020018b81526020018a81526020018981526020018881526020018781526020018681526020018581526020018481526020018381525090506117f48d826118e0565b50505050505050505050505050565b61180b611531565b61181457600080fd5b61181d81611a73565b50565b60606001826040518082805190602001908083835b602083106118545780518252601f199092019160209182019101611835565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188206005018054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561103d5780601f106110125761010080835404028352916020019161103d565b806001836040518082805190602001908083835b602083106119135780518252601f1990920191602091820191016118f4565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810190932084518051919461195494508593500190611b3b565b50602082810151805161196d9260018501920190611b3b565b5060408201518051611989916002840191602090910190611b3b565b50606082015180516119a5916003840191602090910190611b3b565b50608082015180516119c1916004840191602090910190611b3b565b5060a082015180516119dd916005840191602090910190611b3b565b5060c082015180516119f9916006840191602090910190611b3b565b5060e08201518051611a15916007840191602090910190611b3b565b506101008201518051611a32916008840191602090910190611b3b565b506101208201518051611a4f916009840191602090910190611b3b565b506101408201518051611a6c91600a840191602090910190611b3b565b5050505050565b6001600160a01b038116611a8657600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b60405180610160016040528060608152602001606081526020016060815260200160608152602001606081526020016060815260200160608152602001606081526020016060815260200160608152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611b7c57805160ff1916838001178555611ba9565b82800160010185558215611ba9579182015b82811115611ba9578251825591602001919060010190611b8e565b50611bb5929150611bb9565b5090565b61152e91905b80821115611bb55760008155600101611bbf56fea165627a7a72305820bcdff6f00c4cb9271c8f9042c83be1ff5e1799500867b4e5b53310233a75c20d0029`

// DeployBangumiData deploys a new Ethereum contract, binding an instance of BangumiData to it.
func DeployBangumiData(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BangumiData, error) {
	parsed, err := abi.JSON(strings.NewReader(BangumiDataABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BangumiDataBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BangumiData{BangumiDataCaller: BangumiDataCaller{contract: contract}, BangumiDataTransactor: BangumiDataTransactor{contract: contract}, BangumiDataFilterer: BangumiDataFilterer{contract: contract}}, nil
}

// BangumiData is an auto generated Go binding around an Ethereum contract.
type BangumiData struct {
	BangumiDataCaller     // Read-only binding to the contract
	BangumiDataTransactor // Write-only binding to the contract
	BangumiDataFilterer   // Log filterer for contract events
}

// BangumiDataCaller is an auto generated read-only Go binding around an Ethereum contract.
type BangumiDataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BangumiDataTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BangumiDataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BangumiDataFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BangumiDataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BangumiDataSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BangumiDataSession struct {
	Contract     *BangumiData      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BangumiDataCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BangumiDataCallerSession struct {
	Contract *BangumiDataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BangumiDataTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BangumiDataTransactorSession struct {
	Contract     *BangumiDataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BangumiDataRaw is an auto generated low-level Go binding around an Ethereum contract.
type BangumiDataRaw struct {
	Contract *BangumiData // Generic contract binding to access the raw methods on
}

// BangumiDataCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BangumiDataCallerRaw struct {
	Contract *BangumiDataCaller // Generic read-only contract binding to access the raw methods on
}

// BangumiDataTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BangumiDataTransactorRaw struct {
	Contract *BangumiDataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBangumiData creates a new instance of BangumiData, bound to a specific deployed contract.
func NewBangumiData(address common.Address, backend bind.ContractBackend) (*BangumiData, error) {
	contract, err := bindBangumiData(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BangumiData{BangumiDataCaller: BangumiDataCaller{contract: contract}, BangumiDataTransactor: BangumiDataTransactor{contract: contract}, BangumiDataFilterer: BangumiDataFilterer{contract: contract}}, nil
}

// NewBangumiDataCaller creates a new read-only instance of BangumiData, bound to a specific deployed contract.
func NewBangumiDataCaller(address common.Address, caller bind.ContractCaller) (*BangumiDataCaller, error) {
	contract, err := bindBangumiData(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BangumiDataCaller{contract: contract}, nil
}

// NewBangumiDataTransactor creates a new write-only instance of BangumiData, bound to a specific deployed contract.
func NewBangumiDataTransactor(address common.Address, transactor bind.ContractTransactor) (*BangumiDataTransactor, error) {
	contract, err := bindBangumiData(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BangumiDataTransactor{contract: contract}, nil
}

// NewBangumiDataFilterer creates a new log filterer instance of BangumiData, bound to a specific deployed contract.
func NewBangumiDataFilterer(address common.Address, filterer bind.ContractFilterer) (*BangumiDataFilterer, error) {
	contract, err := bindBangumiData(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BangumiDataFilterer{contract: contract}, nil
}

// bindBangumiData binds a generic wrapper to an already deployed contract.
func bindBangumiData(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BangumiDataABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BangumiData *BangumiDataRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BangumiData.Contract.BangumiDataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BangumiData *BangumiDataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BangumiData.Contract.BangumiDataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BangumiData *BangumiDataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BangumiData.Contract.BangumiDataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BangumiData *BangumiDataCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BangumiData.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BangumiData *BangumiDataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BangumiData.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BangumiData *BangumiDataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BangumiData.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BangumiData *BangumiDataCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BangumiData *BangumiDataSession) IsOwner() (bool, error) {
	return _BangumiData.Contract.IsOwner(&_BangumiData.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BangumiData *BangumiDataCallerSession) IsOwner() (bool, error) {
	return _BangumiData.Contract.IsOwner(&_BangumiData.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BangumiData *BangumiDataCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BangumiData *BangumiDataSession) Owner() (common.Address, error) {
	return _BangumiData.Contract.Owner(&_BangumiData.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BangumiData *BangumiDataCallerSession) Owner() (common.Address, error) {
	return _BangumiData.Contract.Owner(&_BangumiData.CallOpts)
}

// QueryEpisode is a free data retrieval call binding the contract method 0xfd4e9025.
//
// Solidity: function queryEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryEpisode(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryEpisode", sn)
	return *ret0, err
}

// QueryEpisode is a free data retrieval call binding the contract method 0xfd4e9025.
//
// Solidity: function queryEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryEpisode(sn string) (string, error) {
	return _BangumiData.Contract.QueryEpisode(&_BangumiData.CallOpts, sn)
}

// QueryEpisode is a free data retrieval call binding the contract method 0xfd4e9025.
//
// Solidity: function queryEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryEpisode(sn string) (string, error) {
	return _BangumiData.Contract.QueryEpisode(&_BangumiData.CallOpts, sn)
}

// QueryHash is a free data retrieval call binding the contract method 0x14d40706.
//
// Solidity: function queryHash(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryHash(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryHash", sn)
	return *ret0, err
}

// QueryHash is a free data retrieval call binding the contract method 0x14d40706.
//
// Solidity: function queryHash(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryHash(sn string) (string, error) {
	return _BangumiData.Contract.QueryHash(&_BangumiData.CallOpts, sn)
}

// QueryHash is a free data retrieval call binding the contract method 0x14d40706.
//
// Solidity: function queryHash(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryHash(sn string) (string, error) {
	return _BangumiData.Contract.QueryHash(&_BangumiData.CallOpts, sn)
}

// QueryName is a free data retrieval call binding the contract method 0xd5de3723.
//
// Solidity: function queryName(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryName(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryName", sn)
	return *ret0, err
}

// QueryName is a free data retrieval call binding the contract method 0xd5de3723.
//
// Solidity: function queryName(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryName(sn string) (string, error) {
	return _BangumiData.Contract.QueryName(&_BangumiData.CallOpts, sn)
}

// QueryName is a free data retrieval call binding the contract method 0xd5de3723.
//
// Solidity: function queryName(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryName(sn string) (string, error) {
	return _BangumiData.Contract.QueryName(&_BangumiData.CallOpts, sn)
}

// QueryPoster is a free data retrieval call binding the contract method 0x8bff2ad7.
//
// Solidity: function queryPoster(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryPoster(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryPoster", sn)
	return *ret0, err
}

// QueryPoster is a free data retrieval call binding the contract method 0x8bff2ad7.
//
// Solidity: function queryPoster(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryPoster(sn string) (string, error) {
	return _BangumiData.Contract.QueryPoster(&_BangumiData.CallOpts, sn)
}

// QueryPoster is a free data retrieval call binding the contract method 0x8bff2ad7.
//
// Solidity: function queryPoster(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryPoster(sn string) (string, error) {
	return _BangumiData.Contract.QueryPoster(&_BangumiData.CallOpts, sn)
}

// QueryRole is a free data retrieval call binding the contract method 0xa496a427.
//
// Solidity: function queryRole(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryRole(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryRole", sn)
	return *ret0, err
}

// QueryRole is a free data retrieval call binding the contract method 0xa496a427.
//
// Solidity: function queryRole(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryRole(sn string) (string, error) {
	return _BangumiData.Contract.QueryRole(&_BangumiData.CallOpts, sn)
}

// QueryRole is a free data retrieval call binding the contract method 0xa496a427.
//
// Solidity: function queryRole(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryRole(sn string) (string, error) {
	return _BangumiData.Contract.QueryRole(&_BangumiData.CallOpts, sn)
}

// QuerySarmid is a free data retrieval call binding the contract method 0x4b9722f5.
//
// Solidity: function querySarmid(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QuerySarmid(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "querySarmid", sn)
	return *ret0, err
}

// QuerySarmid is a free data retrieval call binding the contract method 0x4b9722f5.
//
// Solidity: function querySarmid(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QuerySarmid(sn string) (string, error) {
	return _BangumiData.Contract.QuerySarmid(&_BangumiData.CallOpts, sn)
}

// QuerySarmid is a free data retrieval call binding the contract method 0x4b9722f5.
//
// Solidity: function querySarmid(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QuerySarmid(sn string) (string, error) {
	return _BangumiData.Contract.QuerySarmid(&_BangumiData.CallOpts, sn)
}

// QuerySharpness is a free data retrieval call binding the contract method 0x1fccd5ff.
//
// Solidity: function querySharpness(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QuerySharpness(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "querySharpness", sn)
	return *ret0, err
}

// QuerySharpness is a free data retrieval call binding the contract method 0x1fccd5ff.
//
// Solidity: function querySharpness(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QuerySharpness(sn string) (string, error) {
	return _BangumiData.Contract.QuerySharpness(&_BangumiData.CallOpts, sn)
}

// QuerySharpness is a free data retrieval call binding the contract method 0x1fccd5ff.
//
// Solidity: function querySharpness(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QuerySharpness(sn string) (string, error) {
	return _BangumiData.Contract.QuerySharpness(&_BangumiData.CallOpts, sn)
}

// QuerySwarmAdd is a free data retrieval call binding the contract method 0x557bafc0.
//
// Solidity: function querySwarmAdd(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QuerySwarmAdd(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "querySwarmAdd", sn)
	return *ret0, err
}

// QuerySwarmAdd is a free data retrieval call binding the contract method 0x557bafc0.
//
// Solidity: function querySwarmAdd(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QuerySwarmAdd(sn string) (string, error) {
	return _BangumiData.Contract.QuerySwarmAdd(&_BangumiData.CallOpts, sn)
}

// QuerySwarmAdd is a free data retrieval call binding the contract method 0x557bafc0.
//
// Solidity: function querySwarmAdd(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QuerySwarmAdd(sn string) (string, error) {
	return _BangumiData.Contract.QuerySwarmAdd(&_BangumiData.CallOpts, sn)
}

// QueryTotalEpisode is a free data retrieval call binding the contract method 0x702e8170.
//
// Solidity: function queryTotalEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryTotalEpisode(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryTotalEpisode", sn)
	return *ret0, err
}

// QueryTotalEpisode is a free data retrieval call binding the contract method 0x702e8170.
//
// Solidity: function queryTotalEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryTotalEpisode(sn string) (string, error) {
	return _BangumiData.Contract.QueryTotalEpisode(&_BangumiData.CallOpts, sn)
}

// QueryTotalEpisode is a free data retrieval call binding the contract method 0x702e8170.
//
// Solidity: function queryTotalEpisode(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryTotalEpisode(sn string) (string, error) {
	return _BangumiData.Contract.QueryTotalEpisode(&_BangumiData.CallOpts, sn)
}

// QueryTotalSeason is a free data retrieval call binding the contract method 0x16f16be2.
//
// Solidity: function queryTotalSeason(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryTotalSeason(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryTotalSeason", sn)
	return *ret0, err
}

// QueryTotalSeason is a free data retrieval call binding the contract method 0x16f16be2.
//
// Solidity: function queryTotalSeason(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryTotalSeason(sn string) (string, error) {
	return _BangumiData.Contract.QueryTotalSeason(&_BangumiData.CallOpts, sn)
}

// QueryTotalSeason is a free data retrieval call binding the contract method 0x16f16be2.
//
// Solidity: function queryTotalSeason(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryTotalSeason(sn string) (string, error) {
	return _BangumiData.Contract.QueryTotalSeason(&_BangumiData.CallOpts, sn)
}

// QueryVideoType is a free data retrieval call binding the contract method 0x95463ae0.
//
// Solidity: function queryVideoType(string sn) constant returns(string)
func (_BangumiData *BangumiDataCaller) QueryVideoType(opts *bind.CallOpts, sn string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BangumiData.contract.Call(opts, out, "queryVideoType", sn)
	return *ret0, err
}

// QueryVideoType is a free data retrieval call binding the contract method 0x95463ae0.
//
// Solidity: function queryVideoType(string sn) constant returns(string)
func (_BangumiData *BangumiDataSession) QueryVideoType(sn string) (string, error) {
	return _BangumiData.Contract.QueryVideoType(&_BangumiData.CallOpts, sn)
}

// QueryVideoType is a free data retrieval call binding the contract method 0x95463ae0.
//
// Solidity: function queryVideoType(string sn) constant returns(string)
func (_BangumiData *BangumiDataCallerSession) QueryVideoType(sn string) (string, error) {
	return _BangumiData.Contract.QueryVideoType(&_BangumiData.CallOpts, sn)
}

// InfoInput is a paid mutator transaction binding the contract method 0xe0880590.
//
// Solidity: function _infoInput(string _bangumi, string _poster, string _role, string _hash, string _name, string _sharpness, string _episode, string _totalEpisode, string _season, string _videoType, string _swarmID, string _swarmAdd) returns()
func (_BangumiData *BangumiDataTransactor) InfoInput(opts *bind.TransactOpts, _bangumi string, _poster string, _role string, _hash string, _name string, _sharpness string, _episode string, _totalEpisode string, _season string, _videoType string, _swarmID string, _swarmAdd string) (*types.Transaction, error) {
	return _BangumiData.contract.Transact(opts, "_infoInput", _bangumi, _poster, _role, _hash, _name, _sharpness, _episode, _totalEpisode, _season, _videoType, _swarmID, _swarmAdd)
}

// InfoInput is a paid mutator transaction binding the contract method 0xe0880590.
//
// Solidity: function _infoInput(string _bangumi, string _poster, string _role, string _hash, string _name, string _sharpness, string _episode, string _totalEpisode, string _season, string _videoType, string _swarmID, string _swarmAdd) returns()
func (_BangumiData *BangumiDataSession) InfoInput(_bangumi string, _poster string, _role string, _hash string, _name string, _sharpness string, _episode string, _totalEpisode string, _season string, _videoType string, _swarmID string, _swarmAdd string) (*types.Transaction, error) {
	return _BangumiData.Contract.InfoInput(&_BangumiData.TransactOpts, _bangumi, _poster, _role, _hash, _name, _sharpness, _episode, _totalEpisode, _season, _videoType, _swarmID, _swarmAdd)
}

// InfoInput is a paid mutator transaction binding the contract method 0xe0880590.
//
// Solidity: function _infoInput(string _bangumi, string _poster, string _role, string _hash, string _name, string _sharpness, string _episode, string _totalEpisode, string _season, string _videoType, string _swarmID, string _swarmAdd) returns()
func (_BangumiData *BangumiDataTransactorSession) InfoInput(_bangumi string, _poster string, _role string, _hash string, _name string, _sharpness string, _episode string, _totalEpisode string, _season string, _videoType string, _swarmID string, _swarmAdd string) (*types.Transaction, error) {
	return _BangumiData.Contract.InfoInput(&_BangumiData.TransactOpts, _bangumi, _poster, _role, _hash, _name, _sharpness, _episode, _totalEpisode, _season, _videoType, _swarmID, _swarmAdd)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BangumiData *BangumiDataTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BangumiData.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BangumiData *BangumiDataSession) RenounceOwnership() (*types.Transaction, error) {
	return _BangumiData.Contract.RenounceOwnership(&_BangumiData.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BangumiData *BangumiDataTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BangumiData.Contract.RenounceOwnership(&_BangumiData.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BangumiData *BangumiDataTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BangumiData.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BangumiData *BangumiDataSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BangumiData.Contract.TransferOwnership(&_BangumiData.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BangumiData *BangumiDataTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BangumiData.Contract.TransferOwnership(&_BangumiData.TransactOpts, newOwner)
}

// BangumiDataOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BangumiData contract.
type BangumiDataOwnershipTransferredIterator struct {
	Event *BangumiDataOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BangumiDataOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BangumiDataOwnershipTransferred)
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
		it.Event = new(BangumiDataOwnershipTransferred)
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
func (it *BangumiDataOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BangumiDataOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BangumiDataOwnershipTransferred represents a OwnershipTransferred event raised by the BangumiData contract.
type BangumiDataOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BangumiData *BangumiDataFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BangumiDataOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BangumiData.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BangumiDataOwnershipTransferredIterator{contract: _BangumiData.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BangumiData *BangumiDataFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BangumiDataOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BangumiData.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BangumiDataOwnershipTransferred)
				if err := _BangumiData.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
