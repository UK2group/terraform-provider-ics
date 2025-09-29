package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InventoryDataSource{}

func NewInventoryDataSource() datasource.DataSource {
	return &InventoryDataSource{}
}

// InventoryDataSource defines the data source implementation.
type InventoryDataSource struct {
	client *ICSClient
}

// InventoryDataSourceModel describes the data source data model.
type InventoryDataSourceModel struct {
	Items []InventoryItemModel `tfsdk:"items"`
	ID    types.String         `tfsdk:"id"`
}

type InventoryItemModel struct {
	SkuID                 types.Int64                   `tfsdk:"sku_id"`
	Quantity              types.Int64                   `tfsdk:"quantity"`
	AutoProvisionQuantity types.Int64                   `tfsdk:"auto_provision_quantity"`
	DatacenterID          types.Int64                   `tfsdk:"datacenter_id"`
	RegionID              types.Int64                   `tfsdk:"region_id"`
	LocationCode          types.String                  `tfsdk:"location_code"`
	CPUBrand              types.String                  `tfsdk:"cpu_brand"`
	CPUModel              types.String                  `tfsdk:"cpu_model"`
	CPUClockSpeedGHz      types.Float64                 `tfsdk:"cpu_clock_speed_ghz"`
	CPUCores              types.Int64                   `tfsdk:"cpu_cores"`
	CPUCount              types.Int64                   `tfsdk:"cpu_count"`
	TotalSSDSizeGB        types.Int64                   `tfsdk:"total_ssd_size_gb"`
	TotalHDDSizeGB        types.Int64                   `tfsdk:"total_hdd_size_gb"`
	TotalNVMESizeGB       types.Int64                   `tfsdk:"total_nvme_size_gb"`
	RAIDEnabled           types.Bool                    `tfsdk:"raid_enabled"`
	TotalRAMGB            types.Int64                   `tfsdk:"total_ram_gb"`
	NICSpeedMbps          types.Int64                   `tfsdk:"nic_speed_mbps"`
	QTProductID           types.Int64                   `tfsdk:"qt_product_id"`
	Status                types.String                  `tfsdk:"status"`
	Metadata              []InventoryMetadataModel      `tfsdk:"metadata"`
	CurrencyCode          types.String                  `tfsdk:"currency_code"`
	SkuProductName        types.String                  `tfsdk:"sku_product_name"`
	Price                 types.String                  `tfsdk:"price"`
	PriceHourly           types.String                  `tfsdk:"price_hourly"`
	HourlyEnabled         types.Bool                    `tfsdk:"hourly_enabled"`
}

type InventoryMetadataModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

