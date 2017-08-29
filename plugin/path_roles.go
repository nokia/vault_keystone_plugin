package keystoneauth

import (
	"fmt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"reflect"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name",
			},
			"domain_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "domain_id",
				Default:		 "",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRoleWrite,
			logical.ReadOperation: b.pathRoleRead,
		},
	}
}

func (b *backend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	domain_id := data.Get("domain_id").(string)

	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	conf, err2 := getconfig(req)
	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	keystone_url := conf[0]
	token := conf[1]

	created_role, err3 := CreateRole(name, domain_id, token, keystone_url)
	created_role_id := created_role[1]
	if err3 != nil {
		return nil, fmt.Errorf("creation of the role failed")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":        role.Role_name,
			"domain_id":   role.Role_domain_id,
			"id":					 created_role_id,
		},
	}, nil
}

func pathRolesGroupOnDomain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role") + "/groups/" + framework.GenericNameRegex("group") + "/domains/" + framework.GenericNameRegex("domain"),
		Fields: map[string]*framework.FieldSchema{
			"domain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "domain",
			},
			"group": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "group",
			},
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "role",
			},
			"action": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "grant",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRolesGroupOnDomainWrite,
		},
	}
}

func pathRolesUserOnDomain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role") + "/users/" + framework.GenericNameRegex("user") + "/domains/" + framework.GenericNameRegex("domain"),
		Fields: map[string]*framework.FieldSchema{
			"domain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "domain",
			},
			"user": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "user",
			},
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "role",
			},
			"action": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "grant",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRolesUserOnDomainWrite,
		},
	}
}

func pathRolesGroupOnProject(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role") + "/groups/" + framework.GenericNameRegex("group") + "/projects/" + framework.GenericNameRegex("project"),
		Fields: map[string]*framework.FieldSchema{
			"group": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "group",
			},
			"project": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "project",
			},
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "role",
			},
			"action": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "grant",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRolesGroupOnProjectWrite,
		},
	}
}

func pathRolesUserOnProject(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role") + "/users/" + framework.GenericNameRegex("user") + "/projects/" + framework.GenericNameRegex("project"),
		Fields: map[string]*framework.FieldSchema{
			"user": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "user",
			},
			"project": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "project",
			},
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "role",
			},
			"action": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "grant",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRolesUserOnProjectWrite,
		},
	}
}

func getconfig(req *logical.Request) ([]string, error) {
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

	return []string{string(config.ConnectionURL), string(config.AdminAuthToken)}, nil
}

func (b *backend) pathRolesGroupOnDomainWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	domain_id := data.Get("domain").(string)
	group_id := data.Get("group").(string)
	role_id := data.Get("role").(string)

	conf, err := getconfig(req)

	if err != nil {
		return nil, err
	}

	keystone_url := conf[0]
	token := conf[1]

	// make a request to Keystone

	god, errgod := GroupOnDomain(domain_id, group_id, role_id, token, keystone_url)

	if errgod != nil {
		return nil, errgod
	}

	reflect.TypeOf(god)

	return &logical.Response{
		Data: map[string]interface{}{
			"status": "successful",
		},
	}, nil
}

func (b *backend) pathRolesGroupOnProjectWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	group_id := data.Get("group").(string)
	project_id := data.Get("project").(string)
	role_id := data.Get("role").(string)

	reflect.TypeOf(group_id)
	reflect.TypeOf(project_id)
	reflect.TypeOf(role_id)

	conf, err := getconfig(req)

	if err != nil {
		return nil, err
	}

	keystone_url := conf[0]
	token := conf[1]

	// make a request to Keystone

	gop, errgop := GroupOnProject(group_id, project_id, role_id, token, keystone_url)

	if errgop != nil {
		return nil, errgop
	}

	reflect.TypeOf(gop)

	return &logical.Response{
		Data: map[string]interface{}{
			"status": "successful",
		},
	}, nil
}

func (b *backend) pathRolesUserOnDomainWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	domain_id := data.Get("domain_id").(string)
	user_id := data.Get("user_id").(string)
	role_id := data.Get("role_id").(string)

	reflect.TypeOf(domain_id)
	reflect.TypeOf(user_id)
	reflect.TypeOf(role_id)

	conf, err := getconfig(req)

	if err != nil {
		return nil, err
	}

	keystone_url := conf[0]
	token := conf[1]

	// make a request to Keystone

	uod, erruod := UserOnDomain(user_id, domain_id, role_id, token, keystone_url)

	if erruod != nil {
		return nil, erruod
	}

	reflect.TypeOf(uod)

	return &logical.Response{
		Data: map[string]interface{}{
			"status": "successful",
		},
	}, nil
}

func (b *backend) pathRolesUserOnProjectWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	project_id := data.Get("project_id").(string)
	user_id := data.Get("user_id").(string)
	role_id := data.Get("role_id").(string)

	reflect.TypeOf(project_id)
	reflect.TypeOf(user_id)
	reflect.TypeOf(role_id)

	conf, err := getconfig(req)

	if err != nil {
		return nil, err
	}

	keystone_url := conf[0]
	token := conf[1]

	uop, erruop := UserOnProject(user_id, project_id, role_id, token, keystone_url)

	if erruop != nil {
		return nil, erruop
	}

	reflect.TypeOf(uop)

	if erruop != nil {
		return nil, erruop
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": "successful",
		},
	}, nil
}


func (b *backend) pathRoleWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	fmt.Println(name)
	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, &roleEntry{
		Role_name: name,
	})

	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":               name,
			"created":						true,
		},
	}, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}


type roleEntry struct {
	Role_name string `json:"name" structs:"name" mapstructure:"name"`
	Role_domain_id string `json:"domain_id" structs:"domain_id" mapstructure:"domain_id"`
}
