// This function example is provided as-is with no warranty of any kind
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/oracle/oci-go-sdk/v65/example/helpers"

	fdk "github.com/fnproject/fdk-go"
)

var debug bool

type Message struct {
	EventTime string `json:"eventTime"`
	Source    string `json:"source"`
	Data      Data   `json:"data"`
}

type Data struct {
	CompartmentId     string  `json:"compartmentId"`
	CompartmentName   string  `json:"compartmentName"`
	ResourceName      string  `json:"resourceName"`
	ResourceId        string  `json:"resourceId"`
	AdditionalDetails Details `json:"additionalDetails"`
}

type Details struct {
	BucketName       string `json:"bucketName"`
	PublicAccessType string `json:"publicAccessType"`
	Versioning       string `json:"versioning"`
	Namespace        string `json:"namespace"`
	ETag             string `json:"eTag"`
}

func init() {
	_, debug = os.LookupEnv("DEBUG")
}

func main() {
	if debug {
		fmt.Printf("Value of variable debug: %v\n", debug)
	}
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	fmt.Printf("Context data: %v\n", ctx)
	// Try resource principal configuration and log before exit 1 if error
	//config, err := auth.ResourcePrincipalConfigurationProvider()
	//helpers.FatalIfError(err)

	var v Message
	err := json.NewDecoder(in).Decode(&v)
	helpers.FatalIfError(err)
	if debug {
		fmt.Printf("Message: %+v\n", v)
	}

	fmt.Printf("Attempting to add lifecycle policy to bucket %s\n", v.Data.AdditionalDetails.BucketName)
}
