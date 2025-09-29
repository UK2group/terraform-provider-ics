package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SSHKeyResource{}
var _ resource.ResourceWithImportState = &SSHKeyResource{}

func NewSSHKeyResource() resource.Resource {
	return &SSHKeyResource{}
}

// SSHKeyResource defines the resource implementation.
type SSHKeyResource struct {
	client *ICSClient
}

// SSHKeyResourceModel describes the resource data model.
type SSHKeyResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Label     types.String `tfsdk:"label"`
	PublicKey types.String `tfsdk:"public_key"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	UpdatedAt types.Int64  `tfsdk:"updated_at"`
}

func (r *SSHKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SSHKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SSH key resource for server access",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "SSH key identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Label for the SSH key (must be unique)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "SSH public key content",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Creation timestamp",
				Computed:            true,
			},
			"updated_at": schema.Int64Attribute{
				MarkdownDescription: "Last update timestamp",
				Computed:            true,
			},
		},
	}
}

func (r *SSHKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	label := data.Label.ValueString()
	publicKey := data.PublicKey.ValueString()

	// Create SSH key request
	createReq := SSHKeyCreateRequest{
		Label:     label,
		PublicKey: publicKey,
	}

	tflog.Info(ctx, "Creating SSH key", map[string]interface{}{
		"label": label,
	})

	// Create the SSH key
	_, err := r.client.CreateSSHKey(createReq)
	if err != nil {
		resp.Diagnostics.AddError("SSH Key Creation Failed", fmt.Sprintf("Unable to create SSH key: %s", err))
		return
	}

	// Get the full SSH key details to populate computed fields
	sshKey, err := r.client.GetSSHKeyByLabel(label)
	if err != nil {
		resp.Diagnostics.AddError("SSH Key Retrieval Failed", fmt.Sprintf("SSH key created but unable to retrieve details: %s", err))
		return
	}

	// Update the model with SSH key details
	r.updateModelFromSSHKey(&data, sshKey)

	tflog.Info(ctx, "SSH key created successfully", map[string]interface{}{
		"id":    sshKey.ID,
		"label": sshKey.Label,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state from API using the label (since that's what users work with)
	label := data.Label.ValueString()
	sshKey, err := r.client.GetSSHKeyByLabel(label)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSH key with label '%s', got error: %s", label, err))
		return
	}

	// Update the model with current SSH key state
	r.updateModelFromSSHKey(&data, sshKey)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// SSH keys don't support updates - require replacement
	resp.Diagnostics.AddWarning("Update Not Supported", "SSH key changes require replacement. Please destroy and recreate the resource.")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	keyID := int(data.ID.ValueInt64())
	err := r.client.DeleteSSHKey(keyID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SSH key %d, got error: %s", keyID, err))
		return
	}

	tflog.Info(ctx, "SSH key deleted successfully", map[string]interface{}{
		"id": keyID,
	})
}

func (r *SSHKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by label (more user-friendly than ID)
	label := req.ID

	sshKey, err := r.client.GetSSHKeyByLabel(label)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to find SSH key with label '%s': %s", label, err))
		return
	}

	// Set the state
	var data SSHKeyResourceModel
	r.updateModelFromSSHKey(&data, sshKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// updateModelFromSSHKey updates the Terraform model with SSH key data
func (r *SSHKeyResource) updateModelFromSSHKey(data *SSHKeyResourceModel, sshKey *SSHKey) {
	data.ID = types.Int64Value(int64(sshKey.ID))
	data.Label = types.StringValue(sshKey.Label)
	data.PublicKey = types.StringValue(sshKey.Key)
	data.CreatedAt = types.Int64Value(sshKey.CreatedAt)
	data.UpdatedAt = types.Int64Value(sshKey.UpdatedAt)
}