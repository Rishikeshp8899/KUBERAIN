package controllers

import (
	"encoding/json"
	"fmt"
	"go-mvc-app/models"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// PageData defines the view model passed to Go HTML templates
type PageData struct {
	Title       string
	CurrentPage string
	Submitted   bool
	RFQs        []models.RFQ
}

var customFuncs = template.FuncMap{
	"sub": func(a, b int) int {
		return a - b
	},
	"add": func(a, b int) int {
		return a + b
	},
}

// renderTemplate loads views/layouts/main.html + the targeted page view
func renderTemplate(w http.ResponseWriter, pageFile string, data PageData) {

	layoutPath := filepath.Join("views", "layouts", "main.html")
	pagePath := filepath.Join("views", "pages", pageFile)

	log.Printf("[INFO] Rendering template: Layout [%s] + Page [%s]", layoutPath, pagePath)
	// Use template.New("") with Funcs, then parse both layout and page files into the set
	tmpl, err := template.New("").Funcs(customFuncs).ParseFiles(layoutPath, pagePath)
	if err != nil {
		log.Printf("[ERROR] Template parsing failed (%s): %v", pageFile, err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the layout (or default template)
	// Note: Ensure your layout or execution targets the correct template name if needed (e.g., "main.html" or layout base name)
	layoutName := filepath.Base(layoutPath)
	if err := tmpl.ExecuteTemplate(w, layoutName, data); err != nil {
		log.Printf("[ERROR] Template execution failed (%s): %v", pageFile, err)
	}

}

// getClientIP extracts real client IP behind reverse proxies/load balancers
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}
	return ip
}

// ==========================================
// PUBLIC CONTROLLERS (Unprotected Endpoints)
// ==========================================

// HomePageHandler handles GET /
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:       "Home | B2B Manufacturing Aggregator",
		CurrentPage: "home",
	}
	renderTemplate(w, "home.html", data)
}

// ProductsPageHandler handles GET /products
func ProductsPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Industrial Products & Valves",
		CurrentPage: "products",
	}
	renderTemplate(w, "products.html", data)
}

// ServicesPageHandler handles GET /services
func ServicesPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Services & SME Network Capabilities",
		CurrentPage: "services",
	}
	renderTemplate(w, "service.html", data)
}

// RFQPageHandler handles GET /rfq (Displays form or submission confirmation)
func RFQPageHandler(w http.ResponseWriter, r *http.Request) {
	isSubmitted := r.URL.Query().Get("submitted") == "true"

	data := PageData{
		Title:       "Request Enterprise RFQ",
		CurrentPage: "rfq",
		Submitted:   isSubmitted,
	}
	renderTemplate(w, "rfq.html", data)
}

// SubmitRFQHandler handles multi-file uploads inside request-specific folders
func SubmitRFQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse Multipart Form (30 MB max memory limit)
	err := r.ParseMultipartForm(30 << 20)
	if err != nil {
		log.Printf("[ERROR] Multipart parse error: %v", err)
		http.Error(w, "Failed to parse file upload", http.StatusBadRequest)
		return
	}

	// 2. Generate unique RFQ Code and request-specific folder
	rfqCode := r.FormValue("rfq_code")
	if rfqCode == "" {
		rfqCode = fmt.Sprintf("RFQ-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	}

	// Folder structure: static/uploads/drawings/RFQ-20261020-1234/
	requestFolder := filepath.Join("static", "uploads", "drawings", rfqCode)
	err = os.MkdirAll(requestFolder, os.ModePerm)
	if err != nil {
		log.Printf("[ERROR] Failed to create folder [%s]: %v", requestFolder, err)
		http.Error(w, "Failed to initialize upload directory", http.StatusInternalServerError)
		return
	}

	// 3. Process Multiple Files from input name="drawings"
	var uploadedFilePaths []string
	files := r.MultipartForm.File["drawings"] // Matches <input name="drawings" multiple>

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("[WARNING] Could not open uploaded file [%s]: %v", fileHeader.Filename, err)
			continue
		}

		// Save file into requestFolder
		dstPath := filepath.Join(requestFolder, fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			file.Close()
			log.Printf("[WARNING] Could not create destination file [%s]: %v", dstPath, err)
			continue
		}

		_, err = io.Copy(dst, file)
		file.Close()
		dst.Close()

		if err == nil {
			// Store relative web-accessible path
			webPath := "/" + filepath.ToSlash(dstPath)
			uploadedFilePaths = append(uploadedFilePaths, webPath)
		}
	}

	// Combine multiple file paths into a comma-separated string for SQLite
	combinedDrawingPaths := strings.Join(uploadedFilePaths, ",")

	// 4. Extract standard form fields
	contactPerson := r.FormValue("contact_person")
	if contactPerson == "" {
		contactPerson = r.FormValue("full_name")
	}

	requirements := r.FormValue("requirements")
	if requirements == "" {
		requirements = r.FormValue("message")
	}

	targetPrice, _ := strconv.ParseFloat(r.FormValue("target_price_usd"), 64)

	rfq := models.RFQ{
		RFQCode:          rfqCode,
		CompanyName:      r.FormValue("company_name"),
		ContactPerson:    contactPerson,
		Email:            r.FormValue("email"),
		Phone:            r.FormValue("phone"),
		Category:         r.FormValue("category"),
		Material:         r.FormValue("material"),
		Quantity:         r.FormValue("quantity"),
		Requirements:     requirements,
		DrawingFilePath:  combinedDrawingPaths, // e.g. "/static/uploads/drawings/RFQ-xxx/file1.pdf,/static/uploads/drawings/RFQ-xxx/file2.dwg"
		TargetPriceUSD:   targetPrice,
		DeliveryDeadline: r.FormValue("delivery_deadline"),
		IncobookingTerms: r.FormValue("incobooking_terms"),
		Country:          r.FormValue("country"),
		City:             r.FormValue("city"),
		IPAddress:        getClientIP(r),
	}

	log.Printf("[DEBUG] Processing Multi-File RFQ [%s] | Files Saved: %d | Folder: %s", rfqCode, len(uploadedFilePaths), requestFolder)

	// 5. Persist record to SQLite
	_, err = models.CreateRFQ(rfq)
	if err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to process RFQ submission", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/rfq?submitted=true", http.StatusSeeOther)
}

