package main

import (
	"context"
	"github.com/goji/httpauth"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/liorokman/brokerapi"
	"log"
	"net/http"
)

type KalturaBroker struct {
}

func NewKalturaBroker() *KalturaBroker {
	broker := &KalturaBroker{}
	return broker
}

func (b *KalturaBroker) Services(ctx context.Context) []brokerapi.Service {
	bindingParams := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Human readable name of the remote service being exposed",
				"title":       "Name",
			},
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Remote GUID of the service being exposed/unexposed",
				"title":       "GUID",
			},
			"credentials": map[string]interface{}{
				"type":        "object",
				"description": "Opaque JSON document with credentials to be given to binding instances for this service",
				"title":       "Credentials",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Description of the service being exposed",
				"title":       "Description",
			},
		},
		"required": []string{"id", "name", "description"},
	}
	log.Printf("Got a request to retrieve the catalog")
	planList := []brokerapi.ServicePlan{brokerapi.ServicePlan{
		ID:          "an id",
		Name:        "small-org",
		Description: "Creates a Small CF Org in CloudFoundry",
		Schemas: &brokerapi.ServiceSchemas{
			Instance: brokerapi.ServiceInstanceSchema{
				Create: brokerapi.Schema{
					Schema: map[string]interface{}{
						"$schema": "http://json-schema.org/draft-06/schema#",
						"type":    "object",
						"properties": map[string]interface{}{
							"email": map[string]interface{}{
								"type":        "string",
								"title":       "User Email",
								"description": "The email of the user that already exists in SAP IDP that will be added as an org manager in the created org",
								"default":     "${admin_email}",
							},
						},
						"required": []string{"email"},
					},
				},
				Update: brokerapi.Schema{
					Schema: map[string]interface{}{
						"$schema": "http://json-schema.org/draft-06/schema#",
						"type":    "object",
						"properties": map[string]interface{}{
							"del_bindings": map[string]interface{}{
								"type":  "array",
								"title": "Deleted bindings",
								"items": bindingParams,
							},
							"add_bindings": map[string]interface{}{
								"type":  "array",
								"title": "Added bindings",
								"items": bindingParams,
							},
						},
					},
				},
			},
		},
	}}

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "another id",
			Name:        "CF Org",
			Description: "Create an Org in CF",
			Bindable:    false,
			Metadata:    map[string]interface{}{},
			Plans:       planList,
		},
	}
}

/*
func (b *KalturaBroker) getXuaaToken(xuaaTenantOnBoardingUserName, xuaaTenantOnBoardingPassword string) (string, error) {
	payload := strings.NewReader("grant_type=password&username=" + url.QueryEscape(xuaaTenantOnBoardingUserName) + "&password=" + url.QueryEscape(xuaaTenantOnBoardingPassword) + "&scope=")
	req, _ := http.NewRequest("POST", b.UaaTokenEndpoint+"/oauth/token", payload)
	req.Header.Add("authorization", "Basic eHVhYTo=")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	tokenResponse := GetXuaaTokenResponse{}
	err := json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

*/
type ProvisionParameters struct {
	Email string `json:"email"`
}

func (b *KalturaBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	log.Printf("Got a request to provision a %v for instanceId: %v\n", details, instanceID)
	retval := brokerapi.ProvisionedServiceSpec{}

	return retval, nil
}
func (b *KalturaBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	log.Printf("Got a request to deprovision a %v for instanceId: %v\n", details, instanceID)
	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b KalturaBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	log.Printf("Got a request to bind bindingId %v for instanceId: %v\n", bindingID, instanceID)
	return brokerapi.Binding{
		Credentials: map[string]interface{}{},
	}, nil
}

func (b KalturaBroker) Unbind(ctx context.Context, instanceId, bindingID string, details brokerapi.UnbindDetails) error {
	log.Printf("Got a request to unbind bindingId %v for instanceId: %v\n", bindingID, instanceId)
	return nil
}

func (b *KalturaBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

func (b *KalturaBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	brokerLogger := lager.NewLogger("broker")
	KalturaBroker := NewKalturaBroker()
	router.Use(httpauth.SimpleBasicAuth("user", "pass"))
	brokerapi.AttachRoutes(router, KalturaBroker, brokerLogger)

	//add authentication for broker paths
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))

}
