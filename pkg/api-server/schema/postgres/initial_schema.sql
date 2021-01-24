CREATE OR REPLACE FUNCTION notify_event()
    RETURNS trigger
    LANGUAGE plpgsql
AS
$function$

DECLARE
    data         json;
    notification json;
    status_name build_status.name%type;
BEGIN

    --Get Status name(in case of both status and build table)
    SELECT name INTO status_name FROM build_status
	WHERE id = NEW.status;

    -- Convert the old or new row to JSON, based on the kind of action.
    -- Action = DELETE?             -> OLD row
    -- Action = INSERT or UPDATE?   -> NEW row
    IF (TG_OP = 'DELETE') THEN
--         data = row_to_json(OLD);
        SELECT row_to_json(r) INTO DATA FROM (SELECT OLD.*, status_name) r;
    ELSE
--         data = row_to_json(NEW);
        SELECT row_to_json(r) INTO DATA FROM (SELECT NEW.*, status_name) r;
    END IF;


    -- Contruct the notification as a JSON string.
    notification = json_build_object(
            'table', TG_TABLE_NAME,
            'action', TG_OP,
            'data', data);


    -- Execute pg_notify(channel, notification)
    PERFORM pg_notify('events', notification::text);

    -- Result is ignored since this is an AFTER trigger
    RETURN NULL;
END;
$function$
;

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
	CONSTRAINT build_status_fk FOREIGN KEY (status) REFERENCES build_status(id)
);
drop trigger if exists build_notify_event on build;
create trigger build_notify_event
    after
        insert
        or
        delete
        or
        update
    on
        build
    for each row
execute function notify_event();

CREATE TABLE step(
    id int8 NOT NULL GENERATED ALWAYS AS IDENTITY,
    build_id int8 NOT NULL,
    name VARCHAR(100) NOT NULL,
    status int4 NOT NULL,
    start_ts timestamp NULL,
	finished_ts timestamp NULL,
    created_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_ts timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT step_pk PRIMARY KEY (id),
	CONSTRAINT step_build_fk FOREIGN KEY (build_id) REFERENCES build(id),
	CONSTRAINT step_status_fk FOREIGN KEY (status) REFERENCES build_status(id)
);
drop trigger if exists step_notify_event on step;
create trigger step_notify_event
    after
        insert
        or
        delete
        or
        update
    on
        step
    for each row
execute function notify_event();

CREATE TABLE IF NOT EXISTS logs (
 step_id   int8 NOT NULL,
 log_data  bytea,
 CONSTRAINT logs_pk PRIMARY KEY (step_id)
--  CONSTRAINT logs_step_id FOREIGN KEY (step_id) REFERENCES step(id)
);