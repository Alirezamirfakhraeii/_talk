export type UserProfile = {
    id: number
    name: string
    username: string
    email: string
    bio: string
    avatar_path: string
}

export type UpdateProfileRequest = {
    name: string
    username: string
    bio: string
}

type ApiSuccessResponse<T> = {
    success: true
    message: string
    data: T
}

type ApiErrorResponse = {
    success: false
    error: {
        code: string
        message: string
        details?: unknown
    }
}

export type ApiResponse<T> =
    | ApiSuccessResponse<T>
    | ApiErrorResponse

async function parseResponse<T>(
    response: Response,
): Promise<ApiResponse<T>> {
    return response.json() as Promise<ApiResponse<T>>
}

export async function getProfile() {
    const response = await fetch('/api/v1/me', {
        method: 'GET',
        credentials: 'include',
    })

    return parseResponse<UserProfile>(response)
}

export async function updateProfile(
    input: UpdateProfileRequest,
) {
    const response = await fetch(
        '/api/v1/me/profile',
        {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(input),
        },
    )

    return parseResponse<UserProfile>(response)
}

export async function uploadAvatar(
    avatar: File,
) {
    const formData = new FormData()

    formData.append('avatar', avatar)

    const response = await fetch(
        '/api/v1/me/avatar',
        {
            method: 'POST',
            credentials: 'include',
            body: formData,
        },
    )

    return parseResponse<UserProfile>(response)
}