
import asyncio
from time import sleep
import websockets
import json


async def getThingIdData(token: str, thing_id: str):
    async with websockets.connect("ws://localhost:80/reader/thing") as websocket:
        rec = None
        msg = {
            "token": token,
            "thingId": thing_id,
        }
        msgJson = json.dumps(msg)
        print("WS msg json", msgJson)
        await websocket.send(msgJson)

        while True:
            try:
                rec = await websocket.recv()
                print(rec)
            except Exception as e:
                print("could not receive data from websocket")
                print(e)
                return
            sleep(0.1)


if __name__ == "__main__":
    test_token = "123.123.123"
    test_thing_id = "000.000.001"
    asyncio.run(getThingIdData(test_token, test_thing_id))
