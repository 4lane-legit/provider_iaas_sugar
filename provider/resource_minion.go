package provider

import (
	"fmt"
	"regexp"
	"strings"

	"iaas_sugar/api/cl"
	"iaas_sugar/api/sr"

	"github.com/hashicorp/terraform/helper/schema"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceMinion() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an minion",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceCreateMinion,
		Read:   resourceReadMinion,
		Update: resourceUpdateMinion,
		Delete: resourceDeleteMinion,
		Exists: resourceExistsMinion,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateMinion(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cl.Client)

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	minion := sr.Minion{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        tags,
	}

	err := apiClient.NewMinion(&minion)

	if err != nil {
		return err
	}
	d.SetId(minion.Name)
	return nil
}

func resourceReadMinion(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cl.Client)

	minionId := d.Id()
	minion, err := apiClient.GetMinion(minionId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Minion with ID %s", minionId)
		}
	}

	d.SetId(minion.Name)
	d.Set("name", minion.Name)
	d.Set("description", minion.Description)
	d.Set("tags", minion.Tags)
	return nil
}

func resourceUpdateMinion(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cl.Client)

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	minion := sr.Minion{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        tags,
	}

	err := apiClient.UpdateMinion(&minion)
	if err != nil {
		return err
	}
	return nil
}

func resourceDeleteMinion(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*cl.Client)

	minionId := d.Id()

	err := apiClient.DeleteMinion(minionId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceExistsMinion(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*cl.Client)

	minionId := d.Id()
	_, err := apiClient.GetMinion(minionId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
