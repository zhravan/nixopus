import { NextResponse } from 'next/server';

export async function GET() {
  return NextResponse.json({
    baseUrl: process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:8080/api',
    websocketUrl: process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost:8080/ws',
    webhookUrl: process.env.NEXT_PUBLIC_WEBHOOK_URL || 'http://localhost:8080/webhook',
    port: process.env.NEXT_PUBLIC_PORT || '7443'
  });
} 