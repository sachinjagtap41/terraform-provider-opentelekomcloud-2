package opentelekomcloud

import (
	"fmt"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/hashicorp/terraform/helper/resource"
)

func resourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcRouteV2Create, //providers.go
		Read:   resourceVpcRouteV2Read,
		Delete: resourceVpcRouteV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"nexthop": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required:     true,
				ForceNew: true,
				//Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required:     true,
				ForceNew: true,

			},
		},
	}
}

func resourceVpcRouteV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.vpcRouteV2Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of vpcRouteClient: %#v", vpcRouteClient)

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc route client: %s", err)
	}

	createOpts := routes.CreateOpts{
		Type: d.Get("type").(string),
		NextHop: d.Get("nexthop").(string),
		Destination: d.Get("destination").(string),
		Tenant_Id: d.Get("tenant_id").(string),
		VPC_ID: d.Get("vpc_id").(string),

	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := routes.Create(vpcRouteClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud VPC route: %s", err)
	}
	d.SetId(n.RouteID)

	log.Printf("[INFO] Vpc Route ID: %s", n.RouteID)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Vpc route (%s) to become available", n.RouteID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcRouteActive(vpcRouteClient, n.RouteID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(n.RouteID)

	return resourceVpcRouteV2Read(d, meta)

}

func resourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.vpcRouteV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Vpc route client: %s", err)
	}

	n, err := routes.Get(vpcRouteClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Vpc route: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Vpc Route %s: %+v", d.Id(), n)

	d.Set("type", n.Type)
	d.Set("nexthop", n.NextHop)
	d.Set("destination", n.Destination)
	d.Set("tenant_id", n.Tenant_Id)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("id", n.RouteID)
	d.Set("region", GetRegion(d, config))

	return nil
}


func resourceVpcRouteV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy vpc route: %s", d.Id())

	config := meta.(*Config)
	vpcRouteClient, err := config.vpcRouteV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc route: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcRouteDelete(vpcRouteClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Vpc route: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcRouteActive(vpcRouteClient *golangsdk.ServiceClient, routeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := routes.Get(vpcRouteClient, routeId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud VPC Route Client: %+v", n)


		return n, n.RouteID,nil
	}
}

func waitForVpcRouteDelete(vpcRouteClient *golangsdk.ServiceClient, routeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud vpc route %s.\n", routeId)

		r, err := routes.Get(vpcRouteClient, routeId).Extract()
		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud vpc route %s", routeId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = routes.Delete(vpcRouteClient, routeId).ExtractErr()
		log.Printf("[DEBUG] Value if error: %#v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud vpc route %s", routeId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud vpc route %s still active.\n", routeId)
		return r, "ACTIVE", nil
	}
}
