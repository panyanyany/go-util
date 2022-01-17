package web3_util

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/panyanyany/go-web3"
	"github.com/panyanyany/go-web3/abi"
	"github.com/panyanyany/go-web3/contract"
	"github.com/panyanyany/go-web3/jsonrpc"
)

type Asset struct {
	Name      string
	Symbol    string
	Address   string
	Decimals  int
	ChainName string
}

func NewTargetAsset(address string) *Asset {
	return (&Asset{
		Name:   "Target",
		Symbol: "TARGET",
		//Address:   "0x3753301611c7D2f352d28151D14d915492C6940F",
		Address:   address,
		Decimals:  18,
		ChainName: "",
	}).Init()
}

func (r *Asset) Init() *Asset {
	r.Address = strings.ToLower(r.Address)
	return r
}
func (r *Asset) DownloadAbi() (body string, err error) {
	return
}

func (r *Asset) LoadAbi() (body string, err error) {
	var bs []byte
	bs, err = ioutil.ReadFile(fmt.Sprintf("resources/%s/%s/abi.json",
		r.ChainName,
		r.Name,
	))
	if err != nil {
		return
	}
	body = string(bs)
	return
}

func (r *Asset) ToContract(client jsonrpc.IClient) (c *contract.Contract, err error) {
	var abiStr string
	abiStr, err = r.LoadAbi()
	if err != nil {
		err = fmt.Errorf("r.LoadAbi(): %w", err)
		return
	}

	var abiObj *abi.ABI
	abiObj, err = abi.NewABI(abiStr)
	if err != nil {
		err = fmt.Errorf("abi.NewABI: %w", err)
		return
	}

	c = contract.NewContract(web3.HexToAddress(r.Address), abiObj, &jsonrpc.Eth{Client: client})
	return
}
