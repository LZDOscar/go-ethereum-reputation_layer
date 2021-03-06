// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	"bytes"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	m "math"
	"math/big"
	"runtime"
	"time"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	//"github.com/ethereum/go-ethereum/contracts/minerbook"
	//"github.com/ethereum/go-ethereum/contracts/minerbook/contract"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"

	MBC "github.com/ethereum/go-ethereum/contracts/minerbook/contract"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"strings"
	//"log"
	//"github.com/ethereum/go-ethereum/contracts/minerbook"
	"log"
)

// Ethash proof-of-work protocol constants.
var (
	FrontierBlockReward           = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
	FrontierBlockReputationReward = int64(1)          // Block reward in reputation for successfully mining a block
	ByzantiumBlockReward          = big.NewInt(3e+18) // Block reward in wei for successfully mining a block upward from Byzantium
	ConstantinopleBlockReward     = big.NewInt(2e+18) // Block reward in wei for successfully mining a block upward from Constantinople
	maxUncles                     = 2                 // Maximum number of uncles allowed in a single block
	allowedFutureBlockTime        = 15 * time.Second  // Max time from current time allowed for blocks, before they're considered future blocks

	ReputationLowThreshold              = uint64(0)
	ReputationHighThreshold             = uint64(2000)
	ReputationInit                      = uint64(1000)
	ReputationFrontierBlockCount        = 20 //测试所用值，真实值要根据miner数量来定，大概为8×minerAccount
	ReputationBlackBlockCount           = 40 //测试所用值，真实值要根据miner数量来定，大概为12×minerAccount
	ReputationRwardFormulaOptimizeParam = 100
	ReputationDecayFormulaOptimizeParam = 30
	ReputationCalcDiffBlockCount        = 10 // 测试所用值，真实值要根据miner数量来定，大概4×minerAccount
	ReputationWhiteAddress              = common.HexToAddress("0000000000000000000000000000000000000000")
	ReputationExpectedRewardAccount     = 6
	ReputationExpectedDecayAccount      = 2
	ReputationDifficultyRadio           = uint64(1) // ReputationMax <=> Difficulty 最大信誉值可以转化为1/3的difficulty
	ReputationContinuousBlockUsable     = float64(1.3)
	// calcDifficultyConstantinople is the difficulty adjustment algorithm for Constantinople.
	// It returns the difficulty that a new block should have when created at time given the
	// parent block's time and difficulty. The calculation uses the Byzantium rules, but with
	// bomb offset 5M.
	// Specification EIP-1234: https://eips.ethereum.org/EIPS/eip-1234
	calcDifficultyConstantinople = makeDifficultyCalculator(big.NewInt(5000000))

	// calcDifficultyByzantium is the difficulty adjustment algorithm. It returns
	// the difficulty that a new block should have when created at time given the
	// parent block's time and difficulty. The calculation uses the Byzantium rules.
	// Specification EIP-649: https://eips.ethereum.org/EIPS/eip-649
	calcDifficultyByzantium = makeDifficultyCalculator(big.NewInt(3000000))

	//临时测试所用的
	ReputationContractAddress = common.Address{}
	MinersListTest            = []common.Address{common.HexToAddress("0000000000000000000000000000000000000001"),
		common.HexToAddress("0000000000000000000000000000000000000002"),
		common.HexToAddress("0000000000000000000000000000000000000003"),
		common.HexToAddress("0000000000000000000000000000000000000004"),
		common.HexToAddress("0000000000000000000000000000000000000005"),
		//common.HexToAddress("0000000000000000000000000000000000000006"),
		//common.HexToAddress("0000000000000000000000000000000000000007"),
		//common.HexToAddress("0000000000000000000000000000000000000008"),
	}
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errLargeBlockTime    = errors.New("timestamp too big")
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errTooManyUncles     = errors.New("too many uncles")
	errDuplicateUncle    = errors.New("duplicate uncle")
	errUncleIsAncestor   = errors.New("uncle is ancestor")
	errDanglingUncle     = errors.New("uncle's parent is not ancestor")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (ethash *Ethash) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
func (ethash *Ethash) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Short circuit if the header is known, or it's parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return ethash.verifyHeader(chain, header, parent, false, seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (ethash *Ethash) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = ethash.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (ethash *Ethash) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if chain.GetHeader(headers[index].Hash(), headers[index].Number.Uint64()) != nil {
		return nil // known block
	}
	return ethash.verifyHeader(chain, headers[index], parent, false, seals[index])
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of the stock Ethereum ethash engine.
func (ethash *Ethash) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Verify that there are at most 2 uncles included in this block
	if len(block.Uncles()) > maxUncles {
		return errTooManyUncles
	}
	// Gather the set of past uncles and ancestors
	uncles, ancestors := mapset.NewSet(), make(map[common.Hash]*types.Header)

	number, parent := block.NumberU64()-1, block.ParentHash()
	for i := 0; i < 7; i++ {
		ancestor := chain.GetBlock(parent, number)
		if ancestor == nil {
			break
		}
		ancestors[ancestor.Hash()] = ancestor.Header()
		for _, uncle := range ancestor.Uncles() {
			uncles.Add(uncle.Hash())
		}
		parent, number = ancestor.ParentHash(), number-1
	}
	ancestors[block.Hash()] = block.Header()
	uncles.Add(block.Hash())

	// Verify each of the uncles that it's recent, but not an ancestor
	for _, uncle := range block.Uncles() {
		// Make sure every uncle is rewarded only once
		hash := uncle.Hash()
		if uncles.Contains(hash) {
			return errDuplicateUncle
		}
		uncles.Add(hash)

		// Make sure the uncle has a valid ancestry
		if ancestors[hash] != nil {
			return errUncleIsAncestor
		}
		if ancestors[uncle.ParentHash] == nil || uncle.ParentHash == block.ParentHash() {
			return errDanglingUncle
		}
		if err := ethash.verifyHeader(chain, uncle, ancestors[uncle.ParentHash], true, true); err != nil {
			return err
		}
	}
	return nil
}

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
// See YP section 4.3.4. "Block Header Validity"
func (ethash *Ethash) verifyHeader(chain consensus.ChainReader, header, parent *types.Header, uncle bool, seal bool) error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if uncle {
		if header.Time.Cmp(math.MaxBig256) > 0 {
			return errLargeBlockTime
		}
	} else {
		if header.Time.Cmp(big.NewInt(time.Now().Add(allowedFutureBlockTime).Unix())) > 0 {
			return consensus.ErrFutureBlock
		}
	}
	if header.Time.Cmp(parent.Time) <= 0 {
		return errZeroBlockTime
	}
	// Verify the block's difficulty based in it's timestamp and parent's difficulty
	// New Change: add the reputation of Author to CalcDifficulty function

	// Verify the author miner
	//TODO
	//author, err := ethash.Author(header)
	//if err != nil {
	//	return fmt.Errorf("invalid Author")
	//}
	//reg, err := ethash.CheckRegister(author)
	//if err != nil{
	//	return fmt.Errorf("check register failed")
	//}
	//if !reg {
	//	return fmt.Errorf("block producer hasn't registed, is not a miner!")
	//}

	expected := ethash.CalcDifficulty(chain, header.Time.Uint64(), parent)

	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.GasLimit / params.GasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}

	// Verify the engine specific seal securing the block
	if seal {
		if err := ethash.VerifySeal(chain, header); err != nil {
			return err
		}
	}
	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyDAOHeaderExtraData(chain.Config(), header); err != nil {
		return err
	}
	if err := misc.VerifyForkHashes(chain.Config(), header, uncle); err != nil {
		return err
	}
	return nil
}

