package keystoneauth

import (
	"github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCredentials(t *testing.T) {

	f := fuzz.New()
	var user_id string
	var password string
	var username string
	var description string
	var default_project_id string
	f.Fuzz(&username)
	f.Fuzz(&description)
	f.Fuzz(&default_project_id)
	f.Fuzz(&password)

	var projectname string
	var domain_id string
	f.Fuzz(&projectname)
	f.Fuzz(&domain_id)

	enabled := true
	is_domain := false
	thetype := "ec2"
	blob := "{\"access\":\"181920\",\"secret\":\"secretKey\"}"
	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"

	dom, err := CreateDomain(username, description, enabled, token, keystone_url)
	if err != nil {
		return
	}
	domain_id = dom[1]
	usr, err2 := CreateUser(default_project_id, username, password, enabled, token, domain_id, keystone_url)
	assert.Equal(t, usr[0], username)

	ten, err := CreateProject(projectname, description, domain_id, enabled, is_domain, token, keystone_url)
	assert.Equal(t, ten[0], projectname)

	cre, err2 := CreateCredential(blob, thetype, usr[1], ten[1], token, keystone_url)
	assert.Equal(t, cre[1], user_id)
	assert.Equal(t, err2, nil)
	del, err3 := DeleteCredential(cre[0], token, keystone_url)
	assert.Equal(t, del, "ok")
	assert.Equal(t, err3, nil)

}
