package pluginsys

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// injectPluginSDK registers the `ctx` global on the Goja runtime, exposing
// platform services (db, store, settings, log, http) to the JavaScript plugin.
//
// Each sub-object (ctx.db, ctx.store, …) is built by a dedicated inject* helper
// so that adding a new service is a single-function change in this file.
func injectPluginSDK(vm *goja.Runtime, pctx sdk.PluginContext) {
	ctx := vm.NewObject()
	injectDB(vm, ctx, pctx.DB)
	injectStore(vm, ctx, pctx.Store)
	injectSettings(vm, ctx, pctx.Settings)
	injectLog(vm, ctx, pctx.Log)
	injectHTTP(vm, ctx)
	vm.Set("ctx", ctx)
}

// ── ctx.db ──────────────────────────────────────────────────────────────────

func injectDB(vm *goja.Runtime, ctx *goja.Object, db sdk.DB) {
	o := vm.NewObject()

	// ctx.db.query(sql, ...args) → []row
	o.Set("query", func(call goja.FunctionCall) goja.Value {
		sql := call.Argument(0).String()
		args := exportArgs(call.Arguments[1:])
		rows, err := db.Query(sql, args...)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(rows)
	})

	// ctx.db.execute(sql, ...args) → affectedRows
	o.Set("execute", func(call goja.FunctionCall) goja.Value {
		sql := call.Argument(0).String()
		args := exportArgs(call.Arguments[1:])
		n, err := db.Execute(sql, args...)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(n)
	})

	ctx.Set("db", o)
}

// ── ctx.store ───────────────────────────────────────────────────────────────

func injectStore(vm *goja.Runtime, ctx *goja.Object, store sdk.Store) {
	o := vm.NewObject()

	o.Set("get", func(key string) goja.Value {
		val, err := store.Get(key)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(val)
	})

	o.Set("set", func(key string, value goja.Value) {
		if err := store.Set(key, value.Export()); err != nil {
			panic(vm.NewGoError(err))
		}
	})

	o.Set("delete", func(key string) {
		if err := store.Delete(key); err != nil {
			panic(vm.NewGoError(err))
		}
	})

	o.Set("increment", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		delta := int64(1)
		if len(call.Arguments) > 1 && !goja.IsUndefined(call.Argument(1)) {
			delta = call.Argument(1).ToInteger()
		}
		val, err := store.Increment(key, delta)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(val)
	})

	o.Set("deletePrefix", func(prefix string) goja.Value {
		n, err := store.DeletePrefix(prefix)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(n)
	})

	ctx.Set("store", o)
}

// ── ctx.settings ────────────────────────────────────────────────────────────

func injectSettings(vm *goja.Runtime, ctx *goja.Object, settings sdk.Settings) {
	o := vm.NewObject()

	o.Set("get", func(key string) goja.Value {
		return vm.ToValue(settings.Get(key))
	})

	o.Set("getAll", func() goja.Value {
		return vm.ToValue(settings.GetAll())
	})

	ctx.Set("settings", o)
}

// ── ctx.log ─────────────────────────────────────────────────────────────────

func injectLog(vm *goja.Runtime, ctx *goja.Object, logger sdk.Logger) {
	o := vm.NewObject()
	o.Set("info", func(msg string) { logger.Info(msg) })
	o.Set("warn", func(msg string) { logger.Warn(msg) })
	o.Set("error", func(msg string) { logger.Error(msg) })
	o.Set("debug", func(msg string) { logger.Debug(msg) })
	ctx.Set("log", o)
}

// ── ctx.http ────────────────────────────────────────────────────────────────

// injectHTTP provides ctx.http.fetch(url, options?) for JS plugins.
// Options: { method, headers, body, timeout }.
func injectHTTP(vm *goja.Runtime, ctx *goja.Object) {
	o := vm.NewObject()

	o.Set("fetch", func(call goja.FunctionCall) goja.Value {
		rawURL := call.Argument(0).String()
		method := "GET"
		var body io.Reader
		headers := map[string]string{}
		timeout := 10 * time.Second

		if len(call.Arguments) > 1 && !goja.IsUndefined(call.Argument(1)) {
			opts := call.Argument(1).ToObject(vm)
			if m := opts.Get("method"); m != nil && !goja.IsUndefined(m) {
				method = strings.ToUpper(m.String())
			}
			if b := opts.Get("body"); b != nil && !goja.IsUndefined(b) {
				body = strings.NewReader(b.String())
			}
			if h := opts.Get("headers"); h != nil && !goja.IsUndefined(h) {
				if m, ok := h.Export().(map[string]interface{}); ok {
					for k, v := range m {
						headers[k] = fmt.Sprintf("%v", v)
					}
				}
			}
			if t := opts.Get("timeout"); t != nil && !goja.IsUndefined(t) {
				timeout = time.Duration(t.ToInteger()) * time.Millisecond
			}
		}

		req, err := http.NewRequest(method, rawURL, body)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		client := &http.Client{Timeout: timeout}
		resp, err := client.Do(req)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		result := vm.NewObject()
		result.Set("status", resp.StatusCode)
		result.Set("body", string(respBody))
		respHeaders := vm.NewObject()
		for k := range resp.Header {
			respHeaders.Set(strings.ToLower(k), resp.Header.Get(k))
		}
		result.Set("headers", respHeaders)
		return result
	})

	ctx.Set("http", o)
}

// ── helpers ─────────────────────────────────────────────────────────────────

// exportArgs converts a slice of goja values to Go interface{} values.
func exportArgs(args []goja.Value) []any {
	out := make([]any, len(args))
	for i, a := range args {
		out[i] = a.Export()
	}
	return out
}
