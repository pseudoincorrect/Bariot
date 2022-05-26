import pathlib
import requests

PWD = pathlib.Path(__file__).parent.resolve()
BARIOT_HOST = "localhost"


def get_secret(name):
    dot_env_path = PWD.parent / "docker" / ".env"
    with dot_env_path.open() as f:
        for l in f:
            if name in l:
                return l.split('=')[1].strip()
    f.close()


def get_admin_token(mail, passw):
    url = "http://" + BARIOT_HOST + '/users/login/admin'
    print(url)
    headers = {"Content-Type": "application/json"}
    data = {"Email": mail, "Password": passw, }
    res = requests.post(url, headers=headers, json=data)
    print("Status Code", res.status_code)
    print("JSON Response ", res.json())


admin_mail = get_secret("USER_ADMIN_EMAIL")
admin_pass = get_secret("USER_ADMIN_PASSWORD")

print(admin_mail)
print(admin_pass)

get_admin_token(admin_mail, admin_pass)
