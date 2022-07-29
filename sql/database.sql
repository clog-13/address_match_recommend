DROP TABLE IF EXISTS "public"."address";
-- This script only contains the table creation statements and does not fully represent the table in database. It's still missing: indices, triggers. Do not use it as backup.

-- Squences
CREATE SEQUENCE IF NOT EXISTS address_id_seq

-- Table Definition
CREATE TABLE "public"."address" (
    "id" int8 NOT NULL DEFAULT nextval('address_id_seq'::regclass),
    "address_text" text,
    "road_text" text,
    "road_num" text,
    "building_num" text,
    "province_id" int8,
    "city_id" int8,
    "district_id" int8,
    "street_id" int8,
    "town_id" int8,
    "village_id" int8
);

-- Column Comments
COMMENT ON COLUMN "public"."address"."id" IS '地址ID';
COMMENT ON COLUMN "public"."address"."address_text" IS '完整地址';
COMMENT ON COLUMN "public"."address"."road_text" IS '道路信息';
COMMENT ON COLUMN "public"."address"."road_num" IS '道路号';
COMMENT ON COLUMN "public"."address"."building_num" IS '建筑信息';

DROP TABLE IF EXISTS "public"."document";
-- This script only contains the table creation statements and does not fully represent the table in database. It's still missing: indices, triggers. Do not use it as backup.

-- Squences
CREATE SEQUENCE IF NOT EXISTS document_id_seq

-- Table Definition
CREATE TABLE "public"."document" (
    "id" int8 NOT NULL DEFAULT nextval('document_id_seq'::regclass),
    "town_id" int8,
    "road_num_value" int8,
    "village_id" int8,
    "road_id" int8,
    "road_num_id" int8,
    PRIMARY KEY ("id")
);

-- Column Comments
COMMENT ON COLUMN "public"."document"."id" IS '文档ID';

DROP TABLE IF EXISTS "public"."region";
-- This script only contains the table creation statements and does not fully represent the table in database. It's still missing: indices, triggers. Do not use it as backup.

-- Squences
CREATE SEQUENCE IF NOT EXISTS region_id_seq

-- Table Definition
CREATE TABLE "public"."region" (
    "id" int4 NOT NULL DEFAULT nextval('region_id_seq'::regclass),
    "parent_id" int8,
    "name" text,
    "alias" text,
    "types" int2,
    "division_id" int8,
    PRIMARY KEY ("id")
);

-- Column Comments
COMMENT ON COLUMN "public"."region"."id" IS '行政区域ID';
COMMENT ON COLUMN "public"."region"."parent_id" IS '完整地址';
COMMENT ON COLUMN "public"."region"."name" IS '区域名称';
COMMENT ON COLUMN "public"."region"."alias" IS '区域别名';
COMMENT ON COLUMN "public"."region"."types" IS '区域类型';

DROP TABLE IF EXISTS "public"."term";
-- This script only contains the table creation statements and does not fully represent the table in database. It's still missing: indices, triggers. Do not use it as backup.

-- Squences
CREATE SEQUENCE IF NOT EXISTS "Term_id_seq"

-- Table Definition
CREATE TABLE "public"."term" (
    "id" int8 NOT NULL DEFAULT nextval('"Term_id_seq"'::regclass),
    "text" text,
    "types" int2,
    "idf" numeric,
    "term_id" int8,
    PRIMARY KEY ("id")
);

-- Column Comments
COMMENT ON COLUMN "public"."term"."id" IS '词条ID';
COMMENT ON COLUMN "public"."term"."text" IS '词条字段';
COMMENT ON COLUMN "public"."term"."types" IS '词条类型';
COMMENT ON COLUMN "public"."term"."idf" IS 'IDF';

