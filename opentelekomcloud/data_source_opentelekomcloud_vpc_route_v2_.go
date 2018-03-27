package opentelekomcloud

import (

	"fmt"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nexthop": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.vpcRouteV2Client(GetRegion(d, config))
	log.Printf("[DEBUG] vpcRouteClient %s", vpcRouteClient)
	listOpts := routes.ListOpts{
		Type:       d.Get("type").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:       d.Get("vpc_id").(string),
		Tenant_Id:     d.Get("tenant_id").(string),
		RouteID: d.Get("id").(string),
	}


	pages, err:= routes.List(vpcRouteClient, listOpts).AllPages()
	refinedRoutes, err := routes.ExtractRoutes(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc routes: %s", err)
	}


	Route := refinedRoutes[0]
	log.Printf("[DEBUG] refinedRoutes %s: %+v", refinedRoutes[0])

	log.Printf("[DEBUG] Retrieved Vpc Routes using given filter %s: %+v", Route.RouteID, Route)
	d.SetId(Route.RouteID)

	d.Set("type", Route.Type)
	d.Set("nexthop", Route.NextHop)
	d.Set("destination", Route.Destination)
	d.Set("tenant_id", Route.Tenant_Id)
	d.Set("vpc_id", Route.VPC_ID)
	d.Set("id", Route.RouteID)
	d.Set("region", GetRegion(d, config))

	return nil
}

