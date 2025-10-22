CREATE VIEW view_stranici
    AS
    SELECT stranici.*, c.title, c.body, c.language, c.data_type, c.data_format,
      (SELECT domain FROM domove WHERE (stranici.dom_id=id AND published = 2) LIMIT 1) AS domain
	FROM stranici
    JOIN celini AS c ON (
        stranici.id = c.page_id AND c.pid=0 AND c.permissions LIKE 'd%'
	    AND c.data_type='title'
        AND c.published = stranici.published)

