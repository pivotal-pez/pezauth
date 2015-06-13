package integrations_test

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type mockVcapMongo struct {
	URI string
}

func setMongoEnv(ip, port string) {
	mockVcap := mockVcapMongo{
		URI: fmt.Sprintf("mongodb://%s:%s/59501f9b-5528-421b-84e9-93f9dc6f1080", ip, port),
	}
	vcap_services := `{"p-mongodb": [{"name": "pezauth-mongo","label": "p-mongodb","tags": ["pivotal","mongodb"],"plan": "development","credentials": {"uri": "{{.URI}}"}}]}`
	mongoTmpl, _ := template.New("mongoVcap").Parse(vcap_services)
	var b bytes.Buffer
	mongoTmpl.Execute(&b, mockVcap)
	os.Setenv("VCAP_SERVICES", b.String())
	os.Setenv("MONGO_SERVICE_NAME", "pezauth-mongo")
	os.Setenv("MONGO_URI_NAME", "uri")
	os.Setenv("MONGO_COLLECTION_NAME", "org_users")
}
