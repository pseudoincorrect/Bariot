## MQTT message not received on localhost
  CAUSES: 
    Check whether you have another MQTT Broker running on your machine.
    For instance if you use MQTTX (UI), it will create a Mosquitto broker
    and intercept any MQTT messages. 
  SOLUTIONS: 
    Terminate the task

## NGINX -> "http" directive is not allowed here in
  CAUSE:
    the http{} block is located within another http block (error)
    /etc/nginx/nginx.conf will call other conf file located in /etc/nginx/conf.d and
    include them within its http{} block
  SOLUTION
    replace the /etc/nginx/nginx.conf with your own one
    you can manage the includes there

## NGINX Reverse-proxing influxdb
  CAUSE:
    Reverse-proxing influxdb UI is not really possible at the moment
    ISSUE: https://github.com/influxdata/influxdb/issues/15721
  SOLUTION:
    Please use docker compose ports to access it if needed

## HTTP response 400 : EOF
  CAUSE:
    Golang cannot decode the JSON
  SOLUTION:
    do not forget the "/" at the end of the URL example:
    http://bariot.com/users/   and NOT   http://bariot.com/users
