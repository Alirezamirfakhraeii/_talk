import RegisterForm from './RegisterForm'
import './RegisterPage.css'

function RegisterPage() {
    return (
        <main className="register-page">
            <section className="register-shell">
                <div className="register-intro">
                    <div className="register-brand">
                        <div className="register-brand-icon">
                            S
                        </div>

                        <span>SamaTalk</span>
                    </div>

                    <div className="register-intro-content">
            <span className="register-badge">
              Connect. Chat. Collaborate.
            </span>

                        <h1>
                            Conversations,
                            <br />
                            made simple.
                        </h1>

                        <p>
                            Create your account and start connecting
                            with people in a fast and secure environment.
                        </p>
                    </div>

                    <p className="register-copyright">
                        © 2026 SamaTalk
                    </p>
                </div>

                <div className="register-content">
                    <div className="register-form-wrapper">
                        <div className="register-heading">
              <span className="register-heading-label">
                GET STARTED
              </span>

                            <h2>Create your account</h2>

                            <p>
                                Enter your details below to join SamaTalk.
                            </p>
                        </div>

                        <RegisterForm />

                        <p className="register-login">
                            Already have an account?

                            <a href="/login">
                                Sign in
                            </a>
                        </p>
                    </div>
                </div>
            </section>
        </main>
    )
}

export default RegisterPage