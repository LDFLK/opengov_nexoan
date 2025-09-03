package api

import (
	"fmt"
	"strings"
	"time"

	"orgchart_nexoan/models"
)

// CreateGovernmentNode creates the initial government node
func (c *Client) CreateGovernmentNode() (*models.Entity, error) {
	// Create the government entity
	governmentEntity := &models.Entity{
		ID:      "gov_01",
		Created: "1978-09-07T00:00:00Z",
		Kind: models.Kind{
			Major: "Organisation",
			Minor: "government",
		},
		Name: models.TimeBasedValue{
			StartTime: "1978-09-07T00:00:00Z",
			Value:     "Government of Sri Lanka",
		},
	}

	// Create the entity
	createdEntity, err := c.CreateEntity(governmentEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create government entity: %w", err)
	}

	return createdEntity, nil
}

// GetPresidentByGovernment retrieves a president entity (citizen with AS_PRESIDENT relationship to government) by name
func (c *Client) GetPresidentByGovernment(presidentName string) (*models.Entity, error) {
	// Get the president entity ID - presidents are citizens with AS_PRESIDENT relationship to government
	presidentResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Person",
			Minor: "citizen",
		},
		Name: presidentName,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search for president entity: %w", err)
	}
	if len(presidentResults) == 0 {
		return nil, fmt.Errorf("president entity not found: %s", presidentName)
	}

	// Find the president by checking if they have AS_PRESIDENT relationship to government
	for _, president := range presidentResults {
		// Get government node
		governmentResults, err := c.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: "government",
			},
		})
		if err != nil || len(governmentResults) == 0 {
			continue
		}

		// Check if this citizen has AS_PRESIDENT relationship to government
		presidentRelations, err := c.GetRelatedEntities(governmentResults[0].ID, &models.Relationship{
			Name:            "AS_PRESIDENT",
			RelatedEntityID: president.ID,
		})
		if err != nil {
			continue
		}

		// If there are any AS_PRESIDENT relationships (active or not), return the president
		if len(presidentRelations) > 0 {
			// Convert SearchResult to Entity
			entity := &models.Entity{
				ID:         president.ID,
				Kind:       president.Kind,
				Created:    president.Created,
				Terminated: president.Terminated,
				Name: models.TimeBasedValue{
					Value: president.Name,
				},
				Metadata:      []models.MetadataEntry{},
				Attributes:    []models.AttributeEntry{},
				Relationships: []models.RelationshipEntry{},
			}
			return entity, nil
		}
	}

	return nil, fmt.Errorf("president entity not found or not active: %s", presidentName)
}

// GetMinisterByPresident retrieves a minister entity by president name and minister name
func (c *Client) GetMinisterByPresident(presidentName, ministerName, dateISO string) (*models.Entity, error) {
	// Get the president entity using the helper function
	presidentEntity, err := c.GetPresidentByGovernment(presidentName)
	if err != nil {
		return nil, err
	}
	presidentID := presidentEntity.ID

	// Get all minister relationships for the president
	presidentRelations, err := c.GetRelatedEntities(presidentID, &models.Relationship{
		Name: "AS_MINISTER",
		//ActiveAt: dateISO,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get president's relationships: %w", err)
	}

	// Find the minister with the specified name
	for _, rel := range presidentRelations {
		// Fetch the related entity (minister)
		ministerResults, err := c.SearchEntities(&models.SearchCriteria{
			ID: rel.RelatedEntityID,
		})
		if err != nil || len(ministerResults) == 0 {
			continue
		}
		minister := ministerResults[0]
		if minister.Kind.Minor == "minister" && minister.Name == ministerName {
			// Convert SearchResult to Entity
			entity := &models.Entity{
				ID:         minister.ID,
				Kind:       minister.Kind,
				Created:    minister.Created,
				Terminated: minister.Terminated,
				Name: models.TimeBasedValue{
					Value: minister.Name,
				},
				Metadata:      []models.MetadataEntry{},
				Attributes:    []models.AttributeEntry{},
				Relationships: []models.RelationshipEntry{},
			}
			return entity, nil
		}
	}

	return nil, fmt.Errorf("minister '%s' not found under president '%s'", ministerName, presidentName)
}

// GetActiveMinisterByPresident retrieves an active minister entity by president name and minister name
// Returns an error if multiple active ministers with the same name are found
func (c *Client) GetActiveMinisterByPresident(presidentName, ministerName, dateISO string) (*models.Entity, error) {
	// Get the president entity using the helper function
	presidentEntity, err := c.GetPresidentByGovernment(presidentName)
	if err != nil {
		return nil, err
	}
	presidentID := presidentEntity.ID

	// Get all minister relationships for the president
	presidentRelations, err := c.GetRelatedEntities(presidentID, &models.Relationship{
		Name: "AS_MINISTER",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get president's relationships: %w", err)
	}

	// Find active ministers with the specified name
	var activeMinisters []*models.Entity
	for _, rel := range presidentRelations {
		// Only consider active relationships (EndTime == "")
		if rel.EndTime != "" {
			continue
		}

		// Fetch the related entity (minister)
		ministerResults, err := c.SearchEntities(&models.SearchCriteria{
			ID: rel.RelatedEntityID,
		})
		if err != nil || len(ministerResults) == 0 {
			continue
		}
		minister := ministerResults[0]
		if minister.Kind.Minor == "minister" && minister.Name == ministerName {
			// Convert SearchResult to Entity
			entity := &models.Entity{
				ID:         minister.ID,
				Kind:       minister.Kind,
				Created:    minister.Created,
				Terminated: minister.Terminated,
				Name: models.TimeBasedValue{
					Value: minister.Name,
				},
				Metadata:      []models.MetadataEntry{},
				Attributes:    []models.AttributeEntry{},
				Relationships: []models.RelationshipEntry{},
			}
			activeMinisters = append(activeMinisters, entity)
		}
	}

	// Check for multiple active ministers with the same name
	if len(activeMinisters) > 1 {
		return nil, fmt.Errorf("multiple active ministers found with name '%s' under president '%s'", ministerName, presidentName)
	}

	// Check if no active minister was found
	if len(activeMinisters) == 0 {
		return nil, fmt.Errorf("no active minister found with name '%s' under president '%s'", ministerName, presidentName)
	}

	return activeMinisters[0], nil
}

// AddOrgEntity creates a new entity and establishes its relationship with a parent entity.
// Assumes the parent entity already exists.
func (c *Client) AddOrgEntity(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)
	transactionID := transaction["transaction_id"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Generate new entity ID
	if _, exists := entityCounters[childType]; !exists {
		return 0, fmt.Errorf("unknown child type: %s", childType)
	}

	// Get the part before the first underscore for the prefix
	prefixPart := strings.Split(transactionID, "_")[0]
	prefix := fmt.Sprintf("%s_%s", prefixPart, strings.ToLower(childType[:3]))
	entityCounter := entityCounters[childType] + 1
	newEntityID := fmt.Sprintf("%s_%d", prefix, entityCounter)

	// Get the parent entity ID based on the child type
	var parentID string

	if childType == "minister" {
		// For ministers, parent should be a president (Person type) - presidents are citizens with AS_PRESIDENT relationship
		if parentType != "president" && parentType != "citizen" {
			return 0, fmt.Errorf("minister must be attached to a president, got parent_type: %s", parentType)
		}

		// Removed below: for now if a president creates the same minister again it will create a new entity
		// Check if minister already exists under this president
		// _, err := c.GetMinisterByPresident(parent, child, dateISO)
		// if err == nil {
		// 	// Minister already exists, return error
		// 	return 0, fmt.Errorf("minister '%s' already exists under president '%s'", child, parent)
		// }

		// Get the president entity
		presidentEntity, err := c.GetPresidentByGovernment(parent)
		if err != nil {
			return 0, fmt.Errorf("failed to get parent president entity: %w", err)
		}
		parentID = presidentEntity.ID

	} else if childType == "department" {
		// For departments, parent should be a minister, but we need to verify it's the correct minister
		if parentType != "minister" {
			return 0, fmt.Errorf("department must be attached to a minister, got parent_type: %s", parentType)
		}

		// Get president name from transaction
		presidentName, ok := transaction["president"].(string)
		if !ok || presidentName == "" {
			return 0, fmt.Errorf("president name is required and must be a non-empty string when adding a department")
		}

		// Check if a department with the same name already exists
		existingDepartmentResults, err := c.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: "department",
			},
			Name: child,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to search for existing department: %w", err)
		}
		if len(existingDepartmentResults) > 0 {
			return 0, fmt.Errorf("department with name '%s' already exists", child)
		}

		// Use GetMinisterByPresident to ensure we get the correct minister under the correct president
		ministerEntity, err := c.GetActiveMinisterByPresident(presidentName, parent, dateISO)
		if err != nil {
			return 0, fmt.Errorf("failed to get parent minister entity: %w", err)
		}

		parentID = ministerEntity.ID

	} else {
		// For other entity types, use the original logic
		majorType := "Organisation"
		if parentType == "citizen" {
			majorType = "Person"
		}
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: majorType,
				Minor: parentType,
			},
			Name: parent,
		}

		searchResults, err := c.SearchEntities(searchCriteria)
		if err != nil {
			return 0, fmt.Errorf("failed to search for parent entity: %w", err)
		}

		if len(searchResults) == 0 {
			return 0, fmt.Errorf("parent entity not found: %s", parent)
		}

		parentID = searchResults[0].ID
	}

	// Create the new child entity
	childEntity := &models.Entity{
		ID: newEntityID,
		Kind: models.Kind{
			Major: "Organisation",
			Minor: childType,
		},
		Created:    dateISO,
		Terminated: "",
		Name: models.TimeBasedValue{
			StartTime: dateISO,
			Value:     child,
		},
		Metadata:      []models.MetadataEntry{},
		Attributes:    []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{},
	}

	// Create the child entity
	createdChild, err := c.CreateEntity(childEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to create child entity: %w", err)
	}

	// Update the parent entity to add the relationship to the child
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", parentID, createdChild.ID, currentTimestamp)

	parentEntity := &models.Entity{
		ID:         parentID,
		Kind:       models.Kind{},
		Created:    "",
		Terminated: "",
		Name:       models.TimeBasedValue{},
		Metadata:   []models.MetadataEntry{},
		Attributes: []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: createdChild.ID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            relType,
				},
			},
		},
	}

	_, err = c.UpdateEntity(parentID, parentEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update parent entity: %w", err)
	}

	return entityCounter, nil
}

