export interface FieldOption {
  label: string
  value: string
}

export interface FieldDef {
  key: string
  label: string
  type: 'text' | 'password' | 'switch' | 'select'
  required?: boolean
  placeholder?: string
  options?: FieldOption[]
}

export interface PaymentProvider {
  slug: string
  label: string
  icon: string
  enabled: boolean
  config: Record<string, any>
  fields: FieldDef[]
}

export function usePaymentApi() {
  const { apiFetch } = useApiFetch()

  const listProviders = () =>
    apiFetch<{ items: PaymentProvider[] }>('/admin/payment/providers')

  const getProviderConfig = (slug: string) =>
    apiFetch<PaymentProvider>(`/admin/payment/providers/${slug}/config`)

  const setProviderConfig = (slug: string, config: Record<string, any>) =>
    apiFetch<PaymentProvider>(`/admin/payment/providers/${slug}/config`, {
      method: 'PUT',
      body: { config },
    })

  return { listProviders, getProviderConfig, setProviderConfig }
}
