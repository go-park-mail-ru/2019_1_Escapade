\c erybase

drop table if exists Erythrocytes CASCADE;
drop table if exists UsersInProjects CASCADE;
drop table if exists Projects CASCADE;
drop table if exists ProjectTokens CASCADE;
drop table if exists Diseases CASCADE;
drop table if exists EryObject CASCADE;
drop table if exists Users CASCADE;

-- 255 -> 299

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
		if NEW."photo_title" is NULL or LENGTH(NEW."photo_title") = 0 then
			NEW."photo_title" = '1.png';
		end if;
		if NEW."website" is NULL or LENGTH(NEW."website") = 0 then
			NEW."website" = 'не указан';
		end if;
		if NEW."about" is NULL or LENGTH(NEW."about") = 0 then
			NEW."about" = 'информация не указана';
		end if;
		if NEW."email" is NULL or LENGTH(NEW."email") = 0 then
			NEW."email" = 'не указан';
		end if;
		if NEW."phone" is NULL or LENGTH(NEW."phone") = 0 then
			NEW."phone" = 'не указан';
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fill_name_about() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if NEW."name" is null or LENGTH(NEW."name") = 0 then
			NEW."name" = 'Безымянный';
		end if;
		if NEW."about" is null or LENGTH(NEW."about") = 0 then
			NEW."about" = 'отсутствует';
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


CREATE OR REPLACE FUNCTION fill_position() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if NEW."position" is NULL or LENGTH(NEW."position") = 0 then
			NEW."position" = 'Не указана';
		end if;
    return NEW; 
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION project_update_edit() RETURNS trigger AS $TRIGGER$ 
BEGIN 
		if TG_OP = 'DELETE' then
			UPDATE Projects SET edit = now(), editor_id = old.user_id where id=old.project_id;
			return OLD; 
		else
				UPDATE Projects SET edit = now(), editor_id = new.user_id where id=new.project_id;
		end if;
    return NEW;  
END; 
$TRIGGER$ 
LANGUAGE plpgsql;

CREATE TABLE Projects (
  id SERIAL PRIMARY KEY,
	name text default 'Безымянный',
	public_access bool default false,
	company_access bool default false,
	public_edit bool default false, 
	company_edit bool default false,
	about text default 'нет информации',
	edit  TIMESTAMPTZ default now(),
	editor_id int NOT NULL,
	add TIMESTAMPTZ default now()
);

CREATE TRIGGER projects_name_about
BEFORE INSERT OR UPDATE ON Projects FOR EACH ROW 
EXECUTE PROCEDURE fill_name_about();

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

CREATE TRIGGER up_position
BEFORE INSERT OR UPDATE ON UsersInProjects FOR EACH ROW 
EXECUTE PROCEDURE fill_position();

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

CREATE TRIGGER scene_name_about
BEFORE INSERT OR UPDATE ON Scene FOR EACH ROW 
EXECUTE PROCEDURE fill_name_about();

CREATE TRIGGER scene_update_project
BEFORE INSERT OR UPDATE OR DELETE ON Scene FOR EACH ROW 
EXECUTE PROCEDURE project_update_edit();

-----------------------------------------------

CREATE TABLE Diseases (
  id SERIAL PRIMARY KEY,
	name text NOT NULL,
	about text default 'отсутствует',
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
		if TG_OP = 'DELETE' then
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

CREATE TRIGGER diseases_name_about
BEFORE INSERT OR UPDATE ON Diseases FOR EACH ROW 
EXECUTE PROCEDURE fill_name_about();

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

CREATE TRIGGER eryobj_name_about
BEFORE INSERT OR UPDATE ON EryObject FOR EACH ROW 
EXECUTE PROCEDURE fill_name_about();

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
