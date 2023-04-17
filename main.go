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

	if debug {
		logger.Printf("Client config: %+v", config)
	}

	var v message
	err = json.NewDecoder(in).Decode(&v)
	helpers.FatalIfError(err)
	if debug {
		logger.Printf("Message: %+v\n", v)
	}

	logger.Printf("Attempting to add lifecycle policy to bucket %s in compartment %s\n", v.Data.AdditionalDetails.BucketName, v.Data.CompartmentID)

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

	logger.Printf("Action: %v in %v %v\n", objectLifecyclePolicyDetails.Items[0].Action,
		objectLifecyclePolicyDetails.Items[0].TimeAmount, objectLifecyclePolicyDetails.Items[0].TimeUnit)

	lifecyclePolicyRequest := objectstorage.PutObjectLifecyclePolicyRequest{
		BucketName:                      common.String(v.Data.AdditionalDetails.BucketName),
		NamespaceName:                   common.String(v.Data.AdditionalDetails.Namespace),
		PutObjectLifecyclePolicyDetails: objectLifecyclePolicyDetails,
	}

	if debug {
		logger.Printf("Sending request: %+v\n", lifecyclePolicyRequest)
	} else {
		logger.Printf("Sending Request\n")
	}

	resp, err := client.PutObjectLifecyclePolicy(context.Background(), lifecyclePolicyRequest)
	helpers.FatalIfError(err)

	if debug {
		logger.Printf("Response: %+v\n", resp)
	} else {
		logger.Printf("Response: %v\n", resp.RawResponse.Status)
	}
}
