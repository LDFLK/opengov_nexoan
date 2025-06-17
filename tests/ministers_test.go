package tests

import (
	"fmt"
	"orgchart_nexoan/api"
	"orgchart_nexoan/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *api.Client

func TestMain(m *testing.M) {
	// Set up test environment with correct URLs
	client = api.NewClient("http://localhost:8080/entities", "http://localhost:8081/v1/entities")

	// Create government node using CreateGovernmentNode
	government, err := client.CreateGovernmentNode()
	if err != nil {
		fmt.Printf("Failed to create government node: %v\n", err)
		os.Exit(1)
	}
	if government == nil {
		fmt.Println("Government node is nil")
		os.Exit(1)
	}
	fmt.Printf("Successfully created government node with ID: %s\n", government.ID)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateMinisters(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 0,
	}

	// Test cases for creating ministers
	testCases := []struct {
		transactionID string
		parent        string
		parentType    string
		child         string
		childType     string
		relType       string
		date          string
	}{
		{
			transactionID: "2153/12_tr_01",
			parent:        "Government of Sri Lanka",
			parentType:    "government",
			child:         "Minister of Defence",
			childType:     "minister",
			relType:       "AS_MINISTER",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_02",
			parent:        "Government of Sri Lanka",
			parentType:    "government",
			child:         "Minister of Finance, Economic and Policy Development",
			childType:     "minister",
			relType:       "AS_MINISTER",
			date:          "2019-12-10",
		},
	}

	// Create each minister
	for _, tc := range testCases {
		t.Logf("Creating minister: %s", tc.child)

		// Create transaction map for AddEntity
		transaction := map[string]interface{}{
			"parent":         tc.parent,
			"child":          tc.child,
			"date":           tc.date,
			"parent_type":    tc.parentType,
			"child_type":     tc.childType,
			"rel_type":       tc.relType,
			"transaction_id": tc.transactionID,
		}

		// Use AddEntity to create the minister
		_, err := client.AddOrgEntity(transaction, entityCounters)
		assert.NoError(t, err)

		// Update the counter for the next iteration
		entityCounters[tc.childType]++

		// Verify the minister was created by searching for it
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.childType,
			},
			Name: tc.child,
		}

		results, err := client.SearchEntities(searchCriteria)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, tc.child, results[0].Name)

		// Verify the relationship was created by checking parent's relationships
		parentResults, err := client.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.parentType,
			},
			Name: tc.parent,
		})
		assert.NoError(t, err)
		assert.Len(t, parentResults, 1)

		// Get parent's metadata to verify relationship
		// metadata, err := client.GetEntityMetadata(parentResults[0].ID)
		// assert.NoError(t, err)
		// assert.NotNil(t, metadata)
	}
}

func TestCreateDepartments(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"department": 0,
	}

	// Test cases for creating departments
	testCases := []struct {
		transactionID string
		parent        string
		parentType    string
		child         string
		childType     string
		relType       string
		date          string
	}{
		{
			transactionID: "2153/12_tr_03",
			parent:        "Minister of Defence",
			parentType:    "minister",
			child:         "Sri Lankan Army",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_04",
			parent:        "Minister of Finance, Economic and Policy Development",
			parentType:    "minister",
			child:         "Department of Taxes",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
		{
			transactionID: "2153/12_tr_05",
			parent:        "Minister of Finance, Economic and Policy Development",
			parentType:    "minister",
			child:         "Department of Policies",
			childType:     "department",
			relType:       "AS_DEPARTMENT",
			date:          "2019-12-10",
		},
	}

	// Create each department
	for _, tc := range testCases {
		t.Logf("Creating department: %s under minister: %s", tc.child, tc.parent)

		// Create transaction map for AddEntity
		transaction := map[string]interface{}{
			"parent":         tc.parent,
			"child":          tc.child,
			"date":           tc.date,
			"parent_type":    tc.parentType,
			"child_type":     tc.childType,
			"rel_type":       tc.relType,
			"transaction_id": tc.transactionID,
		}

		// Use AddEntity to create the department
		_, err := client.AddOrgEntity(transaction, entityCounters)
		assert.NoError(t, err)

		// Update the counter for the next iteration
		entityCounters[tc.childType]++

		// Verify the department was created by searching for it
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.childType,
			},
			Name: tc.child,
		}

		results, err := client.SearchEntities(searchCriteria)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, tc.child, results[0].Name)

		// Verify the relationship was created by checking minister's relationships
		ministerResults, err := client.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: tc.parentType,
			},
			Name: tc.parent,
		})
		assert.NoError(t, err)
		assert.Len(t, ministerResults, 1)

		// Get minister's relationships to verify department relationship
		// relations, err := client.GetRelatedEntities(ministerResults[0].ID, &models.Relationship{
		// 	Name:    tc.relType,
		// 	EndTime: "",
		// })
		// assert.NoError(t, err)
		// assert.NotEmpty(t, relations, "Minister should have at least one department relationship")
	}
}

