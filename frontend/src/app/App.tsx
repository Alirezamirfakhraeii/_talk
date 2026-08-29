import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'

import LoginPage from '../features/auth/login/LoginPage'
import RegisterPage from '../features/auth/register/RegisterPage'
import ProtectedRoute from '../features/auth/ProtectedRoute'
import ChatsPage from '../features/chat/ChatsPage'

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/login" replace/>}/>
                <Route path="/login" element={<LoginPage/>}/>
                <Route path="/register" element={<RegisterPage/>}/>
                <Route path="/chats" element={<ProtectedRoute><ChatsPage/></ProtectedRoute>}/>
            </Routes>
        </BrowserRouter>
    )
}

export default App