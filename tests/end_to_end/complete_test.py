import asyncio
import pkg.users_and_things as u_and_t
import pkg.senml as senml
import pkg.mqtt as mqtt
import pkg.websocket as ws
import json


async def send_few_mqtt_msg(how_many, thing_id, thing_token):
    ''' Send few mqtt message to the bariot broker '''
    for _ in range(how_many):
        try:
            msg = create_mqtt_msg(thing_token)
            send_mqtt_msg(thing_id, msg)
            await asyncio.sleep(0.5)
        except Exception as e:
            print(e)


def create_mqtt_msg(thing_token):
    ''' Create a random message with Senml and token '''
    payload = senml.create_random_senml()
    return mqtt.format_message(thing_token, payload)


def send_mqtt_msg(thing_id, msg):
    ''' Send a message to the MQTT broker '''
    mqtt.send_message(thing_id, msg)


def save_to_file(json_data, path_data):
    ''' Save a json object to a file '''
    with open(path_data, 'w', encoding='utf-8') as f:
        json.dump(json_data, f, ensure_ascii=False, indent=4)


def start_websocket_task(thing_id, user_token):
    ''' Start the websocket background task '''
    return asyncio.create_task(
        ws.get_thing_id_data(user_token, thing_id))


async def create_send_delete():
    ''' Create a user and thing, send a message to it. Delete them. '''
    user_id, thing_id = u_and_t.create_user_and_thing()
    user_token = u_and_t.get_user_token()
    thing_token = u_and_t.get_thing_token(user_token, thing_id)

    data = {"user_id": user_id, "thing_id": thing_id}
    save_to_file(data, 'data/data.json')

    task = start_websocket_task(thing_id, user_token)

    await asyncio.sleep(0.5)

    await send_few_mqtt_msg(5, thing_id, thing_token)

    await asyncio.sleep(0.5)

    task.cancel()
    u_and_t.delete_user_and_thing(user_id, thing_id)


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    loop.run_until_complete(create_send_delete())
