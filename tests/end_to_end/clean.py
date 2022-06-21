import asyncio
import pkg.users_and_things as u_and_t
import pkg.senml as senml
import pkg.mqtt as mqtt
import pkg.websocket as ws
from contextlib import suppress
import signal
import sys
import json


def get_from_file():
    with open('data/data.json', 'r') as f:
        data = json.load(f)
    return data


def clean():
    data = get_from_file()
    print(data)
    thing_id = data["thing_id"]
    user_id = data["user_id"]
    u_and_t.delete_user_and_thing(user_id, thing_id)


if __name__ == "__main__":
    clean()
