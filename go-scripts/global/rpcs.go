package global

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"net/url"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

var (
	rpcMutex = &sync.RWMutex{}
)

func ConnectToEndpoints() {
	GeyserRPC = rpc.New("")
	for _, server := range RPCs {
		rpcClient := rpc.New(server)
		if rpcClient == nil {
			log.Println("RPC Connect Client Nil")
		}

		wsServer := ConvertToWSURL(server)
		len := len(wsServer) - 4
		if wsServer[len:] == "8899" {
			wsServer = wsServer[:len] + "8900"
		}
		wsClient, err := ws.Connect(context.Background(), wsServer)
		if err != nil {
			log.Println("RPC WEBSOCKET Connect Error:", err)
		}
		RPCServers = append(RPCServers, rpcClient)
		WSServers = append(WSServers, wsClient)
	}
	for _, server := range JitoServers {
		rpcClient := rpc.New(server)
		if rpcClient == nil {
			log.Fatalln("Broadcast RPC Connect Client Nil")
		}
		JitoRPCs = append(JitoRPCs, rpcClient)
	}
}

func ConvertToWSURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	finalURL := u.String()
	return finalURL
}

//
//func BuildJitoClient() {
//	c, err := client.NewGojitoClient(
//		"https://ny.mainnet.block-engine.jito.wtf",
//		"",
//		"./Jito_Mev_Key.json",
//		false,
//	)
//	if err != nil {
//		log.Println("Jito Connect Err:", err)
//	} else {
//		JitoClient = c
//	}
//}

func GetRPCForRequest() *rpc.Client {
	//rpcMutex.Lock()
	//defer rpcMutex.Unlock()
	//last := RPCLast.Load()
	//if last < RPCLen {
	//	RPCLast.Store(last + 1)
	//	last += 1
	//} else {
	//	RPCLast.Store(0)
	//}
	//return RPCServers[last]
	endpoint := "https://red-radial-morning.solana-mainnet.quiknode.pro/b17ab2e42c9b879e94267c9e4576f396ff0afdc6"
	return rpc.New(endpoint)
}

func GetWSRPCForRequest() (*rpc.Client, *ws.Client) {
	rpcMutex.Lock()
	defer rpcMutex.Unlock()
	last := RPCLast.Load()
	if last < RPCLen {
		RPCLast.Store(last + 1)
		last += 1
	} else {
		RPCLast.Store(0)
	}
	return RPCServers[last], WSServers[last]
}

func GetJitoRPCs() []*rpc.Client {
	rpcMutex.RLock()
	defer rpcMutex.RUnlock()
	return JitoRPCs
}

func CheckRPCsAndReconnect() {
	for n, server := range RPCs {
		rpcMutex.RLock()
		rpcClient := RPCServers[n]
		rpcMutex.RUnlock()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		health, err := rpcClient.GetHealth(ctx)
		if err != nil || health != "ok" {
			fmt.Println("Disconnected from RPC " + server)
			dead2, deadC2 := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
			wsClient, err := ws.Connect(dead2, "ws://"+server+":8900")
			deadC2()
			if err == nil {
				fmt.Println("Reconnected to WS " + server)
				rpcMutex.Lock()
				WSServers[n] = wsClient
				rpcMutex.Unlock()
			} else {
				fmt.Println("Failed to reconnect WS " + server)
			}
			rpcNew := rpc.New("http://" + server + ":8899")
			if rpcNew != nil {
				fmt.Println("Reconnected to HTTP " + server)
				rpcMutex.Lock()
				RPCServers[n] = rpcNew
				rpcMutex.Unlock()
			} else {
				fmt.Println("Failed to reconnect HTTP " + server)
			}
		}
		cancel()
	}
	rpcMutex.RLock()
	jitoClients := JitoRPCs
	rpcMutex.RUnlock()
	for n, c := range jitoClients {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		health, err := c.GetHealth(ctx)
		if err != nil || health != "ok" {
			fmt.Println("Disconnected from Broadcast RPC " + fmt.Sprint(n))
			rpcNew := rpc.New(JitoServers[n])
			if rpcNew != nil {
				fmt.Println("Reconnected to Broadcast RPC " + JitoServers[n])
				rpcMutex.Lock()
				jitoClients[n] = rpcNew
				rpcMutex.Unlock()
			} else {
				fmt.Println("Failed to reconnect Broadcast RPC " + JitoServers[n])
			}
		}
		cancel()
	}
}
