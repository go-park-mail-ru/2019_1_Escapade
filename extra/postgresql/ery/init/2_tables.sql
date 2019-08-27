\c erybase

drop table if exists Erythrocytes CASCADE;
drop table if exists Session CASCADE;
drop table if exists ProjectInCollection CASCADE;
drop table if exists ProjectCollection CASCADE;
drop table if exists ImageInProject CASCADE;
drop table if exists TextureInProject CASCADE;
drop table if exists DiseaseInProject CASCADE;
drop table if exists UsersInProjects CASCADE;
drop table if exists CompaniesInProjects CASCADE;
drop table if exists Projects CASCADE;
drop table if exists ProjectTokens CASCADE;
drop table if exists UsersInCompany CASCADE;
drop table if exists CompanyTokens CASCADE;
drop table if exists Companies CASCADE;
drop table if exists CompanyLogs CASCADE;
drop table if exists Diseases CASCADE;
drop table if exists EryObject CASCADE;
drop table if exists EryImages CASCADE;
drop table if exists EryTextures CASCADE;
drop table if exists EryForms CASCADE;
drop table if exists UserLogs CASCADE;
drop table if exists Logs CASCADE;
drop table if exists Users CASCADE;
drop function if exists logs CASCADE;
drop table if exists Scene CASCADE;

-- 255

CREATE TABLE Users (
  id SERIAL PRIMARY KEY,
  name text NOT NULL,
  password text NOT NULL,
  photo_title text default '1.png',
	website text default 'не указан',
	about text default 'информация не указана',
	email text default 'не указан',
	phone text default 'не указан',
	birthday TIMESTAMPTZ default now(),
	add TIMESTAMPTZ default now(),
	last_seen TIMESTAMPTZ default now(),

	UNIQUE (name)
);

CREATE OR REPLACE FUNCTION users_trigger_func() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if LENGTH(NEW."photo_title") = 0 then
			NEW."photo_title" = OLD."photo_title";
		end if;
		if LENGTH(NEW."password") = 0 then
			NEW."password" = OLD."password";
		end if;
		if LENGTH(NEW."website") = 0 then
			NEW."website" = 'не указан';
		end if;
		if LENGTH(NEW."about") = 0 then
			NEW."about" = 'информация не указана';
		end if;
		if LENGTH(NEW."email") = 0 then
			NEW."email" = 'не указан';
		end if;
		if LENGTH(NEW."phone") = 0 then
			NEW."phone" = 'не указан';
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION users_trigger_insert_func() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if LENGTH(NEW."photo_title") = 0 then
			NEW."photo_title" = '1.png';
		end if;
		if LENGTH(NEW."website") = 0 then
			NEW."website" = 'не указан';
		end if;
		if LENGTH(NEW."about") = 0 then
			NEW."about" = 'информация не указана';
		end if;
		if LENGTH(NEW."email") = 0 then
			NEW."email" = 'не указан';
		end if;
		if LENGTH(NEW."phone") = 0 then
			NEW."phone" = 'не указан';
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE TRIGGER users_trigger_update
BEFORE UPDATE ON Users FOR EACH ROW 
EXECUTE PROCEDURE users_trigger_func();

CREATE TRIGGER users_trigger_insert
BEFORE INSERT ON Users FOR EACH ROW 
EXECUTE PROCEDURE users_trigger_insert_func();

-- CREATE Table Logs (
-- 	id SERIAL PRIMARY KEY,
-- 	user_id int NOT NULL,
-- 	"text" text NOT NULL,
-- 	"date" TIMESTAMPTZ default now(),

-- 	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
-- );

-- CREATE Table UserLogs (
-- 	id SERIAL PRIMARY KEY,
-- 	user_id int NOT NULL,
-- 	log_id int NOT NULL,

-- 	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
-- 	FOREIGN KEY (log_id) REFERENCES Logs (id) ON DELETE CASCADE
-- );

-- create function user_logs() returns trigger language plpgsql as $$
-- 				DECLARE
-- 						log_id int;
--         begin
--         if (tg_op='INSERT') then
--                 log_id := insert into Logs ( user_id, "add ", schema_name, new_row ) values ( new.user_id, tg_table_name, tg_table_schema, new );
--         elsif (tg_op='UPDATE') then
--                 if (old<>new) then
--                         insert into Logs ( user_id, table_name, schema_name, old_row, new_row ) values ( new.user_id, tg_table_name, tg_table_schema, old, new );
--                 end if;
--         elsif (tg_op='DELETE') then
--                 insert into Logs ( user_id, table_name, schema_name, old_row ) values ( old.user_id, tg_table_name, tg_table_schema, old );
--         end if;
--         return null;
--         end;
-- $$;

-------------------------------------------------------------

CREATE Table Companies (
	id SERIAL PRIMARY KEY,
	name text NOT NULL,
	website text,
	about text,
	contacts text,
	add TIMESTAMPTZ default now(),
	public bool default false
);

CREATE Table CompanyTokens (
	id SERIAL PRIMARY KEY,
	owner bool,
	access bool,
	edit_members_list bool,
	edit_members_position bool,
	edit_members_token bool,
	edit_company_name bool,
	edit_company_token bool,
	edit_company_website bool,
	edit_company_about bool,
	edit_company_contacts bool,
	edit_company_private bool
);

