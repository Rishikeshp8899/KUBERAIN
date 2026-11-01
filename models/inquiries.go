/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : inquiry.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing contact form submissions, public
                   client inquiries, tracking metadata, IP logging, and administrative retrieval queries.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateInquiry                                                                ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new public contact form inquiry record into SQLite.                ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding submitter names, emails, phones,    ║
║                     subjects, messages, lead sources, and client IP addresses, returning ID.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateInquiry(inq Inquiry) (int64, error) {
	query := `INSERT INTO inquiries (name, email, phone, subject, message, source, ip_address) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, inq.Name, inq.Email, inq.Phone, inq.Subject, inq.Message, inq.Source, inq.IPAddress)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllInquiries                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all public inquiries sorted descending by ID for administrative review.║
║ 🏛️ ARCHITECTURE   : Broad listing query utilizing COALESCE fallback defaults for safe null handling  ║
║                     on optional phones, subjects, sources ('Contact Form'), and IP addresses.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllInquiries() ([]Inquiry, error) {
	query := `SELECT id, name, email, COALESCE(phone, ''), COALESCE(subject, ''), message, COALESCE(source, 'Contact Form'), COALESCE(ip_address, ''), created_at 
              FROM inquiries ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inquiries []Inquiry
	for rows.Next() {
		var inq Inquiry
		if err := rows.Scan(&inq.ID, &inq.Name, &inq.Email, &inq.Phone, &inq.Subject, &inq.Message, &inq.Source, &inq.IPAddress, &inq.CreatedAt); err != nil {
			continue
		}
		inquiries = append(inquiries, inq)
	}
	return inquiries, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteInquiry                                                                ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a contact inquiry record from the database.              ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the inquiries primary key.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteInquiry(id int) error {
	_, err := database.DB.Exec(`DELETE FROM inquiries WHERE id = ?`, id)
	return err
}
