import grpc from 'k6/experimental/grpc';
import { check, sleep } from 'k6';

// this file contains test to check if request throw errors, the return objects will not be analized


const client = new grpc.Client();
client.load(['../api'], 'main.proto');

export const options = {
    scenarios: {
        unary_call: {
            // name of the executor to use
            executor: 'shared-iterations',
            exec: 'unaryCall',
        },
        stream_call: {
            exec: 'streamCall',
            executor: 'shared-iterations',
        }
    },
};

export function unaryCall() {

    client.connect('localhost:4445', {plaintext: true});

    const data = { "param": 2 };
    const response = client.invoke('main.Service/Method', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    // print result
    if (response.status !== grpc.StatusOK) {

        console.log("status: " + response.status + " error: " + response.error.message);
    } else {
        console.log("status: " + response.status);
    }

    console.log(JSON.stringify(response.message));

    client.close()
    sleep(1)
}

export function streamCall() {

    let err = false;

    client.connect('localhost:4445', {plaintext: true});

    const data = { "empty": {}};
    const stream = new grpc.Stream(client, 'main.Service/Method');

    // print what received
    stream.on('data', (category) => {
        console.log('Received Category: ' + JSON.stringify(category));
    });

    // sets up a handler for the error event (an error occurs)
    stream.on('error', function (e) {
        // An error has occurred and the stream has been closed.
        console.log('Error: ' + JSON.stringify(e));
        err = true;
    });

    // close client when done
    stream.on('end', function () {
        // The server has finished sending
        client.close();
        console.log('All done');
    });

    // send request
    stream.write(data);

    check(err, {
        'status is OK': (r) => !r,
    });

    sleep(1);
}