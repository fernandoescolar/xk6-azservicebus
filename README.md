# xk6-azservicebus

This is a [k6](https://go.k6.io/k6) extension using the [xk6](https://github.com/k6io/xk6) system, that allows to connect with Azure ServiceBus.

|  ❗ This extension isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

- [xk6-azservicebus](#xk6-azservicebus)
	- [Build](#build)
	- [API](#api)
		- [ServiceBus](#servicebus)
		- [Sender](#sender)
		- [Receiver](#receiver)
	- [License](#license)

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Install `xk6` framework for extending `k6`:
```shell
go install go.k6.io/xk6/cmd/xk6@latest
```

2. Build the binary:
```shell
xk6 build --with github.com/fernandoescolar/xk6-azservicebus@latest
```

3. Run a test

```shell
k6 run -e CONNECTION_STRING=your_azure_service_bus_connection_string your_tests.js
```

You can find some example javascript files in the [examples](examples) folder.

## API

### ServiceBus

A `ServiceBus` instance represents the connection with the Azure ServiceBus service and it is created with `new ServiceBus(configuration)`, where configuration attributes are:

| Attribute | Description |
| --- | --- |
| `connectionString` | (mandatory) is the Azure ServiceBus connection string you can find in the azure portal |
| `timeout` | (optional) is the operations timeout in millis  |
| `insecureSkipVerify` | (optional) if `true` it will allow untrusted certificates in connections |

Example:

```ts
import { ServiceBus } from 'k6/x/azservicebus';

const config = {
    connectionString: __ENV.CONNECTION_STRING,
    timeout: 30000,
};

const servicebus = new ServiceBus(config);
```

When you finish using the `ServiceBus` instance, you should close it using the `close()` method:

```ts
export function teardown() {
    servicebus.close();
}
```

### Sender

To send messages to Azure ServiceBus you have to create a new sender using the `createSender(topicOrQueue)` method of the `ServiceBus` instance. The `topicOrQueue` parameter is the name of the topic or the queue where the messages will be sent:

```ts
const sender = servicebus.createSender('test-topic');
// or
const sender = servicebus.createSender('test-queue');
```

Then, you can send messages to a topic or a queue using the following functions:

| Function | Description |
| --- | --- |
| `send(string)` | sends a string message to a topic or a queue |
| `sendMessage(message)` | sends a `Message` object to a topic or a queue |
| `sendBatch(string[])` | sends a batch of string messages to a topic or a queue |
| `sendBatchMessages(message[])` | sends a batch of `Message` objects to a topic or a queue |

Example:

```ts
const sender = servicebus.createSender('test-topic');

sender.send('hello azure service bus!');
sender.sendBatch(['hello azure service bus!', 'hello again!']);
sender.sendMessage({
  subject: 'my subject',
  bodyAsString: 'hello azure service bus!'
});
sender.sendBatchMessages([
  {
    subject: 'my subject',
    bodyAsString: 'hello azure service bus!'
  },
  {
    subject: 'my subject',
    bodyAsString: 'hello again!'
  }
]);

sender.close();
```

Once you have finished using the sender, you should close it using the `close()` method:

```ts
sender.close();
```

And the `Message` object has the following attributes:

| Attribute | Description |
| --- | --- |
| `applicationProperties` | (optional) is a map of string key/value pairs |
| `body` | (mandatory if `bodyAsString` has not been specified) is the message body in byte array format |
| `bodyAsString` | (mandatory if `body` has not been specified) is the message body in string format |
| `contentType` | (optional) is the content type of the message |
| `correlationID` | (optional) is the correlation ID of the message |
| `messageID` | (optional) is the message ID |
| `partitionKey` | (optional) is the partition key |
| `sessionID` | (optional) is the session ID |
| `subject` | (optional) is the subject of the message |
| `timeToLive` | (optional) is the time to live of the message |
| `to` | (optional) is the destination of the message |

### Receiver

To receive messages from Azure ServiceBus you have to create a new queue receiver using the `createQueueReceiver(queue)` method of the `ServiceBus` instance, or `createSubscriptionReceiver(topic, subscription)` method to create a subscription receiver. The `queue` parameter is the name of the queue where the messages will be received, and the `topic` and `subscription` parameters are the name of the topic and the subscription where the messages will be received:

```ts
const receiver = servicebus.createQueueReceiver('test-queue');
// or
const receiver = servicebus.createSubscriptionReceiver('test-topic', 'test-subscription');
```

Then, you can receive messages from a queue or a subscription using the following functions:

| Function | Description |
| --- | --- |
| `getMessage()` | receives a `ReceivedMessage` from a queue or a subscription |
| `getMessages(maxMessages)` | receives an array of `ReceivedMessage` from a queue or a subscription |

```ts
const receiver = servicebus.createQueueReceiver('test-queue');

const message = receiver.getMessage();
check(message, {
  'Is expected message': (m) => m.bodyAsString === expectedMessage,
})

receiver.close();
```

Once you have finished using the receiver, you should close it using the `close()` method:

```ts
receiver.close();
```

```go
type ReceivedMessage struct {
	ApplicationProperties      map[string]string `js:"applicationProperties"`
	Body                       []byte            `js:"body"`
	BodyAsString               string            `js:"bodyAsString"`
	ContentType                string            `js:"contentType"`
	CorrelationID              string            `js:"correlationID"`
	DeadLetterErrorDescription string            `js:"deadLetterErrorDescription"`
	DeadLetterReason           string            `js:"deadLetterReason"`
	DeadLetterSource           string            `js:"deadLetterSource"`
	DeliveryCount              uint32            `js:"deliveryCount"`
	EnqueuedSequenceNumber     int64             `js:"enqueuedSequenceNumber"`
	EnqueuedTime               time.Time         `js:"enqueuedTime"`
	ExpiresAt                  time.Time         `js:"expiresAt"`
	LockedUntil                time.Time         `js:"lockedUntil"`
	MessageID                  string            `js:"messageID"`
	PartitionKey               string            `js:"partitionKey"`
	ReplyTo                    string            `js:"replyTo"`
	ReplyToSessionID           string            `js:"replyToSessionID"`
	ScheduledEnqueueTime       time.Time         `js:"scheduledEnqueueTime"`
	SequenceNumber             int64             `js:"sequenceNumber"`
	SessionID                  string            `js:"sessionID"`
	State                      MessageState      `js:"state"`
	Subject                    string            `js:"subject"`
	TimeToLive                 time.Duration     `js:"timeToLive"`
	To                         string            `js:"to"`
}
```

And the `ReceivedMessage` object has the following attributes:

| Attribute | Description |
| --- | --- |
| `applicationProperties` | is a map of string key/value pairs |
| `body` | is the message body in byte array format |
| `bodyAsString` | is the message body in string format |
| `contentType` | is the content type of the message |
| `correlationID` | is the correlation ID of the message |
| `deadLetterErrorDescription` | is the error description of the dead letter message |
| `deadLetterReason` | is the reason of the dead letter message |
| `deadLetterSource` | is the source of the dead letter message |
| `deliveryCount` | is the delivery count of the message |
| `enqueuedSequenceNumber` | is the enqueued sequence number of the message |
| `enqueuedTime` | is the enqueued time of the message |
| `expiresAt` | is the expiration time of the message |
| `lockedUntil` | is the locked until time of the message |
| `messageID` | is the message ID |
| `partitionKey` | is the partition key |
| `replyTo` | is the reply to destination of the message |
| `replyToSessionID` | is the reply to session ID of the message |
| `scheduledEnqueueTime` | is the scheduled enqueue time of the message |
| `sequenceNumber` | is the sequence number of the message |
| `sessionID` | is the session ID |
| `state` | is the state of the message |
| `subject` | is the subject of the message |
| `timeToLive` | is the time to live of the message |
| `to` | is the destination of the message |

## License

The source code of this project is released under the [MIT License](LICENSE).