// PublicSMEOnboardingHandler handles GET /suppliers/directory (Public Route)
func PublicSMEOnboarding(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("[INFO] [RenderSMERegistrationPage] Rendering supplier registration HTML page...")

	isSubmitted := r.URL.Query().Get("submitted") == "true"

	data := PageData{
		Title:       "SME Manufacturer Onboarding Portal",
		CurrentPage: "smes-onboarding",
		Submitted:   isSubmitted,
	}
	renderTemplate(w, "supplier_onboarding.html", data)

}

// SMESubmitHandler handles POST /suppliers/register (Processes form data & inserts into DB)
func SMESubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method Not Allowed",
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("[ERROR] [APISaveSMERegistration] Failed to parse form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Bad Request: Failed to parse form data",
		})
		return
	}

	// Extract numeric fields safely
	plantArea, _ := strconv.Atoi(r.FormValue("plant_area_sq_ft"))
	workerCount, _ := strconv.Atoi(r.FormValue("worker_count"))
	capacityPerMonth, _ := strconv.Atoi(r.FormValue("capacity_per_month"))
	ndaSigned := r.FormValue("nda_signed") == "true" || r.FormValue("nda_signed") == "on"

	// Construct Supplier struct
	supplier := models.Supplier{
		CompanyName:         r.FormValue("company_name"),
		City:                r.FormValue("city"),
		State:               r.FormValue("state"),
		ContactEmail:        r.FormValue("contact_email"),
		Phone:               r.FormValue("phone"),
		GSTNumber:           r.FormValue("gst_number"),
		MSMERegNumber:       r.FormValue("msme_reg_number"),
		ISOStandard:         r.FormValue("iso_standard"),
		PlantAreaSqFt:       plantArea,
		WorkerCount:         workerCount,
		CapabilitiesSummary: r.FormValue("capabilities_summary"),
		NDASigned:           ndaSigned,
	}

	// Construct SupplierCapability struct
	capability := models.SupplierCapability{
		Category:         r.FormValue("category"),
		MachiningType:    r.FormValue("machining_type"),
		MaxTolerance:     r.FormValue("max_tolerance"),
		MaxComponentDim:  r.FormValue("max_component_dim"),
		MaterialsHandled: r.FormValue("materials_handled"),
		CapacityPerMonth: capacityPerMonth,
	}

	// Extract Multiple Machines from Form Slices
	machineNames := r.Form["machine_name[]"]
	lastCalibrations := r.Form["last_calibration[]"]
	nextServiceDues := r.Form["next_service_due[]"]
	notesList := r.Form["notes[]"]

	var machines []models.MachineStatus
	for i := 0; i < len(machineNames); i++ {
		lastCal := ""
		if i < len(lastCalibrations) {
			lastCal = lastCalibrations[i]
		}
		nextServ := ""
		if i < len(nextServiceDues) {
			nextServ = nextServiceDues[i]
		}
		note := ""
		if i < len(notesList) {
			note = notesList[i]
		}

		machines = append(machines, models.MachineStatus{
			MachineName:     machineNames[i],
			LastCalibration: lastCal,
			NextServiceDue:  nextServ,
			Notes:           note,
		})
	}

	// Save to Database Atomically
	if err := models.RegisterNewSMESupplier(supplier, capability, machines); err != nil {
		log.Printf("[ERROR] [APISaveSMERegistration] Database registration failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database Insertion Failed: " + err.Error(),
		})
		return
	}

	// Return Success JSON Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Supplier registered and saved successfully",
	})

}

// AboutUsPageHandler handles GET /
func AboutUsPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:       "About Us",
		CurrentPage: "about",
	}
	renderTemplate(w, "about_us.html", data)
}

// CompliancePageHandler handles GET /
func ComplianceHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/compliance" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:       "Compliance",
		CurrentPage: "compliance",
	}
	renderTemplate(w, "compliance.html", data)
}

// SmeNetworkPageHandler handles GET /
func SmeNwtworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sme-network" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:       "Sme Network",
		CurrentPage: "sme-network",
	}
	renderTemplate(w, "sme-network.html", data)
}

func RenderDispatchPage(w http.ResponseWriter, r *http.Request) {
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
