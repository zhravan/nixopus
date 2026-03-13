'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import { useAppSelector } from '@/redux/hooks';
import { authClient } from '@/packages/lib/auth-client';
import {
  createAgentClient,
  AGENT_ID,
  streamAgent,
  approveAgentToolCall,
  declineAgentToolCall,
  type StreamChunk
} from '@/packages/lib/agent-client';
import { useAgentConfigured } from '@/packages/hooks/shared/use-config';
import { type ChatContext, formatContextsForAgent } from './chat-context';

export type MessagePart =
  | { type: 'text'; content: string }
  | {
      type: 'tool-call';
      toolName: string;
      toolCallId: string;
      args?: unknown;
      status: 'running' | 'done';
    };

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  parts?: MessagePart[];
  timestamp: Date;
  contexts?: ChatContext[];
  kind?: 'status';
}

interface UseAgentChatOptions {
  threadId: string | null;
  resourceId?: string;
  contexts?: ChatContext[];
  autoRunTools?: boolean;
  onFirstMessage?: (content: string) => void;
  waitForThread?: (id: string) => Promise<void>;
}

async function getAuthHeaders(
  token: string | null,
  organizationId: string | undefined | null
): Promise<Record<string, string>> {
  const headers: Record<string, string> = {};

  let authToken = token;
  if (!authToken) {
    try {
      const session = await authClient.getSession();
      authToken = session?.data?.session?.token ?? null;
    } catch {
      // ignore
    }
  }

  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }
  if (organizationId) {
    headers['X-Organization-Id'] = organizationId;
  }

  return headers;
}

function extractText(content: unknown): string {
  if (typeof content === 'string') return content;
  if (Array.isArray(content)) {
    return content
      .filter((p: { type?: string; text?: string }) => p.type === 'text' && p.text)
      .map((p: { text: string }) => p.text)
      .join('');
  }
  if (content && typeof content === 'object') {
    const obj = content as { content?: string; parts?: { type?: string; text?: string }[] };
    if (typeof obj.content === 'string') return obj.content;
    if (Array.isArray(obj.parts)) return extractText(obj.parts);
  }
  return '';
}

function extractMessageParts(content: unknown): { text: string; parts: MessagePart[] } {
  if (typeof content === 'string') return { text: content, parts: [] };
  if (!content || typeof content !== 'object') return { text: '', parts: [] };

  const obj = content as {
    format?: number;
    content?: string;
    parts?: Array<Record<string, unknown>>;
    toolInvocations?: Array<{
      toolCallId?: string;
      toolName?: string;
      args?: unknown;
      state?: string;
    }>;
  };

  const parts: MessagePart[] = [];
  let text = '';

  if (Array.isArray(obj.parts)) {
    for (const part of obj.parts) {
      if (part.type === 'text') {
        const t = (part.text as string) || (part.content as string) || '';
        if (t) {
          text += t;
          parts.push({ type: 'text', content: t });
        }
      } else if (part.type === 'tool-invocation') {
        const inv = part.toolInvocation as
          | {
              toolCallId?: string;
              toolName?: string;
              args?: unknown;
              state?: string;
            }
          | undefined;
        if (inv) {
          parts.push({
            type: 'tool-call',
            toolName: inv.toolName || 'tool',
            toolCallId: inv.toolCallId || '',
            args: inv.args,
            status: 'done'
          });
        }
      }
    }
  }

  if (parts.every((p) => p.type === 'tool-call') && Array.isArray(obj.toolInvocations)) {
    for (const inv of obj.toolInvocations) {
      if (!parts.some((p) => p.type === 'tool-call' && p.toolCallId === inv.toolCallId)) {
        parts.push({
          type: 'tool-call',
          toolName: inv.toolName || 'tool',
          toolCallId: inv.toolCallId || '',
          args: inv.args,
          status: 'done'
        });
      }
    }
  }

  if (!text && typeof obj.content === 'string') {
    text = obj.content;
  }

  return { text, parts };
}

export interface PendingToolApproval {
  runId: string;
  toolCallId: string;
  toolName: string;
  args: unknown;
}