//
//func (ethash *Ethash) CheckRegister(address common.Address) (bool, error) {
//	var addr = minerbook.MainNetAddress
//	mb, err := contract.NewMinerBook(addr, nil)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a Token contract: %v", err)
//		return false, err
//	}
//	//TODO:
//	reg, err := mb.UsedHashedPubkey(nil, address)
//	if err != nil {
//		log.Fatalf("query registered error :%v", err)
//		return false, err
//	}
//	return reg, nil
//}

func (ethash *Ethash) GetReputationByState(chain consensus.ChainReader, address common.Address) uint64 {
	//conn, _ := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
	//reputation, _ := conn.ReputationAt(nil, address, nil)
	//println(chain)
	// when chain is nil, used to testing, return init reputation = 1000
	if chain == nil {
		return 1000
	}
	s, _ := chain.State()
	//if err != nil {
	//	return 0
	//}
	repCurrent := (s.GetReputation(address))
	return repCurrent
}

// According to the MinerBook contract, obtain the author's reputation.
// TODO:
//func (ethash *Ethash) GetReputationByContract(address common.Address) uint64 {
//	//var abi = contract.MinerBookABI
//	var addr = minerbook.MainNetAddress
//	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
//	//conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
//	//if err != nil {
//	//	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
//	//	return -1
//	//}
//	mb, err := contract.NewMinerBook(addr, nil)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a Token contract: %v", err)
//		return 0
//	}
//
//	//used, err := mb.UsedHashedPubkey(nil, crypto.Keccak256Hash(address[:]))
//	used, err := mb.UsedHashedPubkey(nil, address)
//	if err != nil {
//		log.Fatalf("query registered error :%v", err)
//		return 0
//	}
//	if used != true {
//		log.Fatalf("address is not registered")
//		return 0
//	}
//
//	reputation, err := mb.ReputationList(nil, address)
//	if err != nil {
//		log.Fatalf("query reputation error:%v", err)
//	}
//
//	return reputation
//
//	//var backend = contract.MinerBook
//	//var contract, err = contract.NewMinerBook(minerbook.MainNetAddress,ethash)
//	//minerbookcontract.
//	//return 0
//}
//func (ethash *Ethash) Register(conaddr common.Address, mineraddr common.Address, value uint64) (*types.Transaction, error) {
//	//var addr = minerbook.MainNetAddress
//	var abi = contract.MinerBookABI
//	//conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
//	//if err != nil {
//	//	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
//	//	return nil, err
//	//}
//	mb, err := contract.NewMinerBook(conaddr, nil)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a minerbook contract: %v", err)
//		return nil, err
//	}
//
//	// Create an authorized transactor and spend 1 unicorn
//	auth, err := bind.NewTransactor(strings.NewReader(abi), "123")
//	if err != nil {
//		log.Fatalf("Failed to create authorized transactor: %v", err)
//		return nil, err
//	}
//	tx, err := mb.Register(auth, mineraddr, mineraddr)
//	if err != nil {
//		log.Fatalf("Failed to request minerbook addreputation: %v", err)
//		return nil, err
//	}
//	return tx, nil
//}
//
////TODO:
//func (ethash *Ethash) AddReputation(address common.Address, value uint64) (*types.Transaction, error) {
//	var addr = minerbook.MainNetAddress
//	var abi = contract.MinerBookABI
//	//conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
//	//if err != nil {
//	//	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
//	//	return nil, err
//	//}
//	mb, err := contract.NewMinerBook(addr, nil)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a minerbook contract: %v", err)
//		return nil, err
//	}
//
//	// Create an authorized transactor and spend 1 unicorn
//	auth, err := bind.NewTransactor(strings.NewReader(abi), "123")
//	if err != nil {
//		log.Fatalf("Failed to create authorized transactor: %v", err)
//		return nil, err
//	}
//	tx, err := mb.AddReputation(auth, address[:], value)
//	if err != nil {
//		log.Fatalf("Failed to request minerbook addreputation: %v", err)
//		return nil, err
//	}
//	return tx, nil
//}
//
////TODO:
//func (ethash *Ethash) SubReputation(address common.Address, value uint64) (*types.Transaction, error) {
//	var addr = minerbook.MainNetAddress
//	var abi = contract.MinerBookABI
//	//conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
//	//if err != nil {
//	//	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
//	//	return nil, err
//	//}
//
//	mb, err := contract.NewMinerBook(addr, nil)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a minerbook contract: %v", err)
//		return nil, err
//	}
//
//	// Create an authorized transactor and spend 1 unicorn
//	auth, err := bind.NewTransactor(strings.NewReader(abi), "123")
//	if err != nil {
//		log.Fatalf("Failed to create authorized transactor: %v", err)
//		return nil, err
//	}
//	tx, err := mb.SubReputation(auth, address[:], value)
//	if err != nil {
//		log.Fatalf("Failed to request minerbook addreputation: %v", err)
//		return nil, err
//	}
//	return tx, nil
//}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (ethash *Ethash) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return CalcDifficulty(chain.Config(), time, parent)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(config *params.ChainConfig, time uint64, parent *types.Header) *big.Int {
	next := new(big.Int).Add(parent.Number, big1)
	switch {
	case config.IsConstantinople(next):
		return calcDifficultyConstantinople(time, parent)
	case config.IsByzantium(next):
		return calcDifficultyByzantium(time, parent)
	case config.IsHomestead(next):
		return calcDifficultyHomestead(time, parent)
	default:
		return calcDifficultyFrontier(time, parent)
	}
}

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)
)

