// This function example is provided as-is with no warranty of any kind
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"

	fdk "github.com/fnproject/fdk-go"
)

var (
	debug  bool
	logger *log.Logger
)

type message struct {
	EventTime string `json:"eventTime"`
	Source    string `json:"source"`
	Data      data   `json:"data"`
}

type data struct {
	CompartmentID      string  `json:"compartmentId"`
	CompartmentName    string  `json:"compartmentName"`
	ResourceName       string  `json:"resourceName"`
	ResourceID         string  `json:"resourceId"`
	AvailabilityDomain string  `json:"availabilityDomain"`
	AdditionalDetails  details `json:"additionalDetails"`
}

type details struct {
	BucketName       string `json:"bucketName"`
	PublicAccessType string `json:"publicAccessType"`
	Versioning       string `json:"versioning"`
	Namespace        string `json:"namespace"`
	ETag             string `json:"eTag"`
}

func init() {
	// Ensure logger is writing to std out
	logger = log.New(os.Stdout, "", log.LstdFlags)
	_, debug = os.LookupEnv("DEBUG")
}

func main() {
	if debug {
		log.Printf("Value of variable debug: %v\n", debug)
	}
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	logger.Printf("Context data: %v\n", ctx)
	// Try resource principal configuration and log before exit 1 if error
	config, err := auth.ResourcePrincipalConfigurationProvider()
	helpers.FatalIfError(err)

	var v message
	err = json.NewDecoder(in).Decode(&v)
	helpers.FatalIfError(err)
	if debug {
		logger.Printf("Message: %+v\n", v)
	}

	logger.Printf("Attempting to add lifecycle policy to bucket %s\n", v.Data.AdditionalDetails.BucketName)
	logger.Printf("{\neventTime: %v\ncompartmentName: %v,\ncompartmentId: %v\navailabilityDomain: %v,\neTag: %v\n}",
		v.EventTime, v.Data.CompartmentName, v.Data.CompartmentID,
		v.Data.AvailabilityDomain, v.Data.AdditionalDetails.ETag)

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)
	helpers.FatalIfError(err)

	objectLifecyclePolicyDetails := objectstorage.PutObjectLifecyclePolicyDetails{
		Items: []objectstorage.ObjectLifecycleRule{
			objectstorage.ObjectLifecycleRule{
				Name:       common.String("DefaultLifecycleRule"),
				Action:     common.String("ARCHIVE"),
				TimeAmount: common.Int64(31),
				TimeUnit:   objectstorage.ObjectLifecycleRuleTimeUnitDays,
				IsEnabled:  common.Bool(true),
				// Target: common.String("objects")
			},
		},
	}

	lifecyclePolicyRequest := objectstorage.PutObjectLifecyclePolicyRequest{
		BucketName:                      common.String(v.Data.AdditionalDetails.BucketName),
		NamespaceName:                   common.String(v.Data.AdditionalDetails.Namespace),
		PutObjectLifecyclePolicyDetails: objectLifecyclePolicyDetails,
	}

	logger.Printf("Sending request: %+v\n", lifecyclePolicyRequest)

	resp, err := client.PutObjectLifecyclePolicy(context.Background(), lifecyclePolicyRequest)
	helpers.FatalIfError(err)

	logger.Printf("Response: %+v\n", resp)
}
