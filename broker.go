package main

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/pivotal-cf/brokerapi"
	"log"
	"net/http"
	"os"
)

type KalturaBroker struct {
}

func NewKalturaBroker() *KalturaBroker {
	broker := &KalturaBroker{}
	return broker
}

func (b *KalturaBroker) Services(ctx context.Context) []brokerapi.Service {
	log.Printf("Got a request to retrieve the catalog")

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "5d9bd115-1b05-4f33-920e-ae9c442c0346",
			Name:        "Kaltura",
			Description: "Create Kaltura account",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{brokerapi.ServicePlan{
				ID:          "7a5ab921-e501-409e-917b-4cb4aa87a782",
				Name:        "default",
				Description: "Kaltura Plan",
			}},
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
