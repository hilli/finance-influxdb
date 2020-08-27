# finance-influxdb

Collects finance statistics from Yahoo Finance and emits the data to InfluxDB.

*THIS IS WORK IN PROGRESS!!!*

# Configuring

## Setup ticker symbols in the environment

```bash
docker run --rm -e SYMBOLS=AAPL,TSLA,MSFT hilli/finance-influxdb
```
or with visual output in the log as well while collecting data at every 30s sending it to another server:

```bash
docker run --rm -e DEBUG=true -e COLLECTION_INTERVAL=30 -e INFLUXDB_ENDPOINT=http://127.0.0.1:8086 -e SYMBOLS=AAPL,TSLA,MSFT hilli/finance-influxdb
```

If you prefer, you can create a `.env` file with the same environment as the docker-compose file and map it to the container as a volume.

## InfluxDB

You will need a `InfluxDB` server to accept your data. Set `INFLUXDB_ENTRYPOINT` to point to your server of choice (URL format).

# Docker Compose

`docker-compose.yaml` setting the environment:
```yaml
version: '3.7'
services:
  finance-influxdb:
    image: hilli/finance-influxdb
    container_name: finance-influxdb
    hostname: finance-influxdb
    restart: unless-stopped
    environment:
      INFLUXDB_ENDPOINT: "http://influxdb:8086"
      COLLECTION_INTERVAL: 60 # Seconds
      SYMBOLS: "AAPL,TSLA,MSFT" # Choose your symbols on https://finance.yahoo.com/
      #INFLUX_USER: "myuser" # If needed
      #INFLUX_PASSWORD: "boohoo" # If needed
      #DEBUG: "jearh" # Print the collected results in a human readable format to the docker log
```

`docker-compose.yaml` setting the environment in an `.env` file:
```yaml
version: '3.7'
services:
  finance-influxdb:
    image: hilli/finance-influxdb
    container_name: finance-influxdb
    hostname: finance-influxdb
    restart: unless-stopped
    volumes:
      - "./env:/.env"
```

`.env`:
```bash
INFLUX_ENDPOINT=http://localhost:8086
#INFLUX_USER=myuser # If needed
#INFLUX_PASSWORD=boohoo # If needed
COLLECTION_INTERVAL=15
SYMBOLS=AAPL,TSLA,MSFT,DANSKE.CO
#DEBUG=OK # Output data to stdout as well
```

# License etc

Licensed under MIT License. Source available at https://github.com/hilli/finance-influxdb
