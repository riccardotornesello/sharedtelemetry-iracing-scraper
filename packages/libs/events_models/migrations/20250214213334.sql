-- Rename a column from "name" to "first_name"
ALTER TABLE "public"."competition_drivers" RENAME COLUMN "name" TO "first_name";
-- Modify "competition_drivers" table
ALTER TABLE "public"."competition_drivers" ADD COLUMN "last_name" text NULL;
