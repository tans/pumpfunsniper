func CloseATAs() {
	wallet, _ := solana.PrivateKeyFromBase58("your private key here")
	swap := token.NewRaydiumSwap(global.GetRPCForRequest(), wallet)

	tokens := GetTokenAccountsByOwner(wallet.PublicKey().String())
	count := 0
	for _, token_addr := range tokens {
		out, err_t := global.RPCServers[0].GetTokenLargestAccounts(context.TODO(), solana.MustPublicKeyFromBase58(token_addr), rpc.CommitmentFinalized)
		if err_t != nil {
			fmt.Println("GetLargestIssue", err_t.Error())
		}
		tx, _ := swap.CloseAccount(context.TODO(), out.Value[0].Address.String(), wallet, token_addr, wallet.PublicKey())
		var (
			signature *solana.Signature
			rec       *rpc.GetTransactionResult
			err       error
		)
		fmt.Println("RemoveSent:", tx.Signatures[0].String())
		rec, signature, err = utils.SendTransactionWaitConfirmed(tx)

		if signature == nil || err != nil || rec == nil {
			fmt.Println("Failed to sent transaction")
		} else {
			fmt.Println("Closed:")
			count++
		}
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println("Success:", count)
}

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Parameter struct {
	Pubkey string `json:"pubkey"`
	Config Config `json:"config"`
}

type Config struct {
	Encoding string `json:"encoding"`
}

func GetTokenAccountsByOwner(walletAddress string) []string {
	url := "https://api.mainnet-beta.solana.com" // Change this URL to your Solana cluster endpoint

	requestBody, err := json.Marshal(RpcRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "getTokenAccountsByOwner",
		Params: []interface{}{
			walletAddress,
			map[string]interface{}{
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			},
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Result struct {
		Result struct {
			Value []struct {
				Pubkey  string `json:"pubkey"`
				Account struct {
					Data struct {
						Parsed struct {
							Info struct {
								Mint string `json:"mint"`
							} `json:"info"`
						} `json:"parsed"`
					} `json:"data"`
				} `json:"account"`
			} `json:"value"`
		} `json:"result"`
	}

	var result Result

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return []string{}
	}

	var pubkeys []string
	var tokens []string
	for _, item := range result.Result.Value {
		pubkeys = append(pubkeys, item.Pubkey)
		if item.Account.Data.Parsed.Info.Mint != global.Solana {
			tokens = append(tokens, item.Account.Data.Parsed.Info.Mint)
		}

	}
	fmt.Printf("%+v\n", tokens)
	return tokens
}

func (s *RaydiumSwap) CloseAccount(
	ctx context.Context,
	vault string,
	wallet solana.PrivateKey,
	token_addr string,
	account solana.PublicKey,
) (*solana.Transaction, error) {
	instrs := []solana.Instruction{}
	signers := []solana.PrivateKey{s.account}
	fromAccountDt, _, _ := solana.FindAssociatedTokenAddress(s.account.PublicKey(), solana.MustPublicKeyFromBase58(token_addr))

	account = fromAccountDt

	// Get the balance of the ATA
	balance, _ := utils.GetTokenBalance(wallet.PublicKey(), solana.MustPublicKeyFromBase58(token_addr))

	instrs = append(instrs, computebudget.NewSetComputeUnitPriceInstruction(100).Build())
	instrs = append(instrs, computebudget.NewSetComputeUnitLimitInstruction(10000).Build())

	// Create a transfer instruction to send the balance to the recipient
	if vault != "" && balance.Uint64() > 0 {
		instrs = append(instrs, Transfer(TransferParam{
			From:    account,
			To:      solana.MustPublicKeyFromBase58(vault),
			Auth:    wallet.PublicKey(),
			Signers: nil,
			Amount:  balance.Uint64(),
		}))
	}

	closeInst, err := token.NewCloseAccountInstruction(
		account,
		s.account.PublicKey(),
		s.account.PublicKey(),
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	instrs = append(instrs, closeInst)

	tx, err := BuildTransaction(signers, wallet, instrs...)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
