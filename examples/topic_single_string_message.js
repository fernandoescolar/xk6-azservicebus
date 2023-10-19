import { check, sleep } from 'k6';
import { ServiceBus } from 'k6/x/azservicebus';

const config = {
    connectionString: __ENV.CONNECTION_STRING,
    timeout: 30000,
    insecureSkipVerify: true, // the docker servicebus emulator uses a self-signed cert
};

const servicebus = new ServiceBus(config);
const sender = servicebus.createSender('test-topic');
const receiver = servicebus.createSubscriptionReceiver('test-topic', 'test-sub1');

let counter = 0;
export default function () {
    const data = `${++counter}the message`;
    sender.send(data);

    sleep(1);

    const message = receiver.getMessage();
    check(message, {
        'Is expected message': (m) => m.bodyAsString === data,
    });
}

export function teardown() {
    receiver.close();
    sender.close();
    servicebus.close();
}
