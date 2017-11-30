package keystoneauth

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"reflect"
)

type Create_reponse_struct_credential struct {
	CCredential *CCredential
}

type CCredential struct {
	DomainId  string
	Links     *LinksCredential
	Id        string
	ProjectId string
	Type      string
	UserId    string
}

type LinksCredential struct {
	Self string
}

func CreateCredential(blob string, thetype string, user_id string, project_id string, token string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_credential

	request := gorequest.New()

	var body2 string
	var err []error

	internal_map := map[string]interface{}{"blob": blob, "type": thetype, "user_id": user_id, "project_id": project_id}
	template_map := map[string]map[string]interface{}{"credential": internal_map}

	_, body2, err = request.Post("http://"+keystone_url+"/v3/credentials").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if err != nil {
		return []string{"", ""}, err[0]
	}

	data := &Create_reponse_struct_credential{
		CCredential: &CCredential{
			Links: &LinksCredential{},
		},
	}

	err2 := json.Unmarshal([]byte(body2), data)

	if err2 != nil {
		return []string{"", ""}, err2
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.CCredential.Id), string(data.CCredential.UserId)}
	return reply, err2
}

func DeleteCredential(credential_id string, token string, keystone_url string) (string, error) {

	request := gorequest.New()

	var err []error

	template_map := map[string]interface{}{"credential_id": credential_id}

	_, _, err = request.Delete("http://"+keystone_url+"/v3/credentials").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if err != nil {
		return "", err[0]
	}

	return "ok", nil
}
