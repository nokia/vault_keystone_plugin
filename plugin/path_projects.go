package keystoneauth

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListProjects(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "projects/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathProjectList,
		},
	}
}

func pathProjects(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "projects/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name",
			},
			"is_domain": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "is_domain",
				Default:     false,
			},
			"domain_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "domain_id",
				Default:     "",
			},
			"description": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "description",
				Default:     "",
			},
			"enabled": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "enabled",
				Default:     true,
			},
			"parent_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "parent_id",
				Default:     "",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: 	b.pathProjectWrite,
			logical.ReadOperation:   	b.pathProjectRead,
			logical.DeleteOperation:	b.pathProjectDelete,
		},
	}
}

func (b *backend) Project(s logical.Storage, n string) (*projectEntry, error) {
	entry, err := s.Get("project/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result projectEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathProjectRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	description := data.Get("description").(string)
	domain_id := data.Get("domain_id").(string)
	enabled := data.Get("enabled").(bool)
	is_domain := data.Get("is_domain").(bool)

	project, err := b.Project(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown project: %s", name)), nil
	}

	conf, err2 := getconfig(req)

	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	keystone_url := conf[0]
	token := conf[1]
	created_project, err2 := CreateProject(name, description, domain_id, enabled, is_domain, token, keystone_url)

	if err2 != nil {
		return nil, fmt.Errorf("creation of the project failed")
	}

	created_project_id := created_project[1]

	return &logical.Response{
		Data: map[string]interface{}{
			"name":        project.Project_name,
			"is_domain":   project.Project_is_domain,
			"domain_id":   project.Project_domain_id,
			"description": project.Project_description,
			"enabled":     project.Project_enabled,
			"parent_id":   project.Project_parent_id,
			"id":          created_project_id,
		},
	}, nil
}

func (b *backend) pathProjectList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("project/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathProjectWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	is_domain := data.Get("is_domain").(bool)
	domain_id := data.Get("domain_id").(string)
	description := data.Get("description").(string)
	enabled := data.Get("enabled").(bool)
	parent_id := data.Get("parent_id").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("project/"+name, &projectEntry{
		Project_name:        name,
		Project_is_domain:   is_domain,
		Project_domain_id:   domain_id,
		Project_description: description,
		Project_enabled:     enabled,
		Project_parent_id:   parent_id,
	})

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":    name,
			"created": true,
		},
	}, nil
}

func (b *backend) pathProjectDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

	conf, err := getconfig(req)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	keystone_url := conf[0]
	token := conf[1]

	if err != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	status, err := DeleteProject(keystone_url, token, name)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	if status == "NO_PROJECT" {
		return logical.ErrorResponse(fmt.Sprintf("unknown project: %s", name)), nil
	}

	if status == "" {
		if err := req.Storage.Delete("project/"+name); err != nil {
			return nil, err
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":    name,
			"deleted": true,
		},
	}, nil
}

type projectEntry struct {
	Project_name        string `json:"name" structs:"name" mapstructure:"name"`
	Project_is_domain   bool   `json:"is_domain" structs:"is_domain" mapstructure:"is_domain"`
	Project_domain_id   string `json:"domain_id" structs:"domain_id" mapstructure:"domain_id"`
	Project_description string `json:"description" structs:"description" mapstructure:"description"`
	Project_enabled     bool   `json:"enabled" structs:"enabled" mapstructure:"enabled"`
	Project_parent_id   string `json:"parent_id" structs:"parent_id" mapstructure:"parent_id"`
}
