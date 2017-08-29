package keystoneauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/google/gofuzz"
)

func TestDomains(t *testing.T) {

	f := fuzz.New()
	var username string
	var description string
	f.Fuzz(&username)
	f.Fuzz(&description)
	token := "7a04a385b907caca141f"
	keystone_url := "localhost:35357"
	enabled := true

	dom, err := CreateDomain(username, description, enabled, token, keystone_url)
	if err != nil {
		return
	}
	domain_id := dom[1]
	domain_name := dom[0]
	assert.Equal(t, domain_name, username)
	assert.NotEqual(t, domain_id, nil)
	assert.Equal(t, err, nil)

	del, err2 := DeleteDomain(domain_id, token, keystone_url)
	if err2 != nil {
		return
	}
	assert.Equal(t, del, "ok")
	assert.Equal(t, domain_name, username)
	assert.Equal(t, err2, nil)
}
