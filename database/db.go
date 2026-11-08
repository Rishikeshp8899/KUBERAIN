package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB connects to the SQLite database file and triggers schema generation
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./app.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("[FATAL] Failed to open database file: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("[FATAL] Failed to ping SQLite database: %v", err)
	}

	log.Println("[SUCCESS] Connected to SQLite database successfully!")

	// Initialize all 15 schemas upon database connection
	InitTables()
	seedDefaultAdmin()
}

// seedDefaultAdmin checks if an admin exists; if not, inserts default admin credentials
func seedDefaultAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM admin_users").Scan(&count)
	if err != nil {
		log.Printf("Error checking admin user count: %v", err)
		return
	}

	// Only insert if no admin exists
	if count == 0 {
		defaultEmail := "admin"
		defaultPassword := "Admin@12345" // Change this in production!

		// Hash password using bcrypt
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}

		insertQuery := `INSERT INTO admin_users (email, password_hash, role) VALUES (?, ?, ?)`
		_, err = DB.Exec(insertQuery, defaultEmail, string(hashedBytes), "superadmin")
		if err != nil {
			log.Printf("Failed to seed default admin user: %v", err)
			return
		}

		log.Println("==================================================")
		log.Println(" Default Admin Seeded Successfully:")
		log.Printf(" Email: %s\n", defaultEmail)
		log.Printf(" Password: %s\n", defaultPassword)
		log.Println("==================================================")
	}
}

