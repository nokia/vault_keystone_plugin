package keystoneauth

import (
	"github.com/stretchr/testify/assert"
	"github.com/google/gofuzz"
	"testing"
)

func TestUsers(t *testing.T) {


	f := fuzz.New()
	var username string
	var projectname string
	var description string
	var default_project_id string
	var password string

	f.Fuzz(&username)
	f.Fuzz(&projectname)
	f.Fuzz(&description)
	f.Fuzz(&default_project_id)
	f.Fuzz(&password)

	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"
	enabled := true
	is_domain := true

	dom, err := CreateDomain(username, description, enabled, token, keystone_url)
	if err != nil {
		return
	}
	domain_id := dom[1]
	usr, err2 := CreateUser(default_project_id, username, password, enabled, token, domain_id, keystone_url)
	assert.Equal(t, usr[0], username)
	assert.Equal(t, err2, nil)

  ten, err := CreateProject(projectname, description, domain_id, enabled, is_domain, token, keystone_url)
	assert.Equal(t, ten[0], projectname)
	assert.Equal(t, err, nil)


	ec2, ec2err := UserEC2(usr[1], ten[1], token, keystone_url)
	if ec2err != nil {
	  return
	}
  assert.NotEqual(t, ec2[0], "")

	del, err3 := DeleteUser(usr[1], token, keystone_url)
	assert.Equal(t, del, "")
	assert.Equal(t, err3, nil)

}
