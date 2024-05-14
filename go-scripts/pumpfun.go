package token

import (
	"PumpBot/global"
	"PumpBot/utils"

	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associated_token_account "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	token_program "github.com/gagliardetto/solana-go/programs/token"
)

var (
	PUMPDex     = "pump fun (bonding curve)"
	PUMPManager = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	// pump deployed tokens are 6 decimals by default, so we can use 1*10^6 to quote the value of a single token in SOL
	PUMPQuoteSellAmountIn = big.NewInt(1000000)
	// PUMPBuyMethod         uint64 = 7351630589278743530
	// PUMPSellMethod        uint64 = 3739823480024040365
	PUMPBuyMethod  = []byte{0x66, 0x06, 0x3d, 0x12, 0x01, 0xda, 0xeb, 0xea}
	PUMPSellMethod = []byte{0x33, 0xe6, 0x85, 0xa4, 0x01, 0x7f, 0x83, 0xad}

	SlippageAdjustment int64 = 2
	// 3% slippage for exact in
	DefaultSlippage = float32(3.0)

	// will be cached once it was pulled
	// GlobalAddress *solana.PublicKey     = nil
	// Global        *GlobalSettingsLayout = nil
)

type PUMPBondingCurveData struct {
	BondingCurve             *BondingCurveLayout
	BondingCurvePk           solana.PublicKey
	AssociatedBondingCurvePk solana.PublicKey
	GlobalSettings           *GlobalSettingsLayout
	GlobalSettingsPk         solana.PublicKey
	MintAuthority            solana.PublicKey
}

type BondingCurveLayout struct {
	Blob1                uint64
	VirtualTokenReserves uint64
	VirtualSOLReserves   uint64
	RealTokenReserves    uint64
	RealSOLReserves      uint64
	BLOB4                uint64
	Complete             bool
}

type GlobalSettingsLayout struct {
	Blob1                       [8]byte
	Initialized                 bool
	Authority                   solana.PublicKey
	FeeRecipient                solana.PublicKey
	InitialVirtualTokenReserves uint64
	InitialVirtualSOLReserves   uint64
	InitialRealTokenReserves    uint64
	TokenTotalSupply            uint64
	FeeBasisPoints              uint64
}

type PumpBuyInstruction struct {
	bin.BaseVariant
	MethodId                []byte
	AmountOut               uint64
	MaxAmountIn             uint64
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (inst *PumpBuyInstruction) ProgramID() solana.PublicKey {
	return PUMPManager
}

func (inst *PumpBuyInstruction) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *PumpBuyInstruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *PumpBuyInstruction) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Swap instruction is number 9
	err = encoder.WriteBytes(inst.MethodId, false)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.AmountOut, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.MaxAmountIn, binary.LittleEndian)
	if err != nil {
		return err
	}
	return nil
}

type PumpSellInstruction struct {
	bin.BaseVariant
	MethodId                []byte
	AmountIn                uint64
	AmountOutMin            uint64
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (inst *PumpSellInstruction) ProgramID() solana.PublicKey {
	return PUMPManager
}

func (inst *PumpSellInstruction) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *PumpSellInstruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *PumpSellInstruction) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Swap instruction is number 9
	err = encoder.WriteBytes(inst.MethodId, false)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.AmountIn, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.AmountOutMin, binary.LittleEndian)
	if err != nil {
		return err
	}
	return nil
}

