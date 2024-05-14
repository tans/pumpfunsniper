// const axios = require('axios');
// const { Keypair } = require('@solana/web3.js');
// const bs58 = require('bs58');
// const nacl = require('tweetnacl');
// require('dotenv').config();

// const API_URL = process.env.PUMP_CLIENT_API_URL;
// const PRIVATE_KEY = process.env.PRIVATE_KEY;

// async function login() {
//     // Decode the private key and create a keypair
//     const secretKey = bs58.decode(PRIVATE_KEY);
//     const keypair = Keypair.fromSecretKey(secretKey);

//     // Prepare the message and timestamp
//     const timestamp = new Date().getTime(); // Current time in milliseconds
//     console.log('Timestamp:', timestamp);
//     const message = new TextEncoder().encode(`Sign in to pump.fun: ${timestamp}`);
//     console.log('Message:', message);

//     // Sign the message using nacl
//     const signature = nacl.sign.detached(message, keypair.secretKey);
//     const signatureBase58 = bs58.encode(signature);
//     console.log('Signature:', signatureBase58);

//     // Payload for the POST request
//     const payload = {
//         address: keypair.publicKey.toBase58(),
//         signature: signatureBase58,
//         timestamp: timestamp
//     };
//     console.log('Payload:', payload);

//     // POST request to the login endpoint
//     try {
//         const response = await axios.post(`${API_URL}auth/login`, payload, {
//             headers: { 'Content-Type': 'application/json' }
//         });
//         if (response.status === 200 || response.status === 201) {
//             console.log('Login successful. Access token:', response.data.access_token);
//         } else {
//             throw new Error(`Failed to log in: ${response.status} ${response.statusText}`);
//         }
//     } catch (error) {
//         console.error('Error during login:', error.response ? error.response.data : error.message);
//     }
// }

// login();
