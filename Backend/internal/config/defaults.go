package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("env", string(EnvDevelopment))

	// --- سرور ---
	v.SetDefault("server.port", "5005")
	v.SetDefault("server.domain", "localhost")
	v.SetDefault("server.runMode", "debug")
	v.SetDefault("server.apiBasePath", "/api/v1")
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.readHeaderTimeout", "5s")
	v.SetDefault("server.writeTimeout", "30s")
	v.SetDefault("server.idleTimeout", "120s")
	v.SetDefault("server.shutdownTimeout", "20s")
	v.SetDefault("server.requestTimeout", "25s")
	v.SetDefault("server.maxBodyBytes", 1<<20) // ۱ مگابایت
	v.SetDefault("server.trustedProxies", []string{})
	v.SetDefault("server.insecureCookies", false)

	// --- PostgreSQL ---
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.dbName", "goclean")
	v.SetDefault("postgres.sslMode", "disable")
	v.SetDefault("postgres.timeZone", "UTC")
	v.SetDefault("postgres.maxIdleConns", 10)
	v.SetDefault("postgres.maxOpenConns", 50)
	v.SetDefault("postgres.connMaxLifetime", "30m")
	v.SetDefault("postgres.connMaxIdleTime", "5m")
	v.SetDefault("postgres.slowQueryThreshold", "200ms")

	// --- Redis ---
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.dialTimeout", "5s")
	v.SetDefault("redis.readTimeout", "3s")
	v.SetDefault("redis.writeTimeout", "3s")
	v.SetDefault("redis.poolSize", 10)
	v.SetDefault("redis.minIdleConns", 2)

	// --- JWT ---

	v.SetDefault("jwt.accessTTL", "15m")
	v.SetDefault("jwt.refreshTTL", "168h") // ۷ روز
	v.SetDefault("jwt.issuer", "goclean")
	v.SetDefault("jwt.audience", "goclean-api")

	// --- CORS ---
	v.SetDefault("cors.allowOrigins", []string{})
	v.SetDefault("cors.allowCredentials", true)
	v.SetDefault("cors.maxAge", "6h")

	// --- سیاست دسترسی مسیرها ---
	//
	// پیش‌فرض enterprise است: هر مسیر جدیدی که اضافه کنید، به‌صورت
	// پیش‌فرض بسته است. اگر یادتان برود گروهش را درست انتخاب کنید،
	// بدترین اتفاق «۴۰۱ ناخواسته» است، نه «نشت داده».
	//
	// این همان اصل secure by default است.
	v.SetDefault("auth.posture", string(PostureEnterprise))

	// --- مجوزدهی ---
	v.SetDefault("authz.mode", string(AuthzBoth))
	v.SetDefault("authz.enforce", true)

	// --- لاگ ---
	v.SetDefault("logger.backend", "zerolog")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.encoding", "json")
	v.SetDefault("logger.filePath", "")
	v.SetDefault("logger.logRequestBody", false)

	// --- OTP (ماژول اختیاری — docs/optional-services/otp-auth.md) ---
	//
	// پیش‌فرض خاموش است. یک قالب عمومی نباید فرض کند هر پروژه‌ای
	// پیامک می‌فرستد؛ پیامک هزینه، قرارداد با ارائه‌دهنده و
	// پیچیدگی عملیاتی دارد.
	v.SetDefault("otp.enabled", false)
	v.SetDefault("otp.digits", 6)
	v.SetDefault("otp.expireTime", "2m")
	v.SetDefault("otp.resendDelay", "90s")
	v.SetDefault("otp.maxAttempts", 5)

	// --- رمز عبور ---
	v.SetDefault("password.minLength", 12)
	v.SetDefault("password.maxLength", 128)
	v.SetDefault("password.includeDigits", true)
	v.SetDefault("password.includeUppercase", true)
	v.SetDefault("password.includeLowercase", true)
	v.SetDefault("password.includeSymbols", false)
	v.SetDefault("password.argon2Memory", 65536) // 64 MiB
	v.SetDefault("password.argon2Iterations", 3)
	v.SetDefault("password.argon2Parallelism", 2)

	// --- محدودکننده‌ی نرخ ---
	v.SetDefault("ratelimit.enabled", true)
	v.SetDefault("ratelimit.globalRPS", 20.0)
	v.SetDefault("ratelimit.globalBurst", 40)
	v.SetDefault("ratelimit.authRPS", 0.2) // ۱ درخواست هر ۵ ثانیه
	v.SetDefault("ratelimit.authBurst", 5)
	v.SetDefault("ratelimit.otpRPS", 0.017) // حدود ۱ درخواست در دقیقه
	v.SetDefault("ratelimit.otpBurst", 2)
	v.SetDefault("ratelimit.idleTTL", "15m")
	v.SetDefault("ratelimit.maxKeys", 100000)

	// --- صفحه‌بندی ---
	v.SetDefault("pagination.defaultPageSize", 20)
	v.SetDefault("pagination.maxPageSize", 100)
}

var secretKeys = []string{
	"postgres.password",
	"redis.password",
	"jwt.accessSecret",
	"jwt.refreshSecret",
}

func bindEnvKeys(v *viper.Viper) {
	for _, key := range secretKeys {
		_ = v.BindEnv(key)
	}
}
