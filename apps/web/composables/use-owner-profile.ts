import { MOCK_USER_PROFILES, type MockUserProfile } from '~/fixtures/users.mock'
import { useAuth } from '~/composables/use-auth'

export type OwnerProfile = MockUserProfile

export function useOwnerProfile() {
  const auth = useAuth()

  function resolve(ownerId: string | null | undefined): OwnerProfile | null {
    if (!ownerId) return null
    const mock = MOCK_USER_PROFILES[ownerId]
    if (mock) return mock
    if (auth.user.value?.id === ownerId) {
      return {
        id: auth.user.value.id,
        name: auth.user.value.name,
        email: auth.user.value.email,
        role: auth.user.value.is_super_admin ? 'Admin' : 'Member',
      }
    }
    return {
      id: ownerId,
      name: ownerId.slice(0, 8),
      email: '',
      role: '',
    }
  }

  return { resolve }
}
