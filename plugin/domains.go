package keystoneauth

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"reflect"
)

type Create_reponse_struct_domain struct {
	Domain *Domain
}

type Domain struct {
	Description string
	Links       *LinksDomain
	Enabled     bool
	Id          string
	Name        string
}

type LinksDomain struct {
	Self string
}

func CreateDomain(name string, description string, enabled bool, token string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_domain

	request := gorequest.New()

	internal_map := map[string]interface{}{"name":name}
	template_map := map[string]map[string]interface{}{"domain":internal_map}

	if description != "" {
		template_map["domain"]["description"] = description
	}
	if enabled != true {
		template_map["domain"]["enabled"] = false
	}

	var body2 string
	var err []error

	_, body2, err = request.Post("http://"+keystone_url+"/v3/domains").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if err != nil {
			return []string{"",""},err[0]
	}

	data := &Create_reponse_struct_domain{
		Domain: &Domain{
			Links: &LinksDomain{},
		},
	}

	err2 := json.Unmarshal([]byte(body2), data)

	if err2 != nil {
			return []string{"",""},err2
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.Domain.Name), string(data.Domain.Id)}
	return reply, err2
}

func DeleteDomain(domain_id string, token string, keystone_url string) (string, error) {

	request := gorequest.New()

	var err []error

	_, _, err = request.Delete("http://"+keystone_url+"/v3/projects/"+domain_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if err != nil {
			return "",err[0]
	}

	return "ok",nil
}
