# Instructions and installation of Bariot

### These instructions seems long, but most of them are about installing docker, git and go which most of us developper already have.
### The real deployment of Bariot is not so long/complex. 
### Better deployment automation is planned.

<br/>

---

<br/>

# Installation on AWS EC2 virtual machine

## EC2 VM Preparation

### Start an EC2 machine (T2 medium used here)

#### Change the security group of the VM to include:

| PROTOCOL | PORT | SOURCE |
| --- | --- | --- |
| SSH | 22 | 0.0.0.0/0 |
| HTTP | 80 | 0.0.0.0/0 |
| HTTPS | 443 | 0.0.0.0/0 |
| MQTT | 1883 | 0.0.0.0/0 |
| MQTTS | 8883 | 0.0.0.0/0 |

<br/>

### SSH to your VM (I use putty, [Tuto][PuttyEC2])

<br/>

### Install Docker ([Tuto][Docker])
``` console
  $ sudo yum update -y
  $ sudo amazon-linux-extras install docker
  $ sudo service docker start
  $ sudo usermod -a -G docker ec2-user
```

### Install Docker-compose ([Tuto][Docker-compose])
``` console
  $ cd
  $ sudo curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-`uname -s`-`uname -m` | sudo tee /usr/local/bin/docker-compose > /dev/null
  $ sudo chmod +x /usr/local/bin/docker-compose
  $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
  $ sudo reboot
```

### Optional: Install git and text editor, [Micro][Micro]
``` console
  $ cd && curl https://getmic.ro | bash && sudo mv micro /usr/local/bin/micro
  $ sudo yum instal git
```

### Install Go ([Tuto][InstallGo])
``` console
  $ cd
  $ wget https://go.dev/dl/go1.18.2.linux-amd64.tar.gz
  $ sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
  
  $ mkdir -p go/projects

  # export GO bins and GOPATH
  $ micro ~/.bash_profile
  # add "PATH=$PATH:/usr/local/go/bin" before "export PATH"
  # add "export GOPATH="/home/ec2-user/go" before  "export PATH"

  $ source ~/.bash_profile
```

## Bariot Deployment

### Clone Bariot [repository][Bariot]
``` console
  $ cd go/projects
  $ git clone https://github.com/pseudoincorrect/Bariot.git
```

### Create self-signed SSL certificate ([Tuto][OpenSSL])
``` console
  $ cd /home/ec2-user/go/projects/Bariot/docker/nginx
  $ mkdir ssl && cd ssl
  $ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout self-signed.key -out self-signed.crt
  $ openssl dhparam -out dhparam.pem 2048
  # the dhparam command takes few secs
````

### Change secrets (optional for quick testing)
``` console
  $ cd /home/ec2-user/Bariot/docker
  $ micro .env
  # modify the following secrets:
  # - THINGS_DB_PASSWORD
  # - USERS_DB_PASSWORD
  # - USER_ADMIN_PASSWORD
  # - JWT_SECRET
  # - MQTT_PASS
  # - INFLUXDB_PASSWORD
  # - INFLUXDB_TOKEN
  # You can use https://randomkeygen.com/
  # It is planned to use Vault to manage these secrets later on.
````

### Start the system
``` console
  $ cd /home/ec2-user/go/projects/Bariot/docker
  $ sudo service docker start
  $ docker-compose up
````

It will take quite some time the first time your run it, 10 mins on a T2 medium...

Bariot is downloading every docker image, go package and build all services.

## ! At the moment, timing errors can happen with the start of multiple containers !

Re-running "docker-compose up" from another terminal usually do the trick.

Roadmap: Deployment with precompiled images.

---

## Use the HTTP API with curl

### Create a user

### Create a thing

---

## Send data through MQTT

---

## Vizualize your data with Grafana

--- 

What's next:
  - CLI to replace Curl
  - Get thing data on a user basis

[PuttyEC2]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/putty.html
[Docker]: https://medium.com/appgambit/part-1-running-docker-on-aws-ec2-cbcf0ec7c3f8
[Docker-compose]: https://acloudxpert.com/how-to-install-docker-compose-on-amazon-linux-ami
[Micro]: https://micro-editor.github.io
[Bariot]: https://github.com/pseudoincorrect/Bariot
[OpenSSL]: https://www.howtogeek.com/devops/how-to-create-and-use-self-signed-ssl-on-nginx
[InstallGo]: https://linguinecode.com/post/install-golang-linux-terminal