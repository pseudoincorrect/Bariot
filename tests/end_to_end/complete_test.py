import asyncio
import pkg.users_and_things as u_and_t
import pkg.senml as senml
import pkg.mqtt as mqtt
import pkg.websocket as ws
from contextlib import suppress
import signal
import sys
import json


user_id, thing_id = None, None


def create_mqtt_msg(thing_token):
    ''' Create a random message with Senml and token '''
    payload = senml.create_random_senml()
    return mqtt.format_message(thing_token, payload)


def send_mqtt_msg(thing_id, msg):
    ''' Send a message to the MQTT broker.'''
    mqtt.send_message(thing_id, msg)


def signal_handler(sig, frame):
    u_and_t.delete_user_and_thing(user_id, thing_id)
    sys.exit(0)


def save_to_file(json_data):
    with open('data/data.json', 'w', encoding='utf-8') as f:
        json.dump(json_data, f, ensure_ascii=False, indent=4)


def get_from_file():
    json_data = None
    with open('data/data.json', 'w', encoding='utf-8') as f:
        json.load(json_data, f, ensure_ascii=False, indent=4)
    return json_data


async def create_send_delete():
    ''' Create a user and thing, send a message to it. Delete them. '''
    user_id, thing_id = u_and_t.create_user_and_thing()
    user_token = u_and_t.get_user_token()
    thing_token = u_and_t.get_thing_token(user_token, thing_id)
    data = {"user_id": user_id, "thing_id": thing_id}
    save_to_file(data)

    # task = None
    # try:
    #     task = asyncio.Task(ws.getThingIdData(user_token, thing_id))
    # except Exception as e:
    #     print(e)

    # await asyncio.sleep(3)

    # print("Sending MQTT messages")
    # try:
    #     msg = create_mqtt_msg(thing_token)
    #     send_mqtt_msg(thing_id, msg)
    # except Exception as e:
    #     print(e)

    # print("Waiting for WS responses")

    # await asyncio.sleep(3)

    # task.cancel()
    # with suppress(asyncio.CancelledError):
    #     await task  #

    # u_and_t.delete_user_and_thing(user_id, thing_id)


# async def run_main_loop():

signal.signal(signal.SIGINT, signal_handler)
loop = asyncio.new_event_loop()
asyncio.set_event_loop(loop)
try:
    loop.run_until_complete(create_send_delete())
finally:
    loop.run_until_complete(loop.shutdown_asyncgens())
    loop.close()

# if __name__ == "__main__":
# await run_main_loop()
