package contract

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
	//"math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

var (
	amount33Eth, _ = new(big.Int).SetString("33000000000000000000", 10)
	amount32Eth, _ = new(big.Int).SetString("32000000000000000000", 10)
	amount31Eth, _ = new(big.Int).SetString("31000000000000000000", 10)
	amount00Eth, _ = new(big.Int).SetString("0000000000000000000", 10)
	amount01Eth, _ = new(big.Int).SetString("1000000000000000000", 10)
)

type testAccount struct {
	addr              common.Address
	withdrawalAddress common.Address
	randaoCommitment  [32]byte
	pubKey            []byte
	contract          *MinerBook
	backend           *backends.SimulatedBackend
	txOpts            *bind.TransactOpts

	//callOpts *bind.CallOpts
}

func setup() (*testAccount, error) {
	genesis := make(core.GenesisAlloc)
	privKey, _ := crypto.GenerateKey()
	pubKeyECDSA, ok := privKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(pubKeyECDSA)[4:]
	var pubKey = make([]byte, 48)
	copy(pubKey[:], []byte(publicKeyBytes))

	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	txOpts := bind.NewKeyedTransactor(privKey)
	//callOpts :=bind.new
	genesis[addr] = core.GenesisAccount{Balance: big.NewInt(1000000000), Reputation: 1000}
	backend := backends.NewSimulatedBackend(genesis, 10000000)

	_, _, contract, err := DeployMinerBook(txOpts, backend)
	if err != nil {
		return nil, err
	}
	backend.Commit()

	return &testAccount{
		addr,
		common.Address{},
		[32]byte{},
		pubKey,
		contract,
		backend,
		txOpts,
	}, nil
}

func TestSetupAndContractRegistration(t *testing.T) {
	_, err := setup()
	if err != nil {
		log.Fatalf("Can not deploy validator registration contract: %v", err)
	}
}

// negative test case, test registering with the same public key twice.
func TestRegisterTwice(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	//randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}

	testAccount.txOpts.Value = amount00Eth
	//testAccount.txOpts.GasLimit = math.MaxUint64
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.addr, *withdrawAddr)
	testAccount.backend.Commit()
	if err != nil {
		t.Errorf("Validator registration failed: %v", err)
	}
	//
	testAccount.txOpts.Value = amount00Eth
	//testAccount.txOpts.GasLimit = math.MaxUint64
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.addr, *withdrawAddr)
	testAccount.backend.Commit()
	if err == nil {
		t.Errorf("Registration should have failed with same public key twice")
	}
	if err != nil {
		t.Error(err)
	}
}

//
// normal test case, test depositing 32 ETH and verify validatorRegistered event is correctly emitted.
func TestRegister(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	//randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}
	testAccount.txOpts.Value = amount00Eth

	var hashedPub [20]byte
	copy(hashedPub[:], crypto.Keccak256(testAccount.pubKey))

	println("register account address:" + testAccount.addr.String())
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.addr, *withdrawAddr)
	testAccount.backend.Commit()
	if err == nil {
		println("register no error")
	}
	if err != nil {
		t.Errorf("Validator registration failed: %v", err)
	}
	log, err := testAccount.contract.FilterMinerRegistered(&bind.FilterOpts{}, []common.Address{}, []common.Address{})

	defer func() {
		err = log.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if err != nil {
		t.Fatal(err)
	}
	if log.Error() != nil {
		t.Fatal(log.Error())
	}
	log.Next()
	//if log.Event.RandaoCommitment != *randaoCommitment {
	//	t.Errorf("validatorRegistered event randao commitment miss matched. Want: %v, Got: %v", *randaoCommitment, log.Event.RandaoCommitment)
	//}
	if log.Event.HashedPubkey != testAccount.addr {
		t.Errorf("validatorRegistered event public key miss matched. Want: %v, Got: %v", common.BytesToHash(testAccount.pubKey), log.Event.HashedPubkey)
	}
	if log.Event.WithdrawalAddressbytes48 != *withdrawAddr {
		t.Errorf("validatorRegistered event withdrawal address miss matched. Want: %v, Got: %v", *withdrawAddr, log.Event.WithdrawalAddressbytes48)
	}

	//var miners map[common.Address] bool

	miner, err := testAccount.contract.UsedHashedPubkey(nil, testAccount.addr)
	if err != nil {
		t.Errorf("getminer :%v", err)
	}
	print("getminer: ")
	println(miner)

	miners, err := testAccount.contract.GetMiners(nil)
	if err != nil {
		t.Errorf("getminers :%v", err)
	}
	println("getminers")
	for _, v := range miners {
		println(v.String())
	}

	_, err = testAccount.contract.DeRegister(testAccount.txOpts, testAccount.addr)
	testAccount.backend.Commit()
	if err == nil {
		println("deregister no error")
	}
	if err != nil {
		t.Errorf("Validator registration failed: %v", err)
	}

	miner, err = testAccount.contract.UsedHashedPubkey(nil, testAccount.addr)
	if err != nil {
		t.Errorf("getminer :%v", err)
	}
	print("getminer: ")
	println(miner)

	miners, err = testAccount.contract.GetMiners(nil)
	if err != nil {
		t.Errorf("getminers :%v", err)
	}
	println("getminers")
	for _, v := range miners {
		println(v.String())
	}
}
