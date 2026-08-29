import {useEffect, useState} from 'react'
import type {ReactNode} from 'react'
import {Navigate} from 'react-router-dom'

import {getCurrentUser} from './me/meApi'

type ProtectedRouteProps = {
    children: ReactNode
}

type AuthStatus =
    | 'loading'
    | 'authenticated'
    | 'unauthenticated'

function ProtectedRoute({
                            children,
                        }: ProtectedRouteProps) {
    const [status, setStatus] =
        useState<AuthStatus>('loading')

    useEffect(() => {
        let isActive = true

        async function checkAuthentication() {
            try {
                const response = await getCurrentUser()

                if (!isActive) {
                    return
                }

                if (response.success) {
                    setStatus('authenticated')
                    return
                }

                setStatus('unauthenticated')
            } catch {
                if (isActive) {
                    setStatus('unauthenticated')
                }
            }
        }

        checkAuthentication()

        return () => {
            isActive = false
        }
    }, [])

    if (status === 'loading') {
        return <p>Loading...</p>
    }

    if (status === 'unauthenticated') {
        return (
            <Navigate
                to="/login"
                replace
            />
        )
    }

    return children
}

export default ProtectedRoute