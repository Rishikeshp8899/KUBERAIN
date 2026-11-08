package models

import "time"

type EnterpriseClient struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Domain         string  `json:"domain"`
	BillingCity    string  `json:"billing_city"`
	TaxID          string  `json:"tax_id"`           // GSTIN / VAT identifier for invoicing
	IndustrySector string  `json:"industry_sector"`  // e.g., Valves, Aerospace, Oil & Gas
	CreditLimitUSD float64 `json:"credit_limit_usd"` // Maximum open credit allowed
	AccountStatus  string  `json:"account_status"`   // Pending_KYC, Active, Suspended
	CreatedAt      string  `json:"created_at"`
}

type ClientUser struct {
	ID          int    `json:"id"`
	ClientID    int    `json:"client_id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`      // Direct phone line for alerts
	Role        string `json:"role"`       // Procurement_Manager, Lead_Engineer
	Department  string `json:"department"` // e.g., QA, Supply Chain, R&D
	IsActive    bool   `json:"is_active"`  // Soft disable flag for accounts
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
}

type Supplier struct {
	ID                  int     `json:"id"`
	CompanyName         string  `json:"company_name"`
	City                string  `json:"city"`
	State               string  `json:"state"`
	ContactEmail        string  `json:"contact_email"`
	Phone               string  `json:"phone"`
	GSTNumber           string  `json:"gst_number"`      // Legal tax registration number
	MSMERegNumber       string  `json:"msme_reg_number"` // Government MSME/Udyam registration code
	IsVerified          bool    `json:"is_verified"`
	ISOStandard         string  `json:"iso_standard"`     // e.g., ISO 9001:2015
	PlantAreaSqFt       int     `json:"plant_area_sq_ft"` // Total shop floor space in sq ft
	WorkerCount         int     `json:"worker_count"`     // Total workforce count
	Rating              float64 `json:"rating"`           // Vendor performance score (0-5.0)
	CapabilitiesSummary string  `json:"capabilities_summary"`
	NDASigned           bool    `json:"nda_signed"` // Flag if master platform NDA is signed
	CreatedAt           string  `json:"created_at"`
}

type SupplierCapability struct {
	ID               int    `json:"id"`
	SupplierID       int    `json:"supplier_id"`
	Category         string `json:"category"`           // e.g., CNC Machining, Casting, Forging
	MachiningType    string `json:"machining_type"`     // e.g., 5-Axis Milling, VMC, Lathe
	MaxTolerance     string `json:"max_tolerance"`      // e.g., +/- 0.005mm threshold
	MaxComponentDim  string `json:"max_component_dim"`  // Max workpiece dimensions (LxWxH)
	MaterialsHandled string `json:"materials_handled"`  // e.g., SS316, Duplex, Inconel
	CapacityPerMonth int    `json:"capacity_per_month"` // Monthly unit capacity
}

// MachineStatus tracks shop-floor machine health and calibration
type MachineStatus struct {
	ID              int    `json:"id"`
	SupplierID      int    `json:"supplier_id"`
	MachineName     string `json:"machine_name"`     // e.g., Haas VF-4 4-Axis VMC
	IsOperational   bool   `json:"is_operational"`   // Down for maintenance if false
	LastCalibration string `json:"last_calibration"` // Date of last precision audit
	NextServiceDue  string `json:"next_service_due"` // Service schedule date
	Notes           string `json:"notes"`
}

type Inquiry struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	Source    string `json:"source"` // e.g., Contact Form, Landing Banner
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
}

