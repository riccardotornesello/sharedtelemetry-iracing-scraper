-- Create "cars" table
CREATE TABLE "public"."cars" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "car_name" text NULL,
  "car_name_abbreviated" text NULL,
  "car_make" text NULL,
  "logo" text NULL,
  "small_image" text NULL,
  "sponsor_logo" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_cars_deleted_at" to table: "cars"
CREATE INDEX "idx_cars_deleted_at" ON "public"."cars" ("deleted_at");
-- Create "car_classes" table
CREATE TABLE "public"."car_classes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "short_name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_car_classes_deleted_at" to table: "car_classes"
CREATE INDEX "idx_car_classes_deleted_at" ON "public"."car_classes" ("deleted_at");
-- Create "car_in_classes" table
CREATE TABLE "public"."car_in_classes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "car_id" bigint NULL,
  "car_class_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_car_in_classes_car" FOREIGN KEY ("car_id") REFERENCES "public"."cars" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_car_in_classes_car_class" FOREIGN KEY ("car_class_id") REFERENCES "public"."car_classes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_car_in_classes_deleted_at" to table: "car_in_classes"
CREATE INDEX "idx_car_in_classes_deleted_at" ON "public"."car_in_classes" ("deleted_at");
