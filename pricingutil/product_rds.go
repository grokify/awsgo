// pricingutil provices some helpers for the AWS Pricing API: https://docs.aws.amazon.com/sdk-for-go/api/service/pricing/
package pricingutil

import "time"

type ProductRDS struct {
	Product         Product   `json:"product"`
	PublicationDate time.Time `json:"publicationDate"`
	Version         string    `json:"version"`
}

type Product struct {
	Attributes    map[string]string `json:"attributes"`
	ProductFamily string            `json:"productFamily"`
	SKU           string            `json:"sku"`
}
