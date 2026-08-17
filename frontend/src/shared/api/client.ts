export type ApiErrorResponse = {
    success: false
    error: {
        code: string
        message: string
        fields?: Record<string, string>
    }
}

export type ApiSuccessResponse<T> = {
    success: true
    message?: string
    data: T
}

export type ApiResponse<T> =
    | ApiSuccessResponse<T>
    | ApiErrorResponse

export async function post<T>(
    url: string,
    body: unknown,
): Promise<ApiResponse<T>> {
    const response = await fetch(url, {
        method: 'POST',

        headers: {
            'Content-Type': 'application/json',
        },

        body: JSON.stringify(body),
    })

    const data = await response.json()

    return data as ApiResponse<T>
}