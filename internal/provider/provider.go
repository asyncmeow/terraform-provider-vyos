package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/ganawaj/go-vyos/vyos"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &vyosProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &vyosProvider{
			version: version,
		}
	}
}

// vyosProvider is the provider implementation.
type vyosProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *vyosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vyos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *vyosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host":     schema.StringAttribute{Optional: true},
			"key":      schema.StringAttribute{Optional: true},
			"insecure": schema.BoolAttribute{Optional: true},
		},
	}
}

// Configure prepares a VyOS API client for data sources and resources.
func (p *vyosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config VyosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown VyOS API Host",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API host. "+
				"Either target apply to the source of the value first, set the value statically in the configuration, or use the VYOS_HOST environment variable.",
		)
	}

	if config.Key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown VyOS API Key",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API key. "+
				"Either target apply to the source of the value first, set the value statically in the configuration, or use the VYOS_KEY environment variable.",
		)
	}

	host := os.Getenv("VYOS_HOST")
	key := os.Getenv("VYOS_KEY")
	insecure, err := strconv.ParseBool(os.Getenv("VYOS_INSECURE"))
	if err != nil {
		insecure = false
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Key.IsNull() {
		key = config.Key.ValueString()
	}

	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown VyOS API Host",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API host. "+
				"Either target apply to the source of the value first, set the value statically in the configuration, or use the VYOS_HOST environment variable.",
		)
	}
	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown VyOS API Key",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API key. "+
				"Either target apply to the source of the value first, set the value statically in the configuration, or use the VYOS_KEY environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := vyos.NewClient(nil).WithToken(key).WithURL(host)
	if insecure {
		client = client.Insecure()
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *vyosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEthernetInterfaceDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *vyosProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

type VyosProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Key      types.String `tfsdk:"key"`
	Insecure types.Bool   `tfsdk:"insecure"`
}
