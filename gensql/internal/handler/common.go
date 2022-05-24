package handler

import (
	"context"
	"fmt"
	"github.com/gomelon/tools/gencore"
	"k8s.io/gengo/types"
)

func TplFileNameAndRowType(methodArgs *gencore.MethodArgs) (string, *types.Type) {
	params := methodArgs.Params
	results := methodArgs.Results
	var tplFileName string
	//TODO 要考虑到Alias的情况
	firstResult := results[0]
	var rowType *types.Type
	switch firstResult.Kind {
	case types.Slice:
		tplFileName += "Slice"
		rowType = firstResult.Elem
	case types.Pointer:
		tplFileName += "Get"
		rowType = firstResult
	case types.Map:
		tplFileName += "Map"
	}

	if rowType.Kind == types.Pointer {
		rowType = rowType.Elem
	}
	if len(rowType.Members) == 0 {
		gencore.FatalfOnErr(
			fmt.Errorf("invalid method signature: query method [%s.%s]'s row result must at least 1 column",
				methodArgs.Type.Name.Name, methodArgs.MethodName),
		)
	}

	if gencore.Implements(params[0], (*context.Context)(nil)) {
		tplFileName += "WithCtx"
	}

	if gencore.IsError(results[len(results)-1]) {
		tplFileName += "ReturnErr"
	}
	return tplFileName, rowType
}
