# Simple image with sample Docker EE deployment files

Build the image using multi-stage build Dockerfile:

```
$ docker build -t helloweb:build .  

Sending build context to Docker daemon   72.7kB  
Step 1/8 : FROM golang:1.8-alpine AS compile  
 ---> 4cb86d3661bf  
Step 2/8 : COPY hello-docker-web.go /go  
 ---> 7e51f879e507  
Step 3/8 : RUN go build hello-docker-web.go
 ---> Running in fd3a82eccf4
Removing intermediate container fd3a82eccf42
 ---> 5befd2df76ac
Step 4/8 : FROM alpine:latest
 ---> 3fd9065eaf02
Step 5/8 : COPY --from=compile /go/hello-docker-web /
 ---> 9e518e40c186
Step 6/8 : USER nobody:nobody
 ---> Running in 7d9c7fe14107
Removing intermediate container 7d9c7fe14107
 ---> 750dc540fdae
Step 7/8 : EXPOSE 8080
 ---> Running in d689db3aa3b3
Removing intermediate container d689db3aa3b3
 ---> 395d6cefc5ae
Step 8/8 : ENTRYPOINT ["/hello-docker-web"]
 ---> Running in bcb713719813
Removing intermediate container bcb713719813
 ---> a59cc0402098
Successfully built a59cc0402098
Successfully tagged helloweb:build
```

Tag as required:
```
$ docker tag helloweb:build grahamh/hello-docker-web:3.0
```

Start a local container instance to test:
```
$ docker run -d --rm -p 8080:8080 grahamh/hello-docker-web:3.0
31b349d27ab86d0f3dde80a7c9f73e75f7250c153b3110b207f0e1bd003e3747
```

Test:
```
$ curl localhost:8080

                              ##
                        ## ## ##        ==
                     ## ## ## ## ##    ===
                 /`````````````````\___/ ===
            ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~~/~ === ~~~
                 \______ o           __/
                   \    \         __/
                    \____\_______/
 _           _    _                _            _
| |     ___ | |  | |    ___     __| | ___   ___| | _____ _ __
| |___ / _ \| |  | |   / _ \   / _  |/ _ \ / __| |/ / _ \ '__|
|  _  |  __/| |__| |__| (_) | | (_| | (_) | (__|   <  __/ |
|_| |_|\___/ \___|\___|\___/   \____|\___/ \___|_|\_\___|_|
```

There are some uris in the image too to test layer 7 context routing:
```
$ curl localhost:8080/hello/ignoredtextafteruri
 _           _    _
| |     ___ | |  | |    ___
| |___ / _ \| |  | |   / _ \
|  _  |  __/| |__| |__| (_) |
|_| |_|\___/ \___|\___|\___/
```

# Deploy as a service on Docker EE SWARM orchestrator:  
(Using Docker EE client bundle to securely access remote cluster)
```
$ source env.sh
$ docker service ls

kngptcjtwqay        ucp-agent                 global              11/11               docker/ucp-agent:3.0.2
njs1theeo0sn        ucp-agent-s390x           global              0/0                 docker/ucp-agent-s390x:3.0.2
j9kmxzsr3r8k        ucp-agent-win             global              0/0                 docker/ucp-agent-win:3.0.2
m98e5zdd73md        ucp-interlock             replicated          1/1                 docker/ucp-interlock:3.0.2
3sf1cox19ut6        ucp-interlock-extension   replicated          1/1                 docker/ucp-interlock-extension:3.0.2
kzkte92z6gn6        ucp-interlock-proxy       replicated          2/2                 docker/ucp-interlock-proxy:3.0.2       *:8080->80/tcp, *:8443->443/tcp


$ docker network create -d overlay helloweb-net

$ docker service create \
  --name interlock-helloweb \
  --network helloweb-net \
  --label com.docker.lb.hosts=helloweb.apps.example.org \
  --label com.docker.lb.network=helloweb-net \
  --label com.docker.lb.port=8080 \
  grahamh/hello-docker-web:3.0

$ docker service ls

qop5vxt7n4i0        interlock-helloweb        replicated          1/1                 grahamh/hello-docker-web:3.0
kngptcjtwqay        ucp-agent                 global              11/11               docker/ucp-agent:3.0.2
njs1theeo0sn        ucp-agent-s390x           global              0/0                 docker/ucp-agent-s390x:3.0.2
j9kmxzsr3r8k        ucp-agent-win             global              0/0                 docker/ucp-agent-win:3.0.2
m98e5zdd73md        ucp-interlock             replicated          1/1                 docker/ucp-interlock:3.0.2
3sf1cox19ut6        ucp-interlock-extension   replicated          1/1                 docker/ucp-interlock-extension:3.0.2
kzkte92z6gn6        ucp-interlock-proxy       replicated          2/2                 docker/ucp-interlock-proxy:3.0.2       *:8080->80/tcp, *:8443->443/tcp

$ docker service ps interlock-helloweb

ID           NAME                 IMAGE                        NODE
   DESIRED STATE  CURRENT STATE          ERROR               PORTS
qop5vxt7n4i0 interlock-helloweb.1 grahamh/hello-docker-web:3.0 ucp-worker1.example.org
   Running        Running 1 minute ago
```
With Interlock 2.0 setup on Docker EE, with a load balancer pointing to the ucp-interlock-proxy worker nodes that has a wildcard DNS entry of `*.apps.example.org`:  
(note port 8080 is used below as interlock is setup to listen on it, not that the app port is 8080 also)

```
$ curl helloweb.apps.example.org:8080

                              ##
                        ## ## ##        ==
                     ## ## ## ## ##    ===
                 /`````````````````\___/ ===
            ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~~/~ === ~~~
                 \______ o           __/
                   \    \         __/
                    \____\_______/
 _           _    _                _            _
| |     ___ | |  | |    ___     __| | ___   ___| | _____ _ __
| |___ / _ \| |  | |   / _ \   / _  |/ _ \ / __| |/ / _ \ '__|
|  _  |  __/| |__| |__| (_) | | (_| | (_) | (__|   <  __/ |
|_| |_|\___/ \___|\___|\___/   \____|\___/ \___|_|\_\___|_|
```

With app uri context based routing, testing with /hello uri:
```
$ docker network create -d overlay hello-net

$ docker service create \
  --name interlock-hello \
  --network hello-net \
  --label com.docker.lb.hosts=hello.apps.example.org \
  --label com.docker.lb.network=hello-net \
  --label com.docker.lb.port=8080 \
  --label com.docker.lb.context_root=/hello \
  --label com.docker.lb.context_root_rewrite=true \
  grahamh/hello-docker-web:3.0

$ curl hello.apps.example.org:8080/hello/
 _           _    _
| |     ___ | |  | |    ___
| |___ / _ \| |  | |   / _ \
|  _  |  __/| |__| |__| (_) |
|_| |_|\___/ \___|\___|\___/
```
