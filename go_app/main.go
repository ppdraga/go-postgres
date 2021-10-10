package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PG_HOST     = "127.0.0.1"
	PG_PORT     = "5432"
	PG_USER     = "postgres"
	PG_PASSWORD = "postgres"
	PG_DB_NAME  = "cdn"
)

type Config struct {
	//Mode            int
	MaxConns        int32
	MinConns        int32
	AttackMS        time.Duration
	GoroutinesCount int
}

var sqls []string

func main() {
	rand.Seed(time.Now().UnixNano())

	sqls = []string{
		`select file.name from userfile as uf left join file on uf.file_id = file.id
				where uf.removed is null and uf.user_id in (select id from "user" where name = 'user1');`,
		`select distinct file.name FROM serverfile AS sf left join file on sf.file_id = file.id 
				where sf.removed is null and sf.server_id in (select id from "server" s  
				where s.area_id in (select id from area where name = 'msk.ru'));`,
		`SELECT file.name FROM serverfile AS sf left join file on sf.file_id = file.id 
				where sf.removed is null and sf.server_id in (select id from "server" s  where s.name = 'srv1');`,
	}

	config := &Config{
		MaxConns:        15,
		MinConns:        15,
		AttackMS:        time.Duration(10000),
		GoroutinesCount: 30,
	}
	generateDBLoad(config)

}

func generateDBLoad(cfg *Config) error {
	pool, err := createPGXPool(cfg.MaxConns, cfg.MinConns)
	if err != nil {
		return fmt.Errorf("failed to create a PGX pool: %w", err)
	}
	res, err := Attack(
		context.Background(),
		time.Millisecond*cfg.AttackMS,
		cfg.GoroutinesCount,
		pool,
	)
	if err != nil {
		return fmt.Errorf("attack failed: %w", err)
	}
	//log.Println("Attack result: ", res)
	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
	return nil
}

func createPGXPool(maxConns int32, minConns int32) (*pgxpool.Pool, error) {
	cfg, err := getPoolConfig(maxConns, minConns)
	if err != nil {
		return nil, fmt.Errorf("failed to get the pool config: %w", err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	return pool, nil
}

func getPoolConfig(maxConns int32, minConns int32) (*pgxpool.Config, error) {
	connStr := ComposeConnectionString()
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create a pool config from connection string %s: %w", connStr, err,
		)
	}
	cfg.MaxConns = maxConns
	cfg.MinConns = minConns

	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute

	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour

	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute

	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second

	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	return cfg, nil
}

func ComposeConnectionString() string {
	userspec := fmt.Sprintf("%s:%s", PG_USER, PG_PASSWORD)
	hostspec := fmt.Sprintf("%s:%s", PG_HOST, PG_PORT)
	return fmt.Sprintf("postgresql://%s@%s/%s", userspec, hostspec, PG_DB_NAME)
}

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func (r *AttackResults) String() string {
	return fmt.Sprintf(
		"duration: %dms, threads: %d, queries performed: %d",
		(r.Duration / time.Millisecond),
		r.Threads,
		r.QueriesPerformed,
	)
}

func Attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) (*AttackResults, error) {
	var queries uint64
	attacker := func(stopAt time.Time) {
		for {
			sql := sqls[rand.Intn(len(sqls))]
			rows, err := dbpool.Query(ctx, sql)
			if err != nil {
				log.Printf("failed to fetch the email search hints: %w", err)
				return
			}
			defer rows.Close()
			filenames := []string{}
			for rows.Next() {
				filename := ""
				if err := rows.Scan(&filename); err != nil {
					log.Printf("failed to scan the received email hint: %w", err)
					return
				}
				filenames = append(filenames, filename)
			}

			if err != nil {
				log.Printf("an error occurred while searching by email: %v", err)
				continue
			}
			atomic.AddUint64(&queries, 1)
			if time.Now().After(stopAt) {
				return
			}
		}
	}

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return &AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}, nil
}
