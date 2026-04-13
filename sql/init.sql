-- Active: 1773301732955@@127.0.0.1@3306
-- ============================================================
--  Modern Go Blog — SQLite Schema
--  Generated for SQLite 3.x
--  Usage: sqlite3 blog.db < schema.sql
-- ============================================================

PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA synchronous = NORMAL;

-- ============================================================
--  SYSTEM
-- ============================================================

CREATE TABLE IF NOT EXISTS options (
    key         TEXT        NOT NULL PRIMARY KEY,
    value       TEXT        NOT NULL DEFAULT '{}',   -- JSON
    autoload    INTEGER     NOT NULL DEFAULT 1,       -- 1=load on startup
    updated_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

-- ============================================================
--  USERS
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id              INTEGER     PRIMARY KEY,          -- Snowflake ID (assigned by app)
    username        TEXT        NOT NULL UNIQUE,
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    display_name    TEXT        NOT NULL DEFAULT '',
    avatar_id       INTEGER,                          -- FK → medias.id (nullable)
    bio             TEXT        NOT NULL DEFAULT '',
    role            INTEGER     NOT NULL DEFAULT 1,   -- 1=subscriber 2=editor 3=admin
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=active 2=banned 3=pending
    email_verified  INTEGER     NOT NULL DEFAULT 0,   -- 0=unverified 1=verified
    locale          TEXT        NOT NULL DEFAULT 'zh-CN',
    last_login_at   DATETIME,
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at      DATETIME
);

CREATE INDEX IF NOT EXISTS idx_users_role       ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_status     ON users (status);
CREATE INDEX IF NOT EXISTS idx_users_deleted    ON users (deleted_at);

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id             INTEGER     PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    website             TEXT        NOT NULL DEFAULT '',
    twitter             TEXT        NOT NULL DEFAULT '',
    github              TEXT        NOT NULL DEFAULT '',
    location            TEXT        NOT NULL DEFAULT '',
    social_links        TEXT        NOT NULL DEFAULT '{}',   -- JSON
    notification_prefs  TEXT        NOT NULL DEFAULT '{}',   -- JSON
    checkin_streak      INTEGER     NOT NULL DEFAULT 0,
    last_checkin_date   TEXT        NOT NULL DEFAULT ''
);

-- ============================================================
--  MEDIA
-- ============================================================

CREATE TABLE IF NOT EXISTS medias (
    id              INTEGER     PRIMARY KEY,
    uploader_id     INTEGER     NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    storage_type    INTEGER     NOT NULL DEFAULT 1,   -- 1=local 2=s3 3=oss 4=cos
    storage_key     TEXT        NOT NULL,
    cdn_url         TEXT        NOT NULL DEFAULT '',
    filename        TEXT        NOT NULL,
    mime_type       TEXT        NOT NULL DEFAULT '',
    file_size       INTEGER     NOT NULL DEFAULT 0,   -- bytes
    width           INTEGER,
    height          INTEGER,
    duration        INTEGER,                          -- seconds, for video/audio
    alt_text        TEXT        NOT NULL DEFAULT '',
    title           TEXT        NOT NULL DEFAULT '',
    category        TEXT        NOT NULL DEFAULT 'post', -- user|post|doc|moment|banner
    variants        TEXT        NOT NULL DEFAULT '{}',-- JSON {thumb, medium, large}
    file_meta       TEXT        NOT NULL DEFAULT '{}',-- JSON EXIF etc.
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at      DATETIME
);

CREATE INDEX IF NOT EXISTS idx_medias_uploader  ON medias (uploader_id);
CREATE INDEX IF NOT EXISTS idx_medias_deleted   ON medias (deleted_at);

-- ============================================================
--  POSTS (core)
-- ============================================================

CREATE TABLE IF NOT EXISTS posts (
    id              INTEGER     PRIMARY KEY,
    post_type       INTEGER     NOT NULL DEFAULT 1,   -- 1=post 2=page 3=custom
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=draft 2=published 3=private 4=archived
    title           TEXT        NOT NULL,
    slug            TEXT        NOT NULL UNIQUE,
    content         TEXT        NOT NULL DEFAULT '',
    excerpt         TEXT        NOT NULL DEFAULT '',
    author_id       INTEGER     NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    featured_img_id INTEGER     REFERENCES medias (id),
    comment_status  INTEGER     NOT NULL DEFAULT 1,   -- 0=closed 1=open
    password_hash   TEXT,                             -- nullable, password-protected posts
    locale          TEXT        NOT NULL DEFAULT 'zh-CN',
    published_at    DATETIME,
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at      DATETIME
);

