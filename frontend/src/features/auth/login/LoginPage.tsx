import AuthLayout from '../../../shared/components/AuthLayout'

import LoginForm from './LoginForm'

function LoginPage() {
    return (
        <AuthLayout
            title="Welcome back"
            description="Sign in to continue to your SamaTalk account."
        >
            <LoginForm/>
        </AuthLayout>
    )
}

export default LoginPage