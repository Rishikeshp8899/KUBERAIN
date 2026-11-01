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

	// Public Web Pages
	mux.HandleFunc("/", controllers.HomePageHandler)
	mux.HandleFunc("/products", controllers.ProductsPageHandler)
	mux.HandleFunc("/services", controllers.ServicesPageHandler)
	mux.HandleFunc("/rfq", controllers.RFQPageHandler)
	mux.HandleFunc("/submit-rfq", controllers.SubmitRFQHandler)
	mux.HandleFunc("/smes-onboarding", controllers.PublicSMEOnboarding)
	mux.HandleFunc("/about", controllers.AboutUsPageHandler)
	mux.HandleFunc("/compliance", controllers.ComplianceHandler)
	mux.HandleFunc("/sme-network", controllers.SmeNwtworkHandler)
	mux.HandleFunc("/api/suppliers/register", controllers.SMESubmitHandler)

	// Admin Authentication
	mux.HandleFunc("/admin/login", controllers.AdminLoginHandler)
	mux.HandleFunc("/admin/logout", controllers.AdminLogoutHandler)

	// Admin Protected Operations
	mux.HandleFunc("/admin/rfqs", middleware.RequireAdminAuth(controllers.AdminRFQsHandler))
	mux.HandleFunc("/admin/rfqs/update-rfq-status", middleware.RequireAdminAuth(controllers.AdminUpdateRFQStatusHandler))
	mux.HandleFunc("/admin/sme", middleware.RequireAdminAuth(controllers.AdminSMEHandler))
	mux.HandleFunc("/admin/add-sme", middleware.RequireAdminAuth(controllers.AdminAddSMEHandler))
	mux.HandleFunc("/admin/rfqs/update-sme-status", middleware.RequireAdminAuth(controllers.AdminUpdateSMEStatusHandler))
	mux.HandleFunc("/admin/dispatch", middleware.RequireAdminAuth(controllers.AdminDispatchHandler))
	mux.HandleFunc("/admin/payments", middleware.RequireAdminAuth(controllers.AdminPaymentsHandler))
	mux.HandleFunc("/admin/payments/add", middleware.RequireAdminAuth(controllers.AdminAddPaymentHandler))
	mux.HandleFunc("/admin/update-payment-status", middleware.RequireAdminAuth(controllers.AdminUpdatePaymentStatusHandler))
	mux.HandleFunc("/admin/compliance", middleware.RequireAdminAuth(controllers.AdminComplianceHandler))
	mux.HandleFunc("/admin/rfqs/update-status", middleware.RequireAdminAuth(controllers.AdminRFQDetailAndMatchHandler))
	mux.HandleFunc("/admin/work-packages/details", middleware.RequireAdminAuth(controllers.AdminWorkPackageDetailsHandler))
	mux.HandleFunc("/admin/work-packages/messages/add", middleware.RequireAdminAuth(controllers.AdminSendCommunicationHandler))
	mux.HandleFunc("/admin/work-packages", middleware.RequireAdminAuth(controllers.APIWorkPackagesHandler))
	mux.HandleFunc("/admin/quality-audit", middleware.RequireAdminAuth(controllers.AdminQualityAuditHandler))
	mux.HandleFunc("/admin/audit/create", middleware.RequireAdminAuth(controllers.HandleCreateAudit))
	mux.HandleFunc("/admin/audit-form/add", middleware.RequireAdminAuth(controllers.GetAddAuditFormHandler))

	// Static Assets
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// View Assets
	vs := http.FileServer(http.Dir("./views"))
	mux.Handle("/views/", http.StripPrefix("/views/", vs))

	port := ":9091"
	fmt.Printf("Richim Valve Admin Portal running at http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