// makeDifficultyCalculator creates a difficultyCalculator with the given bomb-delay.
// the difficulty is calculated with Byzantium rules, which differs from Homestead in
// how uncles affect the calculation
// TODO: Calculation formula should add the reputation.
func makeDifficultyCalculator(bombDelay *big.Int) func(time uint64, parent *types.Header) *big.Int {
	// Note, the calculations below looks at the parent number, which is 1 below
	// the block number. Thus we remove one from the delay given
	bombDelayFromParent := new(big.Int).Sub(bombDelay, big1)
	return func(time uint64, parent *types.Header) *big.Int {
		// https://github.com/ethereum/EIPs/issues/100.
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)

		bigTime := new(big.Int).SetUint64(time)
		bigParentTime := new(big.Int).Set(parent.Time)

		// holds intermediate values to make the algo easier to read & audit
		x := new(big.Int)
		y := new(big.Int)

		// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
		x.Sub(bigTime, bigParentTime)
		x.Div(x, big9)
		if parent.UncleHash == types.EmptyUncleHash {
			x.Sub(big1, x)
		} else {
			x.Sub(big2, x)
		}
		// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
		if x.Cmp(bigMinus99) < 0 {
			x.Set(bigMinus99)
		}
		// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
		x.Mul(y, x)
		x.Add(parent.Difficulty, x)

		// minimum difficulty can ever be (before exponential factor)
		if x.Cmp(params.MinimumDifficulty) < 0 {
			x.Set(params.MinimumDifficulty)
		}
		// calculate a fake block number for the ice-age delay
		// Specification: https://eips.ethereum.org/EIPS/eip-1234
		fakeBlockNumber := new(big.Int)
		if parent.Number.Cmp(bombDelayFromParent) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, bombDelayFromParent)
		}
		// for the exponential factor
		periodCount := fakeBlockNumber
		periodCount.Div(periodCount, expDiffPeriod)

		// the exponential factor, commonly referred to as "the bomb"
		// diff = diff + 2^(periodCount - 2)
		if periodCount.Cmp(big1) > 0 {
			y.Sub(periodCount, big2)
			y.Exp(big2, y, nil)
			x.Add(x, y)
		}
		return x
	}
}

