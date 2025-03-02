-- Create "competition_classes" table
CREATE TABLE "public"."competition_classes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "competition_id" bigint NOT NULL,
  "name" text NOT NULL,
  "color" text NOT NULL,
  "index" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_competition_classes_competition" FOREIGN KEY ("competition_id") REFERENCES "public"."competitions" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_competition_classes_deleted_at" to table: "competition_classes"
CREATE INDEX "idx_competition_classes_deleted_at" ON "public"."competition_classes" ("deleted_at");
-- Modify "competition_crews" table
ALTER TABLE "public"."competition_crews" DROP COLUMN "class", ADD COLUMN "class_id" bigint NULL, ADD CONSTRAINT "fk_competition_crews_class" FOREIGN KEY ("class_id") REFERENCES "public"."competition_classes" ("id") ON UPDATE SET NULL ON DELETE SET NULL;
