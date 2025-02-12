package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IoTCentralApplicationResource struct{}

func TestAccIoTCentralApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("ST1"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("template").HasValue("iotc-pnp-preview@1.0.0"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("ST1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").DoesNotExist(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_publicNetworkAccessEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithPublicNetworkAccessEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithPublicNetworkAccessEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (IoTCentralApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apps.ParseIotAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTCentral.AppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (IoTCentralApplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  sku                 = "ST1"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationResource) basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  sku                 = "ST1"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationResource) basicWithPublicNetworkAccessEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  sku                 = "ST1"

  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  display_name        = "some-display-name"
  sku                 = "ST1"
  template            = "iotc-pnp-preview@1.0.0"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sub_domain          = "subdomain-%[1]d"
  display_name        = "some-display-name"
  sku                 = "ST1"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IoTCentralApplicationResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application" "import" {
  name                = azurerm_iotcentral_application.test.name
  resource_group_name = azurerm_iotcentral_application.test.resource_group_name
  location            = azurerm_iotcentral_application.test.location
  sub_domain          = azurerm_iotcentral_application.test.sub_domain
  display_name        = azurerm_iotcentral_application.test.display_name
  sku                 = azurerm_iotcentral_application.test.sku
}
`, template)
}
