package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// ─── Goja plugin adapter ────────────────────────────────────────────────────

// gojaPluginAdapter wraps a Goja JavaScript runtime as an sdk.Plugin.
// All method calls are serialized through mu because goja.Runtime is not
// goroutine-safe.
type gojaPluginAdapter struct {
	mu      sync.Mutex
	vm      *goja.Runtime
	timeout time.Duration // per-call timeout; default 5s

	manifest sdk.Manifest // from plugin.yaml

	// Cached JS callables — extracted once at load time
	activateFn   goja.Callable
	deactivateFn goja.Callable
	onEventFn    goja.Callable

	// Extracted at load time; handler closures re-lock mu on each call
	jsFilters []jsFilterDef
	jsRoutes  []jsRouteDef
}

type jsFilterDef struct {
	Event   string
	Handler goja.Callable
}

type jsRouteDef struct {
	Method  string
	Path    string
	Auth    string
	Handler goja.Callable
}

// Compile-time interface compliance
var _ sdk.Plugin = (*gojaPluginAdapter)(nil)

// ─── Plugin loading ─────────────────────────────────────────────────────────

// loadGojaPlugin creates a Goja runtime, executes the JS source, and returns
// an adapter that implements sdk.Plugin.
func loadGojaPlugin(ctx context.Context, pluginDir, jsFile string, mf sdk.Manifest) (sdk.Plugin, error) {
	srcPath := filepath.Join(pluginDir, jsFile)
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("goja: read %s: %w", srcPath, err)
	}

	// Create runtime
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	// Enable require() scoped to pluginDir
	registry := require.NewRegistryWithLoader(func(path string) ([]byte, error) {
		full := filepath.Join(pluginDir, path)
		data, readErr := os.ReadFile(full)
		if readErr != nil {
			return nil, require.ModuleFileDoesNotExistError
		}
		return data, nil
	})
	registry.Enable(vm)

	// Enable console (wired to plugin logger once activated; for now use a
	// noop printer — the real logger is injected via ctx.log in Activate)
	console.Enable(vm)

	// Set up CommonJS module/exports globals (goja_nodejs/require only creates
	// these for files loaded through require(), not for the main script)
	mod := vm.NewObject()
	mod.Set("exports", vm.NewObject())
	vm.Set("module", mod)
	vm.Set("exports", mod.Get("exports"))

	// Pre-compile and run
	prog, err := goja.Compile(jsFile, string(src), false)
	if err != nil {
		return nil, fmt.Errorf("goja: compile %s: %w", jsFile, err)
	}
	if _, err := vm.RunProgram(prog); err != nil {
		return nil, fmt.Errorf("goja: run %s: %w", jsFile, err)
	}

	// Extract module.exports
	exports, err := extractExports(vm)
	if err != nil {
		return nil, fmt.Errorf("goja: %s: %w", jsFile, err)
	}

	adapter := &gojaPluginAdapter{
		vm:       vm,
		manifest: mf,
		timeout:  5 * time.Second,
	}

	// Extract optional lifecycle functions
	adapter.activateFn = extractCallable(exports, "activate")
	adapter.deactivateFn = extractCallable(exports, "deactivate")
	adapter.onEventFn = extractCallable(exports, "onEvent")

	// Extract filters
	if err := adapter.extractFilters(exports); err != nil {
		return nil, fmt.Errorf("goja: %s filters: %w", jsFile, err)
	}

	// Extract routes
	if err := adapter.extractRoutes(exports); err != nil {
		return nil, fmt.Errorf("goja: %s routes: %w", jsFile, err)
	}

	g.Log().Infof(ctx, "[pluginmgr] goja loaded: %s from %s", mf.ID, srcPath)
	return adapter, nil
}

// ─── sdk.Plugin implementation ──────────────────────────────────────────────

func (a *gojaPluginAdapter) Manifest() sdk.Manifest {
	return a.manifest
}

func (a *gojaPluginAdapter) Activate(pctx sdk.PluginContext) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Inject SDK services as the global `ctx` object
	injectPluginSDK(a.vm, pctx)

	if a.activateFn == nil {
		return nil
	}
	_, err := a.safeCall(a.activateFn, a.vm.Get("ctx"))
	return err
}

