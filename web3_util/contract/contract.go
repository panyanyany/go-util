package contract

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path"

	"tugou-hunter/util/os_util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
	"github.com/panyanyany/go-web3/contract/builtin/erc20"
	"github.com/panyanyany/go-web3/jsonrpc"
)

// Contract is an Ethereum contract
type Contract struct {
	Name      string
	Symbol    string
	Address   common.Address
	Decimals  int
	ChainName string
	From      *common.Address
	Abi       *abi.ABI
	Provider  jsonrpc.IEth
}

func NewTargetContract(address string) *Contract {
	return &Contract{
		Name:   "Target",
		Symbol: "TARGET",
		//Address:   "0x3753301611c7D2f352d28151D14d915492C6940F",
		Address:   common.HexToAddress(address),
		Decimals:  18,
		ChainName: "",
	}
}

func NewContractSimple(name string, symbol string, address string, decimals int, chainName string) (c *Contract) {
	c = new(Contract)
	c.Name = name
	c.Symbol = symbol
	c.Address = common.HexToAddress(address)
	c.Decimals = decimals
	c.ChainName = chainName

	return
}

func (r *Contract) AsErc20() *erc20.ERC20 {
	return erc20.NewERC20(web3.HexToAddress(r.Address.String()), r.Provider)
}

func (r *Contract) MustLoadAbi() {
	err := r.LoadAbi()
	if err != nil {
		panic(err)
		return
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
func (r *Contract) ABI() *abi.ABI {
	return r.Abi
}

// Addr returns the address of the contract
func (r *Contract) Addr() common.Address {
	return r.Address
}

// SetFrom sets the origin of the calls
func (r *Contract) SetFrom(addr common.Address) {
	r.From = &addr
}

// EstimateGas estimates the gas for a contract call
func (r *Contract) EstimateGas(method string, args ...interface{}) (uint64, error) {
	return r.Txn(method, args).EstimateGas()
}

// Call calls a method in the contract
func (r *Contract) Call(method string, block web3.BlockNumber, args ...interface{}) (map[string]interface{}, error) {
	m, ok := r.Abi.Methods[method]
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

	to := web3.HexToAddress(r.Address.String())

	// Call function
	msg := &web3.CallMsg{
		To:   &to,
		Data: data,
	}
	if r.From != nil {
		from := web3.HexToAddress(r.From.String())
		msg.From = from
	}

	rawStr, err := r.Provider.Call(msg, block)
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
func (r *Contract) Txn(method string, args ...interface{}) *Txn {
	m, ok := r.Abi.Methods[method]
	if !ok {
		// TODO, return error
		panic(fmt.Errorf("method %s not found", method))
	}

	to := web3.HexToAddress(r.Address.String())
	from := web3.HexToAddress(r.From.String())
	return &Txn{
		from:     from,
		addr:     &to,
		provider: r.Provider,
		method:   m,
		args:     args,
	}
}

// Event returns a specific event
func (r *Contract) Event(name string) (*Event, bool) {
	event, ok := r.Abi.Events[name]
	if !ok {
		return nil, false
	}
	return &Event{event}, true
}
