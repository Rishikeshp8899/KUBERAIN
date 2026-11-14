package controllers

import (
	// "bytes"

	"database/sql"
	"encoding/json"
	"errors"
	"go-mvc-app/config"
	"go-mvc-app/middleware"
	"go-mvc-app/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func renderAdminView(w http.ResponseWriter, pageFile string, data interface{}) {

	layoutPath := filepath.Join("views", "layouts", "admin.html")
	pagePath := filepath.Join("views", "pages", pageFile)

	tmpl, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		log.Printf("[ERROR] Template parsing failed (%s): %v", pageFile, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("[ERROR] Template execution failed (%s): %v", pageFile, err)
	}

}

// AdminRFQsHandler handles GET /admin/rfqs
func AdminRFQsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	adminEmail := config.GetAdminEmail(r)
	rfqs, err := models.GetAllRFQs()
	if err != nil {
		rfqs = []models.RFQ{}
	}

	data := map[string]interface{}{
		"Title":      "RFQ Monitoring Desk",
		"AdminEmail": adminEmail,
		"Page":       "rfqs",
		"RFQs":       rfqs,
	}

	renderAdminView(w, "admin_rfqs.html", data)
}

// AdminSMEHandler handles GET /admin/sme
func AdminSMEHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse Query Parameters
	filter := r.URL.Query().Get("filter") // "all", "approved", "pending"
	if filter == "" {
		filter = "all"
	}

	search := r.URL.Query().Get("search")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize := 10 // Number of suppliers per page

	// Fetch Paginated & Filtered Data
	paginatedData, err := models.GetPaginatedSuppliers(filter, search, page, pageSize)
	if err != nil {
		log.Printf("[ERROR] [AdminSMEPaginatedHandler] Failed to load paginated suppliers: %v", err)
		paginatedData = models.PaginatedResult{
			Suppliers:   []models.Supplier{},
			TotalPages:  1,
			CurrentPage: 1,
			Filter:      filter,
			Search:      search,
		}
	}

	// Construct View Data Context
	data := map[string]interface{}{
		"Title":      "SME Supplier Directory",
		"Page":       "sme",
		"Suppliers":  paginatedData.Suppliers,
		"Pagination": paginatedData,
		"Filter":     filter,
		"Search":     search,
	}

	// Render View (Assuming renderAdminView uses your admin layout)
	renderAdminView(w, "admin_sme.html", data)

}

// AdminAddSMEHandler handles POST /admin/add-sme (Form submission from sme_modal.html)
func AdminAddSMEHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	workerCount, _ := strconv.Atoi(r.FormValue("worker_count"))
	isVerified := r.FormValue("is_verified") == "true" || r.FormValue("is_verified") == "on"

	supplier := models.Supplier{
		CompanyName:         r.FormValue("company_name"),
		ContactEmail:        r.FormValue("contact_email"),
		Phone:               r.FormValue("phone"),
		City:                r.FormValue("city"),
		State:               r.FormValue("state"),
		GSTNumber:           r.FormValue("gst_number"),
		MSMERegNumber:       r.FormValue("msme_reg_number"),
		ISOStandard:         r.FormValue("iso_standard"),
		IsVerified:          isVerified,
		WorkerCount:         workerCount,
		CapabilitiesSummary: r.FormValue("capabilities_summary"),
	}

	err := models.CreateSMESupplier(supplier)
	if err != nil {
		log.Printf("[ERROR] Failed to create SME supplier: %v", err)
	}

	// Redirect back to the SME directory so the user sees the newly added supplier
	http.Redirect(w, r, "/admin/sme", http.StatusSeeOther)
}

// AdminUpdateSMEStatusHandler handles POST /admin/update-sme-status
func AdminUpdateSMEStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	smeIDStr := r.FormValue("sme_id")
	isVerifiedStr := r.FormValue("is_verified")

	smeID, err := strconv.ParseInt(smeIDStr, 10, 64)
	if err == nil && smeID > 0 {
		isVerified := isVerifiedStr == "true" || isVerifiedStr == "1"
		_ = models.UpdateSMEVerificationStatus(smeID, isVerified)
	}

	http.Redirect(w, r, "/admin/rfqs", http.StatusSeeOther)
}

// AdminPaymentsHandler handles GET /admin/payments
func AdminPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	adminEmail := config.GetAdminEmail(r)
	payments, err := models.GetPaymentLedger()

	log.Printf("[INFO] [AdminPaymentsHandler] Retrieved %d payment records for admin: %s", len(payments), adminEmail)

	if err != nil {
		payments = []models.PaymentTrackerLog{}
	}

	data := map[string]interface{}{
		"Title":      "Double-Entry Payment Ledger",
		"AdminEmail": adminEmail,
		"Page":       "payments",
		"Payments":   payments,
	}

	renderAdminView(w, "admin_payment.html", data)
}

// AdminUpdatePaymentStatusHandler handles POST /admin/update-payment-status
func AdminUpdatePaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	paymentIDStr := r.FormValue("payment_id")
	clientStatus := r.FormValue("client_paid_status")
	supplierStatus := r.FormValue("supplier_paid_status")

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err == nil && paymentID > 0 {
		_ = models.UpdatePaymentStatus(int(paymentID), clientStatus, supplierStatus)
	}

	http.Redirect(w, r, "/admin/payments", http.StatusSeeOther)
}

