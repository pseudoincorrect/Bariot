# BARIOT ROADMAP

<br>

## **TO DO**

|     | Tag    | Description                                       |
| --- | ------ | ------------------------------------------------- |
|     | TEST   | Automate end to end testing with gitlab           |
|     | DOCKER | Pre-compile docker images to speed up deployments |
|     | LOG    | Improve logging                                   |
|     | HTTP   | Add a user/thing data reading HTTP endpoint       |
|     | MAKE   | Create a makefile                                 |
|     |        |                                                   |

<br>
 
## **IN PROCESS**

|     | Tag  | Description                                     |
| --- | ---- | ----------------------------------------------- |
|     | MQTT | Securing and Authorizing MQTT pub/sub with JSON |
|     |      |                                                 |

<br>

## **DONE**

| Date       | Tag      | Description                                     |
| ---------- | -------- | ----------------------------------------------- |
| 21/06/2022 | WS       | Thing data streaming websocket endpoint         |
| 09/06/2022 | TEST     | Test things, users, auth (unit + integration)   |
| 30/05/2022 | REDIS    | Delete thing token from cache + new token case  |
| 29/05/2022 | REDIS    | Use Redis to cache auth data on MQTT endpoint   |
| 28/05/2022 | TEST     | MQTT with python and paho                       |
| 25/05/2022 | RELEASE  | 0.1.0                                           |
| 12/05/2022 | DEPLOY   | Make a working demo with instructions           |
| 11/05/2022 | NGINX    | Add Ngnx reverse proxy with TLS                 |
| 21/04/2022 | GRAFANA  | Add Grafana service                             |
| 17/04/2022 | JWT      | Authorize mqtt with JWT                         |
| 10/03/2022 | INFLUXDB | Put data to influxDb                            |
| 07/03/2022 | HEALTH   | MQTT and NATS healthcheck                       |
| 05/03/2022 | SENML    | Decode SenML messages                           |
| 03/03/2022 | NATS     | Republish mqtt over nats                        |
| 02/03/2022 | NATS     | Add Nats service                                |
| 28/02/2022 | INFLUXDB | Add InfluxDb service                            |
| 26/02/2022 | MQTT     | Get all thing related mqtt messages             |
| 24/02/2022 | GO       | Re-architect the project to use only one go mod |
| 23/02/2022 | MQTT     | Broker and go program works                     |
| 21/02/2022 | AUTH     | Endpoint authentication USER/THINGS             |
| 17/02/2022 | AUTH     | Token generation                                |
| 15/02/2022 | GRPC     | Create a Auth/JWT service with GRPC             |
| 14/02/2022 | CRUD     | User                                            |
| 14/02/2022 | CRUD     | Thing                                           |
| 12/02/2022 | CRUD     | Thing operation in DB                           |
| 11/02/2022 | SQL      | Store thing data in postgres DB                 |
| 07/02/2022 | MQTT     | Add a mqtt broker service                       |
| 05/02/2022 | INIT     | Development environment                         |
|            |          |                                                 |
