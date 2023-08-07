// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebAppExtensionId struct {
	SubscriptionId    string
	ResourceGroup     string
	SiteName          string
	SiteextensionName string
}

func NewWebAppExtensionID(subscriptionId, resourceGroup, siteName, siteextensionName string) WebAppExtensionId {
	return WebAppExtensionId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		SiteName:          siteName,
		SiteextensionName: siteextensionName,
	}
}

func (id WebAppExtensionId) String() string {
	segments := []string{
		fmt.Sprintf("Siteextension Name %q", id.SiteextensionName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web App Extension", segmentsStr)
}

func (id WebAppExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/siteextensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SiteextensionName)
}

// WebAppExtensionID parses a WebAppExtension ID into an WebAppExtensionId struct
func WebAppExtensionID(input string) (*WebAppExtensionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an WebAppExtension ID: %+v", input, err)
	}

	resourceId := WebAppExtensionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.SiteextensionName, err = id.PopSegment("siteextensions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
