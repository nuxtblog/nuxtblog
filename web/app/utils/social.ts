export function getSocialIcon(label: string, url: string): string {
  const s = (label + ' ' + url).toLowerCase()
  if (s.includes('github'))                           return 'i-tabler-brand-github'
  if (s.includes('twitter') || s.includes('x.com'))  return 'i-tabler-brand-x'
  if (s.includes('weibo'))                            return 'i-tabler-brand-weibo'
  if (s.includes('bilibili'))                         return 'i-tabler-brand-bilibili'
  if (s.includes('instagram'))                        return 'i-tabler-brand-instagram'
  if (s.includes('facebook'))                         return 'i-tabler-brand-facebook'
  if (s.includes('youtube'))                          return 'i-tabler-brand-youtube'
  if (s.includes('linkedin'))                         return 'i-tabler-brand-linkedin'
  if (s.includes('telegram'))                         return 'i-tabler-brand-telegram'
  if (s.includes('discord'))                          return 'i-tabler-brand-discord'
  if (s.includes('wechat') || s.includes('微信'))     return 'i-tabler-brand-wechat'
  if (s.includes('zhihu')  || s.includes('知乎'))     return 'i-tabler-brand-zhihu'
  if (s.includes('juejin') || s.includes('掘金'))     return 'i-tabler-brand-juejin'
  if (s.includes('mail')   || s.includes('邮箱'))     return 'i-tabler-mail'
  if (s.includes('rss')    || s.includes('/feed'))    return 'i-tabler-rss'
  return 'i-tabler-link'
}

export const SOCIAL_PLATFORMS = [
  { label: 'GitHub', value: 'GitHub' },
  { label: 'Twitter/X', value: 'Twitter' },
  { label: '微博', value: '微博' },
  { label: 'Bilibili', value: 'Bilibili' },
  { label: '知乎', value: '知乎' },
  { label: '掘金', value: '掘金' },
  { label: 'Instagram', value: 'Instagram' },
  { label: 'YouTube', value: 'YouTube' },
  { label: 'LinkedIn', value: 'LinkedIn' },
  { label: 'Facebook', value: 'Facebook' },
  { label: 'Telegram', value: 'Telegram' },
  { label: 'Discord', value: 'Discord' },
  { label: '微信', value: '微信' },
  { label: '邮箱', value: 'mail' },
  { label: 'RSS', value: 'RSS' },
  { label: '自定义', value: '__custom__' },
]
