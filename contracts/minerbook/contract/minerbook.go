package contract

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

//MinerBookABI is the input ABI used to generate the binding from.
const MinerBookABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getMiners\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_HIGHLIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINER_ADMISSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_LOWLIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"regedAddrs\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_INIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubkey\",\"type\":\"address\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reputationList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawAddrs\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubkey\",\"type\":\"address\"},{\"name\":\"_withdrawalAddressbytes48\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"usedHashedPubkey\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"regedAddrsLen\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reputationBlackList\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"withdrawalAddressbytes48\",\"type\":\"address\"}],\"name\":\"MinerRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"}],\"name\":\"MinerDeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"reputation\",\"type\":\"uint256\"}],\"name\":\"ReputationAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"reputation\",\"type\":\"uint256\"}],\"name\":\"ReputationSubed\",\"type\":\"event\"}]"
const MinerBookBin = `0x6080604052600060025534801561001557600080fd5b50610ebd806100256000396000f3fe6080604052600436106100df576000357c01000000000000000000000000000000000000000000000000000000009004806384ac33ec1161009c578063aa67735411610076578063aa677354146103b1578063dd4e5d1314610415578063e65a38c91461047e578063fd1f473c146104a9576100df565b806384ac33ec146102775780638a3b34a5146102bb5780639a73a3b514610320576100df565b80631633da6e146100e45780631ccb00e714610150578063209bab081461017b5780632cb59fe6146101a6578063373d2035146101d157806363e8499c1461024c575b600080fd5b3480156100f057600080fd5b506100f9610512565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561013c578082015181840152602081019050610121565b505050509050019250505060405180910390f35b34801561015c57600080fd5b506101656105f2565b6040518082815260200191505060405180910390f35b34801561018757600080fd5b506101906105f8565b6040518082815260200191505060405180910390f35b3480156101b257600080fd5b506101bb6105fd565b6040518082815260200191505060405180910390f35b3480156101dd57600080fd5b5061020a600480360360208110156101f457600080fd5b8101908080359060200190929190505050610602565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561025857600080fd5b50610261610640565b6040518082815260200191505060405180910390f35b6102b96004803603602081101561028d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610645565b005b3480156102c757600080fd5b5061030a600480360360208110156102de57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a50565b6040518082815260200191505060405180910390f35b34801561032c57600080fd5b5061036f6004803603602081101561034357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a68565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610413600480360360408110156103c757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a9b565b005b34801561042157600080fd5b506104646004803603602081101561043857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610e4b565b604051808215151515815260200191505060405180910390f35b34801561048a57600080fd5b50610493610e6b565b6040518082815260200191505060405180910390f35b3480156104b557600080fd5b506104f8600480360360208110156104cc57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610e71565b604051808215151515815260200191505060405180910390f35b6060806002546040519080825280602002602001820160405280156105465781602001602082028038833980820191505090505b50905060008090505b6002548110156105ea5760018181548110151561056857fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682828151811015156105a157fe5b9060200190602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050808060010191505061054f565b508091505090565b6107d081565b600081565b600081565b60018181548110151561061157fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600081565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156106e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f496e636f7272656374206d696e65722061646d697373696f6e0000000000000081525060200191505060405180910390fd5b60008190506000808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615156107ad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f5075626c6963206b6579206973206e6f7420757365640000000000000000000081525060200191505060405180910390fd5b6000808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff021916905560008090505b6002548110156109c5578173ffffffffffffffffffffffffffffffffffffffff1660018281548110151561083157fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156109b857600160025414156108cc5760018181548110151561089257fe5b9060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560006002819055506109b3565b600180600254038154811015156108df57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660018281548110151561091957fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506001806002540381548110151561097457fe5b9060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560016002600082825403925050819055505b6109c5565b8080600101915050610801565b50600460008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600090558073ffffffffffffffffffffffffffffffffffffffff167f58aa9ecbe6149ee13276eb8abe67e7c3f422c0a1a90c57dde0bf118664a9dda260405160405180910390a25050565b60046020528060005260406000206000915090505481565b60036020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600034141515610b13576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f496e636f7272656374206d696e65722061646d697373696f6e0000000000000081525060200191505060405180910390fd5b60008290506000808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151515610bd9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f5075626c6963206b657920616c7265616479207573656400000000000000000081525060200191505060405180910390fd5b60016000808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001805490506002541415610cbc5760018190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050600260008154809291906001019190505550610d29565b806001600254815481101515610cce57fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506002600081548092919060010191905055505b81600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fb90c06d0bb2c3934635f60c7418204feae8b8f21b5e030ccf695cb6cec37fe1560405160405180910390a3505050565b60006020528060005260406000206000915054906101000a900460ff1681565b60025481565b60056020528060005260406000206000915054906101000a900460ff168156fea165627a7a723058205f000f0ec4f5b0e8f1f138e36757de064ba13a3a93bafd18a8fce72a0a83bf5c0029`

