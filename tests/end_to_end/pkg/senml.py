from kpn_senml import *
import time
import datetime


def create_random_senml() -> str:
    pack = SenmlPack("")

    pack.add(SenmlRecord("temperature", unit="degree-c", value=23.5))
    pack.add(SenmlRecord("humidity", unit="percents", value=73))
    pack.add(SenmlRecord("heart-rate", unit="bpm", value=86))

    pack.base_time = time.time()

    return pack.to_json


def run_tests():
    json = create_random_senml()
    print(json())


if __name__ == "__main__":
    run_tests()
