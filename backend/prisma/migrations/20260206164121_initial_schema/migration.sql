-- CreateTable
CREATE TABLE "users" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "email" VARCHAR(255) NOT NULL,
    "username" VARCHAR(100) NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "validate_recipe_volumes" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "refresh_tokens" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "token_hash" VARCHAR(255) NOT NULL,
    "expires_at" TIMESTAMP(3) NOT NULL,
    "revoked" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "refresh_tokens_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "invitations" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "code" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255),
    "created_by" UUID NOT NULL,
    "expires_at" TIMESTAMP(3) NOT NULL,
    "used" BOOLEAN NOT NULL DEFAULT false,
    "used_at" TIMESTAMP(3),
    "used_by" UUID,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "invitations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "accords" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "pyramid_position" VARCHAR(10) NOT NULL,
    "volume_ml" DECIMAL(10,2) NOT NULL,
    "volume_drops" INTEGER,
    "supplier" VARCHAR(255),
    "purchase_date" DATE,
    "dilution_percentage" DECIMAL(5,2),
    "notes" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "accords_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "accord_tags" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "accord_id" UUID NOT NULL,
    "tag" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "accord_tags_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "predefined_tags" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "category" VARCHAR(50) NOT NULL,
    "tag" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "predefined_tags_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipes" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "target_volume_ml" DECIMAL(10,2) NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'draft',
    "active_version_id" UUID,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_versions" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "recipe_id" UUID NOT NULL,
    "version_number" INTEGER NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "notes" TEXT,
    "is_active" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_versions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_ingredients" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "version_id" UUID NOT NULL,
    "accord_id" UUID NOT NULL,
    "quantity_ml" DECIMAL(10,2) NOT NULL,
    "quantity_drops" INTEGER,
    "percentage" DECIMAL(5,2),
    "notes" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_ingredients_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_notes" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "recipe_id" UUID NOT NULL,
    "version_id" UUID,
    "content" TEXT NOT NULL,
    "note_type" VARCHAR(20) NOT NULL DEFAULT 'general',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_notes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_tags" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "recipe_id" UUID NOT NULL,
    "tag" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_tags_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_collections" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_collections_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "recipe_collection_members" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "collection_id" UUID NOT NULL,
    "recipe_id" UUID NOT NULL,
    "added_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "recipe_collection_members_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- CreateIndex
CREATE UNIQUE INDEX "refresh_tokens_token_hash_key" ON "refresh_tokens"("token_hash");

-- CreateIndex
CREATE INDEX "refresh_tokens_user_id_idx" ON "refresh_tokens"("user_id");

-- CreateIndex
CREATE INDEX "refresh_tokens_token_hash_idx" ON "refresh_tokens"("token_hash");

-- CreateIndex
CREATE UNIQUE INDEX "invitations_code_key" ON "invitations"("code");

-- CreateIndex
CREATE INDEX "invitations_code_idx" ON "invitations"("code");

-- CreateIndex
CREATE INDEX "invitations_created_by_idx" ON "invitations"("created_by");

-- CreateIndex
CREATE INDEX "accords_user_id_idx" ON "accords"("user_id");

-- CreateIndex
CREATE INDEX "accords_pyramid_position_idx" ON "accords"("pyramid_position");

-- CreateIndex
CREATE INDEX "accords_created_at_idx" ON "accords"("created_at" DESC);

-- CreateIndex
CREATE UNIQUE INDEX "accords_user_id_name_pyramid_position_key" ON "accords"("user_id", "name", "pyramid_position");

-- CreateIndex
CREATE INDEX "accord_tags_accord_id_idx" ON "accord_tags"("accord_id");

-- CreateIndex
CREATE INDEX "accord_tags_tag_idx" ON "accord_tags"("tag");

-- CreateIndex
CREATE UNIQUE INDEX "accord_tags_accord_id_tag_key" ON "accord_tags"("accord_id", "tag");

-- CreateIndex
CREATE UNIQUE INDEX "predefined_tags_tag_key" ON "predefined_tags"("tag");

-- CreateIndex
CREATE INDEX "predefined_tags_category_idx" ON "predefined_tags"("category");

-- CreateIndex
CREATE INDEX "recipes_user_id_idx" ON "recipes"("user_id");

-- CreateIndex
CREATE INDEX "recipes_status_idx" ON "recipes"("status");

-- CreateIndex
CREATE INDEX "recipes_created_at_idx" ON "recipes"("created_at" DESC);

