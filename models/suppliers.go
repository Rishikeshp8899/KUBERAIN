/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : supplier.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing supplier/SME vendor profiles, database
                   queries, atomic transactions, verification workflows, and capability matching.
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
║ 🛠️ FUNCTION       : CreateSupplier                                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a single supplier record into the primary `suppliers` table.           ║
║ 🏛️ ARCHITECTURE   : Lightweight single-table persistence mechanism used when basic vendor stubs  ║
║                     are created independently of complex machinery or capability parameters.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateSupplier(s Supplier) (int64, error) {
	query := `INSERT INTO suppliers (company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, s.CompanyName, s.City, s.State, s.ContactEmail, s.Phone, s.GSTNumber, s.MSMERegNumber, s.IsVerified, s.ISOStandard, s.PlantAreaSqFt, s.WorkerCount, s.Rating, s.CapabilitiesSummary, s.NDASigned)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetSupplierByID                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches a single supplier profile record by its primary key ID.              ║
║ 🏛️ ARCHITECTURE   : Focuses entirely on single-record point lookups using COALESCE wrappers for  ║
║                     null-safety, bypassing list iteration or filtering overhead.                 ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetSupplierByID(id int) (*Supplier, error) {
	query := `SELECT id, company_name, city, state, contact_email, COALESCE(phone, ''), COALESCE(gst_number, ''), COALESCE(msme_reg_number, ''), is_verified, COALESCE(iso_standard, ''), plant_area_sq_ft, worker_count, rating, COALESCE(capabilities_summary, ''), nda_signed, created_at 
              FROM suppliers WHERE id = ?`
	var s Supplier
	err := database.DB.QueryRow(query, id).Scan(&s.ID, &s.CompanyName, &s.City, &s.State, &s.ContactEmail, &s.Phone, &s.GSTNumber, &s.MSMERegNumber, &s.IsVerified, &s.ISOStandard, &s.PlantAreaSqFt, &s.WorkerCount, &s.Rating, &s.CapabilitiesSummary, &s.NDASigned, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : VerifySupplier                                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Grants official verification to a supplier (sets is_verified = 1).           ║
║ 🏛️ ARCHITECTURE   : Executes a targeted SQL update statement for rapid administrative status     ║
║                     override without reloading full entity payloads.                             ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func VerifySupplier(supplierID int) error {
	_, err := database.DB.Exec(`UPDATE suppliers SET is_verified = 1 WHERE id = ?`, supplierID)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteSupplier                                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently removes a supplier profile entry from the database.              ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion execution by primary key ID, contrasting with soft      ║
║                     status update wrappers.                                                      ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteSupplier(supplierID int) error {
	_, err := database.DB.Exec(`DELETE FROM suppliers WHERE id = ?`, supplierID)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetSMESuppliers                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves legacy/secondary MSME vendor listings from the sme_suppliers table.║
║ 🏛️ ARCHITECTURE   : Isolates secondary MSME records from the primary unified supplier directory  ║
║                     to maintain complete backward compatibility with distinct table structures.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetSMESuppliers() ([]Supplier, error) {
	query := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
              FROM sme_suppliers ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var s Supplier
		err := rows.Scan(
			&s.ID,
			&s.CompanyName,
			&s.City,
			&s.State,
			&s.ContactEmail,
			&s.Phone,
			&s.GSTNumber,
			&s.MSMERegNumber,
			&s.IsVerified,
			&s.ISOStandard,
			&s.PlantAreaSqFt,
			&s.WorkerCount,
			&s.Rating,
			&s.CapabilitiesSummary,
			&s.NDASigned,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	return suppliers, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateSMESupplier                                                            ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new MSME vendor profile specifically into SQLite sme_suppliers.    ║
║ 🏛️ ARCHITECTURE   : Specialized insertion procedure ensuring compliance attributes and plant     ║
║                     metrics are bound to secondary MSME data structures.                         ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateSMESupplier(s Supplier) error {
	query := `INSERT INTO sme_suppliers (company_name, contact_email, phone, city, state, gst_number, msme_reg_number, iso_standard, is_verified, plant_area_sq_ft, worker_count, capabilities_summary, nda_signed) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(
		query,
		s.CompanyName,
		s.ContactEmail,
		s.Phone,
		s.City,
		s.State,
		s.GSTNumber,
		s.MSMERegNumber,
		s.ISOStandard,
		s.IsVerified,
		s.PlantAreaSqFt,
		s.WorkerCount,
		s.CapabilitiesSummary,
		s.NDASigned,
	)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpdateSMEVerificationStatus                                                  ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Updates the verification standing of an MSME supplier record.                ║
║ 🏛️ ARCHITECTURE   : Parameterized update operation targeting the secondary sme_suppliers table   ║
║                     using explicit boolean flags.                                                ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpdateSMEVerificationStatus(id int64, isVerified bool) error {
	query := `UPDATE sme_suppliers SET is_verified = ? WHERE id = ?`
	_, err := database.DB.Exec(query, isVerified, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetCapableSuppliersForRFQ                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Intelligent matching engine querying suppliers based on category or materials.║
║ 🏛️ ARCHITECTURE   : Performs relational joins (LEFT JOIN supplier_capabilities) with wildcard  ║
║                     pattern matching, featuring an automatic fallback to top-rated general vendors.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetCapableSuppliersForRFQ(category string, material string) ([]Supplier, error) {
	var suppliers []Supplier

	// Query suppliers whose capabilities or verification match the RFQ needs
	query := `SELECT DISTINCT s.id, s.company_name, s.city, s.state, s.contact_email, s.phone, 
              s.gst_number, s.msme_reg_number, s.is_verified, s.iso_standard, s.plant_area_sq_ft, 
              s.worker_count, s.rating, s.capabilities_summary, s.nda_signed, s.created_at
              FROM suppliers s
              LEFT JOIN supplier_capabilities sc ON s.id = sc.supplier_id
              WHERE s.is_verified = 1 AND (sc.category LIKE ? OR sc.materials_handled LIKE ? OR s.capabilities_summary LIKE ?)`

	searchTerm := "%" + category + "%"
	materialTerm := "%" + material + "%"

	rows, err := database.DB.Query(query, searchTerm, materialTerm, searchTerm)
	if err != nil {
		log.Printf("[ERROR] [GetCapableSuppliersForRFQ] Query execution failed: %v", err)
		return suppliers, err
	}
	defer rows.Close()

	for rows.Next() {
		var sup Supplier
		err := rows.Scan(
			&sup.ID, &sup.CompanyName, &sup.City, &sup.State, &sup.ContactEmail, &sup.Phone,
			&sup.GSTNumber, &sup.MSMERegNumber, &sup.IsVerified, &sup.ISOStandard, &sup.PlantAreaSqFt,
			&sup.WorkerCount, &sup.Rating, &sup.CapabilitiesSummary, &sup.NDASigned, &sup.CreatedAt,
		)
		if err != nil {
			log.Printf("[WARN] [GetCapableSuppliersForRFQ] Row scan error: %v", err)
			continue
		}
		suppliers = append(suppliers, sup)
	}

	// Fallback: If no capability match found, return top verified suppliers as general candidates
	if len(suppliers) == 0 {
		log.Printf("[INFO] [GetCapableSuppliersForRFQ] No strict capability match found. Fetching top verified general suppliers.")
		fallbackQuery := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
                          is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
                          FROM suppliers WHERE is_verified = 1 LIMIT 5`
		fRows, fErr := database.DB.Query(fallbackQuery)
		if fErr == nil {
			defer fRows.Close()
			for fRows.Next() {
				var sup Supplier
				_ = fRows.Scan(
					&sup.ID, &sup.CompanyName, &sup.City, &sup.State, &sup.ContactEmail, &sup.Phone,
					&sup.GSTNumber, &sup.MSMERegNumber, &sup.IsVerified, &sup.ISOStandard, &sup.PlantAreaSqFt,
					&sup.WorkerCount, &sup.Rating, &sup.CapabilitiesSummary, &sup.NDASigned, &sup.CreatedAt,
				)
				suppliers = append(suppliers, sup)
			}
		}
	}

	return suppliers, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllSuppliers                                                              ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all onboarded SME suppliers sorted by priority and rating.         ║
║ 🏛️ ARCHITECTURE   : Employs multi-column sorting prioritization (ORDER BY is_verified DESC,      ║
║                     rating DESC) to bubble up trusted verified vendors in internal views.        ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllSuppliers() ([]Supplier, error) {
	var suppliers []Supplier

	query := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
              is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
              nda_signed, created_at 
              FROM suppliers 
              ORDER BY is_verified DESC, rating DESC, created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetAllSuppliers] Query execution failed: %v", err)
		return suppliers, err
	}
	defer rows.Close()

	for rows.Next() {
		var sup Supplier
		err := rows.Scan(
			&sup.ID, &sup.CompanyName, &sup.City, &sup.State, &sup.ContactEmail, &sup.Phone,
			&sup.GSTNumber, &sup.MSMERegNumber, &sup.IsVerified, &sup.ISOStandard, &sup.PlantAreaSqFt,
			&sup.WorkerCount, &sup.Rating, &sup.CapabilitiesSummary, &sup.NDASigned, &sup.CreatedAt,
		)
		if err != nil {
			log.Printf("[WARN] [GetAllSuppliers] Row scan error: %v", err)
			continue
		}
		suppliers = append(suppliers, sup)
	}

	return suppliers, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : ToggleSupplierVerification                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Dynamically flips a supplier's verification state between 1 and 0.           ║
║ 🏛️ ARCHITECTURE   : Utilizes conditional SQL update expressions (CASE WHEN is_verified = 1 ...)  ║
║                     to toggle states atomically without requiring prior memory-load roundtrips.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func ToggleSupplierVerification(supplierID int) error {
	query := `UPDATE suppliers SET is_verified = CASE WHEN is_verified = 1 THEN 0 ELSE 1 END WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		log.Printf("[ERROR] [ToggleSupplierVerification] Statement preparation failed: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(supplierID)
	if err != nil {
		log.Printf("[ERROR] [ToggleSupplierVerification] Execution failed for Supplier ID %d: %v", supplierID, err)
		return err
	}

	log.Printf("[SUCCESS] [ToggleSupplierVerification] Toggled verification for Supplier ID %d", supplierID)
	return nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPublicVerifiedSuppliers                                                   ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves verified SME suppliers specifically for public-facing portals.     ║
║ 🏛️ ARCHITECTURE   : Restricts query results strictly to verified vendors (WHERE is_verified = 1) ║
║                     sorted by highest ratings, contrasting with unfiltered admin listings.       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPublicVerifiedSuppliers() ([]Supplier, error) {
	var suppliers []Supplier

	query := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
              is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
              nda_signed, created_at 
              FROM suppliers 
              WHERE is_verified = 1 
              ORDER BY rating DESC, created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetPublicVerifiedSuppliers] Query execution failed: %v", err)
		return suppliers, err
	}
	defer rows.Close()

	for rows.Next() {
		var sup Supplier
		err := rows.Scan(
			&sup.ID, &sup.CompanyName, &sup.City, &sup.State, &sup.ContactEmail, &sup.Phone,
			&sup.GSTNumber, &sup.MSMERegNumber, &sup.IsVerified, &sup.ISOStandard, &sup.PlantAreaSqFt,
			&sup.WorkerCount, &sup.Rating, &sup.CapabilitiesSummary, &sup.NDASigned, &sup.CreatedAt,
		)
		if err != nil {
			log.Printf("[WARN] [GetPublicVerifiedSuppliers] Row scan error: %v", err)
			continue
		}
		suppliers = append(suppliers, sup)
	}

	return suppliers, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : RegisterNewSMESupplier                                                       ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Handles multi-entity vendor onboarding (Profile + Capability + Machines)     ║
║                     atomically within a single transaction.                                      ║
║ 🏛️ ARCHITECTURE   : The most comprehensive ACID insertion pipeline in the system. Wraps multiple ║
║                     dependent table inserts in database transactions (Begin, Commit, Rollback)   ║
║                     to prevent partial database states upon failure.                             ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func RegisterNewSMESupplier(sup Supplier, cap SupplierCapability, machines []MachineStatus) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("[ERROR] [RegisterNewSMESupplier] Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// 1. Insert Supplier
	supplierQuery := `INSERT INTO suppliers (company_name, city, state, contact_email, phone, gst_number, 
                      msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, 
                      rating, capabilities_summary, nda_signed) 
                      VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, 4.5, ?, ?)`

	res, err := tx.Exec(supplierQuery, sup.CompanyName, sup.City, sup.State, sup.ContactEmail, sup.Phone,
		sup.GSTNumber, sup.MSMERegNumber, sup.ISOStandard, sup.PlantAreaSqFt, sup.WorkerCount,
		sup.CapabilitiesSummary, sup.NDASigned)
	if err != nil {
		log.Printf("[ERROR] [RegisterNewSMESupplier] Supplier insertion failed: %v", err)
		return err
	}

	supplierID64, err := res.LastInsertId()
	if err != nil {
		log.Printf("[ERROR] [RegisterNewSMESupplier] Failed to retrieve last insert ID: %v", err)
		return err
	}
	supplierID := int(supplierID64)

	// 2. Insert Supplier Capability
	capQuery := `INSERT INTO supplier_capabilities (supplier_id, category, machining_type, max_tolerance, 
                 max_component_dim, materials_handled, capacity_per_month) 
                 VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(capQuery, supplierID, cap.Category, cap.MachiningType, cap.MaxTolerance,
		cap.MaxComponentDim, cap.MaterialsHandled, cap.CapacityPerMonth)
	if err != nil {
		log.Printf("[ERROR] [RegisterNewSMESupplier] Capability insertion failed: %v", err)
		return err
	}

	// 3. Insert Multiple Machine Statuses
	macQuery := `INSERT INTO machine_statuses (supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes) 
                 VALUES (?, ?, 1, ?, ?, ?)`

	for _, mac := range machines {
		// Skip empty entries if any rogue row exists
		if mac.MachineName == "" {
			continue
		}
		_, err = tx.Exec(macQuery, supplierID, mac.MachineName, mac.LastCalibration, mac.NextServiceDue, mac.Notes)
		if err != nil {
			log.Printf("[ERROR] [RegisterNewSMESupplier] Machine status insertion failed for '%s': %v", mac.MachineName, err)
			return err
		}
	}

	// Commit Transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR] [RegisterNewSMESupplier] Transaction commit failed: %v", err)
		return err
	}

	log.Printf("[SUCCESS] [RegisterNewSMESupplier] Successfully registered SME supplier ID %d with %d machines", supplierID, len(machines))
	return nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPaginatedSuppliers                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Advanced dashboard query builder supporting dynamic filtering, multi-column  ║
║                     text searches, row counting, and pagination offsets.                         ║
║ 🏛️ ARCHITECTURE   : Combines dynamic WHERE clause generation with separate count and limit       ║
║                     queries, delivering comprehensive pagination metadata structures.            ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPaginatedSuppliers(filter string, search string, page int, pageSize int) (PaginatedResult, error) {
	var suppliers []Supplier
	offset := (page - 1) * pageSize

	// Base query clauses
	whereClause := "WHERE 1=1"
	var args []interface{}

	if filter == "approved" {
		whereClause += " AND is_verified = 1"
	} else if filter == "pending" {
		whereClause += " AND is_verified = 0"
	}

	if search != "" {
		whereClause += " AND (company_name LIKE ? OR city LIKE ? OR state LIKE ? OR gst_number LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// 1. Get Total Count for Pagination Metadata
	countQuery := "SELECT COUNT(*) FROM suppliers " + whereClause
	var totalRows int
	err := database.DB.QueryRow(countQuery, args...).Scan(&totalRows)
	if err != nil {
		log.Printf("[ERROR] [GetPaginatedSuppliers] Count query failed: %v", err)
		return PaginatedResult{}, err
	}

	totalPages := (totalRows + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// 2. Fetch Paginated Records
	dataQuery := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
                  is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
                  nda_signed, created_at 
                  FROM suppliers ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)

	rows, err := database.DB.Query(dataQuery, args...)
	if err != nil {
		log.Printf("[ERROR] [GetPaginatedSuppliers] Data query failed: %v", err)
		return PaginatedResult{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var sup Supplier
		err := rows.Scan(
			&sup.ID, &sup.CompanyName, &sup.City, &sup.State, &sup.ContactEmail, &sup.Phone,
			&sup.GSTNumber, &sup.MSMERegNumber, &sup.IsVerified, &sup.ISOStandard, &sup.PlantAreaSqFt,
			&sup.WorkerCount, &sup.Rating, &sup.CapabilitiesSummary, &sup.NDASigned, &sup.CreatedAt,
		)
		if err != nil {
			log.Printf("[WARN] [GetPaginatedSuppliers] Row scan error: %v", err)
			continue
		}
		suppliers = append(suppliers, sup)
	}

	return PaginatedResult{
		Suppliers:   suppliers,
		TotalPages:  totalPages,
		CurrentPage: page,
		Filter:      filter,
		Search:      search,
	}, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllSuppliersPaginatedWithFilter                                           ║
║ 👨‍💻 DEVELOPER      : Rishikesh Manohar Patil (KUBERAIN INDUSTRY)                                  ║
║ 🎯 WHAT IT DOES   : Queries supplier records with dynamic filtering, search, and pagination.     ║
║ 🏛️ ARCHITECTURE   : Parameterized pagination query execution with integrated error handling.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/

func GetAllSuppliersPaginatedWithFilter(page int, filter string, search string) ([]Supplier, error) {
	limit := 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at FROM suppliers WHERE 1=1`
	var args []interface{}

	// Apply Filter Tab Logic
	if filter == "approved" {
		query += ` AND is_verified = 1`
	} else if filter == "pending" {
		query += ` AND is_verified = 0`
	}

	// Apply Search Query Logic
	if search != "" {
		query += ` AND (company_name LIKE ? OR city LIKE ? OR gst_number LIKE ?)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	log.Printf("[INFO] Executing paginated & filtered query: %s [Page: %d, Filter: %s, Search: %s]", query, page, filter, search)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("[ERROR] Failed to execute supplier query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var s Supplier
		err := rows.Scan(
			&s.ID, &s.CompanyName, &s.City, &s.State, &s.ContactEmail, &s.Phone,
			&s.GSTNumber, &s.MSMERegNumber, &s.IsVerified, &s.ISOStandard,
			&s.PlantAreaSqFt, &s.WorkerCount, &s.Rating, &s.CapabilitiesSummary,
			&s.NDASigned, &s.CreatedAt,
		)
		if err != nil {
			continue
		}
		suppliers = append(suppliers, s)
	}

	log.Printf("[INFO] Successfully retrieved %d supplier records.", len(suppliers))
	return suppliers, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : SupplierDetailsPage                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Manohar Patil (KUBERAIN INDUSTRY)                                  ║
║ 🎯 WHAT IT DOES   : Queries relational tables for supplier profile, capabilities, and machine    ║
║                     fleet inventory by supplier ID, checking errors and returning ViewModel.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func SupplierDetailsPage(supplierID int) (SupplierDetails, error) {
	if supplierID <= 0 {
		log.Printf("[WARN] [SupplierDetailsPage] Invalid or missing supplier ID: %d", supplierID)
		return SupplierDetails{}, fmt.Errorf("invalid supplier ID: %d", supplierID)
	}

	// 1. Fetch Supplier Profile
	var s Supplier
	querySupplier := `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
                      FROM suppliers WHERE id = ?`

	err := database.DB.QueryRow(querySupplier, supplierID).Scan(
		&s.ID, &s.CompanyName, &s.City, &s.State, &s.ContactEmail, &s.Phone,
		&s.GSTNumber, &s.MSMERegNumber, &s.IsVerified, &s.ISOStandard,
		&s.PlantAreaSqFt, &s.WorkerCount, &s.Rating, &s.CapabilitiesSummary,
		&s.NDASigned, &s.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] [SupplierDetailsPage] Supplier not found for ID %d: %v", supplierID, err)
		return SupplierDetails{}, err
	}

	// 2. Fetch Supplier Capabilities
	var cap SupplierCapability
	queryCap := `SELECT id, supplier_id, category, machining_type, max_tolerance, max_component_dim, materials_handled, capacity_per_month 
                 FROM supplier_capabilities WHERE supplier_id = ?`
	_ = database.DB.QueryRow(queryCap, supplierID).Scan(
		&cap.ID, &cap.SupplierID, &cap.Category, &cap.MachiningType,
		&cap.MaxTolerance, &cap.MaxComponentDim, &cap.MaterialsHandled, &cap.CapacityPerMonth,
	)

	// 3. Fetch Machine Fleet
	queryMachines := `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
                      FROM machine_statuses WHERE supplier_id = ?`
	rows, err := database.DB.Query(queryMachines, supplierID)
	var machines []MachineStatus
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var m MachineStatus
			if scanErr := rows.Scan(&m.ID, &m.SupplierID, &m.MachineName, &m.IsOperational, &m.LastCalibration, &m.NextServiceDue, &m.Notes); scanErr == nil {
				machines = append(machines, m)
			}
		}
	}

	log.Printf("[INFO] [SupplierDetailsPage] Successfully loaded details for supplier ID: %d", supplierID)

	return SupplierDetails{
		Title:      s.CompanyName + " - Audit Profile",
		Supplier:   s,
		Capability: cap,
		Machines:   machines,
	}, nil
}
