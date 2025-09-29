package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OperatingSystemsDataSource{}

func NewOperatingSystemsDataSource() datasource.DataSource {
	return &OperatingSystemsDataSource{}
}

// OperatingSystemsDataSource defines the data source implementation.
type OperatingSystemsDataSource struct {
	client *ICSClient
}

// OperatingSystemsDataSourceModel describes the data source data model.
type OperatingSystemsDataSourceModel struct {
	ServerTypeName types.String                    `tfsdk:"server_type_name"`
	Location       types.String                    `tfsdk:"location"`
	OperatingSystems []OperatingSystemDataModel    `tfsdk:"operating_systems"`
	ID             types.String                    `tfsdk:"id"`
}

type OperatingSystemDataModel struct {
	Name          types.String  `tfsdk:"name"`
	OSType        types.String  `tfsdk:"os_type"`
	ProductCode   types.String  `tfsdk:"product_code"`
	Price         types.Float64 `tfsdk:"price"`
	PriceHourly   types.Float64 `tfsdk:"price_hourly"`
	HourlyEnabled types.Bool    `tfsdk:"hourly_enabled"`
}

func (d *OperatingSystemsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operating_systems"
}

func (d *OperatingSystemsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Operating systems data source provides information about available operating systems for a specific server type and location.",

		Attributes: map[string]schema.Attribute{
			"server_type_name": schema.StringAttribute{
				MarkdownDescription: "Server type name (e.g., 'c1i.small')",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location code (e.g., 'NYC1')",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier",
				Computed:            true,
			},
			"operating_systems": schema.ListNestedAttribute{
				MarkdownDescription: "List of available operating systems",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Operating system name",
							Computed:            true,
						},
						"os_type": schema.StringAttribute{
							MarkdownDescription: "Operating system type (e.g., 'linux', 'windows')",
							Computed:            true,
						},
						"product_code": schema.StringAttribute{
							MarkdownDescription: "Product code used by the API",
							Computed:            true,
						},
						"price": schema.Float64Attribute{
							MarkdownDescription: "Monthly price",
							Computed:            true,
						},
						"price_hourly": schema.Float64Attribute{
							MarkdownDescription: "Hourly price",
							Computed:            true,
						},
						"hourly_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether hourly billing is available",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OperatingSystemsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ICSClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ICSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *OperatingSystemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OperatingSystemsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverTypeName := data.ServerTypeName.ValueString()
	location := data.Location.ValueString()

	// Get addons (which includes operating systems) from API
	addons, err := d.client.GetAddons(serverTypeName, location)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read operating systems for server type '%s' in location '%s', got error: %s", serverTypeName, location, err))
		return
	}

	// Convert API response to Terraform model
	var operatingSystems []OperatingSystemDataModel
	for _, os := range addons.OperatingSystems.Products {
		operatingSystems = append(operatingSystems, OperatingSystemDataModel{
			Name:          types.StringValue(os.Name),
			OSType:        types.StringValue(os.OSType),
			ProductCode:   types.StringValue(os.ProductCode),
			Price:         types.Float64Value(os.Price),
			PriceHourly:   types.Float64Value(os.PriceHourly),
			HourlyEnabled: types.BoolValue(os.HourlyEnabled),
		})
	}

	data.OperatingSystems = operatingSystems
	data.ID = types.StringValue(fmt.Sprintf("%s-%s", serverTypeName, location))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}