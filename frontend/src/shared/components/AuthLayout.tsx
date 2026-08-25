import type {ReactNode} from 'react'

import './AuthLayout.css'

type AuthLayoutProps = {
    title: string
    description: string
    children: ReactNode
}

function AuthLayout({
                        title,
                        description,
                        children,
                    }: AuthLayoutProps) {
    return (
        <main className="auth-page">
            <section className="auth-card">
                <div className="auth-brand">
                    <div className="auth-logo">S</div>
                    <span>SamaTalk</span>
                </div>

                <div className="auth-header">
                    <h1>{title}</h1>
                    <p>{description}</p>
                </div>

                {children}
            </section>
        </main>
    )
}

export default AuthLayout