export interface OmStatus {
  messages: { tokens: number; threshold: number };
  observations: { tokens: number; threshold: number };
  isObserving: boolean;
  observationsText: string | null;
}

export function useAgentChat({
  threadId,
  resourceId,
  contexts = [],
  autoRunTools = false,
  onFirstMessage,
  waitForThread
}: UseAgentChatOptions) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [pendingToolApproval, setPendingToolApproval] = useState<PendingToolApproval | null>(null);
  const [omStatus, setOmStatus] = useState<OmStatus | null>(null);
  const omStatusRef = useRef<OmStatus | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const abortRef = useRef<AbortController | null>(null);
  const pendingApprovalRef = useRef(false);
  const needsStepSeparatorRef = useRef(false);
  const needsNewTextPartRef = useRef(false);

  const token = useAppSelector((state) => state.auth.token);
  const activeOrg = useAppSelector((state) => state.user.activeOrganization);
  const organizationId = activeOrg?.id;

  const isAgentEnabled = useAgentConfigured() === true;

  const scrollToBottom = useCallback(() => {
    if (scrollRef.current) {
      const el = scrollRef.current;
      requestAnimationFrame(() => {
        el.scrollTop = el.scrollHeight;
      });
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  useEffect(() => {
    if (!threadId || !isAgentEnabled) {
      setMessages([]);
      return;
    }

    let cancelled = false;
    setIsLoadingHistory(true);
    setMessages([]);

    (async () => {
      try {
        if (waitForThread) await waitForThread(threadId);
        if (cancelled) return;

        const headers = await getAuthHeaders(token ?? null, organizationId ?? null);
        const client = createAgentClient(headers);
        const thread = client.getMemoryThread({ threadId, agentId: AGENT_ID });
        const result = await thread.listMessages({
          resourceId: resourceId || threadId
        });

        if (cancelled) return;

        const msgs: ChatMessage[] = [];
        const rawMessages = result?.messages ?? [];
        for (const msg of rawMessages) {
          const role = msg.role === 'user' ? 'user' : msg.role === 'assistant' ? 'assistant' : null;
          if (!role) continue;

          if (role === 'assistant') {
            const { text, parts } = extractMessageParts(msg.content);
            if (!text && parts.length === 0) continue;
            msgs.push({
              id: msg.id,
              role,
              content: text,
              ...(parts.length > 0 ? { parts } : {}),
              timestamp: msg.createdAt ? new Date(msg.createdAt) : new Date()
            });
          } else {
            const text = extractText(msg.content);
            if (!text) continue;
            msgs.push({
              id: msg.id,
              role,
              content: text,
              timestamp: msg.createdAt ? new Date(msg.createdAt) : new Date()
            });
          }
        }
        setMessages(msgs);
      } catch {
        // thread may not exist on server yet
      } finally {
        if (!cancelled) setIsLoadingHistory(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [threadId, resourceId, token, organizationId, isAgentEnabled, waitForThread]);

  const handleChunk = useCallback(
    (
      chunk: StreamChunk,
      assistantMessageId: string,
      runIdRef: { current: string },
      abortSignal: AbortSignal
    ) => {
      if (abortSignal.aborted) return;

      if (chunk.type === 'start' && chunk.runId) {
        runIdRef.current = chunk.runId;
      }

      if (chunk.type === 'text-delta') {
        const text = chunk.payload?.text as string | undefined;
        if (text) {
          const insertSep = needsStepSeparatorRef.current;
          const startNewPart = needsNewTextPartRef.current;
          needsStepSeparatorRef.current = false;
          needsNewTextPartRef.current = false;
          setMessages((prev) =>
            prev.map((m) => {
              if (m.id !== assistantMessageId) return m;
              const sep = insertSep && m.content.length > 0 ? '\n\n' : '';
              const newContent = m.content + sep + text;

              let parts = [...(m.parts || [])];
              parts = parts.map((p) =>
                p.type === 'tool-call' && p.status === 'running'
                  ? { ...p, status: 'done' as const }
                  : p
              );
              const lastPart = parts[parts.length - 1];
              if (!startNewPart && lastPart?.type === 'text') {
                parts[parts.length - 1] = { ...lastPart, content: lastPart.content + text };
              } else {
                parts.push({ type: 'text' as const, content: text });
              }

              return { ...m, content: newContent, parts };
            })
          );
        }
      }

      if (chunk.type === 'text-end' || chunk.type === 'step-finish') {
        needsStepSeparatorRef.current = true;
        needsNewTextPartRef.current = true;
      }

      if (
        chunk.type === 'tool-call' ||
        chunk.type === 'tool-call-start' ||
        chunk.type === 'tool-call-approval'
      ) {
        const p = chunk.payload as
          | {
              toolCallId?: string;
              toolName?: string;
              args?: unknown;
              runId?: string;
              id?: string;
            }
          | undefined;
        const toolCallId = p?.toolCallId ?? p?.id;

        if (toolCallId) {
          needsNewTextPartRef.current = true;
          const toolName = (p?.toolName as string) ?? 'tool';
          const tcId = String(toolCallId);
          setMessages((prev) =>
            prev.map((m) => {
              if (m.id !== assistantMessageId) return m;
              const parts = [...(m.parts || [])];
              if (!parts.some((pt) => pt.type === 'tool-call' && pt.toolCallId === tcId)) {
                parts.push({
                  type: 'tool-call' as const,
                  toolName,
                  toolCallId: tcId,
                  args: p?.args,
                  status: 'running' as const
                });
              }
              return { ...m, parts };
            })
          );
        }
      }

      if (chunk.type === 'tool-result') {
        const p = chunk.payload as { toolCallId?: string } | undefined;
        const tcId = p?.toolCallId;
        if (tcId) {
          needsNewTextPartRef.current = true;
          setMessages((prev) =>
            prev.map((m) => {
              if (m.id !== assistantMessageId) return m;
              const parts = (m.parts || []).map((pt) =>
                pt.type === 'tool-call' && pt.toolCallId === tcId
                  ? { ...pt, status: 'done' as const }
                  : pt
              );
              return { ...m, parts };
            })
          );
        }
      }

      if (chunk.type === 'finish' && chunk.payload) {
        const finishReason = (chunk.payload as { finishReason?: string }).finishReason;
        const suspendPayload = (
          chunk.payload as {
            suspendPayload?: {
              toolCallId?: string;
              runId?: string;
              toolName?: string;
              args?: unknown;
            };
          }
        ).suspendPayload;
        if (finishReason === 'suspended' && suspendPayload) {
          pendingApprovalRef.current = true;
          setPendingToolApproval({
            runId: suspendPayload.runId ?? '',
            toolCallId: suspendPayload.toolCallId ?? '',
            toolName: (suspendPayload.toolName as string) ?? 'tool',
            args: suspendPayload.args ?? {}
          });
        }
      }

      if (chunk.type === 'data-om-status') {
        const d = chunk.payload as
          | {
              windows?: {
                active?: {
                  messages?: { tokens?: number; threshold?: number };
                  observations?: { tokens?: number; threshold?: number };
                };
              };
            }
          | undefined;
        const active = d?.windows?.active;
        if (active) {
          const next: OmStatus = {
            messages: {
              tokens: active.messages?.tokens ?? 0,
              threshold: active.messages?.threshold ?? 30000
            },
            observations: {
              tokens: active.observations?.tokens ?? 0,
              threshold: active.observations?.threshold ?? 40000
            },
            isObserving: omStatusRef.current?.isObserving ?? false,
            observationsText: omStatusRef.current?.observationsText ?? null
          };
          omStatusRef.current = next;
          setOmStatus(next);
        }
      }

      if (chunk.type === 'data-om-observation-start') {
        const prev = omStatusRef.current;
        if (prev) {
          const next = { ...prev, isObserving: true };
          omStatusRef.current = next;
          setOmStatus(next);
        }
      }

      if (chunk.type === 'data-om-observation-end' || chunk.type === 'data-om-activation') {
        const d = chunk.payload as { observations?: string } | undefined;
        const prev = omStatusRef.current;
        if (prev) {
          const next: OmStatus = {
            ...prev,
            isObserving: false,
            observationsText: d?.observations ?? prev.observationsText
          };
          omStatusRef.current = next;
          setOmStatus(next);
        }
      }
    },
    []
  );

  const streamResponse = useCallback(
    async (userContent: string) => {
      if (!isAgentEnabled || !threadId) return;

      const headers = await getAuthHeaders(token ?? null, organizationId ?? null);

      const abortController = new AbortController();
      abortRef.current = abortController;
      const runIdRef = { current: '' };

      const assistantMessageId = crypto.randomUUID();
      setMessages((prev) => [
        ...prev,
        { id: assistantMessageId, role: 'assistant', content: '', timestamp: new Date(), parts: [] }
      ]);

      try {
        const contextPrefix = formatContextsForAgent(contexts);
        const stream = streamAgent(
          contextPrefix + userContent,
          threadId,
          resourceId || threadId,
          headers,
          abortController.signal
        );

        for await (const chunk of stream) {
          handleChunk(chunk, assistantMessageId, runIdRef, abortController.signal);
        }
      } catch (err: unknown) {
        if (err instanceof Error && err.name === 'AbortError') return;

        const errorMessage =
          err instanceof Error ? err.message : 'Failed to get response from AI agent';

        setMessages((prev) =>
          prev.map((m) =>
            m.id === assistantMessageId
              ? { ...m, content: m.content || `Error: ${errorMessage}` }
              : m
          )
        );
      } finally {
        if (!pendingApprovalRef.current) {
          setIsStreaming(false);
        }
        abortRef.current = null;
        needsStepSeparatorRef.current = false;
        needsNewTextPartRef.current = false;
        setMessages((prev) =>
          prev.map((m) => {
            if (m.id !== assistantMessageId || !m.parts) return m;
            const hasRunning = m.parts.some(
              (p) => p.type === 'tool-call' && p.status === 'running'
            );
            if (!hasRunning) return m;
            return {
              ...m,
              parts: m.parts.map((p) =>
                p.type === 'tool-call' && p.status === 'running'
                  ? { ...p, status: 'done' as const }
                  : p
              )
            };
          })
        );
      }
    },
    [threadId, resourceId, token, organizationId, isAgentEnabled, contexts, handleChunk]
  );

  const handleApproveToolCall = useCallback(async () => {
    const pending = pendingToolApproval;
    if (!pending || !isAgentEnabled) return;

    setPendingToolApproval(null);
    pendingApprovalRef.current = false;

    const assistantMessageId = messages.filter((m) => m.role === 'assistant').pop()?.id;
    if (!assistantMessageId) {
      setIsStreaming(false);
      return;
    }

    const abortController = new AbortController();
    abortRef.current = abortController;
    const runIdRef = { current: '' };
    needsNewTextPartRef.current = true;

    try {
      const headers = await getAuthHeaders(token ?? null, organizationId ?? null);
      const stream = approveAgentToolCall(
        { runId: pending.runId, toolCallId: pending.toolCallId },
        headers,
        abortController.signal
      );

      for await (const chunk of stream) {
        handleChunk(chunk, assistantMessageId, runIdRef, abortController.signal);
      }
    } catch {
      setMessages((prev) =>
        prev.map((m) =>
          m.id === assistantMessageId
            ? { ...m, content: m.content + '\n\n_Tool execution failed._' }
            : m
        )
      );
    } finally {
      if (!pendingApprovalRef.current) setIsStreaming(false);
      abortRef.current = null;
      needsStepSeparatorRef.current = false;
      needsNewTextPartRef.current = false;
      if (assistantMessageId) {
        setMessages((prev) =>
          prev.map((m) => {
            if (m.id !== assistantMessageId || !m.parts) return m;
            const hasRunning = m.parts.some(
              (p) => p.type === 'tool-call' && p.status === 'running'
            );
            if (!hasRunning) return m;
            return {
              ...m,
              parts: m.parts.map((p) =>
                p.type === 'tool-call' && p.status === 'running'
                  ? { ...p, status: 'done' as const }
                  : p
              )
            };
          })
        );
      }
    }
  }, [pendingToolApproval, isAgentEnabled, token, organizationId, messages, handleChunk]);

  const handleDeclineToolCall = useCallback(async () => {
    const pending = pendingToolApproval;
    if (!pending || !isAgentEnabled) return;

    setPendingToolApproval(null);
    pendingApprovalRef.current = false;

    try {
      const headers = await getAuthHeaders(token ?? null, organizationId ?? null);
      await declineAgentToolCall({ runId: pending.runId, toolCallId: pending.toolCallId }, headers);
      const assistantMessageId = messages.filter((m) => m.role === 'assistant').pop()?.id;
      if (assistantMessageId) {
        setMessages((prev) =>
          prev.map((m) =>
            m.id === assistantMessageId
              ? { ...m, content: m.content + '\n\n_Tool call was declined._' }
              : m
          )
        );
      }
    } catch {
      // ignore
    } finally {
      setIsStreaming(false);
    }
  }, [pendingToolApproval, isAgentEnabled, token, organizationId, messages]);

  const lastAutoApprovedRef = useRef<string | null>(null);
  useEffect(() => {
    if (!pendingToolApproval) {
      lastAutoApprovedRef.current = null;
      return;
    }
    if (!autoRunTools) return;
    const key = `${pendingToolApproval.runId}-${pendingToolApproval.toolCallId}`;
    if (lastAutoApprovedRef.current === key) return;
    lastAutoApprovedRef.current = key;
    handleApproveToolCall();
  }, [pendingToolApproval, autoRunTools, handleApproveToolCall]);

  const handleSubmit = useCallback(
    (e?: React.FormEvent) => {
      e?.preventDefault();
      if (!inputValue.trim() || isStreaming || !threadId) return;

      const content = inputValue.trim();
      const userMessage: ChatMessage = {
        id: crypto.randomUUID(),
        role: 'user',
        content,
        timestamp: new Date(),
        ...(contexts.length > 0 ? { contexts: [...contexts] } : {})
      };

      setMessages((prev) => {
        const isFirst = prev.length === 0;
        if (isFirst && onFirstMessage) {
          onFirstMessage(content);
        }
        return [...prev, userMessage];
      });
      setInputValue('');
      setIsStreaming(true);
      streamResponse(content);
    },
    [inputValue, isStreaming, threadId, streamResponse, onFirstMessage, contexts]
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSubmit();
      }
    },
    [handleSubmit]
  );

  const handleSuggestionClick = useCallback(
    (text: string) => {
      if (!isStreaming && threadId) {
        setInputValue('');
        const userMessage: ChatMessage = {
          id: crypto.randomUUID(),
          role: 'user',
          content: text,
          timestamp: new Date(),
          ...(contexts.length > 0 ? { contexts: [...contexts] } : {})
        };
        setMessages((prev) => {
          const isFirst = prev.length === 0;
          if (isFirst && onFirstMessage) {
            onFirstMessage(text);
          }
          return [...prev, userMessage];
        });
        setIsStreaming(true);
        streamResponse(text);
      }
    },
    [isStreaming, threadId, streamResponse, onFirstMessage, contexts]
  );

  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputValue(e.target.value);
  }, []);

  const stopStreaming = useCallback(() => {
    abortRef.current?.abort();
    setIsStreaming(false);
  }, []);

  return {
    messages,
    inputValue,
    setInputValue,
    isStreaming,
    isLoadingHistory,
    isAgentConfigured: isAgentEnabled,
    pendingToolApproval,
    omStatus,
    scrollRef,
    textareaRef,
    handleSubmit,
    handleKeyDown,
    handleSuggestionClick,
    handleInputChange,
    handleApproveToolCall,
    handleDeclineToolCall,
    stopStreaming
  };
}
