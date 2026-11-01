/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : material_certificate.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing Material Test Certificates (MTC),
                   quality compliance verification records, and heat number traceability.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateMaterialCertificate                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new Material Test Certificate (MTC) compliance record.             ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding work package IDs, heat numbers,       ║
║                     material grades, mill names, and file paths, returning the generated ID.     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateMaterialCertificate(cert MaterialCertificate) (int64, error) {
	query := `INSERT INTO material_certificates (work_package_id, heat_number, material_grade, mill_name, cert_file_path, verified_by) 
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, cert.WorkPackageID, cert.HeatNumber, cert.MaterialGrade, cert.MillName, cert.CertFilePath, cert.VerifiedBy)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetCertificatesByWorkPackageID                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches all MTC compliance records filtered by a specific work package ID.     ║
║ 🏛️ ARCHITECTURE   : Targeted relational query utilizing COALESCE wrappers for safe null-handling   ║
║                     on optional mill names and verification audits.                              ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
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

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteMaterialCertificate                                                    ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a material test certificate entry by its primary key.    ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the material certificates table. ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteMaterialCertificate(id int) error {
	_, err := database.DB.Exec(`DELETE FROM material_certificates WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetMaterialCertificates                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches all Material Test Certificates uploaded across the entire system.     ║
║ 🏛️ ARCHITECTURE   : Broad listing query utilizing COALESCE defaults and descending ID sorting    ║
║                     to display the most recently uploaded compliance documents first.            ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
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
