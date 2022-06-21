from time import sleep
import paho.mqtt.client as mqtt
import json

BARIOT_HOST = "localhost"
SILENT = True


class MyMQTTClass(mqtt.Client):
    def on_connect(self, mqttc, obj, flags, rc):
        if not SILENT:
            print("MQTT: " + "connection result (rc): "+str(rc))

    def on_connect_fail(self, mqttc, obj):
        if not SILENT:
            print("MQTT: " + "Connect failed")

    def on_message(self, mqttc, obj, msg):
        if not SILENT:
            print("MQTT: " + "RECEIVED MSG: topic: " + msg.topic + ", QoS: " +
                  str(msg.qos) + ", payload: " + str(msg.payload))

    def on_publish(self, mqttc, obj, mid):
        if not SILENT:
            print("MQTT: " + "message ID (mid): "+str(mid))

    def on_subscribe(self, mqttc, obj, mid, granted_qos):
        if not SILENT:
            print("MQTT: " + "Subscribed: "+str(mid)+" "+str(granted_qos))

    def on_log(self, mqttc, obj, level, string):
        if not SILENT:
            print("MQTT: " + string)

    def subscribe_all(self):
        self.subscribe(topic="#", qos=0)


def format_message(thing_token, payload):
    ''' Format a message to be sent to the MQTT broker.'''
    msg = {}
    msg["token"] = thing_token
    msg["Records"] = json.loads(payload)
    json_msg = json.dumps(msg)
    print("MQTT msg :", json.dumps(msg, indent=4))
    return json_msg


def make_thing_topic(thing_id):
    return "things/" + thing_id


def send_message(thing_id, msg):
    ''' Send a message to the MQTT broker.'''
    topic = make_thing_topic(thing_id)
    mqttc = MyMQTTClass()
    mqttc.connect(BARIOT_HOST, 1883, 60)
    mqttc.subscribe_all()
    mqttc.loop_start()
    mqttc.publish(topic, msg, qos=2)
    timeout = 2
    while timeout > 0:
        timeout -= 1
        sleep(1)
