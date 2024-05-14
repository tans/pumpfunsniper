const axios = require('axios');

const API_URL = 'https://api.devnet.solana.com';  // Change to your Solana cluster's RPC URL
const transactionSignature = '3PgPht5Hr4Pjcv9HqVgcDupcLduqN7rYRUkVk31pKqoycuQ4EcaUdsoLf9vXZ3Xo7A5expvQMVa1ZDN2DwWzfB3F'; // Replace with your actual transaction signature

async function getTransactionDetails() {
    try {
        const response = await axios.post(API_URL, {
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
