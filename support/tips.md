# USEFUL COMMANDS

### To build a docker image
``` console
$ docker build -t <container_image_name> .
```

### To start and enter a container
``` console
$ docker run -it --rm --name <container_name> <container_image_name>
```

### To start and enter a container and replace initial command
``` console
$ docker run -it --rm --entrypoint=sh <container_image_name>
```

### To remove all volumes
``` console
$ docker volume rm $(docker volume ls -q)
```

### To rebuild a specific service
``` console
$ docker-compose up --no-deps --build <docker_compose_service_name>
```

### To shut down a single container
``` console
$ docker-compose rm -s -v <docker_compose_service_name>
```

### Volumes location in Win 11 WSL2 (windows explorer path):  
``` console
\\wsl$\docker-desktop-data\version-pack-data\community\docker\volumes\
```

### Install Micro (text editor) on a alpine container
``` console  
$ cd && curl https://getmic.ro | bash && sudo mv micro /usr/local/bin/micro
```

## Golang clean test cache (useful for MQTT test)
``` console  
$ go clean -testcache
```

## Golang import (local) package error
file structure : myGoProject/auth/
If your package name is the same as the go file:
  client.go       >>   package client
  you need to use `import "github.com/xxxx/xxx/myproject/auth/client"`
If your package name is different than the go file:
  auth_client.go  >>   package client
  you  need to use `import client "github.com/xxxx/xxx/myproject/auth"`


## add domain to local network on windows
first find you IP in a WSL (linux) terminal
``` console
# hostname -I
```
Copy the IP address outputed
In "C:\Windows\System32\drivers\etc\hosts" add: 
``` txt
# Personal conf
<outputed IP address> balancer.com
```

## Curl an insecure https endpoint (self-signed)
``` console
curl -k -v -L https://balancer.com/w2
# -L : (follow) following redirection (nginx reverse proxy)
# -k : (insecure) self-signed security issue ignoring
# -v : (verbose)
```

## RE-GENERATE PROTOBUFF
with a bash terminal, navigate to utilities/proto
``` console
$ ./generate.sg
```