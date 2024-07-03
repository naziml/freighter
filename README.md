### Clone custom registry code

```bash
git clone git@github.com:johnewart_microsoft/go-containerregistry.git
cd go-containerregistry
git co manifest-store
```

### Build and run freighter server

```bash 
cd cmd/freighter_server
go build 
./freighter_server 
```

### Push an image to the registry 

Note: You may have to use your machine's IP address instead of localhost if docker isn't running locally 

```bash
docker pull influxdb:1.7.4
docker tag influxdb:1.7.4 localhost:1338/influxdb:1.7.4
docker push localhost:1338/influxdb:1.7.4
```

### Build and run freighter client

```bash
cd cmd/freighter_client
go build
mkdir ./containerfs
./freighter_client -repo influxdb -target 1.7.4 ./containerfs
``` 