func (d *InventoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (d *InventoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Inventory data source provides information about available bare metal servers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier",
				Computed:            true,
			},
			"items": schema.ListNestedAttribute{
				MarkdownDescription: "List of available server SKUs",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sku_id": schema.Int64Attribute{
							MarkdownDescription: "SKU identifier",
							Computed:            true,
						},
						"quantity": schema.Int64Attribute{
							MarkdownDescription: "Available quantity",
							Computed:            true,
						},
						"auto_provision_quantity": schema.Int64Attribute{
							MarkdownDescription: "Auto provision quantity",
							Computed:            true,
						},
						"datacenter_id": schema.Int64Attribute{
							MarkdownDescription: "Datacenter identifier",
							Computed:            true,
						},
						"region_id": schema.Int64Attribute{
							MarkdownDescription: "Region identifier",
							Computed:            true,
						},
						"location_code": schema.StringAttribute{
							MarkdownDescription: "Location code",
							Computed:            true,
						},
						"cpu_brand": schema.StringAttribute{
							MarkdownDescription: "CPU brand",
							Computed:            true,
						},
						"cpu_model": schema.StringAttribute{
							MarkdownDescription: "CPU model",
							Computed:            true,
						},
						"cpu_clock_speed_ghz": schema.Float64Attribute{
							MarkdownDescription: "CPU clock speed in GHz",
							Computed:            true,
						},
						"cpu_cores": schema.Int64Attribute{
							MarkdownDescription: "Number of CPU cores",
							Computed:            true,
						},
						"cpu_count": schema.Int64Attribute{
							MarkdownDescription: "Number of CPUs",
							Computed:            true,
						},
						"total_ssd_size_gb": schema.Int64Attribute{
							MarkdownDescription: "Total SSD size in GB",
							Computed:            true,
						},
						"total_hdd_size_gb": schema.Int64Attribute{
							MarkdownDescription: "Total HDD size in GB",
							Computed:            true,
						},
						"total_nvme_size_gb": schema.Int64Attribute{
							MarkdownDescription: "Total NVMe size in GB",
							Computed:            true,
						},
						"raid_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether RAID is enabled",
							Computed:            true,
						},
						"total_ram_gb": schema.Int64Attribute{
							MarkdownDescription: "Total RAM in GB",
							Computed:            true,
						},
						"nic_speed_mbps": schema.Int64Attribute{
							MarkdownDescription: "NIC speed in Mbps",
							Computed:            true,
						},
						"qt_product_id": schema.Int64Attribute{
							MarkdownDescription: "QT product identifier",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Status",
							Computed:            true,
						},
						"metadata": schema.ListNestedAttribute{
							MarkdownDescription: "Metadata",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Metadata name",
										Computed:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "Metadata description",
										Computed:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "Metadata value",
										Computed:            true,
									},
								},
							},
						},
						"currency_code": schema.StringAttribute{
							MarkdownDescription: "Currency code",
							Computed:            true,
						},
						"sku_product_name": schema.StringAttribute{
							MarkdownDescription: "SKU product name",
							Computed:            true,
						},
						"price": schema.StringAttribute{
							MarkdownDescription: "Price",
							Computed:            true,
						},
						"price_hourly": schema.StringAttribute{
							MarkdownDescription: "Hourly price",
							Computed:            true,
						},
						"hourly_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether hourly billing is enabled",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *InventoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InventoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InventoryDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get inventory from API
	inventory, err := d.client.GetInventory()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory, got error: %s", err))
		return
	}

	// Convert API response to Terraform model
	var items []InventoryItemModel
	for _, item := range inventory {
		var metadata []InventoryMetadataModel
		for _, meta := range item.Metadata {
			metadata = append(metadata, InventoryMetadataModel{
				Name:        types.StringValue(meta.Name),
				Description: types.StringValue(meta.Description),
				Value:       types.StringValue(meta.Value),
			})
		}

		items = append(items, InventoryItemModel{
			SkuID:                 types.Int64Value(int64(item.SkuID)),
			Quantity:              types.Int64Value(int64(item.Quantity)),
			AutoProvisionQuantity: types.Int64Value(int64(item.AutoProvisionQuantity)),
			DatacenterID:          types.Int64Value(int64(item.DatacenterID)),
			RegionID:              types.Int64Value(int64(item.RegionID)),
			LocationCode:          types.StringValue(item.LocationCode),
			CPUBrand:              types.StringValue(item.CPUBrand),
			CPUModel:              types.StringValue(item.CPUModel),
			CPUClockSpeedGHz:      types.Float64Value(item.CPUClockSpeedGHz),
			CPUCores:              types.Int64Value(int64(item.CPUCores)),
			CPUCount:              types.Int64Value(int64(item.CPUCount)),
			TotalSSDSizeGB:        types.Int64Value(int64(item.TotalSSDSizeGB)),
			TotalHDDSizeGB:        types.Int64Value(int64(item.TotalHDDSizeGB)),
			TotalNVMESizeGB:       types.Int64Value(int64(item.TotalNVMESizeGB)),
			RAIDEnabled:           types.BoolValue(item.RAIDEnabled),
			TotalRAMGB:            types.Int64Value(int64(item.TotalRAMGB)),
			NICSpeedMbps:          types.Int64Value(int64(item.NICSpeedMbps)),
			QTProductID:           types.Int64Value(int64(item.QTProductID)),
			Status:                types.StringValue(item.Status),
			Metadata:              metadata,
			CurrencyCode:          types.StringValue(item.CurrencyCode),
			SkuProductName:        types.StringValue(item.SkuProductName),
			Price:                 types.StringValue(item.Price),
			PriceHourly:           types.StringValue(item.PriceHourly),
			HourlyEnabled:         types.BoolValue(item.HourlyEnabled),
		})
	}

	data.Items = items
	data.ID = types.StringValue("inventory")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}