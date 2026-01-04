package postgres

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"rest/config"
	errorsPkg "rest/internal/pkg/errorspkg"
	"rest/internal/repo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDeps struct {
	Conf *config.Postgres
}

type postgres struct {
	pool *pgxpool.Pool
}

type PostgresI interface {
	GetWalletByID(ctx context.Context, walletID string) (*models.Wallet, error)
	MakeTransaction(ctx context.Context, walletID, operation string, amount float64) error
	Stop()
}

func New(ctx context.Context, dep PostgresDeps) (PostgresI, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dep.Conf.User, dep.Conf.Password, dep.Conf.Host, dep.Conf.Port, dep.Conf.Database,
	) // путь для подключения к postgres

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25 // Максимум соединений
	config.MinConns = 5  // Минимум соединений в пуле
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config) // Подключение с пулом соединений
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil { // проверка подключения
		pool.Close()
		return nil, err
	}

	return &postgres{
		pool: pool,
	}, nil
}

func (p *postgres) GetWalletByID(ctx context.Context, walletID string) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `
		SELECT id, balance, status, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	err := p.pool.QueryRow(ctx, query, walletID).Scan( // отправляем атомарный запрос, никакие блокировки не нужны
		&wallet.ID,
		&wallet.Balance,
		&wallet.Status,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsPkg.NewWalletNotFoundError(walletID)
		}
		return nil, errorsPkg.NewDatabaseError(err)
	}

	return &wallet, nil
}

func (db *postgres) MakeTransaction(ctx context.Context, walletID, operation string, amount float64) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	}
	var currentBalance float64
	var status string
	walletQuery := `
			SELECT balance, status
			FROM wallets
			WHERE id = $1
			FOR UPDATE
	`

	tx, err := db.pool.BeginTx(ctx, opts) //начинаем транзакцию, чтобы избежать несогласованности данных
	if err != nil {
		return errorsPkg.NewDatabaseError(err)
	}
	defer tx.Rollback(ctx) // откатываемся в случае ошибок и всего прочего(если ошибки не будет, rollback перекроется функцией tx.Commit)

	err = tx.QueryRow(ctx, walletQuery, walletID).Scan(&currentBalance, &status) // командой select for update блокируем нужные нам поля, пока не закончим транзакцию
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorsPkg.NewWalletNotFoundError(walletID)
		}
		return errorsPkg.NewDatabaseError(err)
	}

	if status != "ACTIVE" { // проверка состояния кошелька
		return &errorsPkg.AppError{
			Status: http.StatusForbidden,
			Err:    fmt.Errorf("wallet is not active (current status: %s)", status),
		}
	}

	var newBalance float64

	switch models.OperationType(operation) { // логика подсчёта новой суммы
	case models.OperationDeposit:
		newBalance = currentBalance + amount
	case models.OperationWithdraw:
		if currentBalance < amount {
			return errorsPkg.NewInsufficientFundsError(currentBalance, amount)
		}
		newBalance = currentBalance - amount
	default:
		return errorsPkg.NewBadRequestError(fmt.Errorf("invalid operation type: %s", operation))
	}

	updateQuery := `
			UPDATE wallets
			SET balance = $1
			WHERE id = $2
		`

	updateResult, err := tx.Exec(ctx, updateQuery, newBalance, walletID) // обновляем наш баланс на кошельке и проверяем, всё ли в порядке
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" { // проверяем constraints
			return errorsPkg.NewInsufficientFundsError(currentBalance, amount)
		}
		return errorsPkg.NewDatabaseError(err)
	}

	if updateResult.RowsAffected() == 0 { 
		return errorsPkg.NewDatabaseError(err)
	}

	if err := tx.Commit(ctx); err != nil { // завершаем транзакцию, если всё в порядке
		return errorsPkg.NewDatabaseError(err)
	}

	return nil
}

func (db *postgres) Stop() {
	if db.pool != nil {
		db.pool.Close()
	}
}
