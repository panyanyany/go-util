package contract

import (
	"fmt"
	"math/big"

	unit "github.com/DeOne4eg/eth-unit-converter"
	"github.com/cihub/seelog"
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
	"github.com/panyanyany/go-web3/wallet"
)

// Tx is a transaction object
type Tx struct {
	*web3.Transaction
	Receipt            *web3.Receipt
	Contract           *Contract
	Args               []interface{}
	Method             string
	Key                *wallet.Key
	GasPriceMultiplier uint64
}

func NewTx() *Tx {
	return &Tx{Transaction: &web3.Transaction{}, GasPriceMultiplier: 1}
}

func (t *Tx) SetGasPriceMultiplier(m uint64) *Tx {
	t.GasPriceMultiplier = m
	return t
}
func (t *Tx) SetKey(key *wallet.Key) *Tx {
	t.Key = key
	return t
}

func (t *Tx) SetInput(data []byte) *Tx {
	t.Input = data
	return t
}
func (t *Tx) Validate() (err error) {
	// method & args
	//methodName := "swapExactETHForTokens"
	method := t.Contract.Abi.Methods[t.Method]
	data, err := abi.Encode(t.Args, method.Inputs)
	if err != nil {
		err = fmt.Errorf("aib.Encode(): %w", err)
		return
	}
	t.Input = append(method.ID(), data...)
	return
}

func (t *Tx) SetMethod(method string) *Tx {
	t.Method = method
	return t
}

func (t *Tx) SetContract(c *Contract) *Tx {
	t.Contract = c
	t.To = &c.Address
	return t
}

func (t *Tx) AddArgs(args ...interface{}) *Tx {
	t.Args = args
	return t
}

// SetValue sets the value for the txn
func (t *Tx) SetValue(v *big.Int) *Tx {
	t.Value = new(big.Int).Set(v)
	return t
}

// EstimateGas estimates the gas for the call
func (t *Tx) EstimateGas() (uint64, error) {
	if err := t.Validate(); err != nil {
		err = fmt.Errorf("t.Validate: %w", err)
		return 0, err
	}
	return t.estimateGas()
}

func (t *Tx) estimateGas() (uint64, error) {
	msg := &web3.CallMsg{
		From:  t.From,
		To:    t.To,
		Data:  t.Input,
		Value: t.Value,
	}
	return t.Contract.Provider.EstimateGas(msg)
}

// DoAndWait is a blocking query that combines
// both Do and Wait functions
func (t *Tx) DoAndWait() error {
	if err := t.Do(); err != nil {
		return err
	}
	if err := t.Wait(); err != nil {
		return err
	}
	return nil
}
func (t *Tx) DoRaw() (err error) {
	chErr := make(chan error)
	chCnt := 0
	// estimate gas price
	if t.GasPrice == 0 {
		chCnt++
		go func() {
			var err error
			t.GasPrice, err = t.Contract.Provider.GasPrice()
			t.GasPrice *= t.GasPriceMultiplier
			seelog.Debugf("get GasPrice: %v", unit.NewWei(big.NewInt(int64(t.GasPrice))))
			if err != nil {
				err = fmt.Errorf("t.Contract.Provider.GasPrice(): %w", err)
			}
			chErr <- err
		}()
	}
	// estimate gas limit
	if t.Gas == 0 {
		chCnt++
		go func() {
			var err error
			t.Gas, err = t.estimateGas()
			t.Gas = t.Gas * 150 / 100 // 必须调为 150% 否则可能失败
			seelog.Debugf("get Gas: %v", unit.NewWei(big.NewInt(int64(t.Gas))))
			if err != nil {
				err = fmt.Errorf("t.estimateGas(): %w", err)
			}
			chErr <- err
		}()
	}

	// nonce
	chCnt++
	go func() {
		var err error
		blockNumber, err := t.Contract.Provider.BlockNumber()
		if err != nil {
			err = fmt.Errorf("t.Contract.Provider.BlockNumber(): %w", err)
			chErr <- err
			return
		}

		nonce, err := t.Contract.Provider.GetNonce(t.Key.Address(), web3.BlockNumber(blockNumber))
		if err != nil {
			err = fmt.Errorf("nonce: %w", err)
			chErr <- err
			return
		}
		t.Nonce = nonce
		chErr <- err
	}()

	errCnt := 0
	for goErr := range chErr {
		if goErr != nil {
			err = goErr
			return
		}
		errCnt++
		if errCnt >= chCnt {
			close(chErr)
		}
	}

	// Create the signer object and sign
	t.Transaction, _ = wallet.NewEIP155Signer(56).SignTx(t.Transaction, t.Key)

	// Send the signed transaction
	data := t.Transaction.MarshalRLP()
	//if t.addr != nil {
	//	txn.To = t.addr
	//}
	t.Hash, err = t.Contract.Provider.SendRawTransaction(data)
	if err != nil {
		err = fmt.Errorf("t.Contract.Provider.SendTransaction: %w", err)
		return
	}
	return nil
}

// Do sends the transaction to the network
func (t *Tx) Do() (err error) {
	err = t.Validate()
	if err != nil {
		err = fmt.Errorf("t.Validate: %w", err)
		return err
	}

	return t.DoRaw()
}

// SetGasPrice sets the gas price of the transaction
func (t *Tx) SetGasPrice(gasPrice uint64) *Tx {
	t.GasPrice = gasPrice
	return t
}

// SetGasLimit sets the gas limit of the transaction
func (t *Tx) SetGas(gasLimit uint64) *Tx {
	t.Gas = gasLimit
	return t
}

// Wait waits till the transaction is mined
func (t *Tx) Wait() error {
	if (t.Hash == web3.Hash{}) {
		panic("transaction not executed")
	}

	var err error
	for {
		t.Receipt, err = t.Contract.Provider.GetTransactionReceipt(t.Hash)
		if err != nil {
			if err.Error() != "not found" {
				return err
			}
		}
		if t.Receipt != nil {
			break
		}
	}
	return nil
}
