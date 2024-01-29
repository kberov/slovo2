package model

// SQLMap is a map of name/query. Each entry has a name and an SQL query used
// in some method.
type SQLMap map[string]any

var queryTemplates = SQLMap{
	"INSERT":  `INSERT INTO ${table} (${columns}) VALUES(${placeholders})`,
	"GetByID": `SELECT * FROM ${table} WHERE id=?`,
	"SELECT":  `SELECT ${columns} FROM ${table} ${WHERE} LIMIT ${limit} OFFSET ${offset}`,
	// 0. WHERE alias = $alias failed to match
	// 1. Suppose the user stumbled on a link with the old most recent alias.
	// 2. Look for the most recent id which had such $alias in the given table.
	// To be embedded in other queries
	"ALIAS_IS": `(SELECT new_alias FROM aliases
    	WHERE alias_table='${table}' AND (old_alias=?) ORDER BY ID DESC LIMIT 1)
 	 OR ${table}.id=(SELECT alias_id FROM aliases
    	WHERE alias_table='${table}' AND (old_alias=? OR new_alias=?) ORDER BY ID DESC LIMIT 1)`,
	// To be embedded in other queries
	"PUBLISHED_DOMAIN_ID_BY_NAME_IS": `(SELECT id FROM domove
		WHERE (? LIKE '%' || domain OR aliases LIKE ? OR ips LIKE ?) AND published = 2 LIMIT 1)`,
	// To be embedded in other queries
	"PERMISSIONS_ARE": `(
	-- others can read and execute AND stranica/celina is published
	( ${table}.permissions LIKE '%r_x' AND ${table}.published = 2 ) 
	-- owner can read and execute
	OR ( ${table}.permissions LIKE '_r_x%' AND ${table}.user_id = ? )
	-- user has the group group_id and the page is rx by this group_id and is published(0|1|2)
	OR (
		( ${table}.group_id IN (SELECT group_id from user_group WHERE user_id=?)
			AND ${table}.permissions LIKE '____r_x%' AND ${table}.published = ? 
		)
	)
)`,
	"GET_PAGE_FOR_DISPLAY": `SELECT ${table}.*, c.title, c.body, c.language, c.data_type, c.data_format
		FROM ${table}
		JOIN celini AS c ON (
    		stranici.id = c.page_id AND c.pid=0 AND c.permissions LIKE 'd%'
    		AND c.language=? AND c.data_type='title'
    		AND c.published = stranici.published)
		WHERE (
		${PERMISSIONS_ARE}
		AND ( ${table}.alias = ? OR ${table}.alias  = ${ALIAS_IS})
		AND ${table}.dom_id = ${PUBLISHED_DOMAIN_ID_BY_NAME_IS}
		AND ${table}.hidden = 0 
		${AND_FOR_DISPLAY}
		) LIMIT 1`,
	"GET_CELINA_FOR_DISPLAY": `SELECT * from ${table} WHERE (
		page_id=? AND language LIKE ? AND box=?
		AND ${PERMISSIONS_ARE}
		AND ( ${table}.alias = ? OR ${table}.alias  = ${ALIAS_IS})
		AND ${table}.bad=0 
		${AND_FOR_DISPLAY}
		)
		`,
	"CELINI_FOR_DISPLAY_IN_PAGE": `
	SELECT ${columns} FROM ${table} WHERE (
		page_id = ?	AND pid = ? and box= 'main'
		-- find exact language or at least first part, e.g. (bg-)
		AND (language LIKE ?)
		AND ${PERMISSIONS_ARE}
		${AND_FOR_DISPLAY}
	) 
	ORDER BY featured DESC, id DESC, sorting ASC
	LIMIT ? OFFSET ?
		`,
	"AND_FOR_DISPLAY": `
		AND ${table}.deleted = 0
		AND ( ${table}.start = 0 OR ${table}.start < ? )
		AND ( ${table}.stop = 0 OR ${table}.stop > ? )
		`,
}
