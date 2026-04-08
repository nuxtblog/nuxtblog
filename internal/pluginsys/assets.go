package pluginsys

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ─── Plugin Asset Serving (Phase 2.8) ──────────────────────────────────────
//
// Serves plugin frontend files (admin.mjs, public.js, style.css, etc.) from
// the data/plugins/{id}/ directory.
//
// Route: GET /api/plugins/{id}/assets/{filename}
//
// Only allowed extensions: .js, .mjs, .css, .json
// Sets Cache-Control and ETag headers for optimal caching.

var allowedAssetExts = map[string]bool{
	".js":   true,
	".mjs":  true,
	".css":  true,
	".json": true,
}

var extContentType = map[string]string{
	".js":   "application/javascript; charset=utf-8",
	".mjs":  "application/javascript; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".json": "application/json; charset=utf-8",
}

// vueShimJS is a tiny ES module that re-exports Vue from the host's window global.
// Plugin ESM builds import from '/_shared/vue.mjs' which resolves to this shim.
const vueShimJS = `// Auto-generated Vue shim — re-exports from host window.__nuxtblog_vue
// Uses a Proxy so that imports resolve lazily (the global may not be set
// when the ES module first executes due to dynamic import timing).
function _get(name) {
  const V = window.__nuxtblog_vue;
  if (!V) throw new Error('[plugin-shim] window.__nuxtblog_vue not found — ensure plugin-shared-deps plugin loaded');
  return V[name];
}

// Reactivity
export const ref = (...a) => _get('ref')(...a);
export const reactive = (...a) => _get('reactive')(...a);
export const computed = (...a) => _get('computed')(...a);
export const watch = (...a) => _get('watch')(...a);
export const watchEffect = (...a) => _get('watchEffect')(...a);
export const readonly = (...a) => _get('readonly')(...a);
export const shallowRef = (...a) => _get('shallowRef')(...a);
export const shallowReactive = (...a) => _get('shallowReactive')(...a);
export const triggerRef = (...a) => _get('triggerRef')(...a);
export const toRef = (...a) => _get('toRef')(...a);
export const toRefs = (...a) => _get('toRefs')(...a);
export const toRaw = (...a) => _get('toRaw')(...a);
export const markRaw = (...a) => _get('markRaw')(...a);
export const unref = (...a) => _get('unref')(...a);
export const isRef = (...a) => _get('isRef')(...a);
export const isReactive = (...a) => _get('isReactive')(...a);
export const isReadonly = (...a) => _get('isReadonly')(...a);
export const isProxy = (...a) => _get('isProxy')(...a);

// Lifecycle
export const onMounted = (...a) => _get('onMounted')(...a);
export const onUpdated = (...a) => _get('onUpdated')(...a);
export const onUnmounted = (...a) => _get('onUnmounted')(...a);
export const onBeforeMount = (...a) => _get('onBeforeMount')(...a);
export const onBeforeUpdate = (...a) => _get('onBeforeUpdate')(...a);
export const onBeforeUnmount = (...a) => _get('onBeforeUnmount')(...a);
export const onActivated = (...a) => _get('onActivated')(...a);
export const onDeactivated = (...a) => _get('onDeactivated')(...a);
export const onErrorCaptured = (...a) => _get('onErrorCaptured')(...a);
export const onRenderTracked = (...a) => _get('onRenderTracked')(...a);
export const onRenderTriggered = (...a) => _get('onRenderTriggered')(...a);

// Composition
export const provide = (...a) => _get('provide')(...a);
export const inject = (...a) => _get('inject')(...a);
export const defineComponent = (...a) => _get('defineComponent')(...a);
export const defineAsyncComponent = (...a) => _get('defineAsyncComponent')(...a);
export const defineProps = (...a) => _get('defineProps')(...a);
export const defineEmits = (...a) => _get('defineEmits')(...a);
export const defineExpose = (...a) => _get('defineExpose')(...a);

// Render helpers (used by compiled SFC templates)
export const h = (...a) => _get('h')(...a);
export const createVNode = (...a) => _get('createVNode')(...a);
export const createElementVNode = (...a) => _get('createElementVNode')(...a);
export const createElementBlock = (...a) => _get('createElementBlock')(...a);
export const createBlock = (...a) => _get('createBlock')(...a);
export const createCommentVNode = (...a) => _get('createCommentVNode')(...a);
export const createTextVNode = (...a) => _get('createTextVNode')(...a);
export const createStaticVNode = (...a) => _get('createStaticVNode')(...a);
export const openBlock = (...a) => _get('openBlock')(...a);
export const resolveComponent = (...a) => _get('resolveComponent')(...a);
export const resolveDirective = (...a) => _get('resolveDirective')(...a);
export const resolveDynamicComponent = (...a) => _get('resolveDynamicComponent')(...a);
export const withDirectives = (...a) => _get('withDirectives')(...a);
export const withModifiers = (...a) => _get('withModifiers')(...a);
export const withCtx = (...a) => _get('withCtx')(...a);
export const withKeys = (...a) => _get('withKeys')(...a);
export const renderList = (...a) => _get('renderList')(...a);
export const renderSlot = (...a) => _get('renderSlot')(...a);
export const toDisplayString = (...a) => _get('toDisplayString')(...a);
export const mergeProps = (...a) => _get('mergeProps')(...a);
export const normalizeClass = (...a) => _get('normalizeClass')(...a);
export const normalizeStyle = (...a) => _get('normalizeStyle')(...a);
export const normalizeProps = (...a) => _get('normalizeProps')(...a);
export const guardReactiveProps = (...a) => _get('guardReactiveProps')(...a);
export const cloneVNode = (...a) => _get('cloneVNode')(...a);
export const isVNode = (...a) => _get('isVNode')(...a);

// Special non-function exports (accessed as values, not called)
export let Fragment, Teleport, Suspense, KeepAlive, Transition, TransitionGroup;
export let createApp, createSSRApp, nextTick, getCurrentInstance;
export let useSlots, useAttrs, useCssModule, useCssVars, useModel;
export let effectScope, onScopeDispose, getCurrentScope;

// Eagerly resolve value exports (these are used as identifiers, not called)
const _init = () => {
  const V = window.__nuxtblog_vue;
  if (!V) return;
  Fragment = V.Fragment; Teleport = V.Teleport; Suspense = V.Suspense;
  KeepAlive = V.KeepAlive; Transition = V.Transition; TransitionGroup = V.TransitionGroup;
  createApp = V.createApp; createSSRApp = V.createSSRApp;
  nextTick = V.nextTick; getCurrentInstance = V.getCurrentInstance;
  useSlots = V.useSlots; useAttrs = V.useAttrs;
  useCssModule = V.useCssModule; useCssVars = V.useCssVars; useModel = V.useModel;
  effectScope = V.effectScope; onScopeDispose = V.onScopeDispose; getCurrentScope = V.getCurrentScope;
};
_init();

export default window.__nuxtblog_vue;
`

