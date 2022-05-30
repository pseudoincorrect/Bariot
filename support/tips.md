## USEFUL COMMANDS AND TIPS

## DOCKER: To build a docker image

```console
$ docker build -t <container_image_name> .
```

## DOCKER: To start and enter a container

```console
$ docker run -it --rm --name <container_name> <container_image_name>
```

## DOCKER: To start and enter a container and replace initial command

```console
$ docker run -it --rm --entrypoint=sh <container_image_name>
```

## DOCKER: To start and enter a container with docker-compose

```console
$ docker-compose run --entrypoint=sh  <container_image_name>
```

## DOCKER: To remove all volumes

```console
$ docker volume rm $(docker volume ls -q)
```

## DOCKER: To rebuild a specific service

```console
$ docker-compose up --build <docker_compose_service_name>
```

## DOCKER: To rebuild a specific service completely

```console
$ docker-compose up --no-deps --build <docker_compose_service_name>
```

## DOCKER: To shut down a single container

```console
$ docker-compose rm -s -v <docker_compose_service_name>
```

## DOCKER: Volumes location in Win 11 WSL2 (windows explorer path):

```console
\\wsl$\docker-desktop-data\version-pack-data\community\docker\volumes\
```

## MICRO: Install Micro (text editor) on a alpine container

```console
$ cd && curl https://getmic.ro | bash && sudo mv micro /usr/local/bin/micro
```

## GOLANG: Clean test cache (useful for MQTT test)

```console
$ go clean -testcache
```

## GOLANG: Import (local) package error

file structure : myGoProject/auth/
If your package name is the same as the go file:
client.go >> package client
you need to use `import "github.com/xxxx/xxx/myproject/auth/client"`
If your package name is different than the go file:
auth_client.go >> package client
you need to use `import client "github.com/xxxx/xxx/myproject/auth"`

## NETWORK: Add domain to local network on windows

first find you IP in a WSL (linux) terminal

```console
# hostname -I
```

Copy the IP address outputted
In "C:\Windows\System32\drivers\etc\hosts" add:

```txt
# Personal conf
<outputted IP address> balancer.com
```

## CURL: Reach an insecure https endpoint (self-signed)

```console
curl -k -v -L https://balancer.com/w2
# -L : (follow) following redirection (nginx reverse proxy)
# -k : (insecure) self-signed security issue ignoring
# -v : (verbose)
```

## PROTOBUFF: Regenerate

with a bash terminal, navigate to utilities/proto

```console
$ ./generate.sg
```

## INFLUXDB: delete data from bucket

```console
$ docker exec -it <container_ID> bash
$ influx delete --bucket bariot_bucket --start 2021-01-01T00:00:00Z --stop 2023-01-01T00:00:00Z
```

## GIT: Create a feature branch with code not committed already

You were coding a new feature and then "Oh flute" I forgot to create a new branch

```console
$ git stash
$ git checkout main
$ git checkout -b feature-something
$ git stash apply
$ git commit -m "your message about something"
$ git push -u origin feature-something
```

## GIT: Integrate Main's changes into your feature branch

```console
$ git checkout main
$ git pull
$ git checkout feature-something
$ git rebase main
```

## GIT: Squash your feature branch's commits and push to origin

```console
$ git checkout feature-something
$ git rebase -i HEAD~20  # Squashing up to 20 commit (before pushing to origin)
$ git log                # Check what your commit looks like
$ git commit --amend
$ git push -f
```

if you have already pushed the commits you want to squash
https://stackoverflow.com/questions/5667884/how-to-squash-commits-in-git-after-they-have-been-pushed

```console
git rebase -i origin/feature-something~20 feature-something
 git push --force origin feature-something
```

## PYTHON VENV: create and setup

```console
$ cd ./tests/end_to_end
$ venv -m venv venv
$ ./venv/Scripts/activate
$ pip install -r requirements.txt
```

## POWERSHELL: Setting/getting environment variables

```console
Setting environment variables BARIOT_HOST, THING_TOKEN, THING_ID with PowerShell:
$env:BARIOT_HOST = "xxxxxx";
# calling them
$env:BARIOT_HOST
```
