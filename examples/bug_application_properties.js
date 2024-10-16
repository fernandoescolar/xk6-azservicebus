import { sleep, check } from 'k6';
import { ServiceBus } from 'k6/x/azservicebus';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const config = {
    connectionString: __ENV.CONNECTION_STRING,
    timeout: 30000,
    insecureSkipVerify: true, // the docker servicebus emulator uses a self-signed cert
};

const servicebus = new ServiceBus(config);
const sender = servicebus.createSender('test-topic');
const receiver = servicebus.createSubscriptionReceiver('test-topic', 'test-sub1');
const batch_message_count = 2;

export default function () {
    const data = [];

    for(let i = 0; i < batch_message_count; i++) {
        var properties = {
            'Type': "MessageType"
        };

        data.push({
            messageID: uuidv4(),
            applicationProperties: properties,
            contentType: 'application/json',
            subject: 'example',
            bodyAsString: JSON.stringify({ Title: 'Example'}),
        });
    }

    sender.sendBatchMessages(data);

    sleep(1);

    const messages = receiver.getMessages(batch_message_count);
    console.log(messages);
    console.log(`messages: ${messages.length}`)
    check(messages, {
        'Has the expected size': (m) => m.length === batch_message_count,
    });
}

export function teardown() {
    receiver.close();
    sender.close();
    servicebus.close();
}