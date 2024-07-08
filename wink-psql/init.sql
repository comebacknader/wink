CREATE TABLE users(
	ID serial PRIMARY KEY,
    username VARCHAR(100),
	email VARCHAR(100),
	hash VARCHAR(100),
	usertype VARCHAR(100),
	resetpasstoken VARCHAR(100),
	resetpassexpiry TIMESTAMP,
	coins integer,
	purchtoken VARCHAR(100));

CREATE TABLE frontpage(
	id serial PRIMARY KEY,
	user_id integer REFERENCES users(id));

CREATE TABLE streams(
           ID serial PRIMARY KEY,
           user_id integer REFERENCES users(id),
           title VARCHAR(100),
           game VARCHAR(100),
           withgame boolean,
           online boolean,
           twitter varchar(100),
           siteone varchar(100));

CREATE TABLE sessions(
           id serial PRIMARY KEY,
           user_id integer REFERENCES users(id),
           uuid VARCHAR(256),
           usertype VARCHAR(50),
           activity TIMESTAMP);


CREATE TABLE purchases(
           id serial PRIMARY KEY,
           user_id integer REFERENCES users(id),
           type VARCHAR(100),
           price VARCHAR(50),
           bought TIMESTAMP);


