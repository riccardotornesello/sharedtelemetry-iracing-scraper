-- Create "brands" table
CREATE TABLE "public"."brands" (
  "name" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "icon" text NULL,
  PRIMARY KEY ("name")
);