// TerminateOrgEntity terminates a specific relationship between parent and child at a given date
func (c *Client) TerminateOrgEntity(transaction map[string]interface{}) error {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the parent and child entity IDs based on their types
	var parentID, childID string

	// Handle parent entity retrieval
	if parentType == "president" {
		// Parent is a president - use the helper function
		presidentEntity, err := c.GetPresidentByGovernment(parent)
		if err != nil {
			return fmt.Errorf("failed to get parent president entity: %w", err)
		}
		parentID = presidentEntity.ID

	} else if parentType == "minister" {
		// Parent is a minister, need president context to get the correct minister
		presidentName, ok := transaction["president"].(string)
		if !ok || presidentName == "" {
			return fmt.Errorf("president name is required and must be a non-empty string when terminating minister relationships")
		}

		ministerEntity, err := c.GetActiveMinisterByPresident(presidentName, parent, dateISO)
		if err != nil {
			return fmt.Errorf("failed to get parent minister entity: %w", err)
		}
		parentID = ministerEntity.ID

	} else {
		// For other parent types, use the original logic
		parentMajorType := "Organisation"
		if parentType == "citizen" {
			parentMajorType = "Person"
		}
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: parentMajorType,
				Minor: parentType,
			},
			Name: parent,
		}

		parentResults, err := c.SearchEntities(searchCriteria)
		if err != nil {
			return fmt.Errorf("failed to search for parent entity: %w", err)
		}
		if len(parentResults) == 0 {
			return fmt.Errorf("parent entity not found: %s", parent)
		}
		parentID = parentResults[0].ID
	}

	// Handle child entity retrieval
	if childType == "minister" {
		// Child is a minister, parent is the president's name
		presidentName := parent // parent contains the president's name

		ministerEntity, err := c.GetActiveMinisterByPresident(presidentName, child, dateISO)
		if err != nil {
			return fmt.Errorf("failed to get child minister entity: %w", err)
		}
		childID = ministerEntity.ID

	} else if childType == "department" {
		// Child is a department, need to find it under the correct minister
		presidentName, ok := transaction["president"].(string)
		if !ok || presidentName == "" {
			return fmt.Errorf("president name is required and must be a non-empty string when terminating department relationships")
		}

		// First get the minister that should have this department
		ministerEntity, err := c.GetActiveMinisterByPresident(presidentName, parent, dateISO)
		if err != nil {
			return fmt.Errorf("failed to get minister for department termination: %w", err)
		}

		// Then find the department under this minister
		departmentRelations, err := c.GetRelatedEntities(ministerEntity.ID, &models.Relationship{
			Name: "AS_DEPARTMENT",
		})
		if err != nil {
			return fmt.Errorf("failed to get minister's department relationships: %w", err)
		}

		// Find the department with the matching name
		var foundDepartmentID string
		for _, rel := range departmentRelations {
			if rel.EndTime == "" { // Only active relationships
				departmentResults, err := c.SearchEntities(&models.SearchCriteria{ID: rel.RelatedEntityID})
				if err != nil || len(departmentResults) == 0 {
					continue
				}
				if departmentResults[0].Name == child {
					foundDepartmentID = rel.RelatedEntityID
					break
				}
			}
		}

		if foundDepartmentID == "" {
			return fmt.Errorf("department '%s' not found under minister '%s'", child, parent)
		}
		childID = foundDepartmentID

	} else {
		// For other child types, use the original logic
		childMajorType := "Organisation"
		if childType == "citizen" {
			childMajorType = "Person"
		}

		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: childMajorType,
				Minor: childType,
			},
			Name: child,
		}
		childResults, err := c.SearchEntities(searchCriteria)
		if err != nil {
			return fmt.Errorf("failed to search for child entity: %w", err)
		}
		if len(childResults) == 0 {
			return fmt.Errorf("child entity not found: %s", child)
		}
		childID = childResults[0].ID
	}

	//If we're terminating a minister, check for active departments
	// NOTE: Removing this to allow moving ministers
	// if childType == "minister" {
	// 	// Get all relationships for the minister
	// 	relations, err := c.GetRelatedEntities(childID, &models.Relationship{
	// 		Name: "AS_DEPARTMENT",
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get minister's relationships: %w", err)
	// 	}

	// 	// fmt.Println("relations: ", relations)

	// 	// Manually filter only active (EndTime == "") relationships
	// 	var activeRelations []models.Relationship
	// 	for _, rel := range relations {
	// 		if rel.EndTime == "" {
	// 			activeRelations = append(activeRelations, rel)
	// 		}
	// 	}

	// 	// Check for active departments
	// 	if len(activeRelations) > 0 {
	// 		return fmt.Errorf("cannot terminate minister with active departments")
	// 	}
	// }

	// Get the specific relationship that is still active (no end date) -> this should give us the relationship(s) active for dateISO
	relations, err := c.GetRelatedEntities(parentID, &models.Relationship{
		RelatedEntityID: childID,
		Name:            relType,
	})
	if err != nil {
		return fmt.Errorf("failed to get relationship: %w", err)
	}

	// FIXME: Is it possible to have more than one active relationship? For orgchart case only it won't happen
	// Find the active relationship (no end time)
	// Manually filter for active relationship (i.e., EndTime == "")
	var activeRel *models.Relationship
	for _, rel := range relations {
		if rel.EndTime == "" {
			activeRel = &rel
			break // stop at the first active one
		}
	}

	if activeRel == nil {
		return fmt.Errorf("no active relationship found between %s and %s with type %s", parentID, childID, relType)
	}

	// Update the relationship to set the end date
	_, err = c.UpdateEntity(parentID, &models.Entity{
		ID: parentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: activeRel.ID,
				Value: models.Relationship{
					EndTime: dateISO,
					ID:      activeRel.ID,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate relationship: %w", err)
	}

	return nil
}

// MoveDepartment moves a department from one minister to another
// MoveDepartment moves a department to a new minister
func (c *Client) MoveDepartment(transaction map[string]interface{}) error {
	// Extract details from the transaction
	newParent := transaction["new_parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Search for the department by name
	departmentResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: child,
	})
	if err != nil {
		return fmt.Errorf("failed to search for department: %w", err)
	}
	if len(departmentResults) == 0 {
		return fmt.Errorf("department '%s' not found", child)
	}
	if len(departmentResults) > 1 {
		return fmt.Errorf("multiple departments found with name '%s'", child)
	}
	departmentID := departmentResults[0].ID

	// Check for active incoming relationships to this department
	// Get all relationships where this department is the target
	departmentRelations, err := c.GetRelatedEntities(departmentID, &models.Relationship{
		Name: "AS_DEPARTMENT",
	})
	if err != nil {
		return fmt.Errorf("failed to get department relationships: %w", err)
	}

	// Look for active AS_DEPARTMENT relationships coming into this department
	for _, rel := range departmentRelations {
		if rel.EndTime == "" {
			// Found an active relationship - terminate it directly
			// We have the minister ID (rel.RelatedEntityID) and can terminate the relationship directly
			terminateRelationship := &models.Entity{
				ID: rel.RelatedEntityID, // This is the minister ID
				Relationships: []models.RelationshipEntry{
					{
						Key: rel.ID,
						Value: models.Relationship{
							EndTime: dateISO,
							ID:      rel.ID,
						},
					},
				},
			}

			_, err = c.UpdateEntity(rel.RelatedEntityID, terminateRelationship)
			if err != nil {
				return fmt.Errorf("failed to terminate old relationship: %w", err)
			}
			// Continue to terminate all active relationships, don't break
		}
	}

	// Get the new minister entity ID by president
	// We need the president name to get the correct minister
	newPresidentName, ok := transaction["new_president_name"].(string)
	if !ok || newPresidentName == "" {
		return fmt.Errorf("new_president_name is required and must be a non-empty string")
	}

	newMinisterEntity, err := c.GetActiveMinisterByPresident(newPresidentName, newParent, dateISO)
	if err != nil {
		return fmt.Errorf("failed to get new minister '%s' under president '%s': %w", newParent, newPresidentName, err)
	}
	newMinisterID := newMinisterEntity.ID

	// Create new AS_DEPARTMENT relationship from new minister to department
	// Use transaction ID and timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", newMinisterID, departmentID, currentTimestamp)

	newRelationship := &models.Entity{
		ID: newMinisterID,
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: departmentID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            "AS_DEPARTMENT",
				},
			},
		},
	}

	_, err = c.UpdateEntity(newMinisterID, newRelationship)
	if err != nil {
		return fmt.Errorf("failed to create new relationship: %w", err)
	}

	return nil
}

