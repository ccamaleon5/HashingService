# Hahing Service

This is a Server that hash any type of file, from applications to documents and media to verify the integrity. 

The Server sign the hash to guarantee that document was hashed by this Server 

## Prerequisites

* Go 1.12+ installation or later
* **GOPATH** environment variable is set correctly
* docker version 17.03 or later

## Package overview
#TODO

## Install

```
$ git clone https://github.com/ccamaleon5/CredentialMother

$ export GO111MODULE=on

$ cd CredentialMother
$ go build
```

## Run

```
$ credential-provider-server init [-x PASSWORD]
[PASSWORD] is your keystore password that will be created
$ credential-provider-server start --port=8000 --tlscertificate server.crt --tlskey server.key [-x PASSWORD]
```

where --port is a listen port http

You can try in localhost:8000/swagger-ui/

### Docker

* Clone this repository

```
$ git clone https://github.com/ccamaleon5/CredentialProvider
```

* Create a local directory that saves application data  

```
$ mkdir /CredentialData
```

* Copy the configuration file and swaggerui from repository to your local directory created above:

```
$ cp repo/CredentialProvider/credential-provider-server-config-yaml /CredentialData/
$ cp -r repo/CredentialProvider/swagger/swaggerui  /CredentialData/ 

```

* Now pull the docker image and run the container, setting your node identity and the folder location that will be the volume 

```
$ docker pull aparejaa/credentialprovider:1.0.0
$ docker run -dit -v {CredentialProvider_DIR}:/CredentialProvider -p 8000:8000 -p 8001:8001 aparejaa/credentialprovider:1.0.0 credential-provider-server init [-x PASSWORD]
$ docker run -dit -v {CredentialProvider_DIR}:/CredentialProvider -p 8000:8000 -p 8001:8001 aparejaa/credentialprovider:1.0.0
```

* The container will create KeyStore in your local volume

You can try in localhost:8000/swagger-ui/