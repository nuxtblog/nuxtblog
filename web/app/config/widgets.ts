import { z } from 'zod'

/**
 * Widget Registry — validated by Zod.
 * TypeScript types are inferred from the schemas, not declared separately.
 */

// ── Schemas ────────────────────────────────────────────────────────────────────

export const WidgetConfigFieldSchema = z.enum(['maxCount', 'showRecent', 'showHot'])

export const WidgetFieldDefaultsSchema = z.object({
  maxCount:   z.number().int().positive().optional(),
  showRecent: z.boolean().optional(),
  showHot:    z.boolean().optional(),
})

export const WidgetDefSchema = z.object({
  id:             z.string(),
  label:          z.string(),
  defaultEnabled: z.boolean(),
  defaultTitle:   z.string().optional(),
  configFields:   z.array(WidgetConfigFieldSchema),
  fieldDefaults:  WidgetFieldDefaultsSchema.optional(),
  /**
   * Which sidebar contexts this widget is available in.
   * Omit = available in all contexts (homepage, blog, post).
   */
  context:        z.array(z.enum(['homepage', 'blog', 'post'])).optional(),
})

// ── Inferred types ─────────────────────────────────────────────────────────────

export type WidgetConfigField = z.infer<typeof WidgetConfigFieldSchema>
export type WidgetDef         = z.infer<typeof WidgetDefSchema>

// ── Registry ───────────────────────────────────────────────────────────────────

const REGISTRY_RAW = [
  {
    // UserBox: login card for guests, mini-profile for logged-in users — homepage only
    id: 'user_box', label: 'admin.widgets.id_user_box', defaultEnabled: true,
    configFields: [],
    context: ['homepage'],
  },
  {
    // AuthorBox: site owner profile card — blog/post sidebars only
    id: 'author', label: 'admin.widgets.id_author', defaultEnabled: true,
    configFields: [],
    context: ['blog', 'post'],
  },
  {
    id: 'search', label: 'admin.widgets.id_search', defaultEnabled: true, defaultTitle: 'admin.widgets.id_search',
    configFields: ['showRecent', 'showHot'],
    fieldDefaults: { showRecent: true, showHot: true },
  },
  {
    id: 'tags', label: 'admin.widgets.id_tags', defaultEnabled: true, defaultTitle: 'admin.widgets.id_tags',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 15 },
  },
  {
    id: 'latest_posts', label: 'admin.widgets.id_latest_posts', defaultEnabled: true, defaultTitle: 'admin.widgets.id_latest_posts',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 5 },
  },
  {
    id: 'latest_comments', label: 'admin.widgets.id_latest_comments', defaultEnabled: true, defaultTitle: 'admin.widgets.id_latest_comments',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 5 },
  },
  {
    id: 'recommend', label: 'admin.widgets.id_recommend', defaultEnabled: true, defaultTitle: 'admin.widgets.id_recommend',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 5 },
  },
  {
    id: 'featured', label: 'admin.widgets.id_featured', defaultEnabled: true, defaultTitle: 'admin.widgets.id_featured',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 5 },
  },
  {
    id: 'random_posts', label: 'admin.widgets.id_random_posts', defaultEnabled: true, defaultTitle: 'admin.widgets.id_random_posts',
    configFields: ['maxCount'], fieldDefaults: { maxCount: 5 },
  },
  {
    id: 'toc', label: 'admin.widgets.id_toc', defaultEnabled: true, defaultTitle: 'admin.widgets.id_toc',
    configFields: [],
    context: ['post'],
  },
  {
    id: 'downloads', label: 'admin.widgets.id_downloads', defaultEnabled: true, defaultTitle: 'admin.widgets.id_downloads',
    configFields: [],
    context: ['post'],
  },
] as const

// Parse and validate the registry at module load time.
// This will throw at startup if any entry is malformed.
export const WIDGET_REGISTRY: WidgetDef[] = z.array(WidgetDefSchema).parse(REGISTRY_RAW)

/** Convenience lookup */
export const getWidgetDef = (id: string): WidgetDef | undefined =>
  WIDGET_REGISTRY.find((w) => w.id === id)

/** Filter registry by context */
export const getWidgetsByContext = (ctx: 'homepage' | 'blog' | 'post'): WidgetDef[] =>
  WIDGET_REGISTRY.filter(w => !w.context || w.context.includes(ctx))