type RFQ struct {
	ID                int      `json:"id"`
	RFQCode           string   `json:"rfq_code"` // e.g., RFQ-20261020-1024
	CompanyName       string   `json:"company_name"`
	ContactPerson     string   `json:"contact_person"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	Category          string   `json:"category"`                // Product or process category (e.g., Industrial Valves)
	Material          string   `json:"material"`                // Material specification (e.g., SS316, ASTM A216 WCB)
	Quantity          string   `json:"quantity"`                // Order quantity (e.g., 500 Units)
	Requirements      string   `json:"requirements"`            // BOQ details and technical specs
	DrawingFolderPath string   `json:"drawing_folder_path"`     // Path to request folder: static/uploads/drawings/RFQ-XXXXXX/
	DrawingFilePath   string   `json:"drawing_file_path"`       // Comma-separated web paths of uploaded blueprint files
	DrawingFiles      []string `json:"drawing_files,omitempty"` // Parsed slice of individual file paths for template rendering
	TargetPriceUSD    float64  `json:"target_price_usd"`        // Target price per unit in USD
	DeliveryDeadline  string   `json:"delivery_deadline"`       // Target delivery completion date (YYYY-MM-DD)
	IncobookingTerms  string   `json:"incobooking_terms"`       // e.g., FOB, CIF, EXW, DDP
	Country           string   `json:"country"`                 // Submission country
	City              string   `json:"city"`                    // Submission city / location
	IPAddress         string   `json:"ip_address"`              // Client IP address
	Status            string   `json:"status"`                  // New, Under Review, Dispatched, Closed
	CreatedAt         string   `json:"created_at"`              // Timestamp of creation
}

type DrawingRevision struct {
	ID             int    `json:"id"`
	RFQID          int    `json:"rfq_id"`
	RevisionCode   string `json:"revision_code"`    // e.g., Rev_01, Rev_02B
	OriginalName   string `json:"original_name"`    // File name (e.g., Valve_Body_v2.step)
	RedactedPath   string `json:"redacted_path"`    // Watermarked and identity-stripped file
	ChangeLogNotes string `json:"change_log_notes"` // Engineering ECN/ECO explanation
	IsActive       bool   `json:"is_active"`        // Set to true for the current active version
	CreatedAt      string `json:"created_at"`
}

type WorkPackage struct {
	ID                  int     `json:"id"`
	RFQID               int     `json:"rfq_id"`
	SupplierID          int     `json:"supplier_id"`
	SupplierName        string  `json:"supplier_name"` // Joined from suppliers table
	Title               string  `json:"title"`
	RedactedSpec        string  `json:"redacted_spec"`
	RedactedDrawingPath string  `json:"redacted_drawing_path"`
	BudgetUSD           float64 `json:"budget_usd"`
	TargetDate          string  `json:"target_date"`
	EscrowStatus        string  `json:"escrow_status"`
	Status              string  `json:"status"`
	CreatedAt           string  `json:"created_at"`
}

type ProjectProgress struct {
	ID              int    `json:"id"`
	WorkPackageID   int    `json:"work_package_id"`
	Stage           string `json:"stage"`            // e.g., Sourcing, Machining, Heat Treatment, QA, Dispatched
	PercentComplete int    `json:"percent_complete"` // Milestone completion percentage (0 to 100)
	MilestoneNote   string `json:"milestone_note"`   // Notes from shop floor operations
	PhotosPath      string `json:"photos_path"`      // Photos URL / file path
	InspectionDoc   string `json:"inspection_doc"`   // QA Report PDF file path
	InspectorName   string `json:"inspector_name"`   // QA engineer or third-party inspector name
	UpdatedAt       string `json:"updated_at"`       // Timestamp of last milestone update
}

type PaymentTracker struct {
	ID            int `json:"id"`
	WorkPackageID int `json:"work_package_id"`

	// Client -> Admin Payment
	ClientAmountUSD  float64 `json:"client_amount_usd"`  // Amount buyer owes us
	ClientPaidStatus string  `json:"client_paid_status"` // Pending, Paid, Overdue
	ClientPaidDate   string  `json:"client_paid_date"`

	// Admin -> SME Payment
	SupplierAmountUSD  float64 `json:"supplier_amount_usd"`  // Amount we pay to the SME
	SupplierPaidStatus string  `json:"supplier_paid_status"` // Pending, Paid, Partially_Paid
	SupplierPaidDate   string  `json:"supplier_paid_date"`

	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
}

type CommunicationLog struct {
	ID            int    `json:"id"`
	WorkPackageID int    `json:"work_package_id"`
	SenderRole    string `json:"sender_role"` // Client, Admin, Supplier
	SenderID      int    `json:"sender_id"`
	Message       string `json:"message"`
	IsFlagged     bool   `json:"is_flagged"`  // True if contact leak is detected
	FlagReason    string `json:"flag_reason"` // Reason code (e.g., Phone number shared)
	CreatedAt     string `json:"created_at"`
}

type LeakDetectionRule struct {
	ID           int    `json:"id"`
	RuleName     string `json:"rule_name"`     // e.g., Phone Regex, Email Keyword
	RegexPattern string `json:"regex_pattern"` // Match pattern
	Severity     string `json:"severity"`      // Warning, Immediate_Flag
	IsActive     bool   `json:"is_active"`
}

type AdminUser struct {
	ID           int64  `db:"id" json:"id"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"` // Omitted from JSON serialization for security
	Role         string `db:"role" json:"role"`
	IsActive     bool   `db:"is_active" json:"is_active"`
}

