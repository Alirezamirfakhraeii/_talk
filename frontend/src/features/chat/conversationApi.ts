import {get, post} from '../../shared/api/client'

export type Conversation = {
    id: number
    user_one_id: number
    user_two_id: number
}

export type ConversationSummary = {
    id: number
    user_id: number
    name: string
    username: string
    last_message: string | null
    last_message_time: string | null
}

type StartConversationRequest = {
    user_id: number
}

export function getConversations() {
    return get<ConversationSummary[]>('/api/v1/conversations')
}

export function startConversation(userId: number) {
    const body: StartConversationRequest = {
        user_id: userId,
    }

    return post<Conversation>(
        '/api/v1/conversations',
        body,
    )
}