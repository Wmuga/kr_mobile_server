package logger

const insert = `INSERT INTO public.logs
(uid, "level", request_id, timing, msg, payload)
VALUES($1, $2, $3, $4, $5, $6);`