type PaginatedResult struct {
	Suppliers   []Supplier
	TotalPages  int
	CurrentPage int
	Filter      string
	Search      string
}

type WorkPackageDetailView struct {
	WorkPackage      WorkPackage
	DrawingRevisions []DrawingRevision
	Certificates     []MaterialCertificate
	Inspections      []QualityInspection
	Payments         []PaymentTracker
	Communications   []CommunicationLog
}

type WorkPackagePaginatedResult struct {
	WorkPackages []WorkPackage `json:"work_packages"`
	TotalPages   int           `json:"total_pages"`
	CurrentPage  int           `json:"current_page"`
	PrevPage     int           `json:"prev_page"`
	NextPage     int           `json:"next_page"`
	StatusFilter string        `json:"status_filter"`
	Search       string        `json:"search"`
}

type WorkPackageCompleteDetail struct {
	WorkPackage      WorkPackage
	Supplier         Supplier
	Capability       SupplierCapability
	Progress         []ProjectProgress
	DrawingRevisions []DrawingRevision
	Certificates     []MaterialCertificate
	Inspections      []QualityInspection
	Payments         []PaymentTracker
	Communications   []CommunicationLog
}

type QualityInspection struct {
	ID               int    `json:"id"`
	WorkPackageID    int    `json:"work_package_id"`
	SupplierID       int    `json:"supplier_id"`
	SupplierName     string `json:"supplier_name"`
	InspectionType   string `json:"inspection_type"`
	Passed           bool   `json:"passed"`
	Comments         string `json:"comments"`
	MeasuredValues   string `json:"measured_values"`
	ToleranceAllowed string `json:"tolerance_allowed"`
	ReportFilePath   string `json:"report_file_path"`
	InspectedAt      string `json:"inspected_at"`
}

type MaterialCertificate struct {
	ID            int    `json:"id"`
	WorkPackageID int    `json:"work_package_id"`
	SupplierID    int    `json:"supplier_id"`
	SupplierName  string `json:"supplier_name"`
	HeatNumber    string `json:"heat_number"`
	MaterialGrade string `json:"material_grade"`
	MillName      string `json:"mill_name"`
	CertFilePath  string `json:"cert_file_path"`
	VerifiedBy    string `json:"verified_by"`
	CreatedAt     string `json:"created_at"`
}

type FlaggedLog struct {
	ID            int    `json:"id"`
	WorkPackageID int    `json:"work_package_id"`
	SenderRole    string `json:"sender_role"`
	SenderID      int    `json:"sender_id"`
	FlagReason    string `json:"flag_reason"`
	Message       string `json:"message"`
	CreatedAt     string `json:"created_at"`
}

type WorkPackageSupplier struct {
	WorkPackageID    int       `json:"work_package_id" db:"work_package_id"`
	SupplierID       int       `json:"supplier_id" db:"supplier_id"`
	AssignmentStatus string    `json:"assignment_status" db:"assignment_status"` // e.g., "Invited", "Bidding", "Awarded", "Declined"
	AssignedAt       time.Time `json:"assigned_at" db:"assigned_at"`
	SupplierName     string    `json:"supplier_name,omitempty" db:"supplier_name"`
	WorkTitle        string    `json:"work_title,omitempty" db:"work_title"`
}

type SupplierDetails struct {
	Title      string
	Supplier   Supplier
	Capability SupplierCapability
	Machines   []MachineStatus
}
