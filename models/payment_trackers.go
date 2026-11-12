/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : payment_tracker.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing financial ledgers, double-entry payment
                   tracking, escrow states, upsert synchronizations, and joined quality inspection records.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
	"log"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPaymentsByWorkPackageID                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches payment ledger records filtered by a specific work package ID.       ║
║ 🏛️ ARCHITECTURE   : Relational query joining payment trackers with work package titles, using    ║
║                     COALESCE defaults and descending ID sorting for targeted financial audits.   ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPaymentsByWorkPackageID(wpID int) ([]PaymentTracker, error) {
	log.Printf("[DEBUG] [GetPaymentsByWorkPackageID] Fetching payment records for Work Package ID: %d", wpID)

	query := `SELECT pt.id, pt.work_package_id, COALESCE(wp.title, 'General Project'),
                pt.client_amount_usd, pt.client_paid_status, COALESCE(pt.client_paid_date, ''), 
                pt.supplier_amount_usd, pt.supplier_paid_status, COALESCE(pt.supplier_paid_date, ''), 
                COALESCE(pt.notes, ''), COALESCE(pt.updated_at, '') 
              FROM payment_trackers pt
              LEFT JOIN work_packages wp ON pt.work_package_id = wp.id
              WHERE pt.work_package_id = ?
              ORDER BY pt.id DESC`

	rows, err := database.DB.Query(query, wpID)
	if err != nil {
		log.Printf("[ERROR] [GetPaymentsByWorkPackageID] Database query failed for wpID %d: %v", wpID, err)
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentTracker
	rowCounter := 0

	for rows.Next() {
		var p PaymentTracker
		err := rows.Scan(
			&p.ID,
			&p.WorkPackageID,
			&p.ClientAmountUSD,
			&p.ClientPaidStatus,
			&p.ClientPaidDate,
			&p.SupplierAmountUSD,
			&p.SupplierPaidStatus,
			&p.SupplierPaidDate,
			&p.Notes,
			&p.UpdatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] [GetPaymentsByWorkPackageID] Row scan failed at index %d for wpID %d: %v", rowCounter, wpID, err)
			continue
		}

		rowCounter++
		log.Printf("[DEBUG] [GetPaymentsByWorkPackageID] Scanned payment ID %d (Client USD: %.2f, Supplier USD: %.2f)", p.ID, p.ClientAmountUSD, p.SupplierAmountUSD)
		payments = append(payments, p)
	}

	// Check for iteration errors from rows.Err()
	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetPaymentsByWorkPackageID] Row iteration error for wpID %d: %v", wpID, err)
		return nil, err
	}

	log.Printf("[SUCCESS] [GetPaymentsByWorkPackageID] Successfully fetched %d payment records for Work Package ID %d", len(payments), wpID)
	return payments, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : InsertPaymentRecord                                                          ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new double-entry payment ledger record with live timestamps.       ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding client/supplier amounts and statuses, ║
