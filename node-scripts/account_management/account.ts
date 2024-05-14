import * as dotenv from 'dotenv';
import { Keypair } from '@solana/web3.js';
import * as nacl from 'tweetnacl';
import axios from 'axios';

dotenv.config();

class Account {
    private keypair: Keypair;

    constructor() {
        const privateKey = process.env.PRIVATE_KEY;
        if (!privateKey) {
            throw new Error('Private key not found in environment variables');
        }
        const secretKey = Buffer.from(privateKey, 'base64');
        this.keypair = Keypair.fromSecretKey(secretKey);
    }

    public async signIn(timestamp: Date): Promise<any> {
        const message = new TextEncoder().encode(timestamp.toISOString());
        const signature = nacl.sign.detached(message, this.keypair.secretKey);
        const body = {
            address: this.keypair.publicKey.toBase58(),
            signature: Buffer.from(signature).toString('base64'),
            timestamp: timestamp.toISOString()
        };
        const headers = {
            'Content-Type': 'application/json'
            // Add any other necessary headers here
        };

        try {
            const response = await axios.post('https://client-api.devnet.pump.fun/auth/login', body, { headers });
            return response.data;
        } catch (error) {
            console.error('Error in signIn:', error);
            throw error;
        }
    }
}

async function main() {
    const account = new Account();
    const signInResult = await account.signIn(new Date());
    console.log('Sign In Result:', signInResult);
}

main();