func GetPumpBondingCurveDataIfPoolExists(token string) (*PUMPBondingCurveData, error) {
	tokenMint := solana.MustPublicKeyFromBase58(token)

	rpc := global.GetRPCForRequest()

	bondingCurve, _, err := solana.FindProgramAddress([][]byte{
		[]byte("bonding-curve"),
		tokenMint.Bytes(),
	}, PUMPManager)
	if err != nil {
		return nil, nil
	}

	globalSettings, _, err := solana.FindProgramAddress([][]byte{
		[]byte("global"),
	}, PUMPManager)
	if err != nil {
		return nil, err
	}

	mintAuthority, _, err := solana.FindProgramAddress([][]byte{
		[]byte("mint-authority"),
	}, PUMPManager)
	if err != nil {
		return nil, err
	}

	accountInfos, err := rpc.GetMultipleAccounts(context.Background(), bondingCurve, globalSettings)
	if err != nil {
		return nil, err
	}

	if accountInfos == nil {
		return nil, errors.New("pumpfun: accounts not found")
	}

	if accountInfos.Value[0] == nil || (accountInfos.Value[0] != nil && len(accountInfos.Value[0].Data.GetBinary()) == 0) {
		return nil, errors.New("pumpfun: bonding curve for mint not found")
	}

	if accountInfos.Value[1] == nil || (accountInfos.Value[1] != nil && len(accountInfos.Value[1].Data.GetBinary()) == 0) {
		return nil, errors.New("pumpfun: global settings not found")
	}

	// if bonding curve exists, global has to exists too, we still check
	var bondingCurveLayout BondingCurveLayout
	var globalSettingsLayout GlobalSettingsLayout

	err = decode(accountInfos.Value[0].Data.GetBinary(), &bondingCurveLayout)
	if err != nil {
		return nil, err
	}

	err = decode(accountInfos.Value[1].Data.GetBinary(), &globalSettingsLayout)
	if err != nil {
		return nil, err
	}

	associatedBondingCurve, _, _ := solana.FindAssociatedTokenAddress(bondingCurve, tokenMint)
	return &PUMPBondingCurveData{
		BondingCurve:             &bondingCurveLayout,
		BondingCurvePk:           bondingCurve,
		AssociatedBondingCurvePk: associatedBondingCurve,
		GlobalSettings:           &globalSettingsLayout,
		GlobalSettingsPk:         globalSettings,
		MintAuthority:            mintAuthority,
	}, nil
}

func decode(binary []byte, v interface{}) error {
	borsh := bin.NewBorshDecoder(binary)
	return borsh.Decode(&v)
}

// assumes that bondingCurveData is not nil when called
func getPriceAndLiquidityAndDexFromPump(bondingCurveData *PUMPBondingCurveData) (price *big.Float, liquidity string, dex string) {
	dex = PUMPDex
	// virtual sol * 2 is the virtual liquidity
	liquidity = new(big.Float).Quo(big.NewFloat(float64(bondingCurveData.BondingCurve.VirtualSOLReserves*2)), global.FloatConst[9]).String()
	pricePerTokenInSOLRaw := pumpQuoteSell(PUMPQuoteSellAmountIn, bondingCurveData)
	priceFloat64, _ := pricePerTokenInSOLRaw.Float64()
	pricePerTokenInSOL := new(big.Float).Quo(big.NewFloat(priceFloat64), global.FloatConst[9])
	return pricePerTokenInSOL, liquidity, dex
}

// for a given amountIn, quotes how much sol they are worth
func pumpQuoteSell(amountIn *big.Int, bondingCurveData *PUMPBondingCurveData) *big.Int {
	newReserves := new(big.Int).Add(big.NewInt(int64(bondingCurveData.BondingCurve.VirtualTokenReserves)), amountIn)
	temp := new(big.Int).Mul(amountIn, big.NewInt(int64(bondingCurveData.BondingCurve.VirtualSOLReserves)))
	amountOut := new(big.Int).Div(temp, newReserves)
	fee := pumpGetFee(amountOut, bondingCurveData.GlobalSettings.FeeBasisPoints)
	amountOutAfterFee := new(big.Int).Sub(amountOut, fee)
	return amountOutAfterFee
}

// for a given amountIn, quotes how many tokens can be bought
func pumpQuoteBuy(amountIn *big.Int, bondingCurveData *PUMPBondingCurveData) *big.Int {
	virtualSOLReservesBN := big.NewInt(int64(bondingCurveData.BondingCurve.VirtualSOLReserves))
	virtualTokenReservesBN := big.NewInt(int64(bondingCurveData.BondingCurve.VirtualTokenReserves))

	reservesProduct := new(big.Int).Mul(virtualSOLReservesBN, virtualTokenReservesBN)
	newVirtualSOLReserve := new(big.Int).Add(virtualSOLReservesBN, amountIn)
	newVirtualTokenReserve := new(big.Int).Div(reservesProduct, newVirtualSOLReserve)
	newVirtualTokenReserve = new(big.Int).Add(newVirtualTokenReserve, big.NewInt(1))
	amountOut := new(big.Int).Sub(virtualTokenReservesBN, newVirtualTokenReserve)
	finalAmountOut := amountOut
	if amountOut.Uint64() > bondingCurveData.BondingCurve.RealTokenReserves {
		finalAmountOut = big.NewInt(int64(bondingCurveData.BondingCurve.RealTokenReserves))
	}
	return finalAmountOut
}

func pumpGetFee(amount *big.Int, feeBP uint64) *big.Int {
	temp := new(big.Int).Mul(amount, big.NewInt(int64(feeBP)))
	feeAmount := new(big.Int).Div(temp, global.Big10000)
	return feeAmount
}

