/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : communication_log.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing secure sandbox communications,
                   anti-leak regex detection (phone/email filters), and disintermediation audit logs.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
	"log"
	"regexp"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateCommunicationLog                                                       ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new communication log record into the database.                    ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding work package IDs, sender roles,       ║
║                     sender IDs, message strings, and flag statuses, returning the generated ID.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateCommunicationLog(logData CommunicationLog) (int64, error) {
	query := `INSERT INTO communication_logs (work_package_id, sender_role, sender_id, message, is_flagged, flag_reason) 
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, logData.WorkPackageID, logData.SenderRole, logData.SenderID, logData.Message, logData.IsFlagged, logData.FlagReason)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : AddSecureMessage                                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Scans messages for anti-leak rules (phone/email) and logs secure communications.║
║ 🏛️ ARCHITECTURE   : Pre-insertion filtering pipeline utilizing compiled regular expressions      ║
║                     (RegEx) to automatically flag and block disintermediation attempts before    ║
║                     persisting audit records with live timestamps.                               ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func AddSecureMessage(wpID int, senderRole string, senderID int, message string) error {
	isFlagged := false
	flagReason := ""

	// Anti-leak regex rules for phone numbers and email addresses
	phoneRegex := regexp.MustCompile(`(\+?\d{10,13})|(\b\d{3}[-.\s]?\d{3}[-.\s]?\d{4}\b)`)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	if phoneRegex.MatchString(message) {
		isFlagged = true
		flagReason = "Phone number pattern detected and blocked in secure sandbox chat."
	} else if emailRegex.MatchString(message) {
		isFlagged = true
		flagReason = "Email address pattern detected and blocked in secure sandbox chat."
	}

	query := `INSERT INTO communication_logs (work_package_id, sender_role, sender_id, message, is_flagged, flag_reason, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`

	_, err := database.DB.Exec(query, wpID, senderRole, senderID, message, isFlagged, flagReason)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteCommunicationLog                                                       ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a communication log entry by its primary key ID.         ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the communication logs table.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteCommunicationLog(id int) error {
	_, err := database.DB.Exec(`DELETE FROM communication_logs WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetFlaggedCommunicationLogs                                                  ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches disintermediation and leak detection audit logs.                     ║
║ 🏛️ ARCHITECTURE   : Broad listing query ordered in descending chronological sequence by creation ║
║                     timestamp to surface security violations and flagged messages immediately.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetFlaggedCommunicationLogs() ([]FlaggedLog, error) {
	query := `SELECT id, work_package_id, sender_role, sender_id, flag_reason, message, created_at FROM communication_logs ORDER BY created_at DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetFlaggedCommunicationLogs] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []FlaggedLog
	for rows.Next() {
		var l FlaggedLog
		if err := rows.Scan(&l.ID, &l.WorkPackageID, &l.SenderRole, &l.SenderID, &l.FlagReason, &l.Message, &l.CreatedAt); err != nil {
			continue
		}
		logs = append(logs, l)
	}
	return logs, nil
}
