from api.generator import fun_handler
import asyncio
from account_management.account import Account
import datetime as da
import json

async def main():
    base_url = "https://client-api-2-74b1891ee9f9.herokuapp.com/"
    #base_url = "https://client-api.devnet.pump.fun/"
    r = fun_handler(base_url=base_url)
    a = Account()
    body = await a.sign_in()
    bearer_token = await r.post(path="auth/login", json=body)
    await r.update_headers(bearer_token=bearer_token)
    pubkey = "3pHELg1xpMCy8X5eq73hZwQ3puUyDVu7NiDVqNQptpm9"
    await r.follow(pubkey=pubkey)
    
if __name__ == "__main__":
    asyncio.run(main=main())