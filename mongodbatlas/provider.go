package mongodbatlas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider returns the provider to be use by the code.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_ATLAS_PUBLIC_KEY", ""),
				Description: "MongoDB Atlas Programmatic Public Key",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_ATLAS_PRIVATE_KEY", ""),
				Description: "MongoDB Atlas Programmatic Private Key",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"mongodbatlas_database_user":            dataSourceMongoDBAtlasDatabaseUser(),
			"mongodbatlas_database_users":           dataSourceMongoDBAtlasDatabaseUsers(),
			"mongodbatlas_project":                  dataSourceMongoDBAtlasProject(),
			"mongodbatlas_projects":                 dataSourceMongoDBAtlasProjects(),
			"mongodbatlas_cluster":                  dataSourceMongoDBAtlasCluster(),
			"mongodbatlas_clusters":                 dataSourceMongoDBAtlasClusters(),
			"mongodbatlas_cloud_provider_snapshot":  dataSourceMongoDBAtlasCloudProviderSnapshot(),
			"mongodbatlas_cloud_provider_snapshots": dataSourceMongoDBAtlasCloudProviderSnapshots(),
			"mongodbatlas_network_container":        dataSourceMongoDBAtlasNetworkContainer(),
			"mongodbatlas_network_containers":       dataSourceMongoDBAtlasNetworkContainers(),
			"mongodbatlas_network_peering":          dataSourceMongoDBAtlasNetworkPeering(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"mongodbatlas_database_user":                       resourceMongoDBAtlasDatabaseUser(),
			"mongodbatlas_project_ip_whitelist":                resourceMongoDBAtlasProjectIPWhitelist(),
			"mongodbatlas_project":                             resourceMongoDBAtlasProject(),
			"mongodbatlas_cluster":                             resourceMongoDBAtlasCluster(),
			"mongodbatlas_cloud_provider_snapshot":             resourceMongoDBAtlasCloudProviderSnapshot(),
			"mongodbatlas_network_container":                   resourceMongoDBAtlasNetworkContainer(),
			"mongodbatlas_cloud_provider_snapshot_restore_job": resourceMongoDBAtlasCloudProviderSnapshotRestoreJob(),
			"mongodbatlas_network_peering":                     resourceMongoDBAtlasNetworkPeering(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}
	return config.NewClient(), nil
}
