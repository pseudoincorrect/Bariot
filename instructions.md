# INSTRUCTIONS AND INSTALLATION OF BARIOT

### These instructions seems long, but most of them are about installing docker, git and go.

### The real deployment of Bariot is not so long/complex, see "BARIOT DEPLOYMENT" part.

<br/>

---

---

<br/>

# INSTALLATION ON AWS EC2 VIRTUAL MACHINE

## **EC2 VM PREPARATION**

## Start an EC2 machine (T2 medium used here)

Do not use anything lower than a T2 medium (2 vCPU, 4 GB RAM) !

For the storage (EBS), 8GB is enough.

Change the security group of the VM to include:

| PROTOCOL | PORT | SOURCE    |
| -------- | ---- | --------- |
| SSH      | 22   | 0.0.0.0/0 |
| HTTP     | 80   | 0.0.0.0/0 |
| HTTPS    | 443  | 0.0.0.0/0 |
| MQTT     | 1883 | 0.0.0.0/0 |
| MQTTS    | 8883 | 0.0.0.0/0 |

<br/>

## SSH to your VM (with putty, [Tuto][puttyec2])

<br/>

## Install Docker ([Tuto][docker])

```console
$ sudo yum update -y
$ sudo amazon-linux-extras install docker
$ sudo service docker start
$ sudo usermod -a -G docker ec2-user
```

## Install Docker-compose ([Tuto][docker-compose])

```console
$ cd
$ sudo curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-`uname -s`-`uname -m` | sudo tee /usr/local/bin/docker-compose > /dev/null
$ sudo chmod +x /usr/local/bin/docker-compose
$ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
$ sudo reboot
```

## Optional: install git and a text editor ([Micro][micro])

```console
$ cd && curl https://getmic.ro | bash && sudo mv micro /usr/local/bin/micro
$ sudo yum instal git
```

## Install Go ([Tuto][installgo])

```console
$ cd
$ wget https://go.dev/dl/go1.18.2.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
```

Edit your bash_profile to export GO bins and GOPATH

```console
$ mkdir -p go/projects
$ micro ~/.bash_profile
... text editing, see just bellow ...
$ source ~/.bash_profile
```

Add "PATH=$PATH:/usr/local/go/bin"

Add "export GOPATH="/home/ec2-user/go"

Like so:

```sh
# Optional
alias cdbariot="cd /home/ec2-user/go/projects/Bariot/docker"
bind 'TAB:menu-complete'
bind '"\e[A":history-search-backward'
bind '"\e[B":history-search-forward'
# Mandatory
PATH=$PATH:$HOME/.local/bin:$HOME/bin
PATH=$PATH:/usr/local/go/bin
export PATH
export GOPATH="/home/ec2-user/go"
```

Few alias and settings have been added to simplify the development.

<br/>

---

## **BARIOT DEPLOYMENT**

<br/>

## Clone Bariot [repository][bariot]

```console
$ cd go/projects
$ git clone https://github.com/pseudoincorrect/Bariot.git
```

## Create self-signed SSL certificate ([Tuto][openssl])

```console
$ cd /home/ec2-user/go/projects/Bariot/docker/nginx
$ mkdir ssl && cd ssl
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout self-signed.key -out self-signed.crt
$ openssl dhparam -out dhparam.pem 2048
```

The dhparam command will takes few secs.

## Change secrets (optional for quick testing)

```console
$ cd /home/ec2-user/Bariot/docker
$ micro .env
```

Modify the following secrets:

- THINGS_DB_PASSWORD
- USERS_DB_PASSWORD
- USER_ADMIN_PASSWORD
- JWT_SECRET
- MQTT_PASS
- INFLUXDB_PASSWORD
- INFLUXDB_TOKEN

You can use [randomkeygen][randomkeygen] to generate strong secrets.

It is planned to use Vault to manage these secrets later on.

## Start the system and display logs

```console
$ cd /home/ec2-user/go/projects/Bariot/docker
$ sudo service docker start
$ docker-compose up -d
$ docker-compose logs -f
```

It will take quite some time the first time your run it, 10 mins on a T2 medium...

Bariot is downloading every docker images, go packages and build all services.

**! At the moment, timing errors can happen with the start of multiple containers !**

Re-running "docker-compose up -d" usually do the trick.

Roadmap: deployment with precompiled images.

<br/>

---

---

<br/>

# USAGE / TEST OF BARIOT

<br>

## USING THE HTTP API WITH CURL (CURL) with a linux terminal (Bash)

Save the EC2 public address

```console
$ export BARIOT_HOST=ec2-xx-xx-xxx-xx.eu-west-1.compute.amazonaws.com
```

## Get an admin token

**get your admin email and password from "Bariot/docker/.env" file and replace in the export bellow.**

export ADMIN_MAIL=....
export ADMIN_PASS=....

```console
$ curl -L --request POST \
--header "Content-Type: application/json" \
--data '{"Email" : $ADMIN_MAIL, "Password": $ADMIN_PASS}' \
$BARIOT_HOST/users/login/admin
```

Since the APIs are behing a reverse proxy (nginx) we need to use the **-L** option of Curl.

Example response:

```json
{ "Token": "xxxxxx.xxxxxxx(...)xxxxxxx.xxxxxx" }
```

Save the receives admin token.

