package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BareMetalServerResource{}
var _ resource.ResourceWithImportState = &BareMetalServerResource{}

func NewBareMetalServerResource() resource.Resource {
	return &BareMetalServerResource{}
}

// BareMetalServerResource defines the resource implementation.
type BareMetalServerResource struct {
	client *ICSClient
}

// BareMetalServerResourceModel describes the resource data model.
type BareMetalServerResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	InstanceType       types.String `tfsdk:"instance_type"`     // User-friendly instance type (e.g. c1i.small)
	Location           types.String `tfsdk:"location"`          // Location code (e.g. NYC1)
	OperatingSystem    types.String `tfsdk:"operating_system"`  // OS name (e.g. Ubuntu 24.04)
	Hostname           types.String `tfsdk:"hostname"`
	FriendlyName       types.String `tfsdk:"friendly_name"`
	SSHKeyLabels       types.List   `tfsdk:"ssh_key_labels"`

	// Computed/output fields
	ServiceID          types.Int64  `tfsdk:"service_id"`
	PublicIP           types.String `tfsdk:"public_ip"`
	RootPassword       types.String `tfsdk:"root_password"`
	ServiceDescription types.String `tfsdk:"service_description"`
	PlanID             types.Int64  `tfsdk:"plan_id"`
	DatacenterName     types.String `tfsdk:"datacenter_name"`
	DatacenterID       types.Int64  `tfsdk:"datacenter_id"`
	LocationID         types.Int64  `tfsdk:"location_id"`
	ServerTypeInternal types.String `tfsdk:"server_type"` // Keep for internal use
}

func (r *BareMetalServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bare_metal_server"
}

func (r *BareMetalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Bare metal server resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_type": schema.StringAttribute{
				MarkdownDescription: "Instance type (e.g., 'c1.small', 'c1.medium'). The provider will automatically validate availability and inventory.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location code (e.g., 'NYC1', 'FRA1'). The provider will automatically validate inventory availability for this location.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operating_system": schema.StringAttribute{
				MarkdownDescription: "Operating system name (e.g., 'Ubuntu 24.04', 'Debian 12', 'CentOS 8'). The provider will automatically validate availability for the specified instance type and location.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname for the server",
				Optional:            true,
			},
			"friendly_name": schema.StringAttribute{
				MarkdownDescription: "Friendly name for the server",
				Optional:            true,
			},
			"ssh_key_labels": schema.ListAttribute{
				MarkdownDescription: "List of SSH key labels to add to the server. The SSH keys must already exist.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"service_id": schema.Int64Attribute{
				MarkdownDescription: "Service identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"root_password": schema.StringAttribute{
				MarkdownDescription: "Root password for the server",
				Computed:            true,
				Sensitive:           true,
			},
			"public_ip": schema.StringAttribute{
				MarkdownDescription: "Public IP address",
				Computed:            true,
			},
			"service_description": schema.StringAttribute{
				MarkdownDescription: "Service description",
				Computed:            true,
			},
			"plan_id": schema.Int64Attribute{
				MarkdownDescription: "Plan identifier",
				Computed:            true,
			},
			"datacenter_name": schema.StringAttribute{
				MarkdownDescription: "Datacenter name",
				Computed:            true,
			},
			"datacenter_id": schema.Int64Attribute{
				MarkdownDescription: "Datacenter identifier",
				Computed:            true,
			},
			"location_id": schema.Int64Attribute{
				MarkdownDescription: "Location identifier",
				Computed:            true,
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Server type",
				Computed:            true,
			},
		},
	}
}

