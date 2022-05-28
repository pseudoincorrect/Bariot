from time import sleep
import paho.mqtt.client as mqtt

BARIOT_HOST = "localhost"
THING_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiYWRtaW4iLCJleHAiOjE2NTM1NzQ3NDUsImlhdCI6MTY1MzU3MTE0NSwiaXNzIjoiZGV2X2xvY2FsIiwic3ViIjoiMCJ9.HGkFC8CJo3ZbnRA9YYYJ6qNjtOuMbxbTzoYorkXik-M"
THING_ID = "55084249-1877-4a70-973e-2a5ceaa3e642"
THING_TOPIC = "things/" + THING_ID


class MyMQTTClass(mqtt.Client):

    def on_connect(self, mqttc, obj, flags, rc):
        print("rc: "+str(rc))

    def on_connect_fail(self, mqttc, obj):
        print("Connect failed")

    def on_message(self, mqttc, obj, msg):
        print("RECEIVED MSG: topic: " + msg.topic + ", QoS: " +
              str(msg.qos) + ", payload: " + str(msg.payload))

    def on_publish(self, mqttc, obj, mid):
        print("mid: "+str(mid))

    def on_subscribe(self, mqttc, obj, mid, granted_qos):
        print("Subscribed: "+str(mid)+" "+str(granted_qos))

    def on_log(self, mqttc, obj, level, string):
        print(string)

    def run(self):
        self.connect(BARIOT_HOST, 1883, 60)
        self.subscribe(THING_TOPIC, 0)
        self.loop_start()

        rc = self.publish(THING_TOPIC, "useless message 2", qos=2)

        timeout = 2
        while timeout > 0:
            timeout -= 1
            sleep(1)


def send_a_message(thing_id, thing_token, message):
    pass


def run_tests():
    mqttc = MyMQTTClass()
    rc = mqttc.run()


if __name__ == "__main__":
    run_tests()
