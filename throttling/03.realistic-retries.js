const readline = require('readline');
const AWS = require('aws-sdk');

let success = 0;
let retries = 0;
let failures = 0;

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

const dynamodb = new AWS.DynamoDB({
  region: 'eu-west-1',
  maxRetries: 5,
  retryDelayOptions: {
    base: 300
  }
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
      return data;
    }).catch(function (err) {
      failures++;
    });
}

Main = function () {
  rl.question('Enter db get numbers:', async (numberOfGets) => {
    console.log(`Testing with ${numberOfGets} simultaneous gets`);
    rl.close();

    console.time('full run');
    const promises = new Array();
    promises.length = numberOfGets;
    await Promise.all(promises.fill(undefined).map(x => ReadDynamoPromise()));

    console.log('=============== Summary =================');
    console.log('success:', success);
    console.log('retries:', retries);
    console.log('failures:', failures);
    console.log('=========================================');
    console.timeEnd('full run');
  });
}();