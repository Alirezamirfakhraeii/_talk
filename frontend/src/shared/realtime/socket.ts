export type RealtimeEvent<T = unknown> = {
    type: string
    data: T
}

type MessageHandler = (event: RealtimeEvent) => void

export function connectRealtime(onMessage: MessageHandler) {
    const protocol =
        window.location.protocol === 'https:' ? 'wss' : 'ws'

    const socket = new WebSocket(
        `${protocol}://${window.location.host}/api/v1/ws`,
    )

    socket.onopen = () => {
        console.log('Realtime connected')
    }

    socket.onmessage = (message) => {
        if (typeof message.data !== 'string') {
            return
        }

        try {
            const event = JSON.parse(
                message.data,
            ) as RealtimeEvent

            onMessage(event)
        } catch {
            console.error('Invalid realtime event')
        }
    }

    socket.onerror = () => {
        console.error('Realtime connection error')
    }

    socket.onclose = () => {
        console.log('Realtime disconnected')
    }

    return () => {
        socket.close()
    }
}