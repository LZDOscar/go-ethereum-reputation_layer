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
const MinerBookABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_HIGHLIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINER_ADMISSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_LOWLIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REPUTATION_INIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubkey\",\"type\":\"address\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reputationList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawAddrs\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubkey\",\"type\":\"address\"},{\"name\":\"_withdrawalAddressbytes48\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"usedHashedPubkey\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reputationBlackList\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"withdrawalAddressbytes48\",\"type\":\"address\"}],\"name\":\"MinerRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"}],\"name\":\"MinerDeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"reputation\",\"type\":\"uint256\"}],\"name\":\"ReputationAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hashedPubkey\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"reputation\",\"type\":\"uint256\"}],\"name\":\"ReputationSubed\",\"type\":\"event\"}]"
const MinerBookBin = `0x608060405234801561001057600080fd5b50610560806100206000396000f3fe6080604052600436106100ae576000357c0100000000000000000000000000000000000000000000000000000000900480638a3b34a5116100765780638a3b34a51461012c5780639a73a3b51461015f578063aa677354146101ae578063dd4e5d13146101dc578063fd1f473c14610223576100ae565b80631ccb00e7146100b3578063209bab08146100da5780632cb59fe6146100da57806363e8499c146100ef57806384ac33ec14610104575b600080fd5b3480156100bf57600080fd5b506100c8610256565b60408051918252519081900360200190f35b3480156100e657600080fd5b506100c861025c565b3480156100fb57600080fd5b506100c8610261565b61012a6004803603602081101561011a57600080fd5b5035600160a060020a0316610267565b005b34801561013857600080fd5b506100c86004803603602081101561014f57600080fd5b5035600160a060020a0316610390565b34801561016b57600080fd5b506101926004803603602081101561018257600080fd5b5035600160a060020a03166103a2565b60408051600160a060020a039092168252519081900360200190f35b61012a600480360360408110156101c457600080fd5b50600160a060020a03813581169160200135166103bd565b3480156101e857600080fd5b5061020f600480360360208110156101ff57600080fd5b5035600160a060020a031661050a565b604080519115158252519081900360200190f35b34801561022f57600080fd5b5061020f6004803603602081101561024657600080fd5b5035600160a060020a031661051f565b6107d081565b600081565b6103e881565b33600160a060020a038216146102c7576040805160e560020a62461bcd02815260206004820152601960248201527f496e636f7272656374206d696e65722061646d697373696f6e00000000000000604482015290519081900360640190fd5b600160a060020a038116600090815260208190526040902054819060ff16151561033b576040805160e560020a62461bcd02815260206004820152601660248201527f5075626c6963206b6579206973206e6f74207573656400000000000000000000604482015290519081900360640190fd5b600160a060020a038116600081815260208181526040808320805460ff191690556002909152808220829055517f58aa9ecbe6149ee13276eb8abe67e7c3f422c0a1a90c57dde0bf118664a9dda29190a25050565b60026020526000908152604090205481565b600160205260009081526040902054600160a060020a031681565b3415610413576040805160e560020a62461bcd02815260206004820152601960248201527f496e636f7272656374206d696e65722061646d697373696f6e00000000000000604482015290519081900360640190fd5b600160a060020a038216600090815260208190526040902054829060ff1615610486576040805160e560020a62461bcd02815260206004820152601760248201527f5075626c6963206b657920616c72656164792075736564000000000000000000604482015290519081900360640190fd5b600160a060020a03808216600081815260208181526040808320805460ff191660019081179091558252808320805495881673ffffffffffffffffffffffffffffffffffffffff19909616861790556002909152808220829055517fb90c06d0bb2c3934635f60c7418204feae8b8f21b5e030ccf695cb6cec37fe159190a3505050565b60006020819052908152604090205460ff1681565b60036020526000908152604090205460ff168156fea165627a7a72305820ebda2c40b96a99f8103b70bf1ff623da7bc27a606f18250583ca4c65020059ae0029`

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

func (_MinerBook *MinerBookCaller) GetMiners(opts *bind.CallOpts) (map[common.Address]bool, error) {
	var (
		ret0 = new(map[common.Address]bool)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "getMiners")
	return *ret0, err
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookSession) GetMiners() (map[common.Address]bool, error) {
	return _MinerBook.Contract.GetMiners(&_MinerBook.CallOpts)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes48) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) GetMiners() (map[common.Address]bool, error) {
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
