package keystoneauth

import (
	"fmt"
	"testing"

	"github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestCreateRegion(t *testing.T) {

	f := fuzz.New()
	var id, description, parentRegionID string
	f.Fuzz(&id)
	//f.Fuzz(&parentRegionID)
	f.Fuzz(&description)
	token := "f5cccfb912a8f814189a"
	keystoneURL := "localhost:35357"

	region, err := CreateRegion(id, description, parentRegionID, keystoneURL, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Equal(t, region[0], id)
	assert.Equal(t, err, nil)
}
