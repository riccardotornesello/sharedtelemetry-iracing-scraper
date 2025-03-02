-- Create "driver_stats" table
CREATE TABLE "public"."driver_stats" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "cust_id" bigint NOT NULL,
  "car_category" text NOT NULL,
  "license" text NOT NULL,
  "i_rating" bigint NOT NULL,
  PRIMARY KEY ("cust_id", "car_category")
);
-- Create index "idx_driver_stats_deleted_at" to table: "driver_stats"
CREATE INDEX "idx_driver_stats_deleted_at" ON "public"."driver_stats" ("deleted_at");
-- Create "drivers" table
CREATE TABLE "public"."drivers" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "cust_id" bigserial NOT NULL,
  "name" text NULL,
  "location" text NULL,
  PRIMARY KEY ("cust_id")
);
-- Create index "idx_drivers_deleted_at" to table: "drivers"
CREATE INDEX "idx_drivers_deleted_at" ON "public"."drivers" ("deleted_at");
