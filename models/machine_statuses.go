/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : machine_status.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing shop-floor machine health, calibration
                   schedules, operational status toggles, and public audit listings.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
	"log"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateMachineStatus                                                          ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new shop-floor machine status and calibration record into SQLite.   ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding supplier IDs, machine names,          ║
║                     operational booleans, calibration dates, service due dates, and notes.       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateMachineStatus(m MachineStatus) (int64, error) {
	query := `INSERT INTO machine_statuses (supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes) 
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, m.SupplierID, m.MachineName, m.IsOperational, m.LastCalibration, m.NextServiceDue, m.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetMachinesBySupplierID                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all machine status records associated with a specific supplier ID.   ║
║ 🏛️ ARCHITECTURE   : Filtered relational query utilizing COALESCE wrappers for safe null handling  ║
║                     on calibration dates, service due schedules, and operational notes.          ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetMachinesBySupplierID(supplierID int) ([]MachineStatus, error) {
	query := `SELECT id, supplier_id, machine_name, is_operational, COALESCE(last_calibration, ''), COALESCE(next_service_due, ''), COALESCE(notes, '') 
              FROM machine_statuses WHERE supplier_id = ?`
	rows, err := database.DB.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []MachineStatus
	for rows.Next() {
		var m MachineStatus
		if err := rows.Scan(&m.ID, &m.SupplierID, &m.MachineName, &m.IsOperational, &m.LastCalibration, &m.NextServiceDue, &m.Notes); err != nil {
			continue
		}
		machines = append(machines, m)
	}
	return machines, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpdateMachineOperationalStatus                                               ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Updates the active operational boolean state for a specified machine ID.     ║
║ 🏛️ ARCHITECTURE   : Targeted SQL update statement for rapid machine health state transitions     ║
║                     without requiring full entity payload reloads.                               ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpdateMachineOperationalStatus(machineID int, isOperational bool) error {
	_, err := database.DB.Exec(`UPDATE machine_statuses SET is_operational = ? WHERE id = ?`, isOperational, machineID)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteMachineStatus                                                          ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a machine status record from the database by ID.         ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the machine statuses table.    ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteMachineStatus(id int) error {
	_, err := database.DB.Exec(`DELETE FROM machine_statuses WHERE id = ?`, id)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetPublicMachineStatuses                                                     ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves operational machine health records specifically for public audit.  ║
║ 🏛️ ARCHITECTURE   : Filtered query restricting results strictly to active machines (`is_operational = 1`)║
║                     to display verified live shop-floor capacity to clients.                     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetPublicMachineStatuses() ([]MachineStatus, error) {
	var machines []MachineStatus

	query := `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
              FROM machine_statuses 
              WHERE is_operational = 1`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetPublicMachineStatuses] Query execution failed: %v", err)
		return machines, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MachineStatus
		err := rows.Scan(
			&ms.ID, &ms.SupplierID, &ms.MachineName, &ms.IsOperational,
			&ms.LastCalibration, &ms.NextServiceDue, &ms.Notes,
		)
		if err != nil {
			log.Printf("[WARN] [GetPublicMachineStatuses] Row scan error: %v", err)
			continue
		}
		machines = append(machines, ms)
	}

	return machines, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllMachineStatuses                                                        ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all shop-floor machine health and calibration records.             ║
║ 🏛️ ARCHITECTURE   : Comprehensive administrative query sorted prioritizing operational machines first ║
║                     (`is_operational DESC`) and earliest upcoming service due dates (`ASC`).     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllMachineStatuses() ([]MachineStatus, error) {
	var machines []MachineStatus

	query := `SELECT id, supplier_id, machine_name, is_operational, last_calibration, next_service_due, notes 
              FROM machine_statuses 
              ORDER BY is_operational DESC, next_service_due ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("[ERROR] [GetAllMachineStatuses] Query execution failed: %v", err)
		return machines, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MachineStatus
		err := rows.Scan(
			&ms.ID, &ms.SupplierID, &ms.MachineName, &ms.IsOperational,
			&ms.LastCalibration, &ms.NextServiceDue, &ms.Notes,
		)
		if err != nil {
			log.Printf("[WARN] [GetAllMachineStatuses] Row scan error: %v", err)
			continue
		}
		machines = append(machines, ms)
	}

	return machines, nil
}
