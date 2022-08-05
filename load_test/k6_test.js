import {check, sleep} from 'k6';
import grpc from 'k6/net/grpc';
import encoding from 'k6/encoding';

// host.docker.internal:3606
// docker-compose -f k6-docker-compose.yml run k6 run - <k6_test.js
const client = new grpc.Client();
client.load([], 'api/proto/broker.proto');

const BASE_URL = '192.168.100.35:3606';

export const options = {
    scenarios: {
        publish: {
            executor: 'constant-vus',
            exec: 'testPublish',
            vus: 30,
            duration: '5m',
        },
        // fetch: {
        //     executor: 'constant-vus',
        //     exec: 'testFetch',
        //     vus: 700,
        //     duration: '1m',
        // },
    },
    // thresholds: {
    //     grpc_req_duration: ['p(99)<500'], // 99% of requests must complete below 0.5s
    // },
};

export function testPublish() {
    client.connect(BASE_URL, {
        plaintext: true,
    });

    const data = {
        subject: 'w',
        body: encoding.b64encode("data"),
        expirationSeconds: 100
    };

    for (let i = 0; i < 600; i++) {
        const response = client.invoke('broker.Broker/Publish', data);
        check(response, {
            'status is OK': (r) => r && r.status === grpc.StatusOK,
        });
        // console.log(JSON.stringify(response.message));
    }

    client.close();
    sleep(1);
}

export function testFetch() {
    client.connect(BASE_URL, {
        plaintext: true
    });

    const data = {
        subject: 'qqq',
        id: 421619,
    };
    const response = client.invoke('broker.Broker/Fetch', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    console.log(JSON.stringify(response.message));

    client.close();
    sleep(1);
}