CREATE TABLE build_status (
	id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" varchar(50) NULL,
	CONSTRAINT build_status_pk PRIMARY KEY (id),
	CONSTRAINT build_status_un UNIQUE (name)
);

INSERT INTO build_status (name) VALUES('Pending');
INSERT INTO build_status (name) VALUES('Started');
INSERT INTO build_status (name) VALUES('Finished');

CREATE TABLE build (
	id int8 NOT NULL GENERATED ALWAYS AS IDENTITY,
	project_id int8 NOT NULL,
	status int4 NOT NULL,
	start_ts timestamp NULL,
	finished_ts timestamp NULL,
	created_ts timestamp NOT NULL,
	updated_ts timestamp NOT NULL,
	CONSTRAINT build_pk PRIMARY KEY (id),
	CONSTRAINT build_fk FOREIGN KEY (status) REFERENCES build_status(id)
);