package logging



type Category string

const (
	CategoryGeneral    Category = "General"
	CategoryIO         Category = "IO"
	CategoryPostgres   Category = "Postgres"
	CategoryRedis      Category = "Redis"
	CategoryValidation Category = "Validation"
	CategoryRequest    Category = "RequestResponse"
	CategoryAuth       Category = "Auth"
	CategoryInternal   Category = "Internal"
)

type SubCategory string

const (
	SubStartup    SubCategory = "Startup"
	SubShutdown   SubCategory = "Shutdown"
	SubMigration  SubCategory = "Migration"
	SubSelect     SubCategory = "Select"
	SubInsert     SubCategory = "Insert"
	SubUpdate     SubCategory = "Update"
	SubDelete     SubCategory = "Delete"
	SubAPI        SubCategory = "Api"
	SubLogin      SubCategory = "Login"
	SubRegister   SubCategory = "Register"
	SubToken      SubCategory = "Token"
	SubOtp        SubCategory = "Otp"
	SubRateLimit  SubCategory = "RateLimit"
	SubHashing    SubCategory = "Hashing"
	SubConnection SubCategory = "Connection"
)

func Cat(c Category) Field { return Field{Key: "category", Value: string(c)} }

func Sub(s SubCategory) Field { return Field{Key: "sub_category", Value: string(s)} }

