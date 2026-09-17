SamaTalk

SamaTalk is a real-time private messaging application built with Go, PostgreSQL, WebSocket, React, and TypeScript.

The main goal of this project was to learn backend development with Go using a clean modular structure and Go's standard net/http package.

Features

User registration and login

Session-based authentication with HttpOnly cookies

User profile editing

Avatar upload

User search

One-to-one conversations

Message history

Real-time messaging with WebSocket

Multiple WebSocket connections per user

PostgreSQL persistence

Health check endpoints

React chat interface

Tech Stack

Backend

Go

net/http

PostgreSQL

pgx / pgxpool

Gorilla WebSocket

Goose migrations

Frontend

React

TypeScript

Vite

Lucide React

Architecture

HTTP Request
↓
Handler
↓
Service
↓
Repository
↓
PostgreSQL

Realtime messaging follows this flow:

Send Message
↓
Go API
↓
PostgreSQL
↓
Realtime Hub
↓
WebSocket
↓
Recipient

PostgreSQL is the source of truth, while WebSocket is used for realtime delivery.

Main API

POST /api/v1/auth/register
POST /api/v1/auth/login

GET  /api/v1/me
PUT  /api/v1/me/profile
POST /api/v1/me/avatar

GET  /api/v1/users/search

GET  /api/v1/conversations
POST /api/v1/conversations

GET  /api/v1/conversations/{id}/messages
POST /api/v1/conversations/{id}/messages

GET  /api/v1/ws

Health checks:

GET /health/live
GET /health/ready

Project Structure

samatalk/
├── backend/
│   ├── cmd/api
│   ├── internal/
│   │   ├── app
│   │   ├── auth
│   │   ├── conversation
│   │   ├── message
│   │   ├── realtime
│   │   └── user
│   └── migrations
│
└── frontend/
└── src/

Running Locally

Start PostgreSQL:

cd backend
docker compose up -d postgres

Run backend:

go run ./cmd/api

Run frontend:

cd frontend
npm install
npm run dev

What I Learned

This project helped me practice:

Go project structure

REST API development

authentication and middleware

PostgreSQL and connection pooling

dependency injection

context and graceful shutdown

goroutines and channels

WebSocket architecture

realtime messaging

file uploads

frontend/backend integration

Status

The core private messaging flow is complete:

Register
↓
Login
↓
Search User
↓
Start Conversation
↓
Send Message
↓
PostgreSQL
↓
Realtime WebSocket Delivery

SamaTalk is my first complete Go backend project focused on learning Go, backend architecture, and realtime systems.