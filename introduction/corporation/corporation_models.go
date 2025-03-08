package corporation

import "factorySimulator/commonModels"

type task struct {
	args      []any
	operation *commonModels.Operation
}

// value is the operation result done by worker.
type product struct {
	value any
}

type taskRequest struct {
	response chan task
}

type buyRequest struct {
	response chan product
}
