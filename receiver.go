package azservicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type Receiver struct {
	vu       modules.VU
	receiver *azservicebus.Receiver
	timeout  time.Duration
}

type MessageState string

const (
	MessageStateActive    MessageState = "active"
	MessageStateDeferred  MessageState = "deferred"
	MessageStateScheduled MessageState = "scheduled"
)

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

func (sb *ServiceBus) CreateQueueReceiver(queue string) *sobek.Object {
	rt := sb.vu.Runtime()
	receiver, err := sb.cli.NewReceiverForQueue(queue, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot create queue receiver: %w", err))
	}

	return rt.ToValue(&Receiver{
		vu:       sb.vu,
		receiver: receiver,
		timeout:  sb.timeout,
	}).ToObject(rt)
}

func (sb *ServiceBus) CreateSubscriptionReceiver(topic, subscription string) *sobek.Object {
	rt := sb.vu.Runtime()
	receiver, err := sb.cli.NewReceiverForSubscription(topic, subscription, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot create subscription receiver: %w", err))
	}

	return rt.ToValue(&Receiver{
		vu:       sb.vu,
		receiver: receiver,
	}).ToObject(rt)
}

func (r *Receiver) Close() error {
	ctx, cancel := r.createContext()
	defer cancel()
	return r.receiver.Close(ctx)
}

func (r *Receiver) GetMessage() *ReceivedMessage {
	ctx, cancel := r.createContext()
	defer cancel()
	m, err := r.receiver.ReceiveMessages(ctx, 1, nil)
	if err != nil {
		common.Throw(r.vu.Runtime(), fmt.Errorf("AzServicebus cannot receive message: %w", err))
	}

	if len(m) == 0 {
		return nil
	}

	err = r.receiver.CompleteMessage(ctx, m[0], nil)
	if err != nil {
		common.Throw(r.vu.Runtime(), fmt.Errorf("AzServicebus cannot complete message: %w", err))
	}

	return mapServiceBusMessageToReceivedMessage(m[0])
}

func (r *Receiver) GetMessages(count int) []*ReceivedMessage {
	ctx, cancel := r.createContext()
	defer cancel()
	var receivedMessages []*ReceivedMessage
	for len(receivedMessages) < count {
		messages, err := r.receiver.ReceiveMessages(ctx, count-len(receivedMessages), nil)
		if err != nil {
			common.Throw(r.vu.Runtime(), fmt.Errorf("AzServicebus cannot receive messages: %w", err))
		}

		for _, m := range messages {
			receivedMessages = append(receivedMessages, mapServiceBusMessageToReceivedMessage(m))
		}

		for _, m := range messages {
			err = r.receiver.CompleteMessage(ctx, m, nil)
			if err != nil {
				common.Throw(r.vu.Runtime(), fmt.Errorf("AzServicebus cannot complete message: %w", err))
			}
		}
	}

	return receivedMessages
}

func (r *Receiver) createContext() (context.Context, context.CancelFunc) {
	// if r.timeout > 0 {
	// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	// 	return ctx, cancel
	// }

	return context.Background(), func() {}
}

func mapServiceBusMessageToReceivedMessage(m *azservicebus.ReceivedMessage) *ReceivedMessage {
	receivedMessage := &ReceivedMessage{
		Body:         m.Body,
		BodyAsString: string(m.Body),
		MessageID:    m.MessageID,
	}

	switch m.State {
	case azservicebus.MessageStateActive:
		receivedMessage.State = MessageStateActive
	case azservicebus.MessageStateDeferred:
		receivedMessage.State = MessageStateDeferred
	case azservicebus.MessageStateScheduled:
		receivedMessage.State = MessageStateScheduled
	}

	if m.ContentType != nil {
		receivedMessage.ContentType = *m.ContentType
	}

	if m.CorrelationID != nil {
		receivedMessage.CorrelationID = *m.CorrelationID
	}

	if m.DeadLetterErrorDescription != nil {
		receivedMessage.DeadLetterErrorDescription = *m.DeadLetterErrorDescription
	}

	if m.DeadLetterReason != nil {
		receivedMessage.DeadLetterReason = *m.DeadLetterReason
	}

	if m.DeadLetterSource != nil {
		receivedMessage.DeadLetterSource = *m.DeadLetterSource
	}

	if m.EnqueuedSequenceNumber != nil {
		receivedMessage.EnqueuedSequenceNumber = *m.EnqueuedSequenceNumber
	}

	if m.EnqueuedTime != nil {
		receivedMessage.EnqueuedTime = *m.EnqueuedTime
	}

	if m.ExpiresAt != nil {
		receivedMessage.ExpiresAt = *m.ExpiresAt
	}

	if m.LockedUntil != nil {
		receivedMessage.LockedUntil = *m.LockedUntil
	}

	if m.PartitionKey != nil {
		receivedMessage.PartitionKey = *m.PartitionKey
	}

	if m.ReplyTo != nil {
		receivedMessage.ReplyTo = *m.ReplyTo
	}

	if m.ReplyToSessionID != nil {
		receivedMessage.ReplyToSessionID = *m.ReplyToSessionID
	}

	if m.ScheduledEnqueueTime != nil {
		receivedMessage.ScheduledEnqueueTime = *m.ScheduledEnqueueTime
	}

	if m.SequenceNumber != nil {
		receivedMessage.SequenceNumber = *m.SequenceNumber
	}

	if m.SessionID != nil {
		receivedMessage.SessionID = *m.SessionID
	}

	if m.Subject != nil {
		receivedMessage.Subject = *m.Subject
	}

	if m.TimeToLive != nil {
		receivedMessage.TimeToLive = *m.TimeToLive
	}

	if m.To != nil {
		receivedMessage.To = *m.To
	}

	if m.ApplicationProperties != nil {
		receivedMessage.ApplicationProperties = make(map[string]string)
		for k, v := range m.ApplicationProperties {
			receivedMessage.ApplicationProperties[k] = v.(string)
		}
	}

	return receivedMessage
}