CREATE TABLE UsersInCompany (
  id SERIAL PRIMARY KEY,
	position text,
  user_id int,
	company_id int,
	token_id int,
  "from" TIMESTAMPTZ,
	"to" TIMESTAMPTZ,

	user_confirmed bool default false,
	company_confirmed bool default false,

	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
	FOREIGN KEY (company_id) REFERENCES Companies (id) ON DELETE CASCADE,
	FOREIGN KEY (token_id) REFERENCES CompanyTokens (id) ON DELETE CASCADE
);

-- CREATE Table CompanyLogs (
-- 	id SERIAL PRIMARY KEY,
-- 	user_id int NOT NULL,
-- 	company_id int NOT NULL,
-- 	log_id int NOT NULL,
	
-- 	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
-- 	FOREIGN KEY (company_id) REFERENCES Companies (id) ON DELETE CASCADE,
-- 	FOREIGN KEY (log_id) REFERENCES CompanyLogs (id) ON DELETE CASCADE
-- );

------------------------------------------------

CREATE TABLE Projects (
  id SERIAL PRIMARY KEY,
	name text default 'Безымянный',
	public_access bool default false,
	company_access bool default false,
	public_edit bool default false, 
	company_edit bool default false,
	about text default 'нет информации',
	add TIMESTAMPTZ default now()
);

CREATE Table ProjectTokens (
	id SERIAL PRIMARY KEY,
	owner bool default false,
	edit_name bool default false,
	edit_info bool default false,
	edit_access bool default false,
	edit_scene bool default false,
	edit_members_list bool default false,
	edit_members_token bool default false
);

