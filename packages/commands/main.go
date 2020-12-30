package commands

type CommandID int

const (
	DlFile CommandID = iota
	DAddress
	Ok
	Stream
)

var CommandMap map[string]CommandID = map[string]CommandID{
	"DL_FILE":   DlFile,
	"D_ADDRESS": DAddress,
	"OK":        Ok,
	"STREAM":    Stream,
}
