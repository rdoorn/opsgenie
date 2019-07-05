package main

import (
	"fmt"
)

func (h handler) filterContains(query, timeStr string) error {
	policyID, policyName, err := h.createPolicy("contains", query)
	if err != nil {
		return fmt.Errorf("create policy failed: %s", err)
	}
	return h.maintenanceCreate(policyID, policyName, timeStr)
}

func (h handler) filterRegex(query, timeStr string) {
}
