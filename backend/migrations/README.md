# Migrations

Place ordered SQL migrations here as the enterprise schema is implemented.

Start from the reference design in `../../docs/database-schema.sql`, then split it into versioned migrations such as:

```text
000001_create_identity_tables.up.sql
000001_create_identity_tables.down.sql
000002_create_machine_inventory.up.sql
000002_create_machine_inventory.down.sql
```
