import asyncio
from solana.rpc.async_api import AsyncClient
from dotenv import dotenv_values
import solders
import base58
import time
import json

class Account:
    def __init__(self):
        self.pk = dotenv_values()["PRIVATE_KEY"]
        self.keypair = solders.keypair.Keypair.from_base58_string(self.pk)
        self.public_key = self.keypair.pubkey()
        
    async def sign_in(self, timestamp: int = None)-> dict:
        if timestamp is None:
            timestamp = int(time.time() * 1000)
            timestamp = 1713298780638
            to_sign = f"Sign in to pump.fun: {timestamp}"
        signature = self.keypair.sign_message(to_sign.encode())
        body = {"address":str(self.keypair.pubkey()), "signature":str(signature) , "timestamp": int(timestamp)}
        return body

