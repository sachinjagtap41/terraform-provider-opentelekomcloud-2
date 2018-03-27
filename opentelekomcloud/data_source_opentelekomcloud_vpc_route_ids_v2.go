package opentelekomcloud

import (
	//"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func dataSourceVPCRouteIdsV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteIdsV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVpcRouteIdsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.vpcRouteV2Client(GetRegion(d, config))

	listOpts := routes.ListOpts{
		VPC_ID:   d.Get("vpc_id").(string),
	}

	pages, err:= routes.List(vpcRouteClient, listOpts).AllPages()
	refinedRoutes, err := routes.ExtractRoutes(pages)

	log.Printf("[DEBUG] Value of allRoutes: %#v", refinedRoutes)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc Routes: %s", err)
	}

	//Route := refinedRoutes[0]
	listRoutes := make([]string, 0)

	for _, route := range refinedRoutes {
		listRoutes = append(listRoutes, route.RouteID)

	}
	log.Printf("[DEBUG] listRoutes %s", listRoutes)

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", listRoutes)
	d.Set("region", GetRegion(d, config))

	return nil
}