// RenameMinister renames a minister and transfers all its departments to the new minister
func (c *Client) RenameMinister(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	oldName := transaction["old"].(string)
	newName := transaction["new"].(string)
	dateStr := transaction["date"].(string)
	relType := "AS_MINISTER"
	transactionID := transaction["transaction_id"]

	// Validate president name is provided
	presidentName, ok := transaction["president"].(string)
	if !ok || presidentName == "" {
		return 0, fmt.Errorf("president name is required and must be a non-empty string")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the old minister's ID
	oldMinister, err := c.GetActiveMinisterByPresident(presidentName, oldName, dateISO)
	if err != nil {
		return 0, fmt.Errorf("failed to get old minister: %w", err)
	}
	oldMinisterID := oldMinister.ID

	// Create new minister
	addEntityTransaction := map[string]interface{}{
		"parent":         presidentName,
		"child":          newName,
		"date":           dateStr,
		"parent_type":    "president",
		"child_type":     "minister",
		"rel_type":       relType,
		"transaction_id": transactionID,
		"president":      presidentName,
	}

	// Create the new minister
	newMinisterCounter, err := c.AddOrgEntity(addEntityTransaction, entityCounters)
	if err != nil {
		return 0, fmt.Errorf("failed to create new minister: %w", err)
	}

	// Get the new minister's ID
	newMinister, err := c.GetActiveMinisterByPresident(presidentName, newName, dateISO)
	if err != nil {
		return 0, fmt.Errorf("failed to get new minister: %w", err)
	}
	newMinisterID := newMinister.ID

	// Get all active departments of the old minister
	oldRelations, err := c.GetRelatedEntities(oldMinisterID, &models.Relationship{
		Name: "AS_DEPARTMENT",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get old minister's relationships: %w", err)
	}

	// Manually filter only active relationships (EndTime == "")
	var oldActiveRelations []models.Relationship
	for _, rel := range oldRelations {
		if rel.EndTime == "" {
			oldActiveRelations = append(oldActiveRelations, rel)
		}
	}

	// Transfer each active department to the new minister using MoveDepartment
	for _, rel := range oldActiveRelations {
		// Get the department name using its ID
		departmentResults, err := c.SearchEntities(&models.SearchCriteria{
			ID: rel.RelatedEntityID,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to search for department: %w", err)
		}

		if len(departmentResults) == 0 {
			return 0, fmt.Errorf("failed to find department with ID: %s", rel.RelatedEntityID)
		}

		// Use MoveDepartment to move the department from old minister to new minister
		moveTransaction := map[string]interface{}{
			"old_parent":         oldName,
			"new_parent":         newName,
			"child":              departmentResults[0].Name,
			"date":               dateStr,
			"new_president_name": presidentName,
			"old_president_name": presidentName,
		}

		err = c.MoveDepartment(moveTransaction)
		if err != nil {
			return 0, fmt.Errorf("failed to move department: %w", err)
		}
	}

	// Terminate the old minister's relationship with the president directly
	// We need to get the president ID first
	presidentEntity, err := c.GetPresidentByGovernment(presidentName)
	if err != nil {
		return 0, fmt.Errorf("failed to get president entity: %w", err)
	}
	presidentID := presidentEntity.ID

	// Find the active relationship to terminate it
	presidentRelations, err := c.GetRelatedEntities(presidentID, &models.Relationship{
		Name:            "AS_MINISTER",
		RelatedEntityID: oldMinisterID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get relationship between president and minister: %w", err)
	}

	// Find the active relationship (EndTime == "")
	var activeRel *models.Relationship
	for _, rel := range presidentRelations {
		if rel.EndTime == "" {
			activeRel = &rel
			break
		}
	}

	if activeRel == nil {
		return 0, fmt.Errorf("no active relationship found between president and minister")
	}

	// Terminate the relationship directly
	terminateRelationship := &models.Entity{
		ID: presidentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: activeRel.ID,
				Value: models.Relationship{
					EndTime: dateISO,
					ID:      activeRel.ID,
				},
			},
		},
	}

	_, err = c.UpdateEntity(presidentID, terminateRelationship)
	if err != nil {
		return 0, fmt.Errorf("failed to terminate old minister's government relationship: %w", err)
	}

	// Create RENAMED_TO relationship
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", oldMinisterID, newMinisterID, currentTimestamp)

	renameRelationship := &models.Entity{
		ID: oldMinisterID,
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: newMinisterID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            "RENAMED_TO",
				},
			},
		},
	}

	_, err = c.UpdateEntity(oldMinisterID, renameRelationship)
	if err != nil {
		return 0, fmt.Errorf("failed to create RENAMED_TO relationship: %w", err)
	}

	return newMinisterCounter, nil
}