// AdminDispatchHandler handles GET /admin/dispatch
func AdminDispatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve authenticated email from context
	adminEmail := config.GetAdminEmail(r)

	// Fetch active work packages from database
	packages, err := models.GetWorkPackages()
	if err != nil {
		packages = []models.WorkPackage{}
	}

	data := map[string]interface{}{
		"Title":        "Work Package Dispatch & Bidding",
		"AdminEmail":   adminEmail,
		"Page":         "dispatch",
		"WorkPackages": packages,
	}

	// Render page using layout views/layouts/admin.html and view views/pages/admin_dispatch.html
	renderAdminView(w, "admin_dispatch.html", data)
}

// -----------------------------------------------------------------------------
// MODULE 5: QUALITY & SHIELD COMPLIANCE HANDLER
// -----------------------------------------------------------------------------

// AdminComplianceHandler handles GET /admin/compliance
func AdminComplianceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve authenticated email from context
	adminEmail := config.GetAdminEmail(r)

	// Fetch flagged off-platform messages and Material Test Certificates
	flaggedLogs, err := models.GetFlaggedCommunicationLogs()
	if err != nil {
		flaggedLogs = []models.FlaggedLog{}
	}

	mtcCerts, err := models.GetMaterialCertificates()
	if err != nil {
		mtcCerts = []models.MaterialCertificate{}
	}

	data := map[string]interface{}{
		"Title":           "Quality & Shield Compliance",
		"AdminEmail":      adminEmail,
		"Page":            "compliance",
		"FlaggedLogs":     flaggedLogs,
		"MTCCertificates": mtcCerts,
	}

	// Render page using layout views/layouts/admin.html and view views/pages/admin_compliance.html
	renderAdminView(w, "admin_compliance.html", data)
}

