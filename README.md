# Beacon
Software for exposing an open interface for Lighthouse on cloud providers.
![beacon](http://i.imgur.com/OGIcqyq.png)

## Download
```
go get github.com/lighthouse/beacon
```

## Build
```
go install github.com/lighthouse/beacon
```

## Run
```
$GOPATH/bin/beacon
```

### Arguments
* ```-driver gce``` Driver to use when interfacing with the vm provider.
* ```-h 0.0.0.0:5000``` Address to listen on when hosting the server.
* ```-key server.key``` Path to private key used for hosting TLS connections.
* ```-pem server.pem``` Path to Cert used for hosting TLS connections.
* ```-token 123abc``` Authentication token used to grant access to the beacon api.



## Docker Build
```
docker build -t beacon $GOPATH/src/github.com/lighthouse/beacon
```

## Docker Run
```
docker run -d -p 5000:5000 reg.rob-sheehy.com/beacon -h 0.0.0.0:5000
```

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