func TestTerminateDepartment(t *testing.T) {
	// Create transaction map for terminating the department
	transaction := map[string]interface{}{
		"parent":      "Minister of Defence",
		"child":       "Sri Lankan Army",
		"date":        "2024-01-01",
		"parent_type": "minister",
		"child_type":  "department",
		"rel_type":    "AS_DEPARTMENT",
	}

	// Terminate the department relationship
	err := client.TerminateOrgEntity(transaction)
	assert.NoError(t, err)

	// Find the minister to verify the relationship
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Defence",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	// Find the department
	departmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Sri Lankan Army",
	})
	assert.NoError(t, err)
	assert.Len(t, departmentResults, 1)
	departmentID := departmentResults[0].ID

	// Verify the relationship is terminated
	allRelations, err := client.GetAllRelatedEntities(ministerID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated relationship")
}

func TestTerminateMinister(t *testing.T) {
	// Create transaction map for terminating the minister
	transaction := map[string]interface{}{
		"parent":      "Government of Sri Lanka",
		"child":       "Minister of Defence",
		"date":        "2024-01-01",
		"parent_type": "government",
		"child_type":  "minister",
		"rel_type":    "AS_MINISTER",
	}

	// Terminate the minister relationship
	err := client.TerminateOrgEntity(transaction)
	assert.NoError(t, err)

	// Find the government to verify the relationship
	governmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: "Government of Sri Lanka",
	})
	assert.NoError(t, err)
	assert.Len(t, governmentResults, 1)
	governmentID := governmentResults[0].ID

	// Find the minister
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Defence",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	// Verify the relationship is terminated
	allRelations, err := client.GetAllRelatedEntities(governmentID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == ministerID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated relationship")
}

func TestMoveDepartment(t *testing.T) {
	// First create a new minister
	entityCounters := map[string]int{
		"minister": 2, // Since we already have 2 ministers from previous tests
	}

	// Create transaction map for new minister
	newMinisterTransaction := map[string]interface{}{
		"parent":         "Government of Sri Lanka",
		"child":          "Minister of Education",
		"date":           "2024-01-01",
		"parent_type":    "government",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": "2153/12_tr_06",
	}

	// Create the new minister
	_, err := client.AddOrgEntity(newMinisterTransaction, entityCounters)
	assert.NoError(t, err)

	// Create transaction map for moving the department
	transaction := map[string]interface{}{
		"old_parent": "Minister of Finance, Economic and Policy Development",
		"new_parent": "Minister of Education",
		"child":      "Department of Policies",
		"type":       "department",
		"date":       "2024-01-01",
	}

	// Move the department
	err = client.MoveDepartment(transaction)
	assert.NoError(t, err)

	// Find the new minister to verify the new relationship
	newMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Education",
	})
	assert.NoError(t, err)
	assert.Len(t, newMinisterResults, 1)
	newMinisterID := newMinisterResults[0].ID

	// Find the department
	departmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Department of Policies",
	})
	assert.NoError(t, err)
	assert.Len(t, departmentResults, 1)
	departmentID := departmentResults[0].ID

	// Verify the new relationship exists
	allRelations, err := client.GetAllRelatedEntities(newMinisterID)
	assert.NoError(t, err)
	found := false
	for _, rel := range allRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the new relationship")

	// Find the old minister to verify the old relationship is terminated
	oldMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance, Economic and Policy Development",
	})
	assert.NoError(t, err)
	assert.Len(t, oldMinisterResults, 1)
	oldMinisterID := oldMinisterResults[0].ID

	// Verify the old relationship is terminated
	oldRelations, err := client.GetAllRelatedEntities(oldMinisterID)
	assert.NoError(t, err)
	found = false
	for _, rel := range oldRelations {
		if rel.RelatedEntityID == departmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated old relationship")
}

