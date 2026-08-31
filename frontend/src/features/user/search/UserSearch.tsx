import {useEffect, useState} from 'react'
import {Search} from 'lucide-react'

import {startConversation, type Conversation} from '../../chat/conversationApi'
import {searchUsers, type SearchUser} from './userSearchApi'

type UserSearchProps = {
    onConversationStarted: (
        conversation: Conversation,
        user: SearchUser,
    ) => void
}

function UserSearch({onConversationStarted}: UserSearchProps) {
    const [query, setQuery] = useState('')
    const [users, setUsers] = useState<SearchUser[]>([])
    const [isLoading, setIsLoading] = useState(false)
    const [startingUserId, setStartingUserId] = useState<number | null>(null)
    const [errorMessage, setErrorMessage] = useState('')

    useEffect(() => {
        const searchTerm = query.trim()

        if (searchTerm.length < 2) {
            setUsers([])
            setErrorMessage('')
            return
        }

        const timer = window.setTimeout(async () => {
            setIsLoading(true)
            setErrorMessage('')

            try {
                const response = await searchUsers(searchTerm)

                if (response.success) {
                    setUsers(response.data)
                    return
                }

                setUsers([])
                setErrorMessage(response.error.message)
            } catch {
                setUsers([])
                setErrorMessage('Could not search users.')
            } finally {
                setIsLoading(false)
            }
        }, 300)

        return () => window.clearTimeout(timer)
    }, [query])

    async function handleSelectUser(user: SearchUser) {
        setStartingUserId(user.id)
        setErrorMessage('')

        try {
            const response = await startConversation(user.id)

            if (!response.success) {
                setErrorMessage(response.error.message)
                return
            }

            onConversationStarted(response.data, user)

            setQuery('')
            setUsers([])
        } catch {
            setErrorMessage('Could not start conversation.')
        } finally {
            setStartingUserId(null)
        }
    }

    return (
        <div className="user-search">
            <div className="sidebar-search">
                <Search size={17}/>

                <input
                    type="text"
                    value={query}
                    placeholder="Search people..."
                    onChange={(event) => setQuery(event.target.value)}
                />
            </div>

            {query.trim().length >= 2 && (
                <div className="search-results">
                    {isLoading && (
                        <div className="search-status">
                            Searching...
                        </div>
                    )}

                    {!isLoading && errorMessage && (
                        <div className="search-status">
                            {errorMessage}
                        </div>
                    )}

                    {!isLoading && !errorMessage && users.length === 0 && (
                        <div className="search-status">
                            No users found
                        </div>
                    )}

                    {!isLoading &&
                        users.map((user) => (
                            <button
                                key={user.id}
                                type="button"
                                className="conversation-item"
                                disabled={startingUserId === user.id}
                                onClick={() => handleSelectUser(user)}
                            >
                                <div className="avatar">
                                    {user.name.slice(0, 2).toUpperCase()}
                                </div>

                                <div className="conversation-content">
                                    <div className="conversation-top">
                                        <strong>{user.name}</strong>
                                    </div>

                                    <div className="conversation-bottom">
                                        <p>
                                            {startingUserId === user.id
                                                ? 'Opening conversation...'
                                                : `@${user.username}`}
                                        </p>
                                    </div>
                                </div>
                            </button>
                        ))}
                </div>
            )}
        </div>
    )
}

export default UserSearch