import pathlib
import pkg.http_request as http

PWD = pathlib.Path(__file__).parent.resolve()
USER_NAME = "Jean Bon"
USER_EMAIL = "jean@bon.com"
USER_PASS = "OopsJeanBonHasBeenHacked"
THING_NAME = "smart_plant_1"
THING_KEY = "000001"


def get_secret(name):
    ''' Get a secret from the docker .env file '''
    dot_env_path = PWD.parent.parent / "docker" / ".env"
    with dot_env_path.open() as f:
        for l in f:
            if name in l:
                return l.split('=')[1].strip()
    f.close()


def run_tests():
    admin_mail = get_secret("USER_ADMIN_EMAIL")
    admin_pass = get_secret("USER_ADMIN_PASSWORD")

    admin_token = http.get_admin_token(admin_mail, admin_pass)
    print("admin token: ", admin_token[0:10], "...")

    user_id = http.create_user(admin_token, USER_NAME, USER_EMAIL, USER_PASS)
    print("user id: ", user_id)

    user_token = http.get_user_token(USER_EMAIL, USER_PASS)
    print("user token: ", user_token[0:10], "...")

    thing_id = http.create_thing(user_token, THING_NAME, THING_KEY)
    print("thing id: ", thing_id)

    thing_token = http.get_thing_token(user_token, thing_id)
    print("thing token: ", thing_token[0:10], "...")

    user_id_email = http.get_user_by_email(admin_token, USER_EMAIL)
    print("user id email: ", user_id_email)

    deleted_user_id = http.delete_user(admin_token, user_id_email)
    print("deleted user id: ", deleted_user_id)

    deleted_thing_id = http.delete_thing(user_token, thing_id)
    print("deleted thing id: ", deleted_thing_id)


if __name__ == "__main__":
    run_tests()
