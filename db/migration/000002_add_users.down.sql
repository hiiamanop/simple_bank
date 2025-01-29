ALTER TABLE "account" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE "account" DROP CONSTRAINT IF EXISTS "account_owner_fkey";
DROP TABLE IF EXISTS "users";