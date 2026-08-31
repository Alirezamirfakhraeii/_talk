import {get, post} from '../../shared/api/client'

export type Message = {
    id: number
    conversation_id: number
    sender_id: number
    content: string
    created_at: string
}

type SendMessageRequest = {
    content: string
}

export function getMessages(conversationId: number) {
    return get<Message[]>(
        `/api/v1/conversations/${conversationId}/messages`,
    )
}

export function sendMessage(
    conversationId: number,
    content: string,
) {
    const body: SendMessageRequest = {
        content,
    }

    return post<Message>(
        `/api/v1/conversations/${conversationId}/messages`,
        body,
    )
}