║                     auto-initializing updated_at values using SQLite datetime('now') functions.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func InsertPaymentRecord(wpID int, clientAmt float64, clientStatus, clientDate string, supplierAmt float64, supplierStatus, supplierDate, notes string) error {
	query := `INSERT INTO payment_trackers (work_package_id, client_amount_usd, client_paid_status, client_paid_date, supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`
	_, err := database.DB.Exec(query, wpID, clientAmt, clientStatus, clientDate, supplierAmt, supplierStatus, supplierDate, notes)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpsertPaymentTracker                                                         ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts or updates payment states atomically based on work package conflicts.║
║ 🏛️ ARCHITECTURE   : Utilizes SQLite ON CONFLICT(work_package_id) DO UPDATE clauses to seamlessly    ║
║                     synchronize payment statuses without duplicate row insertions.               ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpsertPaymentTracker(p PaymentTracker) error {
	query := `INSERT INTO payment_trackers (work_package_id, client_amount_usd, client_paid_status, client_paid_date, supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)
              ON CONFLICT(work_package_id) DO UPDATE SET 
              client_paid_status=excluded.client_paid_status, 
              client_paid_date=excluded.client_paid_date,
              supplier_paid_status=excluded.supplier_paid_status,
              supplier_paid_date=excluded.supplier_paid_date,
              notes=excluded.notes,
              updated_at=CURRENT_TIMESTAMP`
	_, err := database.DB.Exec(query, p.WorkPackageID, p.ClientAmountUSD, p.ClientPaidStatus, p.ClientPaidDate, p.SupplierAmountUSD, p.SupplierPaidStatus, p.SupplierPaidDate, p.Notes)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPaymentTrackerByWorkPackageID                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches a single payment tracker record by its associated work package ID.     ║
║ 🏛️ ARCHITECTURE   : Targeted point-lookup query using COALESCE wrappers for safe null-handling   ║
║                     on dates and audit logs.                                                     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPaymentTrackerByWorkPackageID(packageID int) (*PaymentTracker, error) {
	query := `SELECT id, work_package_id, client_amount_usd, client_paid_status, COALESCE(client_paid_date, ''), supplier_amount_usd, supplier_paid_status, COALESCE(supplier_paid_date, ''), COALESCE(notes, ''), updated_at 
              FROM payment_trackers WHERE work_package_id = ?`
	var p PaymentTracker
	err := database.DB.QueryRow(query, packageID).Scan(&p.ID, &p.WorkPackageID, &p.ClientAmountUSD, &p.ClientPaidStatus, &p.ClientPaidDate, &p.SupplierAmountUSD, &p.SupplierPaidStatus, &p.SupplierPaidDate, &p.Notes, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeletePaymentTracker                                                         ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a payment tracker entry from the database.               ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the payment tracker primary key. ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeletePaymentTracker(id int) error {
	_, err := database.DB.Exec(`DELETE FROM payment_trackers WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPaymentLedger                                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches payment ledger records filtered by status and project title search.   ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPaymentLedger(statusArgs ...string) ([]PaymentTrackerLog, error) {
	status := "current"
	if len(statusArgs) > 0 && statusArgs[0] != "" {
		status = statusArgs[0]
	}

	log.Printf("[DEBUG] [GetPaymentLedger] Fetching payment ledger records with status filter: %s", status)

	query := `
        SELECT 
            pt.id, 
            pt.work_package_id, 
            COALESCE(wp.title, 'Untitled Package') as work_package_title,
            pt.client_amount_usd, 
            pt.client_paid_status, 
            COALESCE(pt.client_paid_date, ''), 
            pt.supplier_amount_usd, 
            pt.supplier_paid_status, 
            COALESCE(pt.supplier_paid_date, ''), 
            COALESCE(pt.notes, ''), 
            COALESCE(pt.updated_at, '') 
        FROM payment_trackers pt
        LEFT JOIN work_packages wp ON pt.work_package_id = wp.id
        ORDER BY pt.updated_at DESC
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetPaymentLedger] Database query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentTrackerLog
	for rows.Next() {
		var p PaymentTrackerLog
		err := rows.Scan(
			&p.ID, &p.WorkPackageID, &p.WorkPackageTitle,
			&p.ClientAmountUSD, &p.ClientPaidStatus, &p.ClientPaidDate,
			&p.SupplierAmountUSD, &p.SupplierPaidStatus, &p.SupplierPaidDate,
			&p.Notes, &p.UpdatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] [GetPaymentLedger] Row scan failed: %v", err)
			continue
		}
		payments = append(payments, p)
	}

	log.Printf("[SUCCESS] [GetPaymentLedger] Successfully fetched %d payment ledger records.", len(payments))
	return payments, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetQualityInspections                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches per-SME quality inspections joined with corresponding supplier names.║
║ 🏛️ ARCHITECTURE   : Relational query joining quality inspections and suppliers, featuring COALESCE ║
║                     fallbacks for safety and chronological sorting by inspection timestamps.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetQualityInspections() ([]QualityInspection, error) {
	query := `
        SELECT 
            qi.id, qi.work_package_id, qi.supplier_id, 
            COALESCE(s.name, 'Unknown SME') as supplier_name, 
            qi.inspection_type, qi.passed, COALESCE(qi.comments, ''), 
            COALESCE(qi.measured_values, ''), COALESCE(qi.tolerance_allowed, ''), 
            COALESCE(qi.report_file_path, ''), qi.inspected_at
        FROM quality_inspections qi
        LEFT JOIN suppliers s ON qi.supplier_id = s.id
        ORDER BY qi.inspected_at DESC
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetQualityInspections] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var inspections []QualityInspection
	for rows.Next() {
		var qi QualityInspection
		if err := rows.Scan(&qi.ID, &qi.WorkPackageID, &qi.SupplierID, &qi.SupplierName, &qi.InspectionType, &qi.Passed, &qi.Comments, &qi.MeasuredValues, &qi.ToleranceAllowed, &qi.ReportFilePath, &qi.InspectedAt); err != nil {
			log.Printf("[ERROR] [GetQualityInspections] Row scan failed: %v", err)
			continue
		}
		inspections = append(inspections, qi)
	}
	return inspections, nil
}

// UpdatePaymentStatus updates both client and supplier payment states
func UpdatePaymentStatus(id int, clientStatus, supplierStatus string) error {
	query := `UPDATE payment_ledger 
			  SET client_paid_status = ?, supplier_paid_status = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	_, err := database.DB.Exec(query, clientStatus, supplierStatus, id)
	return err
}
