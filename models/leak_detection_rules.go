/*
====================================================================================================
🏢 ORGANIZATION  : KUBERAIN INDUSTRY
👨‍💻 DEVELOPER     : Rishikesh Patil
📅 DATE          : November 1, 2026
📁 FILE NAME     : leak_detection_rule.go
📦 PACKAGE       : models
📝 FILE OVERVIEW : Core Data Access Object (DAO) layer managing disintermediation leak detection rules,
                   active regex inspection filters, severity ratings, and rule lifecycle configurations.
====================================================================================================
*/

package models

import "go-mvc-app/database"

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : CreateLeakDetectionRule                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Inserts a new leak detection rule record into the database.                  ║
║ 🏛️ ARCHITECTURE   : Parameterized SQL insert command binding rule names, regular expression      ║
║                     patterns, severity levels, and active status flags, returning the generated ID.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func CreateLeakDetectionRule(rule LeakDetectionRule) (int64, error) {
	query := `INSERT INTO leak_detection_rules (rule_name, regex_pattern, severity, is_active) 
              VALUES (?, ?, ?, ?)`
	res, err := database.DB.Exec(query, rule.RuleName, rule.RegexPattern, rule.Severity, rule.IsActive)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : GetActiveLeakRules                                                           ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Retrieves all active leak detection rules currently enforcing chat filters.  ║
║ 🏛️ ARCHITECTURE   : Filtered point-query selecting rules where `is_active = 1` to load live sandbox ║
║                     security validation parameters efficiently into memory.                        ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func GetActiveLeakRules() ([]LeakDetectionRule, error) {
	query := `SELECT id, rule_name, regex_pattern, severity, is_active FROM leak_detection_rules WHERE is_active = 1`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []LeakDetectionRule
	for rows.Next() {
		var r LeakDetectionRule
		if err := rows.Scan(&r.ID, &r.RuleName, &r.RegexPattern, &r.Severity, &r.IsActive); err != nil {
			continue
		}
		rules = append(rules, r)
	}
	return rules, nil
}

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 🛠️ FUNCTION       : DeleteLeakDetectionRule                                                      ║
║ 👨‍💻 DEVELOPER      : Rishikesh Patil (KUBERAIN INDUSTRY)                                          ║
║ 🎯 WHAT IT DOES   : Permanently deletes a leak detection rule record by its primary key ID.      ║
║ 🏛️ ARCHITECTURE   : Direct hard-deletion command executed against the leak detection rules table.║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/
func DeleteLeakDetectionRule(id int) error {
	_, err := database.DB.Exec(`DELETE FROM leak_detection_rules WHERE id = ?`, id)
	return err
}
