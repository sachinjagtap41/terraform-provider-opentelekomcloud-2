package opentelekomcloud
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"
)

// PASS
func TestAccOTCRouteV2_basic(t *testing.T) {
	var route routes.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRouteV2Exists("opentelekomcloud_vpc_route_v2.route_1", &route),
				),
			},
		},
	})
}

// PASS
func TestAccOTCRouteV2_timeout(t *testing.T) {
	var route routes.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRouteV2Exists("opentelekomcloud_vpc_route_v2.route_1", &route),
				),
			},
		},
	})
}

func testAccCheckOTCRouteV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	routeClient, err := config.vpcRouteV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud route client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_v2" {
			continue
		}

		_, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

func testAccCheckOTCRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		routeClient, err := config.vpcRouteV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud route client: %s", err)
		}

		found, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

const testAccRouteV2_basic = `
resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type = "peering"
  nexthop = "7b0cf30f-b5c4-4cc2-979c-e1a964350467"
  destination = "192.168.0.0/16"
  vpc_id ="3127e30b-5f8e-42d1-a3cc-fdadf412c5bf"

}
`

const testAccRouteV2_timeout = `
resource "opentelekomcloud_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = "7b0cf30f-b5c4-4cc2-979c-e1a964350467"
  destination = "192.168.0.0/16"
  vpc_id ="3127e30b-5f8e-42d1-a3cc-fdadf412c5bf"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

