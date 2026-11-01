/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : admin_user.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing administrative authentication,
                   secure password hashing verification (bcrypt), and credential validation.
====================================================================================================
*/

package models

import (
	"database/sql"
	"errors"
	"go-mvc-app/database"

	"golang.org/x/crypto/bcrypt"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : AuthenticateAdmin                                                            ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Verifies administrative email credentials against stored bcrypt password hashes.║
║ 🏛️ ARCHITECTURE   : Secure two-phase authentication pipeline querying admin user records with      ║
║                     COALESCE role fallbacks, handling SQL row absence gracefully, and executing    ║
║                     cryptographic password hash comparisons using golang.org/x/crypto/bcrypt.       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func AuthenticateAdmin(email, password string) (*AdminUser, error) {
	query := `SELECT id, email, password_hash, COALESCE(role, 'admin') FROM admin_users WHERE email = ?`

	var user AdminUser
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare bcrypt password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
