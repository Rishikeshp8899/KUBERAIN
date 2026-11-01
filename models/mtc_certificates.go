/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : mtc_certificate.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing MTC certificates, relational supplier
                   joins, heat number traceability, and chronological compliance audits.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
	"log"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetMTCCertificates                                                           ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches Material Test Certificates joined with corresponding supplier names. ║
║ 🏛️ ARCHITECTURE   : Relational query joining MTC certificates and suppliers tables, utilizing COALESCE║
║                     fallbacks for safe null handling and sorting descending by creation timestamps.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
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