// RenameDepartment renames a department and transfers all its people relationships to the new department
func (c *Client) RenameDepartment(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	oldName := transaction["old"].(string)
	newName := transaction["new"].(string)
	dateStr := transaction["date"].(string)
	relType := "AS_DEPARTMENT"
	transactionID := transaction["transaction_id"].(string)
	presidentName, ok := transaction["president"].(string)
	if !ok || presidentName == "" {
		return 0, fmt.Errorf("president name is required and must be a non-empty string when renaming a department")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the old department's ID
	oldDepartmentResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: oldName,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to search for old department: %w", err)
	}
	if len(oldDepartmentResults) == 0 {
		return 0, fmt.Errorf("old department not found: %s", oldName)
	}
	oldDepartmentID := oldDepartmentResults[0].ID

	// Check if the new department name already exists
	existingDepartmentResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: "department",
		},
		Name: newName,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to search for new department name: %w", err)
	}

	var newDepartmentID string
	var newDepartmentCounter int

	if len(existingDepartmentResults) > 0 {
		// Check if the existing department has any active AS_DEPARTMENT relationships
		existingDepartment := existingDepartmentResults[0]
		existingDepartmentID := existingDepartment.ID

		// Get all AS_DEPARTMENT relationships for this department
		existingDepartmentRelations, err := c.GetRelatedEntities(existingDepartmentID, &models.Relationship{
			Name: "AS_DEPARTMENT",
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get existing department relationships: %w", err)
		}

		// Check if any relationships are still active (EndTime == "")
		hasActiveRelationships := false
		for _, rel := range existingDepartmentRelations {
			if rel.EndTime == "" {
				hasActiveRelationships = true
				break
			}
		}

		if hasActiveRelationships {
			// Department exists and has active relationships, cannot proceed
			return 0, fmt.Errorf("department with name '%s' already exists and has active relationships", newName)
		} else {
			// Department exists but all relationships are terminated, we can reuse it
			newDepartmentID = existingDepartment.ID
			//newDepartmentCounter = 0 // No new entity created, just reusing existing
		}
	} else {
		// Department doesn't exist, we'll create a new one
		newDepartmentID = ""
		//newDepartmentCounter = 0
	}

	// Get all active relationships coming into this department
	// The department can have multiple active relationships to different ministers from different presidents
	departmentRelations, err := c.GetRelatedEntities(oldDepartmentID, &models.Relationship{
		Name: "AS_DEPARTMENT",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get department relationships: %w", err)
	}

	// Find the minister that has an active relationship to this department under the specified president
	var ministerID string
	var ministerName string
	for _, rel := range departmentRelations {
		if rel.EndTime == "" {
			// This is an active relationship, check if the minister is under the correct president
			ministerResults, err := c.SearchEntities(&models.SearchCriteria{ID: rel.RelatedEntityID})
			if err != nil || len(ministerResults) == 0 {
				continue
			}
			minister := ministerResults[0]

			// Check if this minister is under the specified president
			_, err = c.GetActiveMinisterByPresident(presidentName, minister.Name, dateISO)
			if err == nil {
				// Found the minister under the correct president
				ministerID = minister.ID
				ministerName = minister.Name
				break
			}
		}
	}

	if ministerID == "" {
		return 0, fmt.Errorf("no active minister relationship found for department '%s' under president '%s'", oldName, presidentName)
	}

	// Verify that this minister is under the correct president
	// _, err = c.GetMinisterByPresident(presidentName, ministerName, dateISO)
	// if err != nil {
	// 	return 0, fmt.Errorf("minister '%s' not found under president '%s'", ministerName, presidentName)
	// }

	// Create new department or reuse existing inactive department
	if newDepartmentID == "" {
		// Create new department under the same minister
		addEntityTransaction := map[string]interface{}{
			"parent":         ministerName,
			"child":          newName,
			"date":           dateStr,
			"parent_type":    "minister",
			"child_type":     "department",
			"rel_type":       relType,
			"transaction_id": transactionID,
			"president":      presidentName,
		}

		// Create the new department
		newDepartmentCounter, err = c.AddOrgEntity(addEntityTransaction, entityCounters)
		if err != nil {
			return 0, fmt.Errorf("failed to create new department: %w", err)
		}

		// Get the new department's ID
		newDepartmentResults, err := c.SearchEntities(&models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: "department",
			},
			Name: newName,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to search for new department: %w", err)
		}
		if len(newDepartmentResults) == 0 {
			return 0, fmt.Errorf("new department not found: %s", newName)
		}
		if len(newDepartmentResults) > 1 {
			return 0, fmt.Errorf("multiple departments found with name '%s'", newName)
		}
		newDepartmentID = newDepartmentResults[0].ID
	} else {
		// Reusing existing inactive department - create the relationship with the minister
		// Generate a unique relationship ID
		newDepartmentCounter = entityCounters["department"]
		currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
		uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", ministerID, newDepartmentID, currentTimestamp)

		// Create the relationship between minister and the reactivated department
		reactivateRelationship := &models.Entity{
			ID: ministerID,
			Relationships: []models.RelationshipEntry{
				{
					Key: uniqueRelationshipID,
					Value: models.Relationship{
						RelatedEntityID: newDepartmentID,
						StartTime:       dateISO,
						EndTime:         "",
						ID:              uniqueRelationshipID,
						Name:            relType,
					},
				},
			},
		}

		_, err = c.UpdateEntity(ministerID, reactivateRelationship)
		if err != nil {
			return 0, fmt.Errorf("failed to create relationship with reactivated department: %w", err)
		}
	}

	// Terminate the old department's relationship with minister directly
	// Get the specific existing relationship to this department
	existingRelations, err := c.GetRelatedEntities(ministerID, &models.Relationship{
		Name:            "AS_DEPARTMENT",
		RelatedEntityID: oldDepartmentID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get existing relationship: %w", err)
	}

	// Find the active relationship (no end time)
	var existingRel *models.Relationship
	for _, rel := range existingRelations {
		if rel.EndTime == "" {
			existingRel = &rel
			break
		}
	}

	if existingRel == nil {
		return 0, fmt.Errorf("no active relationship found between minister '%s' and department '%s'", ministerID, oldDepartmentID)
	}

	// Terminate the relationship by updating it with the end time
	terminateRelationship := &models.Entity{
		ID: ministerID,
		Relationships: []models.RelationshipEntry{
			{
				Key: existingRel.ID,
				Value: models.Relationship{
					EndTime: dateISO,
					ID:      existingRel.ID,
				},
			},
		},
	}

	_, err = c.UpdateEntity(ministerID, terminateRelationship)
	if err != nil {
		return 0, fmt.Errorf("failed to terminate old department's minister relationship: %w", err)
	}

	// Create RENAMED_TO relationship
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", oldDepartmentID, newDepartmentID, currentTimestamp)

	renameRelationship := &models.Entity{
		ID: oldDepartmentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: newDepartmentID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            "RENAMED_TO",
				},
			},
		},
	}

	_, err = c.UpdateEntity(oldDepartmentID, renameRelationship)
	if err != nil {
		return 0, fmt.Errorf("failed to create RENAMED_TO relationship: %w", err)
	}

	return newDepartmentCounter, nil
}

