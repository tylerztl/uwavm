package exec

// A Resolver resolves global and function symbols imported by wasm code
type Resolver interface {
	ResolveFunc(module, name string) (interface{}, bool)
	ResolveGlobal(module, name string) (int64, bool)
}

// MultiResolver chains multiple Resolvers, symbol looking up is according to the order of resolvers.
// The first found symbol will be returned.
type MultiResolver []Resolver

// NewMultiResolver instance a MultiResolver from resolves
func NewMultiResolver(resolvers ...Resolver) MultiResolver {
	return resolvers
}

// ResolveFunc implements Resolver interface
func (m MultiResolver) ResolveFunc(module, name string) (interface{}, bool) {
	for _, r := range m {
		if f, ok := r.ResolveFunc(module, name); ok {
			return f, true
		}
	}
	return nil, false
}

// ResolveGlobal implements Resolver interface
func (m MultiResolver) ResolveGlobal(module, name string) (int64, bool) {
	for _, r := range m {
		if v, ok := r.ResolveGlobal(module, name); ok {
			return v, true
		}
	}
	return 0, false
}

func applyFuncCall(ctx Context, f interface{}, params []uint32) (uint32, bool) {
	len := len(params)
	switch fun := f.(type) {
	case func(Context) uint32:
		if len != 0 {
			return 0, false
		}
		return fun(ctx), true
	case func(Context, uint32) uint32:
		if len != 1 {
			return 0, false
		}
		return fun(ctx, params[0]), true
	case func(Context, uint32, uint32) uint32:
		if len != 2 {
			return 0, false
		}
		return fun(ctx, params[0], params[1]), true
	case func(Context, uint32, uint32, uint32) uint32:
		if len != 3 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2]), true
	case func(Context, uint32, uint32, uint32, uint32) uint32:
		if len != 4 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2], params[3]), true
	case func(Context, uint32, uint32, uint32, uint32, uint32) uint32:
		if len != 5 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2], params[3], params[4]), true
	case func(Context, uint32, uint32, uint32, uint32, uint32, uint32) uint32:
		if len != 6 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2], params[3], params[4], params[5]), true
	case func(Context, uint32, uint32, uint32, uint32, uint32, uint32, uint32) uint32:
		if len != 7 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2], params[3], params[4], params[5], params[6]), true
	case func(Context, uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32) uint32:
		if len != 8 {
			return 0, false
		}
		return fun(ctx, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7]), true
	default:
		return 0, false
	}
}