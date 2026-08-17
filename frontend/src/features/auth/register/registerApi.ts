import { post } from '../../../shared/api/client'

export type RegisterRequest = {
    name: string
    username: string
    email: string
    password: string
}

export type RegisterResponseData = {
    id: number
    name: string
    username: string
    email: string
    created_at: string
}

export function register(
    input: RegisterRequest,
) {
    return post<RegisterResponseData>(
        '/api/v1/auth/register',
        input,
    )
}