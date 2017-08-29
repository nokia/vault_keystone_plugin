package keystoneauth

import (
	"github.com/stretchr/testify/assert"
	"github.com/google/gofuzz"
	"testing"
)

func TestRoles(t *testing.T) {

	f := fuzz.New()
	var role_name string
  var role_id string
	var user_id string
	var domain_id string
	var description string
	var default_project_id string
	var password string
	var group_id string
	var project_id string

	f.Fuzz(&role_name)
	f.Fuzz(&role_id)
	f.Fuzz(&user_id)
	f.Fuzz(&domain_id)
	f.Fuzz(&description)
	f.Fuzz(&default_project_id)
	f.Fuzz(&password)
	f.Fuzz(&group_id)
	f.Fuzz(&project_id)

	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"

	role,roleerr := CreateRole(role_name, domain_id, token, keystone_url)
	role_name_resp := role[0]
	role_id_resp := role[1]
	assert.Equal(t, role_name, role_name_resp)
	assert.NotEqual(t, role_id_resp, nil)
	assert.Equal(t, roleerr, nil)

	god, errgod := GroupOnDomain(domain_id, group_id, role_id, token, keystone_url)
	assert.Equal(t, god, "ok")
	assert.Equal(t, errgod, nil)

	gop, errgop := GroupOnProject(group_id, project_id, role_id, token, keystone_url)
	assert.Equal(t, gop, "ok")
	assert.Equal(t, errgop, nil)

	uod, erruod := UserOnDomain(user_id, domain_id, role_id, token, keystone_url)
	assert.Equal(t, uod, "ok")
	assert.Equal(t, erruod, nil)

	uop, erruop := UserOnProject(user_id, project_id, role_id, token, keystone_url)
	assert.Equal(t, uop, "ok")
	assert.Equal(t, erruop, nil)
}
