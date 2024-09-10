package notification

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
)

type Client struct {
	Filepath  string
	Context   context.Context
	FcmClient *messaging.Client
}

type Message struct {
	Title      string //Heading
	Subtitle   string //subtitle
	ClientCode string //clientcode
	Data       map[string]string
}

func (s *Client) Initialize(otel bool) error {
	opts := []option.ClientOption{option.WithCredentialsFile(s.Filepath)}
	FirebaseApp, err := firebase.NewApp(s.Context, nil, opts...)
	if err != nil {
		return err
	}
	s.FcmClient, err = FirebaseApp.Messaging(s.Context)
	return err
}

func (s *Client) Send(payload Message, ctx context.Context) {
	tracer := otel.Tracer("Firebase FCM")

	// Start a span for the HTTP request
	ctx, span := tracer.Start(ctx, "FCM Request", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()
	if len(payload.Data) == 0 {
		payload.Data = make(map[string]string, 1)
	}
	payload.Data["senderId"] = "flow_inhouse"
	s.FcmClient.Send(ctx, &messaging.Message{
		Data: payload.Data,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Subtitle,
		},
		Topic: payload.ClientCode,
	})
}
