package keystoneauth

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"connection_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keystone URL",
			},
			"admin_auth_token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keystone admin user auth_token",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionWrite,
			logical.ReadOperation:   b.pathConnectionRead,
		},
	}
}

func (b *backend) pathConnectionRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get("config/connection")
	if err != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}
	if entry == nil {
		return nil, nil
	}

	var config connectionConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: structs.New(config).Map(),
	}, nil
}

func (b *backend) pathConnectionWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	connURL := data.Get("connection_url").(string)
	adminAuthToken := data.Get("admin_auth_token").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
		ConnectionURL:  connURL,
		AdminAuthToken: adminAuthToken,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"created":						true,
		},
	}, nil
}

type connectionConfig struct {
	ConnectionURL  string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
	AdminAuthToken string `json:"admin_auth_token" structs:"admin_auth_token" mapstructure:"admin_auth_token"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to Keystone.
`

const pathConfigConnectionHelpDesc = `
This path configures the URL used to connect to Keystone.
URL should be in the format http://localhost:35357/v3/
`
