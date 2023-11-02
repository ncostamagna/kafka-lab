const express = require("express");
const app = express();
const port = 3000;

const ip = require("ip");

const { Kafka, CompressionTypes, logLevel } = require("kafkajs");

const host = ip.address();

const kafka = new Kafka({
  logLevel: logLevel.DEBUG,
  brokers: [`${host}:9092`],
  clientId: "server1",
});

const producer = kafka.producer();

const sendMessage = () => {
  return producer
    .send({
      topic: "topic-test-server",
      compression: CompressionTypes.GZIP,
      messages: [
        {
          key: `key-123`,
          value: `value-123`,
        },
      ],
    })
    .then(console.log)
    .catch((e) => console.error(`[example/producer] ${e.message}`, e));
};

const run = async () => {
  await producer.connect();
};

run().catch((e) => console.error(`[example/producer] ${e.message}`, e));

app.get("/", (req, res) => {
  res.send("Hello World!");
});

app.get("/send", (req, res) => {
  sendMessage();
  res.send("sent to server2");
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});

const errorTypes = ["unhandledRejection", "uncaughtException"];
const signalTraps = ["SIGTERM", "SIGINT", "SIGUSR2"];

errorTypes.forEach((type) => {
  process.on(type, async () => {
    try {
      console.log(`process.on ${type}`);
      await producer.disconnect();
      process.exit(0);
    } catch (_) {
      process.exit(1);
    }
  });
});

signalTraps.forEach((type) => {
  process.once(type, async () => {
    try {
      await producer.disconnect();
    } finally {
      process.kill(process.pid, type);
    }
  });
});
