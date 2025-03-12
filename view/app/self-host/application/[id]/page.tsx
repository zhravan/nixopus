"use client"
import React, { useEffect, useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import ApplicationLogs from '../../components/application-details/logs';
import Monitor from '../../components/application-details/monitor';
import Configuration from '../../components/application-details/configuration';
import DeploymentsList from '../../components/application-details/deploymentsList';
import { useWebSocket } from '@/hooks/socket_provider';
import { useParams } from 'next/navigation';
import { useGetApplicationByIdQuery } from '@/redux/services/deploy/applicationsApi';

function Page() {
  const { isReady, message, sendJsonMessage } = useWebSocket();
  const { id } = useParams()
  const { data: application, isLoading, error } = useGetApplicationByIdQuery({ id: id as string }, { skip: !id })
  const [crurrentPage,setCurrentPage] = useState(1)
  useEffect(() => {
    sendJsonMessage({
      action: 'subscribe',
      topic: 'monitor_application_deployment',
      data: {
        "resource_id": id
      }
    })
  }, [])

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl mt-10">
      <Tabs defaultValue="monitoring" className="w-full">
        <TabsList>
          <TabsTrigger value="monitoring">Monitoring</TabsTrigger>
          <TabsTrigger value="configuration">Configuration</TabsTrigger>
          <TabsTrigger value="deployments">Deployments</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
        </TabsList>
        <TabsContent value="deployments" className="mt-6">
          <DeploymentsList />
        </TabsContent>

        <TabsContent value="configuration" className="mt-6">
          <Configuration />
        </TabsContent>

        <TabsContent value="logs" className="mt-6">
          <ApplicationLogs logs={application?.logs} onRefresh={()=>{}} currentPage={crurrentPage} setCurrentPage={setCurrentPage} />
        </TabsContent>

        <TabsContent value="monitoring" className="mt-6">
          <Monitor />
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default Page;