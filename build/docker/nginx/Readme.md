# Generate SSL keys and certs

``` console
$ cd docker/nginx/ssl/

$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout self-signed.key -out self-signed.crt

$ openssl dhparam -out dhparam.pem 1024
````