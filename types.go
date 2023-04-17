package main

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
