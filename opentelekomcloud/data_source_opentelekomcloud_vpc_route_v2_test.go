
package opentelekomcloud

import (
	"fmt"
	"testing"


	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOTCRouteV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCRouteV2DataSource_vpcroute,
			},
			resource.TestStep{
				Config: testAccOTCRouteV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCRouteV2DataSourceID("data.opentelekomcloud_vpc_route_v2.route_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.route_1", "vpc_id", "ff45ad82-27de-4a69-bedb-8118f963d82b"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.route_1", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.route_1", "destination", "192.168.0.0/17"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.route_1", "nexthop", "2e0e9c0c-f3c9-4b05-a7ed-aac1347bc290"),

				),
			},
		},
	})
}

// PASS
func TestAccOTCRouteV2DataSource_vpcRouteID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCRouteV2DataSource_vpcroute,
			},
			resource.TestStep{
				Config: testAccOTCRouteV2DataSource_vpcRouteID,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCRouteV2DataSourceID("data.opentelekomcloud_vpc_route_v2.route_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.route_1", "vpc_id", "ff45ad82-27de-4a69-bedb-8118f963d82b"),
				),
			},
		},
	})
}

func testAccOTCRouteV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc route data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Route data source ID not set")
		}

		return nil
	}
}

const testAccOTCRouteV2DataSource_vpcroute = `
resource "opentelekomcloud_vpc_route_v2" "route_1" {  
	destination = "192.168.0.0/17"
	nexthop = "2e0e9c0c-f3c9-4b05-a7ed-aac1347bc290"
	vpc_id = "ff45ad82-27de-4a69-bedb-8118f963d82b"
	type = "peering"
}
`

var testAccOTCRouteV2DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_vpc_route_v2" "route_1" {
	vpc_id = "ff45ad82-27de-4a69-bedb-8118f963d82b"
}
`, testAccOTCRouteV2DataSource_vpcroute)

var testAccOTCRouteV2DataSource_vpcRouteID = fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_route_v2" "route_1" {
	vpc_id = "${opentelekomcloud_vpc_route_v2.route_1.id}"
}
`, testAccOTCRouteV2DataSource_vpcroute)