// MergeMinisters merges multiple ministers into a new minister
func (c *Client) MergeMinisters(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	oldMinistersStr := transaction["old"].(string)
	newMinister := transaction["new"].(string)
	dateStr := transaction["date"].(string)
	transactionID := transaction["transaction_id"].(string)

	// Validate president name is provided
	presidentName, ok := transaction["president"].(string)
	if !ok || presidentName == "" {
		return 0, fmt.Errorf("president name is required and must be a non-empty string")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Parse old ministers list - use semicolons as separators to avoid comma conflicts
	trimmedStr := strings.Trim(oldMinistersStr, "[]")
	oldMinisters := strings.Split(trimmedStr, ";")
	for i := range oldMinisters {
		oldMinisters[i] = strings.TrimSpace(oldMinisters[i])
	}

	// 1. Create new minister using AddEntity
	addEntityTransaction := map[string]interface{}{
		"parent":         presidentName,
		"child":          newMinister,
		"date":           dateStr,
		"parent_type":    "president",
		"child_type":     "minister",
		"rel_type":       "AS_MINISTER",
		"transaction_id": transactionID,
		"president":      presidentName,
	}

	newMinisterCounter, err := c.AddOrgEntity(addEntityTransaction, entityCounters)
	if err != nil {
		return 0, fmt.Errorf("failed to create new minister: %w", err)
	}

	// Get the new minister's ID
	newMinisterEntity, err := c.GetActiveMinisterByPresident(presidentName, newMinister, dateISO)
	if err != nil {
		return 0, fmt.Errorf("failed to get new minister: %w", err)
	}
	newMinisterID := newMinisterEntity.ID

	// For each old minister
	for _, oldMinister := range oldMinisters {
		// Get the old minister's ID
		oldMinisterEntity, err := c.GetActiveMinisterByPresident(presidentName, oldMinister, dateISO)
		if err != nil {
			return 0, fmt.Errorf("failed to get old minister: %w", err)
		}
		oldMinisterID := oldMinisterEntity.ID

		// 2. Move old minister's departments to new minister
		oldRelations, err := c.GetRelatedEntities(oldMinisterID, &models.Relationship{
			Name: "AS_DEPARTMENT",
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get old minister's relationships: %w", err)
		}

		// Manually filter only active relationships (EndTime == "")
		var oldActiveRelations []models.Relationship
		for _, rel := range oldRelations {
			if rel.EndTime == "" {
				oldActiveRelations = append(oldActiveRelations, rel)
			}
		}

		for _, rel := range oldActiveRelations {
			// Get the department name using its ID
			departmentResults, err := c.SearchEntities(&models.SearchCriteria{
				ID: rel.RelatedEntityID,
			})
			if err != nil {
				return 0, fmt.Errorf("failed to search for department: %w", err)
			}
			if len(departmentResults) == 0 {
				return 0, fmt.Errorf("failed to find department with ID: %s", rel.RelatedEntityID)
			}

			// Move department to new minister
			moveTransaction := map[string]interface{}{
				"old_parent":         oldMinister,
				"new_parent":         newMinister,
				"child":              departmentResults[0].Name,
				"date":               dateStr,
				"new_president_name": presidentName,
				"old_president_name": presidentName,
			}

			err = c.MoveDepartment(moveTransaction)
			if err != nil {
				return 0, fmt.Errorf("failed to move department: %w", err)
			}
		}

		// 3. Terminate gov -> old minister relationship
		terminateGovTransaction := map[string]interface{}{
			"parent":      presidentName,
			"child":       oldMinister,
			"date":        dateStr,
			"parent_type": "citizen",
			"child_type":  "minister",
			"rel_type":    "AS_MINISTER",
		}

		err = c.TerminateOrgEntity(terminateGovTransaction)
		if err != nil {
			return 0, fmt.Errorf("failed to terminate old minister's government relationship: %w", err)
		}

		// 4. Create old minister -> new minister MERGED_INTO relationship
		// Use transaction ID and current timestamp to ensure unique relationship ID
		currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
		uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", oldMinisterID, newMinisterID, currentTimestamp)

		mergedIntoRelationship := &models.Entity{
			ID: oldMinisterID,
			Relationships: []models.RelationshipEntry{
				{
					Key: uniqueRelationshipID,
					Value: models.Relationship{
						RelatedEntityID: newMinisterID,
						StartTime:       dateISO,
						EndTime:         "",
						ID:              uniqueRelationshipID,
						Name:            "MERGED_INTO",
					},
				},
			},
		}

		_, err = c.UpdateEntity(oldMinisterID, mergedIntoRelationship)
		if err != nil {
			return 0, fmt.Errorf("failed to create MERGED_INTO relationship: %w", err)
		}
	}

	return newMinisterCounter, nil
}

// AddPersonEntity creates a new person entity and establishes its relationship with a parent entity.
// Assumes the parent entity already exists.
func (c *Client) AddPersonEntity(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)
	transactionID := transaction["transaction_id"].(string)

	// Get president name if parent is a minister -> currently only supports adding people to ministers
	var presidentName string
	if parentType == "minister" {
		var ok bool
		presidentName, ok = transaction["president"].(string)
		if !ok || presidentName == "" {
			return 0, fmt.Errorf("president name is required and must be a non-empty string when adding a person to a minister")
		}
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the parent entity ID
	var parentID string

	if parentType == "minister" {
		// Parent is a minister, need president context to get the correct minister
		ministerEntity, err := c.GetActiveMinisterByPresident(presidentName, parent, dateISO)
		if err != nil {
			return 0, fmt.Errorf("failed to get parent minister entity: %w", err)
		}
		parentID = ministerEntity.ID
	} else {
		// For other parent types, use the original logic
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				Minor: parentType,
			},
			Name: parent,
		}

		searchResults, err := c.SearchEntities(searchCriteria)
		if err != nil {
			return 0, fmt.Errorf("failed to search for parent entity: %w", err)
		}

		if len(searchResults) == 0 {
			return 0, fmt.Errorf("parent entity not found: %s", parent)
		}

		parentID = searchResults[0].ID
	}

	// Check if person already exists (search across all person types)
	personSearchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Person",
		},
		Name: child,
	}

	personResults, err := c.SearchEntities(personSearchCriteria)
	if err != nil {
		return 0, fmt.Errorf("failed to search for person entity: %w", err)
	}

	if len(personResults) > 1 {
		return 0, fmt.Errorf("multiple entities found for person: %s", child)
	}

	var childID string
	if len(personResults) == 1 {
		// Person exists, use existing ID
		childID = personResults[0].ID
	} else {
		// Generate new entity ID
		if _, exists := entityCounters[childType]; !exists {
			return 0, fmt.Errorf("unknown child type: %s", childType)
		}

		// Get the part before the first underscore for the prefix
		prefixPart := strings.Split(transactionID, "_")[0]
		prefix := fmt.Sprintf("%s_%s", prefixPart, strings.ToLower(childType[:3]))
		entityCounters[childType]++ // Increment the counter
		newEntityID := fmt.Sprintf("%s_%d", prefix, entityCounters[childType])

		// Create the new child entity
		childEntity := &models.Entity{
			ID: newEntityID,
			Kind: models.Kind{
				Major: "Person",
				Minor: childType,
			},
			Created:    dateISO,
			Terminated: "",
			Name: models.TimeBasedValue{
				StartTime: dateISO,
				Value:     child,
			},
			Metadata:      []models.MetadataEntry{},
			Attributes:    []models.AttributeEntry{},
			Relationships: []models.RelationshipEntry{},
		}

		// Create the child entity
		createdChild, err := c.CreateEntity(childEntity)
		if err != nil {
			return 0, fmt.Errorf("failed to create child entity: %w", err)
		}
		childID = createdChild.ID
	}

	// Update the parent entity to add the relationship to the child
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", parentID, childID, currentTimestamp)

	parentEntity := &models.Entity{
		ID:         parentID,
		Kind:       models.Kind{},
		Created:    "",
		Terminated: "",
		Name:       models.TimeBasedValue{},
		Metadata:   []models.MetadataEntry{},
		Attributes: []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: childID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            relType,
				},
			},
		},
	}

	_, err = c.UpdateEntity(parentID, parentEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update parent entity: %w", err)
	}

	return entityCounters[childType], nil
}