// calcDifficultyHomestead is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time given the
// parent block's time and difficulty. The calculation uses the Homestead rules.
func calcDifficultyHomestead(time uint64, parent *types.Header) *big.Int {
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 10
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big10)
	x.Sub(big1, x)

	// max(1 - (block_timestamp - parent_timestamp) // 10, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// (parent_diff + parent_diff // 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	// for the exponential factor
	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// calcDifficultyFrontier is the difficulty adjustment algorithm. It returns the
// difficulty that a new block should have when created at time given the parent
// block's time and difficulty. The calculation uses the Frontier rules.
func calcDifficultyFrontier(time uint64, parent *types.Header) *big.Int {
	diff := new(big.Int)
	adjust := new(big.Int).Div(parent.Difficulty, params.DifficultyBoundDivisor)
	bigTime := new(big.Int)
	bigParentTime := new(big.Int)

	bigTime.SetUint64(time)
	bigParentTime.Set(parent.Time)

	if bigTime.Sub(bigTime, bigParentTime).Cmp(params.DurationLimit) < 0 {
		diff.Add(parent.Difficulty, adjust)
	} else {
		diff.Sub(parent.Difficulty, adjust)
	}
	if diff.Cmp(params.MinimumDifficulty) < 0 {
		diff.Set(params.MinimumDifficulty)
	}

	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)
	if periodCount.Cmp(big1) > 0 {
		// diff = diff + 2^(periodCount - 2)
		expDiff := periodCount.Sub(periodCount, big2)
		expDiff.Exp(big2, expDiff, nil)
		diff.Add(diff, expDiff)
		diff = math.BigMax(diff, params.MinimumDifficulty)
	}
	return diff
}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW difficulty requirements.
func (ethash *Ethash) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return ethash.verifySeal(chain, header, false)
}

