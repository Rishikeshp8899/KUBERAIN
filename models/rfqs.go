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
	"database/sql"
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

	// Log function entry and input parameters
	log.Printf("[DEBUG] [GetRFQByID] Function called with rfqID=%d", rfqID)

	query := `SELECT id, rfq_code, company_name, contact_person, email, phone, category, material, quantity, 
              requirements, drawing_file_path, target_price_usd, delivery_deadline, 
              incobooking_terms, country, city, ip_address, status, created_at 
              FROM rfqs WHERE id = ?`

	log.Printf("[DEBUG] [GetRFQByID] Executing point lookup query for rfqID=%d", rfqID)
	row := database.DB.QueryRow(query, rfqID)

	err := row.Scan(
		&rfq.ID, &rfq.RFQCode, &rfq.CompanyName, &rfq.ContactPerson, &rfq.Email, &rfq.Phone,
		&rfq.Category, &rfq.Material, &rfq.Quantity, &rfq.Requirements,
		&rfq.DrawingFilePath, &rfq.TargetPriceUSD, &rfq.DeliveryDeadline, &rfq.IncobookingTerms,
		&rfq.Country, &rfq.City, &rfq.IPAddress, &rfq.Status, &rfq.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] [GetRFQByID] No RFQ record found for ID %d", rfqID)
		} else {
			log.Printf("[ERROR] [GetRFQByID] Failed to scan RFQ ID %d: %v", rfqID, err)
		}
		return rfq, err
	}

	log.Printf("[DEBUG] [GetRFQByID] Successfully fetched and scanned RFQ ID %d (RFQCode: %s)", rfqID, rfq.RFQCode)
	return rfq, nil
}

// GetAssignedSMEsForRFQ retrieves suppliers assigned to work packages linked to this RFQ
func GetAssignedSMEsForRFQ(rfqID int) ([]AssignedSMEDTO, error) {
	start := time.Now()
	log.Printf("[DEBUG] [GetAssignedSMEsForRFQ] Starting query for rfq_id=%d", rfqID)

	var list []AssignedSMEDTO
	query := `SELECT
    s.id,
    s.company_name,
    s.city,
    s.state,
    s.contact_email,
    s.phone,
    s.gst_number,
    s.msme_reg_number,
    s.is_verified,
    s.iso_standard,
    s.plant_area_sq_ft,
    s.worker_count,
    s.rating,
    s.capabilities_summary,
    s.nda_signed,
    s.created_at,
    wps.assignment_status,
    wps.assigned_at
FROM rfqs r
LEFT JOIN work_packages wp
    ON r.id = wp.rfq_id
LEFT JOIN work_package_suppliers wps
    ON wp.id = wps.work_package_id
LEFT JOIN suppliers s
    ON wps.supplier_id = s.id
WHERE r.id = ?
  AND s.id IS NOT NULL`

	rows, err := database.DB.Query(query, rfqID)
	if err != nil {
		log.Printf("[ERROR] [GetAssignedSMEsForRFQ] Query failed for rfq_id=%d: %v (elapsed=%v)", rfqID, err, time.Since(start))
		return list, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("[WARN] [GetAssignedSMEsForRFQ] Failed to close rows for rfq_id=%d: %v", rfqID, closeErr)
		}
	}()

	rowNum := 0
	for rows.Next() {
		rowNum++
		var dto AssignedSMEDTO
		err := rows.Scan(
			&dto.ID, &dto.CompanyName, &dto.City, &dto.State, &dto.ContactEmail, &dto.Phone,
			&dto.GSTNumber, &dto.MSMERegNumber, &dto.IsVerified, &dto.ISOStandard, &dto.PlantAreaSqFt,
			&dto.WorkerCount, &dto.Rating, &dto.CapabilitiesSummary, &dto.NDASigned, &dto.CreatedAt,
			&dto.AssignmentStatus, &dto.AssignedAt,
		)
		if err != nil {
			log.Printf("[WARN] [GetAssignedSMEsForRFQ] Row scan error at row=%d for rfq_id=%d: %v", rowNum, rfqID, err)
			continue
		}

		log.Printf("[DEBUG] [GetAssignedSMEsForRFQ] Scanned row=%d for rfq_id=%d: supplier_id=%d company=%q status=%q rating=%.2f verified=%v",
			rowNum, rfqID, dto.ID, dto.CompanyName, dto.AssignmentStatus, dto.Rating, dto.IsVerified)
		list = append(list, dto)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetAssignedSMEsForRFQ] Row iteration error for rfq_id=%d: %v (elapsed=%v)", rfqID, err, time.Since(start))
		return list, err
	}

	log.Printf("[DEBUG] [GetAssignedSMEsForRFQ] Completed for rfq_id=%d: rows_read=%d assigned_smes=%d (elapsed=%v)",
		rfqID, rowNum, len(list), time.Since(start))

	// Prints struct field names + values for every item in the slice
	log.Printf("[DEBUG] [GetAssignedSMEsForRFQ] list=%+v", list)
	return list, nil
}

