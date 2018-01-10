package keystoneauth

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/parnurzeal/gorequest"
)

type UserEC2response struct {
	Credential *Credential
}

type Credential struct {
	Access string
	Secret string
}

type Create_reponse_struct_user struct {
	User *User
}

type User struct {
	Name      string
	Links     *LinksUser
	Domain_id string
	Enabled   bool
	Id        string
}

type LinksUser struct {
	Self string
}

type UsersResponse struct {
  User []struct {
    Id		string	`json:"id"`
    Name 	string	`json:"name"`
  } `json:"users"`
}

func UserEC2(user_id string, tenant_id string, token string, keystone_url string) ([]string, error) {
	var create_reponse UserEC2response

	request := gorequest.New()
	var body2 string
	var errs []error
	_, body2, errs = request.Post("http://"+keystone_url+"/v3/users/"+user_id+"/credentials/OS-EC2").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(`{"tenant_id":"`+tenant_id+`"}`).End()

		if errs != nil {
			return nil, fmt.Errorf("generation of EC2 credentials has failed")
		}

		data := &UserEC2response{
			Credential: &Credential{},
		}
		errmarshal := json.Unmarshal([]byte(body2), data)

		if errmarshal != nil {
			return nil, errmarshal
		}

		reflect.TypeOf(create_reponse)
		reply := []string{string(data.Credential.Access), string(data.Credential.Secret)}
		return reply, nil

}

func CreateUser(default_project_id string, name string, password string, enabled bool, token string, domain_id string, keystone_url string) ([]string, error) {
	var create_reponse Create_reponse_struct_user

	internal_map := map[string]interface{}{"name":name,"password":password,"enabled":true}
	template_map := map[string]map[string]interface{}{"user":internal_map}

	if domain_id != "" {
		template_map["user"]["domain_id"] = domain_id
	}
	if default_project_id != "" {
		template_map["user"]["default_project_id"] = default_project_id
	}

	request := gorequest.New()
	var body2 string
	var errs []error
	_, body2, errs = request.Post("http://"+keystone_url+"/v3/users").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").
		Send(template_map).End()

	if errs != nil {
		return nil, fmt.Errorf("creation of user has failed")
	}

	data := &Create_reponse_struct_user{
		User: &User{
			Links: &LinksUser{},
		},
	}
	errmarshal := json.Unmarshal([]byte(body2), data)

	if errmarshal != nil {
		return nil, errmarshal
	}

	reflect.TypeOf(create_reponse)
	reply := []string{string(data.User.Name), string(data.User.Id)}
	return reply, nil
}

func ListAllOpenStackUsers(
	name string, token string, keystone_url string) (map[string]string, error) {

		var data string
		var err []error
		to_delete_users := make (map[string]string)

		request := gorequest.New()
		_, data, err = request.Get("http://" + keystone_url + "/v3/users/").
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()
		if err != nil {
			log.Println(err[0])
		}

	  user_struct := new(UsersResponse)
	  err2 := json.Unmarshal([]byte(data), &user_struct)
	  if err2 != nil {
	    log.Fatal(err2)
	  }

		if len(name) > 0 {

			reg_exp := "vault_" + name + "_[a-z0-9]{16}"

		  for i, u := range user_struct.User {
		    fmt.Sprintf("%d. %s, %s\n", i, u.Id, u.Name)
		    match, err := regexp.MatchString(reg_exp, u.Name)

				if err != nil {
					return nil, fmt.Errorf("Matching error: %s", err)
				}

				if match {
		      to_delete_users[u.Id] = u.Name
		    }
		  }
		} else {
			return nil, fmt.Errorf("User Name must be at least one character long")
		}

	if len(to_delete_users) == 0 {
		return nil, fmt.Errorf("unknown user: %s", name)
	}

	return to_delete_users, nil
}

func DeleteUser(user_id string, token string, keystone_url string) (string, error) {
	request := gorequest.New()
	var errs3 []error
	var status string

	_, status, errs3 = request.Delete("http://"+keystone_url+"/v3/users/"+user_id).
		Set("X-Auth-Token", token).
		Set("Content-type", "application/json").End()

	if errs3 != nil {
		return status, errs3[0]
	}

	if errs3 == nil {
		return status, nil
	}

	return status, nil
}
