package pluginsys

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// ─── Plugin Asset Serving ───────────────────────────────────────────────────
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
export const createSlots = (...a) => _get('createSlots')(...a);

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

// uiShimJS is a tiny ES module that re-exports Nuxt UI components from the host's window global.
// Plugin ESM builds import from '/_shared/ui.mjs' which resolves to this shim.
const uiShimJS = `// Auto-generated UI shim — re-exports from host window.__nuxtblog_ui
// Uses lazy Proxy so components resolve at render time, not at module parse time.
// This avoids timing issues where the shim loads before plugin-shared-deps runs.
function _lazy(name) {
  let cached;
  return new Proxy({}, {
    get(_, prop) {
      if (!cached) {
        const U = window.__nuxtblog_ui;
        if (!U) throw new Error('[plugin-shim] window.__nuxtblog_ui not found');
        cached = U[name];
        if (!cached) throw new Error('[plugin-shim] component ' + name + ' not found on window.__nuxtblog_ui');
      }
      return cached[prop];
    }
  });
}
export const UButton = _lazy('UButton');
export const UCard = _lazy('UCard');
export const UBadge = _lazy('UBadge');
export const UIcon = _lazy('UIcon');
export const USkeleton = _lazy('USkeleton');
export const UTable = _lazy('UTable');
export const UInput = _lazy('UInput');
export const USelect = _lazy('USelect');
export const UTabs = _lazy('UTabs');
export const UAlert = _lazy('UAlert');
export const USwitch = _lazy('USwitch');
export const USeparator = _lazy('USeparator');
export const UTextarea = _lazy('UTextarea');
export const UCheckbox = _lazy('UCheckbox');
export const UCheckboxGroup = _lazy('UCheckboxGroup');
export const URadioGroup = _lazy('URadioGroup');
export const UFormField = _lazy('UFormField');
export const UForm = _lazy('UForm');
export const UInputNumber = _lazy('UInputNumber');
export const USelectMenu = _lazy('USelectMenu');
export const USlider = _lazy('USlider');
export const UPinInput = _lazy('UPinInput');
export const UModal = _lazy('UModal');
export const UDrawer = _lazy('UDrawer');
export const USlideover = _lazy('USlideover');
export const UPopover = _lazy('UPopover');
export const UTooltip = _lazy('UTooltip');
export const UContextMenu = _lazy('UContextMenu');
export const UDropdownMenu = _lazy('UDropdownMenu');
export const UAvatar = _lazy('UAvatar');
export const UAvatarGroup = _lazy('UAvatarGroup');
export const UProgress = _lazy('UProgress');
export const UPagination = _lazy('UPagination');
export const UChip = _lazy('UChip');
export const UKbd = _lazy('UKbd');
export const UCollapsible = _lazy('UCollapsible');
export const UAccordion = _lazy('UAccordion');
export const ULink = _lazy('ULink');
export const UBreadcrumb = _lazy('UBreadcrumb');
export const AdminPageContainer = _lazy('AdminPageContainer');
export const AdminPageHeader = _lazy('AdminPageHeader');
export const AdminPageContent = _lazy('AdminPageContent');
export const AdminPageFooter = _lazy('AdminPageFooter');
export const AdminRichEditor = _lazy('AdminRichEditor');
export const PageHeader = _lazy('PageHeader');
export const PageContent = _lazy('PageContent');
export const PageFooter = _lazy('PageFooter');
`

// RegisterAssetRoutes registers the static asset serving endpoint for plugins.
// assetsDir is the root directory where plugin assets are stored (e.g. "data/plugins").
//
// Lookup order for /api/plugins/{id}/assets/{filename}:
//  1. In-memory: if the loaded plugin implements sdk.HasAssets, serve from its
//     Assets() map. This is required for builtin plugins whose assets live only
//     inside the compiled binary (//go:embed) and are never written to disk.
//  2. Disk fallback: read from {assetsDir}/{dirName}/{filename}. This serves
//     JS/yaegi plugins which still write their assets to data/plugins/{id}/.
func (m *Manager) RegisterAssetRoutes(s *ghttp.Server, assetsDir string) {
	// Serve shared dependency shims (e.g. /_shared/vue.mjs)
	s.BindHandler("GET:/_shared/{filename}", func(r *ghttp.Request) {
		filename := r.Get("filename").String()
		switch filename {
		case "vue.mjs":
			r.Response.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			r.Response.Header().Set("Cache-Control", "public, max-age=86400")
			r.Response.Write(vueShimJS)
		case "ui.mjs":
			r.Response.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			r.Response.Header().Set("Cache-Control", "public, max-age=86400")
			r.Response.Write(uiShimJS)
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

		// 1) In-memory lookup via sdk.HasAssets — required for builtin plugins
		//    whose assets are //go:embed'd into the binary and never on disk.
		m.mu.RLock()
		lp := m.plugins[pluginID]
		m.mu.RUnlock()
		if lp != nil && lp.plugin != nil {
			if ap, ok := lp.plugin.(sdk.HasAssets); ok {
				if assets := ap.Assets(); len(assets) > 0 {
					var memData []byte
					if d, hit := assets[cleanName]; hit {
						memData = d
					} else {
						// Defensive: some plugins use keys like "dist/admin.mjs"
						for k, d := range assets {
							if filepath.Base(k) == cleanName {
								memData = d
								break
							}
						}
					}
					if memData != nil {
						if ct, ok := extContentType[ext]; ok {
							r.Response.Header().Set("Content-Type", ct)
						}
						r.Response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
						r.Response.Write(memData)
						return
					}
				}
			}
		}

		// 2) Disk fallback for JS/yaegi plugins.
		filePath := filepath.Join(assetsDir, dirName, cleanName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			r.Response.WriteHeader(http.StatusNotFound)
			return
		}

		if ct, ok := extContentType[ext]; ok {
			r.Response.Header().Set("Content-Type", ct)
		}
		r.Response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		r.Response.Write(data)
	})
}
