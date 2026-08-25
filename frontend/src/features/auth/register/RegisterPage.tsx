import AuthLayout from '../../../shared/components/AuthLayout'

import RegisterForm from './RegisterForm'

function RegisterPage() {
    return (
        <AuthLayout
            title="Create your account"
            description="Create an account and start using SamaTalk."
        >
            <RegisterForm/>
        </AuthLayout>
    )
}

export default RegisterPage