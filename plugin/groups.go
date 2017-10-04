package keystoneauth

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"reflect"
)

type Create_reponse_struct_group struct {
	Group *Group
}

type Group struct {
	Name        string
	Links       *LinksGroup
	Domain_id   string
	Description string
	Id          string
}

type LinksGroup struct {
	Self string
}

func CreateGroup(name string, description string, domain_id string, token string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_group

	internal_map := map[string]interface{}{"name": name}
	template_map := map[string]map[string]interface{}{"group": internal_map}

	if domain_id != "" {
		template_map["group"]["domain_id"] = domain_id
	}
	if description != "" {
		template_map["group"]["description"] = description
	}

	request := gorequest.New()
	var body2 string
	var errs []error
	_, body2, errs = request.Post("http://"+keystone_url+"/v3/groups").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if errs != nil {
		return nil, fmt.Errorf("creation of group has failed")
	}

	data := &Create_reponse_struct_group{
		Group: &Group{
			Links: &LinksGroup{},
		},
	}
	errmarshal := json.Unmarshal([]byte(body2), data)

	if errmarshal != nil {
		return nil, errmarshal
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.Group.Name), string(data.Group.Id)}
	return reply, nil
}

func DeleteGroup(group_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var errs3 []error

	_, _, errs3 = request.Delete("http://"+keystone_url+"/v3/groups/"+group_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if errs3 == nil {
		return "ok", nil
	}

	if errs3 != nil {
		return "", errs3[0]
	}

	return "ok", nil
}

func AddUserToGroup(group_id string, user_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var errs3 []error
	_, _, errs3 = request.Put("http://"+keystone_url+"/v3/groups/"+group_id+"/users/"+user_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if errs3 == nil {
		return "ok", nil
	}
	return "ok", nil
}