// verifySeal checks whether a block satisfies the PoW difficulty requirements,
// either using the usual ethash cache for it, or alternatively using a full DAG
// to make remote mining fast.
func (ethash *Ethash) verifySeal(chain consensus.ChainReader, header *types.Header, fulldag bool) error {
	// If we're running a fake PoW, accept any seal as valid
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
		time.Sleep(ethash.fakeDelay)
		if ethash.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	// If we're running a shared PoW, delegate verification to it
	if ethash.shared != nil {
		return ethash.shared.verifySeal(chain, header, fulldag)
	}
	// Ensure that we have a valid difficulty for the block
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW values
	number := header.Number.Uint64()

	var (
		digest []byte
		result []byte
	)
	// If fast-but-heavy PoW verification was requested, use an ethash dataset
	if fulldag {
		dataset := ethash.dataset(number, true)
		if dataset.generated() {
			digest, result = hashimotoFull(dataset.dataset, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

			// Datasets are unmapped in a finalizer. Ensure that the dataset stays alive
			// until after the call to hashimotoFull so it's not unmapped while being used.
			runtime.KeepAlive(dataset)
		} else {
			// Dataset not yet generated, don't hang, use a cache instead
			fulldag = false
		}
	}
	// If slow-but-light PoW verification was requested (or DAG not yet ready), use an ethash cache
	if !fulldag {
		cache := ethash.cache(number)

		size := datasetSize(number)
		if ethash.config.PowMode == ModeTest {
			size = 32 * 1024
		}
		digest, result = hashimotoLight(size, cache.cache, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

		// Caches are unmapped in a finalizer. Ensure that the cache stays alive
		// until after the call to hashimotoLight so it's not unmapped while being used.
		runtime.KeepAlive(cache)
	}
	// Verify the calculated values against the ones provided in the header
	if !bytes.Equal(header.MixDigest[:], digest) {
		return errInvalidMixDigest
	}

	////NEW change: add reputation
	author, err := ethash.Author(header)
	if err != nil {
		return fmt.Errorf("invalid Author")
	}
	//reputation := ethash.GetReputationByContract(author)
	//s = state.
	reputation := ethash.GetReputationByState(chain, author)
	//reputation := uint64(1000)
	if reputation <= ReputationLowThreshold {
		if author == ReputationWhiteAddress {
			reputation = ReputationInit
		} else {
			return fmt.Errorf("reputation is too low")
		}
	}

	parentHeader := header
	authorAccount := 0
	for i := 0; i < ReputationCalcDiffBlockCount; i++ {
		//println(parentHeader.ParentHash.String())
		//header, _ := conn.HeaderByHash(nil, parentHeader.ParentHash)
		//println("当前："+parentHeader.Coinbase.String())
		parentHeader = chain.GetHeaderByHash(parentHeader.ParentHash)

		if parentHeader == nil {
			break
		}
		//println("前一个："+parentHeader.Coinbase.String())
		iString := parentHeader.Coinbase.String()
		if strings.Compare(author.String(), iString) == 0 {
			authorAccount += 1
		}
	}
	//竞争周期内出块数量与可用信誉度反比，指数递减，
	reputation = uint64(float64(reputation) / (m.Pow(ReputationContinuousBlockUsable, float64(authorAccount))))
	//target  := new(big.Int).Div(two256, header.Difficulty)
	//reputation := ethash.state.
	target := new(big.Int)
	//var repratio  = reputation/ReputationInit

	//target = new(big.Int).Div(two256, new(big.Int).Sub(header.Difficulty, ))
	if reputation == ReputationInit {
		target = new(big.Int).Div(two256, header.Difficulty)
	}
	if reputation > ReputationInit {
		tmp := new(big.Int).Mul(header.Difficulty, new(big.Int).SetUint64(reputation-ReputationInit))
		target = new(big.Int).Div(two256, new(big.Int).Sub(header.Difficulty, new(big.Int).Div(tmp, new(big.Int).SetUint64(ReputationInit*ReputationDifficultyRadio))))
	}
	if reputation < ReputationInit {
		tmp := new(big.Int).Mul(header.Difficulty, new(big.Int).SetUint64(ReputationInit-reputation))
		target = new(big.Int).Div(two256, new(big.Int).Add(header.Difficulty, new(big.Int).Div(tmp, new(big.Int).SetUint64(ReputationInit*ReputationDifficultyRadio))))
	}

	//if reputation >= ReputationInit {
	//	target = new(big.Int).Div(two256, new(big.Int).Sub(header.Difficulty, new(big.Int).SetUint64((reputation-ReputationInit)*repbase)))
	//} else {
	//	target = new(big.Int).Div(two256, new(big.Int).Add(header.Difficulty, new(big.Int).SetUint64((reputation-ReputationInit)*repbase)))
	//}
	//println(target.String())

	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to the ethash protocol. The changes are done inline.
func (ethash *Ethash) Prepare(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	//author, err := ethash.Author(header)
	//if err != nil {
	//	return fmt.Errorf("invalid Author")
	//}
	//authorReputation := GetReputation(author)
	header.Difficulty = ethash.CalcDifficulty(chain, header.Time.Uint64(), parent)
	return nil
}

// Finalize implements consensus.Engine, accumulating the block and uncle rewards,
// setting the final state and assembling the block.
func (ethash *Ethash) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block and uncle rewards and commit the final state root
	accumulateRewards(chain, state, header, uncles)

	////add contract reputation
	//repreward := getReputationRewards(state, header)
	//reptx, _ := ethash.AddReputation(header.Coinbase, repreward)
	//txs = append(txs, reptx)

	//if new(big.Int).Mod(header.Number, new(big.Int).SetInt64(int64(ReputationBlackBlockCount))).Int64() == 0 {
	//	dereptx, _:= ethash.reputationDecayByContract(state, header)
	//	for _,_tx := range dereptx{
	//		txs = append(txs, _tx)
	//	}
	//}

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	//println("con" + header.Number.String())
	//println("con" + header.Root.String())

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, uncles, receipts), nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (ethash *Ethash) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	})
	hasher.Sum(hash[:0])
	return hash
}