CREATE TABLE UsersInProjects (
  id SERIAL PRIMARY KEY,
	position text default 'не указана',
  user_id int NOT NULL,
	token_id int NOT NULL,
	project_id int NOT NULL,
  "from" TIMESTAMPTZ,
	"to" TIMESTAMPTZ,
	user_confirmed bool default false,
	project_confirmed bool default false,

	FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
	FOREIGN KEY (token_id) REFERENCES ProjectTokens(id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,

	UNIQUE (user_id, project_id)
);

CREATE TABLE CompaniesInProjects (
  id SERIAL PRIMARY KEY,
	position text,
  company_id int,
	project_id int,
  "from" TIMESTAMPTZ,
	"to" TIMESTAMPTZ,
	company_confirmed bool default false,
	project_confirmed bool default false,

	FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,

	UNIQUE (company_id, project_id)
);

CREATE TABLE Scene (
  id SERIAL PRIMARY KEY,
	name text default 'Безымянный',
	about text default 'Отсутсвует',
	user_id int NOT NULL,
  project_id int NOT NULL,
	edit TIMESTAMPTZ NOT NULL,
	editor_id  int NOT NULL,
	add TIMESTAMPTZ default now(),

	FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
	FOREIGN KEY (editor_id) REFERENCES Users(id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION scene_trigger_func() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if LENGTH(NEW."name") = 0 then
			NEW."name" = 'Безымянный';
		end if;
		if LENGTH(NEW."about") = 0 then
			NEW."about" = 'Отсутсвует';
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE TRIGGER scene_trigger
BEFORE UPDATE OR INSERT ON Scene FOR EACH ROW 
EXECUTE PROCEDURE scene_trigger_func();

-----------------------------------------------

CREATE TABLE ProjectCollection (
	id SERIAL PRIMARY KEY,
	user_id int,
	name text,
	about text,
	add TIMESTAMPTZ default now(),

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE ProjectInCollection (
	id SERIAL PRIMARY KEY,
	collection_id int,
	project_id int,
	add TIMESTAMPTZ default now(),

	FOREIGN KEY (collection_id) REFERENCES ProjectCollection (id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE, 

	UNIQUE (collection_id, project_id)
);

------------------------------------------------

CREATE TABLE Diseases (
  id SERIAL PRIMARY KEY,
	user_id int,
	scene_id int NOT NULL,
	form float,
	oxygen float,
	gemoglob float,
  add TIMESTAMPTZ,

	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
	FOREIGN KEY (scene_id) REFERENCES Scene (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_scene_edit_time() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if TG_OP = "DELETE" then
			UPDATE Scene SET edit = now(), editor_id = user_id where id=old.scene_id;
			return OLD; 
		else
			UPDATE Scene SET edit = now(), editor_id = user_id where id=new.scene_id;
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE TRIGGER diseases_trigger
AFTER UPDATE OR INSERT OR DELETE ON Diseases FOR EACH ROW 
EXECUTE PROCEDURE update_scene_edit_time();

CREATE TABLE EryObject (
  id SERIAL PRIMARY KEY,
	user_id int NOT NULL,
	scene_id int NOT NULL,
	path text NOT NULL,
	name text NOT NULL,
	about text default 'отсутствует',
	source text default 'не указан',
	public bool default false,
	is_form  bool default false,
	is_texture  bool default false,
	is_image bool default false,
  add TIMESTAMPTZ default now(),

	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
	FOREIGN KEY (scene_id) REFERENCES Scene (id) ON DELETE CASCADE
);

CREATE TRIGGER eryobj_trigger
AFTER UPDATE OR INSERT OR DELETE ON EryObject FOR EACH ROW 
EXECUTE PROCEDURE update_scene_edit_time();

------------------------------------------------

CREATE TABLE Erythrocytes (
  id SERIAL PRIMARY KEY,
  user_id int,
	texture_id int,
	form_id int,
	scene_id int,
	disease_id int,

	size_x float,
	size_y float,
	size_z float,

	angle_x float,
	angle_y float,
	angle_z float,

	scale_x float,
	scale_y float,
	scale_z float,

	position_x float,
	position_y float,
	position_z float,

	form float,
	oxygen float,
	gemoglob float,

  add TIMESTAMPTZ default now(),

	FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
	FOREIGN KEY (texture_id) REFERENCES EryObject (id) ON DELETE CASCADE,
	FOREIGN KEY (form_id) REFERENCES EryObject (id) ON DELETE CASCADE,
	FOREIGN KEY (scene_id) REFERENCES Scene (id) ON DELETE CASCADE
);

CREATE TRIGGER ery_trigger
AFTER UPDATE OR INSERT OR DELETE ON Erythrocytes FOR EACH ROW 
EXECUTE PROCEDURE update_scene_edit_time();

-- create table Logs (
--         id serial primary key,
--         mod_time timestamp default now(),
-- 				user_id int,
--         table_name text,
--         schema_name text,
--         old_row text,
--         new_row text
-- );
-- create function logs() returns trigger language plpgsql as $$
--         begin
--         if (tg_op='INSERT') then
--                 insert into Logs ( user_id, table_name, schema_name, new_row ) values ( new.user_id, tg_table_name, tg_table_schema, new );
--         elsif (tg_op='UPDATE') then
--                 if (old<>new) then
--                         insert into Logs ( user_id, table_name, schema_name, old_row, new_row ) values ( new.user_id, tg_table_name, tg_table_schema, old, new );
--                 end if;
--         elsif (tg_op='DELETE') then
--                 insert into Logs ( user_id, table_name, schema_name, old_row ) values ( old.user_id, tg_table_name, tg_table_schema, old );
--         end if;
--         return null;
--         end;
-- $$;

-- create table t1 ( id integer, name text );
-- create trigger t1 after insert or update or delete on t1 for each row execute procedure hs();

-- create table t2 ( f1 text, f2 date );
-- create trigger t2 after insert or update or delete on t2 for each row execute procedure hs();



-- CREATE TABLE IF NOT EXISTS the_log (
-- 	"timestamp" timestamp with time zone DEFAULT now() NOT NULL,
-- 	"user" int NOT NULL,
-- 	action text NOT NULL,
-- 	table_name text NOT NULL,
-- 	old_row jsonb,
-- 	new_row jsonb,
-- 	CONSTRAINT the_log_check CHECK (
-- 		CASE action
-- 			WHEN 'INSERT' THEN old_row IS NULL
-- 			WHEN 'DELETE' THEN new_row IS NULL
-- 		END
-- 	)
-- );

-- CREATE TRIGGER log_insert
-- AFTER INSERT ON foo
-- REFERENCING NEW TABLE AS new_table
-- FOR EACH STATEMENT
-- BEGIN
-- 	INSERT INTO the_log (
-- 	action, table_name, new_row
-- 	)
-- 	SELECT
-- 	TG_OP, TG_RELNAME, row_to_json(new_table)::
-- 	FROM
-- 	new_table;
-- END;

-- CREATE TRIGGER log_update_public_foo
-- AFTER UPDATE ON foo
-- REFERENCING OLD TABLE AS old_table NEW TABLE AS new_table
-- FOR EACH STATEMENT
-- BEGIN
-- 	INSERT INTO the_log (
-- 	action, table_schema, table_name, old_row, new_row
-- 	)
-- 	SELECT
-- 	TG_OP, TG_TABLE_SCHEMA, TG_RELNAME, old_row, new_row
-- 	FROM
-- 	UNNEST(
-- 	ARRAY(SELECT row_to_json(old_table)::jsonb FROM old_table
-- 	ARRAY(SELECT row_to_json(new_table)::jsonb FROM new_table
-- 	) AS t(old_row, new_row)
-- END;


-- CREATE TRIGGER log_update_public_foo
-- AFTER DELETE ON foo
-- REFERENCING OLD TABLE AS old_table NEW TABLE AS new_table
-- FOR EACH STATEMENT
-- BEGIN
-- 	INSERT INTO the_log (
-- 	action, table_schema, table_name, old_row, new_row
-- 	)
-- 	SELECT
-- 	TG_OP, TG_TABLE_SCHEMA, TG_RELNAME, old_row, new_row
-- 	FROM
-- 	UNNEST(
-- 	ARRAY(SELECT row_to_json(old_table)::jsonb FROM old_table
-- 	ARRAY(SELECT row_to_json(new_table)::jsonb FROM new_table
-- 	) AS t(old_row, new_row)
-- END;



	