
import asyncio
import websockets
import json


async def hello():
    async with websockets.connect("ws://localhost:80/reader/thing") as websocket:
        msg = {
            "token": "123.123.123",
            "thingId": "000.000.001"
        }
        msgJson = json.dumps(msg)
        await websocket.send(msgJson)
        await websocket.recv()

asyncio.run(hello())
