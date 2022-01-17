package contract

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path"

	"tugou-hunter/util/os_util"

	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
	"github.com/panyanyany/go-web3/jsonrpc"
)

// Contract is an Ethereum contract
type Contract struct {
	Name      string
	Symbol    string
	Address   web3.Address
	Decimals  int
	ChainName string
	From      *web3.Address
	Abi       *abi.ABI
	Provider  jsonrpc.IEth
}

func NewTargetContract(address string) *Contract {
	return &Contract{
		Name:   "Target",
		Symbol: "TARGET",
		//Address:   "0x3753301611c7D2f352d28151D14d915492C6940F",
		Address:   web3.HexToAddress(address),
		Decimals:  18,
		ChainName: "",
	}
}

func NewContractSimple(name string, symbol string, address string, decimals int, chainName string) (c *Contract) {
	c = new(Contract)
	c.Name = name
	c.Symbol = symbol
	c.Address = web3.HexToAddress(address)
	c.Decimals = decimals
	c.ChainName = chainName

	return
}

// DeployContract deploys a contract
func DeployContract(provider jsonrpc.IEth, from web3.Address, abi *abi.ABI, bin []byte, args ...interface{}) *Txn {
	return &Txn{
		from:     from,
		provider: provider,
		method:   abi.Constructor,
		args:     args,
		bin:      bin,
	}
}

// NewContract creates a new contract instance
func NewContract(addr web3.Address, abi *abi.ABI, provider jsonrpc.IEth) *Contract {
	return &Contract{
		Address:  addr,
		Abi:      abi,
		Provider: provider,
	}
}

func (r *Contract) LoadAbi() (err error) {
	var bs []byte
	bs, err = ioutil.ReadFile(path.Join(os_util.MustGetExecutableDir(), "..", fmt.Sprintf("resources/%s/%s/abi.json",
		r.ChainName,
		r.Name,
	)))
	if err != nil {
		return
	}

	var abiObj *abi.ABI
	abiObj, err = abi.NewABI(string(bs))
	if err != nil {
		err = fmt.Errorf("abi.NewABI: %w", err)
		return
	}
	r.Abi = abiObj
	return
}

// ABI returns the Abi of the contract
func (c *Contract) ABI() *abi.ABI {
	return c.Abi
}

// Addr returns the address of the contract
func (c *Contract) Addr() web3.Address {
	return c.Address
}

// SetFrom sets the origin of the calls
func (c *Contract) SetFrom(addr web3.Address) {
	c.From = &addr
}

// EstimateGas estimates the gas for a contract call
func (c *Contract) EstimateGas(method string, args ...interface{}) (uint64, error) {
	return c.Txn(method, args).EstimateGas()
}

// Call calls a method in the contract
func (c *Contract) Call(method string, block web3.BlockNumber, args ...interface{}) (map[string]interface{}, error) {
	m, ok := c.Abi.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %s not found in Contract.Abi.Methods[method]", method)
	}

	// Encode input
	data, err := abi.Encode(args, m.Inputs)
	if err != nil {
		err = fmt.Errorf("aib.Encode(): %w", err)
		return nil, err
	}
	data = append(m.ID(), data...)

	// Call function
	msg := &web3.CallMsg{
		To:   &c.Address,
		Data: data,
	}
	if c.From != nil {
		msg.From = *c.From
	}

	rawStr, err := c.Provider.Call(msg, block)
	if err != nil {
		err = fmt.Errorf("Contract.Provider.Call(): %w", err)
		return nil, err
	}

	// Decode output
	raw, err := hex.DecodeString(rawStr[2:])
	if err != nil {
		err = fmt.Errorf("hex.DecodeString: %w", err)
		return nil, err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	respInterface, err := abi.Decode(m.Outputs, raw)
	if err != nil {
		err = fmt.Errorf("Abi.Decode: %w", err)
		return nil, err
	}

	resp := respInterface.(map[string]interface{})
	return resp, nil
}

// Txn creates a new transaction object
func (c *Contract) Txn(method string, args ...interface{}) *Txn {
	m, ok := c.Abi.Methods[method]
	if !ok {
		// TODO, return error
		panic(fmt.Errorf("method %s not found", method))
	}

	return &Txn{
		from:     *c.From,
		addr:     &c.Address,
		provider: c.Provider,
		method:   m,
		args:     args,
	}
}

// Event returns a specific event
func (c *Contract) Event(name string) (*Event, bool) {
	event, ok := c.Abi.Events[name]
	if !ok {
		return nil, false
	}
	return &Event{event}, true
}
