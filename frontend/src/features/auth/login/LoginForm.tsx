import type {FormEvent} from 'react'
import {useState} from 'react'
import {useNavigate} from 'react-router-dom'

import {getCurrentUser} from '../me/meApi'
import {login} from './loginApi'

import '../../../shared/styles/AuthForm.css'

function LoginForm() {
    const navigate = useNavigate()

    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [isLoading, setIsLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState('')

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        setIsLoading(true)
        setErrorMessage('')

        try {
            const response = await login({
                email,
                password,
            })

            if (!response.success) {
                if (response.error.code === 'INVALID_CREDENTIALS') {
                    setErrorMessage('Email or password is incorrect.')
                    return
                }

                setErrorMessage(response.error.message)
                return
            }

            const meResponse = await getCurrentUser()

            if (!meResponse.success) {
                setErrorMessage('Your session could not be verified.')
                return
            }

            navigate('/chats')
        } catch {
            setErrorMessage('Could not connect to the server.')
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <form className="auth-form" onSubmit={handleSubmit}>
            {errorMessage && (
                <div className="auth-alert auth-alert-error">
                    {errorMessage}
                </div>
            )}

            <div className="auth-field">
                <label htmlFor="email">Email address</label>

                <input
                    id="email"
                    type="email"
                    placeholder="name@example.com"
                    value={email}
                    autoComplete="email"
                    required
                    onChange={(event) => setEmail(event.target.value)}
                />
            </div>

            <div className="auth-field">
                <label htmlFor="password">Password</label>

                <input
                    id="password"
                    type="password"
                    placeholder="Enter your password"
                    value={password}
                    autoComplete="current-password"
                    required
                    onChange={(event) => setPassword(event.target.value)}
                />
            </div>

            <button
                className="auth-submit"
                type="submit"
                disabled={isLoading}
            >
                {isLoading ? 'Signing in...' : 'Sign in'}
            </button>

            <div className="auth-footer">
                <span>Don't have an account?</span>

                <button
                    className="auth-link"
                    type="button"
                    onClick={() => navigate('/register')}
                >
                    Create account
                </button>
            </div>
        </form>
    )
}

export default LoginForm