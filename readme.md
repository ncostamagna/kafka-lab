# Bases Kafka JS

## Setup

```
cd bases-kafka-js
npm i
```

You must set the host ip environment with your IP address

```
export HOST_IP=192.X.X.X
```

affter that, you should run the docker compose

```
docker compose  up -d
```

## Run the program

Open 2 console:

- in the first console run the consumer.js

```
node consumer.js
```

- in the other one, run the producer

```
node producer.js
```