// Some weird constants to avoid constant memory allocs for them.
var (
	big8  = big.NewInt(8)
	big32 = big.NewInt(32)
)

func getReputationRewards(chain consensus.ChainReader, state *state.StateDB, header *types.Header) uint64 {

	author := header.Coinbase
	if author == ReputationWhiteAddress {
		return 0
	}
	authorString := author.String()
	authorAcount := 0
	parentHeader := new(types.Header)

	//conn, _ := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
	//ctx := context.Background()
	//count := 0
	//if parentHeader.Number.Cmp(big.NewInt(int64(ReputationFrontierBlockCount))) <= 0{
	//	count = int(parentHeader.Number.Int64())-1
	//}else{
	//	count = ReputationFrontierBlockCount
	//}
	//println(author.String())
	parentHeader = header
	for i := 0; i < ReputationFrontierBlockCount; i++ {
		//println(parentHeader.ParentHash.String())
		//header, _ := conn.HeaderByHash(nil, parentHeader.ParentHash)
		//println("当前："+parentHeader.Coinbase.String())
		parentHeader = chain.GetHeaderByHash(parentHeader.ParentHash)

		if parentHeader == nil {
			break
		}
		//println("前一个："+parentHeader.Coinbase.String())
		iString := parentHeader.Coinbase.String()
		if strings.Compare(authorString, iString) == 0 {
			authorAcount += 1
		}
	}

	repCurrent := int(state.GetReputation(author))
	repRadio := float64(repCurrent) / float64(ReputationHighThreshold)
	expected := float64(0)
	if ReputationCalcDiffBlockCount == 0 || ReputationFrontierBlockCount == 0 {
		expected = 0
	} else {
		expected = float64(authorAcount*ReputationCalcDiffBlockCount) / (float64(ReputationFrontierBlockCount) * repRadio * float64(ReputationExpectedRewardAccount))
	}

	if expected > 1 {
		expected = 1
	}
	//println(repCurrent)
	repReward := float64(1-expected) * float64(int(ReputationHighThreshold)-repCurrent) / float64(ReputationRwardFormulaOptimizeParam)
	if repCurrent+int(repReward) > int(ReputationHighThreshold) {
		return ReputationHighThreshold - uint64(repCurrent)
	}
	//println(authorAcount)
	//println(repReward)
	return uint64(repReward)
}

