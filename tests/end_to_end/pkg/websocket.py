
import asyncio
import websockets
import json
import threading


async def get_thing_id_data(usr_token: str, thing_id: str):
    ''' Will connect to bariot with websocket and poll for received data '''
    async with websockets.connect("ws://localhost:80/reader/thing") as websocket:
        print("WebSocket connected and listening")
        rec = None
        msg = {
            "token": usr_token,
            "thingId": thing_id,
        }
        msgJson = json.dumps(msg)
        await websocket.send(msgJson)
        while True:
            try:
                rec = await websocket.recv()
                print("Received SenML message: ",
                      rec[:20], " ... ", rec[len(rec)-20:])
            except Exception as e:
                print(e)
                return


if __name__ == "__main__":
    test_token = "123.123.123"
    test_thing_id = "000.000.001"
    done = threading.Event()
    asyncio.run(get_thing_id_data(test_token, test_thing_id))
