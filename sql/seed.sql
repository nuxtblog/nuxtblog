-- ============================================================
--  Blog Seed Data — SQLite
--
--  Usage:
--    1. Start the server once (creates schema + admin user)
--    2. Stop the server
--    3. sqlite3 blog.db < sql/seed.sql
--    4. Restart the server
--
--  Admin login: admin / admin123
--
--  NOTE: Seed users below have placeholder password hashes (argon2id format
--  but won't authenticate). Reset their passwords via admin panel if needed.
--  All seed IDs are ≤ 20000 and won't conflict with app Snowflake IDs (18+ digits).
-- ============================================================

PRAGMA foreign_keys = OFF;

-- ============================================================
--  Missing tables (not in init.sql)
-- ============================================================

CREATE TABLE IF NOT EXISTS nav_menus (
    id          INTEGER  PRIMARY KEY,
    name        TEXT     NOT NULL,
    location    TEXT     NOT NULL DEFAULT '',
    description TEXT     NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    updated_at  DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE IF NOT EXISTS nav_menu_items (
    id          INTEGER  PRIMARY KEY,
    menu_id     INTEGER  NOT NULL,
    parent_id   INTEGER  NOT NULL DEFAULT 0,
    object_type TEXT     NOT NULL DEFAULT 'custom',
    object_id   INTEGER  NOT NULL DEFAULT 0,
    label       TEXT     NOT NULL DEFAULT '',
    url         TEXT     NOT NULL DEFAULT '',
    target      TEXT     NOT NULL DEFAULT '',
    css_classes TEXT     NOT NULL DEFAULT '',
    sort_order  INTEGER  NOT NULL DEFAULT 0
);

-- ============================================================
--  OPTIONS
-- ============================================================

INSERT OR REPLACE INTO options (key, value, autoload) VALUES
    ('site_name',        '"代码笔记"',                         1),
    ('site_description', '"分享编程技术、架构设计与开发实践"',  1),
    ('site_url',         '"http://localhost:3000"',             1),
    ('admin_email',      '"admin@example.com"',                 1),
    ('allow_registration', 'true',                              1),
    ('posts_per_page',   '10',                                  1);

-- ============================================================
--  USERS
--  Placeholder hash format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
--  These hashes will NOT match any real password — reset via admin panel.
-- ============================================================

INSERT OR IGNORE INTO users
    (id, username, email, password_hash, display_name, bio, role, status, email_verified, locale, last_login_at, created_at, updated_at)
VALUES
(1001, 'zhangwei', 'zhangwei@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDEx$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '张伟', 'Go 工程师，热衷后端架构与分布式系统设计。', 2, 1, 1, 'zh-CN',
 '2026-03-25 08:30:00', '2025-06-01 09:00:00', '2026-03-25 08:30:00'),

(1002, 'liming', 'liming@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDEy$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '李明', '前端工程师，Vue / React 深度用户，关注性能优化与工程化。', 2, 1, 1, 'zh-CN',
 '2026-03-24 20:15:00', '2025-07-15 14:00:00', '2026-03-24 20:15:00'),

(1003, 'wangfang', 'wangfang@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDEz$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '王芳', '全栈开发，喜欢折腾新技术。', 1, 1, 1, 'zh-CN',
 '2026-03-23 16:00:00', '2025-09-10 10:00:00', '2026-03-23 16:00:00'),

(1004, 'chenjie', 'chenjie@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDE0$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '陈杰', '在读研究生，方向分布式系统。', 1, 1, 0, 'zh-CN',
 '2026-03-22 11:00:00', '2025-10-20 09:00:00', '2026-03-22 11:00:00'),

(1005, 'liuyang', 'liuyang@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDE1$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '刘阳', 'DevOps 工程师，云原生爱好者。', 1, 1, 1, 'zh-CN',
 '2026-03-20 09:30:00', '2025-11-05 08:00:00', '2026-03-20 09:30:00'),

(1006, 'zhaoxue', 'zhaoxue@example.com',
 '$argon2id$v=19$m=65536,t=1,p=4$c2VlZHNhbHQxNjAxMDE2$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
 '赵雪', '设计师转前端，热爱用户体验。', 1, 1, 0, 'zh-CN',
 NULL, '2026-01-18 15:00:00', '2026-01-18 15:00:00');

-- ============================================================
--  USER PROFILES
-- ============================================================

INSERT OR IGNORE INTO user_profiles (user_id, website, github, twitter, location, social_links, notification_prefs, checkin_streak, last_checkin_date)
VALUES
(1001, 'https://zhangwei.dev', 'zhangwei-dev', 'zhangwei_go', '北京', '{}', '{}', 12, '2026-03-25'),
(1002, 'https://liming.io',    'liming-fe',    'liming_vue',  '上海', '{}', '{}', 5,  '2026-03-24'),
(1003, 'https://wangfang.me',  'wangfang-fs',  '',            '杭州', '{}', '{}', 3,  '2026-03-23'),
(1004, '',                     'chenjie-cs',   '',            '武汉', '{}', '{}', 0,  ''),
(1005, '',                     'liuyang-ops',  '',            '深圳', '{}', '{}', 7,  '2026-03-20'),
(1006, '',                     '',             '',            '成都', '{}', '{}', 0,  '');

-- ============================================================
--  MEDIAS  (storage_type=5: external URL, no file stored)
-- ============================================================

INSERT OR IGNORE INTO medias
    (id, uploader_id, storage_type, storage_key, cdn_url, filename, mime_type, file_size, width, height, alt_text, title, category, created_at)
VALUES
(2001, 1001, 5, 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=1200&q=80', 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=1200&q=80', 'go-concurrency.jpg',     'image/jpeg', 189000, 1200, 630, 'Go 并发编程封面',      'Go 并发编程',           'post', '2025-06-01 10:00:00'),
(2002, 1001, 5, 'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb?w=1200&q=80', 'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb?w=1200&q=80', 'docker-cover.jpg',       'image/jpeg', 210000, 1200, 630, 'Docker 容器封面',       'Docker Compose',        'post', '2025-06-15 10:00:00'),
(2003, 1001, 5, 'https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=1200&q=80', 'https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=1200&q=80', 'postgres-cover.jpg',     'image/jpeg', 195000, 1200, 630, 'PostgreSQL 数据库封面', 'PostgreSQL 索引优化',   'post', '2025-07-01 10:00:00'),
(2004, 1002, 5, 'https://images.unsplash.com/photo-1627398242454-45a1465196b3?w=1200&q=80', 'https://images.unsplash.com/photo-1627398242454-45a1465196b3?w=1200&q=80', 'typescript-cover.jpg',   'image/jpeg', 178000, 1200, 630, 'TypeScript 封面',       'TypeScript 高级类型',   'post', '2025-07-20 10:00:00'),
(2005, 1001, 5, 'https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=1200&q=80', 'https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=1200&q=80', 'k8s-cover.jpg',          'image/jpeg', 220000, 1200, 630, 'Kubernetes 封面',       'Kubernetes 实践',       'post', '2025-08-10 10:00:00'),
(2006, 1002, 5, 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=1200&q=80', 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=1200&q=80', 'frontend-cover.jpg',     'image/jpeg', 165000, 1200, 630, '前端开发封面',          '前端性能优化',          'post', '2025-09-05 10:00:00'),
(2007, 1001, 5, 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=1200&q=80', 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=1200&q=80', 'redis-cover.jpg',        'image/jpeg', 182000, 1200, 630, 'Redis 封面',            'Redis 数据结构',        'post', '2025-09-20 10:00:00'),
(2008, 1002, 5, 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=1200&q=80', 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=1200&q=80', 'microservice-cover.jpg', 'image/jpeg', 200000, 1200, 630, '微服务架构封面',        '微服务架构设计',        'post', '2025-10-01 10:00:00'),
(2009, 1001, 5, 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=1200&q=80', 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=1200&q=80', 'doc-goframe.jpg',        'image/jpeg', 175000, 1200, 630, 'GoFrame 文档封面',      'GoFrame 开发指南封面', 'doc',  '2025-08-01 10:00:00'),
(2010, 1002, 5, 'https://images.unsplash.com/photo-1593720213428-28a5b9e94613?w=1200&q=80', 'https://images.unsplash.com/photo-1593720213428-28a5b9e94613?w=1200&q=80', 'doc-frontend.jpg',       'image/jpeg', 168000, 1200, 630, '前端工程化封面',        '前端工程化手册封面',    'doc',  '2025-10-15 10:00:00');

-- ============================================================
--  TERMS
-- ============================================================

INSERT OR IGNORE INTO terms (id, name, slug, created_at) VALUES
(3001, 'Go',         'go',          '2025-06-01 08:00:00'),
(3002, 'Python',     'python',      '2025-06-01 08:00:00'),
(3003, 'JavaScript', 'javascript',  '2025-06-01 08:00:00'),
(3004, 'TypeScript', 'typescript',  '2025-06-01 08:00:00'),
(3005, 'Rust',       'rust',        '2025-06-01 08:00:00'),
(3006, 'Docker',     'docker',      '2025-06-01 08:00:00'),
(3007, 'Kubernetes', 'kubernetes',  '2025-06-01 08:00:00'),
(3008, '后端开发',   'backend',     '2025-06-01 08:00:00'),
(3009, '前端开发',   'frontend',    '2025-06-01 08:00:00'),
(3010, '数据库',     'database',    '2025-06-01 08:00:00'),
(3011, '云原生',     'cloud-native','2025-06-01 08:00:00'),
(3012, '教程',       'tutorial',    '2025-06-01 08:00:00'),
(3013, '最佳实践',   'best-practice','2025-06-01 08:00:00'),
(3014, '性能优化',   'performance', '2025-06-01 08:00:00'),
(3015, '开源',       'open-source', '2025-06-01 08:00:00');

-- ============================================================
--  TAXONOMIES  (categories + tags)
-- ============================================================

-- Categories
INSERT OR IGNORE INTO taxonomies (id, term_id, taxonomy, description, parent_id, post_count) VALUES
(4001, 3008, 'category', '服务端开发语言、框架与工具', NULL, 6),
(4002, 3001, 'category', 'Go 语言相关文章',             4001,  3),
(4003, 3002, 'category', 'Python 相关文章',             4001,  0),
(4004, 3010, 'category', '数据库设计、调优与运维',      4001,  2),
(4005, 3009, 'category', '前端框架、工具与工程化',      NULL,  4),
(4006, 3003, 'category', 'JavaScript 相关文章',         4005,  1),
(4007, 3004, 'category', 'TypeScript 相关文章',         4005,  1),
(4008, 3011, 'category', '容器、编排与云原生实践',      NULL,  3),
(4009, 3006, 'category', 'Docker 相关文章',             4008,  1),
(4010, 3007, 'category', 'Kubernetes 相关文章',         4008,  1);

-- Tags
INSERT OR IGNORE INTO taxonomies (id, term_id, taxonomy, description, parent_id, post_count) VALUES
(4011, 3012, 'tag', '入门教程系列', NULL, 5),
(4012, 3013, 'tag', '工程实践经验', NULL, 7),
(4013, 3014, 'tag', '性能调优技巧', NULL, 4),
(4014, 3015, 'tag', '开源项目推荐', NULL, 3),
(4015, 3005, 'tag', 'Rust 语言',   NULL, 1);

-- ============================================================
--  POSTS
-- ============================================================

INSERT OR IGNORE INTO posts
    (id, post_type, status, title, slug, content, excerpt, author_id, featured_img_id,
     comment_status, locale, published_at, created_at, updated_at)
VALUES

(5001, 1, 2, 'Go 并发编程：Goroutine 与 Channel 深度解析',
 'go-goroutine-channel-deep-dive',
'## 引言

Go 的并发模型是其最具吸引力的特性之一。与传统线程模型不同，Go 采用 **CSP（Communicating Sequential Processes）** 理论，通过 Goroutine 和 Channel 实现高效并发。

## Goroutine：轻量级线程

Goroutine 是 Go 运行时管理的协程，初始栈大小仅 2KB，可按需增长。启动一个 Goroutine 极其简单：

```go
go func() {
    fmt.Println("hello from goroutine")
}()
```

与系统线程（通常 1–8MB 栈）相比，单台机器可轻松跑百万级 Goroutine。

## Channel：通信即同步

Channel 是 Goroutine 之间的通信管道。Go 的哲学是：**不要通过共享内存来通信，而要通过通信来共享内存**。

```go
ch := make(chan int, 10) // 带缓冲 Channel
ch <- 42                 // 发送
v := <-ch                // 接收
```

### select 多路复用

```go
select {
case msg := <-ch1:
    fmt.Println("ch1:", msg)
case msg := <-ch2:
    fmt.Println("ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}
```

## 常见陷阱

- **Goroutine 泄漏**：Channel 永远没有发送方时，接收方 Goroutine 会永久阻塞。
- **数据竞争**：使用 `go build -race` 检测。
- **死锁**：所有 Goroutine 阻塞时 runtime 会 panic。

## 最佳实践

1. 用 `context.Context` 控制生命周期
2. 用 `sync.WaitGroup` 等待 Goroutine 完成
3. 避免在 Goroutine 中直接访问共享变量，优先通过 Channel 传递数据
4. Channel 应由发送方关闭，永远不要由接收方关闭

## 总结

理解 Goroutine 调度器（GMP 模型）和 Channel 的内部实现，是写出高性能 Go 代码的关键。',
 '深入解析 Go 并发模型的核心机制：Goroutine 的调度原理、Channel 的使用模式以及常见并发陷阱与最佳实践。',
 1001, 2001, 1, 'zh-CN', '2025-06-05 09:00:00', '2025-06-03 10:00:00', '2025-06-05 09:00:00'),

(5002, 1, 2, 'Docker Compose 最佳实践指南',
 'docker-compose-best-practices',
'## 为什么需要 Docker Compose

单容器部署简单，但现实项目往往由多个服务组成：Web、数据库、缓存、消息队列……手动管理这些容器既费时又容易出错。Docker Compose 通过一个 YAML 文件声明整个应用栈。

## 项目结构推荐

```
project/
├── docker-compose.yml        # 生产配置
├── docker-compose.override.yml # 本地开发覆盖
├── .env                      # 环境变量（不提交到 Git）
└── services/
    ├── api/Dockerfile
    └── worker/Dockerfile
```

## 核心配置示例

```yaml
version: "3.9"
services:
  api:
    build: ./services/api
    ports: ["9000:9000"]
    environment:
      DATABASE_URL: postgres://user:pass@db:5432/blog
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 5s
      retries: 5

volumes:
  pgdata:
```

## 关键最佳实践

1. **永远不要在镜像里硬编码密钥**，使用 `.env` 或 Docker Secrets
2. **合理设置 restart policy**，生产用 `unless-stopped`
3. **healthcheck + depends_on condition**，确保服务按正确顺序启动
4. **挂载数据卷而非绑定主机路径**，避免权限问题
5. **分离开发与生产配置**，用 `docker-compose.override.yml` 覆盖

## 总结

Docker Compose 是本地开发和中小规模生产部署的利器。掌握这些实践，可以让你的容器化工作流事半功倍。',
 '系统介绍 Docker Compose 的项目结构、核心配置与生产部署最佳实践，帮助你构建可靠的容器化应用栈。',
 1001, 2002, 1, 'zh-CN', '2025-06-20 10:00:00', '2025-06-18 14:00:00', '2025-06-20 10:00:00'),

(5003, 1, 2, 'PostgreSQL 索引优化实战',
 'postgresql-index-optimization',
'## 索引的代价

索引加快读取，但会减慢写入并占用存储。优化的第一步是**只在真正需要的列上建索引**。

## 常用索引类型

| 类型 | 适用场景 |
|------|---------|
| B-Tree | 默认，适合等值查询和范围查询 |
| Hash  | 仅等值查询，不支持范围 |
| GIN   | 数组、JSONB、全文检索 |
| BRIN  | 超大表中有自然排序的列（如时间戳） |

## 用 EXPLAIN ANALYZE 诊断

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM posts
WHERE author_id = 42 AND status = 2
ORDER BY published_at DESC
LIMIT 20;
```

关注 `Seq Scan`（全表扫描）— 这往往是需要索引的信号。

## 组合索引的列顺序

遵循 **"等值过滤在前，范围/排序在后"** 原则：

```sql
-- 查询: WHERE author_id = ? AND status = ? ORDER BY published_at DESC
CREATE INDEX idx_posts_author_status_pub
    ON posts (author_id, status, published_at DESC);
```

## 部分索引节省空间

只索引满足条件的行：

```sql
-- 只索引已发布的帖子
CREATE INDEX idx_posts_published
    ON posts (published_at DESC)
    WHERE status = 2;
```

## 定期维护

```sql
-- 重建膨胀的索引
REINDEX INDEX CONCURRENTLY idx_posts_author_status_pub;

-- 查看索引使用率
SELECT indexrelname, idx_scan, idx_tup_read
FROM pg_stat_user_indexes
WHERE schemaname = ''public''
ORDER BY idx_scan;
```

## 总结

索引优化是一个持续过程：建立→监控→调整。善用 `EXPLAIN ANALYZE` 和 `pg_stat_user_indexes`，避免盲目堆索引。',
 '深入讲解 PostgreSQL 索引类型、组合索引设计原则、部分索引与 EXPLAIN ANALYZE 实战，帮你精准提升查询性能。',
 1001, 2003, 1, 'zh-CN', '2025-07-08 09:00:00', '2025-07-06 15:00:00', '2025-07-08 09:00:00'),

(5004, 1, 2, 'TypeScript 高级类型技巧',
 'typescript-advanced-types',
'## 条件类型

条件类型让类型系统拥有了"if-else"能力：

```ts
type IsArray<T> = T extends any[] ? true : false;
type A = IsArray<number[]>; // true
type B = IsArray<string>;   // false
```

## infer 推断内部类型

```ts
type UnpackPromise<T> = T extends Promise<infer U> ? U : T;
type Result = UnpackPromise<Promise<string>>; // string
```

## 映射类型

```ts
type Nullable<T> = { [K in keyof T]: T[K] | null };
type Optional<T> = { [K in keyof T]?: T[K] };
```

## 模板字面量类型

```ts
type EventName<T extends string> = `on${Capitalize<T>}`;
type Click = EventName<"click">; // "onClick"
```

## 工具类型组合

```ts
// 将某些字段变为必填，其余保持可选
type WithRequired<T, K extends keyof T> =
    Omit<T, K> & Required<Pick<T, K>>;

interface User { id?: number; name?: string; email?: string }
type UserWithId = WithRequired<User, "id">;
// { id: number; name?: string; email?: string }
```

## 实战：安全的 API 响应类型

```ts
type ApiResult<T> =
    | { ok: true;  data: T }
    | { ok: false; error: string };

async function fetchUser(id: number): Promise<ApiResult<User>> {
    // ...
}
```

## 总结

TypeScript 的类型系统是图灵完备的——善用这些工具，可以让错误在编译阶段就暴露出来，而不是在运行时。',
 '系统梳理 TypeScript 条件类型、infer 推断、映射类型、模板字面量类型等高级技巧，配合实战案例加深理解。',
 1002, 2004, 1, 'zh-CN', '2025-07-25 08:00:00', '2025-07-23 10:00:00', '2025-07-25 08:00:00'),

(5005, 1, 2, 'Kubernetes 从入门到实践',
 'kubernetes-getting-started',
'## K8s 是什么

Kubernetes（K8s）是容器编排平台，负责容器的调度、伸缩、自愈与服务发现。核心对象：

- **Pod**：最小调度单元，包含一或多个容器
- **Deployment**：管理 Pod 副本集，支持滚动更新
- **Service**：为 Pod 组提供稳定的网络端点
- **ConfigMap / Secret**：配置与敏感数据管理

## 快速上手 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blog-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: blog-api
  template:
    metadata:
      labels:
        app: blog-api
    spec:
      containers:
      - name: api
        image: blog-api:v1.2.0
        ports:
        - containerPort: 9000
        resources:
          requests: { cpu: "100m", memory: "128Mi" }
          limits:   { cpu: "500m", memory: "512Mi" }
        readinessProbe:
          httpGet: { path: /healthz, port: 9000 }
```

## Service 暴露服务

```yaml
apiVersion: v1
kind: Service
metadata:
  name: blog-api-svc
spec:
  selector:
    app: blog-api
  ports:
  - port: 80
    targetPort: 9000
  type: ClusterIP
```

## 常用命令速查

```bash
kubectl get pods -n production
kubectl describe pod <pod-name>
kubectl logs -f <pod-name> --tail=100
kubectl rollout status deployment/blog-api
kubectl scale deployment blog-api --replicas=5
```

## 总结

K8s 的学习曲线陡峭，但一旦掌握，它能给你的服务带来强大的弹性与可靠性。从本地 minikube 开始练习是最好的路径。',
 '从核心概念到实战配置，带你系统掌握 Kubernetes 的 Deployment、Service、资源限制与常用运维命令。',
 1001, 2005, 1, 'zh-CN', '2025-08-15 09:00:00', '2025-08-12 14:00:00', '2025-08-15 09:00:00'),

(5006, 1, 2, 'Redis 五种数据结构应用场景',
 'redis-data-structures-use-cases',
'## String：万能选手

最常用的类型，除了缓存 KV，还能做：

- **计数器**：`INCR page:view:1001`
- **分布式锁**：`SET lock:resource uuid NX PX 30000`
- **限流**：结合 `INCR` + `EXPIRE`

## List：消息队列

```bash
LPUSH queue:email "send:user:1001"  # 生产者
BRPOP queue:email 0                  # 消费者（阻塞）
```

适合简单的异步任务队列，不需要可靠消息保证时的首选。

## Hash：对象存储

```bash
HSET user:1001 name "张伟" email "zhangwei@example.com"
HGET user:1001 name
HMGET user:1001 name email
```

适合存储对象，比序列化整个 JSON 更节省内存且支持局部更新。

## Set：去重与关系

```bash
SADD post:5001:tags "go" "concurrency"
SMEMBERS post:5001:tags
SINTERSTORE common:tags user:1:tags user:2:tags  # 共同标签
```

## Sorted Set：排行榜

```bash
ZADD hot:posts 9527 "5001"   # 文章 5001 得分 9527
ZREVRANGE hot:posts 0 9 WITHSCORES  # Top 10
ZINCRBY hot:posts 1 "5001"   # 加分
```

这是 Redis 最强大的结构之一，实时排行榜、延迟队列都离不开它。

## 总结

Redis 不只是缓存，五种数据结构各有专长。理解其内部编码（ziplist vs skiplist vs hashtable）有助于写出更高效的代码。',
 '详解 Redis String、List、Hash、Set、Sorted Set 五种核心数据结构的典型应用场景与命令示例。',
 1001, 2007, 1, 'zh-CN', '2025-09-03 08:00:00', '2025-09-01 15:00:00', '2025-09-03 08:00:00'),

(5007, 1, 2, '微服务架构设计原则与实践',
 'microservices-architecture-principles',
'## 什么时候不该用微服务

> "不要在一个团队里搞微服务。" — Sam Newman

微服务的收益（独立部署、技术多样性、故障隔离）在单团队小项目中往往被运维复杂度抵消。先把单体做好，再在自然边界处拆分。

## 服务划分原则

遵循 **DDD（领域驱动设计）** 的限界上下文（Bounded Context）：

- **用户服务**：认证、授权、profile
- **内容服务**：文章、文档、草稿
- **通知服务**：邮件、短信、站内信
- **媒体服务**：上传、存储、CDN

每个服务有自己的数据库——这是隔离的关键，也是最难的部分。

## 服务间通信

| 方式 | 适用场景 |
|------|---------|
| REST/gRPC | 同步请求-响应，延迟敏感 |
| 消息队列  | 异步事件，解耦，削峰 |
| GraphQL Federation | 统一查询层 |

## 可观测性三要素

没有这三样，微服务就是黑盒：

1. **日志**：结构化 JSON，统一 trace_id
2. **指标**：Prometheus + Grafana，监控 RED（Rate/Errors/Duration）
3. **链路追踪**：OpenTelemetry，可视化跨服务调用链

## 总结

微服务不是银弹。架构决策应驱动于业务边界、团队规模和运维能力，而不是技术潮流。',
 '从服务拆分原则、通信方式到可观测性体系，系统梳理微服务架构的设计思路与落地实践。',
 1001, 2008, 1, 'zh-CN', '2025-10-10 09:00:00', '2025-10-08 11:00:00', '2025-10-10 09:00:00'),

(5008, 1, 2, '前端性能优化实战手册',
 'frontend-performance-optimization',
'## 性能指标

Google 的 Core Web Vitals 是衡量用户体验的标准：

- **LCP**（最大内容绘制）< 2.5s
- **FID**（首次输入延迟）< 100ms
- **CLS**（累积布局偏移）< 0.1

## 资源优化

### 图片
- 使用 WebP / AVIF 格式（比 JPEG 小 30-50%）
- `<img loading="lazy">` 懒加载
- 用 `srcset` 提供不同分辨率

### JavaScript
- 代码分割：`import()` 动态导入
- Tree Shaking：确保 ES Module 格式
- 第三方库按需引入，别 `import _ from ''lodash''`

### CSS
- 关键 CSS 内联
- 避免 `@import`（串行加载）
- 使用 CSS containment

## 网络优化

```nginx
# 开启 Brotli/gzip
brotli on;
brotli_comp_level 6;
brotli_types text/html text/css application/javascript;

# 静态资源长缓存
location ~* \.(js|css|woff2)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 渲染优化

- **避免强制重排**：批量读取 DOM 属性，再批量写入
- **使用 CSS transform / opacity 做动画**（走合成线程，不触发重排）
- **Virtual List** 处理长列表（vue-virtual-scroller / react-window）

## 总结

性能优化是工程能力的综合体现。用 Lighthouse + Chrome DevTools 定位瓶颈，再针对性优化，比盲目套方案效率高 10 倍。',
 '系统整理前端性能优化的关键指标、资源压缩、网络配置与渲染优化技巧，附 Nginx 配置与工具推荐。',
 1002, 2006, 1, 'zh-CN', '2025-10-25 08:00:00', '2025-10-23 16:00:00', '2025-10-25 08:00:00'),

(5009, 1, 2, 'Git 工作流最佳实践',
 'git-workflow-best-practices',
'## 选哪种工作流

| 工作流 | 适合场景 |
|--------|---------|
| GitHub Flow | SaaS，持续部署 |
| Git Flow    | 有版本号的软件，定期发版 |
| Trunk-based | 大型团队，特性开关驱动 |

**绝大多数团队用 GitHub Flow 就够了**：main 永远可部署，功能通过 feature branch + PR 合入。

## 提交规范（Conventional Commits）

```
feat(auth): add OAuth2 GitHub login
fix(post): correct pagination off-by-one error
docs(api): update endpoint descriptions
refactor(user): extract email validation logic
```

格式：`<type>(<scope>): <subject>`，用于自动生成 changelog 和触发 semver 版本号。

## 保护主分支

```
main:
  ✓ Require PR review (1 approver)
  ✓ Require CI to pass
  ✓ Require linear history (禁止 merge commit)
  ✓ Dismiss stale reviews on push
```

## rebase vs merge

- `git merge`：保留完整历史，有 merge commit，适合长期分支
- `git rebase -i`：线性历史，清爽，适合 feature branch 合入前整理

**建议**：feature branch 内部随意提交，合 PR 前用 `rebase -i` 整理成逻辑清晰的几个 commit。

## 实用命令

```bash
git log --oneline --graph --all   # 可视化分支图
git stash push -m "wip: login form"
git bisect start                  # 二分查找引入 bug 的 commit
git reflog                        # 找回丢失的 commit
```

## 总结

好的 Git 纪律是团队协作的基础。提交信息是写给未来的自己的信——认真对待它。',
 '从工作流选择、提交规范到分支保护与 rebase 技巧，帮你建立清晰高效的 Git 协作规范。',
 1002, NULL, 1, 'zh-CN', '2025-11-08 09:00:00', '2025-11-06 14:00:00', '2025-11-08 09:00:00'),

(5010, 1, 2, 'Go 错误处理最佳实践',
 'go-error-handling-best-practices',
'## Go 错误处理哲学

Go 将错误视为普通值，而不是异常。这要求你显式处理每一个可能失败的调用——乍看繁琐，实则提高了代码可读性和可靠性。

## 定义有意义的错误类型

```go
// 业务错误（可向用户展示）
type AppError struct {
    Code    int
    Message string
    Err     error
}
func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

var ErrNotFound  = &AppError{Code: 404, Message: "资源不存在"}
var ErrForbidden = &AppError{Code: 403, Message: "无权访问"}
```

## 错误包装（Go 1.13+）

```go
if err := db.QueryRow(q).Scan(&id); err != nil {
    return fmt.Errorf("getUserById(%d): %w", id, err)
}
```

调用方用 `errors.Is` / `errors.As` 解包：

```go
var appErr *AppError
if errors.As(err, &appErr) {
    c.JSON(appErr.Code, gin.H{"message": appErr.Message})
    return
}
```

## sentinel error vs 类型断言

- `errors.Is`：检查错误链中是否包含特定值（sentinel error）
- `errors.As`：检查错误链中是否包含特定类型

## 不要忽略错误

```go
// ❌ 错误
_ = os.Remove(tmpFile)

// ✅ 正确——至少记录日志
if err := os.Remove(tmpFile); err != nil {
    log.Printf("warn: remove tmp file: %v", err)
}
```

## 总结

好的错误处理策略：在最底层返回详细错误，在最顶层（HTTP handler / main）决定如何呈现。中间层只负责包装上下文，不做策略决策。',
 '系统梳理 Go 错误处理的最佳实践：自定义错误类型、错误包装、errors.Is/As 的使用以及常见反模式。',
 1001, NULL, 1, 'zh-CN', '2025-11-22 09:00:00', '2025-11-20 11:00:00', '2025-11-22 09:00:00'),

(5011, 1, 2, 'Rust 所有权机制详解',
 'rust-ownership-explained',
'## 为什么需要所有权

内存管理有三种主流方式：
1. **手动管理**（C/C++）：快但不安全，易出现 use-after-free、内存泄漏
2. **垃圾回收**（Go/Java）：安全但有 GC 停顿
3. **所有权系统**（Rust）：编译期静态保证安全，零运行时开销

## 三条核心规则

1. 每个值有且仅有一个**所有者**（owner）
2. 所有者离开作用域时，值被**drop**（释放内存）
3. 同一时间只能有一个可变引用**或**任意多个不可变引用

## Move 语义

```rust
let s1 = String::from("hello");
let s2 = s1;  // s1 的所有权 move 到 s2
// println!("{}", s1);  // ❌ 编译错误：value moved
println!("{}", s2);     // ✅
```

## 借用（Borrow）

```rust
fn print_len(s: &String) {
    println!("{}", s.len());  // 只借用，不转移所有权
}

let s = String::from("hello");
print_len(&s);
println!("{}", s);  // ✅ s 依然有效
```

## 生命周期

```rust
// 告诉编译器：返回值的生命周期和输入参数一样长
fn longest<''a>(x: &''a str, y: &''a str) -> &''a str {
    if x.len() > y.len() { x } else { y }
}
```

## 总结

所有权是 Rust 的灵魂，也是学习曲线最陡峭的地方。但一旦理解了它，你会发现很多内存安全问题在编译阶段就被消灭了——这正是 Rust 的魔力所在。',
 '系统讲解 Rust 所有权的三条核心规则、Move 语义、借用规则与生命周期标注，帮你跨越最难的学习门槛。',
 1002, NULL, 1, 'zh-CN', '2025-12-10 09:00:00', '2025-12-08 15:00:00', '2025-12-10 09:00:00'),

(5012, 1, 2, 'Vue 3 Composition API 完全指南',
 'vue3-composition-api-complete-guide',
'## 为什么 Composition API

Options API 的问题：同一逻辑关注点的代码散落在 `data`、`methods`、`computed`、`watch` 中，难以复用。Composition API 把**逻辑**而非**选项**作为组织单元。

## setup() 基础

```vue
<script setup>
import { ref, computed, onMounted } from ''vue''

const count = ref(0)
const double = computed(() => count.value * 2)

function increment() { count.value++ }

onMounted(() => console.log(''mounted''))
</script>
```

`<script setup>` 是语法糖，编译后等价于 `setup()` 函数，代码更简洁。

## 响应式核心

- `ref()`：任意类型，通过 `.value` 访问
- `reactive()`：对象类型，直接访问属性
- `readonly()`：防止意外修改

```ts
const state = reactive({ loading: false, data: null })
state.loading = true  // 直接赋值
```

## 组合式函数（Composables）

```ts
// composables/usePost.ts
export function usePost(id: Ref<number>) {
    const post = ref(null)
    const loading = ref(true)

    watchEffect(async () => {
        loading.value = true
        post.value = await fetchPost(id.value)
        loading.value = false
    })

    return { post, loading }
}
```

**与 React Hooks 的关键区别**：Vue composable 内部不需要声明依赖数组，响应式系统自动追踪。

## 总结

Composition API 极大提升了 Vue 组件的可复用性和可维护性，是 Vue 3 项目的首选写法。结合 TypeScript 和 `<script setup>` 使用体验最佳。',
 '全面介绍 Vue 3 Composition API 的核心概念：setup 语法糖、ref/reactive、组合式函数封装与最佳实践。',
 1002, NULL, 1, 'zh-CN', '2026-01-05 09:00:00', '2026-01-03 11:00:00', '2026-01-05 09:00:00'),

(5013, 1, 1, 'GraphQL vs REST：如何做选择',
 'graphql-vs-rest-comparison',
'## 草稿

本文还在撰写中，将对比 GraphQL 与 REST 的适用场景……',
 'GraphQL 与 REST 的深度对比，帮助你在实际项目中做出合理的 API 设计选择。',
 1001, NULL, 1, 'zh-CN', NULL, '2026-02-01 14:00:00', '2026-02-10 09:00:00'),

(5014, 1, 1, 'WebAssembly 实战入门',
 'webassembly-practical-intro',
'## 草稿

WebAssembly 让非 JS 语言在浏览器中运行成为可能……',
 '用实战案例带你了解 WebAssembly 的核心概念、工具链与在浏览器中的应用场景。',
 1002, NULL, 1, 'zh-CN', NULL, '2026-02-20 10:00:00', '2026-02-20 10:00:00'),

(5015, 1, 1, '分布式事务解决方案对比',
 'distributed-transaction-solutions',
'## 草稿

分布式系统中的事务问题是架构师绕不开的难题……',
 '对比 2PC、Saga、TCC 等分布式事务方案的原理、优缺点与适用场景。',
 1001, NULL, 1, 'zh-CN', NULL, '2026-03-05 11:00:00', '2026-03-15 14:00:00');

-- ============================================================
--  POST STATS
-- ============================================================

INSERT OR IGNORE INTO post_stats (post_id, view_count, like_count, comment_count, share_count) VALUES
(5001, 12483, 186, 24, 43),
(5002,  8921, 132, 18, 31),
(5003,  7654,  98, 15, 22),
(5004,  9237, 143, 20, 35),
(5005, 11205, 168, 22, 40),
(5006,  6832,  87, 12, 18),
(5007,  5410,  72,  9, 14),
(5008,  8103, 124, 17, 29),
(5009,  6248,  91, 14, 21),
(5010,  7893, 116, 19, 27),
(5011,  4521,  63,  8, 12),
(5012,  9614, 155, 21, 38),
(5013,    47,   1,  0,  0),
(5014,    23,   0,  0,  0),
(5015,    31,   0,  0,  0);

-- ============================================================
--  POST SEO  (only for key published posts)
-- ============================================================

INSERT OR IGNORE INTO post_seo (post_id, meta_title, meta_desc, og_title, og_image, canonical_url, robots) VALUES
(5001, 'Go 并发编程：Goroutine 与 Channel 完全解析 | 代码笔记',
       '深入理解 Go 并发模型，掌握 Goroutine 和 Channel 的使用技巧，避开常见陷阱。',
       'Go 并发编程深度解析', 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=1200&q=80',
       'http://localhost:3000/posts/go-goroutine-channel-deep-dive', 'index,follow'),

(5004, 'TypeScript 高级类型技巧完全指南 | 代码笔记',
       '条件类型、infer 推断、映射类型、模板字面量——掌握这些高级技巧，让 TypeScript 为你全力工作。',
       'TypeScript 高级类型技巧', 'https://images.unsplash.com/photo-1627398242454-45a1465196b3?w=1200&q=80',
       'http://localhost:3000/posts/typescript-advanced-types', 'index,follow'),

(5005, 'Kubernetes 从入门到实践 | 代码笔记',
       '核心概念、Deployment 配置、Service 暴露与常用 kubectl 命令速查，帮你快速上手 K8s。',
       'Kubernetes 从入门到实践', 'https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=1200&q=80',
       'http://localhost:3000/posts/kubernetes-getting-started', 'index,follow');

-- ============================================================
--  POST METAS
-- ============================================================

INSERT OR IGNORE INTO post_metas (id, post_id, meta_key, meta_value, created_at, updated_at) VALUES
(90001, 5001, 'reading_time', '8',  '2025-06-05 09:00:00', '2025-06-05 09:00:00'),
(90002, 5001, 'is_featured',  '1',  '2025-06-05 09:00:00', '2025-06-05 09:00:00'),
(90003, 5002, 'reading_time', '6',  '2025-06-20 10:00:00', '2025-06-20 10:00:00'),
(90004, 5003, 'reading_time', '7',  '2025-07-08 09:00:00', '2025-07-08 09:00:00'),
(90005, 5003, 'is_featured',  '1',  '2025-07-08 09:00:00', '2025-07-08 09:00:00'),
(90006, 5004, 'reading_time', '7',  '2025-07-25 08:00:00', '2025-07-25 08:00:00'),
(90007, 5005, 'reading_time', '9',  '2025-08-15 09:00:00', '2025-08-15 09:00:00'),
(90008, 5005, 'is_featured',  '1',  '2025-08-15 09:00:00', '2025-08-15 09:00:00'),
(90009, 5006, 'reading_time', '6',  '2025-09-03 08:00:00', '2025-09-03 08:00:00'),
(90010, 5008, 'reading_time', '8',  '2025-10-25 08:00:00', '2025-10-25 08:00:00'),
(90011, 5008, 'is_featured',  '1',  '2025-10-25 08:00:00', '2025-10-25 08:00:00'),
(90012, 5012, 'reading_time', '8',  '2026-01-05 09:00:00', '2026-01-05 09:00:00'),
(90013, 5012, 'is_sticky',    '1',  '2026-01-05 09:00:00', '2026-01-05 09:00:00');

-- ============================================================
--  OBJECT TAXONOMIES  (post → category + tag)
-- ============================================================

INSERT OR IGNORE INTO object_taxonomies (object_id, object_type, taxonomy_id, sort_order) VALUES
(5001, 'post', 4001, 0), (5001, 'post', 4002, 0), (5001, 'post', 4011, 0), (5001, 'post', 4012, 0),
(5002, 'post', 4008, 0), (5002, 'post', 4009, 0), (5002, 'post', 4012, 0),
(5003, 'post', 4001, 0), (5003, 'post', 4004, 0), (5003, 'post', 4013, 0),
(5004, 'post', 4005, 0), (5004, 'post', 4007, 0), (5004, 'post', 4011, 0),
(5005, 'post', 4008, 0), (5005, 'post', 4010, 0), (5005, 'post', 4011, 0), (5005, 'post', 4012, 0),
(5006, 'post', 4001, 0), (5006, 'post', 4004, 0), (5006, 'post', 4011, 0),
(5007, 'post', 4001, 0), (5007, 'post', 4012, 0),
(5008, 'post', 4005, 0), (5008, 'post', 4013, 0),
(5009, 'post', 4012, 0),
(5010, 'post', 4001, 0), (5010, 'post', 4002, 0), (5010, 'post', 4012, 0),
(5011, 'post', 4011, 0), (5011, 'post', 4015, 0),
(5012, 'post', 4005, 0), (5012, 'post', 4006, 0), (5012, 'post', 4011, 0);

-- ============================================================
--  DOC COLLECTIONS
-- ============================================================

INSERT OR IGNORE INTO doc_collections
    (id, slug, title, description, cover_img_id, author_id, status, locale, sort_order, created_at, updated_at)
VALUES
(6001, 'goframe-v2-guide',
 'GoFrame v2 开发指南',
 '从零开始构建 GoFrame v2 项目，涵盖路由、DAO、服务层、中间件与部署配置，适合有 Go 基础的开发者。',
 2009, 1001, 2, 'zh-CN', 1, '2025-08-01 10:00:00', '2025-08-01 10:00:00'),

(6002, 'frontend-engineering',
 '前端工程化实践手册',
 '系统整理前端工程化的核心知识：构建工具、模块系统、代码规范、CI/CD 与性能监控，帮你搭建现代化前端工程体系。',
 2010, 1002, 2, 'zh-CN', 2, '2025-10-15 10:00:00', '2025-10-15 10:00:00');

-- ============================================================
--  DOCS
-- ============================================================

INSERT OR IGNORE INTO docs
    (id, collection_id, parent_id, sort_order, status, title, slug, content, excerpt, author_id,
     comment_status, locale, published_at, created_at, updated_at)
VALUES

(7001, 6001, NULL, 1, 2, '快速开始',              'goframe-quick-start',
 '## 安装 GoFrame CLI

```bash
go install github.com/gogf/gf/cmd/gf/v2@latest
```

## 创建项目

```bash
gf init my-blog
cd my-blog
gf run main.go
```

## 项目结构

GoFrame 遵循 DDD 分层架构，`api/` 定义接口契约，`internal/` 包含具体实现。详见后续章节。',
 '五分钟内安装 GoFrame CLI、创建项目并启动本地服务器。',
 1001, 1, 'zh-CN', '2025-08-05 09:00:00', '2025-08-03 10:00:00', '2025-08-05 09:00:00'),

(7002, 6001, NULL, 2, 2, '路由与控制器',          'goframe-routing-controller',
 '## 路由定义方式

GoFrame 使用结构体标签定义路由，不需要单独的路由注册文件：

```go
type PostGetOneReq struct {
    g.Meta `path:"/posts/{id}" method:"get" tags:"Post"`
    Id     int64 `v:"required|min:1"`
}
```

## 控制器实现

控制器层只负责 **解析请求 → 调用 Service → 返回响应**，不包含业务逻辑。',
 '介绍 GoFrame 的声明式路由、控制器结构与请求验证机制。',
 1001, 1, 'zh-CN', '2025-08-10 09:00:00', '2025-08-08 10:00:00', '2025-08-10 09:00:00'),

(7003, 6001, NULL, 3, 2, 'DAO 与数据库操作',      'goframe-dao-database',
 '## 生成 DAO

```bash
gf gen dao
```

## 基础查询

```go
var post entity.Posts
err := dao.Posts.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&post)
```

## 事务

```go
err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    _, err := dao.Posts.Ctx(ctx).Data(postData).Insert()
    if err != nil { return err }
    _, err = dao.PostStats.Ctx(ctx).Data(statsData).Insert()
    return err
})
```',
 '介绍 GoFrame DAO 的代码生成、常用查询模式与事务处理。',
 1001, 1, 'zh-CN', '2025-08-15 09:00:00', '2025-08-13 10:00:00', '2025-08-15 09:00:00'),

(7004, 6001, NULL, 4, 2, '中间件与认证',          'goframe-middleware-auth',
 '## 注册中间件

```go
s.Group("/api/v1", func(group *ghttp.RouterGroup) {
    group.Middleware(
        middleware.CORS,
        middleware.RequestID,
    )
    group.Group("/", func(g *ghttp.RouterGroup) {
        g.Middleware(middleware.JWTAuth)
        g.Bind(user.NewController())
    })
})
```

## JWT 验证中间件

解析 Bearer Token，将用户 ID 写入 Context，后续 Handler 通过 `service.Auth().GetLoginUserId(ctx)` 获取。',
 '介绍如何在 GoFrame 中注册中间件、实现 JWT 认证与权限控制。',
 1001, 1, 'zh-CN', '2025-08-20 09:00:00', '2025-08-18 14:00:00', '2025-08-20 09:00:00'),

(7005, 6001, NULL, 5, 2, '配置管理与部署',        'goframe-config-deploy',
 '## 配置文件

GoFrame 默认使用 `manifest/config/config.yaml`，支持多环境：

```yaml
# config.yaml
database:
  default:
    link: "sqlite::@file(./blog.db)"

