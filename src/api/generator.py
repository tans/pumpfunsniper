import httpx 
from typing import Dict
import asyncio
import json
import os
import ast

class fun_handler:
    def __init__(self, base_url: str, headers: Dict[str,str]=None, cookies: Dict[str,str]=None):
        if headers is None:
            #path = os.path.abspath("src/api/headers.json")
            path = ("./api/headers.json")
            with open(path, "r") as f:
                self.headers = json.load(fp=f)
        if cookies is None:
            cookies = {}
        self.url = base_url



    async def get(self, relative_path, headers=None, cookies=None):
        if headers is None:
            headers = self.headers
        response = httpx.get(url=self.url + relative_path, headers=headers, cookies=cookies)
        return response.text
    
    
    async def get_full(self, url, headers=None, cookies=None):
        if headers is None:
            headers = self.headers
        response = httpx.get(url=url, headers=headers, cookies=cookies)
        return response.text
    
    async def post(self, path, json: json = None, headers: json = None, cookies=None):
        if headers is None:
            headers = self.headers
        response = httpx.post(url=self.url + path, json=json, headers=headers)
        print (response.text)
        return response.text

    async def update_headers(self, bearer_token: dict):
        bearer_token = json.loads(bearer_token)
        self.headers["Authorization"] = f"Bearer {bearer_token["access_token"]}"

    async def follow(self, pubkey: str):
        response = await self.post(path="following/" + pubkey)
        print (response)
        return response