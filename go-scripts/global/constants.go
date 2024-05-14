package global

// Contains constants variables that set the project's status.

import (
	atomic_ "PumpBot/utils/atomic"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alphbuff/gojito/client"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/joho/godotenv"
)

type LimitType int

var (
	Retries = uint(100)
	// If true, the project runs in dev mode.
	Dev = true
	// Solana Chain NAme to use in messages
	ChainName = "Solana"
	// Gas Unit based on chain
	GasUnitString = "SOL"
	// Solana coin symbol to use in messages
	CoinSymbol = "SOL"
	// Chain explorer account checking page
	ChainExplorerAccount   = "https://solscan.io/account/"
	ChainExplorerTokenLink = "https://solscan.io/token/"
	ChainExplorerTxLink    = "https://solscan.io/tx/"
	// Chain explorer name
	ChainExplorerName = "Solscan"
	// Maximum user's wallet count
	WalletLimit = 3
	// Bit floats that used many times in the code
	// Is TG populated
	Populated       = atomic_.Bool{}
	Float0          = new(big.Float).SetUint64(0)
	Float1          = new(big.Float).SetUint64(1)
	Float100        = new(big.Float).SetUint64(100)
	Float1Lamp      = new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL)
	F1Lamp          = float64(1000000000)
	Solana          = "So11111111111111111111111111111111111111112"
	ZeroAddr        = solana.PublicKey{}
	SolanaPublic, _ = solana.PublicKeyFromBase58(Solana)
	TenThousand     = big.NewInt(10000)
	OneMillion      = big.NewInt(1e6)
	OneBillion      = big.NewInt(1e9)
	OneTrillion     = big.NewInt(1e12)
	BigFloat100     = new(big.Float).SetUint64(100)
	OneT            = new(big.Float).SetInt(OneTrillion)
	OneB            = new(big.Float).SetInt(OneBillion)
	OneM            = new(big.Float).SetInt(OneMillion)
	FloatConst      = []*big.Float{}
	JitoClient      *client.GojitoClient

	LIMIT_TRAILING = 1
	LIMIT_STOP     = 2
	LIMIT_PROFIT   = 3
	LIMITFEE       = 1

	// Solana RPC vars
	RPCLast    = atomic.Int32{}
	WSLast     = atomic.Int32{}
	RPCServers = []*rpc.Client{}
	WSServers  = []*ws.Client{}

	GeyserRPC *rpc.Client

	// Update this
	RPCs        = []string{""}
	JitoServers = []string{""}
	JitoRPCs    []*rpc.Client
	RPCLen      = int32(len(RPCs)) - 1

	FeeAccount = "7tQiiBdKoScWQkB1RmVuML7DBGnR31cuKPEtMM7Vy5SA"
	Pairs      = []string{
		"So11111111111111111111111111111111111111112",
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}
	StablesPub = []solana.PublicKey{}
	PairsPub   = []solana.PublicKey{}
	TipWallets = []solana.PublicKey{
		solana.MustPublicKeyFromBase58("96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5"),
		solana.MustPublicKeyFromBase58("HFqU5x63VTqvQss8hp11i4wVV8bD44PvwucfZ2bU7gRe"),
		solana.MustPublicKeyFromBase58("Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY"),
		solana.MustPublicKeyFromBase58("ADaUMid9yfUytqMBgopwjb2DTLSokTSzL1zt6iGPaS49"),
		solana.MustPublicKeyFromBase58("DfXygSm4jCyNCybVYYK6DwvWqjKee8pbDmJGcLWNDXjh"),
		solana.MustPublicKeyFromBase58("ADuUkR4vqLUMWXxW9gh6D6L8pMSawimctcNZ5pGwDcEt"),
		solana.MustPublicKeyFromBase58("DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL"),
		solana.MustPublicKeyFromBase58("3AVi9Tg9Uo68tJfuvoKvqKNWKkC5wPdSSdeBnizKZ6jT"),
	}
	Stables = []string{
		"Q6XprfkF8RQQKoQVG33xT88H7wi8Uk1B1CC7YAs69Gi",
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}
	FeeAccountBuys, _ = solana.PublicKeyFromBase58("4BBNEVRgrxVKv9f7pMNE788XM1tt379X9vNjpDH2KCL7")
	QuoteAmount       = big.NewInt(1000000)
	Big10000          = big.NewInt(10000)
	SwapFeeBPS = "100"

	UserCfgMutexes = make(map[int64]*sync.RWMutex)
)

// Load envrionment variables from .env
func LoadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	for i := 0; i <= 77; i++ {
		num, _ := new(big.Float).SetString(fmt.Sprintf("1e%d", i))
		FloatConst = append(FloatConst, num)
	}
	for _, pair := range Pairs {
		PairsPub = append(PairsPub, solana.MustPublicKeyFromBase58(pair))
	}
	for _, pair := range Stables {
		StablesPub = append(StablesPub, solana.MustPublicKeyFromBase58(pair))
	}

	// TGAPIKeys = strings.Split(os.Getenv("TGAPIKeys"), ",")
	// RPCs = strings.Split(os.Getenv("RPCs"), ",")

}

func PickRandomTip() solana.PublicKey {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(TipWallets))
	return TipWallets[randomIndex]
}