// TerminatePersonEntity terminates a specific relationship between Person type entity and another entity at a given date
func (c *Client) TerminatePersonEntity(transaction map[string]interface{}) error {
	// Extract details from the transaction
	parent := transaction["parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	parentType := transaction["parent_type"].(string)
	childType := transaction["child_type"].(string)
	relType := transaction["rel_type"].(string)

	// Get president name if parent is a minister -> currently only supports terminating relationships with ministers
	var presidentName string
	if parentType == "minister" {
		var ok bool
		presidentName, ok = transaction["president"].(string)
		if !ok || presidentName == "" {
			return fmt.Errorf("president name is required and must be a non-empty string when terminating relationships with ministers")
		}
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// First, find the person (child) entity
	childSearchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Person",
			Minor: childType,
		},
		Name: child,
	}

	childResults, err := c.SearchEntities(childSearchCriteria)
	if err != nil {
		return fmt.Errorf("failed to search for child entity: %w", err)
	}
	if len(childResults) == 0 {
		return fmt.Errorf("child entity not found: %s", child)
	}
	childID := childResults[0].ID

	// Find the ministry by checking the person's active relationships
	var parentID string
	var activeRel *models.Relationship

	if parentType == "minister" {
		// Get all active relationships from the person to find the ministry
		personRelations, err := c.GetRelatedEntities(childID, &models.Relationship{
			Name: relType,
		})
		if err != nil {
			return fmt.Errorf("failed to get person's relationships: %w", err)
		}

		// Find the active relationship to the specified ministry
		for _, rel := range personRelations {
			if rel.EndTime == "" { // Only active relationships
				// Get the ministry entity to check if it matches the parent name
				ministryResults, err := c.SearchEntities(&models.SearchCriteria{
					ID: rel.RelatedEntityID,
				})
				if err != nil || len(ministryResults) == 0 {
					continue
				}

				ministry := ministryResults[0]
				// Check if this ministry is under the correct president and matches the parent name
				if ministry.Kind.Minor == "minister" && ministry.Name == parent {
					// Verify this minister is under the specified president
					_, err = c.GetActiveMinisterByPresident(presidentName, parent, dateISO)
					if err == nil {
						parentID = ministry.ID
						// We found the ministry, now we need to get the relationship from ministry to person
						// to get the relationship ID for termination
						ministryRelations, err := c.GetRelatedEntities(ministry.ID, &models.Relationship{
							Name:            relType,
							RelatedEntityID: childID,
						})
						if err != nil {
							return fmt.Errorf("failed to get ministry's relationship to person: %w", err)
						}

						// Find the active relationship
						for _, mRel := range ministryRelations {
							if mRel.EndTime == "" {
								activeRel = &mRel
								break
							}
						}
						break
					}
				}
			}
		}

		if parentID == "" {
			return fmt.Errorf("no active relationship found between person '%s' and ministry '%s' under president '%s'", child, parent, presidentName)
		}
	} else {
		// For other parent types, use the original logic
		searchCriteria := &models.SearchCriteria{
			Kind: &models.Kind{
				Major: "Organisation",
				//Minor: parentType,
			},
			Name: parent,
		}
		parentResults, err := c.SearchEntities(searchCriteria)
		if err != nil {
			return fmt.Errorf("failed to search for parent entity: %w", err)
		}
		if len(parentResults) == 0 {
			return fmt.Errorf("parent entity not found: %s", parent)
		}
		parentID = parentResults[0].ID
	}

	// If we haven't found the active relationship yet (for non-minister parent types), search for it
	if activeRel == nil {
		// Get the specific relationship that is still active (no end date)
		relations, err := c.GetRelatedEntities(parentID, &models.Relationship{
			RelatedEntityID: childID,
			Name:            relType,
		})
		if err != nil {
			return fmt.Errorf("failed to get relationship: %w", err)
		}

		// FIXME: Is it possible to have more than one active relationship? For orgchart case only it won't happen
		// Manually filter only active relationships (EndTime == "")
		var activeRelations []models.Relationship
		for _, rel := range relations {
			if rel.EndTime == "" {
				activeRelations = append(activeRelations, rel)
			}
		}

		// Find the active relationship (no end time)
		if len(activeRelations) > 0 {
			activeRel = &activeRelations[0]
		}
	}

	if activeRel == nil {
		return fmt.Errorf("no active relationship found between %s and %s with type %s", parentID, childID, relType)
	}

	// Update the relationship to set the end date
	_, err = c.UpdateEntity(parentID, &models.Entity{
		ID: parentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: activeRel.ID,
				Value: models.Relationship{
					EndTime: dateISO,
					ID:      activeRel.ID,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate relationship: %w", err)
	}

	return nil
}

// MovePerson moves a person from one portfolio to another (limits functionality to only minister)
// TODO: Take the parent type from the transaction such that this function can be used generic
//
//	for moving person from any institution to another
func (c *Client) MovePerson(transaction map[string]interface{}) error {
	// Extract details from the transaction
	newParent := transaction["new_parent"].(string)
	oldParent := transaction["old_parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)
	relType := "AS_APPOINTED"

	// Validate president name is provided
	presidentName, ok := transaction["president"].(string)
	if !ok || presidentName == "" {
		return fmt.Errorf("president name is required and must be a non-empty string")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the new minister (parent) entity ID -> only supports moving person to and from minister
	newParentEntity, err := c.GetActiveMinisterByPresident(presidentName, newParent, dateISO)
	if err != nil {
		return fmt.Errorf("failed to get new parent entity: %w", err)
	}
	newParentID := newParentEntity.ID

	// Get the department (child) entity ID
	childResults, err := c.SearchEntities(&models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Person",
			Minor: "citizen",
		},
		Name: child,
	})
	if err != nil {
		return fmt.Errorf("failed to search for child entity: %w", err)
	}
	if len(childResults) == 0 {
		return fmt.Errorf("child entity not found: %s", child)
	}
	childID := childResults[0].ID

	// Create new relationship between new minister and person
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", newParentID, childID, currentTimestamp)

	newRelationship := &models.Entity{
		ID: newParentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: childID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            relType,
				},
			},
		},
	}

	_, err = c.UpdateEntity(newParentID, newRelationship)
	if err != nil {
		return fmt.Errorf("failed to create new relationship: %w", err)
	}

	// Terminate the old relationship
	terminateTransaction := map[string]interface{}{
		"parent":      oldParent,
		"child":       child,
		"date":        dateStr,
		"parent_type": "minister",
		"child_type":  "citizen",
		"rel_type":    relType,
		"president":   presidentName,
	}

	err = c.TerminatePersonEntity(terminateTransaction)
	if err != nil {
		return fmt.Errorf("failed to terminate old relationship: %w", err)
	}

	return nil
}

