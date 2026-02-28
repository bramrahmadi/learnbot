module github.com/learnbot/api-gateway

go 1.24.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/learnbot/resume-parser v0.0.0
)

require github.com/dslipak/pdf v0.0.2 // indirect

replace github.com/learnbot/resume-parser => ../resume-parser
