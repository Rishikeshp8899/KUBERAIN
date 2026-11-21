package controllers

import (
	"encoding/json"
	"fmt"
	"go-mvc-app/models"
	"html/template"
	"io"
	"log"
	"math/rand"
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
		now := time.Now()
		randomSuffix := rand.Intn(9000) + 1000
		rfqCode = fmt.Sprintf("RFQ-%s-%d", now.Format("20060102"), randomSuffix)
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
	tmpl, err := template.ParseFiles("views/pages/supplier_onboarding.html")

	if err != nil {
		http.Error(w, "Failed to load onboarding page: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
	}

}

// SMESubmitHandler handles POST /suppliers/register (Processes form data & inserts into DB)
func SMESubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("[INFO] [APISaveSMERegistration] Received request: Method=%s, RemoteAddr=%s", r.Method, r.RemoteAddr)

	// Step 1: Validate HTTP Method
	if r.Method != http.MethodPost {
		log.Printf("[WARN] [APISaveSMERegistration] Invalid method received: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method Not Allowed",
		})
		return
	}
	log.Printf("[DEBUG] [APISaveSMERegistration] Method validation passed (POST).")

	// Step 2: Parse Form Data (Multipart vs URL-encoded)
	contentType := r.Header.Get("Content-Type")
	log.Printf("[DEBUG] [APISaveSMERegistration] Parsing incoming form data with Content-Type: %s", contentType)

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("[ERROR] [APISaveSMERegistration] Failed to parse multipart form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Bad Request: Failed to parse multipart form data",
			})
			return
		}
		log.Printf("[DEBUG] [APISaveSMERegistration] Successfully parsed multipart form data.")
	} else {
		if err := r.ParseForm(); err != nil {
			log.Printf("[ERROR] [APISaveSMERegistration] Failed to parse form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Bad Request: Failed to parse form data",
			})
			return
		}
		log.Printf("[DEBUG] [APISaveSMERegistration] Successfully parsed standard URL-encoded form data.")
	}

	// Step 3: Extract and convert numeric/boolean fields safely
	log.Printf("[DEBUG] [APISaveSMERegistration] Extracting numeric and boolean form values...")
	plantArea, errPlant := strconv.Atoi(r.FormValue("plant_area_sq_ft"))
	if errPlant != nil {
		log.Printf("[WARN] [APISaveSMERegistration] Could not parse plant_area_sq_ft ('%s'): %v. Defaulting to 0.", r.FormValue("plant_area_sq_ft"), errPlant)
	}

	workerCount, errWorker := strconv.Atoi(r.FormValue("worker_count"))
	if errWorker != nil {
		log.Printf("[WARN] [APISaveSMERegistration] Could not parse worker_count ('%s'): %v. Defaulting to 0.", r.FormValue("worker_count"), errWorker)
	}

	capacityPerMonth, errCap := strconv.Atoi(r.FormValue("capacity_per_month"))
	if errCap != nil {
		log.Printf("[WARN] [APISaveSMERegistration] Could not parse capacity_per_month ('%s'): %v. Defaulting to 0.", r.FormValue("capacity_per_month"), errCap)
	}

	ndaRaw := r.FormValue("nda_signed")
	ndaSigned := ndaRaw == "true" || ndaRaw == "on"
	log.Printf("[DEBUG] [APISaveSMERegistration] Parsed numeric/boolean fields -> PlantArea: %d, WorkerCount: %d, CapacityPerMonth: %d, NDASigned: %v", plantArea, workerCount, capacityPerMonth, ndaSigned)

	// Step 4: Construct Supplier struct
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
	log.Printf("[DEBUG] [APISaveSMERegistration] Constructed Supplier object: CompanyName=%s, City=%s, Email=%s", supplier.CompanyName, supplier.City, supplier.ContactEmail)

	// Step 5: Construct SupplierCapability struct
	capability := models.SupplierCapability{
		Category:         r.FormValue("category"),
		MachiningType:    r.FormValue("machining_type"),
		MaxTolerance:     r.FormValue("max_tolerance"),
		MaxComponentDim:  r.FormValue("max_component_dim"),
		MaterialsHandled: r.FormValue("materials_handled"),
		CapacityPerMonth: capacityPerMonth,
	}
	log.Printf("[DEBUG] [APISaveSMERegistration] Constructed Capability object: Category=%s, MachiningType=%s", capability.Category, capability.MachiningType)
	log.Printf("[DEBUG] [APISaveSMERegistration] Parsed Supplier: %+v | Capability: %+v", supplier, capability)

	// Step 6: Extract Multiple Machines from Form Slices
	machineNames := r.Form["machine_name[]"]
	lastCalibrations := r.Form["last_calibration[]"]
	nextServiceDues := r.Form["next_service_due[]"]
	notesList := r.Form["notes[]"]
	log.Printf("[DEBUG] [APISaveSMERegistration] Extracted form arrays -> machine_name count: %d, last_calibration count: %d, next_service_due count: %d, notes count: %d", len(machineNames), len(lastCalibrations), len(nextServiceDues), len(notesList))

	// Step 7: Map and build MachineStatus slice
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

		machine := models.MachineStatus{
			MachineName:     machineNames[i],
			LastCalibration: lastCal,
			NextServiceDue:  nextServ,
			Notes:           note,
		}
		machines = append(machines, machine)
		log.Printf("[DEBUG] [APISaveSMERegistration] Processed Machine #%d: Name=%s, LastCal=%s, NextDue=%s", i+1, machine.MachineName, machine.LastCalibration, machine.NextServiceDue)
	}

	log.Printf("[DEBUG] [APISaveSMERegistration] Final Parsed Machines list (total %d): %+v", len(machines), machines)

	// Step 8: Save to Database Atomically
	log.Printf("[INFO] [APISaveSMERegistration] Calling database model RegisterNewSMESupplier...")
	if err := models.RegisterNewSMESupplier(supplier, capability, machines); err != nil {
		log.Printf("[ERROR] [APISaveSMERegistration] Database registration failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database Insertion Failed: " + err.Error(),
		})
		return
	}
	log.Printf("[SUCCESS] [APISaveSMERegistration] Database registration completed successfully.")

	// Step 9: Return Success JSON Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Supplier registered and saved successfully",
	})
	log.Printf("[INFO] [APISaveSMERegistration] Successfully sent 200 OK JSON response to client.")
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
		http.Redirect(w, r, "/admin/work/packages", http.StatusSeeOther)
		return
	}

	idStr := parts[len(parts)-1]
	wpID, err := strconv.Atoi(idStr)
	if err != nil || wpID <= 0 {
		log.Printf("[WARN] [RenderDispatchPage] Invalid dispatch work package ID: %s", idStr)
		http.Redirect(w, r, "/admin/work/packages", http.StatusSeeOther)
		return
	}

	dispatchData, err := models.GetDispatchDetailsByID(wpID)
	if err != nil {
		log.Printf("[ERROR] [RenderDispatchPage] Dispatch record not found for ID %d", wpID)
		http.Redirect(w, r, "/admin/work/packages", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title":    "Dispatch & Bidding Portal - WP #" + strconv.Itoa(wpID),
		"Dispatch": dispatchData,
	}

	renderAdminView(w, "admin_dispatch_details.html", data)
}