-- CreateIndex
CREATE UNIQUE INDEX "recipes_user_id_name_key" ON "recipes"("user_id", "name");

-- CreateIndex
CREATE INDEX "recipe_versions_recipe_id_idx" ON "recipe_versions"("recipe_id");

-- CreateIndex
CREATE INDEX "recipe_versions_is_active_idx" ON "recipe_versions"("is_active");

-- CreateIndex
CREATE UNIQUE INDEX "recipe_versions_recipe_id_version_number_key" ON "recipe_versions"("recipe_id", "version_number");

-- CreateIndex
CREATE INDEX "recipe_ingredients_version_id_idx" ON "recipe_ingredients"("version_id");

-- CreateIndex
CREATE INDEX "recipe_ingredients_accord_id_idx" ON "recipe_ingredients"("accord_id");

-- CreateIndex
CREATE UNIQUE INDEX "recipe_ingredients_version_id_accord_id_key" ON "recipe_ingredients"("version_id", "accord_id");

-- CreateIndex
CREATE INDEX "recipe_notes_recipe_id_idx" ON "recipe_notes"("recipe_id");

-- CreateIndex
CREATE INDEX "recipe_notes_version_id_idx" ON "recipe_notes"("version_id");

-- CreateIndex
CREATE INDEX "recipe_tags_recipe_id_idx" ON "recipe_tags"("recipe_id");

-- CreateIndex
CREATE INDEX "recipe_tags_tag_idx" ON "recipe_tags"("tag");

-- CreateIndex
CREATE UNIQUE INDEX "recipe_tags_recipe_id_tag_key" ON "recipe_tags"("recipe_id", "tag");

-- CreateIndex
CREATE INDEX "recipe_collections_user_id_idx" ON "recipe_collections"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "recipe_collections_user_id_name_key" ON "recipe_collections"("user_id", "name");

-- CreateIndex
CREATE INDEX "recipe_collection_members_collection_id_idx" ON "recipe_collection_members"("collection_id");

-- CreateIndex
CREATE INDEX "recipe_collection_members_recipe_id_idx" ON "recipe_collection_members"("recipe_id");

-- CreateIndex
CREATE UNIQUE INDEX "recipe_collection_members_collection_id_recipe_id_key" ON "recipe_collection_members"("collection_id", "recipe_id");

-- AddForeignKey
ALTER TABLE "refresh_tokens" ADD CONSTRAINT "refresh_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "invitations" ADD CONSTRAINT "invitations_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "invitations" ADD CONSTRAINT "invitations_used_by_fkey" FOREIGN KEY ("used_by") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "accords" ADD CONSTRAINT "accords_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "accord_tags" ADD CONSTRAINT "accord_tags_accord_id_fkey" FOREIGN KEY ("accord_id") REFERENCES "accords"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipes" ADD CONSTRAINT "recipes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipes" ADD CONSTRAINT "recipes_active_version_id_fkey" FOREIGN KEY ("active_version_id") REFERENCES "recipe_versions"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_versions" ADD CONSTRAINT "recipe_versions_recipe_id_fkey" FOREIGN KEY ("recipe_id") REFERENCES "recipes"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_ingredients" ADD CONSTRAINT "recipe_ingredients_version_id_fkey" FOREIGN KEY ("version_id") REFERENCES "recipe_versions"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_ingredients" ADD CONSTRAINT "recipe_ingredients_accord_id_fkey" FOREIGN KEY ("accord_id") REFERENCES "accords"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_notes" ADD CONSTRAINT "recipe_notes_recipe_id_fkey" FOREIGN KEY ("recipe_id") REFERENCES "recipes"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_notes" ADD CONSTRAINT "recipe_notes_version_id_fkey" FOREIGN KEY ("version_id") REFERENCES "recipe_versions"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_tags" ADD CONSTRAINT "recipe_tags_recipe_id_fkey" FOREIGN KEY ("recipe_id") REFERENCES "recipes"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_collections" ADD CONSTRAINT "recipe_collections_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_collection_members" ADD CONSTRAINT "recipe_collection_members_collection_id_fkey" FOREIGN KEY ("collection_id") REFERENCES "recipe_collections"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recipe_collection_members" ADD CONSTRAINT "recipe_collection_members_recipe_id_fkey" FOREIGN KEY ("recipe_id") REFERENCES "recipes"("id") ON DELETE CASCADE ON UPDATE CASCADE;
