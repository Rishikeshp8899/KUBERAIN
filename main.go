package main

import (
	"fmt"
	"go-mvc-app/controllers"
	"go-mvc-app/database"
	"go-mvc-app/middleware"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	mux := http.NewServeMux()

	// ==========================================
	// PUBLIC WEB PAGES & CLIENT APIS
	// ==========================================

	// Renders the main homepage
	mux.HandleFunc("/", controllers.HomePageHandler)

	// Displays product catalog and offerings
	mux.HandleFunc("/products", controllers.ProductsPageHandler)

	// Displays company services
	mux.HandleFunc("/services", controllers.ServicesPageHandler)

	// Renders the Request for Quote (RFQ) submission form
	mux.HandleFunc("/rfq", controllers.RFQPageHandler)

	// Processes incoming public RFQ form submissions
	mux.HandleFunc("/submit-rfq", controllers.SubmitRFQHandler)

	// Renders the public SME supplier onboarding page
	mux.HandleFunc("/smes-onboarding", controllers.PublicSMEOnboarding)

	// Displays company "About Us" information
	mux.HandleFunc("/about", controllers.AboutUsPageHandler)

	// Displays compliance details and certifications
	mux.HandleFunc("/compliance", controllers.ComplianceHandler)

	// Shows the public SME network directory overview
	mux.HandleFunc("/sme-network", controllers.SmeNwtworkHandler)

	// API endpoint for public SME registration submissions
	mux.HandleFunc("/suppliers/register", controllers.SMESubmitHandler)

	// ==========================================
	// ADMIN AUTHENTICATION
	// ==========================================

	// Handles admin login page and authentication requests
	mux.HandleFunc("/admin/login", controllers.AdminLoginHandler)

	// Destroys admin session and logs the user out
	mux.HandleFunc("/admin/logout", controllers.AdminLogoutHandler)

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - RFQ MANAGEMENT
	// ==========================================

	// Displays and manages admin RFQ listings
	mux.HandleFunc("/admin/rfqs", middleware.RequireAdminAuth(controllers.AdminRFQsHandler))

	// Updates the status of an RFQ
	mux.HandleFunc("/admin/rfqs/update-rfq-status", middleware.RequireAdminAuth(controllers.AdminUpdateRFQStatusHandler))

	// Handles RFQ detail view, matching logic, and status updates
	mux.HandleFunc("/admin/rfqs/update-status", middleware.RequireAdminAuth(controllers.AdminRFQDetailAndMatchHandler))

	mux.HandleFunc("/admin/rfqs/screenshot/{id}", middleware.RequireAdminAuth(controllers.AdminRFQDetailAndMatchScreenshotHandler))

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - SME & SUPPLIER MANAGEMENT
	// ==========================================

	// Displays admin view of SME suppliers
	mux.HandleFunc("/admin/sme", middleware.RequireAdminAuth(controllers.AdminSMEHandler))

	// Handles adding a new SME supplier through the admin panel
	mux.HandleFunc("/admin/add-sme", middleware.RequireAdminAuth(controllers.AdminAddSMEHandler))

	// Updates SME status related to RFQs
	mux.HandleFunc("/admin/rfqs/update-sme-status", middleware.RequireAdminAuth(controllers.AdminUpdateSMEStatusHandler))

	// API handler for supplier onboarding status
	mux.HandleFunc("/admin/sme/onboarding-status", middleware.RequireAdminAuth(controllers.APISuppliersHandler))

	// Renders detailed view for a specific supplier
	mux.HandleFunc("/admin/sme/details", middleware.RequireAdminAuth(controllers.RenderSupplierDetailsPage))

	// Renders the supplier onboarding workflow page
	mux.HandleFunc("/admin/sme/onboarding", middleware.RequireAdminAuth(controllers.RenderOnboardingPage))

	mux.HandleFunc("/admin/work-packages/{id}/suppliers", middleware.RequireAdminAuth(controllers.HandleGetWorkPackageSuppliers))

	mux.HandleFunc("/admin/suppliers/status/{id}/{status}", middleware.RequireAdminAuth(controllers.HandleUpdateSupplierStatus))

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - WORK PACKAGES
	// ==========================================

	// Handles work packages filtered by status and pagination
	mux.HandleFunc("/admin/work-packages/status", middleware.RequireAdminAuth(controllers.AdminWorkPackageDetailsHandler))

	// Sends communication/messages within a specific work package
	mux.HandleFunc("/admin/work-packages/messages/add", middleware.RequireAdminAuth(controllers.AdminSendCommunicationHandler))

	// API handler for listing and managing work packages
	mux.HandleFunc("/admin/work-packages", middleware.RequireAdminAuth(controllers.APIWorkPackagesHandler))

	// Renders general work package details view
	mux.HandleFunc("/admin/work-packages/details", middleware.RequireAdminAuth(controllers.RenderDetails))

	// Renders the detailed management page for a specific work package by ID
	mux.HandleFunc("/admin/work-packages/{id}", middleware.RequireAdminAuth(controllers.RenderWorkPackageDetailsPage))

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - PAYMENTS & ESCROW
	// ==========================================

	// Displays the double-entry payment ledger and escrow records
	mux.HandleFunc("/admin/payments", middleware.RequireAdminAuth(controllers.AdminPaymentsHandler))

	// Handles adding a new manual payment entry or multi-SME payout split
	mux.HandleFunc("/admin/payments/add", middleware.RequireAdminAuth(controllers.AdminAddPaymentHandler))

	// Updates payment and escrow status records
	mux.HandleFunc("/admin/update-payment-status", middleware.RequireAdminAuth(controllers.AdminUpdatePaymentStatusHandler))

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - COMPLIANCE & QUALITY AUDITS
	// ==========================================

	// Manages compliance records and audits in the admin panel
	mux.HandleFunc("/admin/compliance", middleware.RequireAdminAuth(controllers.AdminComplianceHandler))

	// Manages quality audit listings and records
	mux.HandleFunc("/admin/quality-audit", middleware.RequireAdminAuth(controllers.AdminQualityAuditHandler))

	// Handles the creation of a new audit record
	mux.HandleFunc("/admin/audit/create", middleware.RequireAdminAuth(controllers.HandleCreateAudit))

	// Renders the form to add a new quality audit
	mux.HandleFunc("/admin/audit-form/add", middleware.RequireAdminAuth(controllers.GetAddAuditFormHandler))

	mux.HandleFunc("/admin/work/quality/{id}", middleware.RequireAdminAuth(controllers.ViewAllPreviousAudits))

	// ==========================================
	// ADMIN PROTECTED OPERATIONS - DISPATCH & LOGISTICS
	// ==========================================

	// Manages overall dispatch and logistics operations
	mux.HandleFunc("/admin/work/packages", middleware.RequireAdminAuth(controllers.AdminDispatchHandler))

	// Renders the dispatch tracking page for a specific project ID
	mux.HandleFunc("/admin/work/packages/{id}", middleware.RequireAdminAuth(controllers.RenderDispatchPerProjectPage))

	// Go 1.22+ standard http.ServeMux
	mux.HandleFunc("/admin/download", middleware.RequireAdminAuth(controllers.DownloadUploadedFilePOST))

	// ==========================================
	// STATIC & VIEW ASSETS
	// ==========================================

	// Serves static files (CSS, JavaScript, images)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serves template views assets
	vs := http.FileServer(http.Dir("./views"))
	mux.Handle("/views/", http.StripPrefix("/views/", vs))

	// Server Initialization
	port := ":9091"
	fmt.Printf("Richim Valve Admin Portal running at http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
