package contract

import (
	"fmt"
	"math/big"

	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
	"github.com/panyanyany/go-web3/jsonrpc"
)

// Txn is a transaction object
type Txn struct {
	from     web3.Address
	addr     *web3.Address
	provider jsonrpc.IEth
	method   *abi.Method
	args     []interface{}
	data     []byte
	bin      []byte
	gasLimit uint64
	gasPrice uint64
	value    *big.Int
	hash     web3.Hash
	receipt  *web3.Receipt
}

func (t *Txn) isContractDeployment() bool {
	return t.bin != nil
}

// AddArgs is used to set the arguments of the transaction
func (t *Txn) AddArgs(args ...interface{}) *Txn {
	t.args = args
	return t
}

// SetValue sets the value for the txn
func (t *Txn) SetValue(v *big.Int) *Txn {
	t.value = new(big.Int).Set(v)
	return t
}

// EstimateGas estimates the gas for the call
func (t *Txn) EstimateGas() (uint64, error) {
	if err := t.Validate(); err != nil {
		return 0, err
	}
	return t.estimateGas()
}

func (t *Txn) estimateGas() (uint64, error) {
	if t.isContractDeployment() {
		return t.provider.EstimateGasContract(t.data)
	}

	msg := &web3.CallMsg{
		From:  t.from,
		To:    t.addr,
		Data:  t.data,
		Value: t.value,
	}
	return t.provider.EstimateGas(msg)
}

// DoAndWait is a blocking query that combines
// both Do and Wait functions
func (t *Txn) DoAndWait() error {
	if err := t.Do(); err != nil {
		return err
	}
	if err := t.Wait(); err != nil {
		return err
	}
	return nil
}

// Do sends the transaction to the network
func (t *Txn) Do() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	// estimate gas price
	if t.gasPrice == 0 {
		t.gasPrice, err = t.provider.GasPrice()
		if err != nil {
			return err
		}
	}
	// estimate gas limit
	if t.gasLimit == 0 {
		t.gasLimit, err = t.estimateGas()
		if err != nil {
			return err
		}
	}

	// send transaction
	txn := &web3.Transaction{
		From:     t.from,
		Input:    t.data,
		GasPrice: t.gasPrice,
		Gas:      t.gasLimit,
		Value:    t.value,
	}
	if t.addr != nil {
		txn.To = t.addr
	}
	t.hash, err = t.provider.SendTransaction(txn)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates the arguments of the transaction
func (t *Txn) Validate() error {
	if t.data != nil {
		// Already validated
		return nil
	}
	if t.isContractDeployment() {
		t.data = append(t.data, t.bin...)
	}
	if t.method != nil {
		data, err := abi.Encode(t.args, t.method.Inputs)
		if err != nil {
			return fmt.Errorf("failed to encode arguments: %v", err)
		}
		if !t.isContractDeployment() {
			t.data = append(t.method.ID(), data...)
		} else {
			t.data = append(t.data, data...)
		}
	}
	return nil
}

// SetGasPrice sets the gas price of the transaction
func (t *Txn) SetGasPrice(gasPrice uint64) *Txn {
	t.gasPrice = gasPrice
	return t
}

// SetGasLimit sets the gas limit of the transaction
func (t *Txn) SetGasLimit(gasLimit uint64) *Txn {
	t.gasLimit = gasLimit
	return t
}

// Wait waits till the transaction is mined
func (t *Txn) Wait() error {
	if (t.hash == web3.Hash{}) {
		panic("transaction not executed")
	}

	var err error
	for {
		t.receipt, err = t.provider.GetTransactionReceipt(t.hash)
		if err != nil {
			if err.Error() != "not found" {
				return err
			}
		}
		if t.receipt != nil {
			break
		}
	}
	return nil
}

// Receipt returns the receipt of the transaction after wait
func (t *Txn) Receipt() *web3.Receipt {
	return t.receipt
}
