package models

import (
	"database/sql"
	"errors"
	"fmt"
	"go-mvc-app/database"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateEnterpriseClient(client EnterpriseClient) (int64, error) {
	query := `INSERT INTO enterprise_clients (name, domain, billing_city, tax_id, industry_sector, credit_limit_usd, account_status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, client.Name, client.Domain, client.BillingCity, client.TaxID, client.IndustrySector, client.CreditLimitUSD, client.AccountStatus)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetAllEnterpriseClients() ([]EnterpriseClient, error) {
	query := `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
	                 COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
	          FROM enterprise_clients ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []EnterpriseClient
	for rows.Next() {
		var c EnterpriseClient
		if err := rows.Scan(&c.ID, &c.Name, &c.Domain, &c.BillingCity, &c.TaxID, &c.IndustrySector, &c.CreditLimitUSD, &c.AccountStatus, &c.CreatedAt); err != nil {
			continue
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func GetEnterpriseClientByID(id int) (*EnterpriseClient, error) {
	query := `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
	                 COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
	          FROM enterprise_clients WHERE id = ?`
	var c EnterpriseClient
	err := database.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Domain, &c.BillingCity, &c.TaxID, &c.IndustrySector, &c.CreditLimitUSD, &c.AccountStatus, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateEnterpriseClientStatus(clientID int, status string) error {
	query := `UPDATE enterprise_clients SET account_status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, clientID)
	return err
}

func DeleteEnterpriseClient(clientID int) error {
	_, err := database.DB.Exec(`DELETE FROM enterprise_clients WHERE id = ?`, clientID)
	return err
}

func CreateClientUser(user ClientUser) (int64, error) {
	query := `INSERT INTO client_users (client_id, full_name, email, phone, role, department, is_active) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, user.ClientID, user.FullName, user.Email, user.Phone, user.Role, user.Department, user.IsActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

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

func DeleteClientUser(userID int) error {
	_, err := database.DB.Exec(`DELETE FROM client_users WHERE id = ?`, userID)
	return err
}

func CreateSupplier(s Supplier) (int64, error) {
	query := `INSERT INTO suppliers (company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, s.CompanyName, s.City, s.State, s.ContactEmail, s.Phone, s.GSTNumber, s.MSMERegNumber, s.IsVerified, s.ISOStandard, s.PlantAreaSqFt, s.WorkerCount, s.Rating, s.CapabilitiesSummary, s.NDASigned)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

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

func VerifySupplier(supplierID int) error {
	_, err := database.DB.Exec(`UPDATE suppliers SET is_verified = 1 WHERE id = ?`, supplierID)
	return err
}

func DeleteSupplier(supplierID int) error {
	_, err := database.DB.Exec(`DELETE FROM suppliers WHERE id = ?`, supplierID)
	return err
}

func CreateSupplierCapability(cap SupplierCapability) (int64, error) {
	query := `INSERT INTO supplier_capabilities (supplier_id, category, machining_type, max_tolerance, max_component_dim, materials_handled, capacity_per_month) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, cap.SupplierID, cap.Category, cap.MachiningType, cap.MaxTolerance, cap.MaxComponentDim, cap.MaterialsHandled, cap.CapacityPerMonth)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetCapabilitiesBySupplierID(supplierID int) ([]SupplierCapability, error) {
	query := `SELECT id, supplier_id, category, machining_type, COALESCE(max_tolerance, ''), COALESCE(max_component_dim, ''), COALESCE(materials_handled, ''), capacity_per_month 
	          FROM supplier_capabilities WHERE supplier_id = ?`
	rows, err := database.DB.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var caps []SupplierCapability
	for rows.Next() {
		var c SupplierCapability
		if err := rows.Scan(&c.ID, &c.SupplierID, &c.Category, &c.MachiningType, &c.MaxTolerance, &c.MaxComponentDim, &c.MaterialsHandled, &c.CapacityPerMonth); err != nil {
			continue
		}
		caps = append(caps, c)
	}
	return caps, nil
}

func DeleteCapability(id int) error {
	_, err := database.DB.Exec(`DELETE FROM supplier_capabilities WHERE id = ?`, id)
	return err
}

func CreateMachineStatus(m MachineStatus) (int64, error) {
	query := `INSERT INTO machine_statuses (supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, m.SupplierID, m.MachineName, m.IsOperational, m.LastCalibration, m.NextServiceDue, m.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetMachinesBySupplierID(supplierID int) ([]MachineStatus, error) {
	query := `SELECT id, supplier_id, machine_name, is_operational, COALESCE(last_calibration, ''), COALESCE(next_service_due, ''), COALESCE(notes, '') 
	          FROM machine_statuses WHERE supplier_id = ?`
	rows, err := database.DB.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []MachineStatus
	for rows.Next() {
		var m MachineStatus
		if err := rows.Scan(&m.ID, &m.SupplierID, &m.MachineName, &m.IsOperational, &m.LastCalibration, &m.NextServiceDue, &m.Notes); err != nil {
			continue
		}
		machines = append(machines, m)
	}
	return machines, nil
}

func UpdateMachineOperationalStatus(machineID int, isOperational bool) error {
	_, err := database.DB.Exec(`UPDATE machine_statuses SET is_operational = ? WHERE id = ?`, isOperational, machineID)
	return err
}

func DeleteMachineStatus(id int) error {
	_, err := database.DB.Exec(`DELETE FROM machine_statuses WHERE id = ?`, id)
	return err
}

func CreateInquiry(inq Inquiry) (int64, error) {
	query := `INSERT INTO inquiries (name, email, phone, subject, message, source, ip_address) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, inq.Name, inq.Email, inq.Phone, inq.Subject, inq.Message, inq.Source, inq.IPAddress)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

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

func DeleteInquiry(id int) error {
	_, err := database.DB.Exec(`DELETE FROM inquiries WHERE id = ?`, id)
	return err
}

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

func UpdateRFQStatus(id int, status string) error {
	_, err := database.DB.Exec(`UPDATE rfqs SET status = ? WHERE id = ?`, status, id)
	return err
}

func DeleteRFQ(id int) error {
	_, err := database.DB.Exec(`DELETE FROM rfqs WHERE id = ?`, id)
	return err
}

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

func DeleteDrawingRevision(id int) error {
	_, err := database.DB.Exec(`DELETE FROM drawing_revisions WHERE id = ?`, id)
	return err
}

func CreateWorkPackage(wp WorkPackage) (int64, error) {
	query := `INSERT INTO work_packages (rfq_id, supplier_id, title, redacted_spec, redacted_drawing_path, budget_usd, target_date, escrow_status, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, wp.RFQID, wp.SupplierID, wp.Title, wp.RedactedSpec, wp.RedactedDrawingPath, wp.BudgetUSD, wp.TargetDate, wp.EscrowStatus, wp.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

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

func UpdateWorkPackageStatus(id int, status string) error {
	_, err := database.DB.Exec(`UPDATE work_packages SET status = ? WHERE id = ?`, status, id)
	return err
}

func DeleteWorkPackage(id int) error {
	_, err := database.DB.Exec(`DELETE FROM work_packages WHERE id = ?`, id)
	return err
}

func CreateProjectProgress(p ProjectProgress) (int64, error) {
	query := `INSERT INTO project_progress (work_package_id, stage, percent_complete, milestone_note, photos_path, inspection_doc, inspector_name) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, p.WorkPackageID, p.Stage, p.PercentComplete, p.MilestoneNote, p.PhotosPath, p.InspectionDoc, p.InspectorName)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetProgressByWorkPackageID(workPackageID int) ([]ProjectProgress, error) {
	query := `SELECT id, work_package_id, stage, percent_complete, COALESCE(milestone_note, ''), COALESCE(photos_path, ''), COALESCE(inspection_doc, ''), COALESCE(inspector_name, ''), updated_at 
	          FROM project_progress WHERE work_package_id = ? ORDER BY updated_at DESC`
	rows, err := database.DB.Query(query, workPackageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ProjectProgress
	for rows.Next() {
		var p ProjectProgress
		if err := rows.Scan(&p.ID, &p.WorkPackageID, &p.Stage, &p.PercentComplete, &p.MilestoneNote, &p.PhotosPath, &p.InspectionDoc, &p.InspectorName, &p.UpdatedAt); err != nil {
			continue
		}
		list = append(list, p)
	}
	return list, nil
}

func DeleteProjectProgress(id int) error {
	_, err := database.DB.Exec(`DELETE FROM project_progress WHERE id = ?`, id)
	return err
}

func CreateMaterialCertificate(cert MaterialCertificate) (int64, error) {
	query := `INSERT INTO material_certificates (work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, cert.WorkPackageID, cert.HeatNumber, cert.MaterialGrade, cert.MillName, cert.CertFilePath, cert.VerifiedBy)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetCertificatesByWorkPackageID(packageID int) ([]MaterialCertificate, error) {
	query := `SELECT id, work_package_id, heat_number, material_grade, COALESCE(mill_name, ''), cert_file_path, COALESCE(verified_by, ''), created_at 
	          FROM material_certificates WHERE work_package_id = ?`
	rows, err := database.DB.Query(query, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []MaterialCertificate
	for rows.Next() {
		var c MaterialCertificate
		if err := rows.Scan(&c.ID, &c.WorkPackageID, &c.HeatNumber, &c.MaterialGrade, &c.MillName, &c.CertFilePath, &c.VerifiedBy, &c.CreatedAt); err != nil {
			continue
		}
		certs = append(certs, c)
	}
	return certs, nil
}

func DeleteMaterialCertificate(id int) error {
	_, err := database.DB.Exec(`DELETE FROM material_certificates WHERE id = ?`, id)
	return err
}

func CreateQualityInspection(qi QualityInspection) (int64, error) {
	query := `INSERT INTO quality_inspections (work_package_id, inspection_type, passed, measured_values, tolerance_allowed, report_file_path) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, qi.WorkPackageID, qi.InspectionType, qi.Passed, qi.MeasuredValues, qi.ToleranceAllowed, qi.ReportFilePath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetInspectionsByWorkPackageID(packageID int) ([]QualityInspection, error) {
	query := `SELECT id, work_package_id, inspection_type, passed, COALESCE(measured_values, ''), COALESCE(tolerance_allowed, ''), COALESCE(report_file_path, ''), inspected_at 
	          FROM quality_inspections WHERE work_package_id = ?`
	rows, err := database.DB.Query(query, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []QualityInspection
	for rows.Next() {
		var q QualityInspection
		if err := rows.Scan(&q.ID, &q.WorkPackageID, &q.InspectionType, &q.Passed, &q.MeasuredValues, &q.ToleranceAllowed, &q.ReportFilePath, &q.InspectedAt); err != nil {
			continue
		}
		inspections = append(inspections, q)
	}
	return inspections, nil
}

func DeleteQualityInspection(id int) error {
	_, err := database.DB.Exec(`DELETE FROM quality_inspections WHERE id = ?`, id)
	return err
}

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

func DeletePaymentTracker(id int) error {
	_, err := database.DB.Exec(`DELETE FROM payment_trackers WHERE id = ?`, id)
	return err
}

func CreateCommunicationLog(logData CommunicationLog) (int64, error) {
	query := `INSERT INTO communication_logs (work_package_id, sender_role, sender_id, message, is_flagged, flag_reason) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, logData.WorkPackageID, logData.SenderRole, logData.SenderID, logData.Message, logData.IsFlagged, logData.FlagReason)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// func GetFlaggedCommunicationLogs() ([]CommunicationLog, error) {
// 	query := `SELECT id, work_package_id, sender_role, COALESCE(sender_id, 0), message, is_flagged, COALESCE(flag_reason, ''), created_at
// 	          FROM communication_logs WHERE is_flagged = 1 ORDER BY id DESC`
// 	rows, err := database.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logsList []CommunicationLog
// 	for rows.Next() {
// 		var l CommunicationLog
// 		if err := rows.Scan(&l.ID, &l.WorkPackageID, &l.SenderRole, &l.SenderID, &l.Message, &l.IsFlagged, &l.FlagReason, &l.CreatedAt); err != nil {
// 			continue
// 		}
// 		logsList = append(logsList, l)
// 	}
// 	return logsList, nil
// }

func DeleteCommunicationLog(id int) error {
	_, err := database.DB.Exec(`DELETE FROM communication_logs WHERE id = ?`, id)
	return err
}

func CreateLeakDetectionRule(rule LeakDetectionRule) (int64, error) {
	query := `INSERT INTO leak_detection_rules (rule_name, regex_pattern, severity, is_active) 
	          VALUES (?, ?, ?, ?)`
	res, err := database.DB.Exec(query, rule.RuleName, rule.RegexPattern, rule.Severity, rule.IsActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetActiveLeakRules() ([]LeakDetectionRule, error) {
	query := `SELECT id, rule_name, regex_pattern, severity, is_active FROM leak_detection_rules WHERE is_active = 1`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []LeakDetectionRule
	for rows.Next() {
		var r LeakDetectionRule
		if err := rows.Scan(&r.ID, &r.RuleName, &r.RegexPattern, &r.Severity, &r.IsActive); err != nil {
			continue
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func DeleteLeakDetectionRule(id int) error {
	_, err := database.DB.Exec(`DELETE FROM leak_detection_rules WHERE id = ?`, id)
	return err
}

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

// CreateSMESupplier inserts a new MSME vendor profile into SQLite
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

func UpdateSMEVerificationStatus(id int64, isVerified bool) error {
	query := `UPDATE sme_suppliers SET is_verified = ? WHERE id = ?`
	_, err := database.DB.Exec(query, isVerified, id)
	return err
}

// UpdatePaymentStatus updates both client and supplier payment states
func UpdatePaymentStatus(id int, clientStatus, supplierStatus string) error {
	query := `UPDATE payment_ledger 
			  SET client_paid_status = ?, supplier_paid_status = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	_, err := database.DB.Exec(query, clientStatus, supplierStatus, id)
	return err
}

// GetWorkPackages retrieves all active work packages from SQLite
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

// GetMaterialCertificates fetches all Material Test Certificates (MTC) uploaded for quality compliance
func GetMaterialCertificates() ([]MaterialCertificate, error) {
	query := `SELECT 
				id, 
				work_package_id, 
				COALESCE(heat_number, ''), 
				COALESCE(material_grade, ''), 
				COALESCE(mill_name, ''), 
				COALESCE(cert_file_path, ''), 
				COALESCE(verified_by, ''), 
				COALESCE(created_at, '') 
			  FROM material_certificates 
			  ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []MaterialCertificate
	for rows.Next() {
		var cert MaterialCertificate
		err := rows.Scan(
			&cert.ID,
			&cert.WorkPackageID,
			&cert.HeatNumber,
			&cert.MaterialGrade,
			&cert.MillName,
			&cert.CertFilePath,
			&cert.VerifiedBy,
			&cert.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

// AuthenticateAdmin verifies admin email and hashed password
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

// GetRFQByID fetches a single RFQ record by its primary key ID
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

// GetCapableSuppliersForRFQ filters SMEs based on matching categories or material capabilities
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

// GetAllSuppliers retrieves all onboarded SME suppliers ordered by rating and verification state
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

// GetAllMachineStatuses retrieves all shop-floor machine health and calibration records
func GetAllMachineStatuses() ([]MachineStatus, error) {
	var machines []MachineStatus

	query := `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
	          FROM machine_statuses 
	          ORDER BY is_operational DESC, next_service_due ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetAllMachineStatuses] Query execution failed: %v", err)
		return machines, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MachineStatus
		err := rows.Scan(
			&ms.ID, &ms.SupplierID, &ms.MachineName, &ms.IsOperational,
			&ms.LastCalibration, &ms.NextServiceDue, &ms.Notes,
		)
		if err != nil {
			log.Printf("[WARN] [GetAllMachineStatuses] Row scan error: %v", err)
			continue
		}
		machines = append(machines, ms)
	}

	return machines, nil
}

// ToggleSupplierVerification flips the verification status of an SME vendor
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

// GetPublicVerifiedSuppliers retrieves all verified SME suppliers for public display
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

// GetPublicMachineStatuses retrieves operational machine health records for public audit display
func GetPublicMachineStatuses() ([]MachineStatus, error) {
	var machines []MachineStatus

	query := `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
	          FROM machine_statuses 
	          WHERE is_operational = 1`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetPublicMachineStatuses] Query execution failed: %v", err)
		return machines, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MachineStatus
		err := rows.Scan(
			&ms.ID, &ms.SupplierID, &ms.MachineName, &ms.IsOperational,
			&ms.LastCalibration, &ms.NextServiceDue, &ms.Notes,
		)
		if err != nil {
			log.Printf("[WARN] [GetPublicMachineStatuses] Row scan error: %v", err)
			continue
		}
		machines = append(machines, ms)
	}

	return machines, nil
}

// RegisterNewSMESupplier handles inserting supplier data, capability, and multiple machine units atomically
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

// GetPaginatedSuppliers fetches suppliers filtered by status ("all", "approved", "pending") and search term with pagination
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

// GetWorkPackageFullDetails fetches everything associated with a Work Package ID
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

// GetFilteredWorkPackages queries work packages with joined supplier names, search filtering (task name or SME name), and pagination
func GetFilteredWorkPackages(statusFilter string, search string, page int, pageSize int) (WorkPackagePaginatedResult, error) {

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
		return WorkPackagePaginatedResult{}, err
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
		return WorkPackagePaginatedResult{}, err
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
	}, nil

}

// GetWorkPackageCompleteDetails fetches all project history, vendor info, progress, and communications
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

// AddSecureMessage checks for contact info leaks (phone/email) and logs the communication
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

// GetPaymentLedger fetches all double-entry records joined with work package titles
func GetPaymentLedger() ([]PaymentTracker, error) {
	query := `SELECT pt.id, pt.work_package_id, COALESCE(wp.title, 'General Project'),
	            pt.client_amount_usd, pt.client_paid_status, COALESCE(pt.client_paid_date, ''), 
	            pt.supplier_amount_usd, pt.supplier_paid_status, COALESCE(pt.supplier_paid_date, ''), 
	            COALESCE(pt.notes, ''), COALESCE(pt.updated_at, '') 
	          FROM payment_trackers pt
	          LEFT JOIN work_packages wp ON pt.work_package_id = wp.id
	          ORDER BY pt.id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentTracker
	for rows.Next() {
		var p PaymentTracker
		err := rows.Scan(
			&p.ID, &p.WorkPackageID,
			&p.ClientAmountUSD, &p.ClientPaidStatus, &p.ClientPaidDate,
			&p.SupplierAmountUSD, &p.SupplierPaidStatus, &p.SupplierPaidDate,
			&p.Notes, &p.UpdatedAt,
		)
		if err != nil {
			continue
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// GetPaymentsByWorkPackageID fetches payment records filtered by a specific work package
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

// InsertPaymentRecord inserts a new double-entry payment ledger record
func InsertPaymentRecord(wpID int, clientAmt float64, clientStatus, clientDate string, supplierAmt float64, supplierStatus, supplierDate, notes string) error {
	query := `INSERT INTO payment_trackers (work_package_id, client_amount_usd, client_paid_status, client_paid_date, supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`
	_, err := database.DB.Exec(query, wpID, clientAmt, clientStatus, clientDate, supplierAmt, supplierStatus, supplierDate, notes)
	return err
}

// GetAllSuppliersForDropdown fetches minimal supplier details for dropdown selection
func GetAllSuppliersForDropdown() ([]Supplier, error) {
	query := `SELECT id, company_name, city, state FROM suppliers ORDER BY company_name ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var s Supplier
		err := rows.Scan(&s.ID, &s.CompanyName, &s.City, &s.State)
		if err != nil {
			continue
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}

// GetAllWorkPackagesForDropdown gets minimal work package details for selection dropdowns
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

// CreateWorkPackageFromRFQ copies details from an approved RFQ into a new WorkPackage record
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

// GetQualityInspections fetches per-SME quality inspections joined with supplier names
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

// GetMTCCertificates fetches Material Test Certificates joined with supplier names
func GetMTCCertificates() ([]MaterialCertificate, error) {
	query := `
		SELECT 
			mtc.id, mtc.work_package_id, mtc.supplier_id, 
			COALESCE(s.name, 'Unknown SME') as supplier_name, 
			mtc.heat_number, mtc.material_grade, mtc.mill_name, 
			mtc.cert_file_path, mtc.verified_by, mtc.created_at
		FROM mtc_certificates mtc
		LEFT JOIN suppliers s ON mtc.supplier_id = s.id
		ORDER BY mtc.created_at DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetMTCCertificates] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var mtcs []MaterialCertificate
	for rows.Next() {
		var m MaterialCertificate
		if err := rows.Scan(&m.ID, &m.WorkPackageID, &m.SupplierID, &m.SupplierName, &m.HeatNumber, &m.MaterialGrade, &m.MillName, &m.CertFilePath, &m.VerifiedBy, &m.CreatedAt); err != nil {
			continue
		}
		mtcs = append(mtcs, m)
	}
	return mtcs, nil
}

// GetFlaggedCommunicationLogs fetches disintermediation / leak detection rules logs
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
