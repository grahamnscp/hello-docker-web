# Setup Docker EE Layer 7 Routing

Docker EE Docs here:  
https://docs.docker.com/ee/ucp/interlock/  
https://docs.docker.com/ee/ucp/interlock/architecture/

First enable it in the Docker EE UCP Admin web UI under:  
`Admin Settings -> Layer 7 Routing -> check: Enable Layer 7 Routing -> Save`

The default configuration interlock services are then spun up:
(Using Docker EE client bundle to securely access remote cluster)
```
$ source env.sh
$ docker service ls

$ docker service ls --filter name=ucp-interlock
ID                  NAME                      MODE                REPLICAS            IMAGE                                  PORTS
m98e5zdd73md        ucp-interlock             replicated          1/1                 docker/ucp-interlock:3.0.2
3sf1cox19ut6        ucp-interlock-extension   replicated          1/1                 docker/ucp-interlock-extension:3.0.2
kzkte92z6gn6        ucp-interlock-proxy       replicated          2/2                 docker/ucp-interlock-proxy:3.0.2       *:8080->80/tcp, *:8443->443/tcp
```

Move the interlock proxy service instances to dedicated (/predictable) cluster nodes.

Firstly label the nodes:  
https://docs.docker.com/ee/ucp/interlock/deploy/production/#apply-labels-to-nodes  
```
docker node update --label-add nodetype=loadbalancer ucp-lb1.example.org
docker node update --label-add nodetype=loadbalancer ucp-lb2.example.org
```
Then reconfigure the interlock service to 'constrain' the proxies to the nodes labeled with the `nodetype=loadbalancer` label:

Interlock config docs:  
https://docs.docker.com/ee/ucp/interlock/deploy/configure/  
```
$ CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0).ConfigName }}' ucp-interlock)

$ docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > config.toml
```
Modify the config file, config.toml to add the new ProxyConstraint:
```
# Edit, find and hange line from:
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux"]

# to:
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux", "node.labels.nodetype==loadbalancer"]

# save modified file
```
Note: full config options reference doc:  
https://docs.docker.com/ee/ucp/interlock/deploy/configuration-reference/

Create a new Docker configuration object from the file you’ve edited:  
```
$ NEW_CONFIG_NAME="com.docker.ucp.interlock.conf-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"

$ docker config create $NEW_CONFIG_NAME config.toml
```

Update the ucp-interlock service to start using the new configuration:  
```
$ docker service update \
  --config-rm $CURRENT_CONFIG_NAME \
  --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
  ucp-interlock
```
Note: Need to deploy a service for the config to take effect and proxies to move to the correct nodes..  


Check where proxy running:
```
$ docker service ps ucp-interlock-proxy | egrep 'Running|ID'
ID                  NAME                        IMAGE                              NODE                           DESIRED STATE       CURRENT STATE           ERROR                              PORTS
75lrdrslzutt        ucp-interlock-proxy.1       docker/ucp-interlock-proxy:3.0.2   ucp-lb2.example.org            Running             Running 5 minutess ago
uk1y5e0jpnj7        ucp-interlock-proxy.2       docker/ucp-interlock-proxy:3.0.2   ucp-lb1.example.org            Running             Running 18 minutess ago
```
Rollback if there is a config error and the service doesn't restart:  
```
$ docker service update \
  --update-failure-action rollback \
  ucp-interlock
```
