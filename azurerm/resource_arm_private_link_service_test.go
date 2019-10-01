package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateLinkService_basic(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_frontend_ip_configuration_ids.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_update(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_frontend_ip_configuration_ids.0"),
				),
			},
			{
				Config: testAccAzureRMPrivateLinkService_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "www.contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_frontend_ip_configuration_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_frontend_ip_configuration_ids.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_complete(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_approval_subscription_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "auto_approval_subscription_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_subscription_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "visibility_subscription_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "www.contoso.com"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.1.18"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.primary", "false"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_frontend_ip_configuration_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMPrivateLinkServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Link Service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Private Link Service %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateLinkServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_link_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateLinkService_basic(rInt int, location string) string {
	standardResources := testAccAzureRMPrivateLinkServiceTemplate_standardResources(rInt, location)
	privateLink := testAccAzureRMPrivateLinkServiceTemplate_basic(rInt)

	return testAccAzureRMPrivateLinkServiceTemplate("", standardResources, privateLink)
}

func testAccAzureRMPrivateLinkService_update(rInt int, location string) string {
	standardResources := testAccAzureRMPrivateLinkServiceTemplate_standardResources(rInt, location)
	privateLink := testAccAzureRMPrivateLinkServiceTemplate_update(rInt)

	return testAccAzureRMPrivateLinkServiceTemplate("", standardResources, privateLink)
}

func testAccAzureRMPrivateLinkService_complete(rInt int, location string) string {
	subscriptionDataSource := testAccAzureRMPrivateLinkServiceTemplate_subscriptionDataSource()
	standardResources := testAccAzureRMPrivateLinkServiceTemplate_standardResources(rInt, location)
	privateLink := testAccAzureRMPrivateLinkServiceTemplate_complete(rInt)

	return testAccAzureRMPrivateLinkServiceTemplate(subscriptionDataSource, standardResources, privateLink)
}

func testAccAzureRMPrivateLinkServiceTemplate(subscriptionDataSourceTemplate string, standardResourcesTemplate string, privateLinkTemplate string) string {
	return fmt.Sprintf(`
%s

%s

%s
`,subscriptionDataSourceTemplate, standardResourcesTemplate, privateLinkTemplate)
}

func testAccAzureRMPrivateLinkServiceTemplate_subscriptionDataSource() string {
	return `data "azurerm_subscription" "current" {}`
}

func testAccAzureRMPrivateLinkServiceTemplate_standardResources(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.1.0/24"
  private_link_service_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`,rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkServiceTemplate_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_private_link_service" "test" {
  name                           = "acctestpls-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, rInt, rInt)
}

func testAccAzureRMPrivateLinkServiceTemplate_update(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_private_link_service" "test" {
  name                           = "acctestpls-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  fqdns                          = ["www.contoso.com"]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, rInt, rInt)
}

func testAccAzureRMPrivateLinkServiceTemplate_complete(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_private_link_service" "test" {
  name                           = "acctestpls-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  fqdns                          = ["www.contoso.com"]
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.1.17"
    private_ip_address_version   = "IPv4"
    private_ip_allocation_method = "Static"
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.1.18"
    private_ip_address_version   = "IPv4"
    private_ip_allocation_method = "Static"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, rInt, rInt, rInt)
}