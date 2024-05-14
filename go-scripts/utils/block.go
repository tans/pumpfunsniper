package utils

import (
	"PumpBot/global"
	atomic_ "PumpBot/utils/atomic"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/valyala/fasthttp"
)

var (
	BlockChan    = make(chan []byte, 4096)
	updateMutex  = &sync.Mutex{}
	currentBlock = &atomic_.Uint64{}
	block        = &CurrentBlock{}
	gasAPIMap    = gasPrio{}
	gasAPI       = ""
	sysvarclock  = solana.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
)

type gasPrio struct {
	VeryHigh uint64
	High     uint64
	Medium   uint64
	Low      uint64
	mu       sync.RWMutex
}

type CurrentBlock struct {
	PrevBlockHash solana.Hash
	BlockNum      uint64
	LastTime      int64
	Time          int64
	Slot          uint64
}

type Clock struct {
	Slot                uint64
	EpochStartTimestamp int64
	Epoch               uint64
	LeaderScheduleEpoch uint64
	UnixTimestamp       int64
}

func BlockSubscribe() {
	rpcClient, wsClient := global.GetWSRPCForRequest()
	sub, err := wsClient.SlotSubscribe()
	if err != nil {
		log.Println("Failed to subscribe to new blocks", err)
	} else {
		fmt.Println("Subbed to new blocks")
	}
	for {
		rec, err := sub.Recv()
		if err != nil {
			fmt.Println("Error block sub", err)
			break
		}
		go UpdateBlock(rpcClient, rec.Slot)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Attempting to reconnect block sub")
	go BlockSubscribe()
}

func UpdateBlock(rpcClient *rpc.Client, slot uint64) {
	updateMutex.Lock()
	defer updateMutex.Unlock()
	clock := Clock{}
	ctx, exp := context.WithTimeout(context.Background(), 3*time.Second)
	defer exp()
	if currentBlock.Load() < slot {
		recent, err := rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			log.Println("Failed grabbing new slot block hash:", err)
			return
		}
		blockMutex.Lock()
		block.BlockNum = slot
		block.Slot = slot
		block.LastTime = clock.UnixTimestamp
		block.Time = clock.UnixTimestamp
		block.PrevBlockHash = recent.Value.Blockhash
		currentBlock.Store(slot)
		blockMutex.Unlock()
		updateGas()
	}
}

func GetSlot() uint64 {
	blockMutex.RLock()
	defer blockMutex.RUnlock()
	return block.Slot
}

func updateGas() {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(gasAPI)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")

	reqBody := []byte(`{
		"jsonrpc": "2.0",
		"id": "1",
		"method": "getPriorityFeeEstimate",
		"params": [{
			"accountKeys": ["JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"],
			"options": {
				"includeAllPriorityFeeLevels": true
			}
		}]
	}`)

	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	respBody := resp.Body()
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Printf("Error decoding JSON: %s\n", err)
		return
	}

	// Assuming `result` is of type `map[string]interface{}`
	if result != nil {
		resultMap, ok := result["result"].(map[string]interface{})
		if ok {
			priorityFeeLevels, ok := resultMap["priorityFeeLevels"].(map[string]interface{})
			if ok {
				for level, fee := range priorityFeeLevels {
					switch level {
					case "low":
						gasAPIMap.mu.Lock()
						gasAPIMap.Low = uint64(fee.(float64))
						gasAPIMap.mu.Unlock()
					case "medium":
						gasAPIMap.mu.Lock()
						gasAPIMap.Medium = uint64(fee.(float64))
						gasAPIMap.mu.Unlock()
					case "high":
						gasAPIMap.mu.Lock()
						gasAPIMap.High = uint64(fee.(float64))
						gasAPIMap.mu.Unlock()
					case "veryHigh":
						gasAPIMap.mu.Lock()
						gasAPIMap.VeryHigh = uint64(fee.(float64))
						gasAPIMap.mu.Unlock()
					}
				}
			}
		}
	}
}

func GetLow() uint64 {
	gasAPIMap.mu.RLock()
	defer gasAPIMap.mu.RUnlock()
	return gasAPIMap.Low
}

func GetMedium() uint64 {
	gasAPIMap.mu.RLock()
	defer gasAPIMap.mu.RUnlock()
	return gasAPIMap.Medium
}

func GetHigh() uint64 {
	gasAPIMap.mu.RLock()
	defer gasAPIMap.mu.RUnlock()
	return gasAPIMap.High
}

func GetVeryHigh() uint64 {
	gasAPIMap.mu.RLock()
	defer gasAPIMap.mu.RUnlock()
	return gasAPIMap.VeryHigh
}
