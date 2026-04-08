/**
 * Enhance toast notifications:
 * - Error/warning toasts stay 15s (instead of default 8s)
 * - Error toasts get a "copy" action button for easy copying of error messages
 */
export default defineNuxtPlugin(() => {
  const toast = useToast()
  const originalAdd = toast.add.bind(toast)

  toast.add = (notification) => {
    const isError = notification.color === 'error'
    const isWarning = notification.color === 'warning'

    // Extend duration for error/warning toasts
    if ((isError || isWarning) && notification.duration === undefined) {
      notification.duration = 15000
    }

    // Add copy action to error toasts
    if (isError && !notification.actions?.length) {
      const text = [notification.title, notification.description].filter(Boolean).join(': ')
      if (text) {
        notification.actions = [{
          label: '复制',
          variant: 'outline' as const,
          onClick: () => {
            navigator.clipboard.writeText(text).catch(() => {})
          },
        }]
      }
    }

    return originalAdd(notification)
  }
})
