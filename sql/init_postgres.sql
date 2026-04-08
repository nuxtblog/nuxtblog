-- ============================================================
--  Modern Go Blog — PostgreSQL Schema
--  Usage: psql -U user -d blog < init_postgres.sql
-- ============================================================

-- ============================================================
--  SYSTEM
-- ============================================================

CREATE TABLE IF NOT EXISTS options (
    key         TEXT        NOT NULL PRIMARY KEY,
    value       TEXT        NOT NULL DEFAULT '{}',
    autoload    INTEGER     NOT NULL DEFAULT 1,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
--  USERS
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id              BIGINT      PRIMARY KEY,            -- Snowflake ID (assigned by app)
    username        TEXT        NOT NULL UNIQUE,
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    display_name    TEXT        NOT NULL DEFAULT '',
    avatar_id       BIGINT,                             -- FK → medias.id (nullable)
    bio             TEXT        NOT NULL DEFAULT '',
    role            INTEGER     NOT NULL DEFAULT 1,     -- 1=subscriber 2=editor 3=admin
    status          INTEGER     NOT NULL DEFAULT 1,     -- 1=active 2=banned 3=pending
    email_verified  INTEGER     NOT NULL DEFAULT 0,     -- 0=unverified 1=verified
    locale          TEXT        NOT NULL DEFAULT 'zh-CN',
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_role    ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_status  ON users (status);
CREATE INDEX IF NOT EXISTS idx_users_deleted ON users (deleted_at);

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id             BIGINT      PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    website             TEXT        NOT NULL DEFAULT '',
    twitter             TEXT        NOT NULL DEFAULT '',
    github              TEXT        NOT NULL DEFAULT '',
    location            TEXT        NOT NULL DEFAULT '',
    social_links        TEXT        NOT NULL DEFAULT '{}',
    notification_prefs  TEXT        NOT NULL DEFAULT '{}',
    checkin_streak      INTEGER     NOT NULL DEFAULT 0,
    last_checkin_date   TEXT        NOT NULL DEFAULT ''
);

-- ============================================================
--  MEDIA
-- ============================================================

