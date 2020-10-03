package provider

import (
	"iaas_sugar/api/cl"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

/* Provider
A Schema which represents the various attributes we can provide to our provider via the provider block of a Terraform file.
Note that if no value is provided we will check if environment variables are set.
This is useful for making sure we don’t need to store secrets in the provider block of terraform files

ResourceMap defines the names of the resources the provider has and where to find the
definition of those resources.
In this case, you can see we have example_item resource, the definition of which is a *schema.Resource
returned by the resourceItem() function,

ConfigureFunc which can do any setup for us. In this case, we have providerConfigure which takes the address,
port and token and returns a client that we’ll use to communicate with the API. Note the providerConfigure returns an interface{}
so we can store anything we like here
*/
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_ADDRESS", ""),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_PORT", ""),
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_item": resourceMinion(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	secret := d.Get("secret").(string)
	return cl.NewClient(host, port, secret), nil
}
