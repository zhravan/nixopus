import { useWebSocket } from '@/hooks/socket-provider';
import { SOCKET_EVENTS } from '@/redux/api-conf';
import { useGetApplicationByIdQuery } from '@/redux/services/deploy/applicationsApi';
import { SubscribeToTopic } from '@/redux/sockets/socket';
import { useParams, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

interface WebSocketMessage {
  action: string;
  data: {
    action: string;
    application_id: string;
    data: {
      application_deployment_id: string;
      application_id: string;
      created_at: string;
      id: string;
      log: string;
      updated_at: string;
    };
    table: string;
  };
  topic: string;
}

function useApplicationDetails() {
  const { id } = useParams();
  const applicationId = id as string;
  const { data: application } = useGetApplicationByIdQuery(
    { id: applicationId },
    { skip: !applicationId }
  );
  const [currentPage, setCurrentPage] = useState(1);
  const [logs, setLogs] = useState(application?.logs || []);
  const { isReady, message, sendJsonMessage } = useWebSocket();
  const searchParams = useSearchParams();
  const defaultTab = searchParams.get('logs') === 'true' ? 'logs' : 'monitoring';

  useEffect(() => {
    if (applicationId) {
      sendJsonMessage(
        SubscribeToTopic(applicationId, SOCKET_EVENTS.MONITOR_APPLICATION_DEPLOYMENT)
      );
    }
  }, [applicationId]);

  useEffect(() => {
    if (message) {
      const parsedMessage: WebSocketMessage = JSON.parse(message);
      if (
        parsedMessage.action === 'message' &&
        parsedMessage.data.table === 'application_logs' &&
        parsedMessage.data.data.application_id === applicationId
      ) {
        setLogs((prevLogs) => [...prevLogs, parsedMessage.data.data]);
      }
    }
  }, [message, applicationId]);

  useEffect(() => {
    if (application?.logs) {
      setLogs(application.logs);
    }
  }, [application?.logs]);

  const parseEnvVariables = (variablesString: string | undefined): Record<string, string> => {
    if (!variablesString) return {};

    try {
      return variablesString
        .split(/\s+/)
        .filter((pair) => pair.includes('='))
        .reduce<Record<string, string>>((acc, curr) => {
          const [key, value] = curr.split('=');
          return key && value ? { ...acc, [key]: value } : acc;
        }, {});
    } catch (error) {
      console.error('Error parsing variables:', error);
      return {};
    }
  };

  const envVariables = application?.environment_variables
    ? parseEnvVariables(application.environment_variables)
    : {};

  const buildVariables = application?.build_variables
    ? parseEnvVariables(application.build_variables)
    : {};

  return {
    currentPage,
    setCurrentPage,
    application: application ? { ...application, logs } : undefined,
    envVariables,
    buildVariables,
    defaultTab
  };
}

export default useApplicationDetails;
