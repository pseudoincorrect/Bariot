import requests

HTTP_PROT = "http://"
BARIOT_HOST = "localhost"


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
