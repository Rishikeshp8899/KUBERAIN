/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : client_user.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing enterprise client user accounts,
                   role-based permissions, contact details, and account lifecycle states.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateClientUser                                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new client user profile record into the database.                  ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding client IDs, full names, emails,     ║
║                     phones, roles, departments, and active flags, returning the generated ID.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateClientUser(user ClientUser) (int64, error) {
	query := `INSERT INTO client_users (client_id, full_name, email, phone, role, department, is_active) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, user.ClientID, user.FullName, user.Email, user.Phone, user.Role, user.Department, user.IsActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetUsersByClientID                                                           ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves client user accounts associated with a specific enterprise client. ║
║ 🏛️ ARCHITECTURE   : Filtered relational query utilizing COALESCE wrappers for safe null handling ║
║                     on phone numbers, departments, and last login timestamps.                    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetUsersByClientID(clientID int) ([]ClientUser, error) {
	query := `SELECT id, client_id, full_name, email, COALESCE(phone, ''), role, COALESCE(department, ''), is_active, COALESCE(last_login_at, ''), created_at 
              FROM client_users WHERE client_id = ?`
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []ClientUser
	for rows.Next() {
		var u ClientUser
		if err := rows.Scan(&u.ID, &u.ClientID, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.Department, &u.IsActive, &u.LastLoginAt, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteClientUser                                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a client user account record by its primary key ID.      ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the client users table.        ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteClientUser(userID int) error {
	_, err := database.DB.Exec(`DELETE FROM client_users WHERE id = ?`, userID)
	return err
}
