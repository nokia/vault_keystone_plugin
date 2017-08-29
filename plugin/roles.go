package keystoneauth

import (
	"encoding/json"
	"reflect"
	"github.com/parnurzeal/gorequest"
)

type Create_reponse_struct_role struct {
	Role *Role
}

type Role struct {
	DomainId 		string
	Links       *LinksRole
	Id			    string
	Name        string
}

type LinksRole struct {
	Self string
}


func CreateRole(name string, domain_id string, token string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_domain

	request := gorequest.New()

	var body2 string
	var err []error

	internal_map := map[string]interface{}{"name":name}
	template_map := map[string]map[string]interface{}{"role":internal_map}

	if domain_id != "" {
		template_map["role"]["domain_id"] = domain_id
	}

	_, body2, err = request.Post("http://"+keystone_url+"/v3/roles").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if err != nil {
			return []string{"",""},err[0]
	}

	data := &Create_reponse_struct_role{
		Role: &Role{
			Links: &LinksRole{},
		},
	}

	err2 := json.Unmarshal([]byte(body2), data)

	if err2 != nil {
			return []string{"",""},err2
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.Role.Name), string(data.Role.Id)}
	return reply, err2
}

func GroupOnDomain(domain_id string, group_id string, role_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var err []error
	_, _, err = request.Put("http://"+keystone_url+"/v3/domains/"+domain_id+"/groups/"+group_id+"/roles/"+role_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err == nil {
		return "ok", nil
	}
	if err != nil {
		return "", err[0]
	}
	return "",nil
}

func UserOnDomain(user_id string, domain_id string, role_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var err []error
	_, _, err = request.Put("http://"+keystone_url+"/v3/domains/"+domain_id+"/users/"+user_id+"/roles/"+role_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err == nil {
		return "ok", nil
	}
	if err != nil {
		return "", err[0]
	}
	return "",nil
}

func GroupOnProject(group_id string, project_id string, role_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var err []error
	_, _, err = request.Put("http://"+keystone_url+"/v3/projects/"+project_id+"/groups/"+group_id+"/roles/"+role_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err == nil {
		return "ok", nil
	}
	if err != nil {
		return "", err[0]
	}
	return "",nil
}

func UserOnProject(user_id string, project_id string, role_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var err []error
	_, _, err = request.Put("http://"+keystone_url+"/v3/projects/"+project_id+"/users/"+user_id+"/roles/"+role_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err == nil {
		return "ok", nil
	}
	if err != nil {
		return "", err[0]
	}
	return "",nil
}
