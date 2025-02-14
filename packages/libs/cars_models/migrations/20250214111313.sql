-- Modify "car_classes" table
ALTER TABLE "public"."car_classes" DROP COLUMN "deleted_at";
-- Modify "car_in_classes" table
ALTER TABLE "public"."car_in_classes" DROP COLUMN "deleted_at";
-- Modify "cars" table
ALTER TABLE "public"."cars" DROP COLUMN "deleted_at";
