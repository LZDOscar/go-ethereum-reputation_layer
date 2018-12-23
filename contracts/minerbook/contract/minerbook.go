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
const MinerBookABI = "[{\"anonymous\": false,\"inputs\": [{\"indexed\": true,\"name\": \"hashedPubkey\",\"type\": \"bytes32\"},{\"indexed\": true, \"name\": \"withdrawalAddressbytes32\", \"type\": \"address\"},{\"indexed\": true, \"name\": \"randaoCommitment\", \"type\": \"bytes32\"}], \"name\": \"MinerRegistered\", \"type\": \"event\"},{\"constant\": false, \"inputs\": [{\"name\": \"_pubkey\", \"type\": \"bytes\"},{\"name\": \"_withdrawalAddressbytes32\", \"type\": \"address\"},{\"name\": \"_randaoCommitment\", \"type\": \"bytes32\"}], \"name\": \"register\", \"outputs\": [], \"payable\": true, \"stateMutability\": \"payable\",\"type\": \"function\"},{\"constant\": true, \"inputs\": [], \"name\": \"MINER_ADMISSION\", \"outputs\": [{\"name\": \"\", \"type\": \"uint256\"}], \"payable\": false, \"stateMutability\": \"view\", \"type\": \"function\"},{\"constant\": true, \"inputs\": [{\"name\": \"\", \"type\": \"bytes32\"}], \"name\": \"usedHashedPubkey\", \"outputs\": [{\"name\": \"\", \"type\": \"bool\"}], \"payable\": false, \"stateMutability\": \"view\", \"type\": \"function\"}]"

//MinerBookBin is the compiled bytecode used for deploying new contracts.
const MinerBookBin = `608060405234801561001057600080fd5b50610401806100206000396000f3fe608060405260043610610050577c01000000000000000000000000000000000000000000000000000000006000350463209bab0881146100555780638618d7781461007c578063930a6539146100ba575b600080fd5b34801561006157600080fd5b5061006a61017d565b60408051918252519081900360200190f35b34801561008857600080fd5b506100a66004803603602081101561009f57600080fd5b503561018a565b604080519115158252519081900360200190f35b61017b600480360360608110156100d057600080fd5b8101906020810181356401000000008111156100eb57600080fd5b8201836020820111156100fd57600080fd5b8035906020019184600183028401116401000000008311171561011f57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505073ffffffffffffffffffffffffffffffffffffffff833516935050506020013561019f565b005b6801bc16d674ec80000081565b60006020819052908152604090205460ff1681565b346801bc16d674ec8000001461021657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e636f7272656374206d696e65722061646d697373696f6e00000000000000604482015290519081900360640190fd5b825160301461028657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f5075626c6963206b6579206973206e6f74203438206279746573000000000000604482015290519081900360640190fd5b6000836040516020018082805190602001908083835b602083106102bb5780518252601f19909201916020918201910161029c565b51815160209384036101000a60001901801990921691161790526040805192909401828103601f190183528452815191810191909120600081815291829052929020549194505060ff16159150610375905057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5075626c6963206b657920616c72656164792075736564000000000000000000604482015290519081900360640190fd5b600081815260208190526040808220805460ff1916600117905551839173ffffffffffffffffffffffffffffffffffffffff86169184917fe64a07a1b5d850fa4668b6f2243f96f9002d7f4bb5ad61a0b5cbec94cbd7d5ff91a45050505056fea165627a7a723058203d0f3db4e9676752c9dc1be83c174f7e40aff140194a455f99dddadd9c8293e00029`

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
// Solidity: function UsedHashedPubkey(bytes32) constant returns(bool)
func (_MinerBook *MinerBookCaller) UsedHashedPubkey(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "usedHashedPubkey", arg0)
	return *ret0, err
}

// UsedHashedPubkey is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function UsedHashedPubkey(bytes32) constant returns(bool)
func (_MinerBook *MinerBookSession) UsedHashedPubkey(arg0 [32]byte) (bool, error) {
	return _MinerBook.Contract.UsedHashedPubkey(&_MinerBook.CallOpts, arg0)
}

// UsedHashedPubkey is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function UsedHashedPubkey(bytes32) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) UsedHashedPubkey(arg0 [32]byte) (bool, error) {
	return _MinerBook.Contract.UsedHashedPubkey(&_MinerBook.CallOpts, arg0)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes32) constant returns(bool)
func (_MinerBook *MinerBookCaller) ReputationList(opts *bind.CallOpts, arg0 [32]byte) (int64, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MinerBook.contract.Call(opts, out, "reputationList", arg0)
	return *ret0, err
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes32) constant returns(bool)
func (_MinerBook *MinerBookSession) ReputationList(arg0 [32]byte) (int64, error) {
	return _MinerBook.Contract.ReputationList(&_MinerBook.CallOpts, arg0)
}

