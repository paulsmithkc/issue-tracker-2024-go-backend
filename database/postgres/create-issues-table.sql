create table
  "public"."issues" (
    "id" UUID not null,
    "created_at" timestamp not null default NOW(),
    "title" varchar(30) not null,
    "description" TEXT not null,
    constraint "issues_pkey" primary key ("id")
  )
