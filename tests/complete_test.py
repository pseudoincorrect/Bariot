import pathlib
import requests

PWD = pathlib.Path(__file__).parent.resolve()
HTTP_PROT = "http://"
BARIOT_HOST = "localhost"
USER_NAME = "Jean Bon"
USER_EMAIL = "jean@bon.com"
USER_PASS = "OopsJeanBonHasBeenHacked"
THING_NAME = "smart_plant_1"
THING_KEY = "000001"


def get_secret(name):
    ''' Get a secret from the docker .env file '''
    dot_env_path = PWD.parent / "docker" / ".env"
    with dot_env_path.open() as f:
        for l in f:
            if name in l:
                return l.split('=')[1].strip()
    f.close()


def get_admin_token(mail, passw):
    url = HTTP_PROT + BARIOT_HOST + '/users/login/admin'
    headers = {"Content-Type": "application/json"}
    data = {"Email": mail, "Password": passw}
    res = requests.post(url, headers=headers, json=data)
    if res.status_code == 200:
        return res.json()["Token"]
    print(res.reason)
    return None


def get_user_token(mail, passw):
    url = HTTP_PROT + BARIOT_HOST + '/users/login'
    headers = {"Content-Type": "application/json"}
    data = {"Email": mail, "Password": passw}
    res = requests.post(url, headers=headers, json=data)
    if res.status_code == 200:
        return res.json()["Token"]
    print(res.reason)
    return None


def get_user_by_email(admin_token, email):
    url = HTTP_PROT + BARIOT_HOST + '/users/email/' + email
    headers = {"Content-Type": "application/json",
               "Authorization": admin_token}
    res = requests.get(url, headers=headers)
    if res.status_code == 200:
        return res.json()["Id"]
    print(res.reason)
    return None


def create_user(admin_token, name, email, passw):
    url = HTTP_PROT + BARIOT_HOST + '/users/'
    headers = {"Content-Type": "application/json",
               "Authorization": admin_token}
    data = {"FullName": name, "Email": email, "Password": passw}
    res = requests.post(url, headers=headers, json=data)
    if res.status_code == 200:
        return res.json()["Id"]
    print(res.reason)
    return None


def delete_user(admin_token, user_id):
    url = HTTP_PROT + BARIOT_HOST + '/users/' + user_id
    headers = {"Content-Type": "application/json",
               "Authorization": admin_token}
    res = requests.delete(url, headers=headers)
    if res.status_code == 200:
        return res.json()["Id"]
    print(res.reason)
    return None


def create_thing(user_token, name, key):
    url = HTTP_PROT + BARIOT_HOST + '/things/'
    headers = {"Content-Type": "application/json",
               "Authorization": user_token}
    data = {"Name": name, "Key": key}
    res = requests.post(url, headers=headers, json=data)
    if res.status_code == 200:
        return res.json()["Id"]
    print(res.reason)
    return None


def delete_thing(user_token, thing_id):
    url = HTTP_PROT + BARIOT_HOST + '/things/' + thing_id
    headers = {"Content-Type": "application/json",
               "Authorization": user_token}
    res = requests.delete(url, headers=headers)
    if res.status_code == 200:
        return res.json()["Id"]
    print(res.reason)
    return None


def get_thing_token(user_token, thing_id):
    url = HTTP_PROT + BARIOT_HOST + '/things/' + thing_id + "/token"
    headers = {"Content-Type": "application/json",
               "Authorization": user_token}
    res = requests.get(url, headers=headers)
    if res.status_code == 200:
        return res.json()["Token"]
    print(res.reason)
    return None


def run_tests():
    admin_mail = get_secret("USER_ADMIN_EMAIL")
    admin_pass = get_secret("USER_ADMIN_PASSWORD")

    admin_token = get_admin_token(admin_mail, admin_pass)
    print("admin token: ", admin_token[0:10], "...")

    user_id = create_user(admin_token, USER_NAME, USER_EMAIL, USER_PASS)
    print("user id: ", user_id)

    user_token = get_user_token(USER_EMAIL, USER_PASS)
    print("user token: ", user_token[0:10], "...")

    thing_id = create_thing(user_token, THING_NAME, THING_KEY)
    print("thing id: ", thing_id)

    thing_token = get_thing_token(user_token, thing_id)
    print("thing token: ", thing_token[0:10], "...")

    user_id_email = get_user_by_email(admin_token, USER_EMAIL)
    print("user id email: ", user_id_email)

    deleted_user_id = delete_user(admin_token, user_id_email)
    print("deleted user id: ", deleted_user_id)

    deleted_thing_id = delete_thing(user_token, thing_id)
    print("deleted thing id: ", deleted_thing_id)


if __name__ == "__main__":
    run_tests()
