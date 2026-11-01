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

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateQualityInspection                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new quality inspection record into SQLite.                         ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding work package IDs, inspection types,   ║
║                     pass/fail boolean statuses, measured values, tolerances, and report paths.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateQualityInspection(qi QualityInspection) (int64, error) {
	query := `INSERT INTO quality_inspections (work_package_id, inspection_type, passed, measured_values, tolerance_allowed, report_file_path) 
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, qi.WorkPackageID, qi.InspectionType, qi.Passed, qi.MeasuredValues, qi.ToleranceAllowed, qi.ReportFilePath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
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
	query := `SELECT id, work_package_id, inspection_type, passed, COALESCE(measured_values, ''), COALESCE(tolerance_allowed, ''), COALESCE(report_file_path, ''), inspected_at 
              FROM quality_inspections WHERE work_package_id = ?`
	rows, err := database.DB.Query(query, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []QualityInspection
	for rows.Next() {
		var q QualityInspection
		if err := rows.Scan(&q.ID, &q.WorkPackageID, &q.InspectionType, &q.Passed, &q.MeasuredValues, &q.ToleranceAllowed, &q.ReportFilePath, &q.InspectedAt); err != nil {
			continue
		}
		inspections = append(inspections, q)
	}
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
	_, err := database.DB.Exec(`DELETE FROM quality_inspections WHERE id = ?`, id)
	return err
}
