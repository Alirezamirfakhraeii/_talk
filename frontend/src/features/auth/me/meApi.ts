import { get } from '../../../shared/api/client'

export type MeResponseData = {
    id: number
    name: string
    username: string
    email: string
}

export function getCurrentUser() {
    return get<MeResponseData>(
        '/api/v1/me',
    )
}