func TestRenameMinister(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 2,
	}

	// Create transaction map for renaming the minister
	transaction := map[string]interface{}{
		"old":            "Minister of Finance, Economic and Policy Development",
		"new":            "Minister of Finance",
		"type":           "minister",
		"date":           "2024-01-01",
		"transaction_id": "2153/13_tr_01",
	}

	// Rename the minister
	newMinisterCounter, err := client.RenameMinister(transaction, entityCounters)
	assert.NoError(t, err)
	assert.Greater(t, newMinisterCounter, 0)

	// Find the new minister to verify it exists
	newMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance",
	})
	assert.NoError(t, err)
	assert.Len(t, newMinisterResults, 1)
	newMinisterID := newMinisterResults[0].ID

	// Find the old minister
	oldMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance, Economic and Policy Development",
	})
	assert.NoError(t, err)
	assert.Len(t, oldMinisterResults, 1)
	oldMinisterID := oldMinisterResults[0].ID

	// Verify the RENAMED_TO relationship exists
	oldRelations, err := client.GetAllRelatedEntities(oldMinisterID)
	assert.NoError(t, err)
	found := false
	for _, rel := range oldRelations {
		if rel.RelatedEntityID == newMinisterID && rel.Name == "RENAMED_TO" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the RENAMED_TO relationship")

	// Verify the old minister's government relationship is terminated
	governmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: "Government of Sri Lanka",
	})
	assert.NoError(t, err)
	assert.Len(t, governmentResults, 1)
	governmentID := governmentResults[0].ID

	govRelations, err := client.GetAllRelatedEntities(governmentID)
	assert.NoError(t, err)
	found = false
	for _, rel := range govRelations {
		if rel.RelatedEntityID == oldMinisterID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated government relationship")

	// Verify the new minister has the government relationship
	found = false
	for _, rel := range govRelations {
		if rel.RelatedEntityID == newMinisterID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the new government relationship")

	// Verify all departments were transferred
	oldDeptRelations, err := client.GetAllRelatedEntities(oldMinisterID)
	assert.NoError(t, err)
	newDeptRelations, err := client.GetAllRelatedEntities(newMinisterID)
	assert.NoError(t, err)

	// Count active department relationships
	oldActiveDepts := 0
	newActiveDepts := 0
	for _, rel := range oldDeptRelations {
		if rel.Name == "AS_DEPARTMENT" && rel.EndTime == "" {
			oldActiveDepts++
		}
	}
	for _, rel := range newDeptRelations {
		if rel.Name == "AS_DEPARTMENT" && rel.EndTime == "" {
			newActiveDepts++
		}
	}

	// All departments should be transferred (no active departments for old minister)
	assert.Equal(t, 0, oldActiveDepts, "Old minister should have no active departments")
	assert.Greater(t, newActiveDepts, 0, "New minister should have active departments")
}

func TestRenameDepartment(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister":   0,
		"department": 0,
	}

	var err error

	// Create government node
	// government, err := client.CreateGovernmentNode()
	// assert.NoError(t, err)
	// assert.NotNil(t, government)

	// Create minister
	ministerTransaction := map[string]interface{}{
		"parent":         "Government of Sri Lanka",
		"child":          "Minister of Finance",
		"date":           "2023-12-01",
		"parent_type":    "government",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": "2153/13_tr_01",
	}

	_, err = client.AddOrgEntity(ministerTransaction, entityCounters)
	assert.NoError(t, err)
	entityCounters["minister"]++

	// Create department
	departmentTransaction := map[string]interface{}{
		"parent":         "Minister of Finance",
		"child":          "Department of Policies",
		"date":           "2023-12-01",
		"parent_type":    "minister",
		"child_type":     "department",
		"rel_type":       "AS_DEPARTMENT",
		"transaction_id": "2153/13_tr_02",
	}

	_, err = client.AddOrgEntity(departmentTransaction, entityCounters)
	assert.NoError(t, err)
	entityCounters["department"]++

	// Create transaction map for renaming the department
	transaction := map[string]interface{}{
		"old":            "Department of Policies",
		"new":            "Department of Policy Development",
		"type":           "department",
		"date":           "2024-01-01",
		"transaction_id": "2153/13_tr_05",
		"president":      ":RW",
	}

	// Rename the department
	newDepartmentCounter, err := client.RenameDepartment(transaction, entityCounters)
	assert.NoError(t, err)
	assert.Greater(t, newDepartmentCounter, 0)

	// Find the new department to verify it exists
	newDepartmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Department of Policy Development",
	})
	assert.NoError(t, err)
	assert.Len(t, newDepartmentResults, 1)
	newDepartmentID := newDepartmentResults[0].ID

	// Find the old department
	oldDepartmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: "Department of Policies",
	})
	assert.NoError(t, err)
	assert.Len(t, oldDepartmentResults, 1)
	oldDepartmentID := oldDepartmentResults[0].ID

	// Verify the RENAMED_TO relationship exists
	oldRelations, err := client.GetAllRelatedEntities(oldDepartmentID)
	assert.NoError(t, err)
	found := false
	for _, rel := range oldRelations {
		if rel.RelatedEntityID == newDepartmentID && rel.Name == "RENAMED_TO" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the RENAMED_TO relationship")

	// Find the minister that owns the department
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	// Verify the old department's minister relationship is terminated
	minRelations, err := client.GetAllRelatedEntities(ministerID)
	assert.NoError(t, err)
	found = false
	for _, rel := range minRelations {
		if rel.RelatedEntityID == oldDepartmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated minister relationship")

	// Verify the new department has the minister relationship
	found = false
	for _, rel := range minRelations {
		if rel.RelatedEntityID == newDepartmentID && rel.Name == "AS_DEPARTMENT" {
			assert.Equal(t, "2024-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the new minister relationship")
}

func TestMergeMinisters(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 0, // Since we already have 3 ministers from previous tests
	}

	// Create transaction map for merging ministers
	transaction := map[string]interface{}{
		"old":            "[Minister of Finance, Minister of Education]",
		"new":            "Minister of Finance and Education",
		"date":           "2025-01-01",
		"transaction_id": "2154/13_tr_01",
	}

	// Merge the ministers
	newMinisterCounter, err := client.MergeMinisters(transaction, entityCounters)
	assert.NoError(t, err)
	assert.Greater(t, newMinisterCounter, 0)

	// Find the new minister to verify it exists
	newMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance and Education",
	})
	assert.NoError(t, err)
	assert.Len(t, newMinisterResults, 1)
	newMinisterID := newMinisterResults[0].ID

	// Find the old ministers
	oldMinisterResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Finance",
	})
	assert.NoError(t, err)
	assert.Len(t, oldMinisterResults, 1)
	oldMinister1ID := oldMinisterResults[0].ID

	oldMinisterResults, err = client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister of Education",
	})
	assert.NoError(t, err)
	assert.Len(t, oldMinisterResults, 1)
	oldMinister2ID := oldMinisterResults[0].ID

	// Verify the MERGED_INTO relationships exist
	oldRelations1, err := client.GetAllRelatedEntities(oldMinister1ID)
	assert.NoError(t, err)
	found := false
	for _, rel := range oldRelations1 {
		if rel.RelatedEntityID == newMinisterID && rel.Name == "MERGED_INTO" {
			assert.Equal(t, "2025-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the MERGED_INTO relationship for first minister")

	oldRelations2, err := client.GetAllRelatedEntities(oldMinister2ID)
	assert.NoError(t, err)
	found = false
	for _, rel := range oldRelations2 {
		if rel.RelatedEntityID == newMinisterID && rel.Name == "MERGED_INTO" {
			assert.Equal(t, "2025-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the MERGED_INTO relationship for second minister")

	// Verify the old ministers' government relationships are terminated
	governmentResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: "Government of Sri Lanka",
	})
	assert.NoError(t, err)
	assert.Len(t, governmentResults, 1)
	governmentID := governmentResults[0].ID

	govRelations, err := client.GetAllRelatedEntities(governmentID)
	assert.NoError(t, err)
	found = false
	for _, rel := range govRelations {
		if rel.RelatedEntityID == oldMinister1ID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2025-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated government relationship for first minister")

	found = false
	for _, rel := range govRelations {
		if rel.RelatedEntityID == oldMinister2ID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2025-01-01T00:00:00Z", rel.EndTime)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the terminated government relationship for second minister")

	// Verify the new minister has the government relationship
	found = false
	for _, rel := range govRelations {
		if rel.RelatedEntityID == newMinisterID && rel.Name == "AS_MINISTER" {
			assert.Equal(t, "2025-01-01T00:00:00Z", rel.StartTime)
			assert.Equal(t, "", rel.EndTime) // Should be active (no end time)
			found = true
			break
		}
	}
	assert.True(t, found, "Should find the new government relationship")

	// Verify all departments were transferred
	newDeptRelations, err := client.GetAllRelatedEntities(newMinisterID)
	assert.NoError(t, err)

	// Count active department relationships
	newActiveDepts := 0
	for _, rel := range newDeptRelations {
		if rel.Name == "AS_DEPARTMENT" && rel.EndTime == "" {
			newActiveDepts++
		}
	}

	// Should have at least 2 departments (one from each old minister)
	assert.GreaterOrEqual(t, newActiveDepts, 2, "New minister should have at least 2 active departments")

	// Verify old ministers have no active departments
	oldDeptRelations1, err := client.GetAllRelatedEntities(oldMinister1ID)
	assert.NoError(t, err)
	oldDeptRelations2, err := client.GetAllRelatedEntities(oldMinister2ID)
	assert.NoError(t, err)

	oldActiveDepts1 := 0
	oldActiveDepts2 := 0
	for _, rel := range oldDeptRelations1 {
		if rel.Name == "AS_DEPARTMENT" && rel.EndTime == "" {
			oldActiveDepts1++
		}
	}
	for _, rel := range oldDeptRelations2 {
		if rel.Name == "AS_DEPARTMENT" && rel.EndTime == "" {
			oldActiveDepts2++
		}
	}

	assert.Equal(t, 0, oldActiveDepts1, "First old minister should have no active departments")
	assert.Equal(t, 0, oldActiveDepts2, "Second old minister should have no active departments")
}

func TestTerminateNonExistentMinister(t *testing.T) {
	// Create transaction map for terminating a non-existent minister
	transaction := map[string]interface{}{
		"parent":      "Government of Sri Lanka",
		"child":       "Non Existent Minister",
		"date":        "2025-01-01",
		"parent_type": "government",
		"child_type":  "minister",
		"rel_type":    "AS_MINISTER",
	}

	// Attempt to terminate the non-existent minister
	err := client.TerminateOrgEntity(transaction)
	assert.Error(t, err)
}

func TestTerminateMinisterWithChildren(t *testing.T) {
	// First create a minister with a department
	entityCounters := map[string]int{
		"minister":   0,
		"department": 0,
	}

	// Create minister
	ministerTransaction := map[string]interface{}{
		"parent":         "Government of Sri Lanka",
		"child":          "Minister to Terminate",
		"date":           "2025-01-01",
		"parent_type":    "government",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": "2154/14_tr_01",
	}

	_, err := client.AddOrgEntity(ministerTransaction, entityCounters)
	assert.NoError(t, err)

	// Create department under the minister
	departmentTransaction := map[string]interface{}{
		"parent":         "Minister to Terminate",
		"child":          "Department Under Minister",
		"date":           "2025-01-01",
		"parent_type":    "minister",
		"child_type":     "department",
		"rel_type":       "AS_DEPARTMENT",
		"transaction_id": "2154/14_tr_02",
	}

	_, err = client.AddOrgEntity(departmentTransaction, entityCounters)
	assert.NoError(t, err)

	// Debug: Print minister's relationships before termination
	ministerResults, err := client.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Minister to Terminate",
	})
	assert.NoError(t, err)
	assert.Len(t, ministerResults, 1)
	ministerID := ministerResults[0].ID

	fmt.Printf("Debug: Minister ID: %s\n", ministerID)
	relations, err := client.GetAllRelatedEntities(ministerID)
	assert.NoError(t, err)
	fmt.Printf("Debug: Minister's relationships before termination: %+v\n", relations)

	// Attempt to terminate the minister
	terminateTransaction := map[string]interface{}{
		"parent":      "Government of Sri Lanka",
		"child":       "Minister to Terminate",
		"date":        "2025-01-02",
		"parent_type": "government",
		"child_type":  "minister",
		"rel_type":    "AS_MINISTER",
	}

	fmt.Printf("Debug: Attempting to terminate minister with transaction: %+v\n", terminateTransaction)
	err = client.TerminateOrgEntity(terminateTransaction)
	assert.Error(t, err)
	// assert.Contains(t, err.Error(), "cannot terminate minister with active departments")

}