CREATE INDEX IF NOT EXISTS idx_posts_author         ON posts (author_id);
CREATE INDEX IF NOT EXISTS idx_posts_slug           ON posts (slug);
CREATE INDEX IF NOT EXISTS idx_posts_status         ON posts (status);
CREATE INDEX IF NOT EXISTS idx_posts_type_status    ON posts (post_type, status, published_at);
CREATE INDEX IF NOT EXISTS idx_posts_deleted        ON posts (deleted_at);

-- stats: split out to avoid hot-row contention on the main table
CREATE TABLE IF NOT EXISTS post_stats (
    post_id         INTEGER     PRIMARY KEY REFERENCES posts (id) ON DELETE CASCADE,
    view_count      INTEGER     NOT NULL DEFAULT 0,
    like_count      INTEGER     NOT NULL DEFAULT 0,
    comment_count   INTEGER     NOT NULL DEFAULT 0,
    share_count     INTEGER     NOT NULL DEFAULT 0,
    updated_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

-- SEO: optional 1:1, only written when needed
CREATE TABLE IF NOT EXISTS post_seo (
    post_id         INTEGER     PRIMARY KEY REFERENCES posts (id) ON DELETE CASCADE,
    meta_title      TEXT        NOT NULL DEFAULT '',
    meta_desc       TEXT        NOT NULL DEFAULT '',
    og_title        TEXT        NOT NULL DEFAULT '',
    og_image        TEXT        NOT NULL DEFAULT '',
    canonical_url   TEXT        NOT NULL DEFAULT '',
    robots          TEXT        NOT NULL DEFAULT 'index,follow',
    structured_data TEXT        NOT NULL DEFAULT '{}'  -- JSON-LD
);

-- metas: flexible key-value metadata per post
CREATE TABLE IF NOT EXISTS post_metas (
    id          INTEGER     PRIMARY KEY,
    post_id     INTEGER     NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    meta_key    TEXT        NOT NULL,
    meta_value  TEXT        NOT NULL DEFAULT '',
    created_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    UNIQUE (post_id, meta_key)
);

CREATE INDEX IF NOT EXISTS idx_post_metas_post    ON post_metas (post_id);
CREATE INDEX IF NOT EXISTS idx_post_metas_key     ON post_metas (meta_key);

-- revisions: snapshot on every publish
CREATE TABLE IF NOT EXISTS post_revisions (
    id          INTEGER     PRIMARY KEY,
    post_id     INTEGER     NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    author_id   INTEGER     NOT NULL REFERENCES users (id),
    title       TEXT        NOT NULL DEFAULT '',
    content     TEXT        NOT NULL DEFAULT '',
    rev_note    TEXT        NOT NULL DEFAULT '',
    created_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_revisions_post   ON post_revisions (post_id);

-- ============================================================
--  DOCS
-- ============================================================

CREATE TABLE IF NOT EXISTS doc_collections (
    id           INTEGER  PRIMARY KEY,
    slug         TEXT     NOT NULL UNIQUE,
    title        TEXT     NOT NULL,
    description  TEXT     NOT NULL DEFAULT '',
    cover_img_id INTEGER  REFERENCES medias (id),
    author_id    INTEGER  NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status       INTEGER  NOT NULL DEFAULT 1,    -- 1=draft 2=published
    locale       TEXT     NOT NULL DEFAULT 'zh-CN',
    sort_order   INTEGER  NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at   DATETIME
);

CREATE INDEX IF NOT EXISTS idx_doc_collections_slug    ON doc_collections (slug);
CREATE INDEX IF NOT EXISTS idx_doc_collections_author  ON doc_collections (author_id);
CREATE INDEX IF NOT EXISTS idx_doc_collections_deleted ON doc_collections (deleted_at);

-- docs: articles within a collection, optionally nested via parent_id (chapters)
CREATE TABLE IF NOT EXISTS docs (
    id             INTEGER  PRIMARY KEY,
    collection_id  INTEGER  NOT NULL REFERENCES doc_collections (id),
    parent_id      INTEGER  REFERENCES docs (id),
    sort_order     INTEGER  NOT NULL DEFAULT 0,
    status         INTEGER  NOT NULL DEFAULT 1,    -- 1=draft 2=published 3=archived
    title          TEXT     NOT NULL,
    slug           TEXT     NOT NULL UNIQUE,
    content        TEXT     NOT NULL DEFAULT '',
    excerpt        TEXT     NOT NULL DEFAULT '',
    author_id      INTEGER  NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    comment_status INTEGER  NOT NULL DEFAULT 1,    -- 0=closed 1=open
    locale         TEXT     NOT NULL DEFAULT 'zh-CN',
    published_at   DATETIME,
    created_at     DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at     DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at     DATETIME
);

CREATE INDEX IF NOT EXISTS idx_docs_collection ON docs (collection_id);
CREATE INDEX IF NOT EXISTS idx_docs_parent     ON docs (parent_id);
CREATE INDEX IF NOT EXISTS idx_docs_slug       ON docs (slug);
CREATE INDEX IF NOT EXISTS idx_docs_author     ON docs (author_id);
CREATE INDEX IF NOT EXISTS idx_docs_status     ON docs (status, published_at);
CREATE INDEX IF NOT EXISTS idx_docs_deleted    ON docs (deleted_at);

CREATE TABLE IF NOT EXISTS doc_stats (
    doc_id        INTEGER  PRIMARY KEY REFERENCES docs (id) ON DELETE CASCADE,
    view_count    INTEGER  NOT NULL DEFAULT 0,
    like_count    INTEGER  NOT NULL DEFAULT 0,
    comment_count INTEGER  NOT NULL DEFAULT 0,
    updated_at    DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS doc_seo (
    doc_id          INTEGER  PRIMARY KEY REFERENCES docs (id) ON DELETE CASCADE,
    meta_title      TEXT     NOT NULL DEFAULT '',
    meta_desc       TEXT     NOT NULL DEFAULT '',
    og_title        TEXT     NOT NULL DEFAULT '',
    og_image        TEXT     NOT NULL DEFAULT '',
    canonical_url   TEXT     NOT NULL DEFAULT '',
    robots          TEXT     NOT NULL DEFAULT 'index,follow',
    structured_data TEXT     NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS doc_revisions (
    id         INTEGER  PRIMARY KEY,
    doc_id     INTEGER  NOT NULL REFERENCES docs (id) ON DELETE CASCADE,
    author_id  INTEGER  NOT NULL REFERENCES users (id),
    title      TEXT     NOT NULL DEFAULT '',
    content    TEXT     NOT NULL DEFAULT '',
    rev_note   TEXT     NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_doc_revisions_doc ON doc_revisions (doc_id);

-- ============================================================
--  MOMENTS
-- ============================================================

-- moments: short-form social content — no title, no slug, no SEO
CREATE TABLE IF NOT EXISTS moments (
    id         INTEGER  PRIMARY KEY,
    author_id  INTEGER  NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    content    TEXT     NOT NULL,
    visibility INTEGER  NOT NULL DEFAULT 1,    -- 1=public 2=private 3=followers
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_moments_author  ON moments (author_id);
CREATE INDEX IF NOT EXISTS idx_moments_vis     ON moments (visibility, created_at);
CREATE INDEX IF NOT EXISTS idx_moments_deleted ON moments (deleted_at);

CREATE TABLE IF NOT EXISTS moment_stats (
    moment_id     INTEGER  PRIMARY KEY REFERENCES moments (id) ON DELETE CASCADE,
    view_count    INTEGER  NOT NULL DEFAULT 0,
    like_count    INTEGER  NOT NULL DEFAULT 0,
    comment_count INTEGER  NOT NULL DEFAULT 0,
    updated_at    DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

-- moment_media: ordered media attachments (images / videos) for a moment
CREATE TABLE IF NOT EXISTS moment_media (
    moment_id  INTEGER NOT NULL REFERENCES moments (id) ON DELETE CASCADE,
    media_id   INTEGER NOT NULL REFERENCES medias (id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (moment_id, media_id)
);

CREATE INDEX IF NOT EXISTS idx_moment_media_moment ON moment_media (moment_id);

-- ============================================================
--  TAXONOMY
-- ============================================================

CREATE TABLE IF NOT EXISTS terms (
    id          INTEGER     PRIMARY KEY,
    name        TEXT        NOT NULL,
    slug        TEXT        NOT NULL UNIQUE,
    created_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS taxonomies (
    id          INTEGER     PRIMARY KEY,
    term_id     INTEGER     NOT NULL REFERENCES terms (id) ON DELETE CASCADE,
    taxonomy    TEXT        NOT NULL,                 -- category / tag / topic …
    description TEXT        NOT NULL DEFAULT '',
    parent_id   INTEGER     REFERENCES taxonomies (id),
    post_count  INTEGER     NOT NULL DEFAULT 0,
    extra       TEXT        NOT NULL DEFAULT '{}',    -- JSON (cover image, color …)
    UNIQUE (term_id, taxonomy)
);

CREATE INDEX IF NOT EXISTS idx_taxonomies_term      ON taxonomies (term_id);
CREATE INDEX IF NOT EXISTS idx_taxonomies_taxonomy  ON taxonomies (taxonomy);
CREATE INDEX IF NOT EXISTS idx_taxonomies_parent    ON taxonomies (parent_id);

-- polymorphic join table: posts, products, videos … all use this one table
CREATE TABLE IF NOT EXISTS object_taxonomies (
    object_id       INTEGER     NOT NULL,
    object_type     TEXT        NOT NULL DEFAULT 'post', -- post / product / video …
    taxonomy_id     INTEGER     NOT NULL REFERENCES taxonomies (id) ON DELETE CASCADE,
    sort_order      INTEGER     NOT NULL DEFAULT 0,
    PRIMARY KEY (object_id, object_type, taxonomy_id)
);

CREATE INDEX IF NOT EXISTS idx_obj_tax_object   ON object_taxonomies (object_id, object_type);
CREATE INDEX IF NOT EXISTS idx_obj_tax_taxonomy ON object_taxonomies (taxonomy_id);

-- ============================================================
--  COMMENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS comments (
    id              INTEGER     PRIMARY KEY,
    object_id       INTEGER     NOT NULL,             -- polymorphic
    object_type     TEXT        NOT NULL DEFAULT 'post',
    parent_id       INTEGER     REFERENCES comments (id) ON DELETE SET NULL,
    user_id         INTEGER     REFERENCES users (id),-- nullable for guests
    author_name     TEXT        NOT NULL DEFAULT '',
    author_email    TEXT        NOT NULL DEFAULT '',
    content         TEXT        NOT NULL,
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=pending 2=approved 3=spam 4=trash
    ip              TEXT        NOT NULL DEFAULT '',
    user_agent      TEXT        NOT NULL DEFAULT '',
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at      DATETIME
);

CREATE INDEX IF NOT EXISTS idx_comments_object        ON comments (object_id, object_type);
CREATE INDEX IF NOT EXISTS idx_comments_object_status ON comments (object_id, object_type, status);
CREATE INDEX IF NOT EXISTS idx_comments_parent        ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_user          ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_status        ON comments (status);
CREATE INDEX IF NOT EXISTS idx_comments_deleted       ON comments (deleted_at);

-- ============================================================
--  TRIGGERS — keep updated_at fresh automatically
-- ============================================================

CREATE TRIGGER IF NOT EXISTS trg_users_updated
    AFTER UPDATE ON users FOR EACH ROW
    BEGIN UPDATE users SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_posts_updated
    AFTER UPDATE ON posts FOR EACH ROW
    BEGIN UPDATE posts SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_post_stats_updated
    AFTER UPDATE ON post_stats FOR EACH ROW
    BEGIN UPDATE post_stats SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE post_id = OLD.post_id; END;

CREATE TRIGGER IF NOT EXISTS trg_post_metas_updated
    AFTER UPDATE ON post_metas FOR EACH ROW
    BEGIN UPDATE post_metas SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_options_updated
    AFTER UPDATE ON options FOR EACH ROW
    BEGIN UPDATE options SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE key = OLD.key; END;

CREATE TRIGGER IF NOT EXISTS trg_doc_collections_updated
    AFTER UPDATE ON doc_collections FOR EACH ROW
    BEGIN UPDATE doc_collections SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_docs_updated
    AFTER UPDATE ON docs FOR EACH ROW
    BEGIN UPDATE docs SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_doc_stats_updated
    AFTER UPDATE ON doc_stats FOR EACH ROW
    BEGIN UPDATE doc_stats SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE doc_id = OLD.doc_id; END;

CREATE TRIGGER IF NOT EXISTS trg_moments_updated
    AFTER UPDATE ON moments FOR EACH ROW
    BEGIN UPDATE moments SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_moment_stats_updated
    AFTER UPDATE ON moment_stats FOR EACH ROW
    BEGIN UPDATE moment_stats SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE moment_id = OLD.moment_id; END;

-- ============================================================
--  NOTIFICATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS notifications (
    id           INTEGER  PRIMARY KEY,
    user_id      INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         TEXT     NOT NULL,
    sub_type     TEXT     NOT NULL DEFAULT '',
    actor_id     INTEGER,
    actor_name   TEXT     NOT NULL DEFAULT '',
    actor_avatar TEXT     NOT NULL DEFAULT '',
    object_type  TEXT     NOT NULL DEFAULT '',
    object_id    INTEGER,
    object_title TEXT     NOT NULL DEFAULT '',
    object_link  TEXT     NOT NULL DEFAULT '',
    content      TEXT     NOT NULL DEFAULT '',
    is_read      INTEGER  NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at   DATETIME
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id      ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read    ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at);

-- user_likes / user_bookmarks: polymorphic — covers post | doc | moment
CREATE TABLE IF NOT EXISTS user_likes (
    user_id     INTEGER NOT NULL,
    object_type TEXT    NOT NULL DEFAULT 'post',
    object_id   INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    PRIMARY KEY (user_id, object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_user_likes_object ON user_likes (object_type, object_id);

CREATE TABLE IF NOT EXISTS user_bookmarks (
    user_id     INTEGER NOT NULL,
    object_type TEXT    NOT NULL DEFAULT 'post',
    object_id   INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    PRIMARY KEY (user_id, object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_user_bookmarks_object ON user_bookmarks (object_type, object_id);

-- Generic user action log: checkin / share / download / login / etc.
CREATE TABLE IF NOT EXISTS user_actions (
    id          INTEGER  PRIMARY KEY,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action      TEXT     NOT NULL,
    object_type TEXT     NOT NULL DEFAULT '',
    object_id   INTEGER  NOT NULL DEFAULT 0,
    extra       TEXT     NOT NULL DEFAULT '{}',
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_user_actions_user_action ON user_actions (user_id, action);
CREATE INDEX IF NOT EXISTS idx_user_actions_created     ON user_actions (user_id, created_at);

-- Personal access tokens (API keys)
CREATE TABLE IF NOT EXISTS user_tokens (
    id           INTEGER  PRIMARY KEY,
    user_id      INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT     NOT NULL,
    prefix       TEXT     NOT NULL,
    token_hash   TEXT     NOT NULL UNIQUE,
    expires_at   DATETIME,
    last_used_at DATETIME,
    created_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_user_tokens_user ON user_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_user_tokens_hash ON user_tokens (token_hash);

-- Follow graph
CREATE TABLE IF NOT EXISTS user_follows (
    follower_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follows_follower  ON user_follows (follower_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_following ON user_follows (following_id);

-- Email/SMS verification codes
CREATE TABLE IF NOT EXISTS verification_codes (
    id         INTEGER  PRIMARY KEY,
    target     TEXT     NOT NULL,
    code       TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at    DATETIME,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_vcodes_target  ON verification_codes (target, type);
CREATE INDEX IF NOT EXISTS idx_vcodes_expires ON verification_codes (expires_at);

-- OAuth provider account links
CREATE TABLE IF NOT EXISTS user_oauth (
    id          INTEGER  PRIMARY KEY,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider    TEXT     NOT NULL,
    provider_id TEXT     NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    UNIQUE(provider, provider_id)
);

CREATE INDEX IF NOT EXISTS idx_user_oauth_user ON user_oauth(user_id);

-- Reports
CREATE TABLE IF NOT EXISTS reports (
    id          INTEGER  PRIMARY KEY,
    reporter_id INTEGER  REFERENCES users(id) ON DELETE SET NULL,
    target_type TEXT     NOT NULL, -- 'post' | 'comment' | 'user'
    target_id   INTEGER  NOT NULL,
    reason      TEXT     NOT NULL,
    detail      TEXT     NOT NULL DEFAULT '',
    status      TEXT     NOT NULL DEFAULT 'pending', -- 'pending' | 'resolved' | 'dismissed'
    notes       TEXT     NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    resolved_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_reporter ON reports(reporter_id);
CREATE INDEX IF NOT EXISTS idx_reports_target ON reports(target_type, target_id);

-- Private messages
-- id is application-assigned (random int64 via idgen) to avoid sequential enumeration.
CREATE TABLE IF NOT EXISTS conversations (
    id           INTEGER  PRIMARY KEY,
    user_a       INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_b       INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_msg     TEXT     NOT NULL DEFAULT '',
    last_msg_at  DATETIME,
    unread_a     INTEGER  NOT NULL DEFAULT 0,
    unread_b     INTEGER  NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    UNIQUE(user_a, user_b)
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_a ON conversations(user_a);
CREATE INDEX IF NOT EXISTS idx_conversations_user_b ON conversations(user_b);

-- id is application-assigned (random int64 via idgen) to avoid sequential enumeration.
CREATE TABLE IF NOT EXISTS messages (
    id              INTEGER  PRIMARY KEY,
    conversation_id INTEGER  NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT     NOT NULL,
    is_read         INTEGER  NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender       ON messages(sender_id);

-- Plugins
CREATE TABLE IF NOT EXISTS plugins (
    id           TEXT     NOT NULL PRIMARY KEY,  -- e.g. "owner/repo"
    title        TEXT     NOT NULL DEFAULT '',
    description  TEXT     NOT NULL DEFAULT '',
    version      TEXT     NOT NULL DEFAULT '',
    author       TEXT     NOT NULL DEFAULT '',
    icon         TEXT     NOT NULL DEFAULT 'i-tabler-plug',
    repo_url     TEXT     NOT NULL DEFAULT '',
    script       TEXT     NOT NULL DEFAULT '',   -- bundled JS entry
    styles       TEXT     NOT NULL DEFAULT '',   -- optional CSS injected into frontend <head>
    priority     INTEGER  NOT NULL DEFAULT 10,   -- execution order: lower runs first
    settings        TEXT     NOT NULL DEFAULT '{}', -- JSON user config (key→value)
    settings_schema TEXT     NOT NULL DEFAULT '[]', -- JSON settings schema from package.json "plugin".settings
    capabilities    TEXT     NOT NULL DEFAULT '{}', -- JSON capability declarations from package.json "plugin".capabilities
    manifest        TEXT     NOT NULL DEFAULT '{}', -- JSON full manifest (pipelines, webhooks, settings schema)
    source          TEXT     NOT NULL DEFAULT 'external', -- 'builtin' (Go native) | 'external' (installed via zip/github)
    enabled         INTEGER  NOT NULL DEFAULT 1,
    installed_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at   DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

-- P-B10: add manifest column for existing databases (silently ignored if column already exists)
ALTER TABLE plugins ADD COLUMN manifest TEXT NOT NULL DEFAULT '{}';

-- Add source column for existing databases
ALTER TABLE plugins ADD COLUMN source TEXT NOT NULL DEFAULT 'external';

-- ============================================================
--  ANNOUNCEMENTS (broadcast notifications)
-- ============================================================

CREATE TABLE IF NOT EXISTS announcements (
    id          INTEGER  PRIMARY KEY,
    title       TEXT     NOT NULL,
    content     TEXT     NOT NULL DEFAULT '',
    type        TEXT     NOT NULL DEFAULT 'info',
    created_by  INTEGER  REFERENCES users(id) ON DELETE SET NULL,
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    deleted_at  DATETIME
);

CREATE INDEX IF NOT EXISTS idx_announcements_created ON announcements(created_at);

-- One row per user: stores the timestamp of their last "mark all read" action.
-- Announcements created after last_read_at are considered unread for that user.
CREATE TABLE IF NOT EXISTS user_announcement_cursor (
    user_id      INTEGER  NOT NULL PRIMARY KEY,
    last_read_at DATETIME NOT NULL
);

-- ============================================================
--  COMMERCE: ORDERS
-- ============================================================

CREATE TABLE IF NOT EXISTS orders (
    id              INTEGER     PRIMARY KEY,
    order_no        TEXT        NOT NULL UNIQUE,
    user_id         INTEGER     NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=pending 2=paid 3=completed 4=cancelled 5=refunded
    total_amount    INTEGER     NOT NULL DEFAULT 0,   -- cents
    paid_amount     INTEGER     NOT NULL DEFAULT 0,   -- cents
    credits_used    INTEGER     NOT NULL DEFAULT 0,
    balance_used    INTEGER     NOT NULL DEFAULT 0,   -- cents
    currency        TEXT        NOT NULL DEFAULT 'CNY',
    expires_at      DATETIME,
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_orders_user     ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status   ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders (order_no);

CREATE TABLE IF NOT EXISTS order_items (
    id          INTEGER     PRIMARY KEY,
    order_id    INTEGER     NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    item_type   TEXT        NOT NULL,   -- 'post_unlock','download','membership','product','topup'
    item_id     TEXT        NOT NULL DEFAULT '',
    title       TEXT        NOT NULL DEFAULT '',
    unit_price  INTEGER     NOT NULL DEFAULT 0,   -- cents
    quantity    INTEGER     NOT NULL DEFAULT 1,
    snapshot    TEXT        NOT NULL DEFAULT '{}'  -- JSON
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items (order_id);

CREATE TABLE IF NOT EXISTS transactions (
    id              INTEGER     PRIMARY KEY,
    order_id        INTEGER     NOT NULL REFERENCES orders (id) ON DELETE RESTRICT,
    user_id         INTEGER     NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    type            INTEGER     NOT NULL DEFAULT 1,   -- 1=payment 2=refund 3=topup
    provider        TEXT        NOT NULL DEFAULT '',   -- 'alipay','paypal','balance','credits'
    provider_tx_id  TEXT        NOT NULL DEFAULT '',
    amount          INTEGER     NOT NULL DEFAULT 0,   -- cents
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=pending 2=success 3=failed
    raw_notify      TEXT        NOT NULL DEFAULT '',
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_transactions_order    ON transactions (order_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user     ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_provider ON transactions (provider, provider_tx_id);

-- ============================================================
--  COMMERCE: WALLET
-- ============================================================

CREATE TABLE IF NOT EXISTS user_wallets (
    user_id     INTEGER     PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    balance     INTEGER     NOT NULL DEFAULT 0,   -- cents
    frozen      INTEGER     NOT NULL DEFAULT 0,   -- cents
    total_topup INTEGER     NOT NULL DEFAULT 0,   -- cents
    total_spent INTEGER     NOT NULL DEFAULT 0,   -- cents
    updated_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS wallet_ledger (
    id              INTEGER     PRIMARY KEY,
    user_id         INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type            INTEGER     NOT NULL,   -- 1=topup 2=spend 3=refund 4=admin_adjust
    amount          INTEGER     NOT NULL,   -- signed: positive=credit, negative=debit
    balance_after   INTEGER     NOT NULL,
    reference_type  TEXT        NOT NULL DEFAULT '',   -- 'order','admin'
    reference_id    TEXT        NOT NULL DEFAULT '',
    note            TEXT        NOT NULL DEFAULT '',
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_wallet_ledger_user ON wallet_ledger (user_id, created_at);

-- ============================================================
--  COMMERCE: CREDITS
-- ============================================================

CREATE TABLE IF NOT EXISTS user_credits (
    user_id      INTEGER     PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    balance      INTEGER     NOT NULL DEFAULT 0,
    total_earned INTEGER     NOT NULL DEFAULT 0,
    total_spent  INTEGER     NOT NULL DEFAULT 0,
    updated_at   DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS credits_ledger (
    id              INTEGER     PRIMARY KEY,
    user_id         INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type            INTEGER     NOT NULL,   -- 1=earn 2=spend 3=expire 4=admin_adjust
    amount          INTEGER     NOT NULL,   -- signed
    balance_after   INTEGER     NOT NULL,
    source          TEXT        NOT NULL DEFAULT '',   -- 'checkin','comment','purchase','admin'
    reference_type  TEXT        NOT NULL DEFAULT '',
    reference_id    TEXT        NOT NULL DEFAULT '',
    note            TEXT        NOT NULL DEFAULT '',
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_credits_ledger_user ON credits_ledger (user_id, created_at);

-- ============================================================
--  COMMERCE: MEMBERSHIP
-- ============================================================

CREATE TABLE IF NOT EXISTS membership_tiers (
    id              INTEGER     PRIMARY KEY,
    name            TEXT        NOT NULL,
    slug            TEXT        NOT NULL UNIQUE,
    description     TEXT        NOT NULL DEFAULT '',
    price           INTEGER     NOT NULL DEFAULT 0,   -- cents
    duration_days   INTEGER     NOT NULL DEFAULT 30,
    discount_pct    INTEGER     NOT NULL DEFAULT 0,   -- 0-100
    access_all      INTEGER     NOT NULL DEFAULT 0,   -- 1=all paid content free
    credits_monthly INTEGER     NOT NULL DEFAULT 0,
    features        TEXT        NOT NULL DEFAULT '[]', -- JSON
    status          INTEGER     NOT NULL DEFAULT 1,   -- 1=active 2=disabled
    sort_order      INTEGER     NOT NULL DEFAULT 0,
    created_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at      DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS user_memberships (
    id          INTEGER     PRIMARY KEY,
    user_id     INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tier_id     INTEGER     NOT NULL REFERENCES membership_tiers (id) ON DELETE RESTRICT,
    status      INTEGER     NOT NULL DEFAULT 1,   -- 1=active 2=expired 3=cancelled
    started_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    expires_at  DATETIME    NOT NULL,
    auto_renew  INTEGER     NOT NULL DEFAULT 0,
    order_id    INTEGER     REFERENCES orders (id),
    created_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_user_memberships_user   ON user_memberships (user_id, status);
CREATE INDEX IF NOT EXISTS idx_user_memberships_expire ON user_memberships (expires_at);

-- ============================================================
--  COMMERCE: PURCHASE RECORDS
-- ============================================================

CREATE TABLE IF NOT EXISTS user_purchases (
    user_id     INTEGER     NOT NULL,
    object_type TEXT        NOT NULL,   -- 'post','download','product'
    object_id   TEXT        NOT NULL,
    order_id    INTEGER     REFERENCES orders (id),
    created_at  DATETIME    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    PRIMARY KEY (user_id, object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_user_purchases_user ON user_purchases (user_id);

-- ============================================================
--  COMMERCE: TRIGGERS
-- ============================================================

CREATE TRIGGER IF NOT EXISTS trg_orders_updated
    AFTER UPDATE ON orders FOR EACH ROW
    BEGIN UPDATE orders SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS trg_user_wallets_updated
    AFTER UPDATE ON user_wallets FOR EACH ROW
    BEGIN UPDATE user_wallets SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE user_id = OLD.user_id; END;

CREATE TRIGGER IF NOT EXISTS trg_user_credits_updated
    AFTER UPDATE ON user_credits FOR EACH ROW
    BEGIN UPDATE user_credits SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE user_id = OLD.user_id; END;

CREATE TRIGGER IF NOT EXISTS trg_membership_tiers_updated
    AFTER UPDATE ON membership_tiers FOR EACH ROW
    BEGIN UPDATE membership_tiers SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE id = OLD.id; END;

-- ============================================================
--  DEFAULT DATA
--  Default admin user is seeded automatically by autoMigrate()
--  on first startup: username=admin  password=admin123  role=3
-- ============================================================
