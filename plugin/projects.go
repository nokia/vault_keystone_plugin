package keystoneauth

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"reflect"
)

type Create_reponse_struct_project struct {
	Project *Project
}

type Project struct {
	Is_domain   bool
	Description string
	Links       *LinksProject
	Enabled     bool
	Id          string
	Domain_id   string
	Name        string
}

type LinksProject struct {
	Self string
}


func CreateProject(name string, description string, domain_id string, enabled bool, is_domain bool, token string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_project

	request := gorequest.New()

	internal_map := map[string]interface{}{"name":name}
	template_map := map[string]map[string]interface{}{"project":internal_map}

	if description != "" {
		template_map["project"]["description"] = description
	}
	if domain_id != "" {
		template_map["project"]["domain_id"] = domain_id
	}
	if enabled != true {
		template_map["project"]["enabled"] = true
	}
	if is_domain != false {
		template_map["project"]["enabled"] = true
	}

	var body2 []byte
	var err3 []error

	_, body2, err3 = request.Post("http://"+keystone_url+"/v3/projects").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).EndStruct(&create_reponse)

	if err3 != nil {
		return []string{"",""}, err3[0]
	}

	data := &Create_reponse_struct_project{
		Project: &Project{
			Links: &LinksProject{},
		},
	}

	errmarshal := json.Unmarshal([]byte(body2), data)

	if errmarshal != nil {
		return nil, errmarshal
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.Project.Name), string(data.Project.Id)}
	return reply, nil
}

func DeleteProject(project_id string, token string, keystone_url string) (string, error) {

	request := gorequest.New()
	var err []error

	_, _, err = request.Delete("http://"+keystone_url+"/v3/projects/"+project_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err != nil {
		return "",err[0]
	}
	return "ok",nil
}
