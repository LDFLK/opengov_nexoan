package tests

import (
	"testing"
)

// TODO: Please add more tests cases when we cover other angels about gazette tracking. 

func TestAddDocumentEntity(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		transaction    map[string]interface{}
		entityCounters map[string]int
		wantErr        bool
	}{
		{
			name: "Add organization document",
			transaction: map[string]interface{}{
				"transaction_id": "2403-53",
				"date":           "2024-09-27",
				"url":            "",
				"description":    "Anura Kumara Dissanayake",
				"child_type":     "extgzt:org",
				"child":          "2403-53",
				"parent_type":    "government",
				"parent":         "Government of Sri Lanka",
				"rel_type":       "HAS_DOCUMENT",
			},
			entityCounters: map[string]int{
				"document": 0,
			},
			wantErr: false,
		},
		{
			name: "Add person document",
			transaction: map[string]interface{}{
				"transaction_id": "2403-03",
				"date":           "2024-08-23",
				"url":            "",
				"description":    "Anura Kumara Dissanayake",
				"child_type":     "extgzt:person",
				"child":          "2403-03",
				"parent_type":    "government",
				"parent":         "Government of Sri Lanka",
				"rel_type":       "HAS_DOCUMENT",
			},
			entityCounters: map[string]int{
				"document": 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, err := client.AddDocumentEntity(tt.transaction, tt.entityCounters)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddDocumentEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && counter == 0 {
				t.Error("AddDocumentEntity() counter = 0, want > 0")
			}
		})
	}
}
