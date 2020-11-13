CREATE TABLE build_status (
	id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" varchar(50) NULL,
	CONSTRAINT build_status_pk PRIMARY KEY (id),
	CONSTRAINT build_status_un UNIQUE (name)
);

INSERT INTO build_status (name) VALUES('Pending');
INSERT INTO build_status (name) VALUES('Started');
INSERT INTO build_status (name) VALUES('Finished');


CREATE TABLE secret_type(
    id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    name varchar(100) NOT NULL,
    CONSTRAINT secret_type_pk PRIMARY KEY (id),
    CONSTRAINT secret_type_name_un UNIQUE (name)
);
INSERT INTO secret_type(name) VALUES('Kubernetes');

CREATE TABLE auth_type(
    id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    name varchar(100) NOT NULL,
    CONSTRAINT auth_type_pk PRIMARY KEY (id),
    CONSTRAINT auth_type_name_un UNIQUE (name)
);
INSERT INTO auth_type(name) VALUES('BasicAuth');
INSERT INTO auth_type(name) VALUES('SSHAuth');

CREATE TABLE repo(
    id int8 NOT NULL GENERATED ALWAYS AS IDENTITY,
    namespace VARCHAR(1000) NOT NULL,
    name VARCHAR(1000) NOT NULL,
    ssh_url VARCHAR(2000) NOT NULL,
    http_url VARCHAR(2000) NOT NULL,
    auth_type int4  NOT NULL,
    secret_type int4  NOT NULL,
    secret_name VARCHAR(2000) NOT NULL,
    created_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT repo_pk PRIMARY KEY (id),
	CONSTRAINT secret_type_fk FOREIGN KEY (secret_type) REFERENCES secret_type(id),
	CONSTRAINT auth_type_fk FOREIGN KEY (auth_type) REFERENCES auth_type(id)
);

CREATE TABLE build(
	id int8 NOT NULL GENERATED ALWAYS AS IDENTITY,
	repo_branch int8 NOT NULL,
	status int4 NOT NULL,
	start_ts timestamp NULL,
	finished_ts timestamp NULL,
	created_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT build_pk PRIMARY KEY (id),
	CONSTRAINT build_fk FOREIGN KEY (status) REFERENCES build_status(id)
);
