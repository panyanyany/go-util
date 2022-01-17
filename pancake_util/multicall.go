package pancake_util

import (
	"fmt"
	"math"
	"math/big"

	"tugou-hunter/util/struct_util"
	"tugou-hunter/util/web3_util/bsc"
	"tugou-hunter/util/web3_util/contract"

	"github.com/panyanyany/go-web3"
)

type MultiCallRepo struct {
	Contract *contract.Contract
}

func NewMultiCallRepo(contract *contract.Contract) *MultiCallRepo {
	mc := &MultiCallRepo{contract}
	return mc
}

func (mc *MultiCallRepo) GetTokenInfo(tokenAddress string) (*Token, error) {
	resp, err := mc.Contract.Call("getTokenInfo", web3.Latest, web3.HexToAddress(tokenAddress))
	if err != nil {
		fmt.Println("getTokenInfo err:", err)
		return nil, err
	}

	token := Token{}
	err = struct_util.Map2Struct(resp["tokenInfo"], &token)
	if err != nil {
		err = fmt.Errorf("MultiCallRepo.GetTokenInfo-Map2Struct: %w", err)
		return nil, err
	}

	return &token, nil
}

func (mc *MultiCallRepo) GetTokenInfos(tokenAddresslist []string) ([]*Token, error) {
	addresslist := make([]web3.Address, 0)
	for _, addr := range tokenAddresslist {
		addresslist = append(addresslist, web3.HexToAddress(addr))
	}

	resp, err := mc.Contract.Call("getTokenInfos", web3.Latest, addresslist)
	if err != nil {
		fmt.Println("getTokenInfos err:", err)
		return nil, err
	}

	tokens := make([]*Token, 0)
	err = struct_util.Map2Struct(resp["tokenInfo"], &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (mc *MultiCallRepo) GetBalances(tokenAddresslist []string, userAddresslist []string) ([]*big.Int, error) {
	tokenAddrList := make([]web3.Address, 0)
	for _, addr := range tokenAddresslist {
		tokenAddrList = append(tokenAddrList, web3.HexToAddress(addr))
	}

	userAddrList := make([]web3.Address, 0)
	for _, addr := range userAddresslist {
		userAddrList = append(userAddrList, web3.HexToAddress(addr))
	}

	resp, err := mc.Contract.Call("getUserBalances", web3.Latest, tokenAddrList, userAddrList)
	if err != nil {
		fmt.Println("getUserBalances err:", err)
		return nil, err
	}

	userBalances := make([]*big.Int, 0)
	err = struct_util.Map2Struct(resp["userBalances"], &userBalances)
	return userBalances, err
}

func (mc *MultiCallRepo) BalanceOf(tokenAddress string, userAddress string) (*big.Int, error) {
	tokenAddr := web3.HexToAddress(tokenAddress)
	userAddr := web3.HexToAddress(userAddress)
	resp, err := mc.Contract.Call("getUserBalance", web3.Latest, tokenAddr, userAddr)
	if err != nil {
		fmt.Println("getUserBalance err:", err)
		return nil, err
	}

	userBalance := new(big.Int)
	err = struct_util.Map2Struct(resp["userBalance"], userBalance)
	return userBalance, err
}

func (mc *MultiCallRepo) GetBnbPrice() (float64, error) {
	resp, err := mc.Contract.Call("getBnbPrice", web3.Latest, bsc.WbnbBusdPair.Address)
	if err != nil {
		fmt.Println("getBnbPrice err:", err)
		return 0, err
	}

	bnbPrice := resp["bnbPriceBusd"].(*big.Int).Int64()
	floatPrice := float64(bnbPrice) / math.Pow(10, 12)
	return floatPrice, err
}

func (mc *MultiCallRepo) GetPairTokenInfo(pairAddress string) ([]*Token, error) {
	addr := web3.HexToAddress(pairAddress)
	resp, err := mc.Contract.Call("getPairTokenInfo", web3.Latest, addr)
	if err != nil {
		fmt.Println("getPairTokenInfo err:", err)
		return nil, err
	}

	tokens := make([]*Token, 0)
	err = struct_util.Map2Struct(resp["tokenInfo"], &tokens)
	return tokens, err
}

func (mc *MultiCallRepo) GetPairInfoWithPrice(pairAddress string) (*Pair, error) {
	pairAddr := web3.HexToAddress(pairAddress)
	resp, err := mc.Contract.Call("getPairInfoWithPrice", web3.Latest, pairAddr, bsc.WbnbBusdPair.Address)
	if err != nil {
		err = fmt.Errorf("getPairInfoWithPrice: %w", err)
		return nil, err
	}

	//bs, err := json.Marshal(resp)
	//if err != nil {
	//    err = fmt.Errorf("json.Marshal: %w", err)
	//    return nil, err
	//}
	//fmt.Println("resp:", string(bs))

	pair := Pair{}
	err = struct_util.Map2Struct(resp["pairInfo"], &pair)
	if err != nil {
	    err = fmt.Errorf("struct_util.Map2Struct: %w", err)
	    return nil, err
	}

	return &pair, err
}

func (mc *MultiCallRepo) GetPairInfoWithFarmTVL(pairAddress string, farmAddress string) (*Pair, error) {
	pairAddr := web3.HexToAddress(pairAddress)
	farmAddr := web3.HexToAddress(farmAddress)
	resp, err := mc.Contract.Call("getPairInfoWithFarmTVL", web3.Latest, pairAddr, farmAddr, bsc.WbnbBusdPair.Address)
	if err != nil {
		return nil, err
	}

	pair := Pair{}
	err = struct_util.Map2Struct(resp["pairInfo"], &pair)
	return &pair, err
}
