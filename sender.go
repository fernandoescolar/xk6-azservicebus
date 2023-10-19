package azservicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type Sender struct {
	vu      modules.VU
	sender  *azservicebus.Sender
	timeout time.Duration
}

type Message struct {
	ApplicationProperties map[string]string `js:"applicationProperties"`
	Body                  []byte            `js:"body"`
	BodyAsString          string            `js:"bodyAsString"`
	ContentType           string            `js:"contentType"`
	CorrelationID         string            `js:"correlationID"`
	MessageID             string            `js:"messageID"`
	PartitionKey          string            `js:"partitionKey"`
	SessionID             string            `js:"sessionID"`
	Subject               string            `js:"subject"`
	TimeToLive            time.Duration     `js:"timeToLive"`
	To                    string            `js:"to"`
}

func (sb *ServiceBus) CreateSender(queueOrTopic string) *goja.Object {
	rt := sb.vu.Runtime()
	if sb.cli == nil {
		common.Throw(rt, fmt.Errorf("AzServicebus connection should not be nil"))
	}

	sender, err := sb.cli.NewSender(queueOrTopic, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot create sender: %w", err))
	}

	return rt.ToValue(&Sender{
		vu:      sb.vu,
		sender:  sender,
		timeout: sb.timeout,
	}).ToObject(rt)
}

func (s *Sender) Close() error {
	ctx, cancel := s.createContext()
	defer cancel()
	return s.sender.Close(ctx)
}

func (s *Sender) Send(message string) {
	rt := s.vu.Runtime()
	ctx, cancel := s.createContext()
	defer cancel()

	err := s.sender.SendMessage(ctx, &azservicebus.Message{
		Body: []byte(message),
	}, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot send message: %w", err))
	}
}

func (s *Sender) SendMessage(message *Message) {
	rt := s.vu.Runtime()
	ctx, cancel := s.createContext()
	defer cancel()

	err := s.sender.SendMessage(ctx, mapMessageToSeriveBus(message), nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot send message: %w", err))
	}
}

func (s *Sender) SendBatch(messages []string) {
	rt := s.vu.Runtime()
	ctx1, cancel1 := s.createContext()
	defer cancel1()

	batch, err := s.sender.NewMessageBatch(ctx1, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot create batch: %w", err))
	}

	for _, message := range messages {
		if err := batch.AddMessage(&azservicebus.Message{Body: []byte(message)}, nil); err != nil {
			common.Throw(rt, fmt.Errorf("AzServicebus cannot add message to batch: %w", err))
		}
	}

	ctx2, cancel2 := s.createContext()
	defer cancel2()
	if err := s.sender.SendMessageBatch(ctx2, batch, nil); err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot send batch: %w", err))
	}
}

func (s *Sender) SendBatchMessages(messages []*Message) {
	rt := s.vu.Runtime()
	ctx1, cancel1 := s.createContext()
	defer cancel1()

	batch, err := s.sender.NewMessageBatch(ctx1, nil)
	if err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot create batch: %w", err))
	}

	for _, message := range messages {
		if err := batch.AddMessage(mapMessageToSeriveBus(message), nil); err != nil {
			common.Throw(rt, fmt.Errorf("AzServicebus cannot add message to batch: %w", err))
		}
	}

	ctx2, cancel2 := s.createContext()
	defer cancel2()
	if err := s.sender.SendMessageBatch(ctx2, batch, nil); err != nil {
		common.Throw(rt, fmt.Errorf("AzServicebus cannot send batch: %w", err))
	}
}

func (s *Sender) createContext() (context.Context, context.CancelFunc) {
	// if s.timeout > 0 {
	// 	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	// 	return ctx, cancel
	// }

	return context.Background(), func() {}
}

func mapMessageToSeriveBus(message *Message) *azservicebus.Message {
	sbMessage := &azservicebus.Message{}
	if message.ApplicationProperties != nil {
		for k, v := range message.ApplicationProperties {
			sbMessage.ApplicationProperties[k] = v
		}
	}

	if message.Body != nil {
		sbMessage.Body = message.Body
	}

	if message.BodyAsString != "" {
		sbMessage.Body = []byte(message.BodyAsString)
	}

	if message.ContentType != "" {
		sbMessage.ContentType = &message.ContentType
	}

	if message.CorrelationID != "" {
		sbMessage.CorrelationID = &message.CorrelationID
	}

	if message.MessageID != "" {
		sbMessage.MessageID = &message.MessageID
	}

	if message.PartitionKey != "" {
		sbMessage.PartitionKey = &message.PartitionKey
	}

	if message.SessionID != "" {
		sbMessage.SessionID = &message.SessionID
	}

	if message.Subject != "" {
		sbMessage.Subject = &message.Subject
	}

	if message.TimeToLive != 0 {
		sbMessage.TimeToLive = &message.TimeToLive
	}

	if message.To != "" {
		sbMessage.To = &message.To
	}
	return sbMessage
}