# config.prod.yaml
database:
  default:
    link: "pgsql:user=blog password=xxx host=db port=5432 dbname=blog"
```

## Docker 部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server main.go

FROM alpine:3.19
COPY --from=builder /app/server /app/server
COPY --from=builder /app/manifest /app/manifest
CMD ["/app/server"]
```',
 '介绍 GoFrame 多环境配置管理与 Docker 容器化部署流程。',
 1001, 1, 'zh-CN', '2025-08-25 09:00:00', '2025-08-23 11:00:00', '2025-08-25 09:00:00'),

(7006, 6002, NULL, 1, 2, 'Vite 构建工具深度解析', 'vite-bundler-deep-dive',
 '## 为什么选择 Vite

传统打包工具（Webpack）在大型项目中启动慢（>30s），因为要先打包所有模块。Vite 利用浏览器原生 ES Module，只在请求时按需编译，冷启动通常 <500ms。

## 核心配置

```ts
// vite.config.ts
export default defineConfig({
  plugins: [vue()],
  resolve: { alias: { ''@'': resolve(__dirname, ''src'') } },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: [''vue'', ''vue-router'', ''pinia''],
          ui: [''@nuxt/ui''],
        }
      }
    }
  }
})
```',
 '深入解析 Vite 的原理、配置技巧与代码分割策略。',
 1002, 1, 'zh-CN', '2025-10-20 09:00:00', '2025-10-18 14:00:00', '2025-10-20 09:00:00'),

