package utils

import (
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/valyala/fasthttp"
)

const (
	tolerance       = 0
	stableTolerance = 0
)

var (
	priceBlock       uint64
	currentPrice     float64
	currentPriceF    *big.Float
	blockMutex       = &sync.RWMutex{}
	priceMutex       = &sync.RWMutex{}
	tokenpriceMutex  = &sync.RWMutex{}
	tokenprice       = make(map[string]float64)
	tokenpriceupdate = make(map[string]uint64)
)

type PriceQuoteTokenData struct {
	ID            string  `json:"id"`
	MintSymbol    string  `json:"mintSymbol"`
	VsToken       string  `json:"vsToken"`
	VsTokenSymbol string  `json:"vsTokenSymbol"`
	Price         float64 `json:"price"`
}

type PriceQuote struct {
	Data      map[string]PriceQuoteTokenData `json:"data"`
	TimeTaken float64                        `json:"timeTaken"`
}

type CryptoPrices struct {
	Dai      CryptoPrice `json:"dai"`
	Ethereum CryptoPrice `json:"ethereum"`
	Tether   CryptoPrice `json:"tether"`
	UsdCoin  CryptoPrice `json:"usd-coin"`
	Solana   CryptoPrice `json:"solana"`
}

type CryptoPrice struct {
	Usd float64 `json:"usd"`
}

func getCurrentBlockFromAPI() uint64 {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("")

	err := fasthttp.Do(req, resp)
	if err != nil {
		//Log.Printf("Error fetching block: %v", err)
		return 0
	}

	var block uint64
	fmt.Sscanf(string(resp.Body()), "%d", &block)
	return block
}

func getPriceFromAPI() float64 {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("")

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Printf("Error fetching price: %v", err)
		return 0
	}

	var price float64
	fmt.Sscanf(string(resp.Body()), "%f", &price)
	return price
}

func GetCurrentPrice() float64 {
	priceMutex.RLock()
	defer priceMutex.RUnlock()
	return currentPrice
}

// func GetTokenPriceByString(token string) float64 {
// 	if token == "" {
// 		fmt.Println("Invalid Active Token to Fetch Price.")
// 		return 0
// 	}
// 	if price, exists := tokenprice[token]; exists && tokenpriceupdate[token] >= GetSlot() {
// 		return price
// 	}
// 	return UpdateTokenPrice(token)
// }

// func UpdateTokenPrices() {
// 	for token := range tokenprice {
// 		go UpdateTokenPrice(token)
// 	}
// }

// func UpdateTokenPrice(token string) float64 {

// 	_, body, err := fasthttp.GetTimeout(nil, "https://price.jup.ag/v6/price?ids="+token, 2*time.Second)

// 	if err != nil {
// 		fmt.Println("Token price update failed for ", token)
// 	}
// 	quote := PriceQuote{}
// 	err = json.Unmarshal(body, &quote)
// 	if err != nil {
// 		fmt.Println("Token price update failed for ", token)
// 		return -1
// 	}
// 	if quote.TimeTaken == 0 {
// 		return -1
// 	}
// 	tokenpriceMutex.Lock()
// 	tokenprice[token] = quote.Data[token].Price
// 	tokenpriceupdate[token] = GetSlot()
// 	tokenpriceMutex.Unlock()
// 	return quote.Data[token].Price
// }

func GetCurrentPriceF() *big.Float {
	priceMutex.RLock()
	defer priceMutex.RUnlock()
	return currentPriceF
}