func reputationDecay(chain consensus.ChainReader, state *state.StateDB, header *types.Header) error {
	// 每100个区块调用一次
	// 调用合约，返回miner列表,减少一些信誉值。
	var addr = ReputationContractAddress
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	//conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
	//if err != nil {
	//	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	//	return err
	//}
	mb, err := MBC.NewMinerBook(addr, nil)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
		return err
	}
	//TODO:获取miner列表，信誉衰减
	miners, err := mb.GetMiners(nil)
	if err != nil {
		log.Fatalf("get miners error :%v", err)
		return err
	}
	var minerList = make(map[common.Address]int)
	for _, miner := range miners {
		minerList[miner] = 0
	}
	parentHeader := header
	for i := 0; i < ReputationBlackBlockCount; i++ {
		parentHeader = chain.GetHeaderByHash(parentHeader.ParentHash)
		if parentHeader == nil {
			break
		}
		//mineriString := parentHeader.Coinbase.String()
		mineraddr := parentHeader.Coinbase
		minerList[mineraddr] += 1
	}
	for miner, mineraccount := range minerList {
		repCurrent := int(state.GetReputation(miner))
		repRadio := float64(repCurrent) / float64(ReputationHighThreshold)
		expected := float64(0)
		if ReputationCalcDiffBlockCount == 0 {
			expected = 0
		} else {
			expected = float64(mineraccount*ReputationCalcDiffBlockCount) / (float64(ReputationBlackBlockCount) * repRadio * float64(ReputationExpectedDecayAccount))
		}

		if expected > 1 {
			expected = 1
		}
		repDecay := int((1 - expected) * float64(repCurrent) / float64(ReputationDecayFormulaOptimizeParam))
		if repCurrent < repDecay {
			state.SubReputation(miner, uint64(repCurrent))
			//TODO：加入黑名单！
		} else {
			state.SubReputation(miner, uint64(repDecay))
		}

	}
	return nil
}

func reputationDecayTest(chain consensus.ChainReader, state *state.StateDB, header *types.Header) error {

	miners := MinersListTest
	var minerList = make(map[common.Address]int)
	for _, miner := range miners {
		minerList[miner] = 0
	}
	parentHeader := header
	for i := 0; i < ReputationBlackBlockCount; i++ {
		parentHeader = chain.GetHeaderByHash(parentHeader.ParentHash)
		if parentHeader == nil {
			break
		}
		//mineriString := parentHeader.Coinbase.String()
		mineraddr := parentHeader.Coinbase
		minerList[mineraddr] += 1
	}
	for miner, mineraccount := range minerList {
		repCurrent := int(state.GetReputation(miner))
		repRadio := float64(repCurrent) / float64(ReputationHighThreshold)
		expected := float64(0)
		if ReputationCalcDiffBlockCount == 0 {
			expected = 0
		} else {
			expected = float64(mineraccount*ReputationCalcDiffBlockCount) / (float64(ReputationBlackBlockCount) * repRadio * float64(ReputationExpectedDecayAccount))
		}
		if expected > 1 {
			expected = 1
		}
		repDecay := int((1 - expected) * float64(repCurrent) / float64(ReputationDecayFormulaOptimizeParam))
		if repCurrent < repDecay {
			state.SubReputation(miner, uint64(repCurrent))
			//TODO：加入黑名单！
		} else {
			state.SubReputation(miner, uint64(repDecay))
		}
	}
	return nil
}

