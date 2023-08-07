package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebAppExtensionResource struct{}

type WebAppExtensionModel struct {
	WebAppName        string `tfschema:"web_app_name"`
	ResourceGroup     string `tfschema:"resource_group_name"`
	SiteExtensionName string `tfschema:"site_extension_name"`
	ApiVersion        string `tfschema:"api_version"`
}

var _ sdk.Resource = WebAppExtensionResource{}

func (r WebAppExtensionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"web_app_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"site_extension_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"api_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (r WebAppExtensionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WebAppExtensionResource) ModelObject() interface{} {
	return nil
}

func (r WebAppExtensionResource) ResourceType() string {
	return "azurerm_web_app_extension"
}

func (r WebAppExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppExtensionID
}

func (r WebAppExtensionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webAppExtension WebAppExtensionModel

			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			extensionId := parse.NewWebAppExtensionID(subscriptionId, webAppExtension.ResourceGroup, webAppExtension.WebAppName, webAppExtension.SiteExtensionName)

			// existing, err := client.Get(ctx, webAppExtension.ResourceGroup, webAppExtension.WebAppName, webAppExtension.SiteExtensionName)
			// if err != nil && !utils.ResponseWasNotFound(existing.Response) {
			// 	return fmt.Errorf("checking presence of existing %s: %+v", extensionId, err)
			// }
			future, err := client.InstallSiteExtension(ctx, webAppExtension.ResourceGroup, webAppExtension.WebAppName, webAppExtension.SiteExtensionName)
			if err != nil {
				return fmt.Errorf("creating Extension %s: %+v", extensionId, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Extension %s: %+v", extensionId, err)
			}

			metadata.SetID(extensionId)

			return nil
		},
	}
}

func (r WebAppExtensionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.WebAppExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			webAppExtension, err := client.GetSiteExtension(ctx, id.ResourceGroup, id.SiteName, id.SiteextensionName)
			if err != nil {
				if utils.ResponseWasNotFound(webAppExtension.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Extension %s: %+v", id, err)
			}

			state := WebAppExtensionModel{}
			if props := webAppExtension.SiteExtensionInfoProperties; props != nil {
				state = WebAppExtensionModel{
					WebAppName:        id.SiteName,
					ResourceGroup:     id.ResourceGroup,
					SiteExtensionName: id.SiteextensionName,
					ApiVersion:        pointer.From(webAppExtension.Version),
				}
			}

			state.WebAppName = id.SiteName
			state.ResourceGroup = id.ResourceGroup
			state.SiteExtensionName = id.SiteextensionName
			state.ApiVersion = pointer.From(webAppExtension.Version)

			return nil
		},
	}
}

func (r WebAppExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r WebAppExtensionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}