// AdminLoginHandler processes GET (render) and POST (submit) requests
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// GET: Render Login Screen
	if r.Method == http.MethodGet {
		// Redirect if already logged in
		cookie, err := r.Cookie(middleware.SessionCookieName)
		if err == nil && cookie.Value != "" {
			http.Redirect(w, r, "/admin/rfqs", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("views/pages/admin_login.html")
		if err != nil {
			http.Error(w, "Login Template Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Error": r.URL.Query().Get("error"),
		}
		tmpl.Execute(w, data)
		return
	}

	// POST: Handle Login Form Submission
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/admin/login?error=missing_fields", http.StatusSeeOther)
			return
		}

		// Verify credentials against database
		admin, err := models.AuthenticateAdmin(email, password)
		if err != nil {
			http.Redirect(w, r, "/admin/login?error=invalid_credentials", http.StatusSeeOther)
			return
		}

		// Set HTTP-only auth cookie
		middleware.SetAuthCookie(w, admin.Email)

		// Redirect to protected dashboard
		http.Redirect(w, r, "/admin/rfqs", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// AdminLogoutHandler clears user session and redirects to login
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	middleware.ClearAuthCookie(w)
	http.Redirect(w, r, "/admin/login?status=logged_out", http.StatusSeeOther)
}

// AdminRFQDetailAndMatchHandler handles GET /admin/rfqs/update-status
// Protected endpoint to view specific RFQ status, audit details, and filter capable SMEs.
func AdminRFQDetailAndMatchHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [AdminRFQDetailAndMatchHandler] Incoming request: Method=%s, RemoteAddr=%s, URI=%s", r.Method, r.RemoteAddr, r.RequestURI)

	// Step 1: Authentication / Session Check
	adminEmail := config.GetAdminEmail(r)
	if adminEmail == "" {
		log.Printf("[WARN] [AdminRFQDetailAndMatchHandler] Unauthorized access attempt detected.")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	// Step 2: Parse Query Parameters for RFQ ID
	rfqIDStr := r.FormValue("rfq_id")
	if rfqIDStr == "" {
		log.Printf("[WARN] [AdminRFQDetailAndMatchHandler] Missing 'id' query parameter.")
		http.Error(w, "Bad Request: Missing RFQ ID", http.StatusBadRequest)
		return
	}

	rfqID, err := strconv.Atoi(rfqIDStr)
	if err != nil || rfqID <= 0 {
		log.Printf("[ERROR] [AdminRFQDetailAndMatchHandler] Invalid RFQ ID format: '%s'", rfqIDStr)
		http.Error(w, "Invalid RFQ ID", http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] [AdminRFQDetailAndMatchHandler] Target RFQ ID parsed successfully: %d", rfqID)

	// Step 3: Fetch Target RFQ Record
	rfq, err := models.GetRFQByID(rfqID)
	if err != nil {
		log.Printf("[ERROR] [AdminRFQDetailAndMatchHandler] RFQ record not found for ID %d: %v", rfqID, err)
		http.Error(w, "RFQ Not Found", http.StatusNotFound)
		return
	}
	log.Printf("[SUCCESS] [AdminRFQDetailAndMatchHandler] Fetched RFQ Code: %s, Category: %s, Material: %s, Current Status: %s",
		rfq.RFQCode, rfq.Category, rfq.Material, rfq.Status)

	// Step 4: Filter Capable SMEs Based on RFQ Specifications
	log.Printf("[INFO] [AdminRFQDetailAndMatchHandler] Filtering capable SMEs for Category='%s', Material='%s'...", rfq.Category, rfq.Material)
	capableSuppliers, err := models.GetCapableSuppliersForRFQ(rfq.Category, rfq.Material)
	if err != nil {
		log.Printf("[ERROR] [AdminRFQDetailAndMatchHandler] Failed to query capable suppliers: %v", err)
		capableSuppliers = []models.Supplier{}
	}
	log.Printf("[SUCCESS] [AdminRFQDetailAndMatchHandler] Found %d matching capable SME suppliers.", len(capableSuppliers))

	// Step 5: Construct View Data Context
	data := map[string]interface{}{
		"Title":            "RFQ Review & SME Matching Desk",
		"AdminEmail":       adminEmail,
		"Page":             "rfq_detail",
		"RFQ":              rfq,
		"CapableSuppliers": capableSuppliers,
	}

	// Step 6: Render View
	log.Printf("[INFO] [AdminRFQDetailAndMatchHandler] Rendering template 'admin_rfq_detail.html'...")
	renderAdminView(w, "admin_rfq_detail.html", data)
	log.Printf("[SUCCESS] [AdminRFQDetailAndMatchHandler] Request completed successfully for RFQ ID=%d", rfqID)
}

// AdminSMEPaginatedHandler handles GET /admin/sme
func AdminSMEPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse Query Parameters
	filter := r.URL.Query().Get("filter") // "all", "approved", "pending"
	if filter == "" {
		filter = "all"
	}

	search := r.URL.Query().Get("search")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize := 6 // Number of suppliers per page

	// Fetch Paginated & Filtered Data
	paginatedData, err := models.GetPaginatedSuppliers(filter, search, page, pageSize)
	if err != nil {
		log.Printf("[ERROR] [AdminSMEPaginatedHandler] Failed to load paginated suppliers: %v", err)
		paginatedData = models.PaginatedResult{
			Suppliers:   []models.Supplier{},
			TotalPages:  1,
			CurrentPage: 1,
			Filter:      filter,
			Search:      search,
		}
	}

	// Construct View Data Context
	data := map[string]interface{}{
		"Title":      "SME Supplier Directory",
		"Page":       "sme",
		"Suppliers":  paginatedData.Suppliers,
		"Pagination": paginatedData,
		"Filter":     filter,
		"Search":     search,
	}

	// Render View (Assuming renderAdminView uses your admin layout)
	renderAdminView(w, "admin_rfq_detail.html", data)
}

// AdminWorkPackageDetailsHandler handles GET /admin/work-packages/details?id=X
func AdminWorkPackageDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	wpID, err := strconv.Atoi(idStr)
	if err != nil || wpID < 1 {
		http.Redirect(w, r, "/admin/work-packages", http.StatusSeeOther)
		return
	}
	log.Printf("[INFO] [AdminWorkPackageDetailsHandler] Fetching details for Work Package ID: %d", wpID)

	details, err := models.GetWorkPackageCompleteDetails(wpID)
	if err != nil {
		log.Printf("[ERROR] [AdminWorkPackageDetailsHandler] Failed to load details: %v", err)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}
	log.Printf("[SUCCESS] [AdminWorkPackageDetailsHandler] Successfully retrieved details for Work Package ID: %d", wpID)
	data := map[string]interface{}{
		"Title":   "Work Package Execution Hub - WP-#" + strconv.Itoa(wpID),
		"Details": details,
	}

	renderAdminView(w, "admin_work_package_details.html", data)
}

// AdminSendCommunicationHandler handles POST /admin/work-packages/messages/add
func AdminSendCommunicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	wpID, _ := strconv.Atoi(r.FormValue("work_package_id"))
	senderRole := r.FormValue("sender_role") // Admin, Client, Supplier
	message := r.FormValue("message")

	if wpID > 0 && message != "" {
		if err := models.AddSecureMessage(wpID, senderRole, 1, message); err != nil {
			log.Printf("[ERROR] Failed to save secure message: %v", err)
		}
	}

	http.Redirect(w, r, "/admin/work-packages/details?id="+strconv.Itoa(wpID), http.StatusSeeOther)
}

// AdminPaymentLedgerHandler handles GET /admin/payments with project filtering
func AdminPaymentLedgerHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status") // Optional filter: "all", "active", "completed"

	if status == "" && (status != "all" || status != "active" || status != "completed") {
		status = "all"
	}
	// Parse optional project/work package filter query param
	packageIDStr := r.URL.Query().Get("id")
	selectedPackageID, _ := strconv.Atoi(packageIDStr)

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var payments []models.PaymentTrackerLog
	var err error

	// if selectedPackageID > 0 {
	// 	payments, err = models.GetPaymentsByWorkPackageID(selectedPackageID)
	// } else {
	payments, err = models.GetPaymentLedger()
	// }

	if err != nil {
		log.Printf("[ERROR] [AdminPaymentLedgerHandler] Failed to fetch payments: %v", err)
	}

	// Fetch all work packages for the dropdown filter and modal selector
	workPackage, _, err := models.GetPaginatedWorkPackages(page, status)

	if err != nil {
		log.Printf("[WARN] [AdminPaymentLedgerHandler] Failed to load work packages dropdown: %v", err)
	}

	data := map[string]interface{}{
		"Title":             "Double-Entry Payment Ledger",
		"Payments":          payments,
		"WorkPackages":      workPackage,
		"SelectedPackageID": selectedPackageID,
	}

	renderAdminView(w, "admin_payment.html", data)
}

