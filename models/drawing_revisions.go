/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : drawing_revision.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing engineering drawing revisions (ECN),
                   version activation toggles, and redacted file path traceability.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateDrawingRevision                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Creates a new drawing revision and deactivates prior active revisions.       ║
║ 🏛️ ARCHITECTURE   : Two-step update-then-insert pipeline that automatically deprecates existing  ║
║                     active revisions (`is_active = 0`) before registering the new revision with    ║
║                     `is_active = 1` for the specified RFQ ID.                                    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateDrawingRevision(rev DrawingRevision) (int64, error) {
	database.DB.Exec(`UPDATE drawing_revisions SET is_active = 0 WHERE rfq_id = ?`, rev.RFQID)

	query := `INSERT INTO drawing_revisions (rfq_id, revision_code, original_name, redacted_path, change_log_notes, is_active) 
              VALUES (?, ?, ?, ?, ?, 1)`
	res, err := database.DB.Exec(query, rev.RFQID, rev.RevisionCode, rev.OriginalName, rev.RedactedPath, rev.ChangeLogNotes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetActiveDrawingByRFQ                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches the currently active engineering drawing revision for an RFQ.        ║
║ 🏛️ ARCHITECTURE   : Point-lookup query filtering by both RFQ ID and active status (`is_active = 1`), ║
║                     utilizing COALESCE fallbacks for safe null handling on change log notes.       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetActiveDrawingByRFQ(rfqID int) (*DrawingRevision, error) {
	query := `SELECT id, rfq_id, revision_code, original_name, redacted_path, COALESCE(change_log_notes, ''), is_active, created_at 
              FROM drawing_revisions WHERE rfq_id = ? AND is_active = 1`
	var rev DrawingRevision
	err := database.DB.QueryRow(query, rfqID).Scan(&rev.ID, &rev.RFQID, &rev.RevisionCode, &rev.OriginalName, &rev.RedactedPath, &rev.ChangeLogNotes, &rev.IsActive, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteDrawingRevision                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a specific drawing revision record by ID.                ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the drawing revisions table.   ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteDrawingRevision(id int) error {
	_, err := database.DB.Exec(`DELETE FROM drawing_revisions WHERE id = ?`, id)
	return err
}