func (a *gojaPluginAdapter) Deactivate() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.deactivateFn == nil {
		return nil
	}
	_, err := a.safeCall(a.deactivateFn)
	return err
}

func (a *gojaPluginAdapter) Filters() []sdk.FilterDef {
	defs := make([]sdk.FilterDef, len(a.jsFilters))
	for i, jf := range a.jsFilters {
		handler := jf.Handler // capture for closure
		defs[i] = sdk.FilterDef{
			Event: jf.Event,
			Handler: func(fc *sdk.FilterContext) {
				a.mu.Lock()
				defer a.mu.Unlock()

				jsFc := a.buildFilterContext(fc)
				if _, err := a.safeCall(handler, jsFc); err != nil {
					g.Log().Warningf(context.Background(),
						"[goja] %s filter(%s) error: %v", a.manifest.ID, fc.Event, err)
				}
				// data and meta are Go maps, mutated in-place by JS
			},
		}
	}
	return defs
}

func (a *gojaPluginAdapter) Routes(r sdk.RouteRegistrar) {
	for _, jr := range a.jsRoutes {
		handler := jr.Handler // capture
		httpHandler := a.makeRouteHandler(handler)

		path := jr.Path
		if path != "" && path[0] != '/' {
			path = "/" + path
		}

		r.Handle(jr.Method, path, httpHandler, sdk.WithAuth(jr.Auth))
	}
}

func (a *gojaPluginAdapter) OnEvent(ctx context.Context, event string, data map[string]any) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.onEventFn == nil {
		return
	}
	if _, err := a.safeCall(a.onEventFn, a.vm.ToValue(event), a.vm.ToValue(data)); err != nil {
		g.Log().Warningf(ctx, "[goja] %s onEvent(%s) error: %v", a.manifest.ID, event, err)
	}
}

func (a *gojaPluginAdapter) Migrations() []sdk.Migration {
	return a.manifest.Migrations
}

// ─── Extraction helpers ─────────────────────────────────────────────────────

// extractExports retrieves the `module.exports` object from the Goja runtime.
func extractExports(vm *goja.Runtime) (*goja.Object, error) {
	mod := vm.Get("module")
	if mod == nil || goja.IsUndefined(mod) {
		return nil, fmt.Errorf("module is not defined")
	}
	exports := mod.ToObject(vm).Get("exports")
	if exports == nil || goja.IsUndefined(exports) || goja.IsNull(exports) {
		return nil, fmt.Errorf("module.exports is not defined")
	}
	return exports.ToObject(vm), nil
}

// extractCallable extracts a function from an object property. Returns nil if missing.
func extractCallable(obj *goja.Object, name string) goja.Callable {
	v := obj.Get(name)
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil
	}
	fn, ok := goja.AssertFunction(v)
	if !ok {
		return nil
	}
	return fn
}

// extractFilters reads the `filters` array from exports.
func (a *gojaPluginAdapter) extractFilters(exports *goja.Object) error {
	v := exports.Get("filters")
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil
	}
	arr := v.ToObject(a.vm)
	length := arr.Get("length")
	if length == nil {
		return nil
	}
	n := length.ToInteger()
	for i := int64(0); i < n; i++ {
		item := arr.Get(fmt.Sprintf("%d", i)).ToObject(a.vm)
		event := item.Get("event")
		handlerVal := item.Get("handler")
		if event == nil || goja.IsUndefined(event) || handlerVal == nil || goja.IsUndefined(handlerVal) {
			continue
		}
		fn, ok := goja.AssertFunction(handlerVal)
		if !ok {
			return fmt.Errorf("filter[%d].handler is not a function", i)
		}
		a.jsFilters = append(a.jsFilters, jsFilterDef{
			Event:   event.String(),
			Handler: fn,
		})
	}
	return nil
}

