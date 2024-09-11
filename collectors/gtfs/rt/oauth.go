package rt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Auth interface {
	GetClient(ctx context.Context) (*http.Client, error)
}

type AzureClientCredentials struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Resource     string
}

func (o *AzureClientCredentials) GetClient(ctx context.Context) (*http.Client, error) {
	conf := &clientcredentials.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		EndpointParams: url.Values{
			"resource": []string{o.Resource},
		},
		TokenURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", o.TenantID),
		AuthStyle: oauth2.AuthStyleInParams,
	}

	return conf.Client(ctx), nil
}