// MoveMinister moves a minister from one president to another
// func (c *Client) MoveMinister(transaction map[string]interface{}) error {
// 	// Extract details from the transaction
// 	newParent := transaction["new_parent"].(string)
// 	oldParent := transaction["old_parent"].(string)
// 	child := transaction["child"].(string)
// 	dateStr := transaction["date"].(string)

// 	// Parse the date
// 	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
// 	if err != nil {
// 		return fmt.Errorf("failed to parse date: %w", err)
// 	}
// 	dateISO := date.Format(time.RFC3339)

// 	// Get the minister entity ID from the old president
// 	ministerEntity, err := c.GetMinisterByPresident(oldParent, child, dateISO)
// 	if err != nil {
// 		return fmt.Errorf("minister entity '%s' not found or not active under old president '%s' on date %s: %w", child, oldParent, dateStr, err)
// 	}
// 	ministerID := ministerEntity.ID

// 	// First: Terminate the old relationship between old president and minister
// 	terminateTransaction := map[string]interface{}{
// 		"parent":      oldParent,
// 		"child":       child,
// 		"date":        dateStr,
// 		"parent_type": "president",
// 		"child_type":  "minister",
// 		"rel_type":    "AS_MINISTER",
// 	}

// 	err = c.TerminateOrgEntity(terminateTransaction)
// 	if err != nil {
// 		return fmt.Errorf("failed to terminate old relationship: %w", err)
// 	}

// 	// Get the new president entity ID
// 	newPresidentEntity, err := c.GetPresidentByGovernment(newParent)
// 	if err != nil {
// 		return fmt.Errorf("failed to get new president entity: %w", err)
// 	}
// 	newPresidentID := newPresidentEntity.ID

// 	// Then: Create new relationship between new president and minister
// 	newRelationship := &models.Entity{
// 		ID: newPresidentID,
// 		Relationships: []models.RelationshipEntry{
// 			{
// 				Key: fmt.Sprintf("%s_%s", newPresidentID, ministerID),
// 				Value: models.Relationship{
// 					RelatedEntityID: ministerID,
// 					StartTime:       dateISO,
// 					EndTime:         "",
// 					ID:              fmt.Sprintf("%s_%s", newPresidentID, ministerID),
// 					Name:            "AS_MINISTER",
// 				},
// 			},
// 		},
// 	}

// 	_, err = c.UpdateEntity(newPresidentID, newRelationship)
// 	if err != nil {
// 		return fmt.Errorf("failed to create new relationship: %w", err)
// 	}

// 	return nil
// }

