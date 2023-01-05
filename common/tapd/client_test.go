package tapd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestListProject(t *testing.T) {
	companyId := os.Getenv("companyid")
	apiUser := os.Getenv("apiuser")
	apiPassword := os.Getenv("apipassword")
	client := NewClient(companyId, apiUser, apiPassword)
	projects, err := client.ListProject()
	assert.Empty(t, err)
	assert.NotEmpty(t, projects)
}
