package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/ThreeDotsLabs/watermill"

	"github.com/go-redis/redis/v8"

	"github.com/mirzakhany/watermill-redisstream/pkg/redisstream"

	"github.com/google/uuid"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/garsue/watermillzap"
	"github.com/golang-tire/pkg/log"
)

type PubSub struct {
	ctx        context.Context
	publisher  message.Publisher
	subscriber message.Subscriber
	router     *message.Router
	logger     watermill.LoggerAdapter
	rc         redis.UniversalClient
}

var pubSub *PubSub

func Init(ctx context.Context, rc redis.UniversalClient) (*PubSub, error) {

	logger := watermillzap.NewLogger(log.Logger())
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	publisher, err := redisstream.NewPublisher(
		ctx,
		rc,
		&redisstream.DefaultMarshaler{},
		logger,
	)

	subscriber, err := redisstream.NewSubscriber(
		ctx,
		redisstream.SubscriberConfig{
			Consumer:        "auth-consumer",
			ConsumerGroup:   "auth-consumer-rbac",
			DoNotDelMessage: false,
		},
		rc,
		&redisstream.DefaultMarshaler{},
		logger,
	)

	pubSub = &PubSub{
		ctx:        ctx,
		router:     router,
		logger:     logger,
		subscriber: subscriber,
		publisher:  publisher,
	}

	router.AddMiddleware(middleware.Recoverer)
	return pubSub, nil
}

func (ps *PubSub) AddHandler(topic string, handler func(*message.Message) error) error {
	ps.router.AddNoPublisherHandler(
		topic+"_"+uuid.New().String(),
		topic,
		ps.subscriber,
		handler,
	)
	return nil
}

func (ps *PubSub) Run(ctx context.Context) {
	go func() {
		_ = ps.router.Run(ctx)
		<-ctx.Done()
	}()
}

func Get() *PubSub {
	return pubSub
}

func (ps *PubSub) Publish(topic string, message *message.Message) error {
	return ps.publisher.Publish(topic, message)
}
