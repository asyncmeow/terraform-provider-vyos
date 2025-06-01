// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ganawaj/go-vyos/vyos"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-vyos/internal/provider/vyos_models"
)

var (
	_ datasource.DataSource              = &ethernetDataSource{}
	_ datasource.DataSourceWithConfigure = &ethernetDataSource{}
)

func NewEthernetInterfaceDataSource() datasource.DataSource {
	return &ethernetDataSource{}
}

type ethernetDataSource struct {
	client *vyos.Client
}

func (d *ethernetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vyos.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected *vyos.Client, got %T. Please report this issue to the provider developers.", req.ProviderData))

		return
	}

	d.client = client
}

func (d *ethernetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ethernet_interface"
}

func (d *ethernetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":      schema.StringAttribute{Required: true},
			"addresses": schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"hw_id":     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *ethernetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var inputModel ethernetInterfaceDataSourceModel
	diags := req.Config.Get(ctx, &inputModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ethernet_interface", inputModel.Name.ValueString())
	tflog.Info(ctx, "fetching ethernet interface")
	out, _, err := d.client.Conf.Get(ctx, "interfaces ethernet "+inputModel.Name.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read ethernet interface", err.Error())
		return
	}
	tflog.Debug(ctx, "got http response", map[string]any{"response": out})
	jsonString, _ := json.Marshal(out.Data)
	tflog.Trace(ctx, "re-marshalled response", map[string]any{"json_str": jsonString})
	stateJson := vyos_models.InterfacesEthernet{}
	err = json.Unmarshal(jsonString, &stateJson)
	if err != nil {
		resp.Diagnostics.AddError("failed to unmarshal response", err.Error())
		return
	}
	tflog.Trace(ctx, "un-marshalled response", map[string]any{"response": stateJson})

	addresses, diags := types.ListValueFrom(ctx, types.StringType, stateJson.Addresses)
	resp.Diagnostics.Append(diags...)
	state := ethernetInterfaceDataSourceModel{
		Address: addresses,
		HwId:    types.StringValue(stateJson.HwId),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type ethernetInterfaceDataSourceModel struct {
	Name    types.String `tfsdk:"name"`
	Address types.List   `tfsdk:"addresses"`
	HwId    types.String `tfsdk:"hw_id"`
}
