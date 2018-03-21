package keystoneauth

import (
	"fmt"
	"context"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
)

func pathListDomains(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "domains/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathDomainList,
		},
	}
}

func pathDomains(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "domains/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "User name",
			},
			"description": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "description",
				Default:     "",
			},
			"enabled": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "enabled",
				Default:     true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDomainWrite,
			logical.ReadOperation:   b.pathDomainRead,
		},
	}
}

func (b *backend) Domain(ctx context.Context, s logical.Storage, n string) (*domainEntry, error) {
	entry, err := s.Get(ctx, "domain/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result domainEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathDomainRead(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	description := data.Get("description").(string)
	enabled := data.Get("enabled").(bool)

	domain, err := b.Domain(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown domain: %s", name)), nil
	}

	conf, err2 := getconfig(ctx, req)
	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	var namepostfix string
	namepostfix, _ = credsutil.RandomAlphaNumeric(20, true)
	name = fmt.Sprintf("%s_%s_%s ","vault", name, namepostfix[4:20])

	keystone_url := conf[0]
	token := conf[1]

	created_domain, err3 := CreateDomain(name, description, enabled, token, keystone_url)
	created_domain_id := created_domain[1]
	if err3 != nil {
		return logical.ErrorResponse(fmt.Sprintf("creation of the domain %s failed ", name)), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":        name,
			"description": domain.Domain_description,
			"enabled":     domain.Domain_enabled,
			"id":					 created_domain_id,
		},
	}, nil
}

func (b *backend) pathDomainList(
	ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "domain/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathDomainWrite(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	description := data.Get("description").(string)
	enabled := data.Get("enabled").(bool)

	// Store it
	entry, err := logical.StorageEntryJSON("domain/"+name, &domainEntry{
		Domain_name:        name,
		Domain_description: description,
		Domain_enabled:     enabled,
	})

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":               name,
			"created":						true,
		},
	}, nil
}

type domainEntry struct {
	Domain_name        string `json:"name" structs:"name" mapstructure:"name"`
	Domain_description string `json:"description" structs:"description" mapstructure:"description"`
	Domain_enabled     bool   `json:"enabled" structs:"enabled" mapstructure:"enabled"`
}