// extractRoutes reads the `routes` array from exports.
func (a *gojaPluginAdapter) extractRoutes(exports *goja.Object) error {
	v := exports.Get("routes")
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil
	}
	arr := v.ToObject(a.vm)
	length := arr.Get("length")
	if length == nil {
		return nil
	}
	n := length.ToInteger()
	for i := int64(0); i < n; i++ {
		item := arr.Get(fmt.Sprintf("%d", i)).ToObject(a.vm)
		method := getString(item, "method", "GET")
		path := getString(item, "path", "")
		auth := getString(item, "auth", "public")
		handlerVal := item.Get("handler")
		if handlerVal == nil || goja.IsUndefined(handlerVal) {
			continue
		}
		fn, ok := goja.AssertFunction(handlerVal)
		if !ok {
			return fmt.Errorf("route[%d].handler is not a function", i)
		}
		a.jsRoutes = append(a.jsRoutes, jsRouteDef{
			Method:  method,
			Path:    path,
			Auth:    auth,
			Handler: fn,
		})
	}
	return nil
}

// getString extracts a string property from a Goja object with a default.
func getString(obj *goja.Object, key, fallback string) string {
	v := obj.Get(key)
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return fallback
	}
	return v.String()
}

// ─── Route handler bridge ───────────────────────────────────────────────────

// makeRouteHandler wraps a JS handler function as an http.HandlerFunc.
// The JS handler receives a request object and should return a JSON-serializable value.
func (a *gojaPluginAdapter) makeRouteHandler(handler goja.Callable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.mu.Lock()
		defer a.mu.Unlock()

		// Build request object for JS
		reqObj := a.vm.NewObject()
		reqObj.Set("method", r.Method)
		reqObj.Set("url", r.URL.String())
		reqObj.Set("path", r.URL.Path)

		// Query parameters
		query := make(map[string]any)
		for k, v := range r.URL.Query() {
			if len(v) == 1 {
				query[k] = v[0]
			} else {
				query[k] = v
			}
		}
		reqObj.Set("query", query)

		// Headers
		headers := make(map[string]string)
		for k := range r.Header {
			headers[k] = r.Header.Get(k)
		}
		reqObj.Set("headers", headers)

		// Body (for POST/PUT)
		if r.Body != nil && (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
			bodyBytes, _ := io.ReadAll(r.Body)
			var bodyMap map[string]any
			if json.Unmarshal(bodyBytes, &bodyMap) == nil {
				reqObj.Set("body", bodyMap)
			} else {
				reqObj.Set("body", string(bodyBytes))
			}
		}

		// User info from context (set by wrapHandler in manager.go)
		reqObj.Set("userId", r.Context().Value("user_id"))
		reqObj.Set("userRole", r.Context().Value("user_role"))

		// Call JS handler
		result, err := a.safeCall(handler, reqObj)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"code": 500, "message": err.Error(),
			})
			return
		}

		// Write response
		if result == nil || goja.IsUndefined(result) || goja.IsNull(result) {
			json.NewEncoder(w).Encode(map[string]any{"code": 0, "message": ""})
			return
		}
		json.NewEncoder(w).Encode(result.Export())
	}
}

// ─── Filter context bridge ──────────────────────────────────────────────────

// buildFilterContext creates a JS object mirroring sdk.FilterContext.
// The `data` and `meta` fields are Go maps passed directly — Goja makes them
// mutable from JS, so changes are reflected in Go without re-export.
func (a *gojaPluginAdapter) buildFilterContext(fc *sdk.FilterContext) goja.Value {
	obj := a.vm.NewObject()
	obj.Set("event", fc.Event)
	obj.Set("data", fc.Data)
	obj.Set("meta", fc.Meta)
	obj.Set("abort", func(reason string) {
		fc.Abort(reason)
	})
	return obj
}

// ─── Safe execution ─────────────────────────────────────────────────────────

// safeCall invokes a JS function with a timeout guard.
// After a timeout, the runtime is interrupted and its state is cleared.
func (a *gojaPluginAdapter) safeCall(fn goja.Callable, args ...goja.Value) (goja.Value, error) {
	timer := time.AfterFunc(a.timeout, func() {
		a.vm.Interrupt("execution timeout")
	})
	defer timer.Stop()

	v, err := fn(goja.Undefined(), args...)
	if err != nil {
		// Clear interrupt for future calls
		a.vm.ClearInterrupt()
		if interrupted, ok := err.(*goja.InterruptedError); ok {
			return nil, fmt.Errorf("js timeout: %s", interrupted.String())
		}
		return nil, fmt.Errorf("js: %s", err.Error())
	}
	return v, nil
}