// RegisterAssetRoutes registers the static asset serving endpoint for plugins.
// assetsDir is the root directory where plugin assets are stored (e.g. "data/plugins").
func RegisterAssetRoutes(s *ghttp.Server, assetsDir string) {
	// Serve shared dependency shims (e.g. /_shared/vue.mjs)
	s.BindHandler("GET:/_shared/{filename}", func(r *ghttp.Request) {
		filename := r.Get("filename").String()
		switch filename {
		case "vue.mjs":
			r.Response.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			r.Response.Header().Set("Cache-Control", "public, max-age=86400")
			r.Response.Write(vueShimJS)
		default:
			r.Response.WriteHeader(http.StatusNotFound)
		}
	})

	s.BindHandler("GET:/api/plugins/{id}/assets/{filename}", func(r *ghttp.Request) {
		pluginID := r.Get("id").String()
		filename := r.Get("filename").String()

		if pluginID == "" || filename == "" {
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate extension.
		ext := strings.ToLower(filepath.Ext(filename))
		if !allowedAssetExts[ext] {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}

		// Prevent path traversal.
		cleanName := filepath.Base(filename)
		if cleanName != filename || strings.Contains(filename, "..") {
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}

		// Sanitize plugin ID to directory name.
		dirName := strings.ReplaceAll(pluginID, "/", "--")
		filePath := filepath.Join(assetsDir, dirName, cleanName)

		data, err := os.ReadFile(filePath)
		if err != nil {
			r.Response.WriteHeader(http.StatusNotFound)
			return
		}

		// Set content type and caching headers.
		if ct, ok := extContentType[ext]; ok {
			r.Response.Header().Set("Content-Type", ct)
		}
		r.Response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		r.Response.Write(data)
	})
}
