# ISSUES ENCOUNTERED DURING THE DEVELOPMENT

## MQTT: Message not received on localhost
  CAUSES
    Check whether you have another MQTT Broker running on your machine.
    For instance if you use MQTTX (UI), it will create a Mosquitto broker
    and intercept any MQTT messages. 
  SOLUTIONS
    Terminate the extra broker task.

## NGINX: "http" directive is not allowed here in
  CAUSE
    The http{} block is located within another http block (error)
    /etc/nginx/nginx.conf will call other conf file located in /etc/nginx/conf.d and
    include them within its http{} block.
  SOLUTION
    Replace the /etc/nginx/nginx.conf with your own one.
    You can manage the includes in it.

## NGINX: Reverse-proxing influxdb
  CAUSE
    Reverse-proxing influxdb UI is not really possible at the moment.
    ISSUE: https://github.com/influxdata/influxdb/issues/15721
  SOLUTION
    Please use docker compose ports to access it if needed.

## HTTP: Response 400, "EOF"
  CAUSE
    Golang cannot decode the JSON.
  SOLUTION
    Do not forget the "/" at the end of the URL example:
    http://bariot.com/users/   and NOT   http://bariot.com/users
