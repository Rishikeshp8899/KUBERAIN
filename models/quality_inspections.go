/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : quality_inspection.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing quality control inspections, dimensional
                   measurements, tolerance compliance, inspection report attachments, and pass/fail audit results.
====================================================================================================
*/

package models

import (
	"go-mvc-app/constants"
	"go-mvc-app/database"
	"log"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateQualityInspection                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new quality inspection record into SQLite with complete audit data ║
║                     including WorkPackageID, SupplierID, inspection metrics, and file paths.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateQualityInspection(qi QualityInspection) (int64, error) {
	// 1. Log entry with all incoming fields including SupplierID
	log.Printf("[DEBUG] [CreateQualityInspection] Entry: WorkPackageID=%d, SupplierID=%d, InspectionType=%q, Passed=%t, MeasuredValues=%q, ToleranceAllowed=%q, ReportFilePath=%q",
		qi.WorkPackageID,
		qi.SupplierID,
		qi.InspectionType,
		qi.Passed,
		qi.MeasuredValues,
		qi.ToleranceAllowed,
		qi.ReportFilePath,
	)

	// 2. Added supplier_id to match the SQLite table schema constraints
	query := `
        INSERT INTO quality_inspections (
            work_package_id, 
            supplier_id, 
            inspection_type, 
            passed, 
            measured_values, 
            tolerance_allowed, 
            report_file_path,
            comments
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	// 3. Detailed payload parameter logging before execution
	log.Printf("[DEBUG] [CreateQualityInspection] Executing Query with Parameters -> WorkPackageID=%d, SupplierID=%d, InspectionType=%s, Passed=%t, MeasuredValues=%s, ToleranceAllowed=%s, ReportFilePath=%s, Comments=%s",
		qi.WorkPackageID,
		qi.SupplierID,
		qi.InspectionType,
		qi.Passed,
		qi.MeasuredValues,
		qi.ToleranceAllowed,
		qi.ReportFilePath,
		qi.Comments,
	)

	// 4. Execute insert statement with all 8 parameters corresponding to the columns above
	res, err := database.DB.Exec(
		query,
		qi.WorkPackageID,
		qi.SupplierID,
		qi.InspectionType,
		qi.Passed,
		qi.MeasuredValues,
		qi.ToleranceAllowed,
		qi.ReportFilePath,
		qi.Comments,
	)
	if err != nil {
		log.Printf("[ERROR] [CreateQualityInspection] DB Exec failed: %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("[ERROR] [CreateQualityInspection] Failed to retrieve LastInsertId: %v", err)
		return 0, err
	}

	log.Printf("[SUCCESS] [CreateQualityInspection] Successfully inserted quality inspection record with ID=%d", id)
	return id, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetInspectionsByWorkPackageID                                                ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all quality inspection records associated with a specific work package.║
║ 🏛️ ARCHITECTURE   : Filtered relational query utilizing COALESCE wrappers for safe null handling  ║
║                     on measured values, allowed tolerances, and inspection report file paths.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetInspectionsByWorkPackageID(packageID int) ([]QualityInspection, error) {
	log.Printf("[DEBUG] [GetInspectionsByWorkPackageID] Entry: packageID=%d", packageID)

	query := `SELECT id, work_package_id, inspection_type, passed, COALESCE(measured_values, ''), COALESCE(tolerance_allowed, ''), COALESCE(report_file_path, ''), inspected_at 
              FROM quality_inspections WHERE work_package_id = ?`

	rows, err := database.DB.Query(query, packageID)
	if err != nil {
		log.Printf("[ERROR] [GetInspectionsByWorkPackageID] Query failed for packageID %d: %v", packageID, err)
		return nil, err
	}
	defer rows.Close()

	var inspections []QualityInspection
	rowCount := 0
	for rows.Next() {
		var q QualityInspection
		if err := rows.Scan(&q.ID, &q.WorkPackageID, &q.InspectionType, &q.Passed, &q.MeasuredValues, &q.ToleranceAllowed, &q.ReportFilePath, &q.InspectedAt); err != nil {
			log.Printf("[WARN] [GetInspectionsByWorkPackageID] Row scan failed at index %d for packageID %d: %v", rowCount, packageID, err)
			continue
		}
		inspections = append(inspections, q)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetInspectionsByWorkPackageID] Row iteration error for packageID %d: %v", packageID, err)
		return nil, err
	}

	log.Printf("[SUCCESS] [GetInspectionsByWorkPackageID] Successfully fetched %d inspection records for packageID %d", len(inspections), packageID)
	return inspections, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteQualityInspection                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a quality inspection record by its primary key ID.        ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the quality inspections table. ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteQualityInspection(id int) error {
	log.Printf("[DEBUG] [DeleteQualityInspection] Entry: ID=%d", id)

	_, err := database.DB.Exec(`DELETE FROM quality_inspections WHERE id = ?`, id)
	if err != nil {
		log.Printf("[ERROR] [DeleteQualityInspection] Deletion failed for ID %d: %v", id, err)
		return err
	}

	log.Printf("[SUCCESS] [DeleteQualityInspection] Successfully deleted quality inspection record with ID=%d", id)
	return nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllQualityInspections                                                     ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all historical quality inspections/audits across all work packages.║
║ 🏛️ ARCHITECTURE   : Relational query joining quality_inspections with work_packages and suppliers  ║
║                     to present a complete audit trail ordered by most recent inspection date.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllQualityInspectionsByWorkPackage(id int) ([]QualityInspectionsDetails, error) {
	log.Printf("[DEBUG] [GetAllQualityInspections] Entry: Fetching all historical quality inspection audits")

	query := constants.QueryGetQualityInspectionsByWorkPackageID

	rows, err := database.DB.Query(query, id)
	if err != nil {
		log.Printf("[ERROR] [GetAllQualityInspections] Query execution failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var audits []QualityInspectionsDetails

	rowCount := 0
	for rows.Next() {
		var a QualityInspectionsDetails
		if err := rows.Scan(
			&a.ID,
			&a.WorkPackageID,
			&a.WorkPackageTitle,
			&a.SupplierID,
			&a.SupplierName,
			&a.InspectionType,
			&a.Passed,
			&a.Comments,
			&a.MeasuredValues,
			&a.ToleranceAllowed,
			&a.ReportFilePath,
			&a.InspectedAt,
		); err != nil {
			log.Printf("[WARN] [GetAllQualityInspections] Row scan failed at index %d: %v", rowCount, err)
			continue
		}
		audits = append(audits, a)
		rowCount++
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ERROR] [GetAllQualityInspections] Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] [GetAllQualityInspections] Successfully retrieved %d historical audit records", len(audits))
	return audits, nil
}
