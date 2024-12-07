CREATE TABLE public.employee (
	employee_id text NOT NULL,
	"name" text NOT NULL,
	"position" text NOT NULL,
	active bool NULL,
	CONSTRAINT employee_pk PRIMARY KEY (employee_id)
);
CREATE INDEX employee_name_idx ON public.employee USING btree (name);


CREATE TABLE public.mac_employee (
	mac text NOT NULL,
	employee_id text NOT NULL,
	CONSTRAINT mac_employee_pk PRIMARY KEY (mac)
);

CREATE TABLE public.passtime (
	uid text NOT NULL,
	passtime timestamptz NOT NULL,
	mac text NOT NULL,
	CONSTRAINT passtime_unique UNIQUE (uid)
);
CREATE INDEX passtime_passtime_idx ON public.passtime USING btree (passtime);
ALTER TABLE public.passtime ADD CONSTRAINT passtime_mac_employee_fk FOREIGN KEY (mac) REFERENCES public.mac_employee(mac);


CREATE TABLE public.logs (
	uid text NOT NULL,
	"level" text NOT NULL,
	request_id text NOT NULL,
	timing text NOT NULL,
	msg text NOT NULL,
	payload jsonb NOT NULL,
	CONSTRAINT logs_pk PRIMARY KEY (uid)
);
CREATE INDEX logs_level_idx ON public.logs USING btree (level);
CREATE INDEX logs_msg_idx ON public.logs USING btree (msg);
CREATE INDEX logs_payload_idx ON public.logs USING btree (payload);
CREATE INDEX logs_request_id_idx ON public.logs USING btree (request_id);
CREATE INDEX logs_timing_idx ON public.logs USING btree (timing);