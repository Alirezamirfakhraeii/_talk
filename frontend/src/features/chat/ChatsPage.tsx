import {useEffect, useState} from 'react'

import {
    Bell,
    ChevronDown,
    MoreHorizontal,
    Paperclip,
    Phone,
    Plus,
    Send,
    Settings,
    Smile,
    Video,
} from 'lucide-react'

import UserSearch from '../user/search/UserSearch'
import type {SearchUser} from '../user/search/userSearchApi'

import {
    getConversations,
    type Conversation,
    type ConversationSummary,
} from './conversationApi'

import {
    getMessages,
    sendMessage,
    type Message,
} from './messageApi'

import './ChatsPage.css'

function ChatsPage() {
    const [activeConversationId, setActiveConversationId] =
        useState<number | null>(null)

    const [activeUser, setActiveUser] =
        useState<SearchUser | null>(null)

    const [conversations, setConversations] =
        useState<ConversationSummary[]>([])

    const [messages, setMessages] =
        useState<Message[]>([])

    const [messageText, setMessageText] =
        useState('')

    const [isLoadingConversations, setIsLoadingConversations] =
        useState(false)

    const [isLoadingMessages, setIsLoadingMessages] =
        useState(false)

    const [isSendingMessage, setIsSendingMessage] =
        useState(false)

    const [conversationsError, setConversationsError] =
        useState('')

    const [messagesError, setMessagesError] =
        useState('')

    async function loadConversations() {
        setIsLoadingConversations(true)
        setConversationsError('')

        try {
            const response = await getConversations()

            if (!response.success) {
                setConversationsError(response.error.message)
                return
            }

            setConversations(response.data)
        } catch {
            setConversationsError('Could not load conversations.')
        } finally {
            setIsLoadingConversations(false)
        }
    }

    useEffect(() => {
        loadConversations()
    }, [])

    useEffect(() => {
        if (!activeConversationId) {
            setMessages([])
            return
        }

        let isActive = true

        async function loadMessages() {
            setIsLoadingMessages(true)
            setMessagesError('')

            try {
                const response = await getMessages(
                    activeConversationId!,
                )

                if (!isActive) {
                    return
                }

                if (!response.success) {
                    setMessagesError(response.error.message)
                    setMessages([])
                    return
                }

                setMessages(response.data)
            } catch {
                if (isActive) {
                    setMessagesError('Could not load messages.')
                    setMessages([])
                }
            } finally {
                if (isActive) {
                    setIsLoadingMessages(false)
                }
            }
        }

        loadMessages()

        return () => {
            isActive = false
        }
    }, [activeConversationId])

    function handleConversationStarted(
        conversation: Conversation,
        user: SearchUser,
    ) {
        setActiveConversationId(conversation.id)
        setActiveUser(user)
        setMessageText('')
        setMessagesError('')

        loadConversations()
    }

    function handleConversationSelect(
        conversation: ConversationSummary,
    ) {
        setActiveConversationId(conversation.id)

        setActiveUser({
            id: conversation.user_id,
            name: conversation.name,
            username: conversation.username,
        })

        setMessageText('')
        setMessagesError('')
    }

    async function handleSendMessage(
        event: React.FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault()

        if (!activeConversationId) {
            return
        }

        const content = messageText.trim()

        if (!content || isSendingMessage) {
            return
        }

        setIsSendingMessage(true)
        setMessagesError('')

        try {
            const response = await sendMessage(
                activeConversationId,
                content,
            )

            if (!response.success) {
                setMessagesError(response.error.message)
                return
            }

            setMessages((currentMessages) => [
                ...currentMessages,
                response.data,
            ])

            setMessageText('')

            await loadConversations()
        } catch {
            setMessagesError('Could not send message.')
        } finally {
            setIsSendingMessage(false)
        }
    }

    function getInitials(name: string) {
        return name
            .split(' ')
            .map((part) => part.charAt(0))
            .join('')
            .slice(0, 2)
            .toUpperCase()
    }

    function formatTime(value: string | null) {
        if (!value) {
            return ''
        }

        return new Date(value).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
        })
    }

    return (
        <main className="chat-app">
            <aside className="chat-sidebar">
                <div className="sidebar-brand">
                    <div className="brand-logo">
                        S
                    </div>

                    <div className="brand-text">
                        <strong>SamaTalk</strong>
                        <span>Workspace</span>
                    </div>

                    <button
                        className="icon-button"
                        type="button"
                    >
                        <ChevronDown size={18}/>
                    </button>
                </div>

                <UserSearch
                    onConversationStarted={
                        handleConversationStarted
                    }
                />

                <div className="sidebar-section-header">
                    <span>Messages</span>

                    <button
                        className="icon-button small"
                        type="button"
                    >
                        <Plus size={17}/>
                    </button>
                </div>

                <div className="conversation-list">
                    {isLoadingConversations && (
                        <div className="conversation-list-status">
                            Loading conversations...
                        </div>
                    )}

                    {!isLoadingConversations &&
                        conversationsError && (
                            <div className="conversation-list-status error">
                                {conversationsError}
                            </div>
                        )}

                    {!isLoadingConversations &&
                        !conversationsError &&
                        conversations.length === 0 && (
                            <div className="conversation-list-status">
                                No conversations yet.
                            </div>
                        )}

                    {!isLoadingConversations &&
                        conversations.map((conversation) => (
                            <button
                                key={conversation.id}
                                type="button"
                                className={`conversation-item ${
                                    activeConversationId === conversation.id
                                        ? 'active'
                                        : ''
                                }`}
                                onClick={() =>
                                    handleConversationSelect(conversation)
                                }
                            >
                                <div className="avatar-wrapper">
                                    <div className="avatar">
                                        {getInitials(conversation.name)}
                                    </div>
                                </div>

                                <div className="conversation-content">
                                    <div className="conversation-top">
                                        <strong>
                                            {conversation.name}
                                        </strong>

                                        <span>
                      {formatTime(
                          conversation.last_message_time,
                      )}
                    </span>
                                    </div>

                                    <div className="conversation-bottom">
                                        <p>
                                            {conversation.last_message ??
                                                'No messages yet'}
                                        </p>
                                    </div>
                                </div>
                            </button>
                        ))}
                </div>

                <div className="sidebar-profile">
                    <div className="avatar avatar-current">
                        AM
                    </div>

                    <div className="profile-info">
                        <strong>Amir</strong>
                        <span>Online</span>
                    </div>

                    <button
                        className="icon-button"
                        type="button"
                    >
                        <Settings size={18}/>
                    </button>
                </div>
            </aside>

            <section className="conversation-panel">
                <header className="conversation-header">
                    <div className="conversation-user">
                        <div className="avatar-wrapper">
                            <div className="avatar">
                                {activeUser
                                    ? getInitials(activeUser.name)
                                    : '?'}
                            </div>
                        </div>

                        <div>
                            <strong>
                                {activeUser?.name ??
                                    'Select a conversation'}
                            </strong>

                            <span>
                {activeUser
                    ? `@${activeUser.username}`
                    : 'Search for someone to start chatting'}
              </span>
                        </div>
                    </div>

                    <div className="conversation-actions">
                        <button
                            className="icon-button action"
                            type="button"
                        >
                            <Phone size={18}/>
                        </button>

                        <button
                            className="icon-button action"
                            type="button"
                        >
                            <Video size={19}/>
                        </button>

                        <button
                            className="icon-button action"
                            type="button"
                        >
                            <Bell size={18}/>
                        </button>

                        <button
                            className="icon-button action"
                            type="button"
                        >
                            <MoreHorizontal size={20}/>
                        </button>
                    </div>
                </header>

                <div className="messages-area">
                    {!activeConversationId && (
                        <div className="messages-status">
                            Select a conversation to start messaging.
                        </div>
                    )}

                    {activeConversationId &&
                        isLoadingMessages && (
                            <div className="messages-status">
                                Loading messages...
                            </div>
                        )}

                    {activeConversationId &&
                        messagesError && (
                            <div className="messages-status messages-status-error">
                                {messagesError}
                            </div>
                        )}

                    {activeConversationId &&
                        !isLoadingMessages &&
                        !messagesError &&
                        messages.length === 0 && (
                            <div className="messages-status">
                                No messages yet. Start the conversation.
                            </div>
                        )}

                    {activeConversationId &&
                        !isLoadingMessages &&
                        messages.map((message) => {
                            const isReceived =
                                message.sender_id === activeUser?.id

                            return (
                                <div
                                    key={message.id}
                                    className={`message-row ${
                                        isReceived
                                            ? 'received'
                                            : 'sent'
                                    }`}
                                >
                                    {isReceived && (
                                        <div className="avatar message-avatar">
                                            {activeUser
                                                ? getInitials(activeUser.name)
                                                : '?'}
                                        </div>
                                    )}

                                    <div className="message-group">
                                        <div
                                            className={`message-bubble ${
                                                isReceived
                                                    ? 'received-bubble'
                                                    : 'sent-bubble'
                                            }`}
                                        >
                                            {message.content}
                                        </div>

                                        <span className="message-time">
                      {formatTime(
                          message.created_at,
                      )}
                    </span>
                                    </div>
                                </div>
                            )
                        })}
                </div>

                <div className="composer-wrapper">
                    <form
                        className="composer"
                        onSubmit={handleSendMessage}
                    >
                        <button
                            type="button"
                            className="composer-button"
                            disabled={!activeConversationId}
                        >
                            <Plus size={19}/>
                        </button>

                        <button
                            type="button"
                            className="composer-button"
                            disabled={!activeConversationId}
                        >
                            <Paperclip size={18}/>
                        </button>

                        <input
                            type="text"
                            value={messageText}
                            placeholder={
                                activeConversationId
                                    ? 'Write a message...'
                                    : 'Select a conversation first'
                            }
                            disabled={
                                !activeConversationId ||
                                isSendingMessage
                            }
                            onChange={(event) =>
                                setMessageText(event.target.value)
                            }
                        />

                        <button
                            type="button"
                            className="composer-button"
                            disabled={!activeConversationId}
                        >
                            <Smile size={19}/>
                        </button>

                        <button
                            type="submit"
                            className="send-button"
                            disabled={
                                !activeConversationId ||
                                !messageText.trim() ||
                                isSendingMessage
                            }
                        >
                            <Send size={18}/>
                        </button>
                    </form>
                </div>
            </section>
        </main>
    )
}

export default ChatsPage