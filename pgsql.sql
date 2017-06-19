-- Connect to Postgres.
-- psql -U goprojects

-- Create a user named `goprojects` with the password `goprojects`. If you ever
-- do this in production, please use a better password.

-- * I had to change the permissions to use md5 for both local
-- * connections and for peer [psql] and create the database in a
-- * separate step

-- # TYPE  DATABASE        USER            ADDRESS                 METHOD
-- local   all             goprojects                              md5
-- host    all             goprojects             127.0.0.1/32     md5

-- 2016-12-15 13:21:32 EST [12064-1] goprojects@goprojects ERROR:  permission denied to create role
-- 2016-12-15 13:21:32 EST [12064-2] goprojects@goprojects STATEMENT:  CREATE USER goprojects WITH PASSWORD 'goprojects';
-- 2016-12-15 13:21:32 EST [12064-3] goprojects@goprojects ERROR:  permission denied to create database
-- 2016-12-15 13:21:32 EST [12064-4] goprojects@goprojects STATEMENT:  CREATE DATABASE goprojects;
-- 2016-12-15 13:21:32 EST [12064-5] goprojects@goprojects ERROR:  syntax error at or near "goprojects" at character 14
-- 2016-12-15 13:21:32 EST [12064-6] goprojects@goprojects STATEMENT:  SET database goprojects;
-- 2016-12-15 13:21:32 EST [12064-7] goprojects@goprojects WARNING:  no privileges were granted for "goprojects"



-- postgres@meson:~$ psql -U goprojects --dbname=goprojects --password  
-- Password for user goprojects: 
-- psql (9.5.5)
-- Type "help" for help.

-- goprojects=> \i /go/src/goprojects/cms3.4/init.sql 

-- goprojects=> \pset pager off
-- Pager usage is off.
-- goprojects=> \dt
--            List of relations
--  Schema |   Name   | Type  |   Owner    
-- --------+----------+-------+------------
--  public | comments | table | goprojects
--  public | pages    | table | goprojects
--  public | posts    | table | goprojects
-- (3 rows)

CREATE USER gorilla WITH PASSWORD 'gorilla'; 

\c gorilla;
CREATE DATABASE gorilla;

GRANT ALL ON DATABASE gorilla TO gorilla;

GRANT ALL PRIVILEGES ON DATABASE gorilla to gorilla;
/*GRANT CREATE ON gorilla TO 'gorilla'@'%';*/


/*
CREATE TABLE pages (
  id  serial primary key,
  page_guid varchar(256) NOT NULL DEFAULT '' unique,
  page_title varchar(256) DEFAULT NULL,
  page_content text,
  page_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;

CREATE OR REPLACE FUNCTION update_page_date_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.page_date = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_ab_changetimestamp BEFORE UPDATE ON pages FOR EACH ROW EXECUTE PROCEDURE update_page_date_column();
*/
/* null not allowed for serial column for primary key id, but no value auto updates */
/*
 INSERT INTO pages (page_guid, page_title, page_content, page_date)
 VALUES ('hello-world', 'Hello, World', 'I''m so glad you found this page!  It''s been sitting patiently on the Internet for some time, just waiting for a visitor.', CURRENT_TIMESTAMP);
*/

/*
CREATE USER goprojects WITH PASSWORD 'goprojects';

-- Create the database we're going to use.
CREATE DATABASE goprojects;

SET database goprojects;
-- Grant all privleges to our user on the DB.
GRANT ALL PRIVILEGES ON DATABASE goprojects to goprojects;

-- Create a new table to store our pages.
CREATE TABLE IF NOT EXISTS PAGES(
  id             SERIAL    PRIMARY KEY,
  title          TEXT      NOT NULL,
  content        TEXT      NOT NULL
);

-- Create a new table to store our posts.
CREATE TABLE IF NOT EXISTS POSTS(
  id             SERIAL    PRIMARY KEY,
  title          TEXT      NOT NULL,
  content        TEXT      NOT NULL,
  date_created   DATE      NOT NULL
);

-- Create a new table to store our comments.
CREATE TABLE IF NOT EXISTS COMMENTS(
  id             SERIAL    PRIMARY KEY,
  author         TEXT      NOT NULL,
  content        TEXT      NOT NULL,
  date_created   DATE      NOT NULL,
  post_id        INT       references POSTS(id)
);
*/
