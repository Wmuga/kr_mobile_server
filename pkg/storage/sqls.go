package storage

const (
	placeholderLIMIT  = "$limit"
	placeholderOFFSET = "$offset"
)

const sqlAddPass = "" +
	"INSERT INTO public.passtime\n" +
	"(uid, passtime, mac)\n" +
	"VALUES($1, $2, $3);"

const sqlSelectAll = `select e."name", e."position", p."mac", p."passtime" from passtime p
left join
mac_employee me on me.mac = p.mac
left join 
employee e on me.employee_id = e.employee_id
LIMIT $limit
OFFSET $offset;`

const sqlSelectToday = `select e."name", e."position", p."mac", p."passtime" from passtime p
left join
mac_employee me on me.mac = p.mac
left join 
employee e on me.employee_id = e.employee_id
WHERE p.passtime >= $1
LIMIT $limit
OFFSET $offset;`

const sqlSelectResult = `select e."name", e."position", p."mac", p."passtime" from passtime p
left join
mac_employee me on me.mac = p.mac
left join 
employee e on me.employee_id = e.employee_id
where p.uid = $1;`
