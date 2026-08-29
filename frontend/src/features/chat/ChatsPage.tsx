import {
    Bell,
    ChevronDown,
    MoreHorizontal,
    Paperclip,
    Phone,
    Plus,
    Search,
    Send,
    Settings,
    Smile,
    Video,
} from 'lucide-react'

import './ChatsPage.css'

function ChatsPage() {
    return (
        <main className="chat-app">
            <aside className="chat-sidebar">
                <div className="sidebar-brand">
                    <div className="brand-logo">S</div>

                    <div className="brand-text">
                        <strong>SamaTalk</strong>
                        <span>Workspace</span>
                    </div>

                    <button className="icon-button">
                        <ChevronDown size={18}/>
                    </button>
                </div>

                <div className="sidebar-search">
                    <Search size={17}/>

                    <input
                        type="text"
                        placeholder="Search conversations"
                    />
                </div>

                <div className="sidebar-section-header">
                    <span>Messages</span>

                    <button className="icon-button small">
                        <Plus size={17}/>
                    </button>
                </div>

                <div className="conversation-list">
                    <button className="conversation-item active">
                        <div className="avatar-wrapper">
                            <div className="avatar">AR</div>
                            <span className="online-dot"/>
                        </div>

                        <div className="conversation-content">
                            <div className="conversation-top">
                                <strong>Ali Reza</strong>
                                <span>14:32</span>
                            </div>

                            <div className="conversation-bottom">
                                <p>Sounds good. See you there!</p>
                                <span className="unread-badge">2</span>
                            </div>
                        </div>
                    </button>

                    <button className="conversation-item">
                        <div className="avatar-wrapper">
                            <div className="avatar avatar-purple">SH</div>
                            <span className="online-dot"/>
                        </div>

                        <div className="conversation-content">
                            <div className="conversation-top">
                                <strong>Shima</strong>
                                <span>13:05</span>
                            </div>

                            <div className="conversation-bottom">
                                <p>Thank you 🤍</p>
                            </div>
                        </div>
                    </button>

                    <button className="conversation-item">
                        <div className="avatar-wrapper">
                            <div className="avatar avatar-orange">SA</div>
                        </div>

                        <div className="conversation-content">
                            <div className="conversation-top">
                                <strong>Sara Ahmadi</strong>
                                <span>Yesterday</span>
                            </div>

                            <div className="conversation-bottom">
                                <p>I sent you the files.</p>
                            </div>
                        </div>
                    </button>
                </div>

                <div className="sidebar-profile">
                    <div className="avatar avatar-current">AM</div>

                    <div className="profile-info">
                        <strong>Amir</strong>
                        <span>Online</span>
                    </div>

                    <button className="icon-button">
                        <Settings size={18}/>
                    </button>
                </div>
            </aside>

            <section className="conversation-panel">
                <header className="conversation-header">
                    <div className="conversation-user">
                        <div className="avatar-wrapper">
                            <div className="avatar">AR</div>
                            <span className="online-dot"/>
                        </div>

                        <div>
                            <strong>Ali Reza</strong>
                            <span>Active now</span>
                        </div>
                    </div>

                    <div className="conversation-actions">
                        <button className="icon-button action">
                            <Phone size={18}/>
                        </button>

                        <button className="icon-button action">
                            <Video size={19}/>
                        </button>

                        <button className="icon-button action">
                            <Bell size={18}/>
                        </button>

                        <button className="icon-button action">
                            <MoreHorizontal size={20}/>
                        </button>
                    </div>
                </header>

                <div className="messages-area">
                    <div className="date-divider">
                        <span>Today</span>
                    </div>

                    <div className="message-row received">
                        <div className="avatar message-avatar">AR</div>

                        <div className="message-group">
                            <div className="message-bubble received-bubble">
                                Hey! How is the SamaTalk project going?
                            </div>

                            <span className="message-time">14:28</span>
                        </div>
                    </div>

                    <div className="message-row sent">
                        <div className="message-group">
                            <div className="message-bubble sent-bubble">
                                Pretty good! I'm working on the chat interface right now.
                            </div>

                            <span className="message-time">
                14:29 · Seen
              </span>
                        </div>
                    </div>

                    <div className="message-row received">
                        <div className="avatar message-avatar">AR</div>

                        <div className="message-group">
                            <div className="message-bubble received-bubble">
                                Nice. Are you planning realtime messaging too?
                            </div>

                            <span className="message-time">14:31</span>
                        </div>
                    </div>

                    <div className="message-row sent">
                        <div className="message-group">
                            <div className="message-bubble sent-bubble">
                                Yes. WebSocket is coming next 🚀
                            </div>

                            <span className="message-time">
                14:32 · Seen
              </span>
                        </div>
                    </div>
                </div>

                <div className="composer-wrapper">
                    <form className="composer">
                        <button type="button" className="composer-button">
                            <Plus size={19}/>
                        </button>

                        <button type="button" className="composer-button">
                            <Paperclip size={18}/>
                        </button>

                        <input
                            type="text"
                            placeholder="Write a message..."
                        />

                        <button type="button" className="composer-button">
                            <Smile size={19}/>
                        </button>

                        <button type="submit" className="send-button">
                            <Send size={18}/>
                        </button>
                    </form>
                </div>
            </section>
        </main>
    )
}

export default ChatsPage