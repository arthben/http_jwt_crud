USE db_todo;
CREATE TABLE todos (
	id varchar(64),
	title varchar(256),
	detail text,
	created_date timestamp,
	updated_date timestamp,
	st_completed char(1),
	completed_date timestamp,
	primary key(id)
) engine=Innodb;
