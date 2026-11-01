/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : supplier_capability.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing SME supplier manufacturing capabilities,
                   machining types, tolerance limits, material proficiencies, and monthly output volumes.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateSupplierCapability                                                     ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new supplier manufacturing capability record into SQLite.          ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding supplier IDs, categories, machining     ║
║                     types, maximum tolerances, component dimensions, materials handled, and capacity.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateSupplierCapability(cap SupplierCapability) (int64, error) {
	query := `INSERT INTO supplier_capabilities (supplier_id, category, machining_type, max_tolerance, max_component_dim, materials_handled, capacity_per_month) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, cap.SupplierID, cap.Category, cap.MachiningType, cap.MaxTolerance, cap.MaxComponentDim, cap.MaterialsHandled, cap.CapacityPerMonth)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetCapabilitiesBySupplierID                                                  ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all capability records associated with a specific supplier ID.       ║
║ 🏛️ ARCHITECTURE   : Filtered relational query utilizing COALESCE wrappers for safe null handling  ║
║                     on maximum tolerances, component dimensions, and handled materials.          ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetCapabilitiesBySupplierID(supplierID int) ([]SupplierCapability, error) {
	query := `SELECT id, supplier_id, category, machining_type, COALESCE(max_tolerance, ''), COALESCE(max_component_dim, ''), COALESCE(materials_handled, ''), capacity_per_month 
              FROM supplier_capabilities WHERE supplier_id = ?`
	rows, err := database.DB.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var caps []SupplierCapability
	for rows.Next() {
		var c SupplierCapability
		if err := rows.Scan(&c.ID, &c.SupplierID, &c.Category, &c.MachiningType, &c.MaxTolerance, &c.MaxComponentDim, &c.MaterialsHandled, &c.CapacityPerMonth); err != nil {
			continue
		}
		caps = append(caps, c)
	}
	return caps, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteCapability                                                             ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a supplier capability record by its primary key ID.      ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the supplier capabilities table.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteCapability(id int) error {
	_, err := database.DB.Exec(`DELETE FROM supplier_capabilities WHERE id = ?`, id)
	return err
}
