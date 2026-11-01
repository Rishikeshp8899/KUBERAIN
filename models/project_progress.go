/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : project_progress.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing project progress milestones, completion
                   percentages, inspection document attachments, and stage-wise progress tracking.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateProjectProgress                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new project progress milestone record into the database.           ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding work package IDs, manufacturing stages, ║
║                     completion percentages, milestone notes, media paths, and inspector credentials.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateProjectProgress(p ProjectProgress) (int64, error) {
	query := `INSERT INTO project_progress (work_package_id, stage, percent_complete, milestone_note, photos_path, inspection_doc, inspector_name) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, p.WorkPackageID, p.Stage, p.PercentComplete, p.MilestoneNote, p.PhotosPath, p.InspectionDoc, p.InspectorName)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetProgressByWorkPackageID                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all progress milestone records associated with a specific work package.║
║ 🏛️ ARCHITECTURE   : Filtered relational query utilizing COALESCE wrappers for safe null handling  ║
║                     on milestone notes, photo paths, inspection documents, and inspector names,    ║
║                     sorted in descending order by update timestamps.                             ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetProgressByWorkPackageID(workPackageID int) ([]ProjectProgress, error) {
	query := `SELECT id, work_package_id, stage, percent_complete, COALESCE(milestone_note, ''), COALESCE(photos_path, ''), COALESCE(inspection_doc, ''), COALESCE(inspector_name, ''), updated_at 
              FROM project_progress WHERE work_package_id = ? ORDER BY updated_at DESC`
	rows, err := database.DB.Query(query, workPackageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ProjectProgress
	for rows.Next() {
		var p ProjectProgress
		if err := rows.Scan(&p.ID, &p.WorkPackageID, &p.Stage, &p.PercentComplete, &p.MilestoneNote, &p.PhotosPath, &p.InspectionDoc, &p.InspectorName, &p.UpdatedAt); err != nil {
			continue
		}
		list = append(list, p)
	}
	return list, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteProjectProgress                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a project progress milestone record by its primary key ID.║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the project progress table.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteProjectProgress(id int) error {
	_, err := database.DB.Exec(`DELETE FROM project_progress WHERE id = ?`, id)
	return err
}