func GetPumpBuyTx(
	signerAndOwner *solana.PrivateKey,
	mint *solana.PublicKey,
	// maxIn without taking any fees
	maxAmountIn *big.Int,
	bondingCurveData *PUMPBondingCurveData,
	// slippage% 0-100
	slippage float32,
	priorityFee uint64,
	fee uint64,
	jitoTip uint64,
) (*solana.Transaction, error) {

	slippage = float32(DefaultSlippage)

	instrs := []solana.Instruction{}
	signers := []solana.PrivateKey{*signerAndOwner}

	amountInAfterOurFee := new(big.Int).Sub(maxAmountIn, big.NewInt(int64(fee)))

	if jitoTip > 0 {
		instrs = append(instrs, system.NewTransferInstruction(jitoTip, signerAndOwner.PublicKey(), global.PickRandomTip()).Build())
	}

	instrs = append(instrs, computebudget.NewSetComputeUnitLimitInstruction(100514).Build())

	if priorityFee > 0 {
		instrs = append(instrs, computebudget.NewSetComputeUnitPriceInstruction(priorityFee).Build())
	}

	if fee > 0 {
		instrs = append(instrs, system.NewTransferInstruction(fee, signerAndOwner.PublicKey(), global.FeeAccountBuys).Build())
	}

	// createSOLAccountOrWrap(&instrs, signerAndOwner.PublicKey(), amountInAfterOurFee)
	createTokenAccountIfNotExists(&instrs, signerAndOwner.PublicKey(), mint)

	addPumpBuyIx(&instrs, signerAndOwner.PublicKey(), mint, amountInAfterOurFee, bondingCurveData, slippage)

	tx, err := BuildTransaction(signers, *signerAndOwner, instrs...)
	return tx, err
}

