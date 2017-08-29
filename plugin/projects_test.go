package keystoneauth

import (
	"github.com/stretchr/testify/assert"
	"github.com/google/gofuzz"
	"testing"
)

func TestProjects(t *testing.T) {

	f := fuzz.New()
	var projectname string
	var description string
	var domain_id string
	var password string
	f.Fuzz(&projectname)
	f.Fuzz(&description)
	f.Fuzz(&domain_id)
	f.Fuzz(&password)

	enabled := true
	is_domain := false
	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"

	ten, err := CreateProject(projectname, description, domain_id, enabled, is_domain, token, keystone_url)
	assert.Equal(t, ten[0], projectname)
	assert.Equal(t, err, nil)
	del, err2 := DeleteProject(ten[1], token, keystone_url)
	assert.Equal(t, del, "ok")
	assert.Equal(t, err2, nil)

}
