# BARIOT

[IoT][IoT] / [IIoT][IIoT] Platform

Bariot is a simple solution to send, manage, secure and visualize data from connected devices/machines. 

<p align="center"><img width=50% src="support/images/bariot_img.jpg"></p>

Micro-services architecture build with [Go][Go] and [Docker][Docker] (compose).

Inspired by the beautiful architecture of [Mainflux][Mainflux].

<br/>

--- 
<br/>

## MOTIVATIONS

Bariot is being developed to offer a simple, complete and scalable solution to wide variety of IoT projects.

Bariot uses the most recent (at year 2022) technologies to create a scalable and cloud agnostic IoT/IIoT platform.

Bariot is opinionated, in the sense that storage and communications technologies are predefined (see COMPONENTS)

Bariot is a chance to understand what technologies are used to create cloud applications with modern standard of security, scalability, performances and devops practices. Whether it is purely serverless or containerized, these components (see below) in one form or another are often involved. 

Bariot is also a fun project to work on !

<br/>

--- 
<br/>

## COMPONENTS

### Implemented
- Transport: [MQTT][MQTT]
- Authentication/Authorization: [JWT][JWT]
- USER, THINGS storage: [PostgreSQL][PostgreSQL]
- THINGS DATA storage: [Influxdb][Influxdb]
- Services intercommunication: [gRpc][gRpc]
- Services messaging: [Nats][Nats]

### Further on the road
- Transport: [OPC-UA][OPC-UA]
- Data presentation: [Grafana][Grafana]
- Reverse proxy: [Ngnix][Ngnix]
- Secret storage/management: [Vault][Vault]
- Caching: [Redis][Redis]
- CI/CD: [Gitlab][Gitlab]
- Permission system: to be decided

<br/>

--- 
<br/>

## WORK IN PROGRESS !!!


<br/>

--- 
<br/>

[IoT]: https://www.zdnet.com/article/what-is-the-internet-of-things-everything-you-need-to-know-about-the-iot-right-now/
[IIoT]: https://www.trendmicro.com/vinfo/us/security/definition/industrial-internet-of-things-iiot
[Go]: https://www.freecodecamp.org/news/what-is-go-programming-language/
[Docker]: https://docs.docker.com/get-started/overview/
[Mainflux]: https://mainflux.com/
[MQTT]: https://mqtt.org/
[JWT]: https://jwt.io/
[PostgreSQL]: https://www.postgresql.org/
[Influxdb]: https://www.influxdata.com/
[gRpc]: https://grpc.io/docs/what-is-grpc/introduction/
[Nats]: https://docs.nats.io/nats-concepts/what-is-nats
[OPC-UA]: https://www.opc-router.com/what-is-opc-ua/
[Grafana]: https://www.scaleyourapp.com/what-is-grafana-why-use-it-everything-you-should-know-about-it/
[Ngnix]: https://medium.com/globant/understanding-nginx-as-a-reverse-proxy-564f76e856b2
[Redis]: https://redis.io/topics/introduction
[Vault]:https://www.vaultproject.io/docs/what-is-vault
[Gitlab]: https://about.gitlab.com/what-is-gitlab/