```console
$ export ADMIN_TOKEN=xxxxxx.xxxx(...)xxxx.xxxx
```

## Create a user

```console
$ curl -L --request POST \
--header "Content-Type: application/json" \
--header "Authorization: $ADMIN_TOKEN" \
--data '{"FullName": "Jacques Cellaire", "Email": "jacques@cellaire.com", "Password": "OopsjacquesHasBeenHacked"}' \
$BARIOT_HOST/users/
```

Example response:

```json
{
  "Id": "0ed50174-ba43-43f8-93f0-a971daad4830",
  "CreatedAt": "2022-05-12T17:22:09Z",
  "Email": "jacques@cellaire.com",
  "FullName": "Jacques Cellaire",
  "Metadata": null
}
```

## Get a user token

```console
$ curl -L --request POST \
--header "Content-Type: application/json" \
--data '{"Email" : "jacques@cellaire.com","Password": "OopsjacquesHasBeenHacked"}' \
$BARIOT_HOST/users/login
```

Example response:

```json
{ "Token": "xxxxxx.xxxxxx(...)xxxxxxxx.xxxxxx" }
```

Save the received user token.

```console
$ export USER_TOKEN=xxxxxx.xxxx(...)xxxx.xxxx
```

## Create a thing

```console
$ curl -L --request POST \
--header "Content-Type: application/json" \
--header "Authorization: $USER_TOKEN" \
--data '{"Name": "smart-bottle-1", "Key": "123456789"}' \
$BARIOT_HOST/things/
```

Example response:

```json
{
  "Id": "9226093b-6b52-43d2-8345-b00e2a682a5d",
  "CreatedAt": "2022-05-12T18:46:43Z",
  "Key": "123456789",
  "Name": "smart-bottle-1",
  "UserId": "0ed50174-ba43-43f8-93f0-a971daad4830",
  "Metadata": null
}
```

Save the thing ID.

```console
$ export THING_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### Get a thing token

```console
$ curl -L --request GET \
--header "Content-Type: application/json" \
--header "Authorization: $USER_TOKEN" \
$BARIOT_HOST/things/$THING_ID/token
```

Example response:

```json
{ "jwt": "xxxxxx.xxxxxxx(...)xxxxxxx.xxxxxx" }
```

Save this token.

```console
$ export THING_TOKEN=xxxxxx.xxxxxxx(...)xxxxxxx.xxxxxx
```

---

<br>

## Alternative to Curl

Head to .../Bariot/support/http/

In there each .http file can be used with **[vscode-restclient](https://github.com/Huachao/vscode-restclient)** to accomplish the same functions as above in a more user-friendly fashion.

---

<br>

## Send data through MQTT

Head to .../Bariot/support/scripts/mqtt.

ATTENTION:
For the MQTT data sending script to work, the following environment variables need to be set:
BARIOT_HOST, THING_TOKEN, THING_ID
(If you followed this tutorial, they should be already set, otherwise, see above)

# OPTIONAL: If you are using PowerShell terminal

```console
Setting environment variables  with PowerShell:
$env:BARIOT_HOST = "xxxxxx";
# calling them
# $env:BARIOT_HOST
```

Open "thing_send_data_mqtt.go" with a text editor.
At the top, replace "JWT" and "THING_ID" with the Thing token and Thing ID obtained previously with curl.
Also, replace MQTT_HOST with the EC2 public DNS address ($BARIOT_HOST).
(example: ec2-xx-xx-xxx-xx.eu-west-1.compute.amazonaws.com)

then run

```console
$ go run thing_send_data_mqtt.go
```

This with send a MQTT with sensor data formatted with SENML and authenticated with a JWT token.

<br>

---

## Vizualize your data with Grafana

Head to $BARIOT_HOST/grafana

ec2-xx-xx-xxx-xx.eu-west-1.compute.amazonaws.com/grafana

```
login:    admin
password: admin
```

Change the password

Add a data source InfluxDB

| Option       | Value                   |
| ------------ | ----------------------- |
| Query        | language: Flux          |
| HTTP         | http://influxdb_db:8086 |
| Organization | bariot_org              |
| Token        | xxxxxxxxx               |
| Bucket       | bariot_bucket           |

The above values can be found in .../Bariot/docker/.env

You can then create a new dashboard.

And add the following query:

```
from(bucket: "bariot_bucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
```

You can then vizualize the data you sent through MQTT.

<br>

**Limit of the current system :**

Please note that the above query is for **ALL things data**, at the moment we cannot restrict the data on a per user or per thing basis.

This is the next step, since Grafana is not made for user authorization and data restriction (it has access to the whole influxdb DB).

In order to enable/restrict each user to access/visualize only his data, a custom HTTP endpoint and visualization tool (frontend) has to be made.

<br>

---

[puttyec2]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/putty.html
[docker]: https://medium.com/appgambit/part-1-running-docker-on-aws-ec2-cbcf0ec7c3f8
[docker-compose]: https://acloudxpert.com/how-to-install-docker-compose-on-amazon-linux-ami
[micro]: https://micro-editor.github.io
[bariot]: https://github.com/pseudoincorrect/Bariot
[openssl]: https://www.howtogeek.com/devops/how-to-create-and-use-self-signed-ssl-on-nginx
[randomkeygen]: https://randomkeygen.com
[installgo]: https://linguinecode.com/post/install-golang-linux-terminal