//func (ethash *Ethash) reputationDecayByContract(state *state.StateDB, header *types.Header) ([]*types.Transaction, error) {
//	// 每100个区块调用一次
//	// 调用合约，返回miner列表,减少一些信誉值。
//	var addr = minerbook.MainNetAddress
//	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
//	conn, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
//	if err != nil {
//		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
//		return nil, err
//	}
//	mb, err := contract.NewMinerBook(addr, conn)
//	if err != nil {
//		log.Fatalf("Failed to instantiate a Token contract: %v", err)
//		return nil, err
//	}
//	//TODO:
//	minerMap, err := mb.GetMiners(nil)
//	if err != nil {
//		log.Fatalf("query registered error :%v", err)
//		return nil, err
//	}
//	var minerList = make(map[common.Address]int)
//	for miner, eable := range minerMap {
//		if eable == true {
//			//mineraddr := crypto.Keccak256Hash(miner[:]).Bytes()
//			minerList[miner] = 0
//		}
//	}
//	parentHeader := new(types.Header)
//	for i := 0; i < ReputationBlackBlockCount; i++ {
//		parentHeader, _ = conn.HeaderByHash(nil, parentHeader.ParentHash)
//		//mineriString := parentHeader.Coinbase.String()
//		mineraddr := parentHeader.Coinbase
//		minerList[mineraddr] += 1
//	}
//	var txs []*types.Transaction
//	for miner, mineraccount := range minerList {
//		repCurrent := int(state.GetReputation(miner))
//		repDecay := (1 - mineraccount/ReputationFrontierBlockCount) * repCurrent / ReputationDecayFormulaOptimizeParam
//		reptx, _ := ethash.SubReputation(header.Coinbase, uint64(repDecay))
//		txs = append(txs, reptx)
//	}
//	return txs, nil
//}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
// TODO: updated reputation
func accumulateRewards(chain consensus.ChainReader, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	// Select the correct block reward based on chain progression
	blockReward := FrontierBlockReward
	config := *chain.Config()
	//blockReputationReward := FrontierBlockReputationReward
	if config.IsByzantium(header.Number) {
		blockReward = ByzantiumBlockReward
	}
	if config.IsConstantinople(header.Number) {
		blockReward = ConstantinopleBlockReward
	}
	// Accumulate the rewards for the miner and any included uncles
	reward := new(big.Int).Set(blockReward)

	r := new(big.Int)
	//rr := 0
	for _, uncle := range uncles {
		r.Add(uncle.Number, big8)
		r.Sub(r, header.Number)
		r.Mul(r, blockReward)
		r.Div(r, big8)
		state.AddBalance(uncle.Coinbase, r)

		r.Div(blockReward, big32)
		reward.Add(reward, r)

		// accumulate the reputation rewards
		//rr.Add(uncle.Number, big8)
		//rr.Sub(rr, header.Number)
		//rr.Mul(rr, blockReputationReward)
		//rr.Div(rr, big8)
		//state.AddReputation(uncle.Coinbase, rr)
		//
		//rr.Div(blockReputationReward, big32)
		//rreward.Add(rreward, rr)
	}
	state.AddBalance(header.Coinbase, reward)

	repReward := getReputationRewards(chain, state, header)
	state.AddReputation(header.Coinbase, repReward)

	//目前所用的是Test方法
	if header.Number.Cmp(new(big.Int).SetInt64(0)) != 0 && ReputationBlackBlockCount != 0 && new(big.Int).Mod(header.Number, new(big.Int).SetInt64(int64(ReputationBlackBlockCount))).Int64() == 0 {
		//println("decay------------------------------")
		reputationDecayTest(chain, state, header)
	}

}
