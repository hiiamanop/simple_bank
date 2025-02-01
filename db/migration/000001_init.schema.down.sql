-- Drop foreign key constraints first
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE "entries" DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";
ALTER TABLE "account" DROP CONSTRAINT IF EXISTS "account_owner_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "transfers_from_account_id_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_idx";
DROP INDEX IF EXISTS "entries_account_id_idx";
DROP INDEX IF EXISTS "account_owner_currency_idx";
DROP INDEX IF EXISTS "account_owner_idx";

-- Drop tables in reverse order
DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "account";
DROP TABLE IF EXISTS "users";

-- Drop custom types
DROP TYPE IF EXISTS "currency";