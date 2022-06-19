# ISSUES ENCOUNTERED DURING THE DEVELOPMENT

## MQTT: Message not received on localhost

**SYMPTOM**
When sending a new MQTT message on localhost, nothing is received on the server
CAUSES
Check whether you have another MQTT Broker running on your machine.
For instance if you use MQTTX (UI), it will create a Mosquitto broker
and intercept any MQTT messages.
**SOLUTION**
Terminate the extra broker task.

## NGINX: "http" directive is not allowed here in

**SYMPTOM**
When launching nginx the following error: "http" directive is not allowed here in
**CAUSE**
The http{} block is located within another http block (error)
/etc/nginx/nginx.conf will call other conf file located in /etc/nginx/conf.d and
include them within its http{} block.
**SOLUTION**
Replace the /etc/nginx/nginx.conf with your own one.
You can manage the includes in it.

## NGINX: Reverse-proxying influxdb

**SYMPTOM & CAUSE**
Reverse-proxying influxdb UI is not really possible at the moment.
ISSUE: https://github.com/influxdata/influxdb/issues/15721
**SOLUTION**
Please use docker compose ports to access it if needed.

## PYTHON: HTTP Response 400, "EOF"

SYMPTOM
Golang cannot decode the JSON.
**CAUSE**
Formatting of the URL
**SOLUTION**
Do not forget the "/" at the end of the URL example:
http://bariot.com/users/ and NOT http://bariot.com/users

## PYTHON: Request and readline, response 500

**SYMPTOM**
Golang cannot decode the JSON.
**CAUSE**
A trailing NewLine char when using readline
**SOLUTION**
Trim the trailing "new line" char

## GOLANG: Interface not implemented, error: compilerInvalidIfaceAssign

**SYMPTOM**
a Struct is not implementing an interface even with the same function
**CAUSE**
functions are not exported
**SOLUTION**
in the interface definition add a Uppercase to the functions

## GOLANG: Interface not implemented, error: method has pointer receiver

**SYMPTOM**
build fail because of a failed interface implementation
**CAUSE**
There is a pointer issue somewhere, but no indication where
**SOLUTION**
Try dereferencing (&service) when using the interface instance

## GOLANG: HTTP Error EOF with json.NewDecoder(req.Body)

**SYMPTOM**
json.NewDecoder(req.Body) ... give an EOF error
**CAUSE**
ioutil.ReadAll read the whole req.Body reader, thus when we try to
read it again with json.NewDecoder, nothing remain, hence EOF error
**SOLUTION**
set req.Body again with the content you just read with ioutil.ReadAll
LINK
https://stackoverflow.com/questions/49745252/reverseproxy-depending-on-the-request-body-in-golang

## TEMPLATE:

**SYMPTOM**

**CAUSE**

**SOLUTION**
