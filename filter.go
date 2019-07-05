package main

import (
	"fmt"
)

func (h handler) filterContains(query, timeStr string) error {

	if timeStr == "" {
		return fmt.Errorf("you must specify for how long this filter should be enabled (e.g. 1h30m or 1d)")
	}
	policyID, policyName, err := h.createPolicy("contains", query)
	if err != nil {
		return fmt.Errorf("create policy failed: %s", err)
	}
	return h.maintenanceCreate(policyID, policyName, timeStr)
}

func (h handler) filterRegex(query, timeStr string) error {
	if timeStr == "" {
		return fmt.Errorf("you must specify for how long this filter should be enabled (e.g. 1h30m or 1d)")
	}
	policyID, policyName, err := h.createPolicy("regex", query)
	if err != nil {
		return fmt.Errorf("create policy failed: %s", err)
	}
	return h.maintenanceCreate(policyID, policyName, timeStr)
}
