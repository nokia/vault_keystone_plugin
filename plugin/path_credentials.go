package keystoneauth

import (
	"fmt"
	"context"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)


func pathCredentials(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "credentials/" + framework.GenericNameRegex("user_id"),
		Fields: map[string]*framework.FieldSchema{
			"user_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "User ID",
			},
			"blob": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "blob",
			},
			"type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "type",
			},
			"project_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "project_id",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathCredentialWrite,
			logical.ReadOperation:   b.pathCredentialRead,
		},
	}
}

func (b *backend) Credential(ctx context.Context, s logical.Storage, n string) (*credentialEntry, error) {
	entry, err := s.Get(ctx, "credential/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result credentialEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathCredentialRead(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	user_id := data.Get("user_id").(string)
	blob := data.Get("blob").(string)
	thetype := data.Get("type").(string)
	project_id := data.Get("project_id").(string)

	credential, err := b.Credential(ctx, req.Storage, blob)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown credential: %s", blob)), nil
	}

	conf, err2 := getconfig(ctx, req)
	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	keystone_url := conf[0]
	token := conf[1]

	created_credential, err3 := CreateCredential(blob, thetype, user_id, project_id, token, keystone_url)
	created_credential_id := created_credential[1]
	if err3 != nil {
		return logical.ErrorResponse(fmt.Sprintf("creation of the credential %s failed ", blob)), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"blob":       credential.Credential_blob,
			"type":       credential.Credential_type,
			"project_id": credential.Credential_project_id,
			"id":         created_credential_id,
		},
	}, nil
}

func (b *backend) pathCredentialList(
	ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "credential/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathCredentialWrite(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	user_id := data.Get("user_id").(string)
	blob := data.Get("blob").(string)
	thetype := data.Get("type").(string)
	project_id := data.Get("project_id").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("credential/"+blob, &credentialEntry{
		Credential_user_id: 		   user_id,
		Credential_blob:               blob,
		Credential_type:               thetype,
		Credential_project_id:         project_id,
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
			"blob":    blob,
			"created": true,
		},
	}, nil
}

type credentialEntry struct {
	Credential_blob       string `json:"blob" structs:"blob" mapstructure:"blob"`
	Credential_user_id    string `json:"user_id" structs:"user_id" mapstructure:"user_id"`
	Credential_type       string `json:"type" structs:"type" mapstructure:"type"`
	Credential_project_id string `json:"project_id" structs:"project_id" mapstructure:"project_id"`
}