func (r *BareMetalServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ICSClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ICSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *BareMetalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BareMetalServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	instanceType := data.InstanceType.ValueString()
	location := data.Location.ValueString()
	osName := data.OperatingSystem.ValueString()

	// Step 1: Validate instance type and location combination with inventory
	tflog.Info(ctx, "Validating instance type and location", map[string]interface{}{
		"instance_type": instanceType,
		"location":      location,
	})

	sku, err := r.client.FindSKUByProductName(instanceType, location)
	if err != nil {
		// Provide helpful error messages with suggestions
		inventory, invErr := r.client.GetInventory()
		if invErr != nil {
			resp.Diagnostics.AddError(
				"Instance Type or Location Invalid",
				fmt.Sprintf("Unable to validate instance type '%s' in location '%s': %s", instanceType, location, err),
			)
			return
		}

		// Find available alternatives
		var availableTypes []string
		var availableLocations []string
		typeLocationMap := make(map[string][]string)

		for _, item := range inventory {
			if item.AutoProvisionQuantity > 0 {
				if item.SkuProductName == instanceType {
					availableLocations = append(availableLocations, item.LocationCode)
				}
				if item.LocationCode == location {
					availableTypes = append(availableTypes, item.SkuProductName)
				}
				if typeLocationMap[item.SkuProductName] == nil {
					typeLocationMap[item.SkuProductName] = []string{}
				}
				typeLocationMap[item.SkuProductName] = append(typeLocationMap[item.SkuProductName], item.LocationCode)
			}
		}

		errorMsg := fmt.Sprintf("Instance type '%s' is not available in location '%s'", instanceType, location)

		if len(availableTypes) > 0 {
			errorMsg += fmt.Sprintf("\n\nAvailable instance types in location '%s': %v", location, availableTypes)
		}

		if len(availableLocations) > 0 {
			errorMsg += fmt.Sprintf("\n\nAvailable locations for instance type '%s': %v", instanceType, availableLocations)
		}

		if len(typeLocationMap) > 0 {
			errorMsg += "\n\nAll available combinations with inventory:"
			for iType, locs := range typeLocationMap {
				errorMsg += fmt.Sprintf("\n  %s: %v", iType, locs)
			}
		}

		resp.Diagnostics.AddError("Invalid Instance Type and Location Combination", errorMsg)
		return
	}

	// Step 2: Validate operating system
	tflog.Info(ctx, "Validating operating system", map[string]interface{}{
		"os_name":       osName,
		"instance_type": instanceType,
		"location":      location,
	})

	addons, osErr := r.client.GetAddons(instanceType, location)
	if osErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Retrieve Operating System Options",
			fmt.Sprintf("Unable to get available operating systems for instance type '%s' in location '%s': %s", instanceType, location, osErr),
		)
		return
	}

	var os *OperatingSystemItem
	var availableOSNames []string

	for _, osOption := range addons.OperatingSystems.Products {
		availableOSNames = append(availableOSNames, osOption.Name)
		if osOption.Name == osName {
			os = &osOption
			break
		}
	}

	if os == nil {
		errorMsg := fmt.Sprintf("Operating system '%s' is not available for instance type '%s' in location '%s'", osName, instanceType, location)
		errorMsg += fmt.Sprintf("\n\nAvailable operating systems: %v", availableOSNames)

		resp.Diagnostics.AddError("Invalid Operating System", errorMsg)
		return
	}

	// All validations passed - proceed with server ordering
	tflog.Info(ctx, "All validations passed, proceeding with server order", map[string]interface{}{
		"instance_type":  instanceType,
		"location":       location,
		"os":             osName,
		"sku_id":         sku.SkuID,
		"os_product_code": os.ProductCode,
	})

	// Create server order request with all required fields
	orderReq := ServerOrderRequest{
		SkuProductName:             instanceType,     // Required: e.g. "c2.small"
		Quantity:                   1,               // Required: hardcoded to 1
		LocationCode:               location,        // Required: e.g. "FRA1"
		OperatingSystemProductCode: os.ProductCode,  // Required: e.g. "UBUNTU_24_04"
		BillHourly:                 true,           // Always bill hourly
	}

	if !data.Hostname.IsNull() {
		orderReq.Hostname = data.Hostname.ValueString()
	}

	// Handle SSH key labels - convert to SSH key IDs
	if !data.SSHKeyLabels.IsNull() {
		var sshKeyLabels []string
		resp.Diagnostics.Append(data.SSHKeyLabels.ElementsAs(ctx, &sshKeyLabels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var sshKeyIDs []int
		for _, label := range sshKeyLabels {
			sshKey, err := r.client.GetSSHKeyByLabel(label)
			if err != nil {
				resp.Diagnostics.AddError(
					"SSH Key Not Found",
					fmt.Sprintf("SSH key with label '%s' not found. Please ensure the SSH key exists before ordering the server. Error: %s", label, err),
				)
				return
			}
			sshKeyIDs = append(sshKeyIDs, sshKey.ID)
		}

		if len(sshKeyIDs) > 0 {
			orderReq.SSHKeyIDs = sshKeyIDs
			tflog.Info(ctx, "Adding SSH keys to server order", map[string]interface{}{
				"ssh_key_labels": sshKeyLabels,
				"ssh_key_ids":    sshKeyIDs,
			})
		}
	}

	// Order the server
	orderResp, err := r.client.OrderServer(orderReq)

	if err != nil {
		// Check if this is a timeout error after a potentially successful order
		if strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Client.Timeout exceeded") {
			resp.Diagnostics.AddError(
				"Server Order Timeout",
				fmt.Sprintf("The server order request timed out, but the order may have been successful. Please check the ICS control panel for any pending orders, or try running 'terraform refresh' to check if a server was created. Error: %s", err),
			)
		} else {
			resp.Diagnostics.AddError("Server Order Failed", fmt.Sprintf("Unable to order server: %s", err))
		}
		return
	}

	if orderResp == nil {
		resp.Diagnostics.AddError("Server Order Failed", "Order response is nil - this indicates an API parsing issue")
		return
	}

	if len(orderResp.OrderServiceIDs) == 0 {
		resp.Diagnostics.AddError("Server Order Failed", fmt.Sprintf("No service IDs returned from server order. Response: %+v", orderResp))
		return
	}

	serviceID := orderResp.OrderServiceIDs[0]
	tflog.Info(ctx, "Server ordered successfully, waiting for provisioning", map[string]interface{}{
		"service_id": serviceID,
	})

	// Wait for server to be provisioned (up to 30 minutes)
	server, err := r.waitForServerProvisioning(ctx, serviceID, 30*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Server Provisioning Timeout",
			fmt.Sprintf("Server was ordered (service ID: %d) but provisioning did not complete within 30 minutes: %s\n\nYou can check the provisioning status in the ICS control panel or try running 'terraform refresh' later.", serviceID, err),
		)
		return
	}

	// Update friendly name if specified (post-provision)
	if !data.FriendlyName.IsNull() {
		friendlyName := data.FriendlyName.ValueString()
		tflog.Info(ctx, "Setting server friendly name", map[string]interface{}{
			"server_id":     server.ID,
			"friendly_name": friendlyName,
		})

		err = r.client.UpdateServerFriendlyName(server.ID, friendlyName)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Friendly Name Update Failed",
				fmt.Sprintf("Server was provisioned successfully but failed to set friendly name: %s", err),
			)
		}
	}

	// Update the model with server details
	r.updateModelFromServer(&data, server)

	tflog.Info(ctx, "Server provisioned successfully", map[string]interface{}{
		"server_id":  server.ID,
		"service_id": server.ServiceID,
		"public_ip":  server.PublicIP,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BareMetalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BareMetalServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state from API
	serviceID := int(data.ServiceID.ValueInt64())
	server, err := r.client.GetServerByServiceID(serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server with service ID %d, got error: %s", serviceID, err))
		return
	}

	// Update the model with current server state
	r.updateModelFromServer(&data, server)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BareMetalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BareMetalServerResourceModel
	var state BareMetalServerResourceModel

	// Read Terraform plan and current state data into the models
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if friendly name changed
	if !plan.FriendlyName.Equal(state.FriendlyName) && !plan.FriendlyName.IsNull() {
		serverID := state.ID.ValueString()
		friendlyName := plan.FriendlyName.ValueString()

		tflog.Info(ctx, "Updating server friendly name", map[string]interface{}{
			"server_id":     serverID,
			"friendly_name": friendlyName,
		})

		err := r.client.UpdateServerFriendlyName(serverID, friendlyName)
		if err != nil {
			resp.Diagnostics.AddError("Friendly Name Update Failed", fmt.Sprintf("Unable to update server friendly name: %s", err))
			return
		}
	}

	// For other changes, require replacement
	if !plan.InstanceType.Equal(state.InstanceType) ||
		!plan.Location.Equal(state.Location) ||
		!plan.OperatingSystem.Equal(state.OperatingSystem) ||
		!plan.Hostname.Equal(state.Hostname) ||
		!plan.SSHKeyLabels.Equal(state.SSHKeyLabels) {
		resp.Diagnostics.AddWarning("Update Requires Replacement", "Changes to instance_type, location, operating_system, hostname, or ssh_key_labels require resource replacement. Please destroy and recreate the resource.")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BareMetalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BareMetalServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Since we always use hourly billing, we can cancel via API
	serverID := data.ID.ValueString()
	err := r.client.CancelServer(serverID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to cancel server %s, got error: %s", serverID, err))
		return
	}

	tflog.Info(ctx, "Server canceled successfully", map[string]interface{}{
		"server_id": serverID,
	})
}

func (r *BareMetalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by service ID
	serviceIDStr := req.ID
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Invalid service ID format: %s", serviceIDStr))
		return
	}

	server, err := r.client.GetServerByServiceID(serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to find server with service ID %d: %s", serviceID, err))
		return
	}

	// Set the state
	var data BareMetalServerResourceModel
	r.updateModelFromServer(&data, server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// waitForServerProvisioning waits for a server to be provisioned
func (r *BareMetalServerResource) waitForServerProvisioning(ctx context.Context, serviceID int, timeout time.Duration) (*Server, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		tflog.Debug(ctx, "Checking server provisioning status", map[string]interface{}{
			"service_id": serviceID,
		})

		server, err := r.client.GetServerByServiceID(serviceID)
		if err == nil {
			// Server found, provisioning complete
			return server, nil
		}

		// If server not found, wait and retry
		time.Sleep(30 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for server with service ID %d to be provisioned", serviceID)
}

// updateModelFromServer updates the Terraform model with server data
func (r *BareMetalServerResource) updateModelFromServer(data *BareMetalServerResourceModel, server *Server) {
	data.ID = types.StringValue(server.ID)
	data.ServiceID = types.Int64Value(int64(server.ServiceID))
	data.PublicIP = types.StringValue(server.PublicIP)
	data.RootPassword = types.StringValue(server.RootPassword)
	data.ServiceDescription = types.StringValue(server.ServiceDescription)
	data.PlanID = types.Int64Value(int64(server.PlanID))
	data.DatacenterName = types.StringValue(server.DatacenterName)
	data.DatacenterID = types.Int64Value(int64(server.DatacenterID))
	data.LocationID = types.Int64Value(int64(server.LocationID))
	data.ServerTypeInternal = types.StringValue(server.ServerType)

	// Update input fields if they were computed
	if data.Hostname.IsNull() && server.Hostname != "" {
		data.Hostname = types.StringValue(server.Hostname)
	}
	if data.FriendlyName.IsNull() && server.FriendlyName != "" {
		data.FriendlyName = types.StringValue(server.FriendlyName)
	}
}