// This function example is provided as-is with no warranty of any kind
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"

	fdk "github.com/fnproject/fdk-go"
)

var (
	logger     *log.Logger
	errLog     *log.Logger
	debug      bool
	action     string
	timeamount string
)

func init() {
	// Ensure logger is writing to Stdout
	// Stdout is preferred by OCI Functions
	logger = log.New(os.Stdout, "", log.LstdFlags)
	errLog = log.New(os.Stdout, "ERR: ", log.Lshortfile)

	// Environment variable capture
	_, debug = os.LookupEnv("DEBUG")

	var ok bool
	// Options "ARCHIVE", "INFREQUENT_ACCESS", "DELETE", & "ABORT"
	action, ok = os.LookupEnv("ACTION")
	if !ok {
		action = "ARCHIVE"
	}

	// Default to 31 days
	timeamount, ok = os.LookupEnv("TIMEAMOUNT")
	if !ok {
		timeamount = "31"
	}
}

func main() {
	if debug {
		log.Printf("Value of environment variables: {debug: %v, action: %v, timeamount: %v\n",
			debug, action, timeamount)
	}

	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	if debug {
		logger.Printf("Context data: %v\n", ctx)
	}

	// Convert timeamount to usable int64
	time, err := strconv.ParseInt(timeamount, 10, 64)
	logFatal(err)

	// Try resource principal configuration and log before exit 1 if error
	config, err := auth.ResourcePrincipalConfigurationProvider()
	logFatal(err)

	if debug {
		logger.Printf("Client config: %+v", config)
	}

	var v message
	err = json.NewDecoder(in).Decode(&v)
	logFatal(err)
	if debug {
		logger.Printf("Message: %+v\n", v)
	}

	logger.Printf("Attempting to add lifecycle policy to bucket %s in compartment %s\n", v.Data.AdditionalDetails.BucketName, v.Data.CompartmentID)

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)
	logFatal(err)

	// Set lifecycle rule with duration time and action to enforce rules on object storage
	// buckets.
	objectLifecyclePolicyDetails := objectstorage.PutObjectLifecyclePolicyDetails{
		Items: []objectstorage.ObjectLifecycleRule{
			{
				Name:       common.String("DefaultLifecycleRule"),
				Action:     common.String(action),
				TimeAmount: common.Int64(time),
				TimeUnit:   objectstorage.ObjectLifecycleRuleTimeUnitDays,
				IsEnabled:  common.Bool(true),
				// Target: common.String("objects")
			},
		},
	}

	logger.Printf("Action: %v in %v %v\n", *objectLifecyclePolicyDetails.Items[0].Action,
		*objectLifecyclePolicyDetails.Items[0].TimeAmount, objectLifecyclePolicyDetails.Items[0].TimeUnit)

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
	logFatal(err)

	if debug {
		logger.Printf("Response: %+v\n", resp)
	} else {
		logger.Printf("Response: %v\n", resp.RawResponse.Status)
	}
}

func logFatal(err error) {
	if err != nil {
		errLog.Fatalln(err)
	}
}
