package web

type MidHandler func(Handler) Handler

func wrapMid(md []MidHandler, handler Handler) Handler {

	for i := len(md) - 1; i >= 0; i-- {
		mFunc := md[i]

		if mFunc != nil {
			handler = mFunc(handler)
		}
	}

	return handler

}
