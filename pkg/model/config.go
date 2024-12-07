package model

type Config struct {
	// Порт http сервера
	Port int `json:"port" default:"8000"`
	// Слушать только локальные соединения
	LocalhostOnly bool `json:"localhost_only" default:"true"`
	// Подписываются ли запросы полем "token"
	CheckAuth bool `json:"check_auth" default:"false"`
	// Драйвер подключения к базе данных
	DBDriver string `json:"db_driver" default:"postgres"`
	// Строка подключения к базе данных
	DBConnectionString string `json:"db_connection_string" default:"postgresql://127.0.0.1:5432"`
	// Максимальное число соединений с базой
	DBMaxConnections int `json:"db_max_connections" default:"20"`
	// Максимальный размер пачки
	DBBatchSize int `json:"db_batch_size" default:"200"`
}
