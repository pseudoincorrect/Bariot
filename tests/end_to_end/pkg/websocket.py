
import asyncio
from time import sleep
import websockets
import json
import threading


async def get_thing_id_data(usr_token: str, thing_id: str):
    async with websockets.connect("ws://localhost:80/reader/thing") as websocket:
        print("WS connect and listen")
        rec = None
        msg = {
            "token": usr_token,
            "thingId": thing_id,
        }
        msgJson = json.dumps(msg)
        # print("WS msg json", msgJson)
        await websocket.send(msgJson)

        while True:
            try:
                rec = await websocket.recv()
                print(rec)
            except Exception as e:
                print(e)
                return


if __name__ == "__main__":
    test_token = "123.123.123"
    test_thing_id = "000.000.001"
    done = threading.Event()
    asyncio.run(get_thing_id_data(test_token, test_thing_id))
