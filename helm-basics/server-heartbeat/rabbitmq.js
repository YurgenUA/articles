import amqp from "amqplib";

const MQ_URL = process.env.MQ_URL || 'amqp://guest:guest@localhost:5672/';
const QUEUE = 'queue#1';

export async function sendHeartbeat() {

  try {
    const connection = await amqp.connect(MQ_URL);
    const channel = await connection.createChannel();
    await channel.assertQueue(QUEUE);

    console.log('Sending new pulse');
    channel.sendToQueue(QUEUE, Buffer.from(`Pulse at '${new Date().toISOString()}'`),
      { expiration: '30000' }
    );

    if (channel) await channel.close();
    if (connection) await connection.close();

  } catch (error) {
    console.error(error);
  }
};

export async function getHeartbeats() {
  const result = [];
  try {
    const connection = await amqp.connect(MQ_URL);
    const channel = await connection.createChannel();
    await channel.assertQueue(QUEUE);

    await channel.consume(QUEUE, (msg) => {
      result.push(msg.content.toString());
      channel.ack(msg);
    })

    console.log('Releasing RabbitMQ objects');
    if (channel) await channel.close();
    if (connection) await connection.close();
    return result;

  } catch (error) {
    console.error(error);
  }
};
