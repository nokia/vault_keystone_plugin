package keystoneauth

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRegions(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "regions/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRegionList,
		},
	}
}

func pathRegions(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "regions/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "id",
			},
			"description": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "description",
			},
			"parent_region_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "parent_region_id",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRegionWrite,
			logical.ReadOperation:   b.pathRegionRead,
		},
	}
}

func (b *backend) Region(s logical.Storage, n string) (*regionEntry, error) {
	entry, err := s.Get("region/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result regionEntry

	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRegionRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("id").(string)
	description := data.Get("description").(string)
	parentRegionID := data.Get("parent_region_id").(string)

	region, err := b.Region(req.Storage, id)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown region: %s", id)), nil
	}

	conf, err2 := getconfig(req)

	if err2 != nil {
		return nil, fmt.Errorf("configure the Keystone connection with config/connection first")
	}

	keystoneURL := conf[0]
	token := conf[1]
	createdRegion, err2 := CreateRegion(id, description, parentRegionID, keystoneURL, token)

	if err2 != nil {
		return nil, fmt.Errorf("creation of the region failed")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"id":               createdRegion[0],
			"description":      createdRegion[1],
			"parent_region_id": createdRegion[2],
		},
	}, nil
}

func (b *backend) pathRegionList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("region/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRegionWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	id := data.Get("id").(string)
	description := data.Get("description").(string)
	parentRegionID := data.Get("parent_region_id").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("region/"+id, &regionEntry{
		Region_id:               id,
		Region_description:      description,
		Region_parent_region_id: parentRegionID,
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
			"id":      id,
			"created": true,
		},
	}, nil
}

type regionEntry struct {
	Region_id               string `json:"id" structs:"id" mapstructure:"id"`
	Region_description      string `json:"description" structs:"description" mapstructure:"description"`
	Region_parent_region_id string `json:"parent_region_id" structs:"parent_region_id" mapstructure:"parent_region_id"`
}
