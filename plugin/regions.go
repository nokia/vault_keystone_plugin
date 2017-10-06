package keystoneauth

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/parnurzeal/gorequest"
)

type CreateResponsStructRegion struct {
	Region *Region
}

//Region represents OpenStack region domain model
type Region struct {
	ID             string
	Description    string
	Links          *LinksRegion
	ParentRegionID string
}

type LinksRegion struct {
	Self string
}

func CreateRegion(id string, description string, parentRegionID string, keystoneUrl string, token string) ([]string, error) {

	var createResponse CreateResponsStructRegion
	internalMap := map[string]interface{}{"id": id, "description": description}
	templateMap := map[string]map[string]interface{}{"region": internalMap}

	if parentRegionID != "" {
		templateMap["region"]["parent_region_id"] = parentRegionID
	}

	req := gorequest.New()
	var body string
	var errs []error

	_, body, errs = req.
		Post("http://"+keystoneUrl+"/v3/regions").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(templateMap).
		End()

	if errs != nil {
		return nil, fmt.Errorf("Failed to create a region")
	}

	data := &CreateResponsStructRegion{
		Region: &Region{
			Links: &LinksRegion{},
		},
	}

	unmarshalErr := json.Unmarshal([]byte(body), data)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	reflect.TypeOf(createResponse)
	reply := []string{string(data.Region.ID), string(data.Region.Description)}

	return reply, nil
}

func DeleteRegion(regionID string, token string, keystoneURL string) (string, error) {

	req := gorequest.New()
	var errs []error

	_, _, errs = req.
		Delete("http://"+keystoneURL+"/v3/regions/"+regionID).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if errs != nil {
		return "", errs[0]
	}

	return "ok", nil
}
