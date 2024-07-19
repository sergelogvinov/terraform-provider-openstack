package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
)

// RecordSetCreateOpts represents the attributes used when creating a new DNS record set.
type RecordSetCreateOpts struct {
	recordsets.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToRecordSetCreateMap casts a CreateOpts struct to a map.
// It overrides recordsets.ToRecordSetCreateMap to add the ValueSpecs field.
func (opts RecordSetCreateOpts) ToRecordSetCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

func dnsRecordSetV2RefreshFunc(dnsClient *gophercloud.ServiceClient, zoneID, recordsetID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recordset, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return recordset, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_recordset_v2 %s current status: %s", recordset.ID, recordset.Status)
		return recordset, recordset.Status, nil
	}
}