(7007, 6002, NULL, 2, 2, 'ESLint + Prettier 代码规范', 'eslint-prettier-setup',
 '## 基础配置

```js
// eslint.config.js
export default [
  { files: ["**/*.{ts,vue}"], rules: { "no-console": "warn" } }
]
```

```json
// .prettierrc
{ "semi": false, "singleQuote": true, "trailingComma": "all" }
```

## pre-commit 钩子

```bash
pnpm add -D husky lint-staged
pnpm husky init
echo "pnpm lint-staged" > .husky/pre-commit
```

```json
// package.json
"lint-staged": {
  "*.{ts,vue}": ["eslint --fix", "prettier --write"]
}
```',
 '从零配置 ESLint、Prettier 和 pre-commit 钩子，自动化代码规范检查。',
 1002, 1, 'zh-CN', '2025-10-28 09:00:00', '2025-10-26 15:00:00', '2025-10-28 09:00:00'),

(7008, 6002, NULL, 3, 2, 'Pinia 状态管理最佳实践', 'pinia-best-practices',
 '## 定义 Store

```ts
export const useUserStore = defineStore(''user'', () => {
    const user = ref<UserResponse | null>(null)
    const isLoggedIn = computed(() => !!user.value)

    async function login(credentials: LoginRequest) {
        const res = await useAuthApi().login(credentials)
        user.value = res.user
        localStorage.setItem(''token'', res.access_token)
    }

    return { user, isLoggedIn, login }
})
```

## 持久化

```ts
// pinia-plugin-persistedstate
defineStore(''user'', { /* ... */ }, { persist: { key: ''user'', storage: localStorage } })
```',
 '介绍 Pinia Setup Store 模式的定义、状态计算与持久化配置。',
 1002, 1, 'zh-CN', '2025-11-05 09:00:00', '2025-11-03 11:00:00', '2025-11-05 09:00:00'),

