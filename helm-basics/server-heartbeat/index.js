import { sendHeartbeat, getHeartbeats } from "./rabbitmq";
import restify from 'restify';

const SERVER_PORT = process.env.SERVER_PORT || '3000';
const server = restify.createServer();

server.use(restify.plugins.acceptParser(server.acceptable));

server.get('/heartbeats', async (req, res, next) => {
  res.send(await getHeartbeats());
});

server.get('*', (req, res, next) => {
  res.send('Default route handler invoked. Please be more specific.');
  return next();
});

server.listen(Number.parseInt(SERVER_PORT), function () {
  console.log('Server is listening at %s', server.url);
});

setInterval(async () => { await sendHeartbeat() }, 1000);
