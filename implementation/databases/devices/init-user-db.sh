#bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO \$$
    BEGIN
       IF NOT EXISTS (
          SELECT FROM pg_catalog.pg_roles  -- SELECT list can be empty for this
          WHERE  rolname = 'angelos') THEN
          CREATE ROLE angelos LOGIN PASSWORD 'example';
       END IF;
    END
    \$$;
    SELECT 'CREATE DATABASE devices_db'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'devices_db')\gexec
    
    GRANT ALL PRIVILEGES ON DATABASE devices_db TO angelos;
    
    \c devices_db
EOSQL