import pkg.users_and_things as u_and_t
import pkg.senml as senml
import pkg.mqtt as mqtt


def create_mqtt_msg(thing_token):
    payload = senml.create_random_senml()
    return mqtt.format_message(thing_token, payload)


def send_mqtt_msg(thing_id, msg):
    mqtt.send_message(thing_id, msg)


def create_send_delete():
    user_id, thing_id = u_and_t.create_user_and_thing()
    user_token = u_and_t.get_user_token()
    thing_token = u_and_t.get_thing_token(user_token, thing_id)
    try:
        msg = create_mqtt_msg(thing_token)
        send_mqtt_msg(thing_id, msg)
        # send_mqtt_msg(thing_id, msg)
        # send_mqtt_msg(thing_id, msg)
    except Exception as e:
        print(e)
    u_and_t.delete_user_and_thing(user_id, thing_id)


if __name__ == "__main__":
    create_send_delete()
