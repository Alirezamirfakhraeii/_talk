import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'

import LoginPage from '../features/auth/login/LoginPage'
import RegisterPage from '../features/auth/register/RegisterPage'

function ChatsPage() {
    return (
        <main>
            <h1>Chats</h1>
            <p>Welcome to SamaTalk.</p>
        </main>
    )
}

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/login" replace/>}/>
                <Route path="/login" element={<LoginPage/>}/>
                <Route path="/register" element={<RegisterPage/>}/>
                <Route path="/chats" element={<ChatsPage/>}/>
            </Routes>
        </BrowserRouter>
    )
}

export default App