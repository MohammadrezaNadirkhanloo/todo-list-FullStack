package database

import (
	"errors"

	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgNotNullViolation    = "23502"
	pgCheckViolation      = "23514"
	pgStringDataTooLong   = "22001"
)

func TranslateError(err error, entity string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NotFound(entity)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgUniqueViolation:
			return apperror.Wrap(err, apperror.CodeConflict,
				"رکوردی با این مشخصات از قبل وجود دارد.")

		case pgForeignKeyViolation:
			return apperror.Wrap(err, apperror.CodeConflict,
				"این رکورد به رکوردهای دیگری وابسته است یا رکورد مرتبط وجود ندارد.")

		case pgNotNullViolation:
			return apperror.Wrap(err, apperror.CodeInvalidInput,
				"یکی از فیلدهای الزامی مقدار ندارد.")

		case pgCheckViolation:
			return apperror.Wrap(err, apperror.CodeInvalidInput,
				"مقدار ارسالی با محدودیت‌های تعریف‌شده سازگار نیست.")

		case pgStringDataTooLong:
			return apperror.Wrap(err, apperror.CodeInvalidInput,
				"طول یکی از مقادیر ارسالی بیش از حد مجاز است.")
		}
	}

	return apperror.Internal(err)
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) ||
		apperror.CodeOf(err) == apperror.CodeNotFound
}
