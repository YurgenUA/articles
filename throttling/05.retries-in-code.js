const readline = require('readline');
const AWS = require('aws-sdk');
const PromiseWhile = require('./05.promise-while');

let success = 0;
let retries = 0;
let failures = 0;

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

const dynamodb = new AWS.DynamoDB({
  region: 'eu-west-1',
  maxRetries: 0
});
const docClient = new AWS.DynamoDB.DocumentClient({
  service: dynamodb
});

AWS.events.on('retry', function (resp) {
  retries++;
});

function ReadDynamoPromise() {
  const params = {
    TableName: 'table-sandbox',
    Key: {
      id: '5c6a8fa0-6014-11e8-b42b-d5ae7edca68c'
    }
  };
  return docClient.get(params).promise()
    .then(function (data) {
      success++;
      if (success % 100 == 0) {
        console.timeEnd('portion run');
        console.time('portion run');
      }
      return data;
    }).catch(function (err) {
      retries++;
      throw err;
    });
}

function InvokeProvisioned(func) {
  const pw = new PromiseWhile({
    timeoutBetweenCalls: 1000
  });
  const act = function () {
    return new Promise((resolve, reject) => {
      const pwThis = this;
      return func()
        .then((data) => {
          pwThis.data = data;
          return resolve(data);
        })
        .catch((err) => {
          pwThis.data = undefined;
          // If it's throttling => let's try more
          if (err.code === 'ProvisionedThroughputExceededException') {
            return resolve();
          }

          return reject(err);
        });
    });
  };

  const condition = function () {
    if (this.data) {
      return false;
    }
    return true;
  };
  return pw.Await(condition, act);
}

Main = function () {
  rl.question('Enter db get numbers:', async (numberOfGets) => {
    console.log(`Testing with ${numberOfGets} simultaneous gets`);
    rl.close();

    console.time('full run');
    console.time('portion run');
    const promises = new Array();
    promises.length = numberOfGets;
    await Promise.all(promises.fill(undefined).map(x => {
      return InvokeProvisioned(ReadDynamoPromise);
    }));

    console.log('=============== Summary =================');
    console.log('success:', success);
    console.log('retries:', retries);
    console.log('failures:', failures);
    console.log('=========================================');
    console.timeEnd('full run');
  });
}();