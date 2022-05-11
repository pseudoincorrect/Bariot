

## Generate the TLS certificates
``` console
# With a linux terminal (windows WSL in my case)
$ cd docker/nginx/ssl/
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout self-signed.key -out self-signed.crt
$ openssl dhparam -out dhparam.pem 2048
```

## Start the system
``` console
$ cd docker
$ docker-compose up
````
