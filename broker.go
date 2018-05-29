package main

import (
	"code.cloudfoundry.org/lager"
	"context"
	"encoding/json"
	"errors"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/pivotal-cf/brokerapi"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type KalturaBroker struct {
	ProvisionedInstances map[string]KalturaPartnerProvision
}

func NewKalturaBroker() *KalturaBroker {
	broker := &KalturaBroker{
		ProvisionedInstances: make(map[string]KalturaPartnerProvision),
	}
	return broker
}

func (b *KalturaBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	log.Printf("Got a request to retrieve the catalog")

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "5d9bd115-1b05-4f33-920e-ae9c442c0346",
			Name:        "kaltura",
			Description: "Create Kaltura account",
			Bindable:    true,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName: "kaltura Video Platform as a Service",
				ImageUrl:    "https://vpaas.kaltura.com/images/VPaaS-logo-full.png",
				LongDescription: `
Kaltura VPaaS (Video Platform as a Service) allows you to build any video experience or workflow, and to integrate rich video experiences into existing applications, business workflows and environments.
Kaltura VPaaS eliminates all complexities involved in handling video at scale: ingestion, transcoding, metadata, playback, distribution, analytics, accessibility, monetization, security, search, interactivity and more. 
Available as an open API, with a set of SDKs, developer tools and dozens of code recipes, we’re making the video experience creation process as easy as it gets.`,
				ProviderDisplayName: "Kaltura Inc.",
				DocumentationUrl:    "https://developer.kaltura.com",
				SupportUrl: "https://forum.kaltura.org",
			},
			Plans: []brokerapi.ServicePlan{brokerapi.ServicePlan{
				ID:          "7a5ab921-e501-409e-917b-4cb4aa87a782",
				Name:        "default",
				Description: "Pay As You Go with base REE package. For more details see: https://vpaas.kaltura.com/pricing",
			}},
		},
	}, nil
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
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
}

type KalturaPartnerProvision struct {
	Id          int    `json:"id"`
	AdminSecret string `json:"adminSecret"`

	ObjectType   string `json:"objectType"`
	ErrorMessage string `json:"message"`
}

func (b *KalturaBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	log.Printf("Got a request to provision instanceId: %v\n", instanceID)
	retval := brokerapi.ProvisionedServiceSpec{
		DashboardURL: "https://kmc.kaltura.com/index.php/kmcng/login",
	}

	var params ProvisionParameters
	err := json.Unmarshal(details.RawParameters, &params)
	if err != nil {
		return retval, err
	}
	if params.Company == "" || params.Email == "" || params.Name == "" {
		return retval, errors.New("Missing parameters")
	}
	values := url.Values{}
	values.Add("partner[objectType]", "KalturaPartner")
	values.Add("partner[description]", "SAP Cloud Platform provisioned")
	values.Add("partner[name]", params.Company)
	values.Add("partner[adminName]", params.Name)
	values.Add("partner[adminEmail]", params.Email)
	values.Add("partner[referenceId]", instanceID)
	values.Add("format", "1")
	resp, err := http.PostForm("https://www.kaltura.com/api_v3/service/partner/action/register", values)
	if err != nil {
		return retval, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Received a return value of %s\n", string(body))

	var kalturaResponse KalturaPartnerProvision
	err = json.Unmarshal(body, &kalturaResponse)
	if err != nil {
		return retval, err
	}
	if kalturaResponse.ObjectType == "KalturaAPIException" {
		return retval, errors.New(kalturaResponse.ErrorMessage)
	}
	log.Printf("Received a return value of %v\n", kalturaResponse)
	b.ProvisionedInstances[instanceID] = kalturaResponse

	return retval, nil
}

func (b *KalturaBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	log.Printf("Got a request to deprovision a %v for instanceId: %v\n", details, instanceID)
	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *KalturaBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	log.Printf("Got a request to bind bindingId %v for instanceId: %v\n", bindingID, instanceID)
	resp, ok := b.ProvisionedInstances[instanceID]
	if !ok {
		return brokerapi.Binding{}, errors.New("No such instance")
	}
	return brokerapi.Binding{
		Credentials: map[string]interface{}{
			"adminSecret": resp.AdminSecret,
			"partnerId":   resp.Id,
		},
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