CREATE TABLE IF NOT EXISTS medias (
    id              BIGINT      PRIMARY KEY,
    uploader_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    storage_type    INTEGER     NOT NULL DEFAULT 1,     -- 1=local 2=s3 3=oss 4=cos
    storage_key     TEXT        NOT NULL,
    cdn_url         TEXT        NOT NULL DEFAULT '',
    filename        TEXT        NOT NULL,
    mime_type       TEXT        NOT NULL DEFAULT '',
    file_size       BIGINT      NOT NULL DEFAULT 0,
    width           INTEGER,
    height          INTEGER,
    duration        INTEGER,
    alt_text        TEXT        NOT NULL DEFAULT '',
    title           TEXT        NOT NULL DEFAULT '',
    category        TEXT        NOT NULL DEFAULT 'post', -- user|post|doc|moment|banner
    variants        TEXT        NOT NULL DEFAULT '{}',
    file_meta       TEXT        NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_medias_uploader ON medias (uploader_id);
CREATE INDEX IF NOT EXISTS idx_medias_deleted  ON medias (deleted_at);

-- ============================================================
--  POSTS
-- ============================================================

CREATE TABLE IF NOT EXISTS posts (
    id              BIGINT      PRIMARY KEY,
    post_type       INTEGER     NOT NULL DEFAULT 1,     -- 1=post 2=page 3=custom
    status          INTEGER     NOT NULL DEFAULT 1,     -- 1=draft 2=published 3=private 4=archived
    title           TEXT        NOT NULL,
    slug            TEXT        NOT NULL UNIQUE,
    content         TEXT        NOT NULL DEFAULT '',
    excerpt         TEXT        NOT NULL DEFAULT '',
    author_id       BIGINT      NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    featured_img_id BIGINT      REFERENCES medias (id),
    comment_status  INTEGER     NOT NULL DEFAULT 1,     -- 0=closed 1=open
    password_hash   TEXT,
    locale          TEXT        NOT NULL DEFAULT 'zh-CN',
    published_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_posts_author      ON posts (author_id);
CREATE INDEX IF NOT EXISTS idx_posts_slug        ON posts (slug);
CREATE INDEX IF NOT EXISTS idx_posts_status      ON posts (status);
CREATE INDEX IF NOT EXISTS idx_posts_type_status ON posts (post_type, status, published_at);
CREATE INDEX IF NOT EXISTS idx_posts_deleted     ON posts (deleted_at);

CREATE TABLE IF NOT EXISTS post_stats (
    post_id       BIGINT      PRIMARY KEY REFERENCES posts (id) ON DELETE CASCADE,
    view_count    BIGINT      NOT NULL DEFAULT 0,
    like_count    BIGINT      NOT NULL DEFAULT 0,
    comment_count BIGINT      NOT NULL DEFAULT 0,
    share_count   BIGINT      NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post_seo (
    post_id         BIGINT  PRIMARY KEY REFERENCES posts (id) ON DELETE CASCADE,
    meta_title      TEXT    NOT NULL DEFAULT '',
    meta_desc       TEXT    NOT NULL DEFAULT '',
    og_title        TEXT    NOT NULL DEFAULT '',
    og_image        TEXT    NOT NULL DEFAULT '',
    canonical_url   TEXT    NOT NULL DEFAULT '',
    robots          TEXT    NOT NULL DEFAULT 'index,follow',
    structured_data TEXT    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS post_metas (
    id          BIGINT      PRIMARY KEY,
    post_id     BIGINT      NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    meta_key    TEXT        NOT NULL,
    meta_value  TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (post_id, meta_key)
);

CREATE INDEX IF NOT EXISTS idx_post_metas_post ON post_metas (post_id);
CREATE INDEX IF NOT EXISTS idx_post_metas_key  ON post_metas (meta_key);

CREATE TABLE IF NOT EXISTS post_revisions (
    id          BIGINT      PRIMARY KEY,
    post_id     BIGINT      NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    author_id   BIGINT      NOT NULL REFERENCES users (id),
    title       TEXT        NOT NULL DEFAULT '',
    content     TEXT        NOT NULL DEFAULT '',
    rev_note    TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_revisions_post ON post_revisions (post_id);

-- ============================================================
--  DOCS
-- ============================================================

CREATE TABLE IF NOT EXISTS doc_collections (
    id           BIGINT      PRIMARY KEY,
    slug         TEXT        NOT NULL UNIQUE,
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL DEFAULT '',
    cover_img_id BIGINT      REFERENCES medias (id),
    author_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status       INTEGER     NOT NULL DEFAULT 1,    -- 1=draft 2=published
    locale       TEXT        NOT NULL DEFAULT 'zh-CN',
    sort_order   INTEGER     NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_doc_collections_slug    ON doc_collections (slug);
CREATE INDEX IF NOT EXISTS idx_doc_collections_author  ON doc_collections (author_id);
CREATE INDEX IF NOT EXISTS idx_doc_collections_deleted ON doc_collections (deleted_at);

-- docs: articles within a collection, optionally nested via parent_id (chapters)
CREATE TABLE IF NOT EXISTS docs (
    id             BIGINT      PRIMARY KEY,
    collection_id  BIGINT      NOT NULL REFERENCES doc_collections (id),
    parent_id      BIGINT      REFERENCES docs (id),
    sort_order     INTEGER     NOT NULL DEFAULT 0,
    status         INTEGER     NOT NULL DEFAULT 1,    -- 1=draft 2=published 3=archived
    title          TEXT        NOT NULL,
    slug           TEXT        NOT NULL UNIQUE,
    content        TEXT        NOT NULL DEFAULT '',
    excerpt        TEXT        NOT NULL DEFAULT '',
    author_id      BIGINT      NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    comment_status INTEGER     NOT NULL DEFAULT 1,    -- 0=closed 1=open
    locale         TEXT        NOT NULL DEFAULT 'zh-CN',
    published_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_docs_collection ON docs (collection_id);
CREATE INDEX IF NOT EXISTS idx_docs_parent     ON docs (parent_id);
CREATE INDEX IF NOT EXISTS idx_docs_slug       ON docs (slug);
CREATE INDEX IF NOT EXISTS idx_docs_author     ON docs (author_id);
CREATE INDEX IF NOT EXISTS idx_docs_status     ON docs (status, published_at);
CREATE INDEX IF NOT EXISTS idx_docs_deleted    ON docs (deleted_at);

CREATE TABLE IF NOT EXISTS doc_stats (
    doc_id        BIGINT      PRIMARY KEY REFERENCES docs (id) ON DELETE CASCADE,
    view_count    BIGINT      NOT NULL DEFAULT 0,
    like_count    BIGINT      NOT NULL DEFAULT 0,
    comment_count BIGINT      NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS doc_seo (
    doc_id          BIGINT  PRIMARY KEY REFERENCES docs (id) ON DELETE CASCADE,
    meta_title      TEXT    NOT NULL DEFAULT '',
    meta_desc       TEXT    NOT NULL DEFAULT '',
    og_title        TEXT    NOT NULL DEFAULT '',
    og_image        TEXT    NOT NULL DEFAULT '',
    canonical_url   TEXT    NOT NULL DEFAULT '',
    robots          TEXT    NOT NULL DEFAULT 'index,follow',
    structured_data TEXT    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS doc_revisions (
    id         BIGINT      PRIMARY KEY,
    doc_id     BIGINT      NOT NULL REFERENCES docs (id) ON DELETE CASCADE,
    author_id  BIGINT      NOT NULL REFERENCES users (id),
    title      TEXT        NOT NULL DEFAULT '',
    content    TEXT        NOT NULL DEFAULT '',
    rev_note   TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_doc_revisions_doc ON doc_revisions (doc_id);

-- ============================================================
--  MOMENTS
-- ============================================================

-- moments: short-form social content — no title, no slug, no SEO
CREATE TABLE IF NOT EXISTS moments (
    id         BIGINT      PRIMARY KEY,
    author_id  BIGINT      NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    content    TEXT        NOT NULL,
    visibility INTEGER     NOT NULL DEFAULT 1,    -- 1=public 2=private 3=followers
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_moments_author  ON moments (author_id);
CREATE INDEX IF NOT EXISTS idx_moments_vis     ON moments (visibility, created_at);
CREATE INDEX IF NOT EXISTS idx_moments_deleted ON moments (deleted_at);

CREATE TABLE IF NOT EXISTS moment_stats (
    moment_id     BIGINT      PRIMARY KEY REFERENCES moments (id) ON DELETE CASCADE,
    view_count    BIGINT      NOT NULL DEFAULT 0,
    like_count    BIGINT      NOT NULL DEFAULT 0,
    comment_count BIGINT      NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- moment_media: ordered media attachments (images / videos) for a moment
CREATE TABLE IF NOT EXISTS moment_media (
    moment_id  BIGINT  NOT NULL REFERENCES moments (id) ON DELETE CASCADE,
    media_id   BIGINT  NOT NULL REFERENCES medias (id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (moment_id, media_id)
);

CREATE INDEX IF NOT EXISTS idx_moment_media_moment ON moment_media (moment_id);

-- ============================================================
--  TAXONOMY
-- ============================================================

CREATE TABLE IF NOT EXISTS terms (
    id          BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    slug        TEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS taxonomies (
    id          BIGINT      PRIMARY KEY,
    term_id     BIGINT      NOT NULL REFERENCES terms (id) ON DELETE CASCADE,
    taxonomy    TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    parent_id   BIGINT      REFERENCES taxonomies (id),
    post_count  INTEGER     NOT NULL DEFAULT 0,
    extra       TEXT        NOT NULL DEFAULT '{}',
    UNIQUE (term_id, taxonomy)
);

CREATE INDEX IF NOT EXISTS idx_taxonomies_term     ON taxonomies (term_id);
CREATE INDEX IF NOT EXISTS idx_taxonomies_taxonomy ON taxonomies (taxonomy);
CREATE INDEX IF NOT EXISTS idx_taxonomies_parent   ON taxonomies (parent_id);

CREATE TABLE IF NOT EXISTS object_taxonomies (
    object_id   BIGINT  NOT NULL,
    object_type TEXT    NOT NULL DEFAULT 'post',
    taxonomy_id BIGINT  NOT NULL REFERENCES taxonomies (id) ON DELETE CASCADE,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (object_id, object_type, taxonomy_id)
);

CREATE INDEX IF NOT EXISTS idx_obj_tax_object   ON object_taxonomies (object_id, object_type);
CREATE INDEX IF NOT EXISTS idx_obj_tax_taxonomy ON object_taxonomies (taxonomy_id);

-- ============================================================
--  COMMENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS comments (
    id           BIGINT      PRIMARY KEY,
    object_id    BIGINT      NOT NULL,
    object_type  TEXT        NOT NULL DEFAULT 'post',
    parent_id    BIGINT      REFERENCES comments (id) ON DELETE SET NULL,
    user_id      BIGINT      REFERENCES users (id),
    author_name  TEXT        NOT NULL DEFAULT '',
    author_email TEXT        NOT NULL DEFAULT '',
    content      TEXT        NOT NULL,
    status       INTEGER     NOT NULL DEFAULT 1,        -- 1=pending 2=approved 3=spam 4=trash
    ip           TEXT        NOT NULL DEFAULT '',
    user_agent   TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_comments_object        ON comments (object_id, object_type);
CREATE INDEX IF NOT EXISTS idx_comments_object_status ON comments (object_id, object_type, status);
CREATE INDEX IF NOT EXISTS idx_comments_parent        ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_user          ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_status        ON comments (status);
CREATE INDEX IF NOT EXISTS idx_comments_deleted       ON comments (deleted_at);

-- ============================================================
--  NOTIFICATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS notifications (
    id           BIGINT      PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         TEXT        NOT NULL,
    sub_type     TEXT        NOT NULL DEFAULT '',
    actor_id     BIGINT,
    actor_name   TEXT        NOT NULL DEFAULT '',
    actor_avatar TEXT        NOT NULL DEFAULT '',
    object_type  TEXT        NOT NULL DEFAULT '',
    object_id    BIGINT,
    object_title TEXT        NOT NULL DEFAULT '',
    object_link  TEXT        NOT NULL DEFAULT '',
    content      TEXT        NOT NULL DEFAULT '',
    is_read      INTEGER     NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id      ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read    ON notifications (user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications (user_id, created_at);

-- ============================================================
--  USER INTERACTIONS
-- ============================================================

-- user_likes / user_bookmarks: polymorphic — covers post | doc | moment
CREATE TABLE IF NOT EXISTS user_likes (
    user_id     BIGINT      NOT NULL,
    object_type TEXT        NOT NULL DEFAULT 'post',
    object_id   BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_user_likes_object ON user_likes (object_type, object_id);

CREATE TABLE IF NOT EXISTS user_bookmarks (
    user_id     BIGINT      NOT NULL,
    object_type TEXT        NOT NULL DEFAULT 'post',
    object_id   BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_user_bookmarks_object ON user_bookmarks (object_type, object_id);

-- Generic user action log: checkin / share / download / login / etc.
CREATE TABLE IF NOT EXISTS user_actions (
    id          BIGINT      PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    action      TEXT        NOT NULL,
    object_type TEXT        NOT NULL DEFAULT '',
    object_id   BIGINT      NOT NULL DEFAULT 0,
    extra       TEXT        NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_actions_user_action ON user_actions (user_id, action);
CREATE INDEX IF NOT EXISTS idx_user_actions_created     ON user_actions (user_id, created_at);

-- Personal access tokens (API keys)
CREATE TABLE IF NOT EXISTS user_tokens (
    id           BIGINT      PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    prefix       TEXT        NOT NULL,
    token_hash   TEXT        NOT NULL UNIQUE,
    expires_at   TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_tokens_user ON user_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_user_tokens_hash ON user_tokens (token_hash);

-- Follow graph
CREATE TABLE IF NOT EXISTS user_follows (
    follower_id  BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    following_id BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follows_follower  ON user_follows (follower_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_following ON user_follows (following_id);

-- Email/SMS verification codes
CREATE TABLE IF NOT EXISTS verification_codes (
    id         BIGINT      PRIMARY KEY,
    target     TEXT        NOT NULL,
    code       TEXT        NOT NULL,
    type       TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vcodes_target  ON verification_codes (target, type);
CREATE INDEX IF NOT EXISTS idx_vcodes_expires ON verification_codes (expires_at);

-- OAuth provider account links
CREATE TABLE IF NOT EXISTS user_oauth (
    id          BIGINT      PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider    TEXT        NOT NULL,
    provider_id TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

CREATE INDEX IF NOT EXISTS idx_user_oauth_user ON user_oauth(user_id);

-- Reports
CREATE TABLE IF NOT EXISTS reports (
    id          BIGINT      PRIMARY KEY,
    reporter_id BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    target_type TEXT        NOT NULL,
    target_id   BIGINT      NOT NULL,
    reason      TEXT        NOT NULL,
    detail      TEXT        NOT NULL DEFAULT '',
    status      TEXT        NOT NULL DEFAULT 'pending',
    notes       TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_reports_status   ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_reporter ON reports(reporter_id);
CREATE INDEX IF NOT EXISTS idx_reports_target   ON reports(target_type, target_id);

-- Private messages
-- id is application-assigned (random int64 via idgen) to avoid sequential enumeration.
CREATE TABLE IF NOT EXISTS conversations (
    id          BIGINT      PRIMARY KEY,
    user_a      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_b      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_msg    TEXT        NOT NULL DEFAULT '',
    last_msg_at TIMESTAMPTZ,
    unread_a    INTEGER     NOT NULL DEFAULT 0,
    unread_b    INTEGER     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_a, user_b)
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_a ON conversations(user_a);
CREATE INDEX IF NOT EXISTS idx_conversations_user_b ON conversations(user_b);

-- id is application-assigned (random int64 via idgen) to avoid sequential enumeration.
CREATE TABLE IF NOT EXISTS messages (
    id              BIGINT      PRIMARY KEY,
    conversation_id BIGINT      NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT        NOT NULL,
    is_read         BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender       ON messages(sender_id);

-- ============================================================
--  TRIGGERS — keep updated_at fresh automatically
-- ============================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_users_updated
    BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_posts_updated
    BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_post_stats_updated
    BEFORE UPDATE ON post_stats FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_post_metas_updated
    BEFORE UPDATE ON post_metas FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_options_updated
    BEFORE UPDATE ON options FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_doc_collections_updated
    BEFORE UPDATE ON doc_collections FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_docs_updated
    BEFORE UPDATE ON docs FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_doc_stats_updated
    BEFORE UPDATE ON doc_stats FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_moments_updated
    BEFORE UPDATE ON moments FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_moment_stats_updated
    BEFORE UPDATE ON moment_stats FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Plugins
CREATE TABLE IF NOT EXISTS plugins (
    id           TEXT        NOT NULL PRIMARY KEY,
    title        TEXT        NOT NULL DEFAULT '',
    description  TEXT        NOT NULL DEFAULT '',
    version      TEXT        NOT NULL DEFAULT '',
    author       TEXT        NOT NULL DEFAULT '',
    icon         TEXT        NOT NULL DEFAULT 'i-tabler-plug',
    repo_url     TEXT        NOT NULL DEFAULT '',
    script       TEXT        NOT NULL DEFAULT '',
    styles       TEXT        NOT NULL DEFAULT '',
    priority     INTEGER     NOT NULL DEFAULT 10,
    settings        TEXT        NOT NULL DEFAULT '{}',
    settings_schema TEXT        NOT NULL DEFAULT '[]',
    capabilities    TEXT        NOT NULL DEFAULT '{}', -- JSON capability declarations
    manifest        TEXT        NOT NULL DEFAULT '{}', -- JSON full manifest (pipelines, webhooks, settings schema)
    source          TEXT        NOT NULL DEFAULT 'external', -- 'builtin' (Go native) | 'external' (installed via zip/github)
    enabled         INTEGER     NOT NULL DEFAULT 1,
    installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- P-B10: add manifest column for existing databases (safe to repeat)
ALTER TABLE plugins ADD COLUMN IF NOT EXISTS manifest TEXT NOT NULL DEFAULT '{}';

-- Add source column for existing databases
ALTER TABLE plugins ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'external';

-- ============================================================
--  ANNOUNCEMENTS (broadcast notifications)
-- ============================================================

CREATE TABLE IF NOT EXISTS announcements (
    id          BIGINT      PRIMARY KEY,
    title       TEXT        NOT NULL,
    content     TEXT        NOT NULL DEFAULT '',
    type        TEXT        NOT NULL DEFAULT 'info',
    created_by  BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_announcements_created ON announcements (created_at);

-- One row per user: stores the timestamp of their last "mark all read" action.
-- Announcements created after last_read_at are considered unread for that user.
CREATE TABLE IF NOT EXISTS user_announcement_cursor (
    user_id      BIGINT      NOT NULL PRIMARY KEY,
    last_read_at TIMESTAMPTZ NOT NULL
);

-- ============================================================
--  DEFAULT DATA
--  Default admin user is seeded automatically by autoMigrate()
--  on first startup: username=admin  password=admin123  role=3
-- ============================================================

