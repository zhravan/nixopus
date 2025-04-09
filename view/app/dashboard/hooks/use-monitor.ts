'use client';
import { useWebSocket } from '@/hooks/socket-provider';
import { useTranslation } from '@/hooks/use-translation';
import { ContainerData, SystemStatsType } from '@/redux/types/monitor';
import { useEffect, useState, useRef } from 'react';
import { toast } from 'sonner';

function use_monitor() {
  const { sendJsonMessage, message, isReady } = useWebSocket();
  const [containersData, setContainersData] = useState<ContainerData[]>([]);
  const [systemStats, setSystemStats] = useState<SystemStatsType | null>(null);
  const [isMonitoring, setIsMonitoring] = useState(false);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const isInitializedRef = useRef(false);
  const { t } = useTranslation();

  useEffect(() => {
    if (message) {
      try {
        const parsedMessage =
          typeof message === 'string' && message.startsWith('{') ? JSON.parse(message) : message;

        if (parsedMessage.topic != 'dashboard_monitor') {
          return;
        }

        if (parsedMessage.action === 'get_containers' && parsedMessage.data) {
          setContainersData(parsedMessage.data);
        } else if (parsedMessage.action === 'get_system_stats' && parsedMessage.data) {
          setSystemStats(parsedMessage.data);
        } else if (parsedMessage.action === 'error') {
          if (isMonitoring) {
            stopMonitoring();
            setTimeout(startMonitoring, 5000);
          }
        }
      } catch (error) {
        toast.error(t('toasts.errors.realtimeMonitor'), {
          description: error instanceof Error ? error.message : 'Unknown error'
        });
      }
    }
  }, [message, t]);

  const startMonitoring = () => {
    if (!isMonitoring) {
      sendJsonMessage({
        action: 'dashboard_monitor',
        data: {
          interval: 10,
          operations: ['get_containers', 'get_system_stats']
        }
      });
      setIsMonitoring(true);
    }
  };

  const stopMonitoring = () => {
    if (isMonitoring) {
      sendJsonMessage({
        action: 'stop_dashboard_monitor'
      });
      setIsMonitoring(false);
    }
  };

  useEffect(() => {
    if (isReady && !isInitializedRef.current) {
      startMonitoring();
      isInitializedRef.current = true;
    }
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [isReady]);

  useEffect(() => {
    if (isReady && !isMonitoring) {
      reconnectTimeoutRef.current = setTimeout(() => {
        startMonitoring();
      }, 5000);
    }
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [isReady, isMonitoring]);

  return {
    containersData,
    systemStats,
    isMonitoring,
    startMonitoring,
    stopMonitoring
  };
}

export default use_monitor;
