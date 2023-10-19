import { check, sleep } from 'k6';
import { ServiceBus } from 'k6/x/azservicebus';

const config = {
    connectionString: __ENV.CONNECTION_STRING,
    timeout: 30000,
    insecureSkipVerify: true, // the docker servicebus emulator uses a self-signed cert
};

const servicebus = new ServiceBus(config);
const sender = servicebus.createSender('test-queue');
const receiver = servicebus.createQueueReceiver('test-queue');
const batch_message_count = 5;

let counter = 0;
export default function () {
    const data = [];
    for(let i = 0; i < batch_message_count; i++) {
        data.push({
            subject: 'test-subject',
            bodyAsString: `${++counter}the message`,
        });
    }

    sender.sendBatchMessages(data);

    sleep(1);

    const messages = receiver.getMessages(batch_message_count);
    check(messages, {
        'Has the expeted size': (m) => m.length === batch_message_count,
    });
}

export function teardown() {
    receiver.close();
    sender.close();
    servicebus.close();
}