// DeployMinerBook deploys a new Ethereum contract, binding an instance of MinerBook to it.
func DeployMinerBook(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MinerBook, error) {
	parsed, err := abi.JSON(strings.NewReader(MinerBookABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MinerBookBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MinerBook{MinerBookCaller: MinerBookCaller{contract: contract}, MinerBookTransactor: MinerBookTransactor{contract: contract}, MinerBookFilterer: MinerBookFilterer{contract: contract}}, nil
}

// MinerBook is an auto generated Go binding around an Ethereum contract.
type MinerBook struct {
	MinerBookCaller     // Read-only binding to the contract
	MinerBookTransactor // Write-only binding to the contract
	MinerBookFilterer   // Log filterer for contract events
}

// MinerBookCaller is an auto generated read-only Go binding around an Ethereum contract.
type MinerBookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinerBookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MinerBookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinerBookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MinerBookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinerBookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MinerBookSession struct {
	Contract     *MinerBook        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MinerBookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MinerBookCallerSession struct {
	Contract *MinerBookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MinerBookransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MinerBookTransactorSession struct {
	Contract     *MinerBookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MinerBookRaw is an auto generated low-level Go binding around an Ethereum contract.
type MinerBookRaw struct {
	Contract *MinerBook // Generic contract binding to access the raw methods on
}

// MinerBookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MinerBookCallerRaw struct {
	Contract *MinerBookCaller // Generic read-only contract binding to access the raw methods on
}

// MinerBookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MinerBookTransactorRaw struct {
	Contract *MinerBookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMinerBook creates a new instance of MinerBook, bound to a specific deployed contract.
func NewMinerBook(address common.Address, backend bind.ContractBackend) (*MinerBook, error) {
	contract, err := bindMinerBook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MinerBook{MinerBookCaller: MinerBookCaller{contract: contract}, MinerBookTransactor: MinerBookTransactor{contract: contract}, MinerBookFilterer: MinerBookFilterer{contract: contract}}, nil
}

// NewMinerBookCaller creates a new read-only instance of MinerBook, bound to a specific deployed contract.
func NewMinerBookCaller(address common.Address, caller bind.ContractCaller) (*MinerBookCaller, error) {
	contract, err := bindMinerBook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MinerBookCaller{contract: contract}, nil
}

// NewMinerBookTransactor creates a new write-only instance of MinerBook, bound to a specific deployed contract.
func NewMinerBookTransactor(address common.Address, transactor bind.ContractTransactor) (*MinerBookTransactor, error) {
	contract, err := bindMinerBook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MinerBookTransactor{contract: contract}, nil
}

// NewMinerBookFilterer creates a new log filterer instance of MinerBook, bound to a specific deployed contract.
func NewMinerBookFilterer(address common.Address, filterer bind.ContractFilterer) (*MinerBookFilterer, error) {
	contract, err := bindMinerBook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MinerBookFilterer{contract: contract}, nil
}

// bindMinerBook binds a generic wrapper to an already deployed contract.
func bindMinerBook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MinerBookABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinerBook *MinerBookRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MinerBook.Contract.MinerBookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinerBook *MinerBookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinerBook.Contract.MinerBookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinerBook *MinerBookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinerBook.Contract.MinerBookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinerBook *MinerBookCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MinerBook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinerBook *MinerBookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinerBook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinerBook *MinerBookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinerBook.Contract.contract.Transact(opts, method, params...)
}

// MinerAdmission is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function MinerAdmission() constant returns(uint256)
func (_MinerBook *MinerBookCaller) MinerAdmission(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "MINER_ADMISSION")
	return *ret0, err
}

// MinerAdmission is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function MinerAdmission() constant returns(uint256)
func (_MinerBook *MinerBookSession) MinerAdmission() (*big.Int, error) {
	return _MinerBook.Contract.MinerAdmission(&_MinerBook.CallOpts)
}

// MinerAdmission is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function MinerAdmission() constant returns(uint256)
func (_MinerBook *MinerBookCallerSession) MinerAdmission() (*big.Int, error) {
	return _MinerBook.Contract.MinerAdmission(&_MinerBook.CallOpts)
}

// UsedHashedPubkey is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function UsedHashedPubkey(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCaller) UsedHashedPubkey(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "usedHashedPubkey", arg0)
	return *ret0, err
}

// UsedHashedPubkey is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function UsedHashedPubkey(bytes48) constant returns(bool)
func (_MinerBook *MinerBookSession) UsedHashedPubkey(arg0 common.Address) (bool, error) {
	return _MinerBook.Contract.UsedHashedPubkey(&_MinerBook.CallOpts, arg0)
}

// UsedHashedPubkey is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function UsedHashedPubkey(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) UsedHashedPubkey(arg0 common.Address) (bool, error) {
	return _MinerBook.Contract.UsedHashedPubkey(&_MinerBook.CallOpts, arg0)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCaller) ReputationList(opts *bind.CallOpts, arg0 common.Address) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "reputationList", arg0)
	return *ret0, err
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookSession) ReputationList(arg0 common.Address) (uint64, error) {
	return _MinerBook.Contract.ReputationList(&_MinerBook.CallOpts, arg0)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) ReputationList(arg0 common.Address) (uint64, error) {
	return _MinerBook.Contract.ReputationList(&_MinerBook.CallOpts, arg0)
}

