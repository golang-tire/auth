package pubsub

import (
	"context"

	"github.com/google/uuid"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/garsue/watermillzap"
	"github.com/golang-tire/pkg/log"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	watermillSql "github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
)

type PubSub struct {
	Subscriber *watermillSql.Subscriber
	Publisher  *watermillSql.Publisher
	router     *message.Router
}

var pubSub *PubSub

func Init(ctx context.Context, db *db.DB) (*PubSub, error) {

	d, err := db.DB().DB()
	if err != nil {
		return nil, err
	}

	logger := watermillzap.NewLogger(log.Logger())
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	pubSub = &PubSub{}
	pubSub.router = router
	subscriber, err := watermillSql.NewSubscriber(
		d,
		watermillSql.SubscriberConfig{
			SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		},
		logger,
	)

	publisher, err := watermillSql.NewPublisher(
		d,
		sql.PublisherConfig{
			SchemaAdapter: sql.DefaultPostgreSQLSchema{},
		},
		logger,
	)

	pubSub.Publisher = publisher
	pubSub.Subscriber = subscriber

	errCh := make(chan error)
	go func() {
		err := pubSub.router.Run(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		select {
		case er := <-errCh:
			log.Error("run pubSub failed", log.Err(er))
		case <-ctx.Done():
			_ = pubSub.router.Close()
			_ = pubSub.Subscriber.Close()
			_ = pubSub.Publisher.Close()
			return
		}
	}()

	return pubSub, nil
}

func AddHandler(topic string, handler func(*message.Message) error) {
	pubSub.router.AddNoPublisherHandler(
		topic+"_"+uuid.New().String(),
		topic,
		pubSub.Subscriber,
		handler,
	)
}

func Publish(topic string, message *message.Message) error {
	return pubSub.Publisher.Publish(topic, message)
}