func TestMoveDepartmentToNonExistentMinister(t *testing.T) {
	// Create transaction map for moving department to non-existent minister
	transaction := map[string]interface{}{
		"old_parent": "Minister of Finance and Education",
		"new_parent": "Non Existent Minister",
		"child":      "Department of Policies",
		"type":       "department",
		"date":       "2025-01-01",
	}

	// Attempt to move the department
	err := client.MoveDepartment(transaction)
	assert.Error(t, err)
}

func TestMergeNonExistentMinister(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 0,
	}

	// Create transaction map for merging non-existent minister
	transaction := map[string]interface{}{
		"old":            "[Non Existent Minister]",
		"new":            "New Merged Minister",
		"date":           "2025-01-01",
		"transaction_id": "2154/14_tr_03",
	}

	// Attempt to merge the ministers
	_, err := client.MergeMinisters(transaction, entityCounters)
	fmt.Printf("Debug: Full error from MergeMinisters: %+v\n", err)
	assert.Error(t, err)
}

func TestCreateDuplicateMinister(t *testing.T) {
	// Initialize entity counters
	entityCounters := map[string]int{
		"minister": 0,
	}

	// Create transaction map for first minister
	firstMinisterTransaction := map[string]interface{}{
		"parent":         "Government of Sri Lanka",
		"child":          "Duplicate Minister",
		"date":           "2025-01-01",
		"parent_type":    "government",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": "2154/15_tr_01",
	}

	// Create the first minister
	firstMinister, err := client.AddOrgEntity(firstMinisterTransaction, entityCounters)
	assert.NoError(t, err)
	assert.NotNil(t, firstMinister)

	// Update counter for second attempt
	entityCounters["minister"]++

	// Create transaction map for second minister with same name
	secondMinisterTransaction := map[string]interface{}{
		"parent":         "Government of Sri Lanka",
		"child":          "Duplicate Minister",
		"date":           "2025-01-02",
		"parent_type":    "government",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": "2154/15_tr_02",
	}

	// Create the second minister with same name
	secondMinister, err := client.AddOrgEntity(secondMinisterTransaction, entityCounters)
	assert.NoError(t, err, "Should be able to create minister with same name but different ID")
	assert.NotNil(t, secondMinister)

	// Verify both ministers exist and have different IDs
	searchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "minister",
		},
		Name: "Duplicate Minister",
	}

	results, err := client.SearchEntities(searchCriteria)
	assert.NoError(t, err)
	assert.Len(t, results, 2, "Should find two ministers with this name")
	assert.NotEqual(t, results[0].ID, results[1].ID, "Ministers should have different IDs")

}