func (_MinerBook *MinerBookCaller) GetMiners(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "getMiners")
	return *ret0, err
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookSession) GetMiners() ([]common.Address, error) {
	return _MinerBook.Contract.GetMiners(&_MinerBook.CallOpts)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) GetMiners() ([]common.Address, error) {
	return _MinerBook.Contract.GetMiners(&_MinerBook.CallOpts)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactor) Register(opts *bind.TransactOpts, _pubkey common.Address, _withdrawalAddressbytes48 common.Address) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "register", _pubkey, _withdrawalAddressbytes48)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookSession) Register(_pubkey common.Address, _withdrawalAddressbytes48 common.Address) (*types.Transaction, error) {
	return _MinerBook.Contract.Register(&_MinerBook.TransactOpts, _pubkey, _withdrawalAddressbytes48)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactorSession) Register(_pubkey common.Address, _withdrawalAddressbytes48 common.Address) (*types.Transaction, error) {
	return _MinerBook.Contract.Register(&_MinerBook.TransactOpts, _pubkey, _withdrawalAddressbytes48)
}

// DeRegister is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: DeRegister (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactor) DeRegister(opts *bind.TransactOpts, _pubkey common.Address) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "deregister", _pubkey)
}

// DeRegister is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: DeRegister (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookSession) DeRegister(_pubkey common.Address) (*types.Transaction, error) {
	return _MinerBook.Contract.DeRegister(&_MinerBook.TransactOpts, _pubkey)
}

// DeRegister is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: DeRegister (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactorSession) DeRegister(_pubkey common.Address) (*types.Transaction, error) {
	return _MinerBook.Contract.DeRegister(&_MinerBook.TransactOpts, _pubkey)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookSession) AddReputation(_pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.Contract.AddReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactorSession) AddReputation(_pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.Contract.AddReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactor) AddReputation(opts *bind.TransactOpts, _pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "addReputation", _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookSession) SubReputation(_pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.Contract.SubReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactorSession) SubReputation(_pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.Contract.SubReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes48 address, _randaoCommitment bytes48) returns()