(7009, 6002, NULL, 4, 2, 'GitHub Actions CI/CD 实战', 'github-actions-cicd',
 '## 基础 CI 流水线

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v3
      - run: pnpm install --frozen-lockfile
      - run: pnpm lint
      - run: pnpm test
      - run: pnpm build
```

## 部署到服务器

```yaml
  deploy:
    needs: test
    if: github.ref == ''refs/heads/main''
    runs-on: ubuntu-latest
    steps:
      - uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: deploy
          key: ${{ secrets.SSH_KEY }}
          script: cd /app && git pull && docker compose up -d --build
```',
 '用 GitHub Actions 搭建前端项目的完整 CI/CD 流水线，从 lint 到自动部署。',
 1002, 1, 'zh-CN', '2025-11-15 09:00:00', '2025-11-13 14:00:00', '2025-11-15 09:00:00'),

(7010, 6002, NULL, 5, 2, '前端监控与性能采集',    'frontend-monitoring',
 '## Web Vitals 采集

```ts
import { onCLS, onFID, onLCP } from ''web-vitals''

function sendToAnalytics({ name, value, id }) {
    fetch(''/api/v1/metrics'', {
        method: ''POST'',
        body: JSON.stringify({ name, value, id, url: location.href })
    })
}

onCLS(sendToAnalytics)
onFID(sendToAnalytics)
onLCP(sendToAnalytics)
```

## 错误监控

```ts
window.addEventListener(''unhandledrejection'', (e) => {
    reportError({ type: ''promise'', message: e.reason?.message, stack: e.reason?.stack })
})
```

接入 Sentry 是更完整的方案：`Sentry.init({ dsn: ''...'', tracesSampleRate: 0.1 })`',
 '介绍前端性能指标采集（Web Vitals）与错误监控的实现方案。',
 1002, 1, 'zh-CN', '2025-11-25 09:00:00', '2025-11-23 15:00:00', '2025-11-25 09:00:00');

-- ============================================================
--  DOC STATS
-- ============================================================

INSERT OR IGNORE INTO doc_stats (doc_id, view_count, like_count, comment_count) VALUES
(7001, 3412, 45, 8),
(7002, 2876, 38, 6),
(7003, 2341, 31, 5),
(7004, 1987, 24, 3),
(7005, 1654, 19, 2),
(7006, 2103, 28, 4),
(7007, 1876, 22, 3),
(7008, 2456, 34, 5),
(7009, 1543, 18, 2),
(7010, 1298, 15, 1);

-- ============================================================
--  MOMENTS
-- ============================================================

INSERT OR IGNORE INTO moments (id, author_id, content, visibility, created_at, updated_at) VALUES
(8001, 1001, '今天写了一个很有趣的 Go channel 模式——用 fan-out/fan-in 处理并发任务，吞吐量提升了 3 倍。核心就是一个 merge 函数，把多个 channel 的输出汇聚成一个。代码发博客了，感兴趣的可以看看 🚀', 1, '2026-03-20 10:30:00', '2026-03-20 10:30:00'),
(8002, 1002, '刚看完 Chrome 124 的 release notes，有个很实用的新特性：CSS @starting-style，可以给元素首次出现时设置入场动画，不需要 JS 了！前端的 CSS 越来越强了。', 1, '2026-03-19 15:00:00', '2026-03-19 15:00:00'),
(8003, 1001, '用 pprof 分析了一个 Go 服务的内存泄漏，最后定位到一个 goroutine 在 channel 上永久阻塞，永远不会被 GC。所以说：goroutine 泄漏 = 内存泄漏。go build -race 和 pprof 是每个 Go 工程师必须掌握的工具。', 1, '2026-03-18 09:00:00', '2026-03-18 09:00:00'),
(8004, 1003, '把公司的老 Vue 2 项目迁移到 Vue 3 完成了，最大的感受：Composition API 写复杂逻辑真的比 Options API 清晰太多了，逻辑复用也方便了很多。TypeScript 支持也好了一个级别。', 1, '2026-03-17 14:00:00', '2026-03-17 14:00:00'),
(8005, 1001, '推荐一本书：《Designing Data-Intensive Applications》（数据密集型应用系统设计），业界公认的分布式系统入门经典。无论你做后端还是架构，这本书都值得反复读。', 1, '2026-03-15 11:00:00', '2026-03-15 11:00:00'),
(8006, 1005, '今天把团队的部署流程迁移到 ArgoCD + GitOps 了，之前用 Jenkins pipeline 每次部署要手动点按钮，现在 push 到 main 就自动同步，减少了很多人为失误。K8s + GitOps 真的是运维的未来。', 1, '2026-03-14 16:00:00', '2026-03-14 16:00:00'),
(8007, 1002, '尝试了一下 Bun 作为 Node.js 替代品，确实快很多（号称 4x faster），但生态兼容性还是差一些。生产环境我还是会继续用 Node，玩具项目可以试试。', 1, '2026-03-12 10:00:00', '2026-03-12 10:00:00'),
(8008, 1004, '刚发现 PostgreSQL 有个很好用的特性：`EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` 可以输出机器可读的执行计划，配合 explain.dalibo.com 可以可视化分析，排查慢查询效率翻倍。', 1, '2026-03-10 09:30:00', '2026-03-10 09:30:00');

-- ============================================================
--  MOMENT STATS
-- ============================================================

INSERT OR IGNORE INTO moment_stats (moment_id, view_count, like_count, comment_count) VALUES
(8001, 1243, 87, 12),
(8002,  892, 63,  8),
(8003, 1105, 79, 10),
(8004,  743, 52,  7),
(8005, 1876, 134, 15),
(8006,  654, 41,  5),
(8007,  832, 58,  9),
(8008,  921, 67, 11);

-- ============================================================
--  COMMENTS
-- ============================================================

INSERT OR IGNORE INTO comments
    (id, object_id, object_type, parent_id, user_id, author_name, author_email, content, status, ip, created_at)
VALUES
(9001, 5001, 'post', NULL, 1003, '王芳',   'wangfang@example.com', '写得非常清晰！之前一直搞不明白 GMP 模型，看了这篇终于理解了 P 的作用。请问有计划写 sync 包系列吗？', 2, '127.0.0.1', '2025-06-06 10:00:00'),
(9002, 5001, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  'select + time.After 的 timeout 模式我之前一直用 context.WithTimeout，原来两种方式都可以。赞一个！', 2, '127.0.0.1', '2025-06-07 14:00:00'),
(9003, 5001, 'post', 9001, 1001, '张伟',   'zhangwei@example.com', '谢谢！sync 包系列已经在计划中了，sync.Map 和 errgroup 都会写到。', 2, '127.0.0.1', '2025-06-07 15:00:00'),
(9004, 5001, 'post', NULL, NULL, '游客读者', 'visitor@test.com',   'GMP 模型那张图能不能补充一下？文字描述有点抽象。', 1, '127.0.0.1', '2025-06-08 09:00:00'),

(9005, 5002, 'post', NULL, 1003, '王芳',   'wangfang@example.com', 'healthcheck + depends_on condition 这个组合之前踩过坑，等服务启动但是没有 ready，加上这个就解决了。很实用！', 2, '127.0.0.1', '2025-06-22 11:00:00'),
(9006, 5002, 'post', NULL, 1005, '刘阳',   'liuyang@example.com',  '生产环境建议再加 `logging` 配置，统一收集到 ELK 方便排查问题。', 2, '127.0.0.1', '2025-06-23 16:00:00'),
(9007, 5002, 'post', 9006, 1001, '张伟',   'zhangwei@example.com', '好建议！下一篇会专门写 Docker 日志收集方案。', 2, '127.0.0.1', '2025-06-23 17:00:00'),

(9008, 5003, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  '请问 BRIN 索引在时序数据场景下真的比 B-Tree 强很多吗？我们有一张日志表 2 亿行，目前用 B-Tree 查询已经很慢了。', 2, '127.0.0.1', '2025-07-10 10:00:00'),
(9009, 5003, 'post', 9008, 1001, '张伟',   'zhangwei@example.com', '如果是按时间范围查询且数据是按时间顺序插入的，BRIN 确实可以显著减少索引大小。但如果有大量随机读，还是 B-Tree 更可靠。建议先用 EXPLAIN ANALYZE 对比两种方案。', 2, '127.0.0.1', '2025-07-10 14:00:00'),

(9010, 5004, 'post', NULL, 1003, '王芳',   'wangfang@example.com', '条件类型 + infer 那段代码我抄下来了，解决了我项目里的一个类型问题，谢谢！', 2, '127.0.0.1', '2025-07-27 10:00:00'),
(9011, 5004, 'post', NULL, 1005, '刘阳',   'liuyang@example.com',  '`WithRequired` 那个工具类型太实用了，第一次见到这种写法，收藏了。', 2, '127.0.0.1', '2025-07-28 15:00:00'),
(9012, 5004, 'post', NULL, NULL, '匿名用户', 'anon@test.com',      '这篇是 TS 入门教程里讲高级类型最清楚的一篇，建议加到收藏。', 2, '127.0.0.1', '2025-07-29 09:00:00'),

(9013, 5005, 'post', NULL, 1005, '刘阳',   'liuyang@example.com',  'readinessProbe 这块还可以加 livenessProbe 和 startupProbe，尤其是启动慢的服务（比如 JVM）用 startupProbe 很关键，否则会被误杀。', 2, '127.0.0.1', '2025-08-17 10:00:00'),
(9014, 5005, 'post', 9013, 1001, '张伟',   'zhangwei@example.com', '好补充！三种 Probe 的区别确实值得单独写一篇，加入计划了。', 2, '127.0.0.1', '2025-08-17 11:30:00'),
(9015, 5005, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  '资源 limits 那里，CPU 压缩（throttling）是 K8s 一个很常见的性能问题，建议补充一下设置 limits 的注意事项。', 2, '127.0.0.1', '2025-08-18 09:00:00'),

(9016, 5008, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  '强制重排那段代码示例很经典，面试常考。虽然现在 React / Vue 框架帮我们批量处理了很多，但原生 JS 操作 DOM 时还是要注意。', 2, '127.0.0.1', '2025-10-27 10:00:00'),
(9017, 5008, 'post', NULL, 1003, '王芳',   'wangfang@example.com', '之前没有系统了解过 Brotli，看了这篇赶紧在公司 Nginx 上开了，页面资源减小了 25%，体感有差异。', 2, '127.0.0.1', '2025-10-28 16:00:00'),

(9018, 5010, 'post', NULL, 1003, '王芳',   'wangfang@example.com', '`errors.As` 这个我一直搞不清楚和 `errors.Is` 的区别，这篇终于讲清楚了！', 2, '127.0.0.1', '2025-11-24 11:00:00'),
(9019, 5010, 'post', NULL, 1006, '赵雪',   'zhaoxue@example.com',  '从 Python 转 Go 的时候，最不适应的就是到处 `if err != nil`，但用久了确实更好维护，错误路径很清晰。', 2, '127.0.0.1', '2025-11-25 14:00:00'),

(9020, 5012, 'post', NULL, 1006, '赵雪',   'zhaoxue@example.com',  '作为从 React 转过来的，Vue 3 的响应式系统确实比 React Hooks 直觉多了，不用 useMemo/useCallback 到处优化。', 2, '127.0.0.1', '2026-01-07 10:00:00'),
(9021, 5012, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  'composable 那个示例很经典，`watchEffect` 自动追踪依赖这个细节很重要，之前踩过 watch 忘加依赖的坑。', 2, '127.0.0.1', '2026-01-08 15:00:00'),
(9022, 5012, 'post', 9020, 1002, '李明',   'liming@example.com',   '对的，Vue 的依赖追踪是编译时+运行时双重机制，React 是纯运行时追踪，所以 Vue 不需要手动声明依赖数组。', 2, '127.0.0.1', '2026-01-08 16:00:00'),

(9023, 8001, 'post', NULL, 1003, '王芳',   'wangfang@example.com', 'fan-out/fan-in 这个模式能发一下代码链接吗？想学习一下具体实现', 2, '127.0.0.1', '2026-03-20 11:00:00'),
(9024, 8001, 'post', NULL, 1005, '刘阳',   'liuyang@example.com',  '在 K8s 里跑这个的话，goroutine 数量要注意，建议和 CPU 核心数挂钩', 2, '127.0.0.1', '2026-03-20 14:00:00'),
(9025, 8005, 'post', NULL, 1004, '陈杰',   'chenjie@example.com',  '这本书看了两遍了，每次都有新收获，强烈推荐！', 2, '127.0.0.1', '2026-03-15 13:00:00');

-- ============================================================
--  USER LIKES
-- ============================================================

INSERT OR IGNORE INTO user_likes (user_id, object_type, object_id, created_at) VALUES
(1003, 'post', 5001, '2025-06-06 10:00:00'),
(1003, 'post', 5002, '2025-06-22 11:00:00'),
(1003, 'post', 5004, '2025-07-27 10:00:00'),
(1003, 'post', 5008, '2025-10-27 10:00:00'),
(1004, 'post', 5001, '2025-06-07 14:00:00'),
(1004, 'post', 5003, '2025-07-10 10:00:00'),
(1004, 'post', 5005, '2025-08-18 09:00:00'),
(1004, 'post', 5010, '2025-11-24 11:00:00'),
(1005, 'post', 5002, '2025-06-23 16:00:00'),
(1005, 'post', 5005, '2025-08-17 10:00:00'),
(1005, 'post', 5007, '2025-10-12 15:00:00'),
(1006, 'post', 5010, '2025-11-25 14:00:00'),
(1006, 'post', 5012, '2026-01-07 10:00:00'),
(1002, 'post', 5001, '2025-06-10 09:00:00'),
(1002, 'post', 5003, '2025-07-09 11:00:00'),
(1003, 'post', 8001, '2026-03-20 11:00:00'),
(1004, 'post', 8005, '2026-03-15 12:00:00'),
(1005, 'post', 8001, '2026-03-20 12:00:00');

-- ============================================================
--  USER BOOKMARKS
-- ============================================================

INSERT OR IGNORE INTO user_bookmarks (user_id, object_type, object_id, created_at) VALUES
(1003, 'post', 5001, '2025-06-06 10:30:00'),
(1003, 'post', 5004, '2025-07-27 10:30:00'),
(1004, 'post', 5003, '2025-07-10 10:30:00'),
(1004, 'post', 5005, '2025-08-18 09:30:00'),
(1005, 'post', 5007, '2025-10-12 15:30:00'),
(1006, 'post', 5012, '2026-01-07 10:30:00'),
(1006, 'post', 5010, '2025-11-25 14:30:00'),
(1002, 'post', 5001, '2025-06-10 09:30:00');

-- ============================================================
--  USER FOLLOWS
-- ============================================================

INSERT OR IGNORE INTO user_follows (follower_id, following_id, created_at) VALUES
(1003, 1001, '2025-09-01 10:00:00'),  -- 王芳 → 张伟
(1004, 1001, '2025-09-05 11:00:00'),  -- 陈杰 → 张伟
(1005, 1001, '2025-09-10 09:00:00'),  -- 刘阳 → 张伟
(1006, 1001, '2026-01-20 14:00:00'),  -- 赵雪 → 张伟
(1003, 1002, '2025-11-01 10:00:00'),  -- 王芳 → 李明
(1004, 1002, '2025-11-05 11:00:00'),  -- 陈杰 → 李明
(1006, 1002, '2026-01-25 15:00:00'),  -- 赵雪 → 李明
(1001, 1002, '2025-07-01 09:00:00'),  -- 张伟 → 李明（互关）
(1002, 1001, '2025-07-02 10:00:00'),  -- 李明 → 张伟（互关）
(1005, 1004, '2025-12-01 09:00:00');  -- 刘阳 → 陈杰

-- ============================================================
-- ============================================================

INSERT OR IGNORE INTO notifications
    (id, user_id, type, sub_type, actor_id, actor_name, actor_avatar,
     object_type, object_id, object_title, object_link, content, is_read, created_at)
VALUES
(10001, 1001, 'comment', '', 1003, '王芳', '', 'post', 5001, 'Go 并发编程：Goroutine 与 Channel 深度解析',
 '/posts/go-goroutine-channel-deep-dive', '王芳 评论了你的文章：写得非常清晰！', 1, '2025-06-06 10:00:00'),

(10002, 1001, 'like', '', 1004, '陈杰', '', 'post', 5001, 'Go 并发编程：Goroutine 与 Channel 深度解析',
 '/posts/go-goroutine-channel-deep-dive', '陈杰 赞了你的文章', 1, '2025-06-07 14:00:00'),

(10003, 1001, 'follow', '', 1003, '王芳', '', '', 0, '', '', '王芳 关注了你', 1, '2025-09-01 10:00:00'),

(10004, 1001, 'follow', '', 1005, '刘阳', '', '', 0, '', '', '刘阳 关注了你', 1, '2025-09-10 09:00:00'),

(10005, 1001, 'comment', '', 1005, '刘阳', '', 'post', 5005, 'Kubernetes 从入门到实践',
 '/posts/kubernetes-getting-started', '刘阳 评论了你的文章：readinessProbe 这块还可以加 livenessProbe…', 0, '2025-08-17 10:00:00'),

(10006, 1002, 'comment', '', 1003, '王芳', '', 'post', 5008, '前端性能优化实战手册',
 '/posts/frontend-performance-optimization', '王芳 评论了你的文章：之前没有系统了解过 Brotli…', 0, '2025-10-28 16:00:00'),

(10007, 1002, 'like', '', 1006, '赵雪', '', 'post', 5012, 'Vue 3 Composition API 完全指南',
 '/posts/vue3-composition-api-complete-guide', '赵雪 赞了你的文章', 0, '2026-01-07 10:00:00'),

(10008, 1003, 'reply', '', 1001, '张伟', '', 'post', 5001, 'Go 并发编程：Goroutine 与 Channel 深度解析',
 '/posts/go-goroutine-channel-deep-dive', '张伟 回复了你的评论：谢谢！sync 包系列已经在计划中了…', 0, '2025-06-07 15:00:00'),

(10009, 1003, 'system', 'announcement', 0, '', '', '', 0, '', '',
 '平台新功能上线：支持文档模块，欢迎体验！', 0, '2025-08-01 09:00:00'),

(10010, 1004, 'reply', '', 1001, '张伟', '', 'post', 5003, 'PostgreSQL 索引优化实战',
 '/posts/postgresql-index-optimization', '张伟 回复了你的评论：如果是按时间范围查询…', 0, '2025-07-10 14:00:00');

-- ============================================================
--  REPORTS
-- ============================================================

INSERT OR IGNORE INTO reports
    (id, reporter_id, target_type, target_id, reason, detail, status, notes, created_at, resolved_at)
VALUES
(11001, 1003, 'comment', 9004, '垃圾信息',    '该评论疑似广告，内容与文章无关。', 'resolved', '已确认为无关评论，已删除。', '2025-06-08 11:00:00', '2025-06-08 16:00:00'),
(11002, 1004, 'post',    5013, '内容重复',    '与已有文章《REST API 设计规范》内容高度重叠。', 'pending', '', '2026-02-15 10:00:00', NULL),
(11003, 1005, 'comment', 9004, '不友善内容', '评论语气不当，对作者不尊重。', 'dismissed', '评论内容属于正常反馈，不符合处理标准。', '2025-06-09 09:00:00', '2025-06-09 14:00:00');

-- ============================================================
--  ANNOUNCEMENTS
-- ============================================================

INSERT OR IGNORE INTO announcements (id, title, content, type, created_by, created_at) VALUES
(12001, '博客平台正式上线！',
 '经过数月开发，**代码笔记**博客平台今日正式上线 🎉\n\n主要功能：\n- 文章发布与管理\n- 文档知识库\n- 社交动态（Moments）\n- 通知与私信系统\n\n感谢所有参与内测的用户！如有问题请通过「举报」功能反馈。',
 'success', 1, '2025-06-01 08:00:00'),

(12002, '新功能上线：文档模块',
 '**文档模块**（Docs）现已上线，支持：\n- 创建文档集合（类似 GitBook）\n- 文章嵌套章节结构\n- 修订历史与版本回滚\n\n欢迎体验并反馈！',
 'info', 1, '2025-08-01 09:00:00'),

(12003, '系统维护通知',
 '**计划维护时间**：2026-04-01 02:00–04:00 UTC+8\n\n维护期间服务暂停访问，数据不受影响。请提前保存草稿。感谢理解与支持！',
 'warning', 1, '2026-03-25 10:00:00');

-- ============================================================
--  NAV MENUS
-- ============================================================

INSERT OR IGNORE INTO nav_menus (id, name, location, description, created_at, updated_at) VALUES
(13001, '主导航', 'header', '顶部主导航菜单', '2025-06-01 08:00:00', '2025-06-01 08:00:00'),
(13002, '底部链接', 'footer', '底部页脚导航', '2025-06-01 08:00:00', '2025-06-01 08:00:00');

INSERT OR IGNORE INTO nav_menu_items (id, menu_id, parent_id, object_type, object_id, label, url, target, css_classes, sort_order) VALUES
(14001, 13001, 0, 'custom', 0, '首页',       '/',                   '', '', 1),
(14002, 13001, 0, 'custom', 0, '文章',       '/posts',              '', '', 2),
(14003, 13001, 0, 'custom', 0, '文档',       '/docs',               '', '', 3),
(14004, 13001, 0, 'custom', 0, '动态',       '/moments',            '', '', 4),
(14005, 13001, 0, 'custom', 0, '关于',       '/about',              '', '', 5),
(14006, 13001, 14002, 'category', 4001, '后端开发', '/categories/backend',  '', '', 1),
(14007, 13001, 14002, 'category', 4005, '前端开发', '/categories/frontend', '', '', 2),
(14008, 13001, 14002, 'category', 4008, '云原生',   '/categories/cloud-native', '', '', 3),
(14009, 13002, 0, 'custom', 0, '关于我们',   '/about',              '', '', 1),
(14010, 13002, 0, 'custom', 0, 'GitHub',     'https://github.com',  '_blank', '', 2),
(14011, 13002, 0, 'custom', 0, 'RSS',        '/feed.xml',           '', '', 3),
(14012, 13002, 0, 'custom', 0, '隐私政策',   '/privacy',            '', '', 4);

-- ============================================================
--  PRIVATE MESSAGES  (one conversation for demo)
-- ============================================================

INSERT OR IGNORE INTO conversations (id, user_a, user_b, last_msg, last_msg_at, unread_a, unread_b, created_at) VALUES
(15001, 1001, 1002, '好的，我看看！', '2026-03-15 16:30:00', 0, 1, '2026-03-15 15:00:00');

INSERT OR IGNORE INTO messages (id, conversation_id, sender_id, content, is_read, created_at) VALUES
(16001, 15001, 1002, '张伟，你好！你那篇 Go 并发的文章写得很棒，我们团队内部分享了。', 1, '2026-03-15 15:00:00'),
(16002, 15001, 1001, '谢谢！有什么问题随时问我。', 1, '2026-03-15 15:05:00'),
(16003, 15001, 1002, '有个问题想请教：在高并发场景下，你们用 channel 还是 sync.Mutex 更多？', 1, '2026-03-15 15:10:00'),
(16004, 15001, 1001, '一般来说：需要通信 → channel，保护共享状态 → Mutex。不要混用。', 1, '2026-03-15 15:20:00'),
(16005, 15001, 1002, '好的，我看看！', 0, '2026-03-15 16:30:00');

-- ============================================================
PRAGMA foreign_keys = ON;

-- ============================================================
--  Done! Summary:
--    6 seed users (IDs 1001–1006, passwords NOT set — reset via admin panel)
--    10 medias (external URLs)
--    15 terms + 15 taxonomies (10 categories + 5 tags)
--    15 posts (12 published + 3 draft) with stats, SEO, metas
--    2 doc collections + 10 docs with stats
--    8 moments with stats
--    25 comments
--    18 likes + 8 bookmarks
--    10 follows
--    10 notifications
--    3 reports + 3 announcements
--    2 nav menus + 12 nav menu items
--    1 conversation + 5 messages
-- ============================================================
