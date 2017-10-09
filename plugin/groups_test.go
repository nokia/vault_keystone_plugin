package keystoneauth

import (
	"github.com/stretchr/testify/assert"
	"github.com/google/gofuzz"
	"testing"
)

func TestGroups(t *testing.T) {

	f := fuzz.New()
	var name string
	var description string
	var domain_id string
	f.Fuzz(&name)
	f.Fuzz(&description)
	f.Fuzz(&domain_id)
	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"

	grp, err := CreateGroup(name, description, domain_id, token, keystone_url)
	if err != nil {
		return
	}
	group_id := grp[1]

	del, err3 := DeleteGroup(group_id, token, keystone_url)
	assert.Equal(t, del, "ok")
	assert.Equal(t, err3, nil)

}