// AdminAddPaymentHandler handles POST /admin/payments/add to store new ledger entries
func AdminAddPaymentHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [AdminAddPaymentHandler] Received %s request from %s", r.Method, r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Printf("[WARN] [AdminAddPaymentHandler] Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("[ERROR] [AdminAddPaymentHandler] Failed to parse form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	wpID, err := strconv.Atoi(r.FormValue("work_package_id"))
	if err != nil {
		log.Printf("[WARN] [AdminAddPaymentHandler] Invalid work_package_id '%s': %v", r.FormValue("work_package_id"), err)
	}

	clientAmt, err := strconv.ParseFloat(r.FormValue("client_amount_usd"), 64)
	if err != nil {
		log.Printf("[WARN] [AdminAddPaymentHandler] Invalid client_amount_usd '%s': %v", r.FormValue("client_amount_usd"), err)
	}
	clientStatus := r.FormValue("client_paid_status")
	clientDate := r.FormValue("client_paid_date")

	supplierAmt, err := strconv.ParseFloat(r.FormValue("supplier_amount_usd"), 64)
	if err != nil {
		log.Printf("[WARN] [AdminAddPaymentHandler] Invalid supplier_amount_usd '%s': %v", r.FormValue("supplier_amount_usd"), err)
	}
	supplierStatus := r.FormValue("supplier_paid_status")
	supplierDate := r.FormValue("supplier_paid_date")
	notes := r.FormValue("notes")

	if wpID > 0 {
		log.Printf("[INFO] [AdminAddPaymentHandler] Attempting to insert payment record for work_package_id=%d", wpID)
		err := models.InsertPaymentRecord(
			wpID, clientAmt, clientStatus, clientDate,
			supplierAmt, supplierStatus, supplierDate, notes,
		)
		if err != nil {
			log.Printf("[ERROR] [AdminAddPaymentHandler] Failed to insert payment for work_package_id=%d: %v", wpID, err)
		} else {
			log.Printf("[INFO] [AdminAddPaymentHandler] Successfully inserted payment record for work_package_id=%d", wpID)
		}
	} else {
		log.Printf("[WARN] [AdminAddPaymentHandler] Skipped database insertion due to invalid work_package_id: %d", wpID)
	}

	log.Printf("[INFO] [AdminAddPaymentHandler] Redirecting to /admin/payments")
	http.Redirect(w, r, "/admin/payments", http.StatusSeeOther)
}

