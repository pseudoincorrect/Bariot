# UNDERSTANDING BARIOT'S CODE

## Where to first look at

### Complete the instruction.md tutorial

### Have a look at the build/docker/ folder and the docker-compose file to a quick glance on the architecture

### Then to every main.go files, as each microservice starts within the cmd/ folder

### Readme.md files will make you understand the folder architecture

### Look at the ...\_tests.go files

---

## Modify the and running the code

Let's say that you ran the following command and all went well

```console
$ cd bariot

Note: Tests will fail if bariot is running (with "docker-compose up")

$ go test ./...
$ cd build/docker
$ docker-compose up -d
$ docker-compose logs -f --tail 5

On another terminal, run:

$ cd bariot/test/end_to_end
$ ./venv/Scripts/activate
$(venv) python complete_test.py
```

Now let's make some modifs to (for instance) pkg/things/client/client.go in order to change the behavior of the _reader_ microservice.

_make some modif to the code_

To stop, rebuild and redeploy (locally) the code of _reader_ microservice, let's run:

```console
$ cd bariot
$ docker-compose rm -s -v -f reader
$ docker-compose build reader
$ docker-compose up -d
$ docker-compose logs -f --tail 5
```

---

## Need some help with the commands and others ?

Have a look at bariot/support/tips.md and gotchas.md
