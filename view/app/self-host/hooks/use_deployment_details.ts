import { useWebSocket } from '@/hooks/socket-provider';
import { SOCKET_EVENTS } from '@/redux/api-conf';
import { useGetApplicationDeploymentByIdQuery } from '@/redux/services/deploy/applicationsApi';
import { SubscribeToTopic } from '@/redux/sockets/socket';
import { useParams } from 'next/navigation';
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

function useDeploymentDetails() {
  const { deployment_id } = useParams();
  const deploymentId = deployment_id?.toString();
  const { data: deployment } = useGetApplicationDeploymentByIdQuery(
    { id: deploymentId as string },
    { skip: !deploymentId }
  );
  const [logs, setLogs] = useState(deployment?.logs || []);
  const { isReady, message, sendJsonMessage } = useWebSocket();

  useEffect(() => {
    if (deploymentId) {
      sendJsonMessage(SubscribeToTopic(deploymentId, SOCKET_EVENTS.MONITOR_APPLICATION_DEPLOYMENT));
    }
  }, [deploymentId]);

  useEffect(() => {
    if (message) {
      const parsedMessage: WebSocketMessage = JSON.parse(message);
      if (
        parsedMessage.action === 'message' &&
        parsedMessage.data.table === 'application_logs' &&
        parsedMessage.data.data.application_deployment_id === deploymentId
      ) {
        setLogs((prevLogs) => [...prevLogs, parsedMessage.data.data]);
      }
    }
  }, [message, deploymentId]);

  useEffect(() => {
    if (deployment?.logs) {
      setLogs(deployment.logs);
    }
  }, [deployment?.logs]);

  return {
    deployment: deployment ? { ...deployment, logs } : undefined,
    logs
  };
}

export default useDeploymentDetails;
