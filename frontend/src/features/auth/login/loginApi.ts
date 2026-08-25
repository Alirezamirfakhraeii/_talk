import { post } from '../../../shared/api/client'

export type LoginRequest = {
    email: string
    password: string
}

export type LoginResponseData = {
    id: number
    name: string
    username: string
    email: string
}

export function login(input: LoginRequest) {
    return post<LoginResponseData>(
        '/api/v1/auth/login',
        input,
    )
}