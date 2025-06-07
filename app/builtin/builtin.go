package builtin

type Builtin interface {
	Run(args []string)
}