// reputationList is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function reputationList(bytes32) constant returns(bool)
func (_MinerBook *MinerBookCallerSession) ReputationList(arg0 [32]byte) (int64, error) {
	return _MinerBook.Contract.ReputationList(&_MinerBook.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactor) Register(opts *bind.TransactOpts, _pubkey []byte, _withdrawalAddressbytes32 common.Address, _randaoCommitment [32]byte) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "register", _pubkey, _withdrawalAddressbytes32, _randaoCommitment)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookSession) Register(_pubkey []byte, _withdrawalAddressbytes32 common.Address, _randaoCommitment [32]byte) (*types.Transaction, error) {
	return _MinerBook.Contract.Register(&_MinerBook.TransactOpts, _pubkey, _withdrawalAddressbytes32, _randaoCommitment)
}

// Register is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: Register (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactorSession) Register(_pubkey []byte, _withdrawalAddressbytes32 common.Address, _randaoCommitment [32]byte) (*types.Transaction, error) {
	return _MinerBook.Contract.Register(&_MinerBook.TransactOpts, _pubkey, _withdrawalAddressbytes32, _randaoCommitment)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookSession) AddReputation(_pubkey []byte, _value int) (*types.Transaction, error) {
	return _MinerBook.Contract.AddReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactorSession) AddReputation(_pubkey []byte, _value int) (*types.Transaction, error) {
	return _MinerBook.Contract.AddReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// addReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: addReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactor) AddReputation(opts *bind.TransactOpts, _pubkey []byte, _value int) (*types.Transaction, error) {
	return _MinerBook.contract.Transact(opts, "addReputation", _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookSession) SubReputation(_pubkey []byte, _value int) (*types.Transaction, error) {
	return _MinerBook.Contract.SubReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactorSession) SubReputation(_pubkey []byte, _value int) (*types.Transaction, error) {
	return _MinerBook.Contract.SubReputation(&_MinerBook.TransactOpts, _pubkey, _value)
}

// SubReputation is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: SubReputation (_pubkey bytes, _withdrawalAddressbytes32 address, _randaoCommitment bytes32) returns()
func (_MinerBook *MinerBookTransactor) SubReputation(opts *bind.TransactOpts, _pubkey []byte, _value int) (*types.Transaction, error) {
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
	HashedPubkey             [32]byte
	WithdrawalAddressbytes32 common.Address
	RandaoCommitment         [32]byte
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterMinerRegistered is a free log retrieval operation binding the contract event 0x2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: event MinerRegistered(hashedPubkey indexed bytes32, withdrawalAddressbytes32 indexed address, randaoCommitment indexed bytes32)
func (_MinerBook *MinerBookFilterer) FilterMinerRegistered(opts *bind.FilterOpts, hashedPubkey [][32]byte, withdrawalAddressbytes32 []common.Address, randaoCommitment [][32]byte) (*MinerBookMinerRegisteredIterator, error) {

	var hashedPubkeyRule []interface{}
	for _, hashedPubkeyItem := range hashedPubkey {
		hashedPubkeyRule = append(hashedPubkeyRule, hashedPubkeyItem)
	}

	var withdrawalAddressbytes32Rule []interface{}
	for _, withdrawalAddressbytes32Item := range withdrawalAddressbytes32 {
		withdrawalAddressbytes32Rule = append(withdrawalAddressbytes32Rule, withdrawalAddressbytes32Item)
	}

	var randaoCommitmentRule []interface{}
	for _, randaoCommitmentItem := range randaoCommitment {
		randaoCommitmentRule = append(randaoCommitmentRule, randaoCommitmentItem)
	}

	logs, sub, err := _MinerBook.contract.FilterLogs(opts, "MinerRegistered", hashedPubkeyRule, withdrawalAddressbytes32Rule, randaoCommitmentRule)
	if err != nil {
		return nil, err
	}

	return &MinerBookMinerRegisteredIterator{contract: _MinerBook.contract, event: "MinerRegistered", logs: logs, sub: sub}, nil
}

// WatchMinerRegistered is a free log subscription operation binding the contract event 0x2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: event MinerRegistered(hashedPubkey indexed bytes32, withdrawalAddressbytes32 indexed address, randaoCommitment indexed bytes32)
func (_MinerBook *MinerBookFilterer) WatchMinerRegistered(opts *bind.WatchOpts, sink chan<- *MinerBookMinerRegistered, hashedPubkey [][32]byte, withdrawalAddressbytes32 []common.Address, randaoCommitment [][32]byte) (event.Subscription, error) {

	var hashedPubkeyRule []interface{}
	for _, hashedPubkeyItem := range hashedPubkey {
		hashedPubkeyRule = append(hashedPubkeyRule, hashedPubkeyItem)
	}

	var withdrawalAddressbytes32Rule []interface{}
	for _, withdrawalAddressbytes32Item := range withdrawalAddressbytes32 {
		withdrawalAddressbytes32Rule = append(withdrawalAddressbytes32Rule, withdrawalAddressbytes32Item)
	}

	var randaoCommitmentRule []interface{}
	for _, randaoCommitmentItem := range randaoCommitment {
		randaoCommitmentRule = append(randaoCommitmentRule, randaoCommitmentItem)
	}

	logs, sub, err := _MinerBook.contract.WatchLogs(opts, "MinerRegistered", hashedPubkeyRule, withdrawalAddressbytes32Rule, randaoCommitmentRule)
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
