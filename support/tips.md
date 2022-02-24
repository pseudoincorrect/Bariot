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

### Volumes location in Win 11 WSL2 (windows explorer path):  
``` console
\\wsl$\docker-desktop-data\version-pack-data\community\docker\volumes\
```

### Install Micro (text editor) on a alpine container
``` console  
$ cd && curl https://getmic.ro | bash && sudo mv micro /usr/local/bin/micro
```