func BuildTransaction(signers []solana.PrivateKey, signer solana.PrivateKey, instrs ...solana.Instruction) (*solana.Transaction, error) {
	tx, err := solana.NewTransaction(
		instrs,
		utils.GetBlockHash(),
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			return &signer
		},
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func addPumpBuyIx(
	instrs *[]solana.Instruction,
	owner solana.PublicKey,
	mint *solana.PublicKey,
	amountInAfterOurFee *big.Int,
	bondingCurveData *PUMPBondingCurveData,
	// slippage% 0-100
	slippage float32,
) {

	pumpFee := pumpGetFee(amountInAfterOurFee, bondingCurveData.GlobalSettings.FeeBasisPoints)

	// this is used for quoting
	amountInAfterPumpFee := new(big.Int).Sub(amountInAfterOurFee, pumpFee)
	amountOut := pumpQuoteBuy(amountInAfterPumpFee, bondingCurveData)
	amountOutWithSlippage := applySlippage(amountOut, slippage)

	// we apply slippage on the amountOut
	// quote buy here, then apply slippage
	// if slippage is 100%, we reduce it

	instruction := &PumpBuyInstruction{
		MethodId:         PUMPBuyMethod,
		MaxAmountIn:      amountInAfterOurFee.Uint64(),
		AmountOut:        amountOutWithSlippage.Uint64(),
		AccountMetaSlice: make(solana.AccountMetaSlice, 10),
	}

	instruction.BaseVariant = bin.BaseVariant{
		Impl: instruction,
	}

	ataUser, _, _ := solana.FindAssociatedTokenAddress(owner, *mint)

	instruction.AccountMetaSlice[0] = solana.Meta(bondingCurveData.GlobalSettingsPk)
	instruction.AccountMetaSlice[1] = solana.Meta(bondingCurveData.GlobalSettings.FeeRecipient).WRITE()
	instruction.AccountMetaSlice[2] = solana.Meta(*mint)
	instruction.AccountMetaSlice[3] = solana.Meta(bondingCurveData.BondingCurvePk).WRITE()
	instruction.AccountMetaSlice[4] = solana.Meta(bondingCurveData.AssociatedBondingCurvePk).WRITE()
	instruction.AccountMetaSlice[5] = solana.Meta(ataUser).WRITE()
	instruction.AccountMetaSlice[6] = solana.Meta(owner).WRITE().SIGNER()
	instruction.AccountMetaSlice[7] = solana.Meta(solana.SystemProgramID)
	instruction.AccountMetaSlice[8] = solana.Meta(solana.TokenProgramID)
	instruction.AccountMetaSlice[9] = solana.Meta(solana.SysVarRentPubkey)

	*instrs = append(*instrs, instruction)
}

func GetTokenBalance(account solana.PublicKey, token solana.PublicKey) (wDecimals *big.Int, woDecimals *big.Float) {
	tokenAccount, _, err := solana.FindAssociatedTokenAddress(account, token)
	if err != nil {
		fmt.Println("Failed to find token account", err)
		return nil, nil
	}
	ctx, exp := context.WithTimeout(context.Background(), 3*time.Second)
	defer exp()
	balance, err := global.GetRPCForRequest().GetTokenAccountBalance(ctx, tokenAccount, "confirmed")
	if err != nil {
		//if !strings.Contains(err.Error(), "could not find account") {
		//fmt.Println("Failed to query token balance", err)
		//}
		return nil, nil
	}
	if balance != nil && balance.Value != nil {
		wDecimals, _ = new(big.Int).SetString(balance.Value.Amount, 10)
		woDecimals, _ = new(big.Float).SetString(balance.Value.Amount)
	}
	return wDecimals, woDecimals
}

func createSOLAccountOrWrap(instrs *[]solana.Instruction, owner solana.PublicKey, amountIn *big.Int) {
	solATA, _, _ := solana.FindAssociatedTokenAddress(owner, solana.WrappedSol)
	bal, _ := GetTokenBalance(owner, solana.WrappedSol)

	var wrapAmountNeeded uint64
	// if ata exists, we just send sol and sync it, otherwise we wrap the difference from amountIn and the current balance and sync it
	if bal == nil || (bal != nil && bal.Uint64() == 0) {
		*instrs = append(*instrs, associated_token_account.NewCreateInstruction(owner, owner, solana.WrappedSol).Build())
		wrapAmountNeeded = amountIn.Uint64()
	} else {
		targetBalance := amountIn.Uint64()
		wsolAccountBalanceU64 := bal.Uint64()
		if wsolAccountBalanceU64 < targetBalance {
			wrapAmountNeeded = targetBalance - wsolAccountBalanceU64
		}
	}

	if wrapAmountNeeded > 0 {
		*instrs = append(*instrs, system.NewTransferInstruction(wrapAmountNeeded, owner, solATA).Build())
		*instrs = append(*instrs, token_program.NewSyncNativeInstruction(solATA).Build())
	}
}

func createTokenAccountIfNotExists(instrs *[]solana.Instruction, owner solana.PublicKey, mint *solana.PublicKey) {
	bal, _ := GetTokenBalance(owner, *mint)
	if bal == nil || (bal != nil && bal.Uint64() == 0) {
		*instrs = append(*instrs, associated_token_account.NewCreateInstruction(owner, owner, *mint).Build())
	}
}

// slippage is a value between 0 - 100
func applySlippage(amount *big.Int, slippage float32) *big.Int {

	slippageBP := (int64(100*slippage) + 25) * SlippageAdjustment
	maxSlippage := new(big.Int).Mul(global.Big10000, big.NewInt(SlippageAdjustment))

	if slippageBP > maxSlippage.Int64() {
		slippageBP = global.Big10000.Int64()
	}

	slippageBPBN := big.NewInt(slippageBP)

	// we adjust slippage so that it caps out at 50%
	slippageNumeratorMul := new(big.Int).Sub(maxSlippage, slippageBPBN)
	slippageNumerator := new(big.Int).Mul(amount, slippageNumeratorMul)
	amountWithSlippage := new(big.Int).Div(slippageNumerator, maxSlippage)
	return amountWithSlippage
}

func PumpGetValueAndPriceImpact(token string, amount *big.Int, decimals uint8) (value float64, priceImpact float64) {

	curve, err := GetPumpBondingCurveDataIfPoolExists(token)
	if err != nil {
		return 0, 0
	}

	// theoretically the lower amount has to be scaled to the other bigger amount
	// but 1 token is usually the lower amount, hence we scale amount down
	amountFloat := big.NewFloat(float64(amount.Uint64()))
	singleTokenFloat := big.NewFloat(float64(PUMPQuoteSellAmountIn.Uint64()))
	valueMulti := new(big.Float).Quo(amountFloat, singleTokenFloat)
	singleTokenPrice := pumpQuoteSell(PUMPQuoteSellAmountIn, curve)

	// value of the tokens
	sellQuote := pumpQuoteSell(amount, curve)

	adj := utils.ReduceDecimals(sellQuote, 9)

	// we have to either scale this up to find out the worth for the same amount
	// or scale the sellquote down to single token and compare those
	singleAdj := utils.ReduceDecimals(singleTokenPrice, 9)
	singleAdjWithMulti := new(big.Float).Mul(valueMulti, singleAdj)

	singleAdjWithMultiFloat, _ := singleAdjWithMulti.Float64()

	value, _ = adj.Float64()
	priceImpact = singleAdjWithMultiFloat / value

	return value, priceImpact
}

func GetPumpSellTx(
	signerAndOwner *solana.PrivateKey,
	mint *solana.PublicKey,
	// maxIn without taking any fees
	amountIn *big.Int,
	bondingCurveData *PUMPBondingCurveData,
	// slippage% 0-100
	slippage float32,
	priorityFee uint64,
	fee uint64,
	jitoTip uint64,
	shouldCloseTokenInAccount bool,
) (*solana.Transaction, error) {
	instrs := []solana.Instruction{}
	signers := []solana.PrivateKey{*signerAndOwner}

	if jitoTip > 0 {
		instrs = append(instrs, system.NewTransferInstruction(jitoTip, signerAndOwner.PublicKey(), global.PickRandomTip()).Build())
	}

	instrs = append(instrs, computebudget.NewSetComputeUnitLimitInstruction(100514).Build())

	if priorityFee > 0 {
		instrs = append(instrs, computebudget.NewSetComputeUnitPriceInstruction(priorityFee).Build())
	}

	if fee > 0 {
		instrs = append(instrs, system.NewTransferInstruction(fee, signerAndOwner.PublicKey(), global.FeeAccountBuys).Build())
	}

	// createSOLAccountOrWrap(&instrs, signerAndOwner.PublicKey(), big.NewInt(0))
	addPumpSellIx(&instrs, signerAndOwner.PublicKey(), mint, amountIn, bondingCurveData, slippage)
	// closeATA(&instrs, signerAndOwner.PublicKey(), solana.WrappedSol)
	if shouldCloseTokenInAccount {
		closeATA(&instrs, signerAndOwner.PublicKey(), *mint)
	}

	tx, err := BuildTransaction(signers, *signerAndOwner, instrs...)
	return tx, err
}

func addPumpSellIx(
	instrs *[]solana.Instruction,
	owner solana.PublicKey,
	mint *solana.PublicKey,
	amountIn *big.Int,
	bondingCurveData *PUMPBondingCurveData,
	// slippage% 0-100
	slippage float32,
) {
	amountOut := pumpQuoteSell(amountIn, bondingCurveData)
	amountOutWithSlippage := applySlippage(amountOut, slippage)
	// we apply slippage on the amountOut
	// quote buy here, then apply slippage
	// if slippage is 100%, we reduce it

	instruction := &PumpSellInstruction{
		MethodId:         PUMPSellMethod,
		AmountIn:         amountIn.Uint64(),
		AmountOutMin:     amountOutWithSlippage.Uint64(),
		AccountMetaSlice: make(solana.AccountMetaSlice, 10),
	}

	instruction.BaseVariant = bin.BaseVariant{
		Impl: instruction,
	}

	ataUser, _, _ := solana.FindAssociatedTokenAddress(owner, *mint)

	instruction.AccountMetaSlice[0] = solana.Meta(bondingCurveData.GlobalSettingsPk)
	instruction.AccountMetaSlice[1] = solana.Meta(bondingCurveData.GlobalSettings.FeeRecipient).WRITE()
	instruction.AccountMetaSlice[2] = solana.Meta(*mint)
	instruction.AccountMetaSlice[3] = solana.Meta(bondingCurveData.BondingCurvePk).WRITE()
	instruction.AccountMetaSlice[4] = solana.Meta(bondingCurveData.AssociatedBondingCurvePk).WRITE()
	instruction.AccountMetaSlice[5] = solana.Meta(ataUser).WRITE()
	instruction.AccountMetaSlice[6] = solana.Meta(owner).WRITE().SIGNER()
	instruction.AccountMetaSlice[7] = solana.Meta(solana.SystemProgramID)
	instruction.AccountMetaSlice[8] = solana.Meta(solana.SPLAssociatedTokenAccountProgramID)
	instruction.AccountMetaSlice[9] = solana.Meta(solana.TokenProgramID)

	*instrs = append(*instrs, instruction)
}

func closeATA(instrs *[]solana.Instruction, owner solana.PublicKey, mint solana.PublicKey) {
	ata, _, _ := solana.FindAssociatedTokenAddress(owner, mint)
	closeInst := token_program.NewCloseAccountInstruction(
		ata,
		owner,
		owner,
		[]solana.PublicKey{},
	).Build()
	*instrs = append(*instrs, closeInst)
}
