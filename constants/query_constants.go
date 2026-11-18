package constants

// ====================================================================================================
// 📁 admin_user.go
// ====================================================================================================
const (
	QueryAuthenticateAdmin = `SELECT id, email, password_hash, COALESCE(role, 'admin') FROM admin_users WHERE email = ?`
)

// ====================================================================================================
// 📁 client_user.go
// ====================================================================================================
const (
	QueryCreateClientUser = `INSERT INTO client_users (client_id, full_name, email, phone, role, department, is_active) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryGetUsersByClientID = `SELECT id, client_id, full_name, email, COALESCE(phone, ''), role, COALESCE(department, ''), is_active, COALESCE(last_login_at, ''), created_at 
              FROM client_users WHERE client_id = ?`

	QueryDeleteClientUser = `DELETE FROM client_users WHERE id = ?`
)

// ====================================================================================================
// 📁 communication_log.go
// ====================================================================================================
const (
	QueryCreateCommunicationLog = `INSERT INTO communication_logs (work_package_id, sender_role, sender_id, message, is_flagged, flag_reason) 
              VALUES (?, ?, ?, ?, ?, ?)`

	QueryAddSecureMessage = `INSERT INTO communication_logs (work_package_id, sender_role, sender_id, message, is_flagged, flag_reason, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`

	QueryDeleteCommunicationLog = `DELETE FROM communication_logs WHERE id = ?`

	QueryGetFlaggedCommunicationLogs = `SELECT id, work_package_id, sender_role, sender_id, flag_reason, message, created_at FROM communication_logs ORDER BY created_at DESC`
)

// ====================================================================================================
// 📁 drawing_revision.go
// ====================================================================================================
const (
	QueryDeactivateDrawingRevisionsByRFQ = `UPDATE drawing_revisions SET is_active = 0 WHERE rfq_id = ?`

	QueryCreateDrawingRevision = `INSERT INTO drawing_revisions (rfq_id, revision_code, original_name, redacted_path, change_log_notes, is_active) 
              VALUES (?, ?, ?, ?, ?, 1)`

	QueryGetActiveDrawingByRFQ = `SELECT id, rfq_id, revision_code, original_name, redacted_path, COALESCE(change_log_notes, ''), is_active, created_at 
              FROM drawing_revisions WHERE rfq_id = ? AND is_active = 1`

	QueryDeleteDrawingRevision = `DELETE FROM drawing_revisions WHERE id = ?`
)

// ====================================================================================================
// 📁 enterprise_client.go
// ====================================================================================================
const (
	QueryCreateEnterpriseClient = `INSERT INTO enterprise_clients (name, domain, billing_city, tax_id, industry_sector, credit_limit_usd, account_status) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryGetAllEnterpriseClients = `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
                     COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
              FROM enterprise_clients ORDER BY id DESC`

	QueryGetEnterpriseClientByID = `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
                     COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
              FROM enterprise_clients WHERE id = ?`

	QueryUpdateEnterpriseClientStatus = `UPDATE enterprise_clients SET account_status = ? WHERE id = ?`

	QueryDeleteEnterpriseClient = `DELETE FROM enterprise_clients WHERE id = ?`
)

// ====================================================================================================
// 📁 inquiry.go
// ====================================================================================================
const (
	QueryCreateInquiry = `INSERT INTO inquiries (name, email, phone, subject, message, source, ip_address) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryGetAllInquiries = `SELECT id, name, email, COALESCE(phone, ''), COALESCE(subject, ''), message, COALESCE(source, 'Contact Form'), COALESCE(ip_address, ''), created_at 
              FROM inquiries ORDER BY id DESC`

	QueryDeleteInquiry = `DELETE FROM inquiries WHERE id = ?`
)

// ====================================================================================================
// 📁 leak_detection_rule.go
// ====================================================================================================
const (
	QueryCreateLeakDetectionRule = `INSERT INTO leak_detection_rules (rule_name, regex_pattern, severity, is_active) 
              VALUES (?, ?, ?, ?)`

	QueryGetActiveLeakRules = `SELECT id, rule_name, regex_pattern, severity, is_active FROM leak_detection_rules WHERE is_active = 1`

	QueryDeleteLeakDetectionRule = `DELETE FROM leak_detection_rules WHERE id = ?`
)

// ====================================================================================================
// 📁 machine_status.go
// ====================================================================================================
const (
	QueryCreateMachineStatus = `INSERT INTO machine_statuses (supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes) 
              VALUES (?, ?, ?, ?, ?, ?)`

	QueryGetMachinesBySupplierID = `SELECT id, supplier_id, machine_name, is_operational, COALESCE(last_calibration, ''), COALESCE(next_service_due, ''), COALESCE(notes, '') 
              FROM machine_statuses WHERE supplier_id = ?`

	QueryUpdateMachineOperationalStatus = `UPDATE machine_statuses SET is_operational = ? WHERE id = ?`

	QueryDeleteMachineStatus = `DELETE FROM machine_statuses WHERE id = ?`

	QueryGetPublicMachineStatuses = `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
              FROM machine_statuses 
              WHERE is_operational = 1`

	QueryGetAllMachineStatuses = `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
              FROM machine_statuses 
              ORDER BY is_operational DESC, next_service_due ASC`
)

// ====================================================================================================
// 📁 material_certificate.go
// ====================================================================================================
const (
	QueryCreateMaterialCertificate = `INSERT INTO material_certificates (work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by) 
              VALUES (?, ?, ?, ?, ?, ?)`

	QueryGetCertificatesByWorkPackageID = `SELECT id, work_package_id, heat_number, material_grade, COALESCE(mill_name, ''), cert_file_path, COALESCE(verified_by, ''), created_at 
              FROM material_certificates WHERE work_package_id = ?`

	QueryDeleteMaterialCertificate = `DELETE FROM material_certificates WHERE id = ?`

	QueryGetMaterialCertificates = `SELECT 
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
)

// ====================================================================================================
// 📁 payment_tracker.go
// ====================================================================================================
const (
	QueryGetPaymentsByWorkPackageID = `SELECT pt.id, pt.work_package_id, COALESCE(wp.title, 'General Project'),
                pt.client_amount_usd, pt.client_paid_status, COALESCE(pt.client_paid_date, ''), 
                pt.supplier_amount_usd, pt.supplier_paid_status, COALESCE(pt.supplier_paid_date, ''), 
                COALESCE(pt.notes, ''), COALESCE(pt.updated_at, '') 
              FROM payment_trackers pt
              LEFT JOIN work_packages wp ON pt.work_package_id = wp.id
              WHERE pt.work_package_id = ?
              ORDER BY pt.id DESC`

	QueryInsertPaymentRecord = `INSERT INTO payment_trackers (work_package_id, client_amount_usd, client_paid_status, client_paid_date, supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`

	QueryUpsertPaymentTracker = `INSERT INTO payment_trackers (work_package_id, client_amount_usd, client_paid_status, client_paid_date, supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)
              ON CONFLICT(work_package_id) DO UPDATE SET 
              client_paid_status=excluded.client_paid_status, 
              client_paid_date=excluded.client_paid_date,
              supplier_paid_status=excluded.supplier_paid_status,
              supplier_paid_date=excluded.supplier_paid_date,
              notes=excluded.notes,
              updated_at=CURRENT_TIMESTAMP`

	QueryGetPaymentTrackerByWorkPackageID = `SELECT id, work_package_id, client_amount_usd, client_paid_status, COALESCE(client_paid_date, ''), supplier_amount_usd, supplier_paid_status, COALESCE(supplier_paid_date, ''), COALESCE(notes, ''), updated_at 
              FROM payment_trackers WHERE work_package_id = ?`

	QueryDeletePaymentTracker = `DELETE FROM payment_trackers WHERE id = ?`

	QueryGetPaymentLedger = `
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

	QueryGetQualityInspections = `
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

	QueryUpdatePaymentStatus = `UPDATE payment_ledger 
              SET client_paid_status = ?, supplier_paid_status = ?, updated_at = CURRENT_TIMESTAMP 
              WHERE id = ?`
)

// ====================================================================================================
// 📁 project_progress.go
// ====================================================================================================
const (
	QueryCreateProjectProgress = `INSERT INTO project_progress (work_package_id, stage, percent_complete, milestone_note, photos_path, inspection_doc, inspector_name) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryGetProgressByWorkPackageID = `SELECT id, work_package_id, stage, percent_complete, COALESCE(milestone_note, ''), COALESCE(photos_path, ''), COALESCE(inspection_doc, ''), COALESCE(inspector_name, ''), updated_at 
              FROM project_progress WHERE work_package_id = ? ORDER BY updated_at DESC`

	QueryDeleteProjectProgress = `DELETE FROM project_progress WHERE id = ?`
)

// ====================================================================================================
// 📁 rfq.go
// ====================================================================================================
const (
	QueryCreateRFQ = `INSERT INTO rfqs (rfq_code, company_name, contact_person, email, phone, category, material, quantity, requirements, drawing_file_path, target_price_usd, delivery_deadline, incobooking_terms, country, city, ip_address, status) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	QueryGetAllRFQs = `SELECT id, COALESCE(rfq_code, 'N/A'), company_name, contact_person, email, COALESCE(phone, ''), COALESCE(category, ''), COALESCE(material, ''), COALESCE(quantity, ''), COALESCE(requirements, ''), COALESCE(drawing_file_path, ''), target_price_usd, COALESCE(delivery_deadline, ''), COALESCE(incobooking_terms, ''), COALESCE(country, 'India'), COALESCE(city, 'Unknown'), COALESCE(ip_address, 'N/A'), COALESCE(status, 'New'), created_at 
              FROM rfqs ORDER BY id DESC`

	QueryUpdateRFQStatus = `UPDATE rfqs SET status = ? WHERE id = ?`

	QueryDeleteRFQ = `DELETE FROM rfqs WHERE id = ?`

	QueryGetRFQByID = `SELECT id, rfq_code, company_name, contact_person, email, phone, category, material, quantity, 
              requirements, drawing_file_path, target_price_usd, delivery_deadline, 
              incobooking_terms, country, city, ip_address, status, created_at 
              FROM rfqs WHERE id = ?`

	QueryGetAssignedSMEsForRFQ = `SELECT s.id, s.company_name, s.city, s.state, s.contact_email, s.phone, 
              s.gst_number, s.msme_reg_number, s.is_verified, s.iso_standard, s.plant_area_sq_ft, 
              s.worker_count, s.rating, s.capabilities_summary, s.nda_signed, s.created_at,
              wps.assignment_status, wps.assigned_at
              FROM work_package_suppliers wps
              JOIN work_packages wp ON wps.work_package_id = wp.id
              JOIN suppliers s ON wps.supplier_id = s.id
              WHERE wp.rfq_id = ?`

	QueryGetWorkPackagesByRFQID = `SELECT id, rfq_id, title, redacted_spec, redacted_drawing_path, budget_usd, target_date, escrow_status, status, created_at 
              FROM work_packages WHERE rfq_id = ?`

	QueryGetPaymentsForRFQ = `SELECT pt.id, pt.work_package_id, pt.client_amount_usd, pt.client_paid_status, pt.client_paid_date, 
              pt.supplier_amount_usd, pt.supplier_paid_status, pt.supplier_paid_date, pt.notes, pt.updated_at
              FROM payment_trackers pt
              JOIN work_packages wp ON pt.work_package_id = wp.id
              WHERE wp.rfq_id = ?`

	QueryGetAuditLogsForRFQ = `SELECT id, rfq_id, action, actor, details, created_at FROM rfq_audit_logs WHERE rfq_id = ? ORDER BY created_at DESC`
)

// ====================================================================================================
// 📁 supplier_capability.go
// ====================================================================================================
const (
	QueryCreateSupplierCapability = `INSERT INTO supplier_capabilities (supplier_id, category, machining_type, max_tolerance, max_component_dim, materials_handled, capacity_per_month) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryGetCapabilitiesBySupplierID = `SELECT id, supplier_id, category, machining_type, COALESCE(max_tolerance, ''), COALESCE(max_component_dim, ''), COALESCE(materials_handled, ''), capacity_per_month 
              FROM supplier_capabilities WHERE supplier_id = ?`

	QueryDeleteCapability = `DELETE FROM supplier_capabilities WHERE id = ?`
)

// ====================================================================================================
// 📁 supplier.go
// ====================================================================================================
const (
	QueryCreateSupplier = `INSERT INTO suppliers (company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	QueryGetSupplierByID = `SELECT id, company_name, city, state, contact_email, COALESCE(phone, ''), COALESCE(gst_number, ''), COALESCE(msme_reg_number, ''), is_verified, COALESCE(iso_standard, ''), plant_area_sq_ft, worker_count, rating, COALESCE(capabilities_summary, ''), nda_signed, created_at 
              FROM suppliers WHERE id = ?`

	QueryVerifySupplier = `UPDATE suppliers SET is_verified = 1 WHERE id = ?`

	QueryDeleteSupplier = `DELETE FROM suppliers WHERE id = ?`

	QueryGetSMESuppliers = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
              FROM sme_suppliers ORDER BY id DESC`

	QueryCreateSMESupplier = `INSERT INTO sme_suppliers (company_name, contact_email, phone, city, state, gst_number, msme_reg_number, iso_standard, is_verified, plant_area_sq_ft, worker_count, capabilities_summary, nda_signed) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	QueryUpdateSMEVerificationStatus = `UPDATE sme_suppliers SET is_verified = ? WHERE id = ?`

	QueryGetCapableSuppliersForRFQ = `SELECT DISTINCT s.id, s.company_name, s.city, s.state, s.contact_email, s.phone, 
              s.gst_number, s.msme_reg_number, s.is_verified, s.iso_standard, s.plant_area_sq_ft, 
              s.worker_count, s.rating, s.capabilities_summary, s.nda_signed, s.created_at
              FROM suppliers s
              LEFT JOIN supplier_capabilities sc ON s.id = sc.supplier_id
              WHERE s.is_verified = 1 AND (sc.category LIKE ? OR sc.materials_handled LIKE ? OR s.capabilities_summary LIKE ?)`

	QueryGetCapableSuppliersFallback = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
                          is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
                          FROM suppliers WHERE is_verified = 1 LIMIT 5`

	QueryGetAllSuppliers = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
              is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
              nda_signed, created_at 
              FROM suppliers 
              ORDER BY is_verified DESC, rating DESC, created_at DESC`

	QueryToggleSupplierVerification = `UPDATE suppliers SET is_verified = CASE WHEN is_verified = 1 THEN 0 ELSE 1 END WHERE id = ?`

	QueryGetPublicVerifiedSuppliers = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
              is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
              nda_signed, created_at 
              FROM suppliers 
              WHERE is_verified = 1 
              ORDER BY rating DESC, created_at DESC`

	QueryRegisterNewSMESupplier = `INSERT INTO suppliers (company_name, city, state, contact_email, phone, gst_number, 
                      msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, 
                      rating, capabilities_summary, nda_signed) 
                      VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, 4.5, ?, ?)`

	QueryRegisterNewSMECapability = `INSERT INTO supplier_capabilities (supplier_id, category, machining_type, max_tolerance, 
                 max_component_dim, materials_handled, capacity_per_month) 
                 VALUES (?, ?, ?, ?, ?, ?, ?)`

	QueryRegisterNewSMEMachine = `INSERT INTO machine_statuses (supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes) 
                 VALUES (?, ?, 1, ?, ?, ?)`

	QueryGetPaginatedSuppliersCountBase = `SELECT COUNT(*) FROM suppliers `

	QueryGetPaginatedSuppliersDataPrefix = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, 
                  is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, 
                  nda_signed, created_at 
                  FROM suppliers `

	QueryGetPaginatedSuppliersDataSuffix = ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	QueryGetAllSuppliersPaginatedWithFilterBase = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at FROM suppliers WHERE 1=1`

	QuerySupplierDetailsProfile = `SELECT id, company_name, city, state, contact_email, phone, gst_number, msme_reg_number, is_verified, iso_standard, plant_area_sq_ft, worker_count, rating, capabilities_summary, nda_signed, created_at 
                      FROM suppliers WHERE id = ?`
)

// ====================================================================================================
// 📁 work_package.go
// ====================================================================================================
const (
	QueryCreateWorkPackage = `INSERT INTO work_packages (rfq_id, supplier_id, title, redacted_spec, redacted_drawing_path, budget_usd, target_date, escrow_status, status) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	QueryGetWorkPackageByID = `SELECT id, rfq_id, COALESCE(supplier_id, 0), title, redacted_spec, COALESCE(redacted_drawing_path, ''), budget_usd, COALESCE(target_date, ''), escrow_status, status, created_at 
              FROM work_packages WHERE id = ?`

	QueryUpdateWorkPackageStatus = `UPDATE work_packages SET status = ? WHERE id = ?`

	QueryDeleteWorkPackage = `DELETE FROM work_packages WHERE id = ?`

	QueryGetWorkPackages = `SELECT 
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

	QueryGetWorkPackageFullDetailsMain = `SELECT wp.id, wp.rfq_id, wp.supplier_id, COALESCE(s.company_name, 'Unassigned SME'), 
                wp.title, wp.redacted_spec, wp.redacted_drawing_path, wp.budget_usd, 
                wp.target_date, wp.escrow_status, wp.status, wp.created_at 
                FROM work_packages wp 
                LEFT JOIN suppliers s ON wp.supplier_id = s.id 
                WHERE wp.id = ?`

	QueryGetWorkPackageDrawingRevisions = `SELECT id, rfq_id, revision_code, original_name, redacted_path, change_log_notes, is_active, created_at 
                            FROM drawing_revisions WHERE rfq_id = ? ORDER BY created_at DESC`

	QueryGetWorkPackageMaterialCertificates = `SELECT id, work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by, created_at 
                             FROM material_certificates WHERE work_package_id = ? ORDER BY created_at DESC`

	QueryGetWorkPackageQualityInspections = `SELECT id, work_package_id, inspection_type, passed, measured_values, tolerance_allowed, report_file_path, inspected_at 
                             FROM quality_inspections WHERE work_package_id = ? ORDER BY inspected_at DESC`

	QueryGetWorkPackagePaymentTrackers = `SELECT id, work_package_id, client_amount_usd, client_paid_status, client_paid_date, 
                             supplier_amount_usd, supplier_paid_status, supplier_paid_date, notes, updated_at 
                             FROM payment_trackers WHERE work_package_id = ? ORDER BY updated_at DESC`

	QueryGetWorkPackageCommunicationLogs = `SELECT id, work_package_id, sender_role, sender_id, message, is_flagged, flag_reason, created_at 
                             FROM communication_logs WHERE work_package_id = ? ORDER BY created_at DESC`

	QueryGetFilteredWorkPackagesCountBase = `SELECT COUNT(*) FROM work_packages wp LEFT JOIN suppliers s ON wp.supplier_id = s.id `

	QueryGetFilteredWorkPackagesDataPrefix = `SELECT wp.id, wp.rfq_id, wp.supplier_id, COALESCE(s.company_name, 'Unassigned SME'), 
                  wp.title, wp.redacted_spec, wp.redacted_drawing_path, wp.budget_usd, 
                  wp.target_date, wp.escrow_status, wp.status, wp.created_at 
                  FROM work_packages wp 
                  LEFT JOIN suppliers s ON wp.supplier_id = s.id `

	QueryGetFilteredWorkPackagesDataSuffix = ` ORDER BY wp.created_at DESC LIMIT ? OFFSET ?`

	QueryGetWorkPackageCompleteProgress = `SELECT id, work_package_id, stage, percent_complete, milestone_note, inspector_name, updated_at 
                             FROM project_progress WHERE work_package_id = ? ORDER BY updated_at DESC`

	QueryGetAllWorkPackagesForDropdown = `
        SELECT id, title, status, created_at 
        FROM work_packages 
        WHERE status = ? 
        ORDER BY created_at DESC
    `

	QueryCreateWorkPackageFromRFQ = `
        INSERT INTO work_packages (
            rfq_id, 
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
            category || ' - ' || company_name, 
            COALESCE(requirements, ''), 
            COALESCE(drawing_file_path, ''), 
            target_price_usd, 
            delivery_deadline, 
            'Unfunded', 
            'current', 
            CURRENT_TIMESTAMP
        FROM rfqs 
        WHERE id = ?
    `

	QueryGetWorkPackageDetailsByID = `
        SELECT 
            wp.id, 
            COALESCE(wp.rfq_id, 0) as rfq_id, 
            COALESCE(wp.supplier_id, 0) as supplier_id, 
            COALESCE(s.company_name, 'Unassigned (Open Bidding)') as supplier_name, 
            wp.title, 
            COALESCE(wp.redacted_spec, '') as redacted_spec, 
            COALESCE(wp.redacted_drawing_path, '') as redacted_drawing_path, 
            COALESCE(wp.budget_usd, 0.0) as budget_usd, 
            COALESCE(wp.target_date, '2026-12-31') as target_date, 
            COALESCE(wp.escrow_status, 'Pending') as escrow_status, 
            wp.status, 
            wp.created_at
        FROM work_packages wp
        LEFT JOIN suppliers s ON wp.supplier_id = s.id
        WHERE wp.id = ?
    `

	QueryGetDispatchDetailsByID = `
        SELECT 
            wp.id, 
            COALESCE(wp.rfq_id, 0), 
            COALESCE(wp.supplier_id, 0), 
            COALESCE(s.company_name, 'Unassigned (Open Bidding Pool)'), 
            wp.title, 
            COALESCE(wp.redacted_spec, ''), 
            COALESCE(wp.redacted_drawing_path, ''), 
            COALESCE(wp.budget_usd, 0.0), 
            COALESCE(wp.target_date, '2026-12-31'), 
            COALESCE(wp.escrow_status, 'Pending Allocation'), 
            wp.status, 
            wp.created_at
        FROM work_packages wp
        LEFT JOIN suppliers s ON wp.supplier_id = s.id
        WHERE wp.id = ?
    `
)

const (
	QueryGetQualityInspectionsByWorkPackageID = `
SELECT 
            qi.id, 
            qi.work_package_id, 
            COALESCE(wp.title, 'General Work Package') as title,
            qi.supplier_id, 
            COALESCE(s.company_name, 'Unassigned SME'), 
            qi.inspection_type, 
            qi.passed, 
            COALESCE(qi.comments, ''), 
            COALESCE(qi.measured_values, ''), 
            COALESCE(qi.tolerance_allowed, ''), 
            COALESCE(qi.report_file_path, ''), 
            COALESCE(qi.inspected_at, '')
        FROM quality_inspections qi
        LEFT JOIN work_packages wp ON qi.work_package_id = wp.id
        LEFT JOIN work_package_suppliers wps ON wp.id = wps.work_package_id AND qi.supplier_id = wps.supplier_id
        LEFT JOIN suppliers s ON qi.supplier_id = s.id
        WHERE qi.work_package_id = ?
        ORDER BY qi.inspected_at DESC
    `
)
