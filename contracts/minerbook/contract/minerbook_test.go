package contract

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
)

var (
	amount33Eth, _ = new(big.Int).SetString("33000000000000000000", 10)
	amount32Eth, _ = new(big.Int).SetString("32000000000000000000", 10)
	amount31Eth, _ = new(big.Int).SetString("31000000000000000000", 10)
	amount02Eth, _ = new(big.Int).SetString("2000000000000000000", 10)
)

type testAccount struct {
	addr              common.Address
	withdrawalAddress common.Address
	randaoCommitment  [32]byte
	pubKey            []byte
	contract          *MinerBook
	backend           *backends.SimulatedBackend
	txOpts            *bind.TransactOpts
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
	startingBalance, _ := new(big.Int).SetString("100000000000000000000", 10)
	genesis[addr] = core.GenesisAccount{Balance: startingBalance}
	backend := backends.NewSimulatedBackend(genesis, 2100000)

	_, _, contract, err := DeployMinerBook(txOpts, backend)
	if err != nil {
		return nil, err
	}

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

// negative test case, public key that is not 48 bytes.
func TestRegisterWithLessThan48BytesPubkey(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	var pubKey = make([]byte, 32)
	copy(pubKey, testAccount.pubKey[:])
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}

	testAccount.txOpts.Value = amount32Eth
	_, err = testAccount.contract.Register(testAccount.txOpts, pubKey, *withdrawAddr, *randaoCommitment)
	if err == nil {
		t.Error("Validator registration should have failed with a 32 bytes pubkey")
	}
	if err != nil {
		println("Test Success!")
	}
}

// negative test case, deposit with less than 32 ETH.
func TestRegisterWithLessThan32Eth(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}

	testAccount.txOpts.Value = amount31Eth
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.pubKey, *withdrawAddr, *randaoCommitment)
	if err == nil {
		t.Error("Validator registration should have failed with insufficient deposit")
	}
	if err != nil {
		t.Error(err)
	}
}

// negative test case, deposit more than 32 ETH.
func TestRegisterWithMoreThan32Eth(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}

	testAccount.txOpts.Value = amount33Eth
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.pubKey, *withdrawAddr, *randaoCommitment)
	if err == nil {
		t.Error("Validator registration should have failed with more than deposit amount")
	}
	if err != nil {
		t.Error(err)
	}
}

// negative test case, test registering with the same public key twice.
func TestRegisterTwice(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}

	testAccount.txOpts.Value = amount32Eth
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.pubKey, *withdrawAddr, *randaoCommitment)
	testAccount.backend.Commit()
	if err != nil {
		t.Errorf("Validator registration failed: %v", err)
	}

	testAccount.txOpts.Value = amount32Eth
	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.pubKey, *withdrawAddr, *randaoCommitment)
	testAccount.backend.Commit()
	if err == nil {
		t.Errorf("Registration should have failed with same public key twice")
	}
	if err != nil {
		t.Error(err)
	}
}

// normal test case, test depositing 32 ETH and verify validatorRegistered event is correctly emitted.
func TestRegister(t *testing.T) {
	testAccount, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	withdrawAddr := &common.Address{'A', 'D', 'D', 'R', 'E', 'S', 'S'}
	randaoCommitment := &[32]byte{'S', 'H', 'H', 'H', 'H', 'I', 'T', 'S', 'A', 'S', 'E', 'C', 'R', 'E', 'T'}
	testAccount.txOpts.Value = amount32Eth

	var hashedPub [32]byte
	copy(hashedPub[:], crypto.Keccak256(testAccount.pubKey))

	_, err = testAccount.contract.Register(testAccount.txOpts, testAccount.pubKey, *withdrawAddr, *randaoCommitment)
	testAccount.backend.Commit()
	if err == nil {
		t.Errorf("no error")
	}
	if err != nil {
		t.Errorf("Validator registration failed: %v", err)
	}
	log, err := testAccount.contract.FilterMinerRegistered(&bind.FilterOpts{}, [][32]byte{}, []common.Address{}, [][32]byte{})

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
	if log.Event.RandaoCommitment != *randaoCommitment {
		t.Errorf("validatorRegistered event randao commitment miss matched. Want: %v, Got: %v", *randaoCommitment, log.Event.RandaoCommitment)
	}
	if log.Event.HashedPubkey != hashedPub {
		t.Errorf("validatorRegistered event public key miss matched. Want: %v, Got: %v", common.BytesToHash(testAccount.pubKey), log.Event.HashedPubkey)
	}
	if log.Event.WithdrawalAddressbytes32 != *withdrawAddr {
		t.Errorf("validatorRegistered event withdrawal address miss matched. Want: %v, Got: %v", *withdrawAddr, log.Event.WithdrawalAddressbytes32)
	}
}