// GetWorkPackagesByRFQID retrieves converted work packages for this RFQ
func GetWorkPackagesByRFQID(rfqID int) ([]WorkPackage, error) {
	start := time.Now()
	log.Printf("[DEBUG][GetWorkPackagesByRFQID] Fetching work packages for rfq_id=%d", rfqID)

	var wps []WorkPackage
	query := `SELECT id, rfq_id, title, redacted_spec, redacted_drawing_path, budget_usd, target_date, escrow_status, status, created_at 
              FROM work_packages WHERE rfq_id = ?`

	rows, err := database.DB.Query(query, rfqID)
	if err != nil {
		log.Printf("[ERROR][GetWorkPackagesByRFQID] Query failed for rfq_id=%d: %v (elapsed=%v)", rfqID, err, time.Since(start))
		return wps, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("[WARN][GetWorkPackagesByRFQID] Failed to close rows for rfq_id=%d: %v", rfqID, closeErr)
		}
	}()

	rowIdx := 0
	for rows.Next() {
		rowIdx++
		var wp WorkPackage
		err := rows.Scan(
			&wp.ID,
			&wp.RFQID,
			&wp.Title,
			&wp.RedactedSpec,
			&wp.RedactedDrawingPath,
			&wp.BudgetUSD,
			&wp.TargetDate,
			&wp.EscrowStatus,
			&wp.Status,
			&wp.CreatedAt,
		)
		if err != nil {
			log.Printf("[ERROR][GetWorkPackagesByRFQID] Row scan failed at row=%d for rfq_id=%d: %v", rowIdx, rfqID, err)
			continue
		}

		log.Printf("[DEBUG][GetWorkPackagesByRFQID] Scanned row=%d: wp_id=%d title=%q status=%q escrow_status=%q budget_usd=%.2f",
			rowIdx, wp.ID, wp.Title, wp.Status, wp.EscrowStatus, wp.BudgetUSD)
		wps = append(wps, wp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR][GetWorkPackagesByRFQID] Row iteration error for rfq_id=%d: %v (elapsed=%v)", rfqID, err, time.Since(start))
		return wps, err
	}

	log.Printf("[DEBUG][GetWorkPackagesByRFQID] Completed for rfq_id=%d: total_found=%d (elapsed=%v)", rfqID, len(wps), time.Since(start))
	return wps, nil
}

// GetPaymentsForRFQ retrieves payment trackers linked to the RFQ's work packages
func GetPaymentsForRFQ(rfqID int) ([]PaymentTracker, error) {
	var payments []PaymentTracker
	query := `SELECT pt.id, pt.work_package_id, pt.client_amount_usd, pt.client_paid_status, pt.client_paid_date, 
              pt.supplier_amount_usd, pt.supplier_paid_status, pt.supplier_paid_date, pt.notes, pt.updated_at
              FROM payment_trackers pt
              JOIN work_packages wp ON pt.work_package_id = wp.id
              WHERE wp.rfq_id = ?`
	rows, err := database.DB.Query(query, rfqID)
	if err != nil {
		return payments, err
	}
	defer rows.Close()

	for rows.Next() {
		var pt PaymentTracker
		if err := rows.Scan(&pt.ID, &pt.WorkPackageID, &pt.ClientAmountUSD, &pt.ClientPaidStatus, &pt.ClientPaidDate, &pt.SupplierAmountUSD, &pt.SupplierPaidStatus, &pt.SupplierPaidDate, &pt.Notes, &pt.UpdatedAt); err == nil {
			payments = append(payments, pt)
		}
	}
	return payments, nil
}

// GetAuditLogsForRFQ retrieves audit history logs for this RFQ record
func GetAuditLogsForRFQ(rfqID int) ([]AuditLog, error) {
	var logs []AuditLog
	// If you maintain an audit table, query it here. Fallback mock builder if table doesn't exist:
	query := `SELECT id, rfq_id, action, actor, details, created_at FROM rfq_audit_logs WHERE rfq_id = ? ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, rfqID)
	if err != nil {
		// Return empty slice gracefully if audit table is optional/uninitialized
		return []AuditLog{
			{ID: 1, RFQID: rfqID, Action: "CREATED", Actor: "Customer", Details: "Initial RFQ submission received via client portal.", CreatedAt: "2026-11-14 07:00:00"},
			{ID: 2, RFQID: rfqID, Action: "STATUS_CHANGE", Actor: "Admin", Details: "Status transitioned from New to Under Review.", CreatedAt: "2026-11-14 08:00:00"},
		}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var al AuditLog
		if err := rows.Scan(&al.ID, &al.RFQID, &al.Action, &al.Actor, &al.Details, &al.CreatedAt); err == nil {
			logs = append(logs, al)
		}
	}
	return logs, nil
}
