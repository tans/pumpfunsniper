package token

import (
	"context"
	"encoding/hex"
	"github.com/blocto/solana-go-sdk/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestPumpFunSell(t *testing.T) {
	endpoint := "https://red-radial-morning.solana-mainnet.quiknode.pro/b17ab2e42c9b879e94267c9e4576f396ff0afdc6"
	tokenMint := "CQnhr6X3B2BW3ez5aWWcr1B2acHFELicQ3Cc7QC1z6Dc"
	bc, _ := GetPumpBondingCurveDataIfPoolExists(tokenMint)
	privateKey := "4wQukByMneN5f4cm6grEWScWEtFb4tjgwdWE96fvvbFczBXqvCQNcEYvk9UrJTibbauNSqUyitDRgj6irPYj69zu"
	feePayer, err := types.AccountFromBase58(privateKey)
	feePayerPub, err := solana.PublicKeyFromBase58("FHzYhzArwZb5WqsHKEbC2YBPkx9TyuoQHDPUynwVyFem")
	assert.Nil(t, err)
	tokenMintPub := solana.MustPublicKeyFromBase58(tokenMint)
	amount := big.NewInt(1200000)
	slippage := float32(5)
	unsignedTx, err := GetPumpSellTx(
		feePayer.PublicKey.String(),
		tokenMint,
		amount,
		bc,
		slippage,
		100000,
		100000,
		100000,
		false,
	)
	assert.Nil(t, err)
	unsignedTxBytes, err := unsignedTx.Message.MarshalBinary()
	assert.Nil(t, err)
	unsignedTxHex := hex.EncodeToString(unsignedTxBytes)
	t.Log(unsignedTxHex)

	// client
	unsignedTxBytes, err = hex.DecodeString(unsignedTxHex)
	assert.Nil(t, err)
	clientTx := &solana.Transaction{Signatures: make([]solana.Signature, 0)}
	assert.Nil(t, err)
	signature := feePayer.Sign(unsignedTxBytes)
	t.Log(base58.Encode(signature))
	clientTx.Signatures = append(clientTx.Signatures, solana.SignatureFromBytes(signature))
	clientMsg := solana.Message{}
	err = clientMsg.UnmarshalWithDecoder(bin.NewBinDecoder(unsignedTxBytes))
	assert.Nil(t, err)
	clientTx.Message = clientMsg

	txs, err := rpc.New(endpoint).SendTransaction(context.Background(), clientTx)
	assert.Nil(t, err)

	w1, w2 := GetTokenBalance(feePayerPub, tokenMintPub)
	t.Log(txs)
	t.Log(w1)
	t.Log(w2)
}

func TestPumpFunBuy(t *testing.T) {
	endpoint := "https://red-radial-morning.solana-mainnet.quiknode.pro/b17ab2e42c9b879e94267c9e4576f396ff0afdc6"
	tokenMint := "CQnhr6X3B2BW3ez5aWWcr1B2acHFELicQ3Cc7QC1z6Dc"
	bc, _ := GetPumpBondingCurveDataIfPoolExists(tokenMint)
	privateKey := "4wQukByMneN5f4cm6grEWScWEtFb4tjgwdWE96fvvbFczBXqvCQNcEYvk9UrJTibbauNSqUyitDRgj6irPYj69zu"
	feePayer, err := types.AccountFromBase58(privateKey)
	tokenMintPub := solana.MustPublicKeyFromBase58(tokenMint)
	feePayerPub, err := solana.PublicKeyFromBase58("FHzYhzArwZb5WqsHKEbC2YBPkx9TyuoQHDPUynwVyFem")
	assert.Nil(t, err)
	unsignedTx, err := GetPumpBuyTx(
		feePayer.PublicKey.String(),
		tokenMint,
		big.NewInt(100),
		bc,
		90,
		590000,
		100000,
		0,
	)
	assert.Nil(t, err)
	unsignedTxBytes, err := unsignedTx.Message.MarshalBinary()
	assert.Nil(t, err)
	unsignedTxHex := hex.EncodeToString(unsignedTxBytes)
	t.Log(unsignedTxHex)

	// client
	unsignedTxBytes, err = hex.DecodeString(unsignedTxHex)
	assert.Nil(t, err)
	clientTx := &solana.Transaction{Signatures: make([]solana.Signature, 0)}
	assert.Nil(t, err)
	signature := feePayer.Sign(unsignedTxBytes)
	t.Log(base58.Encode(signature))
	clientTx.Signatures = append(clientTx.Signatures, solana.SignatureFromBytes(signature))
	clientMsg := solana.Message{}
	err = clientMsg.UnmarshalWithDecoder(bin.NewBinDecoder(unsignedTxBytes))
	assert.Nil(t, err)
	clientTx.Message = clientMsg

	txs, err := rpc.New(endpoint).SendTransaction(context.Background(), clientTx)
	assert.Nil(t, err)
	t.Log(txs)
	w1, w2 := GetTokenBalance(feePayerPub, tokenMintPub)
	t.Log(txs)
	t.Log(w1)
	t.Log(w2)
}

//5SzMBGPZQf3mC2QexdxheqksaGvoLAf8kKjAgXt42EjytVFsy6GWZy2gR83u2G7wQQCbSV8uW7gRDZGpL1j2v7U1
//5RcWuAoVDuAqNUQPeJmiaT8742thfQgnBivr6SJLLAeGHZsrg98giuipnPv8ozkaCYqkqkx5Ap8wmkvYoEuyWd1m
//zHBL6h4ocoqiADRSMKRtR9DNu4MqaryYW5MbS53QPNitteTUNAg8EhjEiPqJcjTczTqe5ZJhFSGsDfr3BcEPwu6
//cixNwWm6X9EeaBT3WqmvAZeVoCXPumMK7X1MqrmHitMiBCRZmoPaBV7p6daTWxgcGYAovBLMbupyPxRD4XGCfaH
//po1FxMSK3DPXJW25m5sV3t858WYv1sYp3sfkmkS6tYJfqqqUFsd6FxPqxofPf3VceMUZZHgm8PcJyEWRpkiGCgW
//3sHZmku8W7341eShfzTNJBQbDtVnzFuKfGeEm7Ztp2oCeicNcNaQNdXRTNYUffedsYZqXbDXtoPKKcLBqaaLVMDA
//5yfnThLKU8p2yHY8HD5oFfmCdP9SAnR9gAoPryZi5hKNS2z6QYjcp77Ua1ajK6pgNVDT3XLZWp7WCaABdBV9ShXw
//5kbK574Wk3qobmg23nGcR9vUCSEV5b3YUgowbspAarx6CeeXHDzRXeCo3ptMc6y1PBuTeN9EyhbBr8xP5TGaZVmi
//3ZTdH2XibioRP7QutKaVfGDSgir9R3J2Wjbu57xwjKJnjP5nJ7BAYscYh6zuCP1TeUhcBLJBVqJ97JZ8DJy7ueJy
//mxNQ2CpNGNrzAfnk9bwMRCSZZxtJCKJQbDC3QUNwUAAMhTvqaRmJ1qSjkeomgAfPaGZFoBs63uPqGyyosQ9nwgn
//269BHFF9cVvAmk5yodE4J3bR6imErzveEkUR96JCprWdhFzMMfYwX5zNqsEHM4vJV3YVnMkEMg992vAjy5W7gFGG
//48kAbBCcJFRKpPrQSQrCfBLbzYSvggArnTxQu8kvfd3GS63URNLFsFCVS8ToNjdYZhe3Bo9VDumXvfZ3enH8K5oH