// InitTables creates all 15 required SQLite tables if they do not exist
func InitTables() {
	schemas := []string{
		// 1. Enterprise Clients
		`CREATE TABLE IF NOT EXISTS enterprise_clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			domain TEXT UNIQUE,
			billing_city TEXT,
			tax_id TEXT,
			industry_sector TEXT,
			credit_limit_usd REAL DEFAULT 0.0,
			account_status TEXT DEFAULT 'Pending_KYC',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 2. Client Users
		`CREATE TABLE IF NOT EXISTS client_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id INTEGER REFERENCES enterprise_clients(id) ON DELETE CASCADE,
			full_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			phone TEXT,
			role TEXT DEFAULT 'Buyer',
			department TEXT,
			is_active BOOLEAN DEFAULT 1,
			last_login_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 3. SME Suppliers
		`CREATE TABLE IF NOT EXISTS suppliers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_name TEXT NOT NULL,
			city TEXT NOT NULL,
			state TEXT NOT NULL,
			contact_email TEXT NOT NULL,
			phone TEXT,
			gst_number TEXT,
			msme_reg_number TEXT,
			is_verified BOOLEAN DEFAULT 0,
			iso_standard TEXT,
			plant_area_sq_ft INTEGER DEFAULT 0,
			worker_count INTEGER DEFAULT 0,
			rating REAL DEFAULT 5.0,
			capabilities_summary TEXT,
			nda_signed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 4. Supplier Capabilities
		`CREATE TABLE IF NOT EXISTS supplier_capabilities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
			category TEXT NOT NULL,
			machining_type TEXT NOT NULL,
			max_tolerance TEXT,
			max_component_dim TEXT,
			materials_handled TEXT,
			capacity_per_month INTEGER DEFAULT 100
		);`,

		// 5. Machine Statuses
		`CREATE TABLE IF NOT EXISTS machine_statuses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
			machine_name TEXT NOT NULL,
			is_operational BOOLEAN DEFAULT 1,
			last_calibration TEXT,
			next_service_due TEXT,
			notes TEXT
		);`,

		// 6. Anonymous Guest Inquiries
		`CREATE TABLE IF NOT EXISTS inquiries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			phone TEXT,
			subject TEXT,
			message TEXT NOT NULL,
			source TEXT DEFAULT 'Contact Form',
			ip_address TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 7. Corporate RFQs
		`CREATE TABLE IF NOT EXISTS rfqs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rfq_code TEXT UNIQUE NOT NULL,
			company_name TEXT NOT NULL,
			contact_person TEXT NOT NULL,
			email TEXT NOT NULL,
			phone TEXT,
			category TEXT,
			material TEXT,
			quantity TEXT,
			requirements TEXT,
			drawing_file_path TEXT,
			target_price_usd REAL DEFAULT 0.0,
			delivery_deadline TEXT,
			incobooking_terms TEXT,
			country TEXT DEFAULT 'India',
			city TEXT DEFAULT 'Unknown',
			ip_address TEXT,
			status TEXT DEFAULT 'New',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 8. Drawing Revision Control
		`CREATE TABLE IF NOT EXISTS drawing_revisions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rfq_id INTEGER REFERENCES rfqs(id) ON DELETE CASCADE,
			revision_code TEXT NOT NULL,
			original_name TEXT NOT NULL,
			redacted_path TEXT NOT NULL,
			change_log_notes TEXT,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 9. Dispatched Work Packages (removed single supplier_id)
		`CREATE TABLE IF NOT EXISTS work_packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rfq_id INTEGER REFERENCES rfqs(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    redacted_spec TEXT NOT NULL,
    redacted_drawing_path TEXT,
    budget_usd REAL DEFAULT 0.0,
    target_date TEXT,
    escrow_status TEXT DEFAULT 'Unfunded',
    status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`,

		// 9b. Junction Table for Many-to-Many Work Package <-> Suppliers
		`CREATE TABLE IF NOT EXISTS work_package_suppliers (
    work_package_id INTEGER REFERENCES work_packages(id) ON DELETE CASCADE,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
    assignment_status TEXT DEFAULT 'Invited', -- e.g., Invited, Bidding, Awarded, Declined
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (work_package_id, supplier_id)
);`,

		// 10. Project Progress Tracking
		`CREATE TABLE IF NOT EXISTS project_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    work_package_id INTEGER REFERENCES work_packages(id) ON DELETE CASCADE,
    stage TEXT NOT NULL DEFAULT 'Sourcing',
    percent_complete INTEGER DEFAULT 0,
    milestone_note TEXT,
    photos_path TEXT,
    inspection_doc TEXT,
    inspector_name TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`,

		// 11. Material Test Certificates (MTC)
		`CREATE TABLE IF NOT EXISTS material_certificates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			work_package_id INTEGER REFERENCES work_packages(id) ON DELETE CASCADE,
			heat_number TEXT NOT NULL,
			material_grade TEXT NOT NULL,
			mill_name TEXT,
			cert_file_path TEXT NOT NULL,
			verified_by TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 12. Quality Inspection Checks
		`CREATE TABLE IF NOT EXISTS quality_inspections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        work_package_id INTEGER NOT NULL,
        supplier_id INTEGER NOT NULL,
        inspection_type TEXT NOT NULL,
		passed BOOLEAN NOT NULL DEFAULT 0,
		comments TEXT,
		measured_values TEXT,
		tolerance_allowed TEXT,
		report_file_path TEXT,
		inspected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (work_package_id) REFERENCES work_packages(id) ON DELETE CASCADE,
		FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE);`,

		// 13. Payment Tracker
		`CREATE TABLE IF NOT EXISTS payment_trackers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			work_package_id INTEGER REFERENCES work_packages(id) ON DELETE CASCADE,
			client_amount_usd REAL DEFAULT 0.0,
			client_paid_status TEXT DEFAULT 'Pending',
			client_paid_date TEXT,
			supplier_amount_usd REAL DEFAULT 0.0,
			supplier_paid_status TEXT DEFAULT 'Pending',
			supplier_paid_date TEXT,
			notes TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 14. Communication Audit Logs
		`CREATE TABLE IF NOT EXISTS communication_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			work_package_id INTEGER REFERENCES work_packages(id) ON DELETE CASCADE,
			sender_role TEXT NOT NULL,
			sender_id INTEGER,
			message TEXT NOT NULL,
			is_flagged BOOLEAN DEFAULT 0,
			flag_reason TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// 15. Leak Detection Rules
		`CREATE TABLE IF NOT EXISTS leak_detection_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rule_name TEXT NOT NULL,
			regex_pattern TEXT NOT NULL,
			severity TEXT DEFAULT 'Immediate_Flag',
			is_active BOOLEAN DEFAULT 1
		);`,

		// 16. Admin System Users
		`CREATE TABLE IF NOT EXISTS admin_users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'Admin',
			is_active BOOLEAN DEFAULT 1
			);`,
		// 17. Table: mtc_certificates
		`CREATE TABLE IF NOT EXISTS mtc_certificates (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
			work_package_id INTEGER NOT NULL,
			supplier_id INTEGER NOT NULL,
			heat_number TEXT NOT NULL,
			material_grade TEXT NOT NULL,
			mill_name TEXT NOT NULL,
			cert_file_path TEXT NOT NULL,
			verified_by TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (work_package_id) REFERENCES work_packages(id) ON DELETE CASCADE,
			FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE);`,

		// 3. Table: communication_logs
		`CREATE TABLE IF NOT EXISTS communication_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    work_package_id INTEGER NOT NULL,
    sender_role TEXT NOT NULL,
    sender_id INTEGER NOT NULL,
    flag_reason TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (work_package_id) REFERENCES work_packages(id) ON DELETE CASCADE);`,
	}

	for _, statement := range schemas {
		_, err := DB.Exec(statement)
		if err != nil {
			log.Fatalf("[FATAL] Error initializing SQLite schema: %v", err)
		}
	}

	log.Println("[SUCCESS] All 15 database tables verified/created successfully.")
}
