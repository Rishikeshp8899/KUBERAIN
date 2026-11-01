/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : rfq.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing Request for Quotation (RFQ) pipelines,
                   unique alphanumeric code generation, status transitions, and client inquiries.
====================================================================================================
*/

package models

import (
	"fmt"
	"go-mvc-app/database"
	"log"
	"time"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateRFQ                                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Generates a unique tracking code and inserts a new RFQ record into SQLite.   ║
║ 🏛️ ARCHITECTURE   : Automatically computes timestamp-based alphanumeric RFQ codes if empty, binding  ║
║                     all corporate requirements, pricing targets, and default 'New' status values.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateRFQ(rfq RFQ) (int64, error) {
	if rfq.RFQCode == "" {
		rfq.RFQCode = fmt.Sprintf("RFQ-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	}

	query := `INSERT INTO rfqs (rfq_code, company_name, contact_person, email, phone, category, material, quantity, requirements, drawing_file_path, target_price_usd, delivery_deadline, incobooking_terms, country, city, ip_address, status) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, rfq.RFQCode, rfq.CompanyName, rfq.ContactPerson, rfq.Email, rfq.Phone, rfq.Category, rfq.Material, rfq.Quantity, rfq.Requirements, rfq.DrawingFilePath, rfq.TargetPriceUSD, rfq.DeliveryDeadline, rfq.IncobookingTerms, rfq.Country, rfq.City, rfq.IPAddress, "New")
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllRFQs                                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all corporate Request for Quotation records sorted descending.     ║
║ 🏛️ ARCHITECTURE   : Broad listing query utilizing COALESCE fallback defaults for safe null handling  ║
║                     on optional contact fields, terms, locations, and status properties.         ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllRFQs() ([]RFQ, error) {
	query := `SELECT id, COALESCE(rfq_code, 'N/A'), company_name, contact_person, email, COALESCE(phone, ''), COALESCE(category, ''), COALESCE(material, ''), COALESCE(quantity, ''), COALESCE(requirements, ''), COALESCE(drawing_file_path, ''), target_price_usd, COALESCE(delivery_deadline, ''), COALESCE(incobooking_terms, ''), COALESCE(country, 'India'), COALESCE(city, 'Unknown'), COALESCE(ip_address, 'N/A'), COALESCE(status, 'New'), created_at 
              FROM rfqs ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rfqs []RFQ
	for rows.Next() {
		var r RFQ
		if err := rows.Scan(&r.ID, &r.RFQCode, &r.CompanyName, &r.ContactPerson, &r.Email, &r.Phone, &r.Category, &r.Material, &r.Quantity, &r.Requirements, &r.DrawingFilePath, &r.TargetPriceUSD, &r.DeliveryDeadline, &r.IncobookingTerms, &r.Country, &r.City, &r.IPAddress, &r.Status, &r.CreatedAt); err != nil {
			continue
		}
		rfqs = append(rfqs, r)
	}
	return rfqs, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpdateRFQStatus                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Updates the workflow status string for a specified RFQ record.               ║
║ 🏛️ ARCHITECTURE   : Targeted SQL update statement for rapid pipeline state transitions without   ║
║                     reloading full entity graphs.                                                ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpdateRFQStatus(id int, status string) error {
	_, err := database.DB.Exec(`UPDATE rfqs SET status = ? WHERE id = ?`, status, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteRFQ                                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently removes an RFQ record from the database.                         ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the primary key identifier.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteRFQ(id int) error {
	_, err := database.DB.Exec(`DELETE FROM rfqs WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetRFQByID                                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches a single RFQ record by its primary key ID with error logging.        ║
║ 🏛️ ARCHITECTURE   : Targeted single-record point lookup query retrieving complete specification  ║
║                     and pricing parameters for detailed evaluation.                              ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetRFQByID(rfqID int) (RFQ, error) {
	var rfq RFQ
	query := `SELECT id, rfq_code, company_name, contact_person, email, phone, category, material, quantity, 
              requirements, drawing_file_path, target_price_usd, delivery_deadline, 
              incobooking_terms, country, city, ip_address, status, created_at 
              FROM rfqs WHERE id = ?`

	row := database.DB.QueryRow(query, rfqID)
	err := row.Scan(
		&rfq.ID, &rfq.RFQCode, &rfq.CompanyName, &rfq.ContactPerson, &rfq.Email, &rfq.Phone,
		&rfq.Category, &rfq.Material, &rfq.Quantity, &rfq.Requirements,
		&rfq.DrawingFilePath, &rfq.TargetPriceUSD, &rfq.DeliveryDeadline, &rfq.IncobookingTerms,
		&rfq.Country, &rfq.City, &rfq.IPAddress, &rfq.Status, &rfq.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] [GetRFQByID] Failed to fetch RFQ ID %d: %v", rfqID, err)
		return rfq, err
	}
	return rfq, nil
}
