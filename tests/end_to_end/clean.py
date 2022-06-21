import pkg.users_and_things as u_and_t
import json


def get_from_file():
    with open('data/data.json', 'r') as f:
        data = json.load(f)
    return data


def clean():
    data = get_from_file()
    print(data)
    if data["thing_id"]:
        u_and_t.delete_thing(data["thing_id"])
    if data["user_id"]:
        u_and_t.delete_user(data["user_id"])


if __name__ == "__main__":
    clean()
