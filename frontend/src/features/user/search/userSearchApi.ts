import {get} from '../../../shared/api/client'

export type SearchUser = {
    id: number
    name: string
    username: string
    avatar_path: string
}

export function searchUsers(query: string) {
    return get<SearchUser[]>(
        `/api/v1/users/search?q=${encodeURIComponent(query)}`,
    )
}

