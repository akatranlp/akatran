package utils

func Map[TInput, TOutput any](input []TInput, f func(TInput) TOutput) []TOutput {
	output := make([]TOutput, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}

func Filter[TInput any](input []TInput, predicate func(TInput) bool) []TInput {
	output := make([]TInput, 0)
	for _, v := range input {
		if predicate(v) {
			output = append(output, v)
		}
	}
	return output
}
