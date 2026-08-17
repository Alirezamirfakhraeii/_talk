import {useState} from 'react'

import {
    register,
    type RegisterResponseData,
} from './registerApi'

import './RegisterForm.css'

function RegisterForm() {
    const [name, setName] = useState('')
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const [isLoading, setIsLoading] = useState(false)

    const [successMessage, setSuccessMessage] =
        useState('')

    const [errorMessage, setErrorMessage] =
        useState('')

    const [fieldErrors, setFieldErrors] = useState<
        Record<string, string>
    >({})

    const [registeredUser, setRegisteredUser] =
        useState<RegisterResponseData | null>(null)

    async function handleSubmit(
        event: React.FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault()

        setIsLoading(true)
        setSuccessMessage('')
        setErrorMessage('')
        setFieldErrors({})
        setRegisteredUser(null)

        try {
            const response = await register({
                name,
                username,
                email,
                password,
            })

            if (!response.success) {
                const errors = {
                    ...(response.error.fields ?? {}),
                }

                if (
                    response.error.code ===
                    'EMAIL_ALREADY_EXISTS'
                ) {
                    errors.email =
                        'This email is already in use.'
                }

                if (
                    response.error.code ===
                    'USERNAME_ALREADY_EXISTS'
                ) {
                    errors.username =
                        'This username is already in use.'
                }

                setFieldErrors(errors)

                if (Object.keys(errors).length === 0) {
                    setErrorMessage(
                        response.error.message,
                    )
                }

                return
            }

            setSuccessMessage(
                response.message ??
                'Account created successfully',
            )

            setRegisteredUser(response.data)

            setName('')
            setUsername('')
            setEmail('')
            setPassword('')
        } catch {
            setErrorMessage(
                'Could not connect to the server.',
            )
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <form
            className="register-form"
            onSubmit={handleSubmit}
        >
            {successMessage && (
                <div className="register-alert register-alert-success">
                    {successMessage}
                </div>
            )}

            {errorMessage && (
                <div className="register-alert register-alert-error">
                    {errorMessage}
                </div>
            )}

            <div className="register-field">
                <label htmlFor="name">
                    Full name
                </label>

                <input
                    id="name"
                    type="text"
                    placeholder="Ali Reza"
                    value={name}
                    onChange={(event) =>
                        setName(event.target.value)
                    }
                />

                {fieldErrors.name && (
                    <span className="register-field-error">
            {fieldErrors.name}
          </span>
                )}
            </div>

            <div className="register-field">
                <label htmlFor="username">
                    Username
                </label>

                <input
                    id="username"
                    type="text"
                    placeholder="alirezamir"
                    value={username}
                    onChange={(event) =>
                        setUsername(event.target.value)
                    }
                />

                {fieldErrors.username && (
                    <span className="register-field-error">
            {fieldErrors.username}
          </span>
                )}
            </div>

            <div className="register-field">
                <label htmlFor="email">
                    Email address
                </label>

                <input
                    id="email"
                    type="email"
                    placeholder="name@example.com"
                    value={email}
                    onChange={(event) =>
                        setEmail(event.target.value)
                    }
                />

                {fieldErrors.email && (
                    <span className="register-field-error">
            {fieldErrors.email}
          </span>
                )}
            </div>

            <div className="register-field">
                <label htmlFor="password">
                    Password
                </label>

                <input
                    id="password"
                    type="password"
                    placeholder="Minimum 8 characters"
                    value={password}
                    onChange={(event) =>
                        setPassword(event.target.value)
                    }
                />

                {fieldErrors.password ? (
                    <span className="register-field-error">
            {fieldErrors.password}
          </span>
                ) : (
                    <span className="register-field-help">
            Use at least 8 characters.
          </span>
                )}
            </div>

            <button
                className="register-submit"
                type="submit"
                disabled={isLoading}
            >
                {isLoading
                    ? 'Creating account...'
                    : 'Create account'}

                {!isLoading && <span>→</span>}
            </button>

            {registeredUser && (
                <div className="register-created-user">
                    Welcome, {registeredUser.name}
                </div>
            )}
        </form>
    )
}

export default RegisterForm