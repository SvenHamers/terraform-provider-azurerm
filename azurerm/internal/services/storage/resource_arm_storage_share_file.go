package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceArmStorageShareFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareFileCreate,
		Update: resourceArmStorageShareFileUpdate,
		Delete: resourceArmStorageShareFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareDirectoryName,
			},
			"share_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"subdir": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"filename": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"sourceFileLocal": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"metadata": MetaDataSchema(),
		},
	}
}

func resourceArmStorageShareFileCreate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("share_name").(string)
	directoryName := d.Get("name").(string)
	subdir := d.Get("subdir").(string)
	fileName := d.Get("filename").(string)
	sourceFileLocal := d.Get("sourceFileLocal").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", accountName, directoryName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Directories Client: %s", err)
	}

	file, err := os.Open(sourceFileLocal)
	if err != nil {
		return fmt.Errorf("failed to load file from disk: %s", err)
	}

	if err := client.PutFile(ctx, accountName, shareName, subdir, fileName, file, 1); err != nil {
		return fmt.Errorf("Error creating file %q (File Share %q / Account %q /file &q): %+v", directoryName, shareName, accountName, fileName, err)
	}

	resourceID := client.GetResourceID(accountName, shareName, subdir, fileName)
	d.SetId(resourceID)

	return resourceArmStorageShareDirectoryRead(d, meta)
}

func resourceArmStorageShareFileUpdate(d *schema.ResourceData, meta interface{}) error {
	//ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	//defer cancel()
	//storageClient := meta.(*clients.Client).Storage

	return nil
}

func resourceArmStorageShareFileDelete(d *schema.ResourceData, meta interface{}) error {
	//ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	//defer cancel()
	//storageClient := meta.(*clients.Client).Storage

	return nil
}
