import asyncio
import threading
import time
import pkg.users_and_things as u_and_t
import pkg.senml as senml
import pkg.mqtt as mqtt
import pkg.websocket as ws
from contextlib import suppress
import json


def create_mqtt_msg(thing_token):
    ''' Create a random message with Senml and token '''
    payload = senml.create_random_senml()
    return mqtt.format_message(thing_token, payload)


def send_mqtt_msg(thing_id, msg):
    ''' Send a message to the MQTT broker.'''
    mqtt.send_message(thing_id, msg)


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

    task = asyncio.create_task(
        ws.get_thing_id_data(user_token, thing_id))
    await asyncio.sleep(3)

    print("Sending MQTT messages")
    for _ in range(10):
        try:
            msg = create_mqtt_msg(thing_token)
            send_mqtt_msg(thing_id, msg)
            await asyncio.sleep(1)
        except Exception as e:
            print(e)

    await asyncio.sleep(3)

    task.cancel()
    u_and_t.delete_user_and_thing(user_id, thing_id)

if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(create_send_delete())
