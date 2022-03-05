# BARIOT

IoT / IIoT Platform

Bariot Provide a simple solution to send, manage, secure and visualize data from connected devices/machines. 

<p align="center"><img width=50% src="support/images/bariot_img.jpg"></p>

Microservices services architecture build with Go and docker (compose).

Inspired by the beautiful architecture of [Mainflux][mainflux].

## MOTIVATION

It uses the most recent (at year 2022) technologies to create a scalable and cloud agnostic IoT/IIoT platform.

It is opinionated, in the sense that storage and communications technologies are predefined (see COMPONENTS)

Bariot is also a chance to understand what technologies are used to create cloud applications with modern standard of security, scalability, performances and devops practices. Whether it is purely serverless or containerized, these components (see below) in one form or another are often involved. 

## COMPONENTS

### Implemented
- Transport: MQTT
- Authentication: JWT
- USER, THINGS storage: Postgres
- THINGS DATA storage: Influxdb
- Services intercommunication: GRPC
- Services messaging: NATS

### Further on the road
- Transport: OPCUA
- Data presentation: Grafana
- Reverse proxy: ngnix
- Secret storage/management: Vault
- Caching: Redis
- CI/CD: Gitlab

## WORK IN PROGRESS

Authorization, error handling, tests, CI/CD, caching, reverse proxy, etc... are on the way.

[mainflux]: https://mainflux.com/