// MoveMinister moves a minister from one president to another
func (c *Client) MoveMinister(transaction map[string]interface{}) error {
	// Extract details from the transaction
	newParent := transaction["new_parent"].(string)
	oldParent := transaction["old_parent"].(string)
	child := transaction["child"].(string)
	dateStr := transaction["date"].(string)

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// --- Get the new president (parent) entity ID ---
	newPresidentEntity, err := c.GetPresidentByGovernment(newParent)
	if err != nil {
		return fmt.Errorf("failed to get new president entity: %w", err)
	}
	newParentID := newPresidentEntity.ID

	// --- Get the old president (parent) entity ID ---
	oldPresidentEntity, err := c.GetPresidentByGovernment(oldParent)
	if err != nil {
		return fmt.Errorf("failed to get old president entity: %w", err)
	}
	oldParentID := oldPresidentEntity.ID

	// Get the minister (child) entity ID connected to the old president
	ministerEntity, err := c.GetMinisterByPresident(oldParent, child, dateISO)
	if err != nil {
		return fmt.Errorf("minister entity '%s' not found or not active under old president '%s' on date %s: %w", child, oldParent, dateStr, err)
	}
	childID := ministerEntity.ID

	// Create new relationship between new president and minister
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", newParentID, childID, currentTimestamp)

	newRelationship := &models.Entity{
		ID: newParentID,
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: childID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            "AS_MINISTER",
				},
			},
		},
	}

	_, err = c.UpdateEntity(newParentID, newRelationship)
	if err != nil {
		return fmt.Errorf("failed to create new relationship: %w", err)
	}

	// Find the active relationship to terminate it.
	oldPresidentRelations, err := c.GetRelatedEntities(oldParentID, &models.Relationship{
		Name:            "AS_MINISTER",
		RelatedEntityID: childID,
	})
	if err != nil {
		return fmt.Errorf("failed to get relationship between old president and minister: %w", err)
	}

	// Manually filter for active relationships (EndTime == "")
	var activeRel *models.Relationship
	for _, rel := range oldPresidentRelations {
		if rel.EndTime == "" {
			activeRel = &rel
			break
		}
	}

	// Only terminate if there is an active relationship
	if activeRel != nil {
		// Terminate the old relationship
		terminateTransaction := map[string]interface{}{
			"parent":      oldParent,
			"child":       child,
			"date":        dateStr,
			"parent_type": "president",
			"child_type":  "minister",
			"rel_type":    "AS_MINISTER",
		}

		err = c.TerminateOrgEntity(terminateTransaction)
		if err != nil {
			return fmt.Errorf("failed to terminate old relationship: %w", err)
		}
	}

	return nil
}

// Document Entity Handling
// Unlike other entities, Documents are not terminated, but there is an aspect to a document being
// regarded in various states. So this needs to be thoroughly thought and represented in the system.
// For now we are only adding the documents and not maintaining any other states.

// AddDocumentEntity creates a new document entity and establishes its relationship with a parent entity.
// The document type is determined by the parent entity type (Organization or Person).
// Assumes the parent entity already exists.
func (c *Client) AddDocumentEntity(transaction map[string]interface{}, entityCounters map[string]int) (int, error) {
	// Extract details from the transaction with validation
	parent, ok := transaction["parent"].(string)
	if !ok || parent == "" {
		return 0, fmt.Errorf("parent is required and must be a string")
	}

	child, ok := transaction["child"].(string)
	if !ok || child == "" {
		return 0, fmt.Errorf("child is required and must be a string")
	}

	dateStr, ok := transaction["date"].(string)
	if !ok || dateStr == "" {
		return 0, fmt.Errorf("date is required and must be a string")
	}

	parentType, ok := transaction["parent_type"].(string)
	if !ok || parentType == "" {
		return 0, fmt.Errorf("parent_type is required and must be a string")
	}

	childType, ok := transaction["child_type"].(string)
	if !ok || childType == "" {
		return 0, fmt.Errorf("child_type is required and must be a string")
	}

	transactionID, ok := transaction["transaction_id"].(string)
	if !ok || transactionID == "" {
		return 0, fmt.Errorf("transaction_id is required and must be a string")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	dateISO := date.Format(time.RFC3339)

	// Get the parent entity ID (which is always gonna be an organisation)
	searchCriteria := &models.SearchCriteria{
		Name: parent,
		Kind: &models.Kind{
			Major: "Organisation",
			Minor: parentType,
		},
	}

	searchResults, err := c.SearchEntities(searchCriteria)
	if err != nil {
		return 0, fmt.Errorf("failed to search for parent entity: %w", err)
	}

	if len(searchResults) == 0 {
		return 0, fmt.Errorf("parent entity not found: %s", parent)
	}

	parentID := searchResults[0].ID

	// Check if document already exists
	documentSearchCriteria := &models.SearchCriteria{
		Kind: &models.Kind{
			Major: "Document",
			Minor: childType,
		},
		Name: child,
	}

	documentResults, err := c.SearchEntities(documentSearchCriteria)
	if err != nil {
		return 0, fmt.Errorf("failed to search for document entity: %w", err)
	}

	if len(documentResults) > 1 {
		return 0, fmt.Errorf("multiple entities found for document: %s", child)
	}

	var childID string
	entityCounter := 0
	if len(documentResults) == 1 {
		// Document exists, use existing ID
		childID = documentResults[0].ID
	} else {
		// Generate new entity ID
		// Get the part before the first underscore for the prefix
		prefixPart := strings.Split(transactionID, "_")[0]
		prefix := fmt.Sprintf("%s_doc", prefixPart)
		entityCounter = entityCounters["document"] + 1
		newEntityID := fmt.Sprintf("%s_%d", prefix, entityCounter)

		// Create the new document entity
		documentEntity := &models.Entity{
			ID: newEntityID,
			Kind: models.Kind{
				Major: "Document",
				Minor: childType,
			},
			Created:    dateISO,
			Terminated: "",
			Name: models.TimeBasedValue{
				StartTime: dateISO,
				Value:     child,
			},
			Metadata:      []models.MetadataEntry{},
			Attributes:    []models.AttributeEntry{},
			Relationships: []models.RelationshipEntry{},
		}

		// Create the document entity
		createdDocument, err := c.CreateEntity(documentEntity)
		if err != nil {
			return 0, fmt.Errorf("failed to create document entity: %w", err)
		}
		childID = createdDocument.ID
	}

	// Update the parent entity to add the relationship to the document
	// Use transaction ID and current timestamp to ensure unique relationship ID
	currentTimestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	uniqueRelationshipID := fmt.Sprintf("%s_%s_%s", parentID, childID, currentTimestamp)

	parentEntity := &models.Entity{
		ID:         parentID,
		Kind:       models.Kind{},
		Created:    "",
		Terminated: "",
		Name:       models.TimeBasedValue{},
		Metadata:   []models.MetadataEntry{},
		Attributes: []models.AttributeEntry{},
		Relationships: []models.RelationshipEntry{
			{
				Key: uniqueRelationshipID,
				Value: models.Relationship{
					RelatedEntityID: childID,
					StartTime:       dateISO,
					EndTime:         "",
					ID:              uniqueRelationshipID,
					Name:            "AS_DOCUMENT",
				},
			},
		},
	}

	_, err = c.UpdateEntity(parentID, parentEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update parent entity: %w", err)
	}

	return entityCounter, nil
}
