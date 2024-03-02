module github.com/sj14/epoch/v4

go 1.17

retract (
	v4.0.0-alpha.1
	v4.0.0-alpha.2
	v4.0.0-alpha.3
	v4.0.0-alpha.4
	v4.0.0-alpha.5
	v4.0.0-alpha.6
	[v4.0.0-0, v4.0.0-alpha.6]
)

require github.com/stretchr/testify v1.8.2

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
