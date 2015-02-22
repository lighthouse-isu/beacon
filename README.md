# Beacon
Software for exposing an open interface for Lighthouse on cloud providers.
![beacon](http://i.imgur.com/OGIcqyq.png)

## Quick Setup With Boot2Docker
### Docker Host
Ensure ```DOCKER_HOST``` env var is set, call ```echo $DOCKER_HOST``` to double check the address.
### Start Boot2Docker
```
boot2docker up
```
### Pull beacon
```
docker pull lighthouse/beacon:latest
```
### Run beacon
```
docker run -t -i -e "DOCKER_HOST=tcp://your.docker.host:2375" -p 5000:5000 \
    lighthouse/beacon:latest -token foobar -h 0.0.0.0:5000 -driver local
```
### Try it out!
* Check current driver ```curl your.docker.host:5000/which -H "Token: foobar"```
* Check available vms ```curl your.docker.host:5000/vms -H "Token: foobar"```

## Running With Go
### Download
```
go get github.com/lighthouse/beacon
```
### Build
```
go install github.com/lighthouse/beacon
```
### Run
```
$GOPATH/bin/beacon
```
### Test
```
go test github.com/lighthouse/beacon/...
```

## Running With Docker
### Download
```
go get github.com/lighthouse/beacon
```
### Build
```
docker build -t beacon $GOPATH/src/github.com/lighthouse/beacon
```
### Run
```
docker run -d -p 5000:5000 reg.rob-sheehy.com/beacon -h 0.0.0.0:5000
```

## Arguments
* ```-driver gce``` Driver to use when interfacing with the vm provider.
* ```-h 0.0.0.0:5000``` Address to listen on when hosting the server.
* ```-key server.key``` Path to private key used for hosting TLS connections.
* ```-pem server.pem``` Path to Cert used for hosting TLS connections.
* ```-token 123abc``` Authentication token used to grant access to the beacon api.


## Authentication
To make successful api calls to beacon from a client you will need the generated auth Token which is logged on app startup.
That token must be in the header of each request as "Token" to preform any api call. Otherwise you will be greeted with a 401 status code.

## Drivers

### Local
Interfaces with [boot2docker](http://boot2docker.io/) and uses the address stored in ```$DOCKER_HOST``` to make requests. This also works if your running inside Docker such that you can use..
```
-e "DOCKER_HOST=tcp://your.docker.i.p:2375"
```
to manually declair the host of the Docker daemon. For example...
```
docker run -d -p 5000:5000 -e "DOCKER_HOST=tcp://192.168.59.103:2375" beacon -h 0.0.0.0:5000 -driver local
```

### GCE
Interfaces with [Compute Engine](https://cloud.google.com/compute/) and requires that the hosting vm has "Compute" read/write privalages to the Project to detect existing vms.

### Digital Ocean
Interfaces with [Digital Ocean](https://www.digitalocean.com/) and requires that an api token be set to the hosting vm's "User Data".  Read more info [here](https://github.com/lighthouse/beacon/pull/4).

### Config
Instead of relying on a provider api you can manually create a config file that list available ips of vms you want beacon to communicate with.
For example all you have to do is drop a ```config.json``` into the running directory of Beacon and it will take care of the rest for you.
A simple  ```config.json``` can look something like this.
```JSON
[
    "192.168.59.103",
    "127.0.0.1"
]

```
