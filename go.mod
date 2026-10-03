module github.com/Ne0nd0g/merlin-cli

go 1.27.0

require (
	github.com/chzyer/readline v1.5.1
	github.com/fatih/color v1.19.0
	github.com/google/uuid v1.6.0
	github.com/mattn/go-shellwords v1.0.16
	github.com/olekukonko/tablewriter v0.0.5
	google.golang.org/grpc v1.83.2 // pinned <1.84: v1.84.0 regressed GO-2026-6443; next stable fix is v1.85.x
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.30 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260928230214-8a89bd6388cc // indirect
)
