require('dotenv').config();
const axios = require('axios');
const { PublicKey, Transaction, SystemProgram, LAMPORTS_PER_SOL, Keypair, Connection, sendAndConfirmTransaction } = require('@solana/web3.js');
const { Token } = require('@solana/spl-token');
const bs58 = require('bs58');

const RPC_URL = 'https://devnet.helius-rpc.com/?api-key=02befe47-b808-4837-8ce3-409c845b79bb';
const PRIVATE_KEY = process.env.PRIVATE_KEY_DEVNET;  // Base58 encoded private key
const SMART_CONTRACT_ADDRESS = '4dsTVfHALQ4vFhwSZmrzuu6jdqM3ytMENEtpKYDDw3Jv'; // Address of the smart contract to handle the swap
const AMOUNT_SOL = 0.1; // Amount in SOL to swap

async function main() {
    const secretKey = bs58.decode(PRIVATE_KEY);
    const fromKeypair = Keypair.fromSecretKey(secretKey);
    const connection = new Connection(RPC_URL, 'confirmed');

    const { blockhash } = await connection.getLatestBlockhash('finalized');

    const transaction = new Transaction({
        feePayer: fromKeypair.publicKey,
        recentBlockhash: blockhash
    });

    transaction.add(
        SystemProgram.transfer({
            fromPubkey: fromKeypair.publicKey,
            toPubkey: new PublicKey(SMART_CONTRACT_ADDRESS),
            lamports: AMOUNT_SOL * LAMPORTS_PER_SOL
        })
    );

    transaction.sign(fromKeypair);

    try {
        const signature = await sendAndConfirmTransaction(connection, transaction, [fromKeypair]);
        console.log('Transaction sent successfully with signature:', signature);
        return signature;
    } catch (error) {
        console.error('Error sending transaction:', error);
    }
}

main();
