require('dotenv').config();
const axios = require('axios');
const { PublicKey, Transaction, SystemProgram, LAMPORTS_PER_SOL, Keypair } = require('@solana/web3.js');
const bs58 = require('bs58');

const API_URL = process.env.PUMP_DEVNET_CLIENT_API_URL;
const JITO_RPC_URL = 'https://devnet.helius-rpc.com/?api-key=02befe47-b808-4837-8ce3-409c845b79bb';
const PRIVATE_KEY = process.env.PRIVATE_KEY_DEVNET;  // Base58 encoded private key
const RECIPIENT_ADDRESS = 'FYiJNJko7R4FvRnB1RZUzncBHXQ445wHEbG8GZWDTnoa';
const AMOUNT_SOL = 0.1; // amount in SOL to transfer

async function getLatestBlockhash() {
    const response = await axios.post('https://devnet.helius-rpc.com/?api-key=02befe47-b808-4837-8ce3-409c845b79bb', {
        jsonrpc: "2.0",
        id: 1,
        method: "getLatestBlockhash",
        params: []
    }, {
        headers: { "Content-Type": "application/json" }
    });

    if (response.data.error) {
        console.error('Failed to fetch latest blockhash:', response.data.error);
        return null;
    }

    return response.data.result.value.blockhash;
}

async function sendTransactionWithPriorityFee() {
    const secretKey = bs58.decode(PRIVATE_KEY);
    if (secretKey.length !== 64) {
        console.error('Invalid secret key length:', secretKey.length);
        return;
    }
    const fromKeypair = Keypair.fromSecretKey(secretKey);

    const blockhash = await getLatestBlockhash();
    if (!blockhash) {
        console.error('Failed to obtain latest blockhash');
        return;
    }

    const transaction = new Transaction();
    transaction.recentBlockhash = blockhash;
    transaction.add(
        SystemProgram.transfer({
            fromPubkey: fromKeypair.publicKey,
            toPubkey: new PublicKey(RECIPIENT_ADDRESS),
            lamports: AMOUNT_SOL * LAMPORTS_PER_SOL
        })
    );

    transaction.sign(fromKeypair);

    try {
        const serializedTransaction = transaction.serialize().toString('base64');
        const result = await axios.post(JITO_RPC_URL, {
            jsonrpc: "2.0",
            id: 1,
            method: "sendTransaction",
            params: [serializedTransaction, {encoding: "base64"}]
        }, {
            headers: { "Content-Type": "application/json" }
        });

        if (result.data.error) {
            console.error('Transaction failed:', result.data.error);
        } else {
            console.log('Transaction sent successfully with signature:', result.data.result);
        }
    } catch (error) {
        console.error('Error sending transaction:', error.response ? error.response.data : error.message);
    }
}

sendTransactionWithPriorityFee();


const EXPLORER_API_URL = 'https://api.devnet.solana.com';  // Change to your Solana cluster's RPC URL
const transactionSignature = '21ixeHNY4aL1NaYxe1fMpcoCtk6jMNiHFELRZSmsVHcZtrfs2CZBnXHznUKHwPGJg7NL2tcbScWdGDDHTbszbDy5'; // Replace with your actual transaction signature

async function getTransactionDetails() {
    try {
        const response = await axios.post(EXPLORER_API_URL, {
            jsonrpc: "2.0",
            id: 1,
            method: "getConfirmedTransaction",
            params: [
                transactionSignature,
                "json"
            ]
        }, {
            headers: { 'Content-Type': 'application/json' }
        });

        if (response.data.result) {
            console.log('Transaction details:', JSON.stringify(response.data.result, null, 2));
        } else {
            console.log('Transaction not found or not confirmed yet.');
        }
    } catch (error) {
        console.error('Error retrieving transaction:', error);
    }
}

getTransactionDetails();