// APISuppliersHandler returns all verified SME suppliers as JSON for multi-vendor payout splitting
func APISuppliersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page := 1
	filter := "all"
	search := ""

	// Decode parameters from POST JSON body
	if r.Body != nil && r.Body != http.NoBody {
		var payload struct {
			Page   int    `json:"page"`
			Filter string `json:"filter"`
			Search string `json:"search"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err == nil {
			if payload.Page > 0 {
				page = payload.Page
			}
			if payload.Filter != "" {
				filter = payload.Filter
			}
			search = payload.Search
		}
	}

	// Call the paginated & filtered model function
	suppliers, err := models.GetAllSuppliersPaginatedWithFilter(page, filter, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "failed to fetch suppliers",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"suppliers": suppliers,
		"page":      page,
		"filter":    filter,
		"search":    search,
	})

}

func APIWorkPackagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] [APIWorkPackagesHandler] Request received from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")

	status := r.URL.Query().Get("status") // Optional filter: "all", "active", "completed"

	workPackages, err := models.GetAllWorkPackagesForDropdown(status)
	if err != nil {
		log.Printf("[ERROR] [APIWorkPackagesHandler] Failed to fetch work packages: %v", err)
		http.Error(w, `{"error":"failed to fetch work packages"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[SUCCESS] [APIWorkPackagesHandler] Successfully encoding %d work packages to JSON", len(workPackages), workPackages)
	json.NewEncoder(w).Encode(workPackages)
}

// AdminUpdateRFQStatusHandler handles POST /admin/update-rfq-status
func AdminUpdateRFQStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Incoming request: Method=%s, RemoteAddr=%s, Path=%s", r.Method, r.RemoteAddr, r.URL.Path)

	// Step 1: Validate HTTP Method
	if r.Method != http.MethodPost {
		log.Printf("[WARN] [AdminUpdateRFQStatusHandler] Method Not Allowed: Received %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 2: Parse Form Values
	rfqIDStr := r.FormValue("rfq_id")
	status := r.FormValue("status")
	log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Form payload extracted: rfq_id='%s', status='%s'", rfqIDStr, status)

	// Step 3: Validate & Update RFQ Status in DB
	rfqID, err := strconv.ParseInt(rfqIDStr, 10, 64)
	if err != nil || rfqID <= 0 || status == "" {
		log.Printf("[ERROR] [AdminUpdateRFQStatusHandler] Invalid inputs: parsed_rfq_id=%d, parse_err=%v, status='%s'", rfqID, err, status)
	} else {
		log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Executing DB update for RFQ ID=%d with new status='%s'", rfqID, status)
		if updateErr := models.UpdateRFQStatus(int(rfqID), status); updateErr != nil {
			log.Printf("[ERROR] [AdminUpdateRFQStatusHandler] Failed to update RFQ status in DB: %v", updateErr)
		} else {
			log.Printf("[SUCCESS] [AdminUpdateRFQStatusHandler] Successfully updated RFQ ID=%d to status='%s'", rfqID, status)

			// Auto-create Work Package when RFQ is Approved
			if status == "accept" {
				log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Status is 'Approved'. Triggering Work Package creation for RFQ ID=%d...", rfqID)

				if wpErr := models.CreateWorkPackageFromRFQ(int(rfqID)); wpErr != nil {
					log.Printf("[ERROR] [AdminUpdateRFQStatusHandler] Failed to create Work Package from RFQ ID=%d: %v", rfqID, wpErr)
				} else {
					log.Printf("[SUCCESS] [AdminUpdateRFQStatusHandler] Work Package successfully created for RFQ ID=%d", rfqID)
				}
			}
		}
	}

	// Step 4: Retrieve updated list of RFQs
	log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Fetching updated list of all RFQs from DB...")
	rfqs, err := models.GetAllRFQs()
	if err != nil {
		log.Printf("[ERROR] [AdminUpdateRFQStatusHandler] Failed to fetch RFQs from DB: %v", err)
		rfqs = []models.RFQ{}
	} else {
		log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Successfully retrieved %d RFQs from DB", len(rfqs))
	}

	// Step 5: Extract Admin Session Data
	adminEmail := config.GetAdminEmail(r)
	log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Admin session identified: email='%s'", adminEmail)

	// Step 6: Construct View Data Map
	data := map[string]interface{}{
		"Title":      "RFQ Monitoring Desk",
		"AdminEmail": adminEmail,
		"Page":       "rfqs",
		"RFQs":       rfqs,
	}

	// Step 7: Render Admin View
	log.Printf("[INFO] [AdminUpdateRFQStatusHandler] Rendering view 'admin_rfqs.html'...")
	renderAdminView(w, "admin_rfqs.html", data)
	log.Printf("[SUCCESS] [AdminUpdateRFQStatusHandler] Completed processing request for RFQ ID=%s", rfqIDStr)
}

// AdminQualityAuditHandler renders the Multi-SME Quality Audit & Compliance Shield view
func AdminQualityAuditHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [AdminQualityAuditHandler] Incoming request: Method=%s, Path=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Fetch Multi-SME Quality Inspections
	inspections, err := models.GetQualityInspections()
	if err != nil {
		log.Printf("[ERROR] [AdminQualityAuditHandler] Failed to load inspections: %v", err)
		inspections = []models.QualityInspection{}
	}

	// 2. Fetch MTC Certificates
	mtcs, err := models.GetMTCCertificates()
	if err != nil {
		log.Printf("[ERROR] [AdminQualityAuditHandler] Failed to load MTCs: %v", err)
		mtcs = []models.MaterialCertificate{}
	}

	// 3. Fetch Flagged Communication Logs
	flaggedLogs, err := models.GetFlaggedCommunicationLogs()
	if err != nil {
		log.Printf("[ERROR] [AdminQualityAuditHandler] Failed to load flagged logs: %v", err)
		flaggedLogs = []models.FlaggedLog{}
	}

	adminEmail := config.GetAdminEmail(r)

	data := map[string]interface{}{
		"Title":              "Quality Audit & Compliance Shield",
		"AdminEmail":         adminEmail,
		"Page":               "quality",
		"QualityInspections": inspections,
		"MTCCertificates":    mtcs,
		"FlaggedLogs":        flaggedLogs,
	}

	log.Printf("[INFO] [AdminQualityAuditHandler] Rendering template 'admin_quality_audit.html'...")
	renderAdminView(w, "admin_quality_audit.html", data)
	log.Printf("[SUCCESS] [AdminQualityAuditHandler] Request completed successfully")
}

// HandleCreateAudit processes POST requests to log a new quality or compliance audit
func HandleCreateAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	wpID, _ := strconv.Atoi(r.FormValue("work_package_id"))
	supplierID, _ := strconv.Atoi(r.FormValue("supplier_id"))
	passed := r.FormValue("passed") == "1"

	inspectionType := r.FormValue("inspection_type")
	if inspectionType == "" {
		inspectionType = r.FormValue("inspection_type_select")
	}

	qi := models.QualityInspection{
		WorkPackageID:    wpID,
		SupplierID:       supplierID,
		InspectionType:   inspectionType,
		Passed:           passed,
		Comments:         r.FormValue("comments"),
		MeasuredValues:   r.FormValue("measured_values"),
		ToleranceAllowed: r.FormValue("tolerance_allowed"),
		ReportFilePath:   r.FormValue("report_file_path"),
	}

	_, err := models.CreateQualityInspection(qi)
	if err != nil {
		log.Printf("[ERROR] [HandleCreateAudit] Failed to save audit: %v", err)
		http.Error(w, "Database insertion error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/compliance", http.StatusSeeOther)
}

// HandleGetWorkPackagesAPI returns filtered work packages (current/completed) with search and pagination
func HandleGetWorkPackagesAPI(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status") // 'current' or 'completed'
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Fetch filtered packages from model layer
	packages, totalPages, err := models.GetFilteredWorkPackages(status, search, page, 10)
	if err != nil {
		log.Printf("[ERROR] [HandleGetWorkPackagesAPI] Query failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"work_packages": packages,
		"current_page":  page,
		"total_pages":   totalPages,
		"status":        status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAddAuditFormHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	adminEmail := config.GetAdminEmail(r)

	// 1. Fetch work packages for the dropdown selection
	// (Using GetFilteredWorkPackages with empty filters to retrieve all active records)
	wpResult, _, err := models.GetFilteredWorkPackages("", "", 1, 100)
	if err != nil {
		log.Printf("[ERROR] [GetAddAuditFormHandler] Failed to fetch work packages: %v", err)
	}

	// 2. Fetch suppliers list for the dropdown selection
	suppliers, err := models.GetAllSuppliers() // Ensure you have this DAO method or adapt as needed
	if err != nil {
		log.Printf("[ERROR] [GetAddAuditFormHandler] Failed to fetch suppliers: %v", err)
	}

	// 3. Populate template data map
	data := map[string]interface{}{
		"Title":        "Quality Audit & Compliance Shield Form",
		"AdminEmail":   adminEmail,
		"Page":         "audit form",
		"WorkPackages": wpResult.WorkPackages,
		"Suppliers":    suppliers,
	}

	log.Printf("[INFO] [AdminQualityAuditHandler] Rendering template 'admin_quality_audit_form.html'...")
	renderAdminView(w, "admin_quality_audit_form.html", data)

}

func RenderSupplierDetailsPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	supplierID, err := strconv.Atoi(idStr)
	if err != nil || supplierID <= 0 {
		log.Printf("[WARN] [RenderSupplierDetailsPage] Invalid or missing supplier ID: %s", idStr)
		http.Redirect(w, r, "/admin/sme", http.StatusSeeOther)
		return
	}
	adminEmail := config.GetAdminEmail(r)
	SupplierDetails, err := models.SupplierDetailsPage(supplierID)
	if err != nil {
		log.Printf("[ERROR] [RenderSupplierDetailsPage] Failed to fetch supplier details for ID %d: %v", supplierID, err)
		http.Redirect(w, r, "/admin/sme", http.StatusSeeOther)
		return
	}

	log.Printf("[INFO] [RenderSupplierDetailsPage] Successfully loaded details for supplier ID: %d", supplierID)

	data := map[string]interface{}{
		"Title":           "Supplier Details",
		"AdminEmail":      adminEmail,
		"SupplierDetails": SupplierDetails,
	}

	// Render using your project's template layout engine
	renderAdminView(w, "admin_supplier_details.html", data)
}

// RenderOnboardingPage renders the SME registration form view via GET
func RenderOnboardingPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "SME Manufacturer Onboarding Portal",
	}
	renderAdminView(w, "admin_sme_onboarding.html", data)
}

func HandleOnboardingSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("[ERROR] [HandleOnboardingSubmission] Failed to parse form: %v", err)
		http.Error(w, "Bad Request: Invalid Form Data", http.StatusBadRequest)
		return
	}

	// 1. Bind to Supplier Model
	supplier := models.Supplier{
		CompanyName:         r.FormValue("company_name"),
		ContactEmail:        r.FormValue("contact_email"),
		City:                r.FormValue("city"),
		State:               r.FormValue("state"),
		Phone:               r.FormValue("phone"),
		GSTNumber:           r.FormValue("gst_number"),
		MSMERegNumber:       r.FormValue("msme_reg_number"),
		ISOStandard:         r.FormValue("iso_standard"),
		PlantAreaSqFt:       parseInt(r.FormValue("plant_area_sq_ft")),
		WorkerCount:         parseInt(r.FormValue("worker_count")),
		CapabilitiesSummary: r.FormValue("capabilities_summary"),
		IsVerified:          false,
		NDASigned:           r.FormValue("nda_signed") == "true",
	}

	// 2. Bind to SupplierCapability Model
	capability := models.SupplierCapability{
		Category:         r.FormValue("category"),
		MachiningType:    r.FormValue("machining_type"),
		MaxTolerance:     r.FormValue("max_tolerance"),
		MaxComponentDim:  r.FormValue("max_component_dim"),
		MaterialsHandled: r.FormValue("materials_handled"),
		CapacityPerMonth: parseInt(r.FormValue("capacity_per_month")),
	}

	// 3. Bind to MachineStatus Slice Models
	names := r.Form["machine_name[]"]
	calibrations := r.Form["last_calibration[]"]
	services := r.Form["next_service_due[]"]
	notes := r.Form["notes[]"]

	var machines []models.MachineStatus
	for i := 0; i < len(names); i++ {
		if names[i] != "" {
			m := models.MachineStatus{
				MachineName:     names[i],
				IsOperational:   true,
				LastCalibration: getSafeSliceVal(calibrations, i),
				NextServiceDue:  getSafeSliceVal(services, i),
				Notes:           getSafeSliceVal(notes, i),
			}
			machines = append(machines, m)
		}
	}

	supplierID, err := models.CreateSupplier(supplier)

	if err != nil {
		log.Printf("[ERROR] [HandleOnboardingSubmission] Failed to create supplier: %v", err)
		data := map[string]interface{}{
			"Title": "SME Manufacturer Onboarding Portal",
			"Error": "Failed to save registration. Please check your inputs and try again.",
		}
		renderAdminView(w, "admin_sme_onboarding.html", data)
		return
	}
	capabilityId, err := models.CreateSupplierCapability(capability)

	if err != nil {
		log.Printf("[ERROR] [HandleOnboardingSubmission] Failed to create supplier capability: %v", err)
		data := map[string]interface{}{
			"Title": "SME Manufacturer Onboarding Portal",
			"Error": "Failed to save registration. Please check your inputs and try again.",
		}
		renderAdminView(w, "admin_sme_onboarding.html", data)
		return
	}
	log.Printf("[INFO] [HandleOnboardingSubmission] Successfully created supplier capability with ID: %d for supplier ID: %d", capabilityId, supplierID)
	for _, machine := range machines {
		machine.SupplierID = int(supplierID)
		machineId, err := models.CreateMachineStatus(machine)
		if err != nil {
			log.Printf("[ERROR] [HandleOnboardingSubmission] Failed to create machine status: %v", err)
			data := map[string]interface{}{
				"Title": "SME Manufacturer Onboarding Portal",
				"Error": "Failed to save registration. Please check your inputs and try again.",
			}
			renderAdminView(w, "admin_sme_onboarding.html", data)
			return
		}
		log.Printf("[INFO] [HandleOnboardingSubmission] Successfully created machine status with ID: %d for supplier ID: %d", machineId, supplierID)
	}

	log.Printf("[SUCCESS] [HandleOnboardingSubmission] Successfully registered supplier ID: %d , ", supplierID)
	http.Redirect(w, r, "/admin/sme?success=true", http.StatusSeeOther)
}

