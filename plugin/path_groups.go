package keystoneauth

import (
	"fmt"
	"context"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
)

func pathListGroups(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "groups/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathGroupList,
		},
	}
}

func pathGroups(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "groups/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name",
			},
			"description": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "description",
				Default:     "",
			},
			"domain_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "domain_id",
				Default:     "",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathGroupWrite,
			logical.ReadOperation:   b.pathGroupRead,
		},
	}
}

func pathGroupsUsers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "groups/" + framework.GenericNameRegex("group_id") + "/users/" + framework.GenericNameRegex("user_id"),
		Fields: map[string]*framework.FieldSchema{
			"group_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "group_id",
			},
			"user_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "user_id",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathGroupAddUser,
		},
	}
}

func (b *backend) Group(ctx context.Context, s logical.Storage, n string) (*groupEntry, error) {
	entry, err := s.Get(ctx, "group/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result groupEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathGroupRead(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	var namepostfix string
	namepostfix, _ = credsutil.RandomAlphaNumeric(20, true)
	name := data.Get("name").(string)
	description := data.Get("description").(string)
	domain_id := data.Get("domain_id").(string)

	group, err := b.Group(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown group: %s", name)), nil
	}

	conf, err2 := getconfig(ctx, req)

	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	name = fmt.Sprintf("%s_%s_%s ", "vault", name, namepostfix[4:20])
	keystone_url := conf[0]
	token := conf[1]

	//Create the group

	created_grp, err3 := CreateGroup(name, description, domain_id, token, keystone_url)
	created_grp_id := created_grp[1]

	if err3 != nil {
		return nil, fmt.Errorf("creation of the group failed")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":        name,
			"description": group.Group_description,
			"domain_id":   group.Group_domain_id,
			"id":          created_grp_id,
		},
	}, nil
}

func (b *backend) pathGroupList(
	ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "group/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathGroupWrite(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	description := data.Get("description").(string)
	domain_id := data.Get("domain_id").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("group/"+name, &groupEntry{
		Group_name:        name,
		Group_description: description,
		Group_domain_id:   domain_id,
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
			"name":    name,
			"created": true,
		},
	}, nil
}

func (b *backend) pathGroupAddUser(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	user_id := data.Get("user_id").(string)
	group_id := data.Get("group_id").(string)

	conf, err2 := getconfig(ctx, req)

	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}
	
	keystone_url := conf[0]
	token := conf[1]

	addedUser, err := AddUserToGroup(group_id, user_id, token, keystone_url)
	if err != nil {
		return nil, err
	}

	if addedUser == "" {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"user":  user_id,
			"added": true,
		},
	}, nil
}

type groupEntry struct {
	Group_name        string `json:"name" structs:"name" mapstructure:"name"`
	Group_description string `json:"description" structs:"description" mapstructure:"description"`
	Group_domain_id   string `json:"domain_id" structs:"domain_id" mapstructure:"domain_id"`
}
