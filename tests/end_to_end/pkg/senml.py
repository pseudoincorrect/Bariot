from kpn_senml import *
import time


def create_random_senml():
    ''' Create a random SenML message with given units. '''
    pack = SenmlPack("")
    pack.add(SenmlRecord("temperature", unit="degree-c", value=23.5))
    pack.add(SenmlRecord("humidity", unit="percents", value=73))
    pack.add(SenmlRecord("heart-rate", unit="bpm", value=86))
    pack.base_time = time.time()
    return pack.to_json()