func parseInt(val string) int {
	res, _ := strconv.Atoi(val)
	return res
}

func getSafeSliceVal(arr []string, index int) string {
	if index < len(arr) {
		return arr[index]
	}
	return ""
}

// 1. RenderDetails handles static/general work package details or auxiliary views
func RenderDetails(w http.ResponseWriter, r *http.Request) {
	adminEmail := config.GetAdminEmail(r)

	data := map[string]interface{}{
		"Title":      "Work Package General Overview",
		"AdminEmail": adminEmail,
	}

	log.Printf("[INFO] [RenderDetails] Rendering general work package details view")
	renderAdminView(w, "admin_work_package_general.html", data)
}

// 2. RenderWorkPackageDetailsPage handles dynamic ID inspection (e.g. /admin/work-packages/203)
func RenderWorkPackageDetailsPage(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path (splits `/admin/work-packages/203` by '/')
	pathParts := strings.Split(strings.TrimSpace(r.URL.Path), "/")
	if len(pathParts) == 0 {
		log.Printf("[WARN] [RenderWorkPackageDetailsPage] Invalid URL path structure: %s", r.URL.Path)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	wpID, err := strconv.Atoi(idStr)
	if err != nil || wpID <= 0 {
		log.Printf("[WARN] [RenderWorkPackageDetailsPage] Invalid work package ID format: %s", idStr)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	adminEmail := config.GetAdminEmail(r)

	// Query database model for work package details and joined supplier info
	wpDetails, err := models.GetWorkPackageDetailsByID(wpID)
	if err != nil {
		log.Printf("[ERROR] [RenderWorkPackageDetailsPage] Work package not found for ID %d: %v", wpID, err)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	log.Printf("[INFO] [RenderWorkPackageDetailsPage] Successfully loaded inspection view for Work Package ID: %d", wpID)

	data := map[string]interface{}{
		"Title":      "Work Package #" + strconv.Itoa(wpID) + " - Audit & Inspection",
		"AdminEmail": adminEmail,
		"Package":    wpDetails,
	}

	renderAdminView(w, "admin_work_package_details.html", data)
}

func RenderDispatchPerProjectPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimSpace(r.URL.Path), "/")
	if len(parts) == 0 {
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	idStr := parts[len(parts)-1]
	wpID, err := strconv.Atoi(idStr)
	if err != nil || wpID <= 0 {
		log.Printf("[WARN] [RenderDispatchPage] Invalid dispatch work package ID: %s", idStr)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	dispatchData, err := models.GetDispatchDetailsByID(wpID)
	if err != nil {
		log.Printf("[ERROR] [RenderDispatchPage] Dispatch record not found for ID %d", wpID)
		http.Redirect(w, r, "/admin/dispatch", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title":    "Dispatch & Bidding Portal - WP #" + strconv.Itoa(wpID),
		"Dispatch": dispatchData,
	}

	renderAdminView(w, "admin_dispatch_details.html", data)
}

func HandleGetWorkPackageSuppliers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[WARN] [HandleGetWorkPackageSuppliers] Invalid HTTP method: %s from %s", r.Method, r.RemoteAddr)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	if idStr == "" {
		log.Printf("[WARN] [HandleGetWorkPackageSuppliers] Missing work package ID in request path from %s", r.RemoteAddr)
		http.Error(w, "Missing work package ID", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] [HandleGetWorkPackageSuppliers] Fetching assigned suppliers for WorkPackageID=%s", idStr)

	suppliers, err := models.GetWorkPackageSuppliers(idStr)
	if err != nil {
		// Return 400 Bad Request if the ID failed integer parsing, otherwise 500 Internal Server Error
		if strings.Contains(err.Error(), "invalid work package ID") {
			log.Printf("[WARN] [HandleGetWorkPackageSuppliers] Bad request for WorkPackageID=%s: %v", idStr, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("[ERROR] [HandleGetWorkPackageSuppliers] Database error fetching suppliers for WorkPackageID=%s: %v", idStr, err)
		http.Error(w, "Database error querying suppliers", http.StatusInternalServerError)
		return
	}

	// Ensure non-nil slice so empty results encode as [] instead of null
	if suppliers == nil {
		suppliers = []models.SupplierResponse{}
	}

	log.Printf("[INFO] [HandleGetWorkPackageSuppliers] Successfully retrieved %d supplier(s) for WorkPackageID=%s", len(suppliers), idStr)

	// Return JSON matching frontend expectations
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"suppliers": suppliers,
	}); err != nil {
		log.Printf("[ERROR] [HandleGetWorkPackageSuppliers] Failed to encode JSON response for WorkPackageID=%s: %v", idStr, err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : HandleUpdateSupplierStatus                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Manohar Patil (KUBERAIN INDUSTRY)                                  ║
║ 🎯 WHAT IT DOES   : Handles GET /admin/suppliers/status/{id}/{status} with detailed debug logs,  ║
║                     updates supplier status, fetches refreshed paginated list, and renders view. ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func HandleUpdateSupplierStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Incoming request: Method=%s Path=%s RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

	filter := "all"

	if r.Method != http.MethodGet {
		log.Printf("[WARN] [HandleUpdateSupplierStatus] Invalid HTTP method: %s from %s", r.Method, r.RemoteAddr)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	statusStr := r.PathValue("status")
	log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Extracted PathValues: id=%q, status=%q", idStr, statusStr)

	if idStr == "" || statusStr == "" {
		log.Printf("[WARN] [HandleUpdateSupplierStatus] Missing id or status path parameter from %s", r.RemoteAddr)
		http.Error(w, "Missing supplier ID or status", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] [HandleUpdateSupplierStatus] Calling models.UpdateSupplierStatus(id=%s, status=%s)", idStr, statusStr)
	err := models.UpdateSupplierStatus(idStr, statusStr)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[WARN] [HandleUpdateSupplierStatus] SupplierID=%s not found in DB", idStr)
			http.Error(w, "Supplier not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "invalid supplier") {
			log.Printf("[WARN] [HandleUpdateSupplierStatus] Validation failed for SupplierID=%s: %v", idStr, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("[ERROR] [HandleUpdateSupplierStatus] Database error updating SupplierID=%s: %v", idStr, err)
		http.Error(w, "Failed to update supplier status", http.StatusInternalServerError)
		return
	}

	log.Printf("[SUCCESS] [HandleUpdateSupplierStatus] Successfully updated SupplierID=%s to Status=%s", idStr, statusStr)

	search := ""
	log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Fetching refreshed paginated suppliers: filter=%q, search=%q, page=1, limit=10", filter, search)

	// Fetch Paginated & Filtered Data
	paginatedData, err := models.GetPaginatedSuppliers(filter, search, 1, 10)
	if err != nil {
		log.Printf("[ERROR] [HandleUpdateSupplierStatus] Failed to load paginated suppliers: %v", err)
		paginatedData = models.PaginatedResult{
			Suppliers:   []models.Supplier{},
			TotalPages:  1,
			CurrentPage: 1,
			Filter:      filter,
			Search:      search,
		}
	} else {
		log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Retrieved %d suppliers (CurrentPage=%d, TotalPages=%d)", len(paginatedData.Suppliers), paginatedData.CurrentPage, paginatedData.TotalPages)
	}

	// Construct View Data Context
	data := map[string]interface{}{
		"Title":      "SME Supplier Directory",
		"Page":       "sme",
		"Suppliers":  paginatedData.Suppliers,
		"Pagination": paginatedData,
		"Filter":     filter,
		"Search":     search,
	}

	log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Rendering template 'admin_sme.html' for SupplierID=%s", idStr)
	renderAdminView(w, "admin_sme.html", data)
	log.Printf("[DEBUG] [HandleUpdateSupplierStatus] Completed request for SupplierID=%s", idStr)
}
