package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasAtlasUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAtlasUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"org_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAtlasUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	username := d.Get("username").(string)

	atlasUser, _, err := conn.AtlasUsers.GetByName(ctx, username)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Atlas user information: %s", err))
	}

	if err := d.Set("username", atlasUser.Username); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `username` for Atlas user (%s): %s", d.Id(), err))
	}

	if err := d.Set("country_code", atlasUser.Country); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `country_code` for Atlas user (%s): %s", d.Id(), err))
	}

	if err := d.Set("email_address", atlasUser.EmailAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `email_address` for Atlas user (%s): %s", d.Id(), err))
	}

	if err := d.Set("first_name", atlasUser.FirstName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `first_name` for Atlas user (%s): %s", d.Id(), err))
	}

	if err := d.Set("last_name", atlasUser.LastName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `last_name` for Atlas user (%s): %s", d.Id(), err))
	}

	if err := d.Set("roles", flattenAtlasRoles(atlasUser.Roles)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `roles` for Atlas user (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username": username,
	}))

	return nil
}