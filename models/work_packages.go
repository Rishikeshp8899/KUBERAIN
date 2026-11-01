/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : work_package.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing work packages, lifecycle stages,
                   comprehensive aggregated detail views, paginated dashboards, and RFQ conversions.
====================================================================================================
*/

package models

import (
	"fmt"
	"go-mvc-app/database"
	"log"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateWorkPackage                                                            ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new work package record with budgetary and escrow terms.           ║
║ 🏛️ ARCHITECTURE   : Parameterized execution binding RFQ references, supplier IDs, specs, and     ║
║                     escrow statuses, returning the generated SQLite primary key ID.              ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateWorkPackage(wp WorkPackage) (int64, error) {
	query := `INSERT INTO work_packages (rfq_id, supplier_id, title, redacted_spec, redacted_drawing_path, budget_usd, target_date, escrow_status, status) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, wp.RFQID, wp.SupplierID, wp.Title, wp.RedactedSpec, wp.RedactedDrawingPath, wp.BudgetUSD, wp.TargetDate, wp.EscrowStatus, wp.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetWorkPackageByID                                                           ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches a single work package record by its unique primary key ID.           ║
║ 🏛️ ARCHITECTURE   : Single-record point lookup query using COALESCE fallbacks for safe null      ║
║                     handling on optional supplier assignments and drawing paths.                 ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetWorkPackageByID(id int) (*WorkPackage, error) {
	query := `SELECT id, rfq_id, COALESCE(supplier_id, 0), title, redacted_spec, COALESCE(redacted_drawing_path, ''), budget_usd, COALESCE(target_date, ''), escrow_status, status, created_at 
              FROM work_packages WHERE id = ?`
	var wp WorkPackage
	err := database.DB.QueryRow(query, id).Scan(&wp.ID, &wp.RFQID, &wp.SupplierID, &wp.Title, &wp.RedactedSpec, &wp.RedactedDrawingPath, &wp.BudgetUSD, &wp.TargetDate, &wp.EscrowStatus, &wp.Status, &wp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &wp, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpdateWorkPackageStatus                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Updates the active lifecycle status string for a specified work package.     ║
║ 🏛️ ARCHITECTURE   : Targeted SQL update statement for rapid workflow state transitions without   ║
║                     reloading full entity graph records.                                         ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpdateWorkPackageStatus(id int, status string) error {
	_, err := database.DB.Exec(`UPDATE work_packages SET status = ? WHERE id = ?`, status, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteWorkPackage                                                            ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a work package record from the database.                 ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command targeting the primary key identifier.           ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteWorkPackage(id int) error {
	_, err := database.DB.Exec(`DELETE FROM work_packages WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetWorkPackages                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all active work packages sorted in descending order by ID.         ║
║ 🏛️ ARCHITECTURE   : Comprehensive listing query utilizing COALESCE defaults and descending sort    ║
║                     priorities for standard administrative dashboards.                             ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetWorkPackages() ([]WorkPackage, error) {
	query := `SELECT 
                id, 
                rfq_id, 
                COALESCE(supplier_id, 0), 
                COALESCE(title, ''), 
                COALESCE(redacted_spec, ''), 
                COALESCE(redacted_drawing_path, ''), 
                COALESCE(budget_usd, 0.0), 
                COALESCE(target_date, ''), 
                COALESCE(escrow_status, 'Unfunded'), 
                COALESCE(status, 'Bidding'), 
                COALESCE(created_at, '') 
              FROM work_packages 
              ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []WorkPackage
	for rows.Next() {
		var wp WorkPackage
		err := rows.Scan(
			&wp.ID,
			&wp.RFQID,
			&wp.SupplierID,
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
			return nil, err
		}
		packages = append(packages, wp)
	}

	return packages, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetWorkPackageFullDetails                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Aggregates complete related compliance records linked to a Work Package ID.  ║
║ 🏛️ ARCHITECTURE   : Relational aggregation wrapper joining supplier metadata, drawing revisions, ║
║                     material certificates, inspections, payments, and communication logs.        ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetWorkPackageFullDetails(wpID int) (WorkPackageDetailView, error) {
	var details WorkPackageDetailView

	// 1. Fetch Work Package & Assigned Supplier
	wpQuery := `SELECT wp.id, wp.rfq_id, wp.supplier_id, COALESCE(s.company_name, 'Unassigned SME'), 
                wp.title, wp.redacted_spec, wp.redacted_drawing_path, wp.budget_usd, 
                wp.target_date, wp.escrow_status, wp.status, wp.created_at 
                FROM work_packages wp 
                LEFT JOIN suppliers s ON wp.supplier_id = s.id 
                WHERE wp.id = ?`

	err := database.DB.QueryRow(wpQuery, wpID).Scan(
		&details.WorkPackage.ID, &details.WorkPackage.RFQID, &details.WorkPackage.SupplierID,
		&details.WorkPackage.SupplierName, &details.WorkPackage.Title, &details.WorkPackage.RedactedSpec,
		&details.WorkPackage.RedactedDrawingPath, &details.WorkPackage.BudgetUSD, &details.WorkPackage.TargetDate,
		&details.WorkPackage.EscrowStatus, &details.WorkPackage.Status, &details.WorkPackage.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] [GetWorkPackageFullDetails] WP fetch failed: %v", err)
		return details, err
	}

	// 2. Fetch Drawing Revisions (linked via RFQID)
	revRows, _ := database.DB.Query(`SELECT id, rfq_id, revision_code, original_name, redacted_path, change_log_notes, is_active, created_at 
                            FROM drawing_revisions WHERE rfq_id = ? ORDER BY created_at DESC`, details.WorkPackage.RFQID)
	if revRows != nil {
		defer revRows.Close()
		for revRows.Next() {
			var dr DrawingRevision
			revRows.Scan(&dr.ID, &dr.RFQID, &dr.RevisionCode, &dr.OriginalName, &dr.RedactedPath, &dr.ChangeLogNotes, &dr.IsActive, &dr.CreatedAt)
			details.DrawingRevisions = append(details.DrawingRevisions, dr)
		}
	}

	// 3. Fetch Material Certificates
	certRows, _ := database.DB.Query(`SELECT id, work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by, created_at 
                             FROM material_certificates WHERE work_package_id = ? ORDER BY created_at DESC`, wpID)
	if certRows != nil {
		defer certRows.Close()
		for certRows.Next() {
			var mc MaterialCertificate
			certRows.Scan(&mc.ID, &mc.WorkPackageID, &mc.HeatNumber, &mc.MaterialGrade, &mc.MillName, &mc.CertFilePath, &mc.VerifiedBy, &mc.CreatedAt)
			details.Certificates = append(details.Certificates, mc)
		}
	}

	// 4. Fetch Quality Inspections
	inspRows, _ := database.DB.Query(`SELECT id, work_package_id, inspection_type, passed, measured_values, tolerance_allowed, report_file_path, inspected_at 
                             FROM quality_inspections WHERE work_package_id = ? ORDER BY inspected_at DESC`, wpID)
	if inspRows != nil {
		defer inspRows.Close()
		for inspRows.Next() {
			var qi QualityInspection
			inspRows.Scan(&qi.ID, &qi.WorkPackageID, &qi.InspectionType, &qi.Passed, &qi.MeasuredValues, &qi.ToleranceAllowed, &qi.ReportFilePath, &qi.InspectedAt)
			details.Inspections = append(details.Inspections, qi)
		}
	}

	// 5. Fetch Payment Ledgers
	payRows, _ := database.DB.Query(`SELECT id, work_package_id, client_amount_usd, client_paid_status, client_paid_date, 
                             supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at 
                             FROM payment_trackers WHERE work_package_id = ? ORDER BY updated_at DESC`, wpID)
	if payRows != nil {
		defer payRows.Close()
		for payRows.Next() {
			var pt PaymentTracker
			payRows.Scan(&pt.ID, &pt.WorkPackageID, &pt.ClientAmountUSD, &pt.ClientPaidStatus, &pt.ClientPaidDate,
				&pt.SupplierAmountUSD, &pt.SupplierPaidStatus, &pt.SupplierPaidDate, &pt.Notes, &pt.UpdatedAt)
			details.Payments = append(details.Payments, pt)
		}
	}

	// 6. Fetch Communication Logs
	comRows, _ := database.DB.Query(`SELECT id, work_package_id, sender_role, sender_id, message, is_flagged, flag_reason, created_at 
                             FROM communication_logs WHERE work_package_id = ? ORDER BY created_at DESC`, wpID)
	if comRows != nil {
		defer comRows.Close()
		for comRows.Next() {
			var cl CommunicationLog
			comRows.Scan(&cl.ID, &cl.WorkPackageID, &cl.SenderRole, &cl.SenderID, &cl.Message, &cl.IsFlagged, &cl.FlagReason, &cl.CreatedAt)
			details.Communications = append(details.Communications, cl)
		}
	}

	return details, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetFilteredWorkPackages                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Queries work packages with joined supplier names, search, and pagination.    ║
║ 🏛️ ARCHITECTURE   : Advanced dashboard query builder combining dynamic status filters ("current" vs  ║
║                     "completed"), multi-table search joins, total count execution, and pagination. ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetFilteredWorkPackages(statusFilter string, search string, page int, pageSize int) (WorkPackagePaginatedResult, int, error) {

	var list []WorkPackage
	offset := (page - 1) * pageSize

	whereClause := "WHERE 1=1"
	var args []interface{}

	// Status filtering matching your filter tabs (all, current, completed)
	if statusFilter == "current" {
		whereClause += " AND wp.status NOT IN ('Delivered', 'Completed', 'Closed')"
	} else if statusFilter == "completed" {
		whereClause += " AND wp.status IN ('Delivered', 'Completed', 'Closed')"
	}

	// Search filter across task name (wp.title) and SME name (s.company_name)
	if search != "" {
		whereClause += " AND (wp.title LIKE ? OR s.company_name LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	log.Printf("[DEBUG] [GetFilteredWorkPackages] Input params -> statusFilter: %q, search: %q, page: %d, pageSize: %d", statusFilter, search, page, pageSize)

	// 1. Count Total Rows
	countQuery := `SELECT COUNT(*) FROM work_packages wp LEFT JOIN suppliers s ON wp.supplier_id = s.id ` + whereClause
	log.Printf("[DEBUG] [GetFilteredWorkPackages] Count Query: %s | Args: %v", countQuery, args)

	var totalRows int
	err := database.DB.QueryRow(countQuery, args...).Scan(&totalRows)
	if err != nil {
		log.Printf("[ERROR] [GetFilteredWorkPackages] Count query failed: %v", err)
		return WorkPackagePaginatedResult{}, 0, err
	}

	totalPages := (totalRows + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	log.Printf("[DEBUG] [GetFilteredWorkPackages] Count result -> totalRows: %d, totalPages: %d", totalRows, totalPages)

	// 2. Fetch Paginated Records
	dataQuery := `SELECT wp.id, wp.rfq_id, wp.supplier_id, COALESCE(s.company_name, 'Unassigned SME'), 
                  wp.title, wp.redacted_spec, wp.redacted_drawing_path, wp.budget_usd, 
                  wp.target_date, wp.escrow_status, wp.status, wp.created_at 
                  FROM work_packages wp 
                  LEFT JOIN suppliers s ON wp.supplier_id = s.id ` + whereClause + ` 
                  ORDER BY wp.created_at DESC LIMIT ? OFFSET ?`

	// Create a separate slice for data query arguments so we don't mutate `args` permanently
	dataArgs := append([]interface{}{}, args...)
	dataArgs = append(dataArgs, pageSize, offset)

	log.Printf("[DEBUG] [GetFilteredWorkPackages] Data Query: %s | Args: %v", dataQuery, dataArgs)

	rows, err := database.DB.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("[ERROR] [GetFilteredWorkPackages] Data query failed: %v", err)
		return WorkPackagePaginatedResult{}, 0, err
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		var wp WorkPackage
		if err := rows.Scan(
			&wp.ID, &wp.RFQID, &wp.SupplierID, &wp.SupplierName,
			&wp.Title, &wp.RedactedSpec, &wp.RedactedDrawingPath, &wp.BudgetUSD,
			&wp.TargetDate, &wp.EscrowStatus, &wp.Status, &wp.CreatedAt,
		); err != nil {
			log.Printf("[WARN] [GetFilteredWorkPackages] Row scan error at index %d: %v", rowCount, err)
			continue
		}
		list = append(list, wp)
		rowCount++
	}

	log.Printf("[DEBUG] [GetFilteredWorkPackages] Successfully scanned %d records for page %d", len(list), page)

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return WorkPackagePaginatedResult{
		WorkPackages: list,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PrevPage:     prevPage,
		NextPage:     nextPage,
		StatusFilter: statusFilter,
		Search:       search,
	}, totalPages, nil

}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetWorkPackageCompleteDetails                                                ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches all project history, vendor profiles, capabilities, and logs.        ║
║ 🏛️ ARCHITECTURE   : Deep multi-entity hydration wrapper gathering work packages, supplier profiles,  ║
║                     supplier capabilities, progress milestones, ECNs, MTCs, inspections, and ledgers.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetWorkPackageCompleteDetails(wpID int) (WorkPackageCompleteDetail, error) {
	var details WorkPackageCompleteDetail

	// 1. Fetch Work Package & Assigned Supplier
	wpQuery := `SELECT wp.id, wp.rfq_id, wp.supplier_id, COALESCE(s.company_name, 'Unassigned SME'), 
                wp.title, wp.redacted_spec, wp.redacted_drawing_path, wp.budget_usd, 
                wp.target_date, wp.escrow_status, wp.status, wp.created_at 
                FROM work_packages wp 
                LEFT JOIN suppliers s ON wp.supplier_id = s.id 
                WHERE wp.id = ?`

	err := database.DB.QueryRow(wpQuery, wpID).Scan(
		&details.WorkPackage.ID, &details.WorkPackage.RFQID, &details.WorkPackage.SupplierID,
		&details.WorkPackage.SupplierName, &details.WorkPackage.Title, &details.WorkPackage.RedactedSpec,
		&details.WorkPackage.RedactedDrawingPath, &details.WorkPackage.BudgetUSD, &details.WorkPackage.TargetDate,
		&details.WorkPackage.EscrowStatus, &details.WorkPackage.Status, &details.WorkPackage.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] [GetWorkPackageCompleteDetails] WP fetch failed: %v", err)
		return details, err
	}

	// 2. Fetch Supplier Full Profile & Capabilities
	if details.WorkPackage.SupplierID > 0 {
		database.DB.QueryRow(`SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
                     is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
                     FROM suppliers WHERE id = ?`, details.WorkPackage.SupplierID).Scan(
			&details.Supplier.ID, &details.Supplier.CompanyName, &details.Supplier.City, &details.Supplier.State,
			&details.Supplier.ContactEmail, &details.Supplier.Phone, &details.Supplier.GSTNumber, &details.Supplier.MSMERegNumber,
			&details.Supplier.IsVerified, &details.Supplier.ISOStandard, &details.Supplier.PlantAreaSqFt, &details.Supplier.WorkerCount,
			&details.Supplier.Rating, &details.Supplier.CapabilitiesSummary, &details.Supplier.NDASigned, &details.Supplier.CreatedAt,
		)

		database.DB.QueryRow(`SELECT id, supplier_id, category, machining_type, max_tolerance, max_component_dim, materials_handled, capacity_per_month 
                     FROM supplier_capabilities WHERE supplier_id = ?`, details.WorkPackage.SupplierID).Scan(
			&details.Capability.ID, &details.Capability.SupplierID, &details.Capability.Category, &details.Capability.MachiningType,
			&details.Capability.MaxTolerance, &details.Capability.MaxComponentDim, &details.Capability.MaterialsHandled, &details.Capability.CapacityPerMonth,
		)
	}

	// 3. Fetch Progress Milestones
	progRows, _ := database.DB.Query(`SELECT id, work_package_id, stage, percent_complete, milestone_note, inspector_name, updated_at 
                             FROM project_progress WHERE work_package_id = ? ORDER BY updated_at DESC`, wpID)
	if progRows != nil {
		defer progRows.Close()
		for progRows.Next() {
			var pp ProjectProgress
			progRows.Scan(&pp.ID, &pp.WorkPackageID, &pp.Stage, &pp.PercentComplete, &pp.MilestoneNote, &pp.InspectorName, &pp.UpdatedAt)
			details.Progress = append(details.Progress, pp)
		}
	}

	// 4. Fetch Drawing Revisions (ECN)
	revRows, _ := database.DB.Query(`SELECT id, rfq_id, revision_code, original_name, redacted_path, change_log_notes, is_active, created_at 
                            FROM drawing_revisions WHERE rfq_id = ? ORDER BY created_at DESC`, details.WorkPackage.RFQID)
	if revRows != nil {
		defer revRows.Close()
		for revRows.Next() {
			var dr DrawingRevision
			revRows.Scan(&dr.ID, &dr.RFQID, &dr.RevisionCode, &dr.OriginalName, &dr.RedactedPath, &dr.ChangeLogNotes, &dr.IsActive, &dr.CreatedAt)
			details.DrawingRevisions = append(details.DrawingRevisions, dr)
		}
	}

	// 5. Fetch Material Certificates (MTC)
	certRows, _ := database.DB.Query(`SELECT id, work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by, created_at 
                             FROM material_certificates WHERE work_package_id = ? ORDER BY created_at DESC`, wpID)
	if certRows != nil {
		defer certRows.Close()
		for certRows.Next() {
			var mc MaterialCertificate
			certRows.Scan(&mc.ID, &mc.WorkPackageID, &mc.HeatNumber, &mc.MaterialGrade, &mc.MillName, &mc.CertFilePath, &mc.VerifiedBy, &mc.CreatedAt)
			details.Certificates = append(details.Certificates, mc)
		}
	}

	// 6. Fetch Quality Inspections
	inspRows, _ := database.DB.Query(`SELECT id, work_package_id, inspection_type, passed, measured_values, tolerance_allowed, report_file_path, inspected_at 
                             FROM quality_inspections WHERE work_package_id = ? ORDER BY inspected_at DESC`, wpID)
	if inspRows != nil {
		defer inspRows.Close()
		for inspRows.Next() {
			var qi QualityInspection
			inspRows.Scan(&qi.ID, &qi.WorkPackageID, &qi.InspectionType, &qi.Passed, &qi.MeasuredValues, &qi.ToleranceAllowed, &qi.ReportFilePath, &qi.InspectedAt)
			details.Inspections = append(details.Inspections, qi)
		}
	}

	// 7. Fetch Payment Ledgers
	payRows, _ := database.DB.Query(`SELECT id, work_package_id, client_amount_usd, client_paid_status, client_paid_date, 
                             supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at 
                             FROM payment_trackers WHERE work_package_id = ? ORDER BY updated_at DESC`, wpID)
	if payRows != nil {
		defer payRows.Close()
		for payRows.Next() {
			var pt PaymentTracker
			payRows.Scan(&pt.ID, &pt.WorkPackageID, &pt.ClientAmountUSD, &pt.ClientPaidStatus, &pt.ClientPaidDate,
				&pt.SupplierAmountUSD, &pt.SupplierPaidStatus, &pt.SupplierPaidDate, &pt.Notes, &pt.UpdatedAt)
			details.Payments = append(details.Payments, pt)
		}
	}

	// 8. Fetch Communication Logs
	comRows, _ := database.DB.Query(`SELECT id, work_package_id, sender_role, sender_id, message, is_flagged, flag_reason, created_at 
                             FROM communication_logs WHERE work_package_id = ? ORDER BY created_at DESC`, wpID)
	if comRows != nil {
		defer comRows.Close()
		for comRows.Next() {
			var cl CommunicationLog
			comRows.Scan(&cl.ID, &cl.WorkPackageID, &cl.SenderRole, &cl.SenderID, &cl.Message, &cl.IsFlagged, &cl.FlagReason, &cl.CreatedAt)
			details.Communications = append(details.Communications, cl)
		}
	}

	return details, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllWorkPackagesForDropdown                                                ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Gets minimal work package details for selection dropdowns.                   ║
║ 🏛️ ARCHITECTURE   : Lightweight projection query returning only primary keys and task titles     ║
║                     ordered descending by ID to optimize UI dropdown selectors.                  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllWorkPackagesForDropdown() ([]WorkPackage, error) {
	log.Printf("[DEBUG] [GetAllWorkPackagesForDropdown] Executing SQLite query on 'work_packages' table...")
	query := `SELECT id, title FROM work_packages ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetAllWorkPackagesForDropdown] Database query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []WorkPackage
	scanCount := 0

	for rows.Next() {
		var wp WorkPackage
		if err := rows.Scan(&wp.ID, &wp.Title); err != nil {
			log.Printf("[ERROR] [GetAllWorkPackagesForDropdown] Row scan failed at index %d: %v", scanCount, err)
			continue
		}
		scanCount++
		log.Printf("[DEBUG] [GetAllWorkPackagesForDropdown] Scanned WP ID: %d, Title: %s", wp.ID, wp.Title)
		list = append(list, wp)
	}

	log.Printf("[SUCCESS] [GetAllWorkPackagesForDropdown] Successfully retrieved %d work packages from database", len(list))
	return list, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateWorkPackageFromRFQ                                                     ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Copies details from an approved RFQ into a new WorkPackage record.           ║
║ 🏛️ ARCHITECTURE   : Cross-table INSERT-SELECT pipeline mapping approved RFQ attributes directly  ║
║                     into active work package records with default escrow and active states.       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateWorkPackageFromRFQ(rfqID int) error {
	log.Printf("[DEBUG] [CreateWorkPackageFromRFQ] Executing insertion for RFQ ID=%d...", rfqID)

	// SQL query inserting RFQ data directly into work_packages
	query := `
        INSERT INTO work_packages (
            rfq_id, 
            supplier_id, 
            title, 
            redacted_spec, 
            redacted_drawing_path, 
            budget_usd, 
            target_date, 
            escrow_status, 
            status, 
            created_at
        )
        SELECT 
            id, 
            COALESCE(supplier_id, 0), 
            title, 
            COALESCE(specifications, ''), 
            COALESCE(drawing_path, ''), 
            budget_usd, 
            target_date, 
            'Pending', 
            'Active', 
            datetime('now')
        FROM rfqs 
        WHERE id = ?
    `

	res, err := database.DB.Exec(query, rfqID)
	if err != nil {
		log.Printf("[ERROR] [CreateWorkPackageFromRFQ] DB Exec failed for RFQ ID=%d: %v", rfqID, err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[WARN] [CreateWorkPackageFromRFQ] No RFQ found with ID=%d to convert", rfqID)
		return fmt.Errorf("no rfq found with id %d", rfqID)
	}

	log.Printf("[SUCCESS] [CreateWorkPackageFromRFQ] Created WorkPackage for RFQ ID=%d (Rows affected: %d)", rfqID, rowsAffected)
	return nil
}