func (_MinerBook *MinerBookTransactor) SubReputation(opts *bind.TransactOpts, _pubkey []byte, _value uint64) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "subReputation", _pubkey, _value)
}

// MinerBookMinerRegisteredIterator is returned from FilterMinerRegistered and is used to iterate over the raw logs and unpacked data for MinerRegistered events raised by the MinerBook contract.
type MinerBookMinerRegisteredIterator struct {
	Event *MinerBookMinerRegistered // Event containing the contract specifics and raw log

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
func (it *MinerBookMinerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MinerBookMinerRegistered)
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
		it.Event = new(MinerBookMinerRegistered)
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

// Error retruned any retrieval or parsing error occurred during filtering.
func (it *MinerBookMinerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MinerBookMinerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MinerBookOverdraft represents a Overdraft event raised by the MinerBook contract.
type MinerBookMinerRegistered struct {
	HashedPubkey             common.Address
	WithdrawalAddressbytes48 common.Address
	//RandaoCommitment         [48]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMinerRegistered is a free log retrieval operation binding the contract event 0x2250e2993c15843b48621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: event MinerRegistered(hashedPubkey indexed bytes48, withdrawalAddressbytes48 indexed address, randaoCommitment indexed bytes48)
func (_MinerBook *MinerBookFilterer) FilterMinerRegistered(opts *bind.FilterOpts, hashedPubkey []common.Address, withdrawalAddressbytes48 []common.Address) (*MinerBookMinerRegisteredIterator, error) {

	var hashedPubkeyRule []interface{}
	for _, hashedPubkeyItem := range hashedPubkey {
		hashedPubkeyRule = append(hashedPubkeyRule, hashedPubkeyItem)
	}

	var withdrawalAddressbytes48Rule []interface{}
	for _, withdrawalAddressbytes48Item := range withdrawalAddressbytes48 {
		withdrawalAddressbytes48Rule = append(withdrawalAddressbytes48Rule, withdrawalAddressbytes48Item)
	}

	//var randaoCommitmentRule []interface{}
	//for _, randaoCommitmentItem := range randaoCommitment {
	//	randaoCommitmentRule = append(randaoCommitmentRule, randaoCommitmentItem)
	//}

	logs, sub, err := _MinerBook.contract.FilterLogs(opts, "MinerRegistered", hashedPubkeyRule, withdrawalAddressbytes48Rule)
	if err != nil {
		return nil, err
	}

	return &MinerBookMinerRegisteredIterator{contract: _MinerBook.contract, event: "MinerRegistered", logs: logs, sub: sub}, nil
}

// WatchMinerRegistered is a free log subscription operation binding the contract event 0x2250e2993c15843b48621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: event MinerRegistered(hashedPubkey indexed bytes48, withdrawalAddressbytes48 indexed address, randaoCommitment indexed bytes48)
func (_MinerBook *MinerBookFilterer) WatchMinerRegistered(opts *bind.WatchOpts, sink chan<- *MinerBookMinerRegistered, hashedPubkey []common.Address, withdrawalAddressbytes48 []common.Address) (event.Subscription, error) {

	var hashedPubkeyRule []interface{}
	for _, hashedPubkeyItem := range hashedPubkey {
		hashedPubkeyRule = append(hashedPubkeyRule, hashedPubkeyItem)
	}

	var withdrawalAddressbytes48Rule []interface{}
	for _, withdrawalAddressbytes48Item := range withdrawalAddressbytes48 {
		withdrawalAddressbytes48Rule = append(withdrawalAddressbytes48Rule, withdrawalAddressbytes48Item)
	}

	//var randaoCommitmentRule []interface{}
	//for _, randaoCommitmentItem := range randaoCommitment {
	//	randaoCommitmentRule = append(randaoCommitmentRule, randaoCommitmentItem)
	//}

	logs, sub, err := _MinerBook.contract.WatchLogs(opts, "MinerRegistered", hashedPubkeyRule, withdrawalAddressbytes48Rule)
	if err != nil {
		return nil, err
	}

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MinerBookMinerRegistered)
				if err := _MinerBook.contract.UnpackLog(event, "MinerRegistered", log); err != nil {
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
