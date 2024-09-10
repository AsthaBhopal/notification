package notification

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/AsthaBhopal/pkgGoAsthaLogs/asthahttp/v2"
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

func (s *Client) Initialize(otel bool, verboseSpans bool, transport http.RoundTripper) error {
	opts := []option.ClientOption{option.WithCredentialsFile(s.Filepath), option.WithHTTPClient(asthahttp.InitHttpClient(asthahttp.HttpConfig{
		OtelHttpTransport: otel,
		VerboseSpans:      verboseSpans,
		ServiceName:       "Firebase-Client",
		BaseTransport:     transport,
	}))}
	FirebaseApp, err := firebase.NewApp(s.Context, nil, opts...)
	if err != nil {
		return err
	}
	s.FcmClient, err = FirebaseApp.Messaging(s.Context)
	return err
}

func (s *Client) Send(payload Message, ctx context.Context) {
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
