/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : enterprise_client.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing enterprise client accounts, billing
                   metadata, credit limits, account statuses, and lifecycle queries.
====================================================================================================
*/

package models

import (
	"go-mvc-app/database"
)

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateEnterpriseClient                                                       ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new enterprise client profile record into the database.            ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding corporate names, domains, billing     ║
║                     cities, tax IDs, industry sectors, credit limits, and statuses, returning ID.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateEnterpriseClient(client EnterpriseClient) (int64, error) {
	query := `INSERT INTO enterprise_clients (name, domain, billing_city, tax_id, industry_sector, credit_limit_usd, account_status) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := database.DB.Exec(query, client.Name, client.Domain, client.BillingCity, client.TaxID, client.IndustrySector, client.CreditLimitUSD, client.AccountStatus)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetAllEnterpriseClients                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all enterprise client records sorted descending by primary key.    ║
║ 🏛️ ARCHITECTURE   : Broad listing query utilizing COALESCE fallback defaults for safe null handling  ║
║                     on optional domains, billing cities, tax IDs, and industry sectors.            ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetAllEnterpriseClients() ([]EnterpriseClient, error) {
	query := `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
                     COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
              FROM enterprise_clients ORDER BY id DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []EnterpriseClient
	for rows.Next() {
		var c EnterpriseClient
		if err := rows.Scan(&c.ID, &c.Name, &c.Domain, &c.BillingCity, &c.TaxID, &c.IndustrySector, &c.CreditLimitUSD, &c.AccountStatus, &c.CreatedAt); err != nil {
			continue
		}
		clients = append(clients, c)
	}
	return clients, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetEnterpriseClientByID                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Fetches a single enterprise client profile record by its primary key ID.     ║
║ 🏛️ ARCHITECTURE   : Targeted single-record point lookup query using COALESCE wrappers for safe   ║
║                     null handling on optional fields and metadata.                               ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetEnterpriseClientByID(id int) (*EnterpriseClient, error) {
	query := `SELECT id, name, COALESCE(domain, ''), COALESCE(billing_city, ''), COALESCE(tax_id, ''), 
                     COALESCE(industry_sector, ''), credit_limit_usd, account_status, created_at 
              FROM enterprise_clients WHERE id = ?`
	var c EnterpriseClient
	err := database.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Domain, &c.BillingCity, &c.TaxID, &c.IndustrySector, &c.CreditLimitUSD, &c.AccountStatus, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : UpdateEnterpriseClientStatus                                                 ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Updates the active account status string for a specified enterprise client.  ║
║ 🏛️ ARCHITECTURE   : Targeted SQL update statement for rapid status transitions without reloading ║
║                     full client entity graphs.                                                   ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func UpdateEnterpriseClientStatus(clientID int, status string) error {
	query := `UPDATE enterprise_clients SET account_status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, clientID)
	return err
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteEnterpriseClient                                                       ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently removes an enterprise client record from the database.           ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the enterprise clients table.  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteEnterpriseClient(clientID int) error {
	_, err := database.DB.Exec(`DELETE FROM enterprise_clients WHERE id = ?`, clientID)
	return err
}
