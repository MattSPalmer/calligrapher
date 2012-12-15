package main

func getCallbacks(calls []CallRecord) []string {
	needsCallback := make(map[string]bool)
	callbackNumbers := make([]string, 0)

	for _, c := range calls {
		if c.IsMissed {
			if _, ok := needsCallback[c.IncomingNumber]; !ok {
				needsCallback[c.IncomingNumber] = true
			}
		} else {
			needsCallback[c.IncomingNumber] = false
		}
	}
	for k, v := range needsCallback {
		if v {
			callbackNumbers = append(callbackNumbers, k)
		}
	}

	return callbackNumbers
}
