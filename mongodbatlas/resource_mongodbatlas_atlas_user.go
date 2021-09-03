package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasAtlasUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasAtlasUserCreate,
		ReadContext:   resourceMongoDBAtlasAtlasUserRead,
		DeleteContext: resourceMongoDBAtlasAtlasUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasAtlasUserImportState,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"org_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasAtlasUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	username := ids["username"]

	atlasUser, resp, err := conn.AtlasUsers.GetByName(context.Background(), username)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

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

func resourceMongoDBAtlasAtlasUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	expandedRoles, err := expandAtlasRoles(d)
	if err != nil {
		return diag.FromErr(err)
	}

	AtlasUserReq := &matlas.AtlasUser{
		Roles:    expandedRoles,
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Country:  d.Get("country_code").(string),
	}

	AtlasUserRes, _, err := conn.AtlasUsers.Create(ctx, AtlasUserReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Atlas user: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"username": AtlasUserRes.Username,
	}))

	return resourceMongoDBAtlasAtlasUserRead(ctx, d, meta)
}

func resourceMongoDBAtlasAtlasUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceMongoDBAtlasAtlasUserImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	atlasUser, _, err := conn.AtlasUsers.GetByName(ctx, d.Id())
	if err != nil {
		return nil, fmt.Errorf("couldn't import Atlas user(%s), error: %s", d.Id(), err)
	}

	if err := d.Set("username", atlasUser.Username); err != nil {
		return nil, fmt.Errorf("error setting `username` for Atlas user (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"username": atlasUser.Username,
	}))

	return []*schema.ResourceData{d}, nil
}

func expandAtlasRoles(d *schema.ResourceData) ([]matlas.AtlasRole, error) {
	var roles []matlas.AtlasRole
	proj_roles := []string{
		"GROUP_OWNER",
		"GROUP_CLUSTER_MANAGER",
		"GROUP_READ_ONLY",
		"GROUP_DATA_ACCESS_ADMIN",
		"GROUP_DATA_ACCESS_READ_WRITE",
		"GROUP_DATA_ACCESS_READ_ONLY",
	}
	org_roles := []string{
		"ORG_OWNER",
		"ORG_GROUP_CREATOR",
		"ORG_BILLING_ADMIN",
		"ORG_READ_ONLY",
		"ORG_MEMBER",
	}

	if v, ok := d.GetOk("roles"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			roles = make([]matlas.AtlasRole, rs.Len())

			for k, r := range rs.List() {
				roleMap := r.(map[string]interface{})
				if roleMap["org_id"] != "" {
					if !containedinSlice(org_roles, roleMap["role_name"].(string)) {
						return nil, fmt.Errorf("error creating Atlas user Org role: " + roleMap["role_name"].(string) + " is not a vaild role")
					}
					roles[k] = matlas.AtlasRole{
						RoleName: roleMap["role_name"].(string),
						OrgID:    roleMap["org_id"].(string),
					}
				} else {
					if !containedinSlice(proj_roles, roleMap["role_name"].(string)) {
						return nil, fmt.Errorf("error creating Atlas user Project role: " + roleMap["role_name"].(string) + " is not a vaild role")
					}
					roles[k] = matlas.AtlasRole{
						RoleName: roleMap["role_name"].(string),
						GroupID:  roleMap["project_id"].(string),
					}
				}
			}
		}
	}

	return roles, nil
}

func flattenAtlasRoles(roles []matlas.AtlasRole) []interface{} {
	roleList := make([]interface{}, 0)
	for _, v := range roles {
		if v.OrgID != "" {
			roleList = append(roleList, map[string]interface{}{
				"role_name": v.RoleName,
				"org_id":    v.OrgID,
			})
		} else {
			roleList = append(roleList, map[string]interface{}{
				"role_name":  v.RoleName,
				"project_id": v.GroupID,
			})
		}
	}

	return roleList
}

func containedinSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
