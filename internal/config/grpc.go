package config

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

type Grpc struct {
	Host            string `yaml:"host" env-default:"localhost" env-required:"true"`
	Port            string `yaml:"port" env-default:"443" env-required:"true"`
	CredentialsPath string `yaml:"credentials_path"`
	Timeout         int    `yaml:"timeout" env-default:"5"`
}

func (g *Grpc) Address() string {
	return g.Host + ":" + g.Port
}

func (g *Grpc) Credentials() (credentials.TransportCredentials, error) {
	if g.CredentialsPath != "" {
		creds, err := credentials.NewClientTLSFromFile(g.CredentialsPath, "")
		if err != nil {
			return nil, err
		}

		return creds, nil
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(config), nil
}
