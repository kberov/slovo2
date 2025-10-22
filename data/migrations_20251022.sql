-- 202510222301 up
PRAGMA foreign_keys=OFF;
-- No NULL values in existing table.
-- Set default values for existing NULLs.
UPDATE domove SET description='' WHERE description IS NULL;
UPDATE domove SET owner_id=0 WHERE owner_id IS NULL;
UPDATE domove SET group_id=0 WHERE group_id IS NULL;
UPDATE domove SET aliases='' WHERE aliases IS NULL;
UPDATE domove SET templates='' WHERE templates IS NULL;
-- Recreate domove to alter some columns to not be NULL by default.
CREATE TABLE domove_new (
-- domove is the plural form of 'dom' in Bulgarian, meaning 'home'.
-- The similarity with domains is not a coincidence
--  'ID, referenced by stranici, that belong to this domain.'
  id INTEGER PRIMARY KEY AUTOINCREMENT,
--  'Domain name as in $ENV{HTTP_HOST}.'
  domain VARCHAR(63) UNIQUE NOT NULL,
--  'The name of this site.'
  site_name VARCHAR(63) NOT NULL,
--  'Site description'
  description VARCHAR(2000) NOT NULL DEFAULT '',
--   'User for which the permissions apply (owner).'
  owner_id INTEGER NOT NULL DEFAULT 0 REFERENCES users(id),
--  'Group for which the permissions apply.'
  group_id INTEGER NOT NULL DEFAULT 0 REFERENCES groups(id),
--  'Domain permissions like in a unix filesystem, eg -rwxr-xr-x.'
  permissions VARCHAR(10) DEFAULT '-rwxr-xr-x' ,
--  '0:not published, 1:for review, >=2:published'
  published INT(1) DEFAULT 0,
  -- IPs from which this domain may be served, eg localhost can be on '127.0.0.1,127.0.1.1'
  ips VARCHAR DEFAULT '127.0.0.1,127.0.1.1',
  aliases VARCHAR(2000) NOT NULL DEFAULT '',
  templates VARCHAR(255) NOT NULL DEFAULT ''
);
INSERT INTO domove_new (id, domain, site_name, description, owner_id,
    group_id, permissions, published)
    SELECT id, domain, site_name, description, owner_id,
    group_id, permissions, published FROM domove;
DROP TABLE domove;
DROP INDEX IF EXISTS domove_published;
DROP INDEX IF EXISTS domove_ips;
PRAGMA foreign_key_checks;
ALTER TABLE domove_new RENAME TO domove;
CREATE INDEX domove_published ON domove(published);
CREATE INDEX domove_ips ON domove(ips);
PRAGMA foreign_keys=ON;